package facility

import (
	"context"

	"github.com/juju/juju/rpc/params"

	"github.com/otterscale/otterscale/internal/core/machine"
)

type (
	// Status represents a Juju ApplicationStatus resource.
	Status = params.ApplicationStatus

	// UnitStatus represents a Juju UnitStatus resource.
	UnitStatus = params.UnitStatus

	// DetailedStatus represents a Juju DetailedStatus resource.
	DetailedStatus = params.DetailedStatus

	// MachineStatus represents a Juju MachineStatus resource.
	MachineStatus = params.MachineStatus
)

type Facility struct {
	Name   string
	Status *Status
}

//nolint:revive // allows this exported interface name for specific domain clarity.
type FacilityRepo interface {
	List(ctx context.Context, scope, jujuID string) ([]Facility, error)
	Create(ctx context.Context, scope, name, configYAML, charmName, channel, placementScope string, subordinate bool, directive, series string) error
	Update(ctx context.Context, scope, name, configYAML string) error
	Delete(ctx context.Context, scope, name string, destroyStorage, force bool) error
	Resolve(ctx context.Context, scope, unitName string) error
	Config(ctx context.Context, scope string, name string) (map[string]any, error)
}

type RelationRepo interface {
	Create(ctx context.Context, scope string, endpoints []string) error
	Delete(ctx context.Context, scope string, id int) error
	Consume(ctx context.Context, scope, url string) error
}

type UseCase struct {
	facility FacilityRepo

	machine machine.MachineRepo
}

func NewUseCase(facility FacilityRepo, machine machine.MachineRepo) *UseCase {
	return &UseCase{
		facility: facility,
		machine:  machine,
	}
}

func (uc *UseCase) ListFacilities(ctx context.Context, scope string) ([]Facility, error) {
	return uc.facility.List(ctx, scope, "")
}

func (uc *UseCase) ResolveFacilityUnitErrors(ctx context.Context, scope, unitName string) error {
	return uc.facility.Resolve(ctx, scope, unitName)
}

func (uc *UseCase) JujuIDMachineMap(ctx context.Context, scope string) (map[string]machine.Machine, error) {
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}

	ret := map[string]machine.Machine{}

	for i := range machines {
		machineScope, err := uc.machine.ExtractScope(&machines[i])
		if err != nil {
			continue
		}

		if machineScope != scope {
			continue
		}

		machineJujuID, err := uc.machine.ExtractJujuID(&machines[i])
		if err != nil {
			continue
		}

		ret[machineJujuID] = machines[i]
	}

	return ret, nil
}
