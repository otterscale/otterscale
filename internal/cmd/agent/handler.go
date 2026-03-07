package agent

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	utilproxy "k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"

	"github.com/otterscale/otterscale/internal/config"
)

// proxyPathPrefix is the URL prefix reserved for service proxies
// managed by the agent (e.g. Prometheus). Requests under this prefix
// are routed to a dedicated reverse proxy instead of kube-apiserver.
const proxyPathPrefix = "/__otterscale/proxy/"

// Handler is a reverse proxy that forwards incoming HTTP requests to
// the local kube-apiserver and, optionally, to configured in-cluster
// services (e.g. Prometheus). It is the spoke-side component in the
// hub-and-spoke architecture.
type Handler struct {
	cfg                *rest.Config
	proxyPrometheusURL string
}

// NewHandler creates a Handler backed by the given Kubernetes client
// config. proxyPrometheusURL, when non-empty, enables the Prometheus
// metrics proxy under the /__otterscale/proxy/ path prefix.
func NewHandler(cfg *rest.Config, conf *config.Config) *Handler {
	return &Handler{cfg: cfg, proxyPrometheusURL: conf.AgentProxyPrometheusURL()}
}

// Mount registers HTTP handlers on mux. A catch-all reverse proxy to
// kube-apiserver is always registered. When a Prometheus URL is
// configured, a dedicated reverse proxy is registered under
// /__otterscale/proxy/ so that the hub can relay metrics queries
// without going through kube-apiserver.
func (h *Handler) Mount(mux *http.ServeMux) error {
	if h.proxyPrometheusURL != "" {
		if err := h.mountPrometheusProxy(mux); err != nil {
			return err
		}
	}

	return h.mountKubeProxy(mux)
}

// mountPrometheusProxy registers a reverse proxy that forwards
// requests under /__otterscale/proxy/ to the configured Prometheus
// service URL.
func (h *Handler) mountPrometheusProxy(mux *http.ServeMux) error {
	target, err := url.Parse(h.proxyPrometheusURL)
	if err != nil {
		return fmt.Errorf("failed to parse proxy prometheus URL %q: %w", h.proxyPrometheusURL, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorLog = slog.NewLogLogger(slog.Default().Handler(), slog.LevelWarn)

	mux.Handle(proxyPathPrefix, http.StripPrefix("/__otterscale/proxy", proxy))
	slog.Info("prometheus proxy enabled", "target", h.proxyPrometheusURL)
	return nil
}

// mountKubeProxy registers the catch-all reverse proxy to
// kube-apiserver with WebSocket/SPDY upgrade support.
func (h *Handler) mountKubeProxy(mux *http.ServeMux) error {
	host := h.cfg.Host
	if !strings.HasSuffix(host, "/") {
		host += "/"
	}
	targetURL, err := url.Parse(host)
	if err != nil {
		return fmt.Errorf("failed to parse k8s host URL: %w", err)
	}

	rt, err := rest.TransportFor(h.cfg)
	if err != nil {
		return fmt.Errorf("failed to create rest transport: %w", err)
	}

	upgradeRT, err := makeUpgradeTransport(h.cfg, rt)
	if err != nil {
		return fmt.Errorf("failed to create upgrade transport: %w", err)
	}

	proxy := utilproxy.NewUpgradeAwareHandler(targetURL, rt, false, false, &errorResponder{})
	proxy.UpgradeTransport = upgradeRT
	proxy.UseRequestLocation = true
	proxy.UseLocationHost = true
	mux.Handle("/", proxy)
	return nil
}

// errorResponder logs proxy errors and returns 502 Bad Gateway.
type errorResponder struct{}

func (r *errorResponder) Error(w http.ResponseWriter, _ *http.Request, err error) {
	slog.Error("proxy error", "error", err)
	http.Error(w, "bad gateway", http.StatusBadGateway)
}

// makeUpgradeTransport builds an UpgradeRequestRoundTripper that
// carries the same TLS/auth credentials as rt but handles HTTP
// upgrade negotiations (WebSocket, SPDY) used by exec/attach/port-forward.
func makeUpgradeTransport(cfg *rest.Config, rt http.RoundTripper) (utilproxy.UpgradeRequestRoundTripper, error) {
	transportConfig, err := cfg.TransportConfig()
	if err != nil {
		return nil, err
	}

	upgrader, err := transport.HTTPWrappersForConfig(transportConfig, utilproxy.MirrorRequest)
	if err != nil {
		return nil, err
	}

	return utilproxy.NewUpgradeRequestRoundTripper(rt, upgrader), nil
}
