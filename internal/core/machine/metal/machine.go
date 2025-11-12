package metal

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/jaypipes/pcidb"
	"github.com/juju/juju/core/base"
	"golang.org/x/sync/errgroup"
)

type Machine struct {
	ID                  string
	Hostname            string
	AgentStatus         string
	FQDN                string
	WorkloadAnnotations map[string]string
	GPUs                []GPU

	LastCommissionedAt time.Time
}

type GPU struct {
	VendorID    string
	ProductID   string
	VendorName  string
	ProductName string
}

type MachineRepo interface {
	List(ctx context.Context) ([]Machine, error)
	Get(ctx context.Context, id string) (*Machine, error)
	Release(ctx context.Context, id string, force bool) (*Machine, error)
	PowerOff(ctx context.Context, id string, comment string) (*Machine, error)
	Commission(ctx context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool) (*Machine, error)
}

type MachineManagerRepo interface {
	AddMachines(ctx context.Context, scope, uuid, fqdn, baseOS, baseChannel string) error
	DestroyMachines(ctx context.Context, scope string, force, keep, dryRun bool, maxWait *time.Duration, machines ...string) error
}

type NodeDeviceRepo interface {
	ListGPUs(ctx context.Context, machineID string) ([]GPU, error)
}

func (uc *MetalUseCase) ListMachines(ctx context.Context, scope string) ([]Machine, error) {
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}

	eg, ctx := errgroup.WithContext(ctx)

	for i := range machines {
		eg.Go(func() error {
			gpus, err := uc.nodeDevice.ListGPUs(ctx, machines[i].ID)
			if err != nil {
				return err
			}

			if err = uc.setGPUName(gpus); err != nil {
				return err
			}

			machines[i].GPUs = gpus
			return nil
		})

		eg.Go(func() error {
			lastCommissionedAt, err := uc.lastCommissionedAt(ctx, machines[i].ID)
			if err != nil {
				return err
			}
			machines[i].LastCommissionedAt = lastCommissionedAt
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return slices.DeleteFunc(machines, func(machine Machine) bool {
		machineScope, _ := uc.extractScopeFromMachine(&machine)
		return !strings.Contains(machineScope, scope) // empty
	}), nil
}

func (uc *MetalUseCase) GetMachine(ctx context.Context, id string) (*Machine, error) {
	machine, err := uc.machine.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// only get agent status when needed
	machine.AgentStatus, err = uc.agentStatus(ctx, machine)
	if err != nil {
		return nil, err
	}

	machine.GPUs, err = uc.nodeDevice.ListGPUs(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = uc.setGPUName(machine.GPUs); err != nil {
		return nil, err
	}

	machine.LastCommissionedAt, err = uc.lastCommissionedAt(ctx, id)
	if err != nil {
		return nil, err
	}

	return machine, nil
}

func (uc *MetalUseCase) CreateMachine(ctx context.Context, machineID, scopeName string) (*Machine, error) {
	// validate scope exists
	scope, err := uc.scope.Get(ctx, scopeName)
	if err != nil {
		return nil, err
	}

	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return nil, err
	}

	series, err := uc.provisioner.Get(ctx, "default_distro_series")
	if err != nil {
		return nil, err
	}

	base, err := base.GetBaseFromSeries(series)
	if err != nil {
		return nil, err
	}

	if err := uc.machineManager.AddMachines(ctx, scope.Name, scope.UUID, machine.FQDN, base.OS, base.Channel.String()); err != nil {
		return nil, err
	}

	return machine, nil
}

// Note: Delete from MAAS only.
func (uc *MetalUseCase) DeleteMachine(ctx context.Context, id string, force, purgeDisk bool) error {
	if purgeDisk {
		if err := uc.purgeDisk(ctx, id); err != nil {
			return err
		}
	}
	if _, err := uc.machine.Release(ctx, id, force); err != nil {
		return err
	}
	return nil
}

func (uc *MetalUseCase) CommissionMachine(ctx context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool) error {
	if _, err := uc.machine.Commission(ctx, id, enableSSH, skipBMCConfig, skipNetworking, skipStorage); err != nil {
		return err
	}
	return nil
}

func (uc *MetalUseCase) PowerOffMachine(ctx context.Context, id, comment string) (*Machine, error) {
	return uc.machine.PowerOff(ctx, id, comment)
}

func (uc *MetalUseCase) extractScopeFromMachine(m *Machine) (string, error) {
	v, ok := m.WorkloadAnnotations["juju-machine-id"]
	if !ok {
		return "", fmt.Errorf("machine %q missing workload annotation juju-machine-id", m.Hostname)
	}

	token := strings.Split(v, "-")
	if len(token) < 1 {
		return "", fmt.Errorf("invalid juju-machine-id format for machine %q", m.Hostname)
	}

	return token[0], nil
}

func (uc *MetalUseCase) extractJujuIDFromMachine(m *Machine) (string, error) {
	v, ok := m.WorkloadAnnotations["juju-machine-id"]
	if !ok {
		return "", fmt.Errorf("machine %q missing workload annotation juju-machine-id", m.Hostname)
	}

	token := strings.Split(v, "-")
	if len(token) < 2 {
		return "", fmt.Errorf("invalid juju-machine-id format for machine %q", m.Hostname)
	}

	return token[1], nil
}

func (uc *MetalUseCase) agentStatus(ctx context.Context, machine *Machine) (string, error) {
	scope, err := uc.extractScopeFromMachine(machine)
	if err != nil {
		return "", nil // ignore
	}

	jujuID, err := uc.extractJujuIDFromMachine(machine)
	if err != nil {
		return "", nil // ignore
	}

	return uc.orchestrator.AgentStatus(ctx, scope, jujuID)
}

func (uc *MetalUseCase) setGPUName(gpus []GPU) error {
	pci, err := pcidb.New()
	if err != nil {
		return err
	}

	for i := range gpus {
		if vendor, ok := pci.Vendors[gpus[i].VendorID]; ok {
			gpus[i].VendorName = vendor.Name
		}

		if product, ok := pci.Products[gpus[i].VendorID+gpus[i].ProductID]; ok {
			gpus[i].ProductName = product.Name
		}
	}

	return nil
}

// TODO: improve performance by parallel execution
// TODO: osd devices in unit config or app config
func (uc *MetalUseCase) purgeDisk(ctx context.Context, machineID string) error {
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return err
	}

	scope, err := uc.extractScopeFromMachine(machine)
	if err != nil {
		return err
	}

	jujuID, err := uc.extractJujuIDFromMachine(machine)
	if err != nil {
		return err
	}

	apps, err := uc.orchestrator.ListApplications(ctx, scope, jujuID)
	if err != nil {
		return err
	}

	for _, app := range apps {
		if !strings.Contains(app.Charm, "ceph-osd") {
			continue
		}

		config, err := uc.facility.Config(ctx, scope, app.Name)
		if err != nil {
			continue
		}

		info, ok := config["osd-devices"].(map[string]interface{})
		if !ok {
			continue
		}

		val, ok := info["value"]
		if !ok || val == nil {
			continue
		}

		osdDevices := strings.SplitSeq(val.(string), " ")

		for osdDevice := range osdDevices {
			for _, unitName := range app.UnitNames {
				if _, err := uc.action.Execute(ctx, scope, unitName, fmt.Sprintf("sudo dd if=/dev/zero of=%s bs=1M count=200000", osdDevice)); err != nil {
					continue
				}
				break
			}
		}
	}

	return nil
}
