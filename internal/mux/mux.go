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
	largelanguagemodelv1 "github.com/otterscale/otterscale/api/large_language_model/v1/pbconnect"
	machinev1 "github.com/otterscale/otterscale/api/machine/v1/pbconnect"
	networkv1 "github.com/otterscale/otterscale/api/network/v1/pbconnect"
	premiumv1 "github.com/otterscale/otterscale/api/premium/v1/pbconnect"
	scopev1 "github.com/otterscale/otterscale/api/scope/v1/pbconnect"
	storagev1 "github.com/otterscale/otterscale/api/storage/v1/pbconnect"
	tagv1 "github.com/otterscale/otterscale/api/tag/v1/pbconnect"
	virtualmachinev1 "github.com/otterscale/otterscale/api/virtual_machine/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/app"
)

var services = []string{
	applicationv1.ApplicationServiceName,
	bistv1.BISTServiceName,
	configurationv1.ConfigurationServiceName,
	environmentv1.EnvironmentServiceName,
	facilityv1.FacilityServiceName,
	largelanguagemodelv1.LargeLanguageModelServiceName,
	essentialv1.EssentialServiceName,
	machinev1.MachineServiceName,
	networkv1.NetworkServiceName,
	premiumv1.PremiumServiceName,
	scopev1.ScopeServiceName,
	storagev1.StorageServiceName,
	tagv1.TagServiceName,
	virtualmachinev1.VirtualMachineServiceName,
}

func New(helper bool, app *app.ApplicationService, bist *app.BISTService, config *app.ConfigurationService, environment *app.EnvironmentService, facility *app.FacilityService, largelanguagemodel *app.LargeLanguageModelService, essential *app.EssentialService, machine *app.MachineService, network *app.NetworkService, premium *app.PremiumService, storage *app.StorageService, scope *app.ScopeService, tag *app.TagService, virtualmachine *app.VirtualMachineService) (*http.ServeMux, error) {
	// interceptor
	opts, err := newInterceptorOptions()
	if err != nil {
		return nil, err
	}

	// multiplexer
	mux := http.NewServeMux()
	mux.Handle(applicationv1.NewApplicationServiceHandler(app, opts...))
	mux.Handle(bistv1.NewBISTServiceHandler(bist, opts...))
	mux.Handle(configurationv1.NewConfigurationServiceHandler(config, opts...))
	mux.Handle(environmentv1.NewEnvironmentServiceHandler(environment, opts...))
	mux.Handle(facilityv1.NewFacilityServiceHandler(facility, opts...))
	mux.Handle(largelanguagemodelv1.NewLargeLanguageModelServiceHandler(largelanguagemodel, opts...))
	mux.Handle(essentialv1.NewEssentialServiceHandler(essential, opts...))
	mux.Handle(machinev1.NewMachineServiceHandler(machine, opts...))
	mux.Handle(premiumv1.NewPremiumServiceHandler(premium, opts...))
	mux.Handle(networkv1.NewNetworkServiceHandler(network, opts...))
	mux.Handle(storagev1.NewStorageServiceHandler(storage, opts...))
	mux.Handle(scopev1.NewScopeServiceHandler(scope, opts...))
	mux.Handle(tagv1.NewTagServiceHandler(tag, opts...))
	mux.Handle(virtualmachinev1.NewVirtualMachineServiceHandler(virtualmachine, opts...))

	// proxy
	handlePrometheusProxy(mux, environment)

	// helper
	if helper {
		handleHealth(mux)
		handleReflect(mux)
		if err := handleMetrics(mux); err != nil {
			return nil, err
		}
	}

	return mux, nil
}

func newInterceptorOptions() ([]connect.HandlerOption, error) {
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, err
	}
	return []connect.HandlerOption{
		connect.WithInterceptors(otelInterceptor),
	}, nil
}

func handlePrometheusProxy(mux *http.ServeMux, environment *app.EnvironmentService) {
	proxy := httputil.NewSingleHostReverseProxy(environment.GetPrometheusURL())
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		return nil
	}
	mux.Handle("/prometheus/", http.StripPrefix("/prometheus", proxy))
}

func handleHealth(mux *http.ServeMux) {
	checker := grpchealth.NewStaticChecker(services...)
	mux.Handle(grpchealth.NewHandler(checker))
}

func handleReflect(mux *http.ServeMux) {
	reflector := grpcreflect.NewStaticReflector(services...)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}

func handleMetrics(mux *http.ServeMux) error {
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}
	otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(exporter)))
	mux.Handle("/metrics", promhttp.Handler())
	return nil
}
