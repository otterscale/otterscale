package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/canonical/gomaasclient/entity/node"
	"github.com/juju/juju/rpc/params"

	"github.com/openhdc/otterscale/internal/domain/model"
)

func (s *NexusService) ListMachines(ctx context.Context, scopeUUID string) ([]model.Machine, error) {
	ms, err := s.machine.List(ctx)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(ms, func(m model.Machine) bool {
		modelUUID, _ := getJujuModelUUID(m.WorkloadAnnotations)
		return !strings.Contains(modelUUID, scopeUUID) // empty
	}), nil
}

func (s *NexusService) GetMachine(ctx context.Context, id string) (*model.Machine, error) {
	return s.machine.Get(ctx, id)
}

func (s *NexusService) CreateMachine(ctx context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool, uuid string, tags []string) (*model.Machine, error) {
	if err := s.AddMachineTags(ctx, id, tags); err != nil {
		return nil, err
	}

	commissionParams := &model.MachineCommissionParams{
		TestingScripts: "",
		EnableSSH:      boolToInt(enableSSH),
		SkipBMCConfig:  boolToInt(skipBMCConfig),
		SkipNetworking: boolToInt(skipNetworking),
		SkipStorage:    boolToInt(skipStorage),
	}
	machine, err := s.machine.Commission(ctx, id, commissionParams)
	if err != nil {
		return nil, err
	}

	if err := s.waitForMachineReady(ctx, id); err != nil {
		return nil, err
	}

	base, err := s.imageBase(ctx)
	if err != nil {
		return nil, err
	}

	os, channel := "", ""
	if base != nil {
		os = base.OS
		channel = base.Channel.String()
	}

	addMachineParams := []model.MachineAddParams{
		{
			Placement: &model.Placement{Scope: uuid, Directive: machine.FQDN},
			Jobs:      []model.MachineJob{model.JobHostUnits},
			Base:      &params.Base{Name: os, Channel: channel},
		},
	}
	results, err := s.machineManager.AddMachines(ctx, uuid, addMachineParams)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return machine, nil
}

func (s *NexusService) DeleteMachine(ctx context.Context, id string, force bool) error {
	// maas
	m, err := s.machine.Release(ctx, id, force)
	if err != nil {
		return err
	}

	// juju
	uuid, _ := getJujuModelUUID(m.WorkloadAnnotations)
	machine, _ := getJujuMachineID(m.WorkloadAnnotations)
	if uuid == "" || machine == "" {
		return nil
	}
	results, err := s.machineManager.DestroyMachines(ctx, uuid, force, machine)
	if err != nil {
		return err
	}
	for _, result := range results {
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
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

func (s *NexusService) waitForMachineReady(ctx context.Context, id string) error {
	const tickInterval = 10 * time.Second
	const timeoutDuration = 10 * time.Minute

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	timeout := time.After(timeoutDuration)
	for {
		select {
		case <-ticker.C:
			machine, err := s.machine.Get(ctx, id)
			if err != nil {
				return err
			}

			if machine.Status == node.StatusReady {
				break
			}
			continue

		case <-timeout:
			return fmt.Errorf("timeout waiting for machine %s to become ready", id)

		case <-ctx.Done():
			return ctx.Err()
		}

		break
	}

	return nil
}

func (s *NexusService) JujuToMAASMachineMap(ctx context.Context, uuid string) (map[string]string, error) {
	msm, err := s.listJujuMachines(ctx, uuid)
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	for key := range msm {
		m[key] = string(msm[key].InstanceId)
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
		m[string(msm[key].InstanceId)] = key
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

func getJujuModelUUID(m map[string]string) (string, error) {
	v, ok := m["juju-model-uuid"]
	if !ok {
		return "", errors.New("juju model uuid not found")
	}
	return v, nil
}

func getJujuMachineID(m map[string]string) (string, error) {
	v, ok := m["juju-machine-id"]
	if !ok {
		return "", errors.New("juju machine uuid not found")
	}
	token := strings.Split(v, "-")
	return token[len(token)-1], nil
}
