package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/otterscale/otterscale/internal/core"
)

func serveShared(t *testing.T, access core.AccessReviewer, backend string, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	h := NewProxyHandler(&mockTunnelForProxy{address: backend}, access)
	rec := httptest.NewRecorder()
	newProxyMux(h).ServeHTTP(rec, req)
	return rec
}

func sharedGet(t *testing.T, name string, params url.Values) *http.Request {
	t.Helper()
	return httptest.NewRequestWithContext(signedIn(t), http.MethodGet,
		"/proxy/prod/prometheus-shared/"+name+"?"+params.Encode(), http.NoBody)
}

func manyUUIDs(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = fmt.Sprintf("GPU-%d", i)
	}
	return out
}

// TestShared_ForwardsTheHubsQuery: what reaches the agent is the catalog's
// PromQL with the caller's values, at the endpoint the query belongs to, and
// nothing else the caller sent.
func TestShared_ForwardsTheHubsQuery(t *testing.T) {
	cases := []struct {
		desc   string
		name   string
		access *fakeAccess
		params url.Values
		path   string
		want   url.Values
	}{
		{
			desc: "instant with uuids and time", name: "gpu-util", access: tenantOf("ws-a"),
			params: url.Values{"namespace": {"ws-a"}, "uuid": {"GPU-a", "GPU-b"}, "time": {"1700000000"}},
			path:   "/__otterscale/proxy/api/v1/query",
			want: url.Values{
				"query": {`max by (UUID) (DCGM_FI_DEV_GPU_UTIL{UUID=~"GPU-a|GPU-b"})`},
				"time":  {"1700000000"},
			},
		},
		{
			desc: "range", name: "host-gpu-fb-used", access: tenantOf("ws-a"),
			params: url.Values{"namespace": {"ws-a"}, "hostname": {"hci"}, "start": {"1700000000"}, "end": {"1700003600"}, "step": {"120"}},
			path:   "/__otterscale/proxy/api/v1/query_range",
			want: url.Values{
				"query": {`sum by (hostname) (DCGM_FI_DEV_FB_USED{hostname=~"hci"})`},
				"start": {"1700000000"}, "end": {"1700003600"}, "step": {"120"},
			},
		},
		{
			desc: "no parameters, cluster viewer without namespace", name: "kubevirt-handlers", access: clusterViewer(),
			params: url.Values{},
			path:   "/__otterscale/proxy/api/v1/query",
			want:   url.Values{"query": {`count(kubevirt_info)`}},
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			up, srv := newUpstream(t)
			if rec := serveShared(t, c.access, srv.URL, sharedGet(t, c.name, c.params)); rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200 (%s)", rec.Code, rec.Body)
			}
			got := up.all()
			if len(got) != 1 {
				t.Fatalf("upstream called %d times, want 1", len(got))
			}
			if got[0].method != http.MethodGet || got[0].path != c.path {
				t.Errorf("upstream %s %s, want GET %s", got[0].method, got[0].path, c.path)
			}
			if !reflect.DeepEqual(got[0].query, c.want) {
				t.Errorf("upstream query = %v, want %v", got[0].query, c.want)
			}
			if len(got[0].form) != 0 {
				t.Errorf("upstream received a body %v", got[0].form)
			}
		})
	}
}

