package service

import (
	"context"

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
		directives, err := s.maasToJujuMachineMap(ctx, uuid)
		if err != nil {
			return nil, err
		}
		params = append(params, model.MachineAddParams{
			Placement:   toPlacement(factor.MachinePlacement, directives[factor.MachineID]),
			Constraints: toConstraint(factor.MachineConstraint),
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

func (s *NexusService) listJujuMachines(ctx context.Context, uuid string) (map[string]model.MachineStatus, error) {
	status, err := s.client.Status(ctx, uuid, []string{"machine", "*"})
	if err != nil {
		return nil, err
	}
	return status.Machines, nil
}

func (s *NexusService) JujuToMAASMachineMap(ctx context.Context, uuid string) (map[string]string, error) {
	msm, err := s.listJujuMachines(ctx, uuid)
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	for key := range msm {
		m[key] = msm[key].InstanceId.String()
	}
	return m, nil
}

func (s *NexusService) maasToJujuMachineMap(ctx context.Context, uuid string) (map[string]string, error) {
	msm, err := s.listJujuMachines(ctx, uuid)
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	for key := range msm {
		m[msm[key].InstanceId.String()] = key
	}
	return m, nil
}

func toPlacement(p *model.MachinePlacement, directive string) *model.Placement {
	pla := &model.Placement{}
	if p.LXD {
		pla.Scope = "lxd"
	} else if p.KVM {
		pla.Scope = "kvm"
	} else if p.Machine {
		pla.Scope = "#"
		pla.Directive = directive
	} else {
		return nil
	}
	return pla
}

func toConstraint(c *model.MachineConstraint) model.Constraint {
	con := model.Constraint{}
	if c.Architecture != "" {
		con.Arch = &c.Architecture
	}
	if c.CPUCores > 0 {
		con.CpuCores = &c.CPUCores
	}
	if c.MemoryMB > 0 {
		con.Mem = &c.MemoryMB
	}
	if len(c.Tags) > 0 {
		con.Tags = &c.Tags
	}
	return con
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
