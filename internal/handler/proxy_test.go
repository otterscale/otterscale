package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/otterscale/otterscale/internal/core"
)

// mockTunnelForProxy implements core.TunnelProvider for these tests.
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

// signedIn is a request context as the OIDC middleware leaves it.
func signedIn(t *testing.T) context.Context {
	t.Helper()
	return core.WithUserInfo(t.Context(), core.UserInfo{Subject: "alice", Groups: []string{"oidc:ws-a"}})
}

// newProxyMux routes the three scopes the way internal/cmd/server does.
func newProxyMux(h *ProxyHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", h)
	mux.HandleFunc("/proxy/{cluster}/prometheus-tenancy/{path...}", h.ServeTenancy)
	mux.HandleFunc("/proxy/{cluster}/prometheus-shared/{name...}", h.ServeShared)
	return mux
}

// scopeQueries is one valid query per scope, for tests that hold for all of them.
var scopeQueries = map[string]string{
	"cluster scope":   "/proxy/prod/prometheus/api/v1/query?query=up",
	"namespace scope": "/proxy/prod/prometheus-tenancy/api/v1/query?namespace=ws-a&query=up",
	"shared scope":    "/proxy/prod/prometheus-shared/kubevirt-handlers?namespace=ws-a",
}

// assertPrometheusError requires Prometheus' JSON error shape, which api.md
// promises for every error at all three scopes.
func assertPrometheusError(t *testing.T, rec *httptest.ResponseRecorder, status int, errorType string) {
	t.Helper()
	if rec.Code != status {
		t.Errorf("status = %d, want %d", rec.Code, status)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var body struct {
		Status    string `json:"status"`
		ErrorType string `json:"errorType"`
		Error     string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil ||
		body.Status != "error" || body.ErrorType != errorType || body.Error == "" {
		t.Errorf("body = %q, want Prometheus' error shape with errorType %q", rec.Body, errorType)
	}
}

func TestProxyHandler_ForbiddenPath(t *testing.T) {
	handler := NewProxyHandler(&mockTunnelForProxy{address: "http://127.0.0.1:8080"}, clusterViewer())

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

			req := httptest.NewRequestWithContext(signedIn(t), "GET", "/proxy/prod/prometheus"+tt.path, http.NoBody)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			assertPrometheusError(t, rec, http.StatusForbidden, "forbidden")
		})
	}
}

