package service

import (
	"context"

	"github.com/juju/juju/rpc/params"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (s *NexusService) ListMachines(ctx context.Context) ([]model.Machine, error) {
	return s.machine.List(ctx)
}

func (s *NexusService) GetMachine(ctx context.Context, id string) (*model.Machine, error) {
	return s.machine.Get(ctx, id)
}

func (s *NexusService) AddMachines(ctx context.Context, uuid string, factors []model.MachineFactor) ([]string, error) {
	params := []model.MachineAddParams{}
	for _, factor := range factors {
		directive, err := s.maasToJujuMachine(ctx, uuid, factor.MachineID)
		if err != nil {
			return nil, err
		}
		params = append(params, model.MachineAddParams{
			Placement:   toPlacement(&factor, directive),
			Constraints: toConstraint(&factor),
		})
	}
	results, err := s.machineManager.AddMachines(ctx, uuid, params)
	if err != nil {
		return nil, err
	}
	machines := make([]string, len(results))
	for i, r := range results {
		machines[i] = r.Machine
	}
	return machines, nil
}

func (s *NexusService) CommissionMachine(ctx context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool) (*model.Machine, error) {
	params := &model.MachineCommissionParams{
		EnableSSH:      boolToInt(enableSSH),
		SkipBMCConfig:  boolToInt(skipBMCConfig),
		SkipNetworking: boolToInt(skipNetworking),
		SkipStorage:    boolToInt(skipStorage),
	}
	return s.machine.Commission(ctx, id, params)
}

func (s *NexusService) PowerOnMachine(ctx context.Context, id, comment string) (*model.Machine, error) {
	params := &model.MachinePowerOnParams{
		Comment: comment,
	}
	return s.machine.PowerOn(ctx, id, params)
}

func (s *NexusService) PowerOffMachine(ctx context.Context, id, comment string) (*model.Machine, error) {
	params := &model.MachinePowerOffParams{
		Comment: comment,
	}
	return s.machine.PowerOff(ctx, id, params)
}

func (s *NexusService) listJujuMachines(ctx context.Context, uuid string) (map[string]params.MachineStatus, error) {
	status, err := s.client.Status(ctx, uuid, []string{"machine", "*"})
	if err != nil {
		return nil, err
	}
	return status.Machines, nil
}

func (s *NexusService) maasToJujuMachine(ctx context.Context, uuid, machineID string) (string, error) {
	if machineID == "" {
		return "", nil
	}
	mss, err := s.listJujuMachines(ctx, uuid)
	if err != nil {
		return "", err
	}
	for i := range mss {
		if mss[i].InstanceId.String() == machineID {
			return mss[i].Id, nil
		}
	}
	return "", status.Errorf(codes.NotFound, "maas machine %q not found", machineID)
}

func toPlacement(f *model.MachineFactor, directive string) *model.Placement {
	pla := &model.Placement{}
	if f.LXD {
		pla.Scope = "lxd"
	} else if f.KVM {
		pla.Scope = "kvm"
	} else if f.Machine {
		pla.Scope = "#"
		pla.Directive = directive
	} else {
		return nil
	}
	return pla
}

func toConstraint(f *model.MachineFactor) model.Constraint {
	con := model.Constraint{}
	if f.Architecture != "" {
		con.Arch = &f.Architecture
	}
	if f.CPUCores > 0 {
		con.CpuCores = &f.CPUCores
	}
	if f.MemoryMB > 0 {
		con.Mem = &f.MemoryMB
	}
	if len(f.Tags) > 0 {
		con.Tags = &f.Tags
	}
	return con
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
