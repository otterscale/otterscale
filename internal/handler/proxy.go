package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/otterscale/otterscale/internal/core"
)

// impersonationHeaderPrefix covers Impersonate-User, Impersonate-Group,
// Impersonate-Uid and the open-ended Impersonate-Extra-* family.
const impersonationHeaderPrefix = "Impersonate-"

// agentProxyPrefix is where the agent serves its Prometheus proxy; forward
// puts the approved Prometheus path under it.
const agentProxyPrefix = "/__otterscale/proxy"

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
//     pods in <ns> or have cluster metrics access, rewritten to
//     namespace="<ns>" by prom-label-proxy's enforcer.
//   - Shared scope, /proxy/{cluster}/prometheus-shared/{name} (ServeShared,
//     proxy_shared.go): the hub's own named queries for cluster-wide numbers
//     that name no workload — per GPU card, per host, or a total — for the
//     namespace scope's callers too. Callers send no PromQL.
//
// This replaces the earlier rule that metrics were shared operational data any
// authenticated user could read. It mirrors OpenShift's monitoring tenancy
// (Thanos Querier :9091 and :9092), with the check in the hub rather than in
// each cluster, so the namespace and shared scope checks need no change to
// agents or managed clusters.
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

	h.forward(w, r, target, promPath, r.URL.RawQuery, "")
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
// contentType, when not empty, is the Content-Type of the body the hub wrote
// itself. It is set on the outgoing request, after the reverse proxy has
// removed the hop-by-hop headers, so the body is always read as what it is.
func (h *ProxyHandler) forward(w http.ResponseWriter, r *http.Request, target *url.URL, promPath, rawQuery, contentType string) {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			out := pr.Out
			out.URL.Scheme = target.Scheme
			if out.URL.Scheme == "" {
				out.URL.Scheme = "http"
			}
			out.URL.Host = target.Host
			out.URL.Path = agentProxyPrefix + promPath
			out.URL.RawPath = ""
			out.URL.RawQuery = rawQuery
			out.Host = target.Host
			if contentType != "" {
				out.Header.Set("Content-Type", contentType)
			}

			// Rewrite, unlike the deprecated Director, does not forward the
			// client address on its own. SetXForwarded overwrites rather than
			// appends, so an inbound X-Forwarded-For cannot be used to forge a
			// chain — the same reasoning as the header stripping below.
			pr.SetXForwarded()

			// Drop the caller's credentials and the impersonation headers.
			// Authorization and Cookie belong to the OIDC middleware and the
			// dashboard, not to Prometheus, and the Impersonate-* family
			// must never be caller-controlled: whatever trusted those headers
			// would be trusting whoever made this request. Nothing on this
			// path reads them today — which is exactly why they are cheap to
			// remove now rather than after something does. Other headers go
			// through; kube-rbac-proxy takes identity from the agent's token
			// alone.
			//
			// net/http canonicalizes header names as it parses them, so the
			// canonical prefix also catches Impersonate-Extra-*.
			out.Header.Del("Authorization")
			out.Header.Del("Cookie")
			for name := range out.Header {
				if strings.HasPrefix(name, impersonationHeaderPrefix) {
					out.Header.Del(name)
				}
			}
		},
		// The default answers an empty 502. The error, which names the
		// tunnel address, goes to the log only.
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			// A caller that went away, a page navigating off say, cancels the
			// request: nothing is wrong with the agent.
			if r.Context().Err() != nil {
				slog.Debug("metrics proxy request canceled by the caller", "cluster", r.PathValue("cluster"), "error", err)
			} else {
				slog.Warn("metrics proxy cannot reach the agent", "cluster", r.PathValue("cluster"), "error", err)
			}
			writePrometheusError(w, http.StatusBadGateway, "unavailable", "cannot reach this cluster's agent")
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
