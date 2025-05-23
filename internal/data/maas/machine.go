package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/otterscale/internal/core"
)

type machine struct {
	maas *MAAS
}

func NewMachine(maas *MAAS) core.MachineRepo {
	return &machine{
		maas: maas,
	}
}

var _ core.MachineRepo = (*machine)(nil)

func (r *machine) List(_ context.Context) ([]entity.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machines.Get(&entity.MachinesParams{})
}

func (r *machine) Get(ctx context.Context, systemID string) (*core.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machine.Get(systemID)
}

func (r *machine) Release(ctx context.Context, systemID string, params *entity.MachineReleaseParams) (*core.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machine.Release(systemID, params)
}

func (r *machine) PowerOff(_ context.Context, systemID string, params *entity.MachinePowerOffParams) (*core.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machine.PowerOff(systemID, params)
}

func (r *machine) Commission(_ context.Context, systemID string, params *entity.MachineCommissionParams) (*core.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.Machine.Commission(systemID, params)
}
