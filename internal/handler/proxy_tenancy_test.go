package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/otterscale/otterscale/internal/core"
)

// fakeAccess answers per attribute and records every question it was asked.
type fakeAccess struct {
	mu    sync.Mutex
	allow map[core.AccessAttributes]bool
	err   error
	asked []core.AccessAttributes
}

func (f *fakeAccess) Allowed(_ context.Context, _ string, attrs core.AccessAttributes) (bool, error) { //nolint:gocritic // hugeParam: signature fixed by core.AccessReviewer
	f.mu.Lock()
	defer f.mu.Unlock()
	f.asked = append(f.asked, attrs)
	if f.err != nil {
		return false, f.err
	}
	return f.allow[attrs], nil
}

func (f *fakeAccess) questions() []core.AccessAttributes {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]core.AccessAttributes(nil), f.asked...)
}

// tenantOf may read ns's metrics but not the cluster's.
func tenantOf(ns string) *fakeAccess {
	return &fakeAccess{allow: map[core.AccessAttributes]bool{core.NamespaceMetricsAccess(ns): true}}
}

// clusterViewer has cluster metrics access.
func clusterViewer() *fakeAccess {
	return &fakeAccess{allow: map[core.AccessAttributes]bool{core.ClusterMetricsAccess: true}}
}

// upstream records what reached the agent side of the tunnel.
type upstream struct {
	mu       sync.Mutex
	requests []upstreamRequest
}

type upstreamRequest struct {
	method string
	path   string
	query  url.Values
	form   url.Values
	header http.Header
}

