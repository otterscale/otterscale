package mux

import (
	"net/http"
	"net/http/httputil"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	appv1 "github.com/otterscale/otterscale/api/application/v1/pbconnect"
	configv1 "github.com/otterscale/otterscale/api/configuration/v1/pbconnect"
	envv1 "github.com/otterscale/otterscale/api/environment/v1/pbconnect"
	facilityv1 "github.com/otterscale/otterscale/api/facility/v1/pbconnect"
	instancev1 "github.com/otterscale/otterscale/api/instance/v1/pbconnect"
	k8sv1 "github.com/otterscale/otterscale/api/kubernetes/v1/pbconnect"
	machinev1 "github.com/otterscale/otterscale/api/machine/v1/pbconnect"
	modelv1 "github.com/otterscale/otterscale/api/model/v1/pbconnect"
	networkv1 "github.com/otterscale/otterscale/api/network/v1/pbconnect"
	orchv1 "github.com/otterscale/otterscale/api/orchestrator/v1/pbconnect"
	registryv1 "github.com/otterscale/otterscale/api/registry/v1/pbconnect"
	resourcev1 "github.com/otterscale/otterscale/api/resource/v1/pbconnect"
	scopev1 "github.com/otterscale/otterscale/api/scope/v1/pbconnect"
	storagev1 "github.com/otterscale/otterscale/api/storage/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/app"
)

type Serve struct {
	*http.ServeMux

	application   *app.ApplicationService
	configuration *app.ConfigurationService
	environment   *app.EnvironmentService
	facility      *app.FacilityService
	instance      *app.InstanceService
	kubernetes    *app.KubernetesService
	machine       *app.MachineService
	model         *app.ModelService
	network       *app.NetworkService
	orchestrator  *app.OrchestratorService
	registry      *app.RegistryService
	resource      *app.ResourceService
	storage       *app.StorageService
	scope         *app.ScopeService
}

func NewServe(application *app.ApplicationService, configuration *app.ConfigurationService, environment *app.EnvironmentService, facility *app.FacilityService, instance *app.InstanceService, kubernetes *app.KubernetesService, machine *app.MachineService, model *app.ModelService, network *app.NetworkService, orchestrator *app.OrchestratorService, registry *app.RegistryService, resource *app.ResourceService, storage *app.StorageService, scope *app.ScopeService) *Serve {
	return &Serve{
		ServeMux:      &http.ServeMux{},
		application:   application,
		configuration: configuration,
		environment:   environment,
		facility:      facility,
		instance:      instance,
		kubernetes:    kubernetes,
		machine:       machine,
		model:         model,
		network:       network,
		orchestrator:  orchestrator,
		registry:      registry,
		resource:      resource,
		storage:       storage,
		scope:         scope,
	}
}

func (s *Serve) RegisterHandlers(opts []connect.HandlerOption) error {
	// Prepare service names for reflection and health check
	services := []string{
		appv1.ApplicationServiceName,
		configv1.ConfigurationServiceName,
		envv1.EnvironmentServiceName,
		facilityv1.FacilityServiceName,
		instancev1.InstanceServiceName,
		machinev1.MachineServiceName,
		k8sv1.KubernetesServiceName,
		modelv1.ModelServiceName,
		networkv1.NetworkServiceName,
		orchv1.OrchestratorServiceName,
		registryv1.RegistryServiceName,
		resourcev1.ResourceServiceName,
		scopev1.ScopeServiceName,
		storagev1.StorageServiceName,
	}

	// Register gRPC reflection
	reflector := grpcreflect.NewStaticReflector(services...)
	s.Handle(grpcreflect.NewHandlerV1(reflector))
	s.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Register gRPC health check
	checker := grpchealth.NewStaticChecker(services...)
	s.Handle(grpchealth.NewHandler(checker))

	// Register metrics endpoint
	if err := s.registerMetrics(); err != nil {
		return err
	}

	// Register Prometheus proxy
	s.registerProxy()

	// Register WebSocket handler
	s.registerWebSocket()

	// Register service handlers
	s.Handle(appv1.NewApplicationServiceHandler(s.application, opts...))
	s.Handle(configv1.NewConfigurationServiceHandler(s.configuration, opts...))
	s.Handle(envv1.NewEnvironmentServiceHandler(s.environment, opts...))
	s.Handle(facilityv1.NewFacilityServiceHandler(s.facility, opts...))
	s.Handle(instancev1.NewInstanceServiceHandler(s.instance, opts...))
	s.Handle(k8sv1.NewKubernetesServiceHandler(s.kubernetes, opts...))
	s.Handle(machinev1.NewMachineServiceHandler(s.machine, opts...))
	s.Handle(modelv1.NewModelServiceHandler(s.model, opts...))
	s.Handle(networkv1.NewNetworkServiceHandler(s.network, opts...))
	s.Handle(orchv1.NewOrchestratorServiceHandler(s.orchestrator, opts...))
	s.Handle(registryv1.NewRegistryServiceHandler(s.registry, opts...))
	s.Handle(resourcev1.NewResourceServiceHandler(s.resource, opts...))
	s.Handle(storagev1.NewStorageServiceHandler(s.storage, opts...))
	s.Handle(scopev1.NewScopeServiceHandler(s.scope, opts...))

	return nil
}

func (s *Serve) registerMetrics() error {
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}
	otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(exporter)))
	s.Handle("/metrics", promhttp.Handler())
	return nil
}

func (s *Serve) registerProxy() {
	proxy := httputil.NewSingleHostReverseProxy(s.environment.GetPrometheusURL())
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		return nil
	}
	s.Handle("/prometheus/", http.StripPrefix("/prometheus", proxy))
}

func (s *Serve) registerWebSocket() {
	s.HandleFunc(s.instance.VNCPathPrefix(), s.instance.VNCHandler())
}
