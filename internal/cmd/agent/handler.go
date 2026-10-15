package agent

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	utilproxy "k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"

	"github.com/otterscale/otterscale/internal/config"
)

// proxyPathPrefix is reserved for agent-managed service proxies (Prometheus,
// for one). Requests under it skip kube-apiserver's proxy: they go straight to
// the service. An https target (kube-rbac-proxy) authorizes the agent's own
// ServiceAccount token itself; an http target gets no token and is trusted as
// reachable only from inside the cluster.
const proxyPathPrefix = "/__otterscale/proxy/"

// Handler is the spoke-side reverse proxy: it forwards requests to the local
// kube-apiserver and, optionally, to configured in-cluster services.
type Handler struct {
	cfg                     *rest.Config
	proxyPrometheusURL      string
	proxyPrometheusCASecret string
}

// NewHandler enables the metrics proxy under /__otterscale/proxy/ when a
// Prometheus URL is configured.
func NewHandler(cfg *rest.Config, conf *config.Config) *Handler {
	return &Handler{
		cfg:                     cfg,
		proxyPrometheusURL:      conf.AgentProxyPrometheusURL(),
		proxyPrometheusCASecret: conf.AgentProxyPrometheusCASecret(),
	}
}

// Mount always registers the catch-all proxy to kube-apiserver, plus one
// under /__otterscale/proxy/: the metrics proxy, so the hub can relay metrics
// queries without going through kube-apiserver, or, when it is off, a 503
// that says why.
func (h *Handler) Mount(mux *http.ServeMux) error {
	if h.proxyPrometheusURL == "" {
		mountDisabledPrometheusProxy(mux, "no Prometheus URL is configured")
	} else if err := h.mountPrometheusProxy(mux); err != nil {
		return err
	}

	return h.mountKubeProxy(mux)
}

// mountDisabledPrometheusProxy answers the metrics proxy's prefix with 503
// and why, which would otherwise fall through to the kube-apiserver proxy and
// come back as its 404.
func mountDisabledPrometheusProxy(mux *http.ServeMux, reason string) {
	message := "the metrics proxy is disabled on this cluster's agent: " + reason
	mux.Handle(proxyPathPrefix, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writePrometheusError(w, http.StatusServiceUnavailable, message)
	}))
}

func (h *Handler) mountPrometheusProxy(mux *http.ServeMux) error {
	target, err := url.Parse(h.proxyPrometheusURL)
	if err != nil {
		return fmt.Errorf("failed to parse proxy prometheus URL %q: %w", h.proxyPrometheusURL, err)
	}

	var loadCA caLoader
	if target.Scheme == "https" {
		if h.cfg.BearerToken == "" && h.cfg.BearerTokenFile == "" {
			slog.Warn("prometheus proxy disabled: kube-rbac-proxy needs the agent's ServiceAccount token, and this kube config carries none",
				"target", h.proxyPrometheusURL)
			mountDisabledPrometheusProxy(mux, "its kube config carries no ServiceAccount token")
			return nil
		}
		client, err := kubernetes.NewForConfig(h.cfg)
		if err != nil {
			return fmt.Errorf("failed to create kubernetes client: %w", err)
		}
		if loadCA, err = secretCALoader(client, h.proxyPrometheusCASecret); err != nil {
			return err
		}
	} else {
		slog.Warn("prometheus proxy target is not https: queries go out unauthenticated and the agent's token is withheld",
			"target", h.proxyPrometheusURL)
	}

	rt, err := newPrometheusTransport(target, h.cfg, loadCA)
	if err != nil {
		return err
	}

	proxy := newPrometheusProxy(target, rt, slog.NewLogLogger(slog.Default().Handler(), slog.LevelWarn))
	mux.Handle(proxyPathPrefix, http.StripPrefix("/__otterscale/proxy", proxy))
	slog.Info("prometheus proxy enabled", "target", h.proxyPrometheusURL)
	return nil
}

// mountKubeProxy adds WebSocket/SPDY upgrade support to the catch-all proxy.
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

// makeUpgradeTransport carries the same TLS/auth credentials as rt while
// handling the upgrade negotiations exec, attach and port-forward need.
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