// TestShared_Parameters: anything a query does not declare, or declares in
// another shape, is 400 before any cluster is asked.
func TestShared_Parameters(t *testing.T) {
	cases := []struct {
		desc   string
		name   string
		params url.Values
	}{
		{"PromQL", "gpu-util", url.Values{"uuid": {"GPU-a"}, "query": {"up"}}},
		{"match[]", "kubevirt-handlers", url.Values{"match[]": {"up"}}},
		{"a parameter the query does not take", "kubevirt-handlers", url.Values{"uuid": {"GPU-a"}}},
		{"the other repeatable parameter", "gpu-util", url.Values{"uuid": {"GPU-a"}, "hostname": {"hci"}}},
		{"missing uuid", "gpu-util", url.Values{}},
		{"empty uuid", "gpu-util", url.Values{"uuid": {""}}},
		{"uuid with a quote", "gpu-util", url.Values{"uuid": {`GPU-a"}`}}},
		{"uuid with a pipe", "gpu-util", url.Values{"uuid": {"GPU-a|.*"}}},
		{"uuid too long", "gpu-util", url.Values{"uuid": {strings.Repeat("a", 129)}}},
		{"too many uuids", "gpu-util", url.Values{"uuid": manyUUIDs(sharedMaxValues + 1)}},
		{"hostname not a node name", "host-gpu-count", url.Values{"hostname": {"Bad_Host"}}},
		{"time on a range query", "host-gpu-fb-used", url.Values{"hostname": {"hci"}, "start": {"1"}, "end": {"2"}, "step": {"1"}, "time": {"1"}}},
		{"start on an instant query", "kubevirt-handlers", url.Values{"start": {"1"}}},
		{"range without step", "host-gpu-fb-used", url.Values{"hostname": {"hci"}, "start": {"1"}, "end": {"2"}}},
		{"step zero", "host-gpu-fb-used", url.Values{"hostname": {"hci"}, "start": {"1"}, "end": {"2"}, "step": {"0"}}},
		{"time not a number", "kubevirt-handlers", url.Values{"time": {"yesterday"}}},
		{"time twice", "kubevirt-handlers", url.Values{"time": {"1", "2"}}},
		{"namespace twice", "kubevirt-handlers", url.Values{"namespace": {"ws-a", "ws-b"}}},
		{"namespace not a label", "kubevirt-handlers", url.Values{"namespace": {"WS_A"}}},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			if !c.params.Has("namespace") {
				c.params.Set("namespace", "ws-a")
			}
			up, srv := newUpstream(t)
			access := tenantOf("ws-a")
			rec := serveShared(t, access, srv.URL, sharedGet(t, c.name, c.params))
			assertPrometheusError(t, rec, http.StatusBadRequest, "bad_data")
			if n := len(up.all()); n != 0 {
				t.Errorf("upstream called %d times", n)
			}
			if asked := access.questions(); len(asked) != 0 {
				t.Errorf("asked %v before the parameters were checked", asked)
			}
		})
	}

	t.Run("at the cap", func(t *testing.T) {
		up, srv := newUpstream(t)
		params := url.Values{"namespace": {"ws-a"}, "uuid": manyUUIDs(sharedMaxValues)}
		if rec := serveShared(t, tenantOf("ws-a"), srv.URL, sharedGet(t, "gpu-util", params)); rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (%s)", rec.Code, rec.Body)
		}
		if n := len(up.all()); n != 1 {
			t.Errorf("upstream called %d times, want 1", n)
		}
	})
}

// sharedURILen is the length of the request URI the hub forwards for name
// and params: path, "?" and encoded query.
func sharedURILen(t *testing.T, name string, params url.Values) int {
	t.Helper()
	q := sharedCatalog[name]
	forwarded, _, err := parseSharedRequest(sharedGet(t, name, params), q)
	if err != nil {
		t.Fatal(err)
	}
	return len(agentProxyPrefix+q.endpoint()) + len("?") + len(forwarded.Encode())
}

