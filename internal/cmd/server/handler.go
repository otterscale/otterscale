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

// Handler mounts the gRPC service handlers, interceptors, and operational
// endpoints (health, reflection, metrics) onto an HTTP mux.
type Handler struct {
	link     *handler.LinkService
	resource *handler.ResourceService
	runtime  *handler.RuntimeService
	proxy    *handler.ProxyHandler
}

func NewHandler(link *handler.LinkService, resource *handler.ResourceService, runtime *handler.RuntimeService, proxy *handler.ProxyHandler) *Handler {
	return &Handler{
		link:     link,
		resource: resource,
		runtime:  runtime,
		proxy:    proxy,
	}
}

// LongRunningPaths names the procedures whose response is a long-lived stream.
// The transport lifts its request timeouts for these, so a watch or an exec
// session is not cut off mid-flight.
func (h *Handler) LongRunningPaths() []string {
	return []string{
		resourcev1.ResourceServiceWatchProcedure,
		runtimev1.RuntimeServicePodLogProcedure,
		runtimev1.RuntimeServiceExecuteTTYProcedure,
		runtimev1.RuntimeServicePortForwardProcedure,
		runtimev1.RuntimeServiceVNCProcedure,
	}
}

func (h *Handler) Mount(mux *http.ServeMux) error {
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return err
	}

	interceptors := connect.WithInterceptors(
		otelInterceptor,
	)

	services := []string{
		linkv1.LinkServiceName,
		resourcev1.ResourceServiceName,
		runtimev1.RuntimeServiceName,
	}

	if err := h.registerOpsHandlers(mux, services); err != nil {
		return err
	}

	// RPCs with idempotency_level = NO_SIDE_EFFECTS accept HTTP GET through the
	// generated connect.WithIdempotency option.
	mux.Handle(linkv1.NewLinkServiceHandler(h.link, interceptors))
	mux.Handle(resourcev1.NewResourceServiceHandler(h.resource, interceptors))
	mux.Handle(runtimev1.NewRuntimeServiceHandler(h.runtime, interceptors))

	// Requests arrive as /proxy/{cluster}/prometheus/api/v1/query?... and are
	// forwarded through the tunnel to the agent's /__otterscale/proxy/. The
	// path is absent from the public paths list, so OIDC protects it.
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", h.proxy)

	return nil
}

// registerOpsHandlers sets up reflection, health checks, and metrics scraping.
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
	// Set globally, not injected via Wire: otelconnect and other libraries
	// discover the MeterProvider only through the global.
	otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(exporter)))

	// /metrics is absent from the public paths, so OIDC guards it like any
	// other route: a scrape must present a bearer token, which Prometheus can
	// do with an oauth2: section pointed at the same Keycloak client.
	//
	// Protected deliberately — these metrics carry cluster names and
	// per-procedure call patterns across every managed cluster. A deployment
	// needing open scraping should expose it on a separate listener, not add
	// this path to WithPublicPaths, which would also open it on the
	// internet-facing API port.
	mux.Handle("/metrics", promhttp.Handler())

	return nil
}
