package mux

import (
	"net/http"
	"net/http/httputil"

	"connectrpc.com/connect"

	appv1 "github.com/otterscale/otterscale/api/application/v1/pbconnect"
	configv1 "github.com/otterscale/otterscale/api/configuration/v1/pbconnect"
	envv1 "github.com/otterscale/otterscale/api/environment/v1/pbconnect"
	facilityv1 "github.com/otterscale/otterscale/api/facility/v1/pbconnect"
	llmv1 "github.com/otterscale/otterscale/api/large_language_model/v1/pbconnect"
	machinev1 "github.com/otterscale/otterscale/api/machine/v1/pbconnect"
	networkv1 "github.com/otterscale/otterscale/api/network/v1/pbconnect"
	orchv1 "github.com/otterscale/otterscale/api/orchestrator/v1/pbconnect"
	premiumv1 "github.com/otterscale/otterscale/api/premium/v1/pbconnect"
	scopev1 "github.com/otterscale/otterscale/api/scope/v1/pbconnect"
	storagev1 "github.com/otterscale/otterscale/api/storage/v1/pbconnect"
	vmv1 "github.com/otterscale/otterscale/api/virtual_machine/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/app"
)

type Serve struct {
	*http.ServeMux

	services serviceRegistry
}

type serviceRegistry struct {
	app      *app.ApplicationService
	config   *app.ConfigurationService
	env      *app.EnvironmentService
	facility *app.FacilityService
	llm      *app.LargeLanguageModelService
	machine  *app.MachineService
	network  *app.NetworkService
	orch     *app.OrchestratorService
	premium  *app.PremiumService
	storage  *app.StorageService
	scope    *app.ScopeService
	vm       *app.VirtualMachineService
}

func (s *Serve) serviceNames() []string {
	return []string{
		appv1.ApplicationServiceName,
		configv1.ConfigurationServiceName,
		envv1.EnvironmentServiceName,
		facilityv1.FacilityServiceName,
		llmv1.LargeLanguageModelServiceName,
		machinev1.MachineServiceName,
		networkv1.NetworkServiceName,
		orchv1.OrchestratorServiceName,
		premiumv1.PremiumServiceName,
		scopev1.ScopeServiceName,
		storagev1.StorageServiceName,
		vmv1.VirtualMachineServiceName,
	}
}

func (s *Serve) registerHandlers(opts []connect.HandlerOption) {
	s.Handle(appv1.NewApplicationServiceHandler(s.services.app, opts...))
	s.Handle(configv1.NewConfigurationServiceHandler(s.services.config, opts...))
	s.Handle(envv1.NewEnvironmentServiceHandler(s.services.env, opts...))
	s.Handle(facilityv1.NewFacilityServiceHandler(s.services.facility, opts...))
	s.Handle(llmv1.NewLargeLanguageModelServiceHandler(s.services.llm, opts...))
	s.Handle(machinev1.NewMachineServiceHandler(s.services.machine, opts...))
	s.Handle(premiumv1.NewPremiumServiceHandler(s.services.premium, opts...))
	s.Handle(networkv1.NewNetworkServiceHandler(s.services.network, opts...))
	s.Handle(orchv1.NewOrchestratorServiceHandler(s.services.orch, opts...))
	s.Handle(storagev1.NewStorageServiceHandler(s.services.storage, opts...))
	s.Handle(scopev1.NewScopeServiceHandler(s.services.scope, opts...))
	s.Handle(vmv1.NewVirtualMachineServiceHandler(s.services.vm, opts...))
}

func (s *Serve) registerProxy() {
	proxy := httputil.NewSingleHostReverseProxy(s.services.env.GetPrometheusURL())
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		return nil
	}
	s.Handle("/prometheus/", http.StripPrefix("/prometheus", proxy))
}

func (s *Serve) registerWebSocket() {
	s.HandleFunc(s.services.vm.WebSocketPathPrefix, s.services.vm.VNCHandler)
}

func NewServe(app *app.ApplicationService, config *app.ConfigurationService, env *app.EnvironmentService, facility *app.FacilityService, llm *app.LargeLanguageModelService, machine *app.MachineService, network *app.NetworkService, orch *app.OrchestratorService, premium *app.PremiumService, storage *app.StorageService, scope *app.ScopeService, vm *app.VirtualMachineService, opts []connect.HandlerOption) (*Serve, error) {
	serve := &Serve{
		ServeMux: &http.ServeMux{},
		services: serviceRegistry{
			app:      app,
			config:   config,
			env:      env,
			facility: facility,
			llm:      llm,
			machine:  machine,
			network:  network,
			orch:     orch,
			premium:  premium,
			storage:  storage,
			scope:    scope,
			vm:       vm,
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
