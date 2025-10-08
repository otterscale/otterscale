package mux

import (
	"net/http"
	"net/http/httputil"

	"connectrpc.com/connect"
	appv1 "github.com/otterscale/otterscale/api/application/v1/pbconnect"
	configv1 "github.com/otterscale/otterscale/api/configuration/v1/pbconnect"
	envv1 "github.com/otterscale/otterscale/api/environment/v1/pbconnect"
	essentialv1 "github.com/otterscale/otterscale/api/essential/v1/pbconnect"
	facilityv1 "github.com/otterscale/otterscale/api/facility/v1/pbconnect"
	llmv1 "github.com/otterscale/otterscale/api/large_language_model/v1/pbconnect"
	machinev1 "github.com/otterscale/otterscale/api/machine/v1/pbconnect"
	networkv1 "github.com/otterscale/otterscale/api/network/v1/pbconnect"
	premiumv1 "github.com/otterscale/otterscale/api/premium/v1/pbconnect"
	scopev1 "github.com/otterscale/otterscale/api/scope/v1/pbconnect"
	storagev1 "github.com/otterscale/otterscale/api/storage/v1/pbconnect"
	tagv1 "github.com/otterscale/otterscale/api/tag/v1/pbconnect"
	vmv1 "github.com/otterscale/otterscale/api/virtual_machine/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/app"
)

type Serve struct {
	*http.ServeMux

	services serviceRegistry
}

type serviceRegistry struct {
	app       *app.ApplicationService
	config    *app.ConfigurationService
	env       *app.EnvironmentService
	facility  *app.FacilityService
	llm       *app.LargeLanguageModelService
	essential *app.EssentialService
	machine   *app.MachineService
	network   *app.NetworkService
	premium   *app.PremiumService
	storage   *app.StorageService
	scope     *app.ScopeService
	tag       *app.TagService
	vm        *app.VirtualMachineService
}

func (s *Serve) serviceNames() []string {
	return []string{
		appv1.ApplicationServiceName,
		configv1.ConfigurationServiceName,
		envv1.EnvironmentServiceName,
		facilityv1.FacilityServiceName,
		llmv1.LargeLanguageModelServiceName,
		essentialv1.EssentialServiceName,
		machinev1.MachineServiceName,
		networkv1.NetworkServiceName,
		premiumv1.PremiumServiceName,
		scopev1.ScopeServiceName,
		storagev1.StorageServiceName,
		tagv1.TagServiceName,
		vmv1.VirtualMachineServiceName,
	}
}

func (s *Serve) registerHandlers(opts []connect.HandlerOption) {
	s.ServeMux.Handle(appv1.NewApplicationServiceHandler(s.services.app, opts...))
	s.ServeMux.Handle(configv1.NewConfigurationServiceHandler(s.services.config, opts...))
	s.ServeMux.Handle(envv1.NewEnvironmentServiceHandler(s.services.env, opts...))
	s.ServeMux.Handle(facilityv1.NewFacilityServiceHandler(s.services.facility, opts...))
	s.ServeMux.Handle(llmv1.NewLargeLanguageModelServiceHandler(s.services.llm, opts...))
	s.ServeMux.Handle(essentialv1.NewEssentialServiceHandler(s.services.essential, opts...))
	s.ServeMux.Handle(machinev1.NewMachineServiceHandler(s.services.machine, opts...))
	s.ServeMux.Handle(premiumv1.NewPremiumServiceHandler(s.services.premium, opts...))
	s.ServeMux.Handle(networkv1.NewNetworkServiceHandler(s.services.network, opts...))
	s.ServeMux.Handle(storagev1.NewStorageServiceHandler(s.services.storage, opts...))
	s.ServeMux.Handle(scopev1.NewScopeServiceHandler(s.services.scope, opts...))
	s.ServeMux.Handle(tagv1.NewTagServiceHandler(s.services.tag, opts...))
	s.ServeMux.Handle(vmv1.NewVirtualMachineServiceHandler(s.services.vm, opts...))
}

func (s *Serve) registerProxy() {
	proxy := httputil.NewSingleHostReverseProxy(s.services.env.GetPrometheusURL())
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		return nil
	}
	s.ServeMux.Handle("/prometheus/", http.StripPrefix("/prometheus", proxy))
}

func (s *Serve) registerWebSocket() {
	s.ServeMux.HandleFunc(s.services.vm.WebSocketPathPrefix, s.services.vm.VNCHandler)
}

func NewServe(
	app *app.ApplicationService,
	config *app.ConfigurationService,
	env *app.EnvironmentService,
	facility *app.FacilityService,
	llm *app.LargeLanguageModelService,
	essential *app.EssentialService,
	machine *app.MachineService,
	network *app.NetworkService,
	premium *app.PremiumService,
	storage *app.StorageService,
	scope *app.ScopeService,
	tag *app.TagService,
	vm *app.VirtualMachineService,
	opts []connect.HandlerOption) (*Serve, error) {
	serve := &Serve{
		ServeMux: &http.ServeMux{},
		services: serviceRegistry{
			app:       app,
			config:    config,
			env:       env,
			facility:  facility,
			llm:       llm,
			essential: essential,
			machine:   machine,
			network:   network,
			premium:   premium,
			storage:   storage,
			scope:     scope,
			tag:       tag,
			vm:        vm,
		},
	}
	serve.registerHandlers(opts)
	serve.registerProxy()
	serve.registerWebSocket()

	handleGRPCCompatible(serve.ServeMux, serve.serviceNames()...)

	if err := handleMetrics(serve.ServeMux); err != nil {
		return nil, err
	}
	return serve, nil
}
