package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/machine"
)

type machineRepo struct {
	maas *MAAS
}

func NewMachineRepo(maas *MAAS) machine.MachineRepo {
	return &machineRepo{
		maas: maas,
	}
}

var _ machine.MachineRepo = (*machineRepo)(nil)

func (r *machineRepo) List(_ context.Context) ([]machine.Machine, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.MachinesParams{}

	return client.Machines.Get(params)
}

func (r *machineRepo) Get(_ context.Context, id string) (*machine.Machine, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.Machine.Get(id)

}

func (r *machineRepo) Release(_ context.Context, id string, force bool) (*machine.Machine, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.MachineReleaseParams{
		Force: force,
	}

	return client.Machine.Release(id, params)
}

func (r *machineRepo) PowerOff(_ context.Context, id, comment string) (*machine.Machine, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.MachinePowerOffParams{
		Comment: comment,
	}

	return client.Machine.PowerOff(id, params)
}

func (r *machineRepo) Commission(_ context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool) (*machine.Machine, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.MachineCommissionParams{
		TestingScripts: "",
		EnableSSH:      r.boolToInt(enableSSH),
		SkipBMCConfig:  r.boolToInt(skipBMCConfig),
		SkipNetworking: r.boolToInt(skipNetworking),
		SkipStorage:    r.boolToInt(skipStorage),
	}

	return client.Machine.Commission(id, params)
}

func (r *machineRepo) boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
