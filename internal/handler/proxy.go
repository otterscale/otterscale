package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/otterscale/otterscale/internal/core"
)

// ProxyHandler is a raw HTTP reverse proxy that relays Prometheus
// queries from the dashboard frontend through the chisel tunnel to
// the in-cluster Prometheus service running alongside the agent. It
// validates paths against a read-only whitelist before forwarding.
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
	promPath := "/" + r.PathValue("path")

	if !core.IsAllowedPrometheusPath(promPath) {
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
			// Strip auth headers — they are for the OIDC
			// middleware, not for Prometheus.
			req.Header.Del("Authorization")
		},
	}
	proxy.ServeHTTP(w, r) // #nosec G704
}
