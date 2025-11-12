package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/machine/metal"
)

type machineRepo struct {
	maas *MAAS
}

func NewMachineRepo(maas *MAAS) metal.MachineRepo {
	return &machineRepo{
		maas: maas,
	}
}

var _ metal.MachineRepo = (*machineRepo)(nil)

func (r *machineRepo) List(_ context.Context) ([]metal.Machine, error) {
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

func (r *machineRepo) Get(_ context.Context, id string) (*metal.Machine, error) {
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

func (r *machineRepo) Release(_ context.Context, id string, force bool) (*metal.Machine, error) {
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

func (r *machineRepo) PowerOff(_ context.Context, id, comment string) (*metal.Machine, error) {
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

func (r *machineRepo) Commission(_ context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool) (*metal.Machine, error) {
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

func (r *machineRepo) toMachine(m *entity.Machine) *metal.Machine {
	return &metal.Machine{
		ID: m.SystemID,
	}
}

func (r *machineRepo) toMachines(ms []entity.Machine) []metal.Machine {
	list := []metal.Machine{}
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
