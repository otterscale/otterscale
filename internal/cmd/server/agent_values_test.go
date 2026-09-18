package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/handler"
	"github.com/otterscale/otterscale/internal/providers/cache"
)

const testJoinSecret = "root-secret"

// stubRenderer returns a recognizable file, to tell a served body from an
// error page.
type stubRenderer struct{}

func (stubRenderer) Render(values *core.AgentValues) (string, error) {
	return "cluster: " + values.Cluster + "\n", nil
}

// stubHarbor answers without contacting anything.
type stubHarbor struct{}

func (stubHarbor) EnsureRobotAccount(
	_ context.Context, cluster, secret string,
) (core.HarborRobotCredentials, error) {
	return core.HarborRobotCredentials{Name: "robot$" + cluster, Secret: secret}, nil
}

// newAgentValuesMux mounts just the agent values route, returning the use case
// so a test can mint a real id for it.
func newAgentValuesMux(t *testing.T) (*http.ServeMux, *core.AgentValuesUseCase) {
	t.Helper()

	join, err := core.NewJoinAuthority(testJoinSecret)
	if err != nil {
		t.Fatalf("NewJoinAuthority() error = %v", err)
	}
	useCase := core.NewAgentValuesUseCase(
		&core.AgentValuesConfig{
			ExternalURL:     "https://otterscale.example.com/api/",
			TunnelServerURL: "https://192.0.2.1:30300",
			HarborURL:       "https://harbor.example.com:8443",
		},
		core.Version("v1.5.0"), join, cache.NewAgentValuesStore(), stubRenderer{}, stubHarbor{},
	)

	h := &Handler{agentValues: handler.NewAgentValuesHandler(useCase)}

	mux := http.NewServeMux()
	mux.HandleFunc("GET "+agentValuesPath+"{id}", h.handleAgentValues)

	return mux, useCase
}

func issueTestTicket(t *testing.T, useCase *core.AgentValuesUseCase) string {
	t.Helper()

	result, err := useCase.Issue(t.Context(), &core.AgentValuesRequest{
		Cluster: "prod",
		Subject: "admin-1",
		ClusterInfo: &core.AgentClusterInfo{
			ExternalAddress: "192.0.2.1",
			NodePortRange:   "30000-32767",
		},
	})
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	_, id, found := strings.Cut(result.URL, agentValuesPath)
	if !found {
		t.Fatalf("URL %q carries no id", result.URL)
	}
	return id
}

// Also pins no-store: the body is a join token and a registry secret.
func TestHandleAgentValues_ServesYAML(t *testing.T) {
	mux, useCase := newAgentValuesMux(t)
	id := issueTestTicket(t, useCase)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, agentValuesPath+id, http.NoBody))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d (body %q)", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got, want := rec.Header().Get("Content-Type"), "text/yaml; charset=utf-8"; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
	if got, want := rec.Header().Get("Cache-Control"), "no-store"; got != want {
		t.Errorf("Cache-Control = %q, want %q", got, want)
	}
	if got, want := rec.Body.String(), "cluster: prod\n"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestHandleAgentValues_RejectsUnknownID(t *testing.T) {
	mux, useCase := newAgentValuesMux(t)
	valid := issueTestTicket(t, useCase)

	tests := []struct {
		name string
		id   string
	}{
		{name: "never issued", id: "AAAAAAAAAAAAAAAAAAAAAA"},
		{name: "truncated", id: valid[:len(valid)/2]},
		{name: "not base64", id: "!!!!!!!!!!!!!!!!!!!!!!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, agentValuesPath+tt.id, http.NoBody))

			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
			}
			if got, want := strings.TrimSpace(rec.Body.String()), "invalid or expired token"; got != want {
				t.Errorf("body = %q, want the uniform %q", got, want)
			}
			// Nor may it say whether the id was never issued or has expired.
			if strings.Contains(rec.Body.String(), "cluster:") {
				t.Error("the refusal leaked values")
			}
		})
	}
}

// The route must not answer its own prefix with an empty id.
func TestHandleAgentValues_NoID(t *testing.T) {
	mux, _ := newAgentValuesMux(t)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, agentValuesPath, http.NoBody))

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
