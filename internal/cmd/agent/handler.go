package agent

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	utilproxy "k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
)

// Handler is a reverse proxy that forwards all incoming HTTP requests
// to the local kube-apiserver. It is the spoke-side component in the
// hub-and-spoke architecture: the hub tunnels user requests into the
// spoke cluster, and Handler proxies them to the real API server.
type Handler struct {
	cfg *rest.Config
}

// NewHandler creates a Handler backed by the given Kubernetes client config.
func NewHandler(cfg *rest.Config) *Handler {
	return &Handler{cfg: cfg}
}

// Mount registers a catch-all reverse proxy on mux that supports both
// regular HTTP requests and WebSocket/SPDY upgrades (required by
// kubectl exec, attach, and port-forward).
func (h *Handler) Mount(mux *http.ServeMux) error {
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
