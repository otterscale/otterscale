package server

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	linkv1 "github.com/otterscale/otterscale/api/link/v1"
	resourcev1 "github.com/otterscale/otterscale/api/resource/v1"
	runtimev1 "github.com/otterscale/otterscale/api/runtime/v1"

	"github.com/otterscale/otterscale/internal/handler"
)

// Handler is responsible for mounting all gRPC service handlers,
// interceptors, and operational endpoints (health, reflection,
// metrics) onto an HTTP mux.
type Handler struct {
	link     *handler.LinkService
	resource *handler.ResourceService
	runtime  *handler.RuntimeService
	proxy    *handler.ProxyHandler
}

// NewHandler returns a Handler for the given gRPC services and the
// Prometheus reverse proxy handler.
func NewHandler(link *handler.LinkService, resource *handler.ResourceService, runtime *handler.RuntimeService, proxy *handler.ProxyHandler) *Handler {
	return &Handler{
		link:     link,
		resource: resource,
		runtime:  runtime,
		proxy:    proxy,
	}
}

// LongRunningPaths returns the procedures whose response is a
// long-lived stream. The transport lifts its request timeouts for
// these so that a watch or an exec session is not cut off mid-flight.
func (h *Handler) LongRunningPaths() []string {
	return []string{
		resourcev1.ResourceServiceWatchProcedure,
		runtimev1.RuntimeServicePodLogProcedure,
		runtimev1.RuntimeServiceExecuteTTYProcedure,
		runtimev1.RuntimeServicePortForwardProcedure,
		runtimev1.RuntimeServiceVNCProcedure,
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
		linkv1.LinkServiceName,
		resourcev1.ResourceServiceName,
		runtimev1.RuntimeServiceName,
	}

	if err := h.registerOpsHandlers(mux, services); err != nil {
		return err
	}

	// Application service handlers.
	// RPCs with idempotency_level = NO_SIDE_EFFECTS automatically
	// accept HTTP GET requests via the generated
	// connect.WithIdempotency(connect.IdempotencyNoSideEffects) option.
	mux.Handle(linkv1.NewLinkServiceHandler(h.link, interceptors))
	mux.Handle(resourcev1.NewResourceServiceHandler(h.resource, interceptors))
	mux.Handle(runtimev1.NewRuntimeServiceHandler(h.runtime, interceptors))

	// Prometheus reverse proxy. Requests arrive as
	// /proxy/{cluster}/prometheus/api/v1/query?... and are
	// forwarded through the tunnel to the agent's
	// /__otterscale/proxy/ endpoint. OIDC middleware protects this
	// path (it is not in the public paths list).
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", h.proxy)

	return nil
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

	// /metrics is not in the server's public paths, so the OIDC
	// middleware guards it like any other route: a scrape must present
	// a valid bearer token. Prometheus can do that with an oauth2:
	// section in its scrape config, pointed at the same Keycloak
	// client — a plain unauthenticated scrape gets 401.
	//
	// Left protected deliberately: these metrics carry cluster names
	// and per-procedure call patterns across every managed cluster. If
	// a deployment needs open scraping, expose it on a separate
	// listener rather than adding this path to WithPublicPaths, which
	// would also open it to the internet-facing API port.
	mux.Handle("/metrics", promhttp.Handler())

	return nil
}
