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

// ProxyHandler relays Prometheus queries from the dashboard through the chisel
// tunnel to the agent, which queries the cluster's Prometheus with its own
// credentials. Who may see what is therefore decided here, per cluster, by
// asking that cluster's RBAC as the caller (core.AccessReviewer).
//
// ServeHTTP is cluster scope, /proxy/{cluster}/prometheus/...; ServeTenancy is
// namespace scope, /proxy/{cluster}/prometheus-tenancy/...?namespace=<ns>
// (proxy_tenancy.go).
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
		http.Error(w, "forbidden prometheus path", http.StatusForbidden)
		return
	}

	target, ok := h.resolve(w, r, cluster)
	if !ok {
		return
	}

	h.forward(w, r, target, promPath, r.URL.RawQuery)
}

// resolve writes the error itself and reports whether to go on.
func (h *ProxyHandler) resolve(w http.ResponseWriter, r *http.Request, cluster string) (*url.URL, bool) {
	address, err := h.tunnel.ResolveAddress(r.Context(), cluster)
	if err != nil {
		http.Error(w, "cluster not found", http.StatusNotFound)
		return nil, false
	}

	target, err := url.Parse(address)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return nil, false
	}
	return target, true
}

// authorize asks the questions in order and lets the request through on the
// first yes. It writes the refusal itself. A review that fails is 503 — never
// a pass, and never a "no" either: the caller may well be allowed.
func (h *ProxyHandler) authorize(w http.ResponseWriter, r *http.Request, cluster string, attrs ...core.AccessAttributes) bool {
	for _, a := range attrs {
		allowed, err := h.access.Allowed(r.Context(), cluster, a)
		if err != nil {
			var notFound *core.ErrClusterNotFound
			switch code, _ := core.DomainErrorCode(err); {
			case errors.As(err, &notFound):
				writePrometheusError(w, http.StatusNotFound, "not_found", "cluster not found")
			case code == core.ErrorCodeUnauthenticated:
				writePrometheusError(w, http.StatusUnauthorized, "unauthorized", "no authenticated user")
			default:
				slog.Warn("metrics access review failed", "cluster", cluster, "resource", a.Resource, "namespace", a.Namespace, "error", err)
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
