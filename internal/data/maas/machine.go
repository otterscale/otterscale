package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core"
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

func (r *machine) List(_ context.Context) ([]core.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	machines, err := client.Machines.Get(&entity.MachinesParams{})
	if err != nil {
		return nil, err
	}
	list := []core.Machine{}
	for i := range machines {
		list = append(list, core.Machine{
			Machine: &machines[i],
		})
	}
	return list, nil
}

func (r *machine) Get(ctx context.Context, systemID string) (*core.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	machine, err := client.Machine.Get(systemID)
	if err != nil {
		return nil, err
	}
	return &core.Machine{
		Machine: machine,
	}, nil
}

func (r *machine) Release(ctx context.Context, systemID string, params *entity.MachineReleaseParams) (*core.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	machine, err := client.Machine.Release(systemID, params)
	if err != nil {
		return nil, err
	}
	return &core.Machine{
		Machine: machine,
	}, nil
}

func (r *machine) PowerOff(_ context.Context, systemID string, params *entity.MachinePowerOffParams) (*core.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	machine, err := client.Machine.PowerOff(systemID, params)
	if err != nil {
		return nil, err
	}
	return &core.Machine{
		Machine: machine,
	}, nil
}

func (r *machine) Commission(_ context.Context, systemID string, params *entity.MachineCommissionParams) (*core.Machine, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	machine, err := client.Machine.Commission(systemID, params)
	if err != nil {
		return nil, err
	}
	return &core.Machine{
		Machine: machine,
	}, nil
}
