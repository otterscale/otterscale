package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type machine struct {
	maas *MAAS
}

func NewMachine(maas *MAAS) service.MAASMachine {
	return &machine{
		maas: maas,
	}
}

var _ service.MAASMachine = (*machine)(nil)

func (r *machine) List(_ context.Context) ([]entity.Machine, error) {
	return r.maas.Machines.Get(&entity.MachinesParams{})
}

func (r *machine) Get(ctx context.Context, systemID string) (*entity.Machine, error) {
	return r.maas.Machine.Get(systemID)
}

func (r *machine) Release(ctx context.Context, systemID string, force bool) (*entity.Machine, error) {
	params := &entity.MachineReleaseParams{
		Force: force,
	}
	return r.maas.Machine.Release(systemID, params)
}

func (r *machine) PowerOn(_ context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error) {
	return r.maas.Machine.PowerOn(systemID, params)
}

func (r *machine) PowerOff(_ context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error) {
	return r.maas.Machine.PowerOff(systemID, params)
}

func (r *machine) Commission(_ context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error) {
	return r.maas.Machine.Commission(systemID, params)
}
