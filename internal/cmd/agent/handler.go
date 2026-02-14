package agent

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	utilproxy "k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
)

// Handler sets up the HTTP routes served by the agent. Its sole route
// is a reverse proxy to the local Kubernetes API server.
type Handler struct {
	cfg *rest.Config
}

// NewHandler returns a new agent Handler.
func NewHandler(cfg *rest.Config) *Handler {
	return &Handler{cfg: cfg}
}

// Mount registers a catch-all reverse proxy to the Kubernetes API
// server on the given mux. The proxy uses the in-cluster service
// account credentials (or falls back to KUBECONFIG) and rewrites
// the Host header so that the upstream kube-apiserver recognises
// the request.
func (h *Handler) Mount(mux *http.ServeMux) error {
	targetURL, err := url.Parse(h.cfg.Host)
	if err != nil {
		return fmt.Errorf("failed to parse k8s host URL: %w", err)
	}

	transport, err := rest.TransportFor(h.cfg)
	if err != nil {
		return fmt.Errorf("failed to create rest transport: %w", err)
	}

	proxy := utilproxy.NewUpgradeAwareHandler(targetURL, transport, false, false, &errorResponder{})
	mux.Handle("/", proxy)
	return nil
}

// errorResponder implements k8s.io/apimachinery/pkg/util/proxy.ErrorResponder.
// It logs errors and returns a 502 Bad Gateway response to the client.
type errorResponder struct{}

func (r *errorResponder) Error(w http.ResponseWriter, _ *http.Request, err error) {
	slog.Error("proxy error", "error", err)
	http.Error(w, "bad gateway", http.StatusBadGateway)
}
