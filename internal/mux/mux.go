package mux

import (
	"net/http"
	"net/http/httputil"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	applicationv1 "github.com/otterscale/otterscale/api/application/v1/pbconnect"
	bistv1 "github.com/otterscale/otterscale/api/bist/v1/pbconnect"
	configurationv1 "github.com/otterscale/otterscale/api/configuration/v1/pbconnect"
	environmentv1 "github.com/otterscale/otterscale/api/environment/v1/pbconnect"
	essentialv1 "github.com/otterscale/otterscale/api/essential/v1/pbconnect"
	facilityv1 "github.com/otterscale/otterscale/api/facility/v1/pbconnect"
	machinev1 "github.com/otterscale/otterscale/api/machine/v1/pbconnect"
	networkv1 "github.com/otterscale/otterscale/api/network/v1/pbconnect"
	premiumv1 "github.com/otterscale/otterscale/api/premium/v1/pbconnect"
	scopev1 "github.com/otterscale/otterscale/api/scope/v1/pbconnect"
	storagev1 "github.com/otterscale/otterscale/api/storage/v1/pbconnect"
	tagv1 "github.com/otterscale/otterscale/api/tag/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/app"
)

var Services = []string{
	applicationv1.ApplicationServiceName,
	bistv1.BISTServiceName,
	configurationv1.ConfigurationServiceName,
	environmentv1.EnvironmentServiceName,
	facilityv1.FacilityServiceName,
	essentialv1.EssentialServiceName,
	machinev1.MachineServiceName,
	networkv1.NetworkServiceName,
	premiumv1.PremiumServiceName,
	scopev1.ScopeServiceName,
	storagev1.StorageServiceName,
	tagv1.TagServiceName,
}

func New(helper bool, app *app.ApplicationService, bist *app.BISTService, config *app.ConfigurationService, environment *app.EnvironmentService, facility *app.FacilityService, essential *app.EssentialService, machine *app.MachineService, network *app.NetworkService, premium *app.PremiumService, storage *app.StorageService, scope *app.ScopeService, tag *app.TagService) (*http.ServeMux, error) {
	// interceptor
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, err
	}
	opts := []connect.HandlerOption{
		connect.WithInterceptors(otelInterceptor),
	}

	// mux
	mux := http.NewServeMux()
	mux.Handle(applicationv1.NewApplicationServiceHandler(app, opts...))
	mux.Handle(bistv1.NewBISTServiceHandler(bist, opts...))
	mux.Handle(configurationv1.NewConfigurationServiceHandler(config, opts...))
	mux.Handle(environmentv1.NewEnvironmentServiceHandler(environment, opts...))
	mux.Handle(facilityv1.NewFacilityServiceHandler(facility, opts...))
	mux.Handle(essentialv1.NewEssentialServiceHandler(essential, opts...))
	mux.Handle(machinev1.NewMachineServiceHandler(machine, opts...))
	mux.Handle(premiumv1.NewPremiumServiceHandler(premium, opts...))
	mux.Handle(networkv1.NewNetworkServiceHandler(network, opts...))
	mux.Handle(storagev1.NewStorageServiceHandler(storage, opts...))
	mux.Handle(scopev1.NewScopeServiceHandler(scope, opts...))
	mux.Handle(tagv1.NewTagServiceHandler(tag, opts...))

	// prometheus proxy
	proxy := httputil.NewSingleHostReverseProxy(environment.GetPrometheusURL())
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		return nil
	}

	mux.Handle("/prometheus/", http.StripPrefix("/prometheus", proxy))

	if helper {
		// health
		checker := grpchealth.NewStaticChecker(Services...)
		mux.Handle(grpchealth.NewHandler(checker))

		// reflect
		reflector := grpcreflect.NewStaticReflector(Services...)
		mux.Handle(grpcreflect.NewHandlerV1(reflector))
		mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

		// metrics
		exporter, err := prometheus.New()
		if err != nil {
			return nil, err
		}
		otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(exporter)))
		mux.Handle("/metrics", promhttp.Handler())
	}
	return mux, nil
}
