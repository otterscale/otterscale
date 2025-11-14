package gpu

import (
	"context"
	"slices"
	"strconv"
	"time"

	"github.com/Project-HAMi/HAMi/pkg/device"
	"github.com/otterscale/otterscale/internal/core/application/cluster"
	"github.com/otterscale/otterscale/internal/core/application/workload"
	"github.com/otterscale/otterscale/internal/core/machine"
	"github.com/otterscale/otterscale/internal/core/model"
)

const (
	hamiNodeNvidiaRegisterAnnotation   = "hami.io/node-nvidia-register"
	hamiVGPUNodeAnnotation             = "hami.io/vgpu-node"
	hamiVGPUDevicesAllocatedAnnotation = "hami.io/vgpu-devices-allocated"
	hamiBindTimeAnnotation             = "hami.io/bind-time"
	hamiBindPhaseAnnotation            = "hami.io/bind-phase"
)

type Relations struct {
	Machines []machine.Machine
	GPUs     []GPURelation
	Pods     []PodRelation
}

type GPURelation struct {
	ID          string
	Index       uint32
	Count       int32
	Cores       int32
	MemoryBytes int64
	Type        string
	Health      bool
	MachineID   string
}

type PodRelation struct {
	Name         string
	Namespace    string
	ModelName    string
	BindingPhase string
	BoundAt      time.Time
	PodDevices   []PodDevice
}

type PodDevice struct {
	GPUID           string
	UsedCores       int32
	UsedMemoryBytes int64
}

type GPUUseCase struct {
	machine machine.MachineRepo
	node    cluster.NodeRepo
	pod     workload.PodRepo
}

func NewGPUUseCase(machine machine.MachineRepo, node cluster.NodeRepo, pod workload.PodRepo) *GPUUseCase {
	return &GPUUseCase{
		machine: machine,
		node:    node,
		pod:     pod,
	}
}

func (uc *GPUUseCase) ListGPURelationsByMachine(ctx context.Context, scope, machineID string) (*Relations, error) {
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return nil, err
	}

	labelSelector := hamiVGPUNodeAnnotation + "=" + machine.Hostname

	return uc.listRelations(ctx, scope, "", labelSelector)
}

func (uc *GPUUseCase) ListGPURelationsByModel(ctx context.Context, scope, namespace, modelName string) (*Relations, error) {
	labelSelector := model.ModelNameAnnotation + "=" + modelName

	return uc.listRelations(ctx, scope, namespace, labelSelector)
}

func (uc *GPUUseCase) listRelations(ctx context.Context, scope, namespace, selector string) (*Relations, error) {
	pods, err := uc.pod.List(ctx, scope, namespace, selector)
	if err != nil {
		return nil, err
	}

	nodes, err := uc.node.List(ctx, scope, "")
	if err != nil {
		return nil, err
	}

	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}

	return uc.buildRelations(pods, nodes, machines)
}

func (uc *GPUUseCase) buildRelations(pods []workload.Pod, nodes []cluster.Node, machines []machine.Machine) (*Relations, error) {
	filteredMachines := uc.filterMachines(pods, machines)
	filteredNodes := uc.filterNodes(pods, nodes)

	gpuRelations, err := uc.buildGPURelations(filteredMachines, filteredNodes)
	if err != nil {
		return nil, err
	}

	podRelations, err := uc.buildPodRelations(pods)
	if err != nil {
		return nil, err
	}

	return &Relations{
		Machines: filteredMachines,
		GPUs:     gpuRelations,
		Pods:     podRelations,
	}, nil
}

func (uc *GPUUseCase) filterMachines(pods []workload.Pod, machines []machine.Machine) []machine.Machine {
	nodeNames := []string{}

	for i := range pods {
		if nodeName, ok := pods[i].Annotations[hamiVGPUNodeAnnotation]; ok {
			nodeNames = append(nodeNames, nodeName)
		}
	}

	ret := []machine.Machine{}

	for i := range machines {
		if !slices.Contains(nodeNames, machines[i].Hostname) {
			continue
		}

		ret = append(ret, machines[i])
	}

	return ret
}

func (uc *GPUUseCase) filterNodes(pods []workload.Pod, nodes []cluster.Node) []cluster.Node {
	nodeNames := []string{}

	for i := range pods {
		if nodeName, ok := pods[i].Annotations[hamiVGPUNodeAnnotation]; ok {
			nodeNames = append(nodeNames, nodeName)
		}
	}

	ret := []cluster.Node{}

	for i := range nodes {
		if !slices.Contains(nodeNames, nodes[i].Name) {
			continue
		}

		ret = append(ret, nodes[i])
	}

	return ret
}

func (uc *GPUUseCase) buildGPURelations(machines []machine.Machine, nodes []cluster.Node) ([]GPURelation, error) {
	machineMap := map[string]machine.Machine{}
	for i := range machines {
		machineMap[machines[i].Hostname] = machines[i]
	}

	gpus := []GPURelation{}

	for i := range nodes {
		annotation, ok := nodes[i].Annotations[hamiNodeNvidiaRegisterAnnotation]
		if !ok {
			continue
		}

		machine, ok := machineMap[nodes[i].Name]
		if !ok {
			continue
		}

		nodeDevices, err := device.DecodeNodeDevices(annotation)
		if err != nil {
			return nil, err
		}

		for _, nodeDevice := range nodeDevices {
			gpus = append(gpus, GPURelation{
				ID:          nodeDevice.ID,
				Index:       uint32(nodeDevice.Index), //nolint:gosec // uint to uint32
				Count:       nodeDevice.Count,
				Cores:       nodeDevice.Devcore,
				MemoryBytes: int64(nodeDevice.Devmem) * 1024 * 1024, // gigabytes to bytes
				Type:        nodeDevice.Type,
				Health:      nodeDevice.Health,
				MachineID:   machine.SystemID,
			})
		}
	}

	return gpus, nil
}

func (uc *GPUUseCase) extractPodDevices(pod *workload.Pod, checkList map[string]string) ([]PodDevice, error) {
	podDevices, err := device.DecodePodDevices(checkList, pod.Annotations)
	if err != nil {
		return nil, err
	}

	devices := []PodDevice{}

	for _, podDevice := range podDevices {
		for _, containerDevices := range podDevice {
			for _, containerDevice := range containerDevices {
				devices = append(devices, PodDevice{
					GPUID:           containerDevice.UUID,
					UsedCores:       containerDevice.Usedcores,
					UsedMemoryBytes: int64(containerDevice.Usedmem) * 1024, // gigabytes to bytes
				})
			}
		}
	}

	return devices, nil
}

func (uc *GPUUseCase) buildPodRelations(pods []workload.Pod) ([]PodRelation, error) {
	checkList := map[string]string{
		"NVIDIA": hamiVGPUDevicesAllocatedAnnotation,
	}

	relationsPods := []PodRelation{}

	for i := range pods {
		podDevices, err := uc.extractPodDevices(&pods[i], checkList)
		if err != nil {
			return nil, err
		}

		boundAt, _ := unixTimestampStringToTime(pods[i].Annotations[hamiBindTimeAnnotation])

		relationsPods = append(relationsPods, PodRelation{
			Name:         pods[i].Name,
			Namespace:    pods[i].Namespace,
			ModelName:    pods[i].Annotations[model.ModelNameAnnotation],
			BindingPhase: pods[i].Annotations[hamiBindPhaseAnnotation],
			BoundAt:      boundAt,
			PodDevices:   podDevices,
		})
	}

	return relationsPods, nil
}

func unixTimestampStringToTime(str string) (time.Time, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(i, 0), nil
}