func TestProxyHandler_ClusterNotFound(t *testing.T) {
	handler := NewProxyHandler(&mockTunnelForProxy{
		addressErr: &core.ErrClusterNotFound{Cluster: "missing"},
	}, clusterViewer())

	mux := http.NewServeMux()
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", handler)

	req := httptest.NewRequestWithContext(signedIn(t), "GET", "/proxy/missing/prometheus/api/v1/query?query=up", http.NoBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assertPrometheusError(t, rec, http.StatusNotFound, "not_found")
}

// TestProxyHandler_BadClusterAddress: an address the tunnel hands back that
// does not parse is the server's fault, and still answered as Prometheus would.
func TestProxyHandler_BadClusterAddress(t *testing.T) {
	for scope, target := range scopeQueries {
		t.Run(scope, func(t *testing.T) {
			h := NewProxyHandler(&mockTunnelForProxy{address: "http://[::1"}, clusterViewer())
			rec := httptest.NewRecorder()
			newProxyMux(h).ServeHTTP(rec, httptest.NewRequestWithContext(signedIn(t), http.MethodGet, target, http.NoBody))
			assertPrometheusError(t, rec, http.StatusInternalServerError, "internal")
		})
	}
}

// TestProxyHandler_NoAuthenticatedUser: 401 is for a request the hub itself
// has no user for, decided before any cluster is asked — whatever a reviewer
// would have answered.
func TestProxyHandler_NoAuthenticatedUser(t *testing.T) {
	contexts := map[string]func(*testing.T) context.Context{
		"no user": func(t *testing.T) context.Context { return t.Context() },
		"no subject": func(t *testing.T) context.Context {
			return core.WithUserInfo(t.Context(), core.UserInfo{Groups: []string{"oidc:admin"}})
		},
	}
	for scope, target := range scopeQueries {
		for name, ctx := range contexts {
			t.Run(scope+"/"+name, func(t *testing.T) {
				up, srv := newUpstream(t)
				access := &fakeAccess{allow: map[core.AccessAttributes]bool{
					core.ClusterMetricsAccess: true, core.NamespaceMetricsAccess("ws-a"): true,
				}}
				h := NewProxyHandler(&mockTunnelForProxy{address: srv.URL}, access)
				rec := httptest.NewRecorder()
				newProxyMux(h).ServeHTTP(rec, httptest.NewRequestWithContext(ctx(t), http.MethodGet, target, http.NoBody))

				assertPrometheusError(t, rec, http.StatusUnauthorized, "unauthorized")
				if n := len(up.all()); n != 0 {
					t.Errorf("upstream called %d times", n)
				}
				if got := access.questions(); len(got) != 0 {
					t.Errorf("asked %v, want no review without a user", got)
				}
			})
		}
	}
}

// TestProxyHandler_ReviewClusterNotFound: a cluster that is gone by the time
// it is asked is 404, as when the tunnel cannot resolve it — not 503.
func TestProxyHandler_ReviewClusterNotFound(t *testing.T) {
	for scope, target := range scopeQueries {
		t.Run(scope, func(t *testing.T) {
			up, srv := newUpstream(t)
			access := &fakeAccess{err: &core.ErrClusterNotFound{Cluster: "prod"}}
			h := NewProxyHandler(&mockTunnelForProxy{address: srv.URL}, access)
			rec := httptest.NewRecorder()
			newProxyMux(h).ServeHTTP(rec, httptest.NewRequestWithContext(signedIn(t), http.MethodGet, target, http.NoBody))

			assertPrometheusError(t, rec, http.StatusNotFound, "not_found")
			if n := len(up.all()); n != 0 {
				t.Errorf("upstream called %d times", n)
			}
		})
	}
}

// TestProxyHandler_ClusterSideUnauthenticated: a cluster that answers the
// review with 401 has rejected the agent's credentials, not the caller — the
// caller is signed in. That is a failed review like any other: 503.
func TestProxyHandler_ClusterSideUnauthenticated(t *testing.T) {
	for scope, target := range scopeQueries {
		t.Run(scope, func(t *testing.T) {
			up, srv := newUpstream(t)
			access := &fakeAccess{err: &core.DomainError{Code: core.ErrorCodeUnauthenticated, Message: "Unauthorized"}}
			h := NewProxyHandler(&mockTunnelForProxy{address: srv.URL}, access)
			rec := httptest.NewRecorder()
			newProxyMux(h).ServeHTTP(rec, httptest.NewRequestWithContext(signedIn(t), http.MethodGet, target, http.NoBody))

			assertPrometheusError(t, rec, http.StatusServiceUnavailable, "unavailable")
			if n := len(up.all()); n != 0 {
				t.Errorf("upstream called %d times", n)
			}
		})
	}
}

// TestProxyHandler_AgentUnreachable: a request that never reached the agent
// is 502 in Prometheus' error shape, with nothing of the underlying error —
// the tunnel address — in it.
func TestProxyHandler_AgentUnreachable(t *testing.T) {
	closed := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	closed.Close()
	for scope, target := range scopeQueries {
		t.Run(scope, func(t *testing.T) {
			access := &fakeAccess{allow: map[core.AccessAttributes]bool{
				core.ClusterMetricsAccess: true, core.NamespaceMetricsAccess("ws-a"): true,
			}}
			h := NewProxyHandler(&mockTunnelForProxy{address: closed.URL}, access)
			rec := httptest.NewRecorder()
			newProxyMux(h).ServeHTTP(rec, httptest.NewRequestWithContext(signedIn(t), http.MethodGet, target, http.NoBody))

			assertPrometheusError(t, rec, http.StatusBadGateway, "unavailable")
			if strings.Contains(rec.Body.String(), "127.0.0.1") {
				t.Errorf("body %q carries the tunnel address", rec.Body)
			}
		})
	}
}

func TestProxyHandler_AllowedPath_ForwardsToBackend(t *testing.T) {
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

	handler := NewProxyHandler(&mockTunnelForProxy{address: backend.URL}, clusterViewer())

	mux := http.NewServeMux()
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", handler)

	req := httptest.NewRequestWithContext(signedIn(t), "GET", "/proxy/prod/prometheus/api/v1/query?query=up", http.NoBody)
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
			if name == "Authorization" || name == "Cookie" || strings.HasPrefix(name, "Impersonate-") {
				got = append(got, name)
			}
		}
		leaked <- got
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	handler := NewProxyHandler(&mockTunnelForProxy{address: backend.URL}, clusterViewer())

	mux := http.NewServeMux()
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", handler)

	req := httptest.NewRequestWithContext(signedIn(t), "GET", "/proxy/prod/prometheus/api/v1/query?query=up", http.NoBody)
	req.Header.Set("Authorization", "Bearer some-oidc-token")
	req.Header.Set("Impersonate-User", "system:admin")
	req.Header.Add("Impersonate-Group", "system:masters")
	req.Header.Set("Impersonate-Uid", "0")
	req.Header.Set("Impersonate-Extra-Scopes", "everything")
	req.Header.Set("Cookie", "session=some-session")

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

	handler := NewProxyHandler(&mockTunnelForProxy{address: backend.URL}, clusterViewer())

	mux := http.NewServeMux()
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", handler)

	// Resolves back inside the allowlist, so it is served — but it must
	// be forwarded in its resolved form, not as written. Encoded so
	// that ServeMux passes it through rather than redirecting; see
	// TestProxyHandler_ForbiddenPath.
	req := httptest.NewRequestWithContext(signedIn(t), "GET",
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

func serveCluster(t *testing.T, access core.AccessReviewer, backend, rawPathAndQuery string) *httptest.ResponseRecorder {
	t.Helper()
	h := NewProxyHandler(&mockTunnelForProxy{address: backend}, access)
	mux := http.NewServeMux()
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", h)
	req := httptest.NewRequestWithContext(signedIn(t), http.MethodGet, "/proxy/prod/prometheus"+rawPathAndQuery, http.NoBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

// TestProxyHandler_ClusterScopeRefusesRegularUsers: without cluster metrics
// access every cluster-scope request is 403, whatever the PromQL — aggregates
// leak other namespaces as surely as series do.
func TestProxyHandler_ClusterScopeRefusesRegularUsers(t *testing.T) {
	requests := []string{
		"/api/v1/query?" + url.Values{"query": {`sum(ceph_cluster_total_bytes)`}}.Encode(),
		"/api/v1/query?" + url.Values{"query": {`count(kube_pod_info{namespace="victim"})`}}.Encode(),
		"/api/v1/query?" + url.Values{"query": {`sum(kube_pod_info) - sum(kube_pod_info{namespace!="victim"})`}}.Encode(),
		"/api/v1/query?" + url.Values{"query": {`count(kube_deployment_created{deployment="secret-app"})`}}.Encode(),
		"/api/v1/query?" + url.Values{"query": {`kube_pod_info`}}.Encode(),
		"/api/v1/query?" + url.Values{"query": {`kube_pod_info{namespace="ws-a"}`}}.Encode(), // even their own
		"/api/v1/query_range?" + url.Values{"query": {`up`}, "start": {"0"}, "end": {"60"}, "step": {"15"}}.Encode(),
		"/api/v1/series?" + url.Values{"match[]": {`up`}}.Encode(),
		"/api/v1/labels",
		"/api/v1/targets",
	}
	for _, rq := range requests {
		t.Run(rq, func(t *testing.T) {
			up, srv := newUpstream(t)
			access := tenantOf("ws-a")
			rec := serveCluster(t, access, srv.URL, rq)
			if rec.Code != http.StatusForbidden {
				t.Errorf("status = %d, want 403", rec.Code)
			}
			if n := len(up.all()); n != 0 {
				t.Errorf("upstream called %d times", n)
			}
			if got := access.questions(); len(got) != 1 || got[0] != core.ClusterMetricsAccess {
				t.Errorf("asked %v, want only the cluster question", got)
			}
		})
	}
}

// TestProxyHandler_ClusterScopeUnchangedForViewers: with cluster metrics access
// the request goes through exactly as sent.
func TestProxyHandler_ClusterScopeUnchangedForViewers(t *testing.T) {
	up, srv := newUpstream(t)
	raw := "query=sum%28ceph_cluster_total_bytes%29+or+vector%280%29&time=1700000000"
	if rec := serveCluster(t, clusterViewer(), srv.URL, "/api/v1/query?"+raw); rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	got := up.all()
	if len(got) != 1 || got[0].query.Encode() != mustParse(t, raw).Encode() {
		t.Errorf("upstream got %+v, want the query unchanged", got)
	}
}

func TestProxyHandler_ClusterScopeReviewFailure(t *testing.T) {
	up, srv := newUpstream(t)
	rec := serveCluster(t, &fakeAccess{err: errors.New("apiserver down")}, srv.URL, "/api/v1/query?query=up")
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", rec.Code)
	}
	if n := len(up.all()); n != 0 {
		t.Errorf("upstream called %d times: a failed review must not pass", n)
	}
}

func mustParse(t *testing.T, raw string) url.Values {
	t.Helper()
	v, err := url.ParseQuery(raw)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
