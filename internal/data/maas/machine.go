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
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machines.Get(&entity.MachinesParams{})
}

func (r *machine) Get(ctx context.Context, systemID string) (*entity.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machine.Get(systemID)
}

func (r *machine) Release(ctx context.Context, systemID string, params *entity.MachineReleaseParams) (*entity.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machine.Release(systemID, params)
}

func (r *machine) PowerOn(_ context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machine.PowerOn(systemID, params)
}

func (r *machine) PowerOff(_ context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machine.PowerOff(systemID, params)
}

func (r *machine) Commission(_ context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machine.Commission(systemID, params)
}
