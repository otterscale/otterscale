package mux

import (
	"net/http"
	"net/http/httputil"

	"connectrpc.com/connect"

	appv1 "github.com/otterscale/otterscale/api/application/v1/pbconnect"
	configv1 "github.com/otterscale/otterscale/api/configuration/v1/pbconnect"
	envv1 "github.com/otterscale/otterscale/api/environment/v1/pbconnect"
	facilityv1 "github.com/otterscale/otterscale/api/facility/v1/pbconnect"
	instancev1 "github.com/otterscale/otterscale/api/instance/v1/pbconnect"
	machinev1 "github.com/otterscale/otterscale/api/machine/v1/pbconnect"
	modelv1 "github.com/otterscale/otterscale/api/model/v1/pbconnect"
	networkv1 "github.com/otterscale/otterscale/api/network/v1/pbconnect"
	orchv1 "github.com/otterscale/otterscale/api/orchestrator/v1/pbconnect"
	scopev1 "github.com/otterscale/otterscale/api/scope/v1/pbconnect"
	storagev1 "github.com/otterscale/otterscale/api/storage/v1/pbconnect"
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
	instance *app.InstanceService
	machine  *app.MachineService
	model    *app.ModelService
	network  *app.NetworkService
	orch     *app.OrchestratorService
	storage  *app.StorageService
	scope    *app.ScopeService
}

func (s *Serve) serviceNames() []string {
	return []string{
		appv1.ApplicationServiceName,
		configv1.ConfigurationServiceName,
		envv1.EnvironmentServiceName,
		facilityv1.FacilityServiceName,
		instancev1.InstanceServiceName,
		machinev1.MachineServiceName,
		modelv1.ModelServiceName,
		networkv1.NetworkServiceName,
		orchv1.OrchestratorServiceName,
		scopev1.ScopeServiceName,
		storagev1.StorageServiceName,
	}
}

func (s *Serve) registerHandlers(opts []connect.HandlerOption) {
	s.Handle(appv1.NewApplicationServiceHandler(s.services.app, opts...))
	s.Handle(configv1.NewConfigurationServiceHandler(s.services.config, opts...))
	s.Handle(envv1.NewEnvironmentServiceHandler(s.services.env, opts...))
	s.Handle(facilityv1.NewFacilityServiceHandler(s.services.facility, opts...))
	s.Handle(instancev1.NewInstanceServiceHandler(s.services.instance, opts...))
	s.Handle(machinev1.NewMachineServiceHandler(s.services.machine, opts...))
	s.Handle(modelv1.NewModelServiceHandler(s.services.model, opts...))
	s.Handle(networkv1.NewNetworkServiceHandler(s.services.network, opts...))
	s.Handle(orchv1.NewOrchestratorServiceHandler(s.services.orch, opts...))
	s.Handle(storagev1.NewStorageServiceHandler(s.services.storage, opts...))
	s.Handle(scopev1.NewScopeServiceHandler(s.services.scope, opts...))
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
	s.HandleFunc(s.services.instance.VNCPathPrefix(), s.services.instance.VNCHandler())
}

func NewServe(app *app.ApplicationService, config *app.ConfigurationService, env *app.EnvironmentService, facility *app.FacilityService, instance *app.InstanceService, machine *app.MachineService, model *app.ModelService, network *app.NetworkService, orch *app.OrchestratorService, storage *app.StorageService, scope *app.ScopeService, opts []connect.HandlerOption) (*Serve, error) {
	// Initialize ServeMux and register all handlers
	serve := &Serve{
		ServeMux: &http.ServeMux{},
		services: serviceRegistry{
			app:      app,
			config:   config,
			env:      env,
			facility: facility,
			instance: instance,
			machine:  machine,
			model:    model,
			network:  network,
			orch:     orch,
			storage:  storage,
			scope:    scope,
		},
	}
	serve.registerHandlers(opts)
	serve.registerProxy()
	serve.registerWebSocket()

	// Register gRPC reflection, health check, and Prometheus metrics
	registerGRPCReflection(serve.ServeMux, serve.serviceNames()...)
	registerGRPCHealth(serve.ServeMux, serve.serviceNames()...)
	if err := registerPrometheusMetrics(serve.ServeMux); err != nil {
		return nil, err
	}

	return serve, nil
}
