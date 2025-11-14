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

	machines, err := client.Machines.Get(params)
	if err != nil {
		return nil, err
	}

	return r.toMachines(machines), nil
}

func (r *machineRepo) Get(_ context.Context, id string) (*machine.Machine, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	machine, err := client.Machine.Get(id)
	if err != nil {
		return nil, err
	}

	return r.toMachine(machine), nil

}

func (r *machineRepo) Release(_ context.Context, id string, force bool) (*machine.Machine, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.MachineReleaseParams{
		Force: force,
	}

	machine, err := client.Machine.Release(id, params)
	if err != nil {
		return nil, err
	}

	return r.toMachine(machine), nil
}

func (r *machineRepo) PowerOff(_ context.Context, id, comment string) (*machine.Machine, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.MachinePowerOffParams{
		Comment: comment,
	}

	machine, err := client.Machine.PowerOff(id, params)
	if err != nil {
		return nil, err
	}

	return r.toMachine(machine), nil
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

	machine, err := client.Machine.Commission(id, params)
	if err != nil {
		return nil, err
	}

	return r.toMachine(machine), nil
}

func (r *machineRepo) toMachine(m *entity.Machine) *machine.Machine {
	return &machine.Machine{
		Machine: m,
	}
}

func (r *machineRepo) toMachines(ms []entity.Machine) []machine.Machine {
	list := []machine.Machine{}
	for _, m := range ms {
		list = append(list, *r.toMachine(&m))
	}
	return list
}

func (r *machineRepo) boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
