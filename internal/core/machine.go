package core

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/canonical/gomaasclient/entity"
	"github.com/jaypipes/pcidb"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/rpc/params"
	"golang.org/x/sync/errgroup"
)

type (
	BlockDevice      = entity.BlockDevice
	Event            = entity.Event
	NetworkInterface = entity.NetworkInterface
	NodeDevice       = entity.NodeDevice
	NUMANode         = entity.NUMANode
	Tag              = entity.Tag
)

type Machine struct {
	*entity.Machine
	AgentStatus      params.DetailedStatus // from Juju agent
	GPUs             []NodeDevice
	LastCommissioned time.Time
}

type MachinePlacement struct {
	LXD       bool
	KVM       bool
	Machine   bool
	MachineID string
}

type MachineConstraint struct {
	Architecture string
	CPUCores     uint64
	MemoryMB     uint64
	Tags         []string
}

type MachineFactor struct {
	*MachinePlacement
	*MachineConstraint
}

type ClientRepo interface {
	Status(ctx context.Context, scope string, patterns []string) (*params.FullStatus, error)
}

type EventRepo interface {
	Get(ctx context.Context, systemID string) ([]Event, error)
}

type MachineRepo interface {
	List(ctx context.Context) ([]Machine, error)
	Get(ctx context.Context, systemID string) (*Machine, error)
	Release(ctx context.Context, systemID string, params *entity.MachineReleaseParams) (*Machine, error)
	PowerOff(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*Machine, error)
	Commission(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*Machine, error)
}

type MachineManagerRepo interface {
	AddMachines(ctx context.Context, scope string, params []params.AddMachineParams) error
	DestroyMachines(ctx context.Context, scope string, force, keep, dryRun bool, maxWait *time.Duration, machines ...string) error
}

type NodeDeviceRepo interface {
	List(ctx context.Context, systemID, hardwareType string) ([]NodeDevice, error)
}

type ServerRepo interface {
	Get(ctx context.Context, name string) (string, error)
	Update(ctx context.Context, name, value string) error
}

type TagRepo interface {
	List(ctx context.Context) ([]Tag, error)
	Get(ctx context.Context, name string) (*Tag, error)
	Create(ctx context.Context, name, comment string) (*Tag, error)
	Delete(ctx context.Context, name string) error
	AddMachines(ctx context.Context, name string, machineIDs []string) error
	RemoveMachines(ctx context.Context, name string, machineIDs []string) error
}

type MachineUseCase struct {
	action         ActionRepo
	client         ClientRepo
	event          EventRepo
	facility       FacilityRepo
	machine        MachineRepo
	machineManager MachineManagerRepo
	nodeDevice     NodeDeviceRepo
	scope          ScopeRepo
	server         ServerRepo
	tag            TagRepo
}

func NewMachineUseCase(action ActionRepo, client ClientRepo, event EventRepo, facility FacilityRepo, machine MachineRepo, machineManager MachineManagerRepo, nodeDevice NodeDeviceRepo, scope ScopeRepo, server ServerRepo, tag TagRepo) *MachineUseCase {
	return &MachineUseCase{
		action:         action,
		client:         client,
		event:          event,
		facility:       facility,
		machine:        machine,
		machineManager: machineManager,
		nodeDevice:     nodeDevice,
		scope:          scope,
		server:         server,
		tag:            tag,
	}
}

