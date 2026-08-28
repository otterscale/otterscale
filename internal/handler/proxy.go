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

// ProxyHandler is a raw HTTP reverse proxy that relays Prometheus
// queries from the dashboard frontend through the chisel tunnel to
// the in-cluster Prometheus service running alongside the agent. It
// validates paths against a read-only whitelist before forwarding.
//
// Authorization here is deliberately cluster-agnostic: any
// authenticated user may query any registered cluster's Prometheus.
// Metrics are treated as shared operational data, unlike the Kubernetes
// API paths, which enforce per-cluster RBAC by impersonating the caller
// against the target cluster. This is a decision, not an oversight — do
// not "fix" it without changing the product rule first.
//
// Namespace-level isolation is not provided and cannot be added here:
// the allowlist gates endpoints, not PromQL, so a caller past the gate
// can select any series. Restricting tenants to their own namespaces
// requires rewriting queries to carry an enforced label matcher
// (prom-label-proxy or equivalent) in front of Prometheus itself.
type ProxyHandler struct {
	tunnel core.TunnelProvider
}

// NewProxyHandler returns a ProxyHandler backed by the given
// TunnelProvider.
func NewProxyHandler(tunnel core.TunnelProvider) *ProxyHandler {
	return &ProxyHandler{tunnel: tunnel}
}

// ServeHTTP handles requests of the form
// /proxy/{cluster}/prometheus/{path...} by forwarding them through the
// tunnel to the agent's /__otterscale/proxy/{path} endpoint. The
// frontend configures prometheus-query with
// endpoint="/proxy/{cluster}/prometheus" and baseURL="/api/v1", so
// requests arrive as e.g. /proxy/my-cluster/prometheus/api/v1/query.
func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cluster := r.PathValue("cluster")

	// PathValue returns the decoded remainder, so it can still carry
	// "." and ".." segments; the allowlist checks the normalized form
	// and hands it back for forwarding.
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
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			if req.URL.Scheme == "" {
				req.URL.Scheme = "http"
			}
			req.URL.Host = target.Host
			req.URL.Path = "/__otterscale/proxy" + promPath
			req.URL.RawQuery = originalQuery
			req.Host = target.Host

			// Drop everything the caller sent that something
			// downstream could read as a credential or as an identity
			// assertion. Authorization belongs to the OIDC middleware,
			// not to Prometheus. The Impersonate-* family must never be
			// caller-controlled: anything that trusts those headers
			// would be trusting whoever made this request. Nothing on
			// this path reads them today, which is exactly why they are
			// cheap to remove now rather than after something does.
			//
			// net/http canonicalizes header names as it parses them, so
			// matching the canonical prefix also catches
			// Impersonate-Extra-*.
			req.Header.Del("Authorization")
			for name := range req.Header {
				if strings.HasPrefix(name, impersonationHeaderPrefix) {
					req.Header.Del(name)
				}
			}
		},
	}
	proxy.ServeHTTP(w, r) // #nosec G704
}
