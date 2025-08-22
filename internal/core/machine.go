package core

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/node"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/rpc/params"
	"golang.org/x/sync/errgroup"
)

const JobHostUnits = model.JobHostUnits

type (
	NUMANode         = entity.NUMANode
	BlockDevice      = entity.BlockDevice
	NetworkInterface = entity.NetworkInterface
	Event            = entity.Event
)

type Machine struct {
	*entity.Machine
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

type MachineRepo interface {
	List(ctx context.Context) ([]Machine, error)
	Get(ctx context.Context, systemID string) (*Machine, error)
	Release(ctx context.Context, systemID string, params *entity.MachineReleaseParams) (*Machine, error)
	PowerOff(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*Machine, error)
	Commission(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*Machine, error)
}

type MachineManagerRepo interface {
	AddMachines(ctx context.Context, uuid string, params []params.AddMachineParams) error
	DestroyMachines(ctx context.Context, uuid string, force, keep, dryRun bool, maxWait *time.Duration, machines ...string) error
}

type ServerRepo interface {
	Get(ctx context.Context, name string) (string, error)
	Update(ctx context.Context, name, value string) error
}

type ClientRepo interface {
	Status(ctx context.Context, uuid string, patterns []string) (*params.FullStatus, error)
}

type EventRepo interface {
	Get(ctx context.Context, systemID string) ([]Event, error)
}

type MachineUseCase struct {
	machine        MachineRepo
	machineManager MachineManagerRepo
	server         ServerRepo
	client         ClientRepo
	tag            TagRepo
	action         ActionRepo
	facility       FacilityRepo
	event          EventRepo
}

func NewMachineUseCase(machine MachineRepo, machineManager MachineManagerRepo, server ServerRepo, client ClientRepo, tag TagRepo, action ActionRepo, facility FacilityRepo, event EventRepo) *MachineUseCase {
	return &MachineUseCase{
		machine:        machine,
		machineManager: machineManager,
		server:         server,
		client:         client,
		tag:            tag,
		action:         action,
		facility:       facility,
		event:          event,
	}
}

func (uc *MachineUseCase) ListMachines(ctx context.Context, scopeUUID string) ([]Machine, error) {
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}
	eg, ctx := errgroup.WithContext(ctx)
	for i := range machines {
		eg.Go(func() error {
			var err error
			machines[i].LastCommissioned, err = uc.getLastCommissioned(ctx, machines[i].SystemID)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return uc.filterMachines(machines, scopeUUID), nil
}

func (uc *MachineUseCase) GetMachine(ctx context.Context, id string) (*Machine, error) {
	machine, err := uc.machine.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	machine.LastCommissioned, err = uc.getLastCommissioned(ctx, id)
	if err != nil {
		return nil, err
	}
	return machine, nil
}

func (uc *MachineUseCase) CreateMachine(ctx context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool, uuid string, tags []string) (*Machine, error) {
	// tag
	eg, egctx := errgroup.WithContext(ctx)
	for _, tag := range tags {
		eg.Go(func() error {
			return uc.tag.AddMachines(egctx, tag, []string{id})
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// commission
	commissionParams := &entity.MachineCommissionParams{
		TestingScripts: "",
		EnableSSH:      boolToInt(enableSSH),
		SkipBMCConfig:  boolToInt(skipBMCConfig),
		SkipNetworking: boolToInt(skipNetworking),
		SkipStorage:    boolToInt(skipStorage),
	}
	machine, err := uc.machine.Commission(ctx, id, commissionParams)
	if err != nil {
		return nil, err
	}

	// wait
	if err := uc.waitForMachineReady(ctx, id); err != nil {
		return nil, err
	}

	// add
	base, err := defaultBase(ctx, uc.server)
	if err != nil {
		return nil, err
	}
	addMachineParams := []params.AddMachineParams{
		{
			Placement: &instance.Placement{Scope: uuid, Directive: machine.FQDN},
			Jobs:      []model.MachineJob{JobHostUnits},
			Base:      &params.Base{Name: base.OS, Channel: base.Channel.String()},
		},
	}
	if err := uc.machineManager.AddMachines(ctx, uuid, addMachineParams); err != nil {
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

func (uc *MachineUseCase) filterMachines(machines []Machine, scopeUUID string) []Machine {
	return slices.DeleteFunc(machines, func(machine Machine) bool {
		modelUUID, _ := getJujuModelUUID(machine.WorkloadAnnotations)
		return !strings.Contains(modelUUID, scopeUUID) // empty
	})
}

func (uc *MachineUseCase) waitForMachineReady(ctx context.Context, id string) error {
	const tickInterval = 10 * time.Second
	const timeoutDuration = 10 * time.Minute

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	timeout := time.After(timeoutDuration)
	for {
		select {
		case <-ticker.C:
			machine, err := uc.machine.Get(ctx, id)
			if err != nil {
				return err
			}
			if machine.Status == node.StatusReady {
				return nil
			}
			continue

		case <-timeout:
			return fmt.Errorf("timeout waiting for machine %s to become ready", id)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (uc *MachineUseCase) purgeDisk(ctx context.Context, machineSystemID string) error {
	machine, err := uc.machine.Get(ctx, machineSystemID)
	if err != nil {
		return err
	}
	uuid, err := getJujuModelUUID(machine.WorkloadAnnotations)
	if err != nil {
		return err
	}
	id, err := getJujuMachineID(machine.WorkloadAnnotations)
	if err != nil {
		return err
	}
	s, err := uc.client.Status(ctx, uuid, []string{"machine", id})
	if err != nil {
		return err
	}
	for name := range s.Applications {
		if !strings.Contains(s.Applications[name].Charm, "ceph-osd") {
			continue
		}
		config, err := uc.facility.GetConfig(ctx, uuid, name)
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
				id, err := uc.action.RunCommand(ctx, uuid, uname, fmt.Sprintf("sudo dd if=/dev/zero of=%s bs=1M count=200000", osdDevice))
				if err != nil {
					continue
				}
				if _, err := waitForActionCompleted(ctx, uc.action, uuid, id, time.Second*10, time.Minute*10); err != nil { //nolint:mnd
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
