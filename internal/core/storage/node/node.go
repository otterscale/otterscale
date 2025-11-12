package node

import (
	"context"
	"fmt"
	"strings"

	"github.com/otterscale/otterscale/internal/core/machine"
)

type Monitor struct {
	Name    string
	Machine *machine.Machine
}

type ObjectStorageDaemon struct {
	Hostname string
	Machine  *machine.Machine
}

type NodeRepo interface {
	ListMonitors(ctx context.Context, scope string) ([]Monitor, error)
	ListObjectStorageDaemons(ctx context.Context, scope string) ([]ObjectStorageDaemon, error)
	DoSMART(ctx context.Context, scope string, who string) (map[string][]string, error)
}

type NodeUseCase struct {
	node NodeRepo

	machine machine.MachineRepo
}

func NewNodeUseCase(node NodeRepo, machine machine.MachineRepo) *NodeUseCase {
	return &NodeUseCase{
		node:    node,
		machine: machine,
	}
}

func (uc *NodeUseCase) ListMonitors(ctx context.Context, scope string) ([]Monitor, error) {
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

func (uc *NodeUseCase) ListObjectStorageDaemons(ctx context.Context, scope string) ([]ObjectStorageDaemon, error) {
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

func (uc *NodeUseCase) DoSMART(ctx context.Context, scope string, osd string) (map[string][]string, error) {
	return uc.node.DoSMART(ctx, scope, osd)
}

func (uc *NodeUseCase) extractScopeFromMachine(m *machine.Machine) (string, error) {
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

func (uc *NodeUseCase) extractJujuIDFromMachine(m *machine.Machine) (string, error) {
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

func (uc *NodeUseCase) extractJujuIDFromMonitor(m *Monitor) (string, error) {
	token := strings.Split(m.Name, "-")
	if len(token) < 3 {
		return "", fmt.Errorf("invalid monitor name format for monitor %q", m.Name)
	}

	return token[2], nil
}

func (uc *NodeUseCase) setObjectStorageDaemonMachine(osds []ObjectStorageDaemon, machines []machine.Machine) {
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

func (uc *NodeUseCase) setMonitorMachine(scope string, monitors []Monitor, machines []machine.Machine) {
	for i := range monitors {
		monitorJujuID, err := uc.extractJujuIDFromMonitor(&monitors[i])
		if err != nil {
			continue
		}

		for j := range machines {
			machineScope, err := uc.extractScopeFromMachine(&machines[j])
			if err != nil {
				continue
			}

			if machineScope != scope {
				continue
			}

			machineJujuID, err := uc.extractJujuIDFromMachine(&machines[j])
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
