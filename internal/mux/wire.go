package mux

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"

	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

var ProviderSet = wire.NewSet(NewBootstrap, NewServe, NewInterceptorOptions)

func NewInterceptorOptions() ([]connect.HandlerOption, error) {
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, err
	}
	return []connect.HandlerOption{
		connect.WithInterceptors(otelInterceptor),
	}, nil
}

func registerPrometheusMetrics(mux *http.ServeMux) error {
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}
	otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(exporter)))
	mux.Handle("/metrics", promhttp.Handler())
	return nil
}

func registerGRPCReflection(mux *http.ServeMux, services ...string) {
	reflector := grpcreflect.NewStaticReflector(services...)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}

func registerGRPCHealth(mux *http.ServeMux, services ...string) {
	checker := grpchealth.NewStaticChecker(services...)
	mux.Handle(grpchealth.NewHandler(checker))
}
