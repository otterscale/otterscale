package service

import (
	"context"
	"slices"

	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
)

// StackService coordinates operations across multiple MAAS resources
type StackService struct {
	server            MAASServer
	packageRepository MAASPackageRepository
	fabric            MAASFabric
	vlan              MAASVLAN
	subnet            MAASSubnet
	ipRange           MAASIPRange
	bootResource      MAASBootResource
	client            JujuClient
	machine           MAASMachine
	jujuMachine       JujuMachine
	model             JujuModel
	modelConfig       JujuModelConfig
	application       JujuApplication
	action            JujuAction
}

// NewStackService creates a new instance of StackService
func NewStackService(
	server MAASServer,
	packageRepository MAASPackageRepository,
	fabric MAASFabric,
	vlan MAASVLAN,
	subnet MAASSubnet,
	ipRange MAASIPRange,
	bootResource MAASBootResource,
	machine MAASMachine,
	jujuMachine JujuMachine,
	client JujuClient,
	model JujuModel,
	modelConfig JujuModelConfig,
	application JujuApplication,
	action JujuAction,
) *StackService {
	return &StackService{
		server:            server,
		packageRepository: packageRepository,
		fabric:            fabric,
		vlan:              vlan,
		subnet:            subnet,
		ipRange:           ipRange,
		bootResource:      bootResource,
		machine:           machine,
		client:            client,
		jujuMachine:       jujuMachine,
		model:             model,
		modelConfig:       modelConfig,
		application:       application,
		action:            action,
	}
}

func (s *StackService) ListApplications(ctx context.Context, uuid string, filters ...string) (map[string]params.ApplicationStatus, error) {
	patterns := []string{"application", "*"}
	if len(filters) > 0 {
		patterns = append(patterns[:1], filters...)
	}
	status, err := s.client.Status(ctx, uuid, patterns)
	if err != nil {
		return nil, err
	}
	if len(filters) == 0 {
		return status.Applications, nil
	}
	ret := make(map[string]params.ApplicationStatus)
	for k := range status.Applications {
		if slices.Contains(filters, k) {
			ret[k] = status.Applications[k]
		}
	}
	return ret, nil
}

func (s *StackService) ListJujuMachines(ctx context.Context, uuid string, filters ...string) (map[string]params.MachineStatus, error) {
	patterns := []string{"machine", "*"}
	if len(filters) > 0 {
		patterns = append(patterns[:1], filters...)
	}
	status, err := s.client.Status(ctx, uuid, patterns)
	if err != nil {
		return nil, err
	}
	return status.Machines, nil
}

func (s *StackService) ListApplicationConfigs(ctx context.Context, uuid string, appStatuses map[string]params.ApplicationStatus) (map[string]map[string]any, error) {
	names := []string{}
	for name := range appStatuses {
		names = append(names, name)
	}
	configs, err := s.application.GetConfigs(ctx, uuid, names...)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *StackService) CreateApplication(ctx context.Context, uuid, charmName, appName, channel string, revision, number int, config map[string]string, constraint constraints.Value, placements []instance.Placement, trust bool) (map[string]params.ApplicationStatus, error) {
	if err := s.application.Create(ctx, uuid, charmName, appName, channel, revision, number, config, constraint, placements, trust); err != nil {
		return nil, err
	}
	return s.ListApplications(ctx, uuid, appName)
}

func (s *StackService) UpdateApplication(ctx context.Context, uuid, name string, config map[string]string) (map[string]params.ApplicationStatus, error) {
	if err := s.application.Update(ctx, uuid, name, config); err != nil {
		return nil, err
	}
	return s.ListApplications(ctx, uuid, name)
}

func (s *StackService) DeleteApplication(ctx context.Context, uuid, name string, destroyStorage, force bool) error {
	return s.application.Delete(ctx, uuid, name, destroyStorage, force)
}

func (s *StackService) ExposeApplication(ctx context.Context, uuid, name string, endpoints map[string]params.ExposedEndpoint) error {
	return s.application.Expose(ctx, uuid, name, endpoints)
}

func (s *StackService) AddApplicationsUnits(ctx context.Context, uuid, name string, number int, placements []instance.Placement) ([]params.MachineStatus, error) {
	ids, err := s.application.AddUnits(ctx, uuid, name, number, placements)
	if err != nil {
		return nil, err
	}
	ms, err := s.ListJujuMachines(ctx, uuid, ids...)
	if err != nil {
		return nil, err
	}
	ret := []params.MachineStatus{}
	for i := range ms {
		ret = append(ret, ms[i])
	}
	return ret, nil
}

func (s *StackService) ListActions(ctx context.Context, uuid, appName string) (map[string]action.ActionSpec, error) {
	return s.action.List(ctx, uuid, appName)
}
