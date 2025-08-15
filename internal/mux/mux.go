package mux

import (
	"net/http"
	"net/http/httputil"

	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"

	applicationv1 "github.com/openhdc/otterscale/api/application/v1/pbconnect"
	bistv1 "github.com/openhdc/otterscale/api/bist/v1/pbconnect"
	configurationv1 "github.com/openhdc/otterscale/api/configuration/v1/pbconnect"
	environmentv1 "github.com/openhdc/otterscale/api/environment/v1/pbconnect"
	essentialv1 "github.com/openhdc/otterscale/api/essential/v1/pbconnect"
	facilityv1 "github.com/openhdc/otterscale/api/facility/v1/pbconnect"
	machinev1 "github.com/openhdc/otterscale/api/machine/v1/pbconnect"
	networkv1 "github.com/openhdc/otterscale/api/network/v1/pbconnect"
	premiumv1 "github.com/openhdc/otterscale/api/premium/v1/pbconnect"
	scopev1 "github.com/openhdc/otterscale/api/scope/v1/pbconnect"
	storagev1 "github.com/openhdc/otterscale/api/storage/v1/pbconnect"
	tagv1 "github.com/openhdc/otterscale/api/tag/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/app"
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

func New(helper bool, app *app.ApplicationService, bist *app.BISTService, config *app.ConfigurationService, environment *app.EnvironmentService, facility *app.FacilityService, essential *app.EssentialService, machine *app.MachineService, network *app.NetworkService, premium *app.PremiumService, storage *app.StorageService, scope *app.ScopeService, tag *app.TagService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(applicationv1.NewApplicationServiceHandler(app))
	mux.Handle(bistv1.NewBISTServiceHandler(bist))
	mux.Handle(configurationv1.NewConfigurationServiceHandler(config))
	mux.Handle(environmentv1.NewEnvironmentServiceHandler(environment))
	mux.Handle(facilityv1.NewFacilityServiceHandler(facility))
	mux.Handle(essentialv1.NewEssentialServiceHandler(essential))
	mux.Handle(machinev1.NewMachineServiceHandler(machine))
	mux.Handle(premiumv1.NewPremiumServiceHandler(premium))
	mux.Handle(networkv1.NewNetworkServiceHandler(network))
	mux.Handle(storagev1.NewStorageServiceHandler(storage))
	mux.Handle(scopev1.NewScopeServiceHandler(scope))
	mux.Handle(tagv1.NewTagServiceHandler(tag))

	// prometheus proxy
	proxy := httputil.NewSingleHostReverseProxy(environment.GetPrometheusURL())
	mux.Handle("/prometheus/", http.StripPrefix("/prometheus", proxy))

	if helper {
		// health
		checker := grpchealth.NewStaticChecker(Services...)
		mux.Handle(grpchealth.NewHandler(checker))

		// reflect
		reflector := grpcreflect.NewStaticReflector(Services...)
		mux.Handle(grpcreflect.NewHandlerV1(reflector))
		mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	}
	return mux
}
