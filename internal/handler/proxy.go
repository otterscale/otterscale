package handler

import (
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
// tunnel to the in-cluster Prometheus beside the agent, validating paths
// against a read-only allowlist first.
//
// Authorization is deliberately cluster-agnostic: any authenticated user may
// query any registered cluster's Prometheus. Metrics are treated as shared
// operational data, unlike the Kubernetes API paths, which enforce per-cluster
// RBAC by impersonation. This is a decision, not an oversight — do not "fix" it
// without changing the product rule first.
//
// Namespace-level isolation is not provided and cannot be added here: the
// allowlist gates endpoints, not PromQL, so a caller past the gate can select
// any series. Restricting tenants to their own namespaces takes an enforced
// label matcher (prom-label-proxy or equivalent) in front of Prometheus.
type ProxyHandler struct {
	tunnel core.TunnelProvider
}

func NewProxyHandler(tunnel core.TunnelProvider) *ProxyHandler {
	return &ProxyHandler{tunnel: tunnel}
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
		http.Error(w, "forbidden prometheus path", http.StatusForbidden)
		return
	}

	address, err := h.tunnel.ResolveAddress(r.Context(), cluster)
	if err != nil {
		http.Error(w, "cluster not found", http.StatusNotFound)
		return
	}

	target, err := url.Parse(address)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	originalQuery := r.URL.RawQuery
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			out := pr.Out
			out.URL.Scheme = target.Scheme
			if out.URL.Scheme == "" {
				out.URL.Scheme = "http"
			}
			out.URL.Host = target.Host
			out.URL.Path = "/__otterscale/proxy" + promPath
			out.URL.RawQuery = originalQuery
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
