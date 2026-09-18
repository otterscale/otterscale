package handler

import (
	"context"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"

	pb "github.com/otterscale/otterscale/api/link/v1"

	"github.com/otterscale/otterscale/internal/core"
)

// stubTunnelProvider satisfies core.TunnelProvider without doing anything.
// Issuing agent values never reaches the tunnel, so no call here should happen.
type stubTunnelProvider struct{}

func (stubTunnelProvider) CACertPEM() []byte { return nil }

func (stubTunnelProvider) ListLinks() map[string]core.Link { return map[string]core.Link{} }

func (stubTunnelProvider) ResolveAddress(context.Context, string) (string, error) {
	return "", nil
}

func (stubTunnelProvider) RegisterLink(
	context.Context, string, string, string, []byte,
) (core.TunnelGrant, error) {
	return core.TunnelGrant{}, nil
}

const testJoinSecret = "test-root-secret"

// The rendered shape is tested where it is built; these only have to answer.
type stubRenderer struct{}

func (stubRenderer) Render(*core.AgentValues) (string, error) {
	return "rendered: true\n", nil
}

// stubTicketStore is an AgentValuesStore that keeps tickets in a map. The real
// one lives in the providers layer.
type stubTicketStore struct {
	tickets map[string]*core.AgentValuesRequest
}

func newStubTicketStore() *stubTicketStore {
	return &stubTicketStore{tickets: make(map[string]*core.AgentValuesRequest)}
}

func (s *stubTicketStore) Put(
	_ context.Context, id string, req *core.AgentValuesRequest, _ time.Time,
) error {
	s.tickets[id] = req
	return nil
}

func (s *stubTicketStore) Get(_ context.Context, id string) (*core.AgentValuesRequest, error) {
	req, ok := s.tickets[id]
	if !ok {
		return nil, core.ErrAgentValuesTicketNotFound
	}
	return req, nil
}

type stubHarbor struct{}

func (stubHarbor) EnsureRobotAccount(
	_ context.Context, cluster, secret string,
) (core.HarborRobotCredentials, error) {
	return core.HarborRobotCredentials{Name: "robot$" + cluster, Secret: secret}, nil
}

// testAgentValuesConfig is a fully configured deployment, so a test that wants
// a precondition failure has to remove something explicitly.
func testAgentValuesConfig() *core.AgentValuesConfig {
	return &core.AgentValuesConfig{
		ExternalURL:     "https://otterscale.example.com/api/",
		TunnelServerURL: "https://192.0.2.1:30300",
		HarborURL:       "https://harbor.example.com:8443",
		TrustedCASecret: core.TrustedCASecretName,
		TrustedCAKey:    core.DefaultTrustedCAKey,
	}
}

func newTestLinkService(t *testing.T) *LinkService {
	t.Helper()
	return newTestLinkServiceWith(t, testAgentValuesConfig())
}

func newTestLinkServiceWith(t *testing.T, cfg *core.AgentValuesConfig) *LinkService {
	t.Helper()

	join, err := core.NewJoinAuthority(testJoinSecret)
	if err != nil {
		t.Fatalf("NewJoinAuthority() error = %v", err)
	}
	version := core.Version("v1.0.0")
	return NewLinkService(
		core.NewLinkUseCase(stubTunnelProvider{}, version, join),
		core.NewAgentValuesUseCase(cfg, join, newStubTicketStore(), stubRenderer{}, stubHarbor{}),
	)
}

func issueAgentValuesRequest(cluster string, info *pb.AgentClusterInfo) *pb.IssueAgentValuesRequest {
	req := &pb.IssueAgentValuesRequest{}
	req.SetCluster(cluster)
	if info != nil {
		req.SetClusterInfo(info)
	}
	return req
}

func testClusterInfo() *pb.AgentClusterInfo {
	info := &pb.AgentClusterInfo{}
	info.SetExternalAddress("192.0.2.1")
	info.SetNodePortRange("30000-32767")
	return info
}

func adminContext(t *testing.T) context.Context {
	t.Helper()
	return core.WithUserInfo(t.Context(), core.UserInfo{
		Subject: "admin-1",
		Groups:  []string{"system:authenticated", "oidc:admin"},
	})
}

// TestLinkService_IssueAgentValues_RequiresAdmin is the regression test for the
// gate on this procedure: the file it returns embeds a join token, which claims
// a cluster, and binds the caller to cluster-admin on it. An ordinary
// authenticated user must not be able to obtain one.
func TestLinkService_IssueAgentValues_RequiresAdmin(t *testing.T) {
	tests := []struct {
		name string
		ctx  func(t *testing.T) context.Context
		want connect.Code
	}{
		{
			name: "no authenticated caller",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				return t.Context()
			},
			want: connect.CodeUnauthenticated,
		},
		{
			name: "authenticated but not an admin",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				return core.WithUserInfo(t.Context(), core.UserInfo{
					Subject: "user-1",
					Groups:  []string{"system:authenticated", "oidc:developer"},
				})
			},
			want: connect.CodePermissionDenied,
		},
		{
			// "admin" unprefixed is what a Kubernetes-native group looks like.
			// The OIDC middleware prefixes every claim it forwards, so matching
			// an unprefixed name would honor a group it never issued.
			name: "admin without the oidc: prefix",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				return core.WithUserInfo(t.Context(), core.UserInfo{
					Subject: "user-1",
					Groups:  []string{"system:authenticated", "admin"},
				})
			},
			want: connect.CodePermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestLinkService(t)

			resp, err := s.IssueAgentValues(tt.ctx(t), issueAgentValuesRequest("my-cluster", testClusterInfo()))
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if got := connect.CodeOf(err); got != tt.want {
				t.Errorf("code = %v, want %v", got, tt.want)
			}
			if resp != nil {
				t.Error("a refused request must return no response")
			}
		})
	}
}

