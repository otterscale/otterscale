package handler

import (
	"context"
	"errors"
	"strings"
	"testing"

	"connectrpc.com/connect"
	linkv1 "github.com/otterscale/api/link/v1"

	"github.com/otterscale/otterscale/internal/core"
)

const testRancherProjectID = "local:p-test"

type handlerTunnel struct {
	rancherProjectID string
}

func (*handlerTunnel) CACertPEM() []byte               { return []byte("ca") }
func (*handlerTunnel) ListLinks() map[string]core.Link { return nil }
func (t *handlerTunnel) RegisterLink(_ context.Context, _, _, _, rancherProjectID string, _ []byte) (endpoint string, certificate []byte, err error) {
	t.rancherProjectID = rancherProjectID
	return "127.0.0.1:16598", []byte("cert"), nil
}
func (*handlerTunnel) ResolveAddress(context.Context, string) (string, error) { return "", nil }

type handlerProjectStore struct {
	projects []core.RancherProject
	has      bool
	err      error
}

func (s *handlerProjectStore) ListProjects(context.Context) ([]core.RancherProject, error) {
	return s.projects, s.err
}

func (s *handlerProjectStore) HasProject(context.Context, string) (bool, error) {
	return s.has, s.err
}

type capturingRenderer struct {
	params *core.ManifestParams
}

func (r *capturingRenderer) RenderAgentManifest(params *core.ManifestParams) (string, error) {
	copiedParams := *params
	r.params = &copiedParams
	return "manifest:" + params.RancherProjectID, nil
}

func newLinkServiceForTest(t *testing.T, store *handlerProjectStore) (*LinkService, *handlerTunnel, *capturingRenderer) {
	t.Helper()
	tunnel := &handlerTunnel{}
	renderer := &capturingRenderer{}
	useCase, err := core.NewLinkUseCase(tunnel, "v1.0.0", core.AgentManifestConfig{
		ServerURL: "https://server.example.com",
		TunnelURL: "https://tunnel.example.com",
		HMACKey:   []byte("test-hmac-key-must-be-32-bytes!!"),
	}, renderer, nil, store)
	if err != nil {
		t.Fatalf("NewLinkUseCase: %v", err)
	}
	return NewLinkService(useCase), tunnel, renderer
}

func adminContext() context.Context {
	return core.WithUserInfo(context.Background(), core.UserInfo{
		Subject: "admin@example.com",
		Groups:  []string{"oidc:admin"},
	})
}

func TestLinkServiceListRancherProjectsAuthorizationAndCache(t *testing.T) {
	store := &handlerProjectStore{projects: []core.RancherProject{{ID: testRancherProjectID, DisplayName: "Test"}}}
	service, _, _ := newLinkServiceForTest(t, store)
	req := &linkv1.ListRancherProjectsRequest{}

	_, err := service.ListRancherProjects(context.Background(), req)
	assertConnectCode(t, err, connect.CodeUnauthenticated)

	nonAdmin := core.WithUserInfo(context.Background(), core.UserInfo{Groups: []string{"oidc:user"}})
	_, err = service.ListRancherProjects(nonAdmin, req)
	assertConnectCode(t, err, connect.CodePermissionDenied)

	resp, err := service.ListRancherProjects(adminContext(), req)
	if err != nil {
		t.Fatalf("ListRancherProjects: %v", err)
	}
	if len(resp.GetProjects()) != 1 || resp.GetProjects()[0].GetId() != testRancherProjectID {
		t.Fatalf("projects = %#v", resp.GetProjects())
	}

	store.err = core.ErrRancherProjectCacheNotReady
	_, err = service.ListRancherProjects(adminContext(), req)
	assertConnectCode(t, err, connect.CodeUnavailable)
}

func TestLinkServiceGetAgentManifestRancherProjectValidation(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		has      bool
		storeErr error
		wantCode connect.Code
	}{
		{name: "empty bypasses unavailable cache", storeErr: core.ErrRancherProjectCacheNotReady},
		{name: "cache hit", id: testRancherProjectID, has: true},
		{name: "cache miss", id: "local:p-missing", wantCode: connect.CodeInvalidArgument},
		{name: "cache not ready", id: testRancherProjectID, storeErr: core.ErrRancherProjectCacheNotReady, wantCode: connect.CodeUnavailable},
		{name: "malformed before cache", id: "bad", storeErr: core.ErrRancherProjectCacheNotReady, wantCode: connect.CodeInvalidArgument},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &handlerProjectStore{has: tt.has, err: tt.storeErr}
			service, _, renderer := newLinkServiceForTest(t, store)
			req := &linkv1.GetAgentManifestRequest{}
			req.SetCluster("cluster")
			if tt.id != "" {
				req.SetRancherProjectId(tt.id)
			}
			resp, err := service.GetAgentManifest(adminContext(), req)
			if tt.wantCode != 0 {
				assertConnectCode(t, err, tt.wantCode)
				if renderer.params != nil {
					t.Fatal("renderer must not run after Project validation failure")
				}
				return
			}
			if err != nil {
				t.Fatalf("GetAgentManifest: %v", err)
			}
			if renderer.params.RancherProjectID != tt.id || resp.GetManifest() != "manifest:"+tt.id {
				t.Fatalf("manifest Project ID = %q, response = %q", renderer.params.RancherProjectID, resp.GetManifest())
			}

			token := strings.TrimPrefix(resp.GetUrl(), "https://server.example.com/link/manifest/")
			claims, err := service.link.VerifyManifestToken(t.Context(), token)
			if err != nil || claims.RancherProjectID != tt.id {
				t.Fatalf("token claims = %#v, %v", claims, err)
			}
			store.err = core.ErrRancherProjectCacheNotReady
			renderer.params = nil
			rawManifest, err := NewManifestHandler(service.link).RenderManifest(t.Context(), &claims)
			if err != nil || rawManifest != "manifest:"+tt.id || renderer.params == nil || renderer.params.RancherProjectID != tt.id {
				t.Fatalf("raw manifest snapshot = %q, %#v, %v", rawManifest, renderer.params, err)
			}
		})
	}
}

func TestLinkServiceRegisterRancherProjectRoundTrip(t *testing.T) {
	service, tunnel, _ := newLinkServiceForTest(t, &handlerProjectStore{})
	req := &linkv1.RegisterRequest{}
	req.SetCluster("cluster")
	req.SetAgentId("agent")
	req.SetCsr([]byte("csr"))
	req.SetRancherProjectId(testRancherProjectID)
	if _, err := service.Register(t.Context(), req); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if tunnel.rancherProjectID != testRancherProjectID {
		t.Fatalf("registered Project ID = %q", tunnel.rancherProjectID)
	}

	req.SetRancherProjectId("bad")
	_, err := service.Register(t.Context(), req)
	assertConnectCode(t, err, connect.CodeInvalidArgument)

	req.ClearRancherProjectId()
	if _, err := service.Register(t.Context(), req); err != nil {
		t.Fatalf("legacy RegisterRequest: %v", err)
	}
	if tunnel.rancherProjectID != "" {
		t.Fatalf("legacy registered Project ID = %q", tunnel.rancherProjectID)
	}

	link := toProtoLink("cluster", core.Link{RancherProjectID: testRancherProjectID})
	if link.GetRancherProjectId() != testRancherProjectID {
		t.Fatalf("ListLinks Project ID = %q", link.GetRancherProjectId())
	}
}

func assertConnectCode(t *testing.T, err error, want connect.Code) {
	t.Helper()
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) || connectErr.Code() != want {
		t.Fatalf("error = %v, want Connect code %v", err, want)
	}
}