func newUpstream(t *testing.T) (*upstream, *httptest.Server) {
	t.Helper()
	u := &upstream{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		form, _ := url.ParseQuery(string(body))
		u.mu.Lock()
		u.requests = append(u.requests, upstreamRequest{
			method: r.Method, path: r.URL.Path, query: r.URL.Query(), form: form, header: r.Header.Clone(),
		})
		u.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[]}}`))
	}))
	t.Cleanup(srv.Close)
	return u, srv
}

func (u *upstream) all() []upstreamRequest {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]upstreamRequest(nil), u.requests...)
}

func serveTenancy(t *testing.T, access core.AccessReviewer, backend string, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	h := NewProxyHandler(&mockTunnelForProxy{address: backend}, access)
	mux := http.NewServeMux()
	mux.HandleFunc("/proxy/{cluster}/prometheus-tenancy/{path...}", h.ServeTenancy)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func tenancyGet(t *testing.T, path string, params url.Values) *http.Request {
	t.Helper()
	return httptest.NewRequestWithContext(signedIn(t), http.MethodGet,
		"/proxy/prod/prometheus-tenancy"+path+"?"+params.Encode(), http.NoBody)
}

func tenancyPost(t *testing.T, path string, params, form url.Values) *http.Request {
	t.Helper()
	req := httptest.NewRequestWithContext(signedIn(t), http.MethodPost,
		"/proxy/prod/prometheus-tenancy"+path+"?"+params.Encode(), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

// assertPinned parses what reached upstream and requires every selector in it
// to carry exactly one namespace matcher, namespace="ns".
func assertPinned(t *testing.T, query, ns string) {
	t.Helper()
	expr, err := parser.NewParser(parser.Options{}).ParseExpr(query)
	if err != nil {
		t.Fatalf("upstream query %q does not parse: %v", query, err)
	}
	selectors := 0
	parser.Inspect(expr, func(node parser.Node, _ []parser.Node) error {
		vs, ok := node.(*parser.VectorSelector)
		if !ok {
			return nil
		}
		selectors++
		assertMatchersPinned(t, vs.String(), vs.LabelMatchers, ns)
		return nil
	})
	if selectors == 0 {
		t.Errorf("upstream query %q has no selector to check", query)
	}
}

func assertMatchersPinned(t *testing.T, what string, matchers []*labels.Matcher, ns string) {
	t.Helper()
	var got []*labels.Matcher
	for _, m := range matchers {
		if m.Name == "namespace" {
			got = append(got, m)
		}
	}
	if len(got) != 1 || got[0].Type != labels.MatchEqual || got[0].Value != ns {
		t.Errorf("%s: namespace matchers %v, want exactly namespace=%q", what, got, ns)
	}
}

func assertSelectorPinned(t *testing.T, selector, ns string) {
	t.Helper()
	ms, err := parser.NewParser(parser.Options{}).ParseMetricSelector(selector)
	if err != nil {
		t.Fatalf("upstream match[] %q does not parse: %v", selector, err)
	}
	assertMatchersPinned(t, selector, ms, ns)
}

func TestTenancy_NamespaceParameter(t *testing.T) {
	tests := []struct {
		name   string
		params url.Values
	}{
		{"missing", url.Values{"query": {"up"}}},
		{"two", url.Values{"query": {"up"}, "namespace": {"ws-a", "ws-b"}}},
		{"empty", url.Values{"query": {"up"}, "namespace": {""}}},
		{"upper case", url.Values{"query": {"up"}, "namespace": {"WS-A"}}},
		{"dot", url.Values{"query": {"up"}, "namespace": {"ws.a"}}},
		{"regex", url.Values{"query": {"up"}, "namespace": {"ws-.*"}}},
		{"too long", url.Values{"query": {"up"}, "namespace": {strings.Repeat("a", 64)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			up, srv := newUpstream(t)
			rec := serveTenancy(t, clusterViewer(), srv.URL, tenancyGet(t, "/api/v1/query", tt.params))
			if rec.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", rec.Code)
			}
			if n := len(up.all()); n != 0 {
				t.Errorf("upstream called %d times", n)
			}
		})
	}

	t.Run("also in the POST body", func(t *testing.T) {
		up, srv := newUpstream(t)
		req := tenancyPost(t, "/api/v1/query", url.Values{"namespace": {"ws-a"}},
			url.Values{"query": {"up"}, "namespace": {"ws-a"}})
		if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, req); rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", rec.Code)
		}
		if n := len(up.all()); n != 0 {
			t.Errorf("upstream called %d times", n)
		}
	})
}

func TestTenancy_Access(t *testing.T) {
	params := url.Values{"query": {"kube_pod_info"}, "namespace": {"ws-a"}}

	t.Run("tenant of that namespace", func(t *testing.T) {
		_, srv := newUpstream(t)
		access := tenantOf("ws-a")
		if rec := serveTenancy(t, access, srv.URL, tenancyGet(t, "/api/v1/query", params)); rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		want := []core.AccessAttributes{core.ClusterMetricsAccess, core.NamespaceMetricsAccess("ws-a")}
		if got := access.questions(); len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
			t.Errorf("asked %v, want %v", got, want)
		}
	})

	t.Run("cluster viewer skips the namespace question", func(t *testing.T) {
		_, srv := newUpstream(t)
		access := clusterViewer()
		if rec := serveTenancy(t, access, srv.URL, tenancyGet(t, "/api/v1/query", params)); rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if got := access.questions(); len(got) != 1 || got[0] != core.ClusterMetricsAccess {
			t.Errorf("asked %v, want only the cluster question", got)
		}
	})

	t.Run("tenant of another namespace", func(t *testing.T) {
		up, srv := newUpstream(t)
		if rec := serveTenancy(t, tenantOf("ws-b"), srv.URL, tenancyGet(t, "/api/v1/query", params)); rec.Code != http.StatusForbidden {
			t.Errorf("status = %d, want 403", rec.Code)
		}
		if n := len(up.all()); n != 0 {
			t.Errorf("upstream called %d times", n)
		}
	})

	t.Run("review failed", func(t *testing.T) {
		up, srv := newUpstream(t)
		access := &fakeAccess{err: errors.New("apiserver down")}
		if rec := serveTenancy(t, access, srv.URL, tenancyGet(t, "/api/v1/query", params)); rec.Code != http.StatusServiceUnavailable {
			t.Errorf("status = %d, want 503", rec.Code)
		}
		if n := len(up.all()); n != 0 {
			t.Errorf("upstream called %d times: a failed review must not pass", n)
		}
	})
}

func TestTenancy_RejectsConflictingNamespace(t *testing.T) {
	queries := []string{
		`kube_pod_info{namespace="other"}`,
		`kube_pod_info{namespace!="ws-a"}`,
		`kube_pod_info{namespace=~"other-.*"}`,
		`kube_pod_info{namespace!~"ws-.*"}`,
		`sum(kube_pod_info) - sum(kube_pod_info{namespace="other"})`,
		`sum(`, // not PromQL
	}
	for _, q := range queries {
		t.Run(q, func(t *testing.T) {
			for _, endpoint := range []string{"/api/v1/query", "/api/v1/query_range", "/api/v1/query_exemplars"} {
				up, srv := newUpstream(t)
				params := url.Values{"query": {q}, "namespace": {"ws-a"}}
				if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, tenancyGet(t, endpoint, params)); rec.Code != http.StatusBadRequest {
					t.Errorf("%s: status = %d, want 400", endpoint, rec.Code)
				}
				if n := len(up.all()); n != 0 {
					t.Errorf("%s: upstream called %d times", endpoint, n)
				}
			}
		})
	}

	selectors := []string{`up{namespace="other"}`, `{namespace!="ws-a"}`, `{__name__="up",namespace=~"x.*"}`, `up{`}
	for _, s := range selectors {
		t.Run("match[] "+s, func(t *testing.T) {
			for _, endpoint := range []string{"/api/v1/series", "/api/v1/labels", "/api/v1/label/pod/values"} {
				up, srv := newUpstream(t)
				params := url.Values{"match[]": {s}, "namespace": {"ws-a"}}
				if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, tenancyGet(t, endpoint, params)); rec.Code != http.StatusBadRequest {
					t.Errorf("%s: status = %d, want 400", endpoint, rec.Code)
				}
				if n := len(up.all()); n != 0 {
					t.Errorf("%s: upstream called %d times", endpoint, n)
				}
			}
		})
	}
}

func TestTenancy_PinsNamespace(t *testing.T) {
	queries := []string{
		`kube_pod_info`,
		`kube_pod_info{namespace="ws-a"}`,
		`kube_pod_info{namespace=~"ws-.*"}`,
		`kube_pod_info{namespace!="other"}`,
		`sum(rate(container_cpu_usage_seconds_total{container!=""}[5m])) or vector(0)`,
		`count(kube_pod_info) / count(kube_pod_info{namespace="ws-a"})`,
	}
	for _, q := range queries {
		t.Run(q, func(t *testing.T) {
			for _, endpoint := range []string{"/api/v1/query", "/api/v1/query_range", "/api/v1/query_exemplars"} {
				up, srv := newUpstream(t)
				params := url.Values{"query": {q}, "namespace": {"ws-a"}, "time": {"1700000000"}}
				if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, tenancyGet(t, endpoint, params)); rec.Code != http.StatusOK {
					t.Fatalf("%s: status = %d, want 200 (%s)", endpoint, rec.Code, rec.Body)
				}
				got := up.all()
				if len(got) != 1 {
					t.Fatalf("%s: upstream called %d times, want 1", endpoint, len(got))
				}
				assertPinned(t, got[0].query.Get("query"), "ws-a")
				if got[0].query.Get("time") != "1700000000" {
					t.Errorf("%s: other parameters lost: %v", endpoint, got[0].query)
				}
			}

			// The same through a POST form, which is how prometheus-query sends
			// long queries when preferPost is on.
			up, srv := newUpstream(t)
			req := tenancyPost(t, "/api/v1/query", url.Values{"namespace": {"ws-a"}}, url.Values{"query": {q}})
			if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, req); rec.Code != http.StatusOK {
				t.Fatalf("POST: status = %d, want 200 (%s)", rec.Code, rec.Body)
			}
			got := up.all()
			if len(got) != 1 || got[0].method != http.MethodPost {
				t.Fatalf("POST: upstream got %+v", got)
			}
			assertPinned(t, got[0].form.Get("query"), "ws-a")
		})
	}
}

func TestTenancy_MatchSelectors(t *testing.T) {
	t.Run("each match[] pinned", func(t *testing.T) {
		up, srv := newUpstream(t)
		params := url.Values{"match[]": {`up`, `{job="kubelet"}`, `kube_pod_info{namespace="ws-a"}`}, "namespace": {"ws-a"}}
		if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, tenancyGet(t, "/api/v1/series", params)); rec.Code != http.StatusOK {
			t.Fatalf("status = %d (%s)", rec.Code, rec.Body)
		}
		got := up.all()[0].query["match[]"]
		if len(got) != 3 {
			t.Fatalf("upstream match[] = %v, want 3", got)
		}
		for _, s := range got {
			assertSelectorPinned(t, s, "ws-a")
		}
	})

	for _, endpoint := range []string{"/api/v1/labels", "/api/v1/label/namespace/values", "/api/v1/series"} {
		t.Run("none given "+endpoint, func(t *testing.T) {
			up, srv := newUpstream(t)
			params := url.Values{"namespace": {"ws-a"}}
			if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, tenancyGet(t, endpoint, params)); rec.Code != http.StatusOK {
				t.Fatalf("status = %d (%s)", rec.Code, rec.Body)
			}
			got := up.all()[0].query["match[]"]
			if len(got) != 1 {
				t.Fatalf("upstream match[] = %v, want the injected one", got)
			}
			assertSelectorPinned(t, got[0], "ws-a")
		})
	}

	t.Run("POST body", func(t *testing.T) {
		up, srv := newUpstream(t)
		req := tenancyPost(t, "/api/v1/series", url.Values{"namespace": {"ws-a"}}, url.Values{"match[]": {`up`}})
		if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, req); rec.Code != http.StatusOK {
			t.Fatalf("status = %d (%s)", rec.Code, rec.Body)
		}
		got := up.all()[0]
		for _, s := range append(got.form["match[]"], got.query["match[]"]...) {
			assertSelectorPinned(t, s, "ws-a")
		}
		if len(got.form["match[]"]) != 1 {
			t.Errorf("body match[] = %v, want the enforced one", got.form["match[]"])
		}
	})
}

func TestTenancy_Metadata(t *testing.T) {
	up, srv := newUpstream(t)
	params := url.Values{"metric": {"kube_resourcequota"}, "namespace": {"ws-a"}}
	if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, tenancyGet(t, "/api/v1/metadata", params)); rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	got := up.all()[0]
	if got.path != "/__otterscale/proxy/api/v1/metadata" || got.query.Get("metric") != "kube_resourcequota" {
		t.Errorf("upstream got %s %v", got.path, got.query)
	}
}

func TestTenancy_ForbiddenEndpoints(t *testing.T) {
	paths := []string{
		"/api/v1/targets",
		"/api/v1/status/config",
		"/api/v1/status/tsdb",
		"/api/v1/rules",
		"/api/v1/alerts",
		"/api/v1/admin/tsdb/delete_series",
		"/api/v1/query/",
		"/api/v1/label/a/b/values",
		"/api/v1/label/pod/values/x",
		"/api/v1/query/%2e%2e/targets",
		"/federate",
	}
	for _, p := range paths {
		t.Run(p, func(t *testing.T) {
			up, srv := newUpstream(t)
			params := url.Values{"query": {"up"}, "namespace": {"ws-a"}}
			if rec := serveTenancy(t, clusterViewer(), srv.URL, tenancyGet(t, p, params)); rec.Code != http.StatusForbidden {
				t.Errorf("status = %d, want 403", rec.Code)
			}
			if n := len(up.all()); n != 0 {
				t.Errorf("upstream called %d times", n)
			}
		})
	}
}

// TestTenancy_ForwardHygiene: the namespace parameter is the hub's, not
// Prometheus'; and nothing the caller sent may reach the agent as a credential.
func TestTenancy_ForwardHygiene(t *testing.T) {
	up, srv := newUpstream(t)
	req := tenancyGet(t, "/api/v1/query", url.Values{"query": {"up"}, "namespace": {"ws-a"}})
	req.Header.Set("Authorization", "Bearer some-oidc-token")
	req.Header.Set("Impersonate-User", "system:admin")
	req.Header.Add("Impersonate-Group", "system:masters")
	req.Header.Set("Impersonate-Extra-Scopes", "everything")

	if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, req); rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	got := up.all()[0]
	if got.path != "/__otterscale/proxy/api/v1/query" {
		t.Errorf("upstream path = %q", got.path)
	}
	if got.query.Has("namespace") {
		t.Errorf("namespace parameter forwarded: %v", got.query)
	}
	for name := range got.header {
		if name == "Authorization" || strings.HasPrefix(name, "Impersonate-") {
			t.Errorf("caller header %s reached upstream", name)
		}
	}
}

// TestTenancy_RequestShape: every way of handing Prometheus a query the hub
// did not check is refused.
func TestTenancy_RequestShape(t *testing.T) {
	ns := url.Values{"namespace": {"ws-a"}}

	t.Run("query twice", func(t *testing.T) {
		_, srv := newUpstream(t)
		params := url.Values{"query": {"up", `kube_pod_info{namespace="other"}`}, "namespace": {"ws-a"}}
		if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, tenancyGet(t, "/api/v1/query", params)); rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", rec.Code)
		}
	})

	t.Run("query in URL and body", func(t *testing.T) {
		_, srv := newUpstream(t)
		req := tenancyPost(t, "/api/v1/query", url.Values{"query": {"up"}, "namespace": {"ws-a"}},
			url.Values{"query": {`kube_pod_info{namespace="other"}`}})
		if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, req); rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", rec.Code)
		}
	})

	t.Run("no query", func(t *testing.T) {
		_, srv := newUpstream(t)
		if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, tenancyGet(t, "/api/v1/query", ns)); rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", rec.Code)
		}
	})

	for _, ct := range []string{"multipart/form-data; boundary=x", "text/plain", "application/json", ""} {
		t.Run("POST as "+ct, func(t *testing.T) {
			up, srv := newUpstream(t)
			req := httptest.NewRequestWithContext(signedIn(t), http.MethodPost,
				"/proxy/prod/prometheus-tenancy/api/v1/query?namespace=ws-a",
				strings.NewReader("--x\r\nContent-Disposition: form-data; name=\"query\"\r\n\r\nkube_pod_info\r\n--x--\r\n"))
			if ct != "" {
				req.Header.Set("Content-Type", ct)
			}
			if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, req); rec.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", rec.Code)
			}
			if n := len(up.all()); n != 0 {
				t.Errorf("upstream called %d times", n)
			}
		})
	}

	t.Run("malformed query string", func(t *testing.T) {
		_, srv := newUpstream(t)
		req := httptest.NewRequestWithContext(signedIn(t), http.MethodGet,
			"/proxy/prod/prometheus-tenancy/api/v1/query?namespace=ws-a&query=up;x=1", http.NoBody)
		if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, req); rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", rec.Code)
		}
	})

	for _, method := range []string{http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			_, srv := newUpstream(t)
			req := httptest.NewRequestWithContext(signedIn(t), method,
				"/proxy/prod/prometheus-tenancy/api/v1/query?namespace=ws-a&query=up", http.NoBody)
			if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, req); rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("status = %d, want 405", rec.Code)
			}
		})
	}

	t.Run("GET body is not forwarded", func(t *testing.T) {
		up, srv := newUpstream(t)
		req := httptest.NewRequestWithContext(signedIn(t), http.MethodGet,
			"/proxy/prod/prometheus-tenancy/api/v1/query?namespace=ws-a&query=up",
			strings.NewReader("query=kube_pod_info"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if rec := serveTenancy(t, tenantOf("ws-a"), srv.URL, req); rec.Code != http.StatusOK {
			t.Fatalf("status = %d", rec.Code)
		}
		if got := up.all()[0].form; len(got) != 0 {
			t.Errorf("upstream received a body %v on GET", got)
		}
	})
}

func TestTenancy_ClusterNotFound(t *testing.T) {
	h := NewProxyHandler(&mockTunnelForProxy{addressErr: &core.ErrClusterNotFound{Cluster: "missing"}}, tenantOf("ws-a"))
	mux := http.NewServeMux()
	mux.HandleFunc("/proxy/{cluster}/prometheus-tenancy/{path...}", h.ServeTenancy)
	req := httptest.NewRequestWithContext(signedIn(t), http.MethodGet,
		"/proxy/missing/prometheus-tenancy/api/v1/query?namespace=ws-a&query=up", http.NoBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	assertPrometheusError(t, rec, http.StatusNotFound, "not_found")
}