func TestLinkService_IssueAgentValues_AllowsAdmin(t *testing.T) {
	s := newTestLinkService(t)

	resp, err := s.IssueAgentValues(adminContext(t), issueAgentValuesRequest("my-cluster", testClusterInfo()))
	if err != nil {
		t.Fatalf("IssueAgentValues() error = %v", err)
	}
	if resp.GetValues() == "" {
		t.Error("values are empty")
	}
	if want := "https://otterscale.example.com/api/link/values/"; !strings.HasPrefix(resp.GetUrl(), want) {
		t.Errorf("url = %q, want the prefix %q", resp.GetUrl(), want)
	}
	if got := resp.GetUrlExpiresAt().AsTime(); !got.After(time.Now()) {
		t.Errorf("url_expires_at = %v, want a time in the future", got)
	}
}

// TestLinkService_IssueAgentValues_RejectsBadRequests covers what the chart
// would otherwise only reject at install time, on another cluster.
func TestLinkService_IssueAgentValues_RejectsBadRequests(t *testing.T) {
	tests := []struct {
		name    string
		cluster string
		info    *pb.AgentClusterInfo
		want    connect.Code
	}{
		{
			// Absent rather than defaulted: the chart's own defaults here are
			// example addresses, and defaulting would publish those.
			name:    "cluster info omitted",
			cluster: "my-cluster",
			info:    nil,
			want:    connect.CodeInvalidArgument,
		},
		{
			name:    "cluster name is not a label value",
			cluster: "My Cluster",
			info:    testClusterInfo(),
			want:    connect.CodeInvalidArgument,
		},
		{
			name:    "external address carries a scheme",
			cluster: "my-cluster",
			info: func() *pb.AgentClusterInfo {
				info := testClusterInfo()
				info.SetExternalAddress("https://192.0.2.1")
				return info
			}(),
			want: connect.CodeInvalidArgument,
		},
		{
			name:    "external address carries a port",
			cluster: "my-cluster",
			info: func() *pb.AgentClusterInfo {
				info := testClusterInfo()
				info.SetExternalAddress("192.0.2.1:8443")
				return info
			}(),
			want: connect.CodeInvalidArgument,
		},
		{
			name:    "node port range is inverted",
			cluster: "my-cluster",
			info: func() *pb.AgentClusterInfo {
				info := testClusterInfo()
				info.SetNodePortRange("32767-30000")
				return info
			}(),
			want: connect.CodeInvalidArgument,
		},
		{
			name:    "inference url is not absolute",
			cluster: "my-cluster",
			info: func() *pb.AgentClusterInfo {
				info := testClusterInfo()
				info.SetInferenceUrl("inference.example.com")
				return info
			}(),
			want: connect.CodeInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestLinkService(t)

			resp, err := s.IssueAgentValues(adminContext(t), issueAgentValuesRequest(tt.cluster, tt.info))
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if got := connect.CodeOf(err); got != tt.want {
				t.Errorf("code = %v, want %v", got, tt.want)
			}
			if resp != nil {
				t.Error("a refused request must return no response")
			}
		})
	}
}

// TestLinkService_IssueAgentValues_ReportsMissingConfig is why the providers
// do not fail at startup: an unconfigured setting has to surface here, on the
// one procedure that needs it, rather than stop the server from serving.
func TestLinkService_IssueAgentValues_ReportsMissingConfig(t *testing.T) {
	tests := []struct {
		name  string
		clear func(cfg *core.AgentValuesConfig)
	}{
		{
			name:  "no external url",
			clear: func(cfg *core.AgentValuesConfig) { cfg.ExternalURL = "" },
		},
		{
			name:  "no external tunnel url",
			clear: func(cfg *core.AgentValuesConfig) { cfg.TunnelServerURL = "" },
		},
		{
			name:  "no harbor url",
			clear: func(cfg *core.AgentValuesConfig) { cfg.HarborURL = "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testAgentValuesConfig()
			tt.clear(cfg)
			s := newTestLinkServiceWith(t, cfg)

			_, err := s.IssueAgentValues(adminContext(t), issueAgentValuesRequest("my-cluster", testClusterInfo()))
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if got := connect.CodeOf(err); got != connect.CodeFailedPrecondition {
				t.Errorf("code = %v, want %v", got, connect.CodeFailedPrecondition)
			}
		})
	}
}
