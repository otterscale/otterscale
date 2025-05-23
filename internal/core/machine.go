package core

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/node"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/rpc/params"
	"golang.org/x/sync/errgroup"
)

const JobHostUnits = model.JobHostUnits

type (
	Machine                 = entity.Machine
	MachineCommissionParams = entity.MachineCommissionParams
	MachinePowerOnParams    = entity.MachinePowerOnParams
	MachinePowerOffParams   = entity.MachinePowerOffParams
	NUMANode                = entity.NUMANode
	BlockDevice             = entity.BlockDevice
	NetworkInterface        = entity.NetworkInterface
)

type (
	MachineStatus    = params.MachineStatus
	MachineAddParams = params.AddMachineParams
	Placement        = instance.Placement
	Constraint       = constraints.Value
	MachineJob       = model.MachineJob
)

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
	List(ctx context.Context) ([]entity.Machine, error)
	Get(ctx context.Context, systemID string) (*entity.Machine, error)
	Release(ctx context.Context, systemID string, params *entity.MachineReleaseParams) (*entity.Machine, error)
	PowerOn(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error)
	PowerOff(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error)
	Commission(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error)
}

type MachineManagerRepo interface {
	AddMachines(ctx context.Context, uuid string, params []params.AddMachineParams) ([]params.AddMachinesResult, error)
	DestroyMachines(ctx context.Context, uuid string, force, keep, dryRun bool, maxWait *time.Duration, machines ...string) ([]params.DestroyMachineResult, error)
}

type ServerRepo interface {
	Get(ctx context.Context, name string) (string, error)
	Update(ctx context.Context, name, value string) error
}

type ClientRepo interface {
	Status(ctx context.Context, uuid string, patterns []string) (*params.FullStatus, error)
}

type MachineUseCase struct {
	machine        MachineRepo
	machineManager MachineManagerRepo
	server         ServerRepo
	client         ClientRepo
	tag            TagRepo
}

func NewMachineUseCase(machine MachineRepo, machineManager MachineManagerRepo, server ServerRepo, client ClientRepo, tag TagRepo) *MachineUseCase {
	return &MachineUseCase{
		machine:        machine,
		machineManager: machineManager,
		server:         server,
		client:         client,
		tag:            tag,
	}
}

func (uc *MachineUseCase) ListMachines(ctx context.Context, scopeUUID string) ([]Machine, error) {
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}
	return uc.filterMachines(machines, scopeUUID), nil
}

func (uc *MachineUseCase) GetMachine(ctx context.Context, id string) (*Machine, error) {
	return uc.machine.Get(ctx, id)
}

func (uc *MachineUseCase) CreateMachine(ctx context.Context, id string, enableSSH, skipBMCConfig, skipNetworking, skipStorage bool, uuid string, tags []string) (*Machine, error) {
	// tag
	eg, ctx := errgroup.WithContext(ctx)
	for _, tag := range tags {
		tag := tag // fixed on go 1.22
		eg.Go(func() error {
			return uc.tag.AddMachines(ctx, tag, []string{id})
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// commission
	commissionParams := &MachineCommissionParams{
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
	addMachineParams := []MachineAddParams{
		{
			Placement: &Placement{Scope: uuid, Directive: machine.FQDN},
			Jobs:      []MachineJob{JobHostUnits},
			Base:      &params.Base{Name: base.OS, Channel: base.Channel.String()},
		},
	}
	results, err := uc.machineManager.AddMachines(ctx, uuid, addMachineParams)
	if err != nil {
		return nil, err
	}
	errs := []error{}
	for _, result := range results {
		errs = append(errs, result.Error)
	}
	err = errors.Join(errs...)
	if err != nil {
		return nil, err
	}

	return machine, nil
}

// Note: Delete from MAAS only.
func (uc *MachineUseCase) DeleteMachine(ctx context.Context, id string, force bool) error {
	params := &entity.MachineReleaseParams{
		Force: force,
	}
	_, err := uc.machine.Release(ctx, id, params)
	return err
}

func (uc *MachineUseCase) PowerOffMachine(ctx context.Context, id, comment string) (*Machine, error) {
	params := &MachinePowerOffParams{
		Comment: comment,
	}
	return uc.machine.PowerOff(ctx, id, params)
}

func (uc *MachineUseCase) JujuToMAASMachineMap(ctx context.Context, uuid string) (map[string]string, error) {
	status, err := uc.client.Status(ctx, uuid, []string{"machine", "*"})
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	for name := range status.Machines {
		m[name] = string(status.Machines[name].InstanceId)
	}
	return m, nil
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
