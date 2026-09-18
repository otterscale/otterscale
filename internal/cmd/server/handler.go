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

// agentValuesPath prefixes the route serving agent values behind an opaque id.
// The trailing slash matters: it is registered as a public prefix, and without
// it a sibling path could be swept in with it.
const agentValuesPath = "/link/values/"

// Handler mounts the gRPC service handlers, interceptors, and operational
// endpoints (health, reflection, metrics) onto an HTTP mux.
type Handler struct {
	link        *handler.LinkService
	resource    *handler.ResourceService
	runtime     *handler.RuntimeService
	agentValues *handler.AgentValuesHandler
	proxy       *handler.ProxyHandler
}

func NewHandler(
	link *handler.LinkService,
	resource *handler.ResourceService,
	runtime *handler.RuntimeService,
	agentValues *handler.AgentValuesHandler,
	proxy *handler.ProxyHandler,
) *Handler {
	return &Handler{
		link:        link,
		resource:    resource,
		runtime:     runtime,
		agentValues: agentValues,
		proxy:       proxy,
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

	// Raw YAML for `curl <url> | helm install -f -`. The opaque id in the path
	// is the only credential, which is why server.go registers this prefix as
	// public.
	mux.HandleFunc("GET "+agentValuesPath+"{id}", h.handleAgentValues)

	// Requests arrive as /proxy/{cluster}/prometheus/api/v1/query?... and are
	// forwarded through the tunnel to the agent's /__otterscale/proxy/. The
	// path is absent from the public paths list, so OIDC protects it.
	mux.Handle("/proxy/{cluster}/prometheus/{path...}", h.proxy)

	return nil
}

// handleAgentValues serves the values file a URL stands for, so it can be
// piped straight into a Helm install. The body is a credential bundle, hence
// no-store — and hence the procedure that mints these URLs is not GET-able
// either, since a Connect handler cannot set that header.
func (h *Handler) handleAgentValues(w http.ResponseWriter, r *http.Request) {
	values, err := h.agentValues.Render(r.Context(), r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if _, err := w.Write([]byte(values)); err != nil { // #nosec G705
		http.Error(w, "failed to write agent values response", http.StatusInternalServerError)
	}
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