func (uc *MachineUseCase) ListMachines(ctx context.Context, scope string) ([]Machine, error) {
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}
	eg, ctx := errgroup.WithContext(ctx)
	for i := range machines {
		eg.Go(func() error {
			gpus, err := uc.nodeDevice.List(ctx, machines[i].SystemID, "gpu")
			if err != nil {
				return err
			}
			machines[i].GPUs, err = setGPUVendorProduct(gpus)
			return err
		})
		eg.Go(func() error {
			var err error
			machines[i].LastCommissioned, err = uc.getLastCommissioned(ctx, machines[i].SystemID)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return uc.filterMachines(machines, scope), nil
}

func (uc *MachineUseCase) GetMachine(ctx context.Context, id string) (*Machine, error) {
	machine, err := uc.machine.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	agentStatus, ok := uc.getAgentStatus(ctx, machine)
	if ok {
		machine.AgentStatus = agentStatus
	}
	gpus, err := uc.nodeDevice.List(ctx, id, "gpu")
	if err != nil {
		return nil, err
	}
	machine.GPUs, err = setGPUVendorProduct(gpus)
	if err != nil {
		return nil, err
	}
	machine.LastCommissioned, err = uc.getLastCommissioned(ctx, id)
	if err != nil {
		return nil, err
	}
	return machine, nil
}

func (uc *MachineUseCase) CreateMachine(ctx context.Context, machineID, scopeName string) (*Machine, error) {
	scope, err := uc.scope.Get(ctx, scopeName) // validate scope exists
	if err != nil {
		return nil, err
	}

	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return nil, err
	}

	base, err := defaultBase(ctx, uc.server)
	if err != nil {
		return nil, err
	}

	// add
	addMachineParams := []params.AddMachineParams{
		{
			Placement: &instance.Placement{Scope: scope.UUID, Directive: machine.FQDN},
			Jobs:      []model.MachineJob{model.JobHostUnits},
			Base:      &params.Base{Name: base.OS, Channel: base.Channel.String()},
		},
	}
	if err := uc.machineManager.AddMachines(ctx, scopeName, addMachineParams); err != nil {
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
	params := &entity.MachineReleaseParams{
		Force: force,
	}
	if _, err := uc.machine.Release(ctx, id, params); err != nil {
		return err
	}
	return nil
}

func (uc *MachineUseCase) CommissionMachine(ctx context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool) error {
	commissionParams := &entity.MachineCommissionParams{
		TestingScripts: "",
		EnableSSH:      boolToInt(enableSSH),
		SkipBMCConfig:  boolToInt(skipBMCConfig),
		SkipNetworking: boolToInt(skipNetworking),
		SkipStorage:    boolToInt(skipStorage),
	}
	if _, err := uc.machine.Commission(ctx, id, commissionParams); err != nil {
		return err
	}
	return nil
}

func (uc *MachineUseCase) PowerOffMachine(ctx context.Context, id, comment string) (*Machine, error) {
	params := &entity.MachinePowerOffParams{
		Comment: comment,
	}
	return uc.machine.PowerOff(ctx, id, params)
}

func (uc *MachineUseCase) AddMachineTags(ctx context.Context, id string, tags []string) error {
	eg, egctx := errgroup.WithContext(ctx)
	for _, tag := range tags {
		eg.Go(func() error {
			return uc.tag.AddMachines(egctx, tag, []string{id})
		})
	}
	return eg.Wait()
}

func (uc *MachineUseCase) RemoveMachineTags(ctx context.Context, id string, tags []string) error {
	eg, egctx := errgroup.WithContext(ctx)
	for _, tag := range tags {
		eg.Go(func() error {
			return uc.tag.RemoveMachines(egctx, tag, []string{id})
		})
	}
	return eg.Wait()
}

func (uc *MachineUseCase) ListTags(ctx context.Context) ([]Tag, error) {
	return uc.tag.List(ctx)
}

func (uc *MachineUseCase) GetTag(ctx context.Context, name string) (*Tag, error) {
	return uc.tag.Get(ctx, name)
}

func (uc *MachineUseCase) CreateTag(ctx context.Context, name, comment string) (*Tag, error) {
	return uc.tag.Create(ctx, name, comment)
}

func (uc *MachineUseCase) DeleteTag(ctx context.Context, name string) error {
	return uc.tag.Delete(ctx, name)
}

func (uc *MachineUseCase) filterMachines(machines []Machine, scope string) []Machine {
	return slices.DeleteFunc(machines, func(machine Machine) bool {
		model, _ := getJujuModelName(machine.WorkloadAnnotations)
		return !strings.Contains(model, scope) // empty
	})
}

func (uc *MachineUseCase) getAgentStatus(ctx context.Context, machine *Machine) (params.DetailedStatus, bool) {
	scope, err := getJujuModelName(machine.WorkloadAnnotations)
	if err != nil {
		return params.DetailedStatus{}, false
	}
	jujuMachineID, err := getJujuMachineID(machine.WorkloadAnnotations)
	if err != nil {
		return params.DetailedStatus{}, false
	}
	fullStatus, err := uc.client.Status(ctx, scope, []string{"machine", jujuMachineID})
	if err != nil {
		return params.DetailedStatus{}, false
	}
	machineStatus, ok := fullStatus.Machines[jujuMachineID]
	if !ok {
		return params.DetailedStatus{}, false
	}
	return machineStatus.AgentStatus, true
}

func (uc *MachineUseCase) purgeDisk(ctx context.Context, machineSystemID string) error {
	machine, err := uc.machine.Get(ctx, machineSystemID)
	if err != nil {
		return err
	}
	scope, err := getJujuModelName(machine.WorkloadAnnotations)
	if err != nil {
		return err
	}
	id, err := getJujuMachineID(machine.WorkloadAnnotations)
	if err != nil {
		return err
	}
	s, err := uc.client.Status(ctx, scope, []string{"machine", id})
	if err != nil {
		return err
	}
	for name := range s.Applications {
		if !strings.Contains(s.Applications[name].Charm, "ceph-osd") {
			continue
		}
		config, err := uc.facility.GetConfig(ctx, scope, name)
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
			for uname := range s.Applications[name].Units {
				id, err := uc.action.RunCommand(ctx, scope, uname, fmt.Sprintf("sudo dd if=/dev/zero of=%s bs=1M count=200000", osdDevice))
				if err != nil {
					continue
				}
				if _, err := waitForActionCompleted(ctx, uc.action, scope, id, time.Second*10, time.Minute*10); err != nil { //nolint:mnd // default
					return err
				}
				break
			}
		}
	}
	return nil
}

func (uc *MachineUseCase) getLastCommissioned(ctx context.Context, machineSystemID string) (time.Time, error) {
	events, err := uc.event.Get(ctx, machineSystemID)
	if err != nil {
		return time.Time{}, err
	}
	for i := range events { // desc
		if events[i].Type == "Commissioning" {
			return time.Parse("Mon, 02 Jan. 2006 15:04:05", events[i].Created)
		}
	}
	return time.Time{}, nil
}

func setGPUVendorProduct(gpus []NodeDevice) ([]NodeDevice, error) {
	pci, err := pcidb.New()
	if err != nil {
		return nil, err
	}
	newGPUs := make([]NodeDevice, len(gpus))
	copy(newGPUs, gpus)
	for i := range gpus {
		if vendor, ok := pci.Vendors[gpus[i].VendorID]; ok {
			newGPUs[i].VendorName = vendor.Name
		}
		if product, ok := pci.Products[gpus[i].VendorID+gpus[i].ProductID]; ok {
			newGPUs[i].ProductName = product.Name
		}
	}
	return newGPUs, nil
}
