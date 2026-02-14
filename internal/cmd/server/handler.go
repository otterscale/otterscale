package server

import (
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	fleetv1 "github.com/otterscale/otterscale/api/fleet/v1/pbconnect"
	resourcev1 "github.com/otterscale/otterscale/api/resource/v1/pbconnect"
	runtimev1 "github.com/otterscale/otterscale/api/runtime/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/handler"
)

// Handler is responsible for mounting all gRPC service handlers,
// interceptors, and operational endpoints (health, reflection,
// metrics) onto an HTTP mux.
type Handler struct {
	fleet    *handler.FleetService
	resource *handler.ResourceService
	runtime  *handler.RuntimeService
	manifest *handler.ManifestHandler
}

// NewHandler returns a Handler for the given gRPC services and the
// raw HTTP manifest handler.
func NewHandler(fleet *handler.FleetService, resource *handler.ResourceService, runtime *handler.RuntimeService, manifest *handler.ManifestHandler) *Handler {
	return &Handler{
		fleet:    fleet,
		resource: resource,
		runtime:  runtime,
		manifest: manifest,
	}
}

// Mount registers all gRPC service handlers, OTel interceptors, and
// operational endpoints onto the provided mux.
func (h *Handler) Mount(mux *http.ServeMux) error {
	// OpenTelemetry interceptor for automatic tracing and metrics.
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return err
	}

	interceptors := connect.WithInterceptors(
		otelInterceptor,
	)

	// Operational endpoints: gRPC reflection, health checks, Prometheus.
	services := []string{
		fleetv1.FleetServiceName,
		resourcev1.ResourceServiceName,
		runtimev1.RuntimeServiceName,
	}

	if err := h.registerOpsHandlers(mux, services); err != nil {
		return err
	}

	// Application service handlers.
	// RPCs with idempotency_level = NO_SIDE_EFFECTS (e.g. GetAgentManifest)
	// automatically accept HTTP GET requests via the generated
	// connect.WithIdempotency(connect.IdempotencyNoSideEffects) option.
	mux.Handle(fleetv1.NewFleetServiceHandler(h.fleet, interceptors))
	mux.Handle(resourcev1.NewResourceServiceHandler(h.resource, interceptors))
	mux.Handle(runtimev1.NewRuntimeServiceHandler(h.runtime, interceptors))

	// Raw YAML endpoint for kubectl apply -f. Authentication is
	// handled by the HMAC token embedded in the URL path, so this
	// route is registered as a public path prefix in server.go.
	mux.HandleFunc("GET /fleet/manifest/{token}", h.handleRawManifest)

	return nil
}

// handleRawManifest verifies the HMAC token in the URL path and
// returns the agent installation manifest as raw YAML. This enables
// `kubectl apply -f <url>` without additional authentication headers.
func (h *Handler) handleRawManifest(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	cluster, userName, err := h.manifest.VerifyManifestToken(r.Context(), token)
	if err != nil {
		slog.Debug("manifest token verification failed", "error", err)
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	manifest, err := h.manifest.RenderManifest(r.Context(), cluster, userName)
	if err != nil {
		slog.Debug("manifest render failed", "cluster", cluster, "user", userName, "error", err)
		http.Error(w, "failed to render manifest", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	if _, err := w.Write([]byte(manifest)); err != nil {
		slog.Warn("failed to write manifest response", "error", err)
	}
}

// registerOpsHandlers sets up gRPC reflection, health checks, and
// Prometheus metrics scraping.
func (h *Handler) registerOpsHandlers(mux *http.ServeMux, serviceNames []string) error {
	reflector := grpcreflect.NewStaticReflector(serviceNames...)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	checker := grpchealth.NewStaticChecker(serviceNames...)
	mux.Handle(grpchealth.NewHandler(checker))

	exporter, err := prometheus.New()
	if err != nil {
		return err
	}
	// NOTE: This intentionally sets the global OTel MeterProvider so
	// that otelconnect interceptors and other libraries can discover
	// it without explicit injection. Ideally this would be injected
	// via Wire, but otelconnect relies on the global provider.
	otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(exporter)))
	mux.Handle("/metrics", promhttp.Handler())

	return nil
}
