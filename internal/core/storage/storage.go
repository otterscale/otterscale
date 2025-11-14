package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/otterscale/otterscale/internal/core/machine"
)

type Monitor struct {
	Leader        bool
	Name          string
	Rank          uint64
	PublicAddress string
	Machine       *machine.Machine
}

type ObjectStorageDaemon struct {
	ID          int64
	Name        string
	Up          bool
	In          bool
	Exists      bool
	DeviceClass string
	Size        uint64
	Used        uint64
	PGCount     uint64
	Hostname    string
	Machine     *machine.Machine
}

type NodeRepo interface {
	ListMonitors(ctx context.Context, scope string) ([]Monitor, error)
	ListObjectStorageDaemons(ctx context.Context, scope string) ([]ObjectStorageDaemon, error)
	DoSMART(ctx context.Context, scope string, who string) (map[string][]string, error)
	Config(scope string) (host string, id string, key string, err error)
}

type StorageUseCase struct {
	node NodeRepo
	pool PoolRepo

	machine machine.MachineRepo
}

func NewStorageUseCase(node NodeRepo, pool PoolRepo, machine machine.MachineRepo) *StorageUseCase {
	return &StorageUseCase{
		node:    node,
		pool:    pool,
		machine: machine,
	}
}

func (uc *StorageUseCase) ListMonitors(ctx context.Context, scope string) ([]Monitor, error) {
	monitors, err := uc.node.ListMonitors(ctx, scope)
	if err != nil {
		return nil, err
	}

	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}

	uc.setMonitorMachine(scope, monitors, machines)

	return monitors, nil
}

func (uc *StorageUseCase) ListObjectStorageDaemons(ctx context.Context, scope string) ([]ObjectStorageDaemon, error) {
	osds, err := uc.node.ListObjectStorageDaemons(ctx, scope)
	if err != nil {
		return nil, err
	}

	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}

	uc.setObjectStorageDaemonMachine(osds, machines)

	return osds, nil
}

func (uc *StorageUseCase) DoSMART(ctx context.Context, scope string, osd string) (map[string][]string, error) {
	return uc.node.DoSMART(ctx, scope, osd)
}

func (uc *StorageUseCase) extractJujuIDFromMonitor(m *Monitor) (string, error) {
	token := strings.Split(m.Name, "-")
	if len(token) < 3 {
		return "", fmt.Errorf("invalid monitor name format for monitor %q", m.Name)
	}

	return token[2], nil
}

func (uc *StorageUseCase) setObjectStorageDaemonMachine(osds []ObjectStorageDaemon, machines []machine.Machine) {
	for i := range osds {
		for j := range machines {
			if osds[i].Hostname != machines[j].Hostname {
				continue
			}

			osds[i].Machine = &machines[j]
			break
		}
	}
}

func (uc *StorageUseCase) setMonitorMachine(scope string, monitors []Monitor, machines []machine.Machine) {
	for i := range monitors {
		monitorJujuID, err := uc.extractJujuIDFromMonitor(&monitors[i])
		if err != nil {
			continue
		}

		for j := range machines {
			machineScope, err := uc.machine.ExtractScope(&machines[j])
			if err != nil {
				continue
			}

			if machineScope != scope {
				continue
			}

			machineJujuID, err := uc.machine.ExtractJujuID(&machines[j])
			if err != nil {
				continue
			}

			if machineJujuID != monitorJujuID {
				continue
			}

			monitors[i].Machine = &machines[j]
			break
		}
	}
}
