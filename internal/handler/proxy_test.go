package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/otterscale/otterscale/internal/core"
)

// mockTunnelForProxy implements core.TunnelProvider for proxy tests.
type mockTunnelForProxy struct {
	address    string
	addressErr error
}

func (m *mockTunnelForProxy) CACertPEM() []byte { return nil }
func (m *mockTunnelForProxy) ListLinks() map[string]core.Link {
	return nil
}

func (m *mockTunnelForProxy) RegisterLink(context.Context, string, string, string, []byte) (core.TunnelGrant, error) {
	return core.TunnelGrant{}, nil
}

func (m *mockTunnelForProxy) ResolveAddress(_ context.Context, _ string) (string, error) {
	return m.address, m.addressErr
}

func TestProxyHandler_ForbiddenPath(t *testing.T) {
	handler := NewProxyHandler(&mockTunnelForProxy{address: "http://127.0.0.1:8080"})

	tests := []struct {
		name string
		path string
	}{
		{"admin path", "/api/v1/admin/tsdb/delete_series"},
		{"reload", "/-/reload"},
		{"quit", "/-/quit"},
		{"write", "/api/v1/write"},
		{"root", "/"},
		{"random", "/random"},

		// Traversal out of an allowed prefix. It has to be
		// percent-encoded to reach this handler at all: ServeMux
		// resolves literal ".." segments and answers with a redirect,
		// but "%2e%2e" survives that and is only decoded later, by
		// PathValue. Written this way the path satisfies the
		// "/api/v1/query" prefix while addressing an admin endpoint.
		{"traversal to reload", "/api/v1/query/%2e%2e/%2e%2e/%2e%2e/-/reload"},
		{"traversal to admin", "/api/v1/query/%2e%2e/admin/tsdb/delete_series"},
		{"traversal to write", "/api/v1/labels/%2e%2e/%2e%2e/v1/write"},
		{"encoded separator", "/api/v1/query/..%2f..%2f..%2f-/reload"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.Handle("/proxy/{cluster}/prometheus/{path...}", handler)

			req := httptest.NewRequestWithContext(t.Context(), "GET", "/proxy/prod/prometheus"+tt.path, http.NoBody)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusForbidden)
			}
		})
	}
}

func TestProxyHandler_ClusterNotFound(t *testing.T) {
	handler := NewProxyHandler(&mockTunnelForProxy{
		addressErr: &core.ErrClusterNotFound{Cluster: "missing"},
	})

	mux := http.NewServeMux()
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", handler)

	req := httptest.NewRequestWithContext(t.Context(), "GET", "/proxy/missing/prometheus/api/v1/query?query=up", http.NoBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestProxyHandler_AllowedPath_ForwardsToBackend(t *testing.T) {
	// Start a fake Prometheus backend.
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/__otterscale/proxy/api/v1/query" {
			t.Errorf("backend received path %q, want %q", r.URL.Path, "/__otterscale/proxy/api/v1/query")
		}
		if r.URL.RawQuery != "query=up" {
			t.Errorf("backend received query %q, want %q", r.URL.RawQuery, "query=up")
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("Authorization header should be stripped")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success"}`))
	}))
	defer backend.Close()

	handler := NewProxyHandler(&mockTunnelForProxy{address: backend.URL})

	mux := http.NewServeMux()
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", handler)

	req := httptest.NewRequestWithContext(t.Context(), "GET", "/proxy/prod/prometheus/api/v1/query?query=up", http.NoBody)
	req.Header.Set("Authorization", "Bearer some-oidc-token")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != `{"status":"success"}` {
		t.Errorf("body = %q, want %q", rec.Body.String(), `{"status":"success"}`)
	}
}

// TestProxyHandler_StripsCallerIdentityHeaders checks that nothing the
// caller sent can reach the backend as a credential or as an identity
// assertion. Nothing on this path reads impersonation headers today;
// the point is that they cannot become load-bearing by accident, since
// anything that trusted them would be trusting the caller.
func TestProxyHandler_StripsCallerIdentityHeaders(t *testing.T) {
	leaked := make(chan []string, 1)

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var got []string
		for name := range r.Header {
			if name == "Authorization" || strings.HasPrefix(name, "Impersonate-") {
				got = append(got, name)
			}
		}
		leaked <- got
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	handler := NewProxyHandler(&mockTunnelForProxy{address: backend.URL})

	mux := http.NewServeMux()
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", handler)

	req := httptest.NewRequestWithContext(t.Context(), "GET", "/proxy/prod/prometheus/api/v1/query?query=up", http.NoBody)
	req.Header.Set("Authorization", "Bearer some-oidc-token")
	req.Header.Set("Impersonate-User", "system:admin")
	req.Header.Add("Impersonate-Group", "system:masters")
	req.Header.Set("Impersonate-Uid", "0")
	req.Header.Set("Impersonate-Extra-Scopes", "everything")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := <-leaked; len(got) > 0 {
		t.Errorf("backend received caller-controlled headers %v, want none", got)
	}
}

// TestProxyHandler_ForwardsNormalizedPath pins that the path reaching
// the backend is the one the allowlist approved. Checking one path and
// forwarding another is how an allowlist gets bypassed.
func TestProxyHandler_ForwardsNormalizedPath(t *testing.T) {
	forwarded := make(chan string, 1)

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forwarded <- r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	handler := NewProxyHandler(&mockTunnelForProxy{address: backend.URL})

	mux := http.NewServeMux()
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", handler)

	// Resolves back inside the allowlist, so it is served — but it must
	// be forwarded in its resolved form, not as written. Encoded so
	// that ServeMux passes it through rather than redirecting; see
	// TestProxyHandler_ForbiddenPath.
	req := httptest.NewRequestWithContext(t.Context(), "GET",
		"/proxy/prod/prometheus/api/v1/status/%2e%2e/query?query=up", http.NoBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got, want := <-forwarded, "/__otterscale/proxy/api/v1/query"; got != want {
		t.Errorf("backend received path %q, want %q", got, want)
	}
}
