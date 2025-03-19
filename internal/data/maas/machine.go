package maas

import (
	"context"

	"github.com/openhdc/openhdc/internal/domain/model"
	"github.com/openhdc/openhdc/internal/domain/service"
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

func (r *machine) List(ctx context.Context) ([]*model.Machine, error) {
	ms, err := r.maas.Machines.Get(&model.MachinesParams{})
	if err != nil {
		return nil, err
	}

	ret := make([]*model.Machine, len(ms))
	for i := range ms {
		ret[i] = &ms[i]
	}
	return ret, nil
}

func (r *machine) PowerOn(_ context.Context, systemID string, params *model.MachinePowerOnParams) (*model.Machine, error) {
	return r.maas.Machine.PowerOn(systemID, params)
}

func (r *machine) PowerOff(_ context.Context, systemID string, params *model.MachinePowerOffParams) (*model.Machine, error) {
	return r.maas.Machine.PowerOff(systemID, params)
}

func (r *machine) Commission(_ context.Context, systemID string, params *model.MachineCommissionParams) (*model.Machine, error) {
	return r.maas.Machine.Commission(systemID, params)
}
