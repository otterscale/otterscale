package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/otterscale/otterscale/internal/core"
)

// impersonationHeaderPrefix covers Impersonate-User, Impersonate-Group,
// Impersonate-Uid and the open-ended Impersonate-Extra-* family.
const impersonationHeaderPrefix = "Impersonate-"

// ProxyHandler relays Prometheus queries from the dashboard through the chisel
// tunnel to the agent, which queries the cluster's Prometheus with its own
// credentials. Who may see what is therefore decided here, per cluster, by
// asking that cluster's RBAC as the caller (core.AccessReviewer, a cached
// SelfSubjectAccessReview):
//
//   - Cluster scope, /proxy/{cluster}/prometheus/... (ServeHTTP): any PromQL,
//     for callers with core.ClusterMetricsAccess (cluster-admin, or a binding
//     to otterscale:metrics-cluster-view). Everyone else gets 403 whatever the
//     query: aggregates leak other namespaces as surely as series do.
//   - Namespace scope, /proxy/{cluster}/prometheus-tenancy/...?namespace=<ns>
//     (ServeTenancy, proxy_tenancy.go): any PromQL, for callers who may get
//     pods in <ns> or have cluster metrics access, with namespace="<ns>"
//     forced into every selector.
//
// This replaces the earlier rule that metrics were shared operational data any
// authenticated user could read. It mirrors OpenShift's monitoring tenancy
// (Thanos Querier :9091 and :9092), with the check in the hub rather than in
// each cluster, so agents and managed clusters need no change.
type ProxyHandler struct {
	tunnel core.TunnelProvider
	access core.AccessReviewer
}

func NewProxyHandler(tunnel core.TunnelProvider, access core.AccessReviewer) *ProxyHandler {
	return &ProxyHandler{tunnel: tunnel, access: access}
}

// ServeHTTP forwards /proxy/{cluster}/prometheus/{path...} through the tunnel
// to the agent's /__otterscale/proxy/{path}. The frontend configures
// prometheus-query with endpoint="/proxy/{cluster}/prometheus" and
// baseURL="/api/v1", so requests arrive as
// /proxy/my-cluster/prometheus/api/v1/query.
func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cluster := r.PathValue("cluster")

	// PathValue returns the decoded remainder, which can still carry "." and
	// ".." segments; the allowlist checks the normalized form and hands it back
	// for forwarding.
	promPath, allowed := core.AllowedPrometheusPath("/" + r.PathValue("path"))
	if !allowed {
		writePrometheusError(w, http.StatusForbidden, "forbidden", "forbidden prometheus path")
		return
	}

	target, ok := h.resolve(w, r, cluster)
	if !ok {
		return
	}

	// Cluster scope is all or nothing: the hub does not read the PromQL here,
	// because an aggregate over other namespaces (count, or a difference of
	// sums) leaks them as surely as their series do.
	if !h.authorize(w, r, cluster, core.ClusterMetricsAccess) {
		return
	}

	h.forward(w, r, target, promPath, r.URL.RawQuery)
}

// resolve writes the error itself and reports whether to go on.
func (h *ProxyHandler) resolve(w http.ResponseWriter, r *http.Request, cluster string) (*url.URL, bool) {
	address, err := h.tunnel.ResolveAddress(r.Context(), cluster)
	if err != nil {
		writePrometheusError(w, http.StatusNotFound, "not_found", "cluster not found")
		return nil, false
	}

	target, err := url.Parse(address)
	if err != nil {
		writePrometheusError(w, http.StatusInternalServerError, "internal", "internal error")
		return nil, false
	}
	return target, true
}

// authorize asks the questions in order and lets the request through on the
// first yes. It writes the refusal itself.
//
// 401 means the hub has no user for this request — read from the same
// UserInfo the reviews impersonate, before any cluster is asked. Every review
// error after that is 503, a cluster's own 401 included (the caller is signed
// in; it is the agent's credentials that were refused): never a pass, and
// never a "no" either, since the caller may well be allowed. The reviewer
// logs the failure, once per review rather than once per query.
func (h *ProxyHandler) authorize(w http.ResponseWriter, r *http.Request, cluster string, attrs ...core.AccessAttributes) bool {
	if user, ok := core.UserInfoFromContext(r.Context()); !ok || user.Subject == "" {
		writePrometheusError(w, http.StatusUnauthorized, "unauthorized", "no authenticated user")
		return false
	}
	for _, a := range attrs {
		allowed, err := h.access.Allowed(r.Context(), cluster, a)
		if err != nil {
			var notFound *core.ErrClusterNotFound
			if errors.As(err, &notFound) {
				writePrometheusError(w, http.StatusNotFound, "not_found", "cluster not found")
			} else {
				writePrometheusError(w, http.StatusServiceUnavailable, "unavailable", "cannot check metrics access on this cluster; retry shortly")
			}
			return false
		}
		if allowed {
			return true
		}
	}
	writePrometheusError(w, http.StatusForbidden, "forbidden", "not allowed to read these metrics")
	return false
}

// forward relays to the agent with the path and query given, which are the
// ones the checks approved — never the request's own, unchecked ones.
func (h *ProxyHandler) forward(w http.ResponseWriter, r *http.Request, target *url.URL, promPath, rawQuery string) {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			out := pr.Out
			out.URL.Scheme = target.Scheme
			if out.URL.Scheme == "" {
				out.URL.Scheme = "http"
			}
			out.URL.Host = target.Host
			out.URL.Path = "/__otterscale/proxy" + promPath
			out.URL.RawPath = ""
			out.URL.RawQuery = rawQuery
			out.Host = target.Host

			// Rewrite, unlike the deprecated Director, does not forward the
			// client address on its own. SetXForwarded overwrites rather than
			// appends, so an inbound X-Forwarded-For cannot be used to forge a
			// chain — the same reasoning as the header stripping below.
			pr.SetXForwarded()

			// Drop anything downstream could read as a credential or an
			// identity assertion. Authorization belongs to the OIDC
			// middleware, not to Prometheus, and the Impersonate-* family
			// must never be caller-controlled: whatever trusted those headers
			// would be trusting whoever made this request. Nothing on this
			// path reads them today — which is exactly why they are cheap to
			// remove now rather than after something does.
			//
			// net/http canonicalizes header names as it parses them, so the
			// canonical prefix also catches Impersonate-Extra-*.
			out.Header.Del("Authorization")
			for name := range out.Header {
				if strings.HasPrefix(name, impersonationHeaderPrefix) {
					out.Header.Del(name)
				}
			}
		},
	}
	proxy.ServeHTTP(w, r) // #nosec G704
}

// writePrometheusError answers in Prometheus' own error shape, which
// prometheus-query surfaces with its message instead of "unexpected client
// error".
func writePrometheusError(w http.ResponseWriter, status int, errorType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(struct {
		Status    string `json:"status"`
		ErrorType string `json:"errorType"`
		Error     string `json:"error"`
	}{"error", errorType, message})
}