// TestShared_ForwardedURILimit: the PromQL travels in the URL, and the agent,
// like the hub, refuses request headers over 8 KiB with 431. The hub refuses
// a request whose forwarded URI would be over maxSharedRequestURIBytes with
// its own 400 instead, and the largest it accepts still gets through.
func TestShared_ForwardedURILimit(t *testing.T) {
	var reached []int
	var mu sync.Mutex
	agent := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		reached = append(reached, len(r.RequestURI))
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[]}}`))
	}))
	agent.Config.MaxHeaderBytes = 8 << 10 // as the agent's server sets it
	agent.Start()
	t.Cleanup(agent.Close)

	// The server is as strict as the agent's: well past its limit, 431.
	long, err := http.NewRequestWithContext(t.Context(), http.MethodGet, agent.URL+"/?pad="+strings.Repeat("a", 16<<10), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := agent.Client().Do(long)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusRequestHeaderFieldsTooLarge {
		t.Fatalf("a 16 KiB URI got %d, want 431: the server does not stand in for the agent's", resp.StatusCode)
	}

	// The most cards whose forwarded URI fits, then the last UUID padded to
	// the limit exactly: one more byte is over it.
	params := url.Values{"namespace": {"ws-a"}, "time": {"1700000000"}}
	for i := 0; ; i++ {
		params.Add("uuid", fmt.Sprintf("GPU-%08d-0000-0000-0000-000000000000", i)) // as long as a real one
		if sharedURILen(t, "gpu-util", params) > maxSharedRequestURIBytes {
			uuids := params["uuid"]
			params["uuid"] = uuids[:len(uuids)-1]
			break
		}
	}
	uuids := params["uuid"]
	uuids[len(uuids)-1] += strings.Repeat("a", maxSharedRequestURIBytes-sharedURILen(t, "gpu-util", params))
	if n := sharedURILen(t, "gpu-util", params); n != maxSharedRequestURIBytes {
		t.Fatalf("built a %d-byte URI, want %d", n, maxSharedRequestURIBytes)
	}

	t.Run("at the limit", func(t *testing.T) {
		rec := serveShared(t, tenantOf("ws-a"), agent.URL, sharedGet(t, "gpu-util", params))
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (%s)", rec.Code, rec.Body)
		}
		mu.Lock()
		defer mu.Unlock()
		if len(reached) != 1 || reached[0] != maxSharedRequestURIBytes {
			t.Errorf("agent got request URIs of %v bytes, want one of %d", reached, maxSharedRequestURIBytes)
		}
	})

	t.Run("one byte over", func(t *testing.T) {
		uuids[len(uuids)-1] += "a"
		access := tenantOf("ws-a")
		rec := serveShared(t, access, agent.URL, sharedGet(t, "gpu-util", params))
		assertPrometheusError(t, rec, http.StatusBadRequest, "bad_data")
		mu.Lock()
		defer mu.Unlock()
		if len(reached) != 1 {
			t.Errorf("agent called again, with URIs of %v bytes", reached[1:])
		}
		if asked := access.questions(); len(asked) != 0 {
			t.Errorf("asked %v before the request was checked", asked)
		}
	})
}

// TestShared_Access: namespace only proves membership, checked like the
// namespace scope; without it only cluster viewers get through.
func TestShared_Access(t *testing.T) {
	cases := []struct {
		desc      string
		access    *fakeAccess
		namespace string
		status    int
		errorType string
		asked     []core.AccessAttributes
	}{
		{
			"member of the namespace, never asked about the cluster", tenantOf("ws-a"), "ws-a", http.StatusOK, "",
			[]core.AccessAttributes{core.NamespaceMetricsAccess("ws-a")},
		},
		{
			"another namespace", tenantOf("ws-a"), "ws-b", http.StatusForbidden, "forbidden",
			[]core.AccessAttributes{core.NamespaceMetricsAccess("ws-b"), core.ClusterMetricsAccess},
		},
		{
			"no namespace, not a cluster viewer", tenantOf("ws-a"), "", http.StatusForbidden, "forbidden",
			[]core.AccessAttributes{core.ClusterMetricsAccess},
		},
		{
			"no namespace, cluster viewer", clusterViewer(), "", http.StatusOK, "",
			[]core.AccessAttributes{core.ClusterMetricsAccess},
		},
		{
			"cluster viewer naming any namespace", clusterViewer(), "ws-b", http.StatusOK, "",
			[]core.AccessAttributes{core.NamespaceMetricsAccess("ws-b"), core.ClusterMetricsAccess},
		},
		{
			"review failure", &fakeAccess{err: errors.New("tunnel down")}, "ws-a", http.StatusServiceUnavailable, "unavailable",
			[]core.AccessAttributes{core.NamespaceMetricsAccess("ws-a")},
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			params := url.Values{}
			if c.namespace != "" {
				params.Set("namespace", c.namespace)
			}
			up, srv := newUpstream(t)
			rec := serveShared(t, c.access, srv.URL, sharedGet(t, "kubevirt-handlers", params))
			if c.status == http.StatusOK {
				if rec.Code != http.StatusOK {
					t.Fatalf("status = %d, want 200 (%s)", rec.Code, rec.Body)
				}
				if n := len(up.all()); n != 1 {
					t.Errorf("upstream called %d times, want 1", n)
				}
			} else {
				assertPrometheusError(t, rec, c.status, c.errorType)
				if n := len(up.all()); n != 0 {
					t.Errorf("upstream called %d times", n)
				}
			}
			if got := c.access.questions(); !reflect.DeepEqual(got, c.asked) {
				t.Errorf("asked %v, want %v", got, c.asked)
			}
		})
	}
}

// TestShared_RequestShape: unknown names and other methods, in Prometheus'
// error shape, and nothing upstream.
func TestShared_RequestShape(t *testing.T) {
	for _, path := range []string{
		"/proxy/prod/prometheus-shared/no-such-query?namespace=ws-a",
		"/proxy/prod/prometheus-shared/gpu-util/extra?namespace=ws-a&uuid=GPU-a",
		"/proxy/prod/prometheus-shared/?namespace=ws-a",
	} {
		t.Run(path, func(t *testing.T) {
			up, srv := newUpstream(t)
			req := httptest.NewRequestWithContext(signedIn(t), http.MethodGet, path, http.NoBody)
			assertPrometheusError(t, serveShared(t, tenantOf("ws-a"), srv.URL, req), http.StatusNotFound, "not_found")
			if n := len(up.all()); n != 0 {
				t.Errorf("upstream called %d times", n)
			}
		})
	}
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			up, srv := newUpstream(t)
			req := httptest.NewRequestWithContext(signedIn(t), method,
				"/proxy/prod/prometheus-shared/kubevirt-handlers?namespace=ws-a", http.NoBody)
			rec := serveShared(t, tenantOf("ws-a"), srv.URL, req)
			assertPrometheusError(t, rec, http.StatusMethodNotAllowed, "bad_data")
			if got := rec.Header().Get("Allow"); got != http.MethodGet {
				t.Errorf("Allow = %q, want %q", got, http.MethodGet)
			}
			if n := len(up.all()); n != 0 {
				t.Errorf("upstream called %d times", n)
			}
		})
	}
}

// TestShared_ForwardHygiene: no caller credential or identity header reaches
// the agent.
func TestShared_ForwardHygiene(t *testing.T) {
	up, srv := newUpstream(t)
	req := sharedGet(t, "kubevirt-handlers", url.Values{"namespace": {"ws-a"}})
	req.Header.Set("Authorization", "Bearer some-oidc-token")
	req.Header.Set("Impersonate-User", "system:admin")
	req.Header.Add("Impersonate-Group", "system:masters")
	req.Header.Set("Cookie", "session=some-session")
	if rec := serveShared(t, tenantOf("ws-a"), srv.URL, req); rec.Code != http.StatusOK {
		t.Fatalf("status = %d (%s)", rec.Code, rec.Body)
	}
	for name := range up.all()[0].header {
		if name == "Authorization" || name == "Cookie" || strings.HasPrefix(name, "Impersonate-") {
			t.Errorf("caller header %s reached upstream", name)
		}
	}
}

// TestShared_PrometheusErrorsPassThrough: Prometheus' own refusal — here a
// range too fine for its 11000-point limit — reaches the caller unchanged.
func TestShared_PrometheusErrorsPassThrough(t *testing.T) {
	const body = `{"status":"error","errorType":"bad_data","error":"exceeded maximum resolution of 11,000 points per timeseries. Try decreasing the query resolution (?step=XX)"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	params := url.Values{"namespace": {"ws-a"}, "hostname": {"hci"}, "start": {"0"}, "end": {"86400000"}, "step": {"1"}}
	rec := serveShared(t, tenantOf("ws-a"), srv.URL, sharedGet(t, "host-gpu-fb-used", params))
	if rec.Code != http.StatusBadRequest || rec.Body.String() != body {
		t.Errorf("got %d %s, want 400 %s", rec.Code, rec.Body, body)
	}
}

func TestShared_ClusterNotFound(t *testing.T) {
	h := NewProxyHandler(&mockTunnelForProxy{addressErr: &core.ErrClusterNotFound{Cluster: "missing"}}, tenantOf("ws-a"))
	req := httptest.NewRequestWithContext(signedIn(t), http.MethodGet,
		"/proxy/missing/prometheus-shared/kubevirt-handlers?namespace=ws-a", http.NoBody)
	rec := httptest.NewRecorder()
	newProxyMux(h).ServeHTTP(rec, req)
	assertPrometheusError(t, rec, http.StatusNotFound, "not_found")
}
