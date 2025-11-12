package metal

import (
	"github.com/otterscale/otterscale/internal/core/configuration"
	"github.com/otterscale/otterscale/internal/core/facility"
	"github.com/otterscale/otterscale/internal/core/facility/action"
	"github.com/otterscale/otterscale/internal/core/scope"
)

type MetalUseCase struct {
	event          EventRepo
	machine        MachineRepo
	machineManager MachineManagerRepo
	nodeDevice     NodeDeviceRepo

	action       action.ActionRepo
	facility     facility.FacilityRepo
	orchestrator configuration.OrchestratorRepo
	provisioner  configuration.ProvisionerRepo
	scope        scope.ScopeRepo
}

func NewMetalUseCase(event EventRepo, machine MachineRepo, machineManager MachineManagerRepo, nodeDevice NodeDeviceRepo, action action.ActionRepo, facility facility.FacilityRepo, orchestrator configuration.OrchestratorRepo, provisioner configuration.ProvisionerRepo, scope scope.ScopeRepo) *MetalUseCase {
	return &MetalUseCase{
		event:          event,
		machine:        machine,
		machineManager: machineManager,
		nodeDevice:     nodeDevice,
		action:         action,
		facility:       facility,
		orchestrator:   orchestrator,
		provisioner:    provisioner,
		scope:          scope,
	}
}
