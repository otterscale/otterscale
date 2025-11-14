package machine

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/canonical/gomaasclient/entity"
	"github.com/jaypipes/pcidb"
	"github.com/juju/juju/core/base"
	"github.com/otterscale/otterscale/internal/core/configuration"
	"github.com/otterscale/otterscale/internal/core/scope"
	"golang.org/x/sync/errgroup"
)

type (
	// Machine represents a MAAS Machine resource.
	Machine = entity.Machine

	// NUMANode represents a MAAS NUMANode resource.
	NUMANode = entity.NUMANode

	// BlockDevice represents a MAAS BlockDevice resource.
	BlockDevice = entity.BlockDevice

	// NetworkInterface represents a MAAS NetworkInterface resource.
	NetworkInterface = entity.NetworkInterface

	// GPU represents a MAAS NodeDevice resource.
	GPU = entity.NodeDevice
)

type MachineData struct {
	*Machine

	GPUs               []GPU
	LastCommissionedAt time.Time
	AgentStatus        string
}

type MachineRepo interface {
	List(ctx context.Context) ([]Machine, error)
	Get(ctx context.Context, id string) (*Machine, error)
	Release(ctx context.Context, id string, force bool) (*Machine, error)
	PowerOff(ctx context.Context, id string, comment string) (*Machine, error)
	Commission(ctx context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool) (*Machine, error)
	ExtractScope(m *Machine) (string, error)
	ExtractJujuID(m *Machine) (string, error)
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

	provisioner configuration.ProvisionerRepo
	scope       scope.ScopeRepo
}

func NewMachineUseCase(event EventRepo, machine MachineRepo, machineManager MachineManagerRepo, nodeDevice NodeDeviceRepo, orchestrator OrchestratorRepo, provisioner configuration.ProvisionerRepo, scope scope.ScopeRepo) *MachineUseCase {
	return &MachineUseCase{
		event:          event,
		machine:        machine,
		machineManager: machineManager,
		nodeDevice:     nodeDevice,
		orchestrator:   orchestrator,
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
		machineScope, _ := uc.machine.ExtractScope(machine.Machine)
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

func (uc *MachineUseCase) CreateMachine(ctx context.Context, machineID, scopeName string) (*MachineData, error) {
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

	return &MachineData{
		Machine: machine,
	}, nil
}

// Note: Delete from MAAS only.
func (uc *MachineUseCase) DeleteMachine(ctx context.Context, id string, force bool) error {
	_, err := uc.machine.Release(ctx, id, force)
	return err
}

func (uc *MachineUseCase) CommissionMachine(ctx context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool) error {
	_, err := uc.machine.Commission(ctx, id, enableSSH, skipBMCConfig, skipNetworking, skipStorage)
	return err
}

func (uc *MachineUseCase) PowerOffMachine(ctx context.Context, id, comment string) (*MachineData, error) {
	machine, err := uc.machine.PowerOff(ctx, id, comment)
	if err != nil {
		return nil, err
	}

	return &MachineData{
		Machine: machine,
	}, nil
}

func (uc *MachineUseCase) agentStatus(ctx context.Context, machine *Machine) (string, error) {
	scope, err := uc.machine.ExtractScope(machine)
	if err != nil {
		return "", nil // ignore
	}

	jujuID, err := uc.machine.ExtractJujuID(machine)
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
