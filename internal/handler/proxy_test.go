package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
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

func (m *mockTunnelForProxy) RegisterLink(context.Context, string, string, string, []byte) (addr string, cert []byte, err error) {
	return "", nil, nil
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
