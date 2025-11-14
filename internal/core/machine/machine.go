package machine

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/canonical/gomaasclient/entity"
	"github.com/jaypipes/pcidb"
	"github.com/juju/juju/core/base"
	"github.com/otterscale/otterscale/internal/core/configuration"
	"github.com/otterscale/otterscale/internal/core/facility"
	"github.com/otterscale/otterscale/internal/core/facility/action"
	"github.com/otterscale/otterscale/internal/core/scope"
	"golang.org/x/sync/errgroup"
)

// Machine represents a MAAS machine.
type Machine = entity.Machine

type MachineData struct {
	*Machine

	GPUs               []GPU
	LastCommissionedAt time.Time
	AgentStatus        string
}

type GPU = entity.NodeDevice

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

type OrchestratorRepo interface {
	AgentStatus(ctx context.Context, scope string, jujuID string) (string, error)
}

type MachineUseCase struct {
	event          EventRepo
	machine        MachineRepo
	machineManager MachineManagerRepo
	nodeDevice     NodeDeviceRepo
	orchestrator   OrchestratorRepo

	action      action.ActionRepo
	facility    facility.FacilityRepo
	provisioner configuration.ProvisionerRepo
	scope       scope.ScopeRepo
}

func NewMachineUseCase(event EventRepo, machine MachineRepo, machineManager MachineManagerRepo, nodeDevice NodeDeviceRepo, orchestrator OrchestratorRepo, action action.ActionRepo, facility facility.FacilityRepo, provisioner configuration.ProvisionerRepo, scope scope.ScopeRepo) *MachineUseCase {
	return &MachineUseCase{
		event:          event,
		machine:        machine,
		machineManager: machineManager,
		nodeDevice:     nodeDevice,
		orchestrator:   orchestrator,
		action:         action,
		facility:       facility,
		provisioner:    provisioner,
		scope:          scope,
	}
}

func (uc *MachineUseCase) ListMachines(ctx context.Context, scope string) ([]MachineData, error) {
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]MachineData, len(machines))
	for i := range machines {
		ret[i].Machine = &machines[i]
	}

	eg, ctx := errgroup.WithContext(ctx)

	for i := range ret {
		eg.Go(func() error {
			gpus, err := uc.nodeDevice.ListGPUs(ctx, machines[i].SystemID)
			if err == nil {
				if err = uc.setGPUName(gpus); err != nil {
					return err
				}
				ret[i].GPUs = gpus
			}
			return err
		})

		eg.Go(func() error {
			lastCommissionedAt, err := uc.lastCommissionedAt(ctx, machines[i].SystemID)
			if err == nil {
				ret[i].LastCommissionedAt = lastCommissionedAt
			}
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return slices.DeleteFunc(ret, func(machine MachineData) bool {
		machineScope, _ := uc.extractScopeFromMachine(machine.Machine)
		return !strings.Contains(machineScope, scope) // empty
	}), nil
}

func (uc *MachineUseCase) GetMachine(ctx context.Context, id string) (*MachineData, error) {
	machine, err := uc.machine.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// only get agent status when needed
	agentStatus, err := uc.agentStatus(ctx, machine)
	if err != nil {
		return nil, err
	}

	gpus, err := uc.nodeDevice.ListGPUs(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = uc.setGPUName(gpus); err != nil {
		return nil, err
	}

	lastCommissionedAt, err := uc.lastCommissionedAt(ctx, id)
	if err != nil {
		return nil, err
	}

	return &MachineData{
		Machine:            machine,
		GPUs:               gpus,
		LastCommissionedAt: lastCommissionedAt,
		AgentStatus:        agentStatus,
	}, nil
}

func (uc *MachineUseCase) CreateMachine(ctx context.Context, machineID, scopeName string) (*Machine, error) {
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
func (uc *MachineUseCase) DeleteMachine(ctx context.Context, id string, force, purgeDisk bool) error {
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

func (uc *MachineUseCase) CommissionMachine(ctx context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool) error {
	if _, err := uc.machine.Commission(ctx, id, enableSSH, skipBMCConfig, skipNetworking, skipStorage); err != nil {
		return err
	}
	return nil
}

func (uc *MachineUseCase) PowerOffMachine(ctx context.Context, id, comment string) (*Machine, error) {
	return uc.machine.PowerOff(ctx, id, comment)
}

func (uc *MachineUseCase) extractScopeFromMachine(m *Machine) (string, error) {
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

func (uc *MachineUseCase) extractJujuIDFromMachine(m *Machine) (string, error) {
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

func (uc *MachineUseCase) agentStatus(ctx context.Context, machine *Machine) (string, error) {
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

func (uc *MachineUseCase) setGPUName(gpus []GPU) error {
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
func (uc *MachineUseCase) purgeDisk(ctx context.Context, machineID string) error {
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

	apps, err := uc.facility.List(ctx, scope, jujuID)
	if err != nil {
		return err
	}

	for _, app := range apps {
		if app.Name != scope+"-ceph-osd" {
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
			for unitName := range app.Status.Units {
				if _, err := uc.action.Execute(ctx, scope, unitName, fmt.Sprintf("sudo dd if=/dev/zero of=%s bs=1M count=200000", osdDevice)); err != nil {
					continue
				}
				break
			}
		}
	}

	return nil
}
