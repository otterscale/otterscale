package core

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/Project-HAMi/HAMi/pkg/device"
)

type GPURelations struct {
	Machines []Machine
	GPUs     []GPURelationsGPU
	Pods     []GPURelationsPod
}

type GPURelationsGPU struct {
	ID          string
	Index       uint32
	Count       int32
	Cores       int32
	MemoryBytes int64
	Type        string
	Health      bool
	MachineID   string
}

type GPURelationsPod struct {
	Name         string
	Namespace    string
	ModelName    string
	BindingPhase string
	BoundAt      time.Time
	PodDevices   []GPURelationPodDevice
}

type GPURelationPodDevice struct {
	GPUID           string
	UsedCores       int32
	UsedMemoryBytes int64
}

func (uc *EssentialUseCase) ListGPURelationsByMachine(ctx context.Context, scopeUUID, facilityName, machineID string) (*GPURelations, error) {
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return nil, err
	}
	labelSelector := fmt.Sprintf("%s=%s", annotationHAMIVGPUNode, machine.Hostname)
	return uc.listGPURelations(ctx, scopeUUID, facilityName, "", labelSelector)
}

func (uc *EssentialUseCase) ListGPURelationsByModel(ctx context.Context, scopeUUID, facilityName, namespace, modelName string) (*GPURelations, error) {
	labelSelector := fmt.Sprintf("%s=%s", ApplicationReleaseLLMDModelNameLabel, modelName)
	return uc.listGPURelations(ctx, scopeUUID, facilityName, namespace, labelSelector)
}

func (uc *EssentialUseCase) listGPURelations(ctx context.Context, scopeUUID, facilityName, namespace, labelSelector string) (*GPURelations, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scopeUUID, facilityName)
	if err != nil {
		return nil, err
	}

	pods, err := uc.kubeCore.ListPodsByLabel(ctx, config, namespace, labelSelector)
	if err != nil {
		return nil, err
	}

	nodes, err := uc.kubeCore.ListNodes(ctx, config)
	if err != nil {
		return nil, err
	}

	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}

	return uc.buildGPURelations(pods, nodes, machines)
}

func (uc *EssentialUseCase) buildGPURelations(pods []Pod, nodes []Node, machines []Machine) (*GPURelations, error) {
	nodeNames := extractNodeNamesFromPods(pods)
	filteredMachines := filterMachinesByNodeNames(machines, nodeNames)
	filteredNodes := filterNodesByNames(nodes, nodeNames)
	machineMap := buildMachineMap(filteredMachines)

	gpus, err := uc.buildGPUsFromNodes(filteredNodes, machineMap)
	if err != nil {
		return nil, err
	}

	relationsPods, err := buildRelationPodsFromPods(pods)
	if err != nil {
		return nil, err
	}

	return &GPURelations{
		Machines: filteredMachines,
		GPUs:     gpus,
		Pods:     relationsPods,
	}, nil
}

func extractNodeNamesFromPods(pods []Pod) []string {
	nodeNames := make([]string, 0, len(pods))
	for i := range pods {
		if nodeName, ok := pods[i].Annotations[annotationHAMIVGPUNode]; ok {
			nodeNames = append(nodeNames, nodeName)
		}
	}
	return nodeNames
}

func filterMachinesByNodeNames(machines []Machine, nodeNames []string) []Machine {
	return slices.DeleteFunc(machines, func(m Machine) bool {
		return !slices.Contains(nodeNames, m.Hostname)
	})
}

func filterNodesByNames(nodes []Node, nodeNames []string) []Node {
	return slices.DeleteFunc(nodes, func(n Node) bool {
		return !slices.Contains(nodeNames, n.Name)
	})
}

func buildMachineMap(machines []Machine) map[string]Machine {
	machineMap := make(map[string]Machine, len(machines))
	for i := range machines {
		machineMap[machines[i].Hostname] = machines[i]
	}
	return machineMap
}

func (uc *EssentialUseCase) buildGPUsFromNodes(nodes []Node, machineMap map[string]Machine) ([]GPURelationsGPU, error) {
	gpus := []GPURelationsGPU{}
	for i := range nodes {
		annotation, ok := nodes[i].Annotations[annotationHAMINodeNvidiaRegister]
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
			gpus = append(gpus, GPURelationsGPU{
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

func buildRelationPodsFromPods(pods []Pod) ([]GPURelationsPod, error) {
	checkList := map[string]string{
		"NVIDIA": annotationHAMIVGPUDevicesAllocated,
	}

	relationsPods := make([]GPURelationsPod, 0, len(pods))
	for i := range pods {
		podDevices, err := extractPodDevices(&pods[i], checkList)
		if err != nil {
			return nil, err
		}
		boundAt, _ := unixTimestampStringToTime(pods[i].Annotations[annotationHAMIBindTime])
		relationsPods = append(relationsPods, GPURelationsPod{
			Name:         pods[i].Name,
			Namespace:    pods[i].Namespace,
			ModelName:    pods[i].Labels[ApplicationReleaseLLMDModelNameLabel],
			BindingPhase: pods[i].Annotations[annotationHAMIBindPhase],
			BoundAt:      boundAt,
			PodDevices:   podDevices,
		})
	}
	return relationsPods, nil
}

func extractPodDevices(pod *Pod, checkList map[string]string) ([]GPURelationPodDevice, error) {
	podDevices, err := device.DecodePodDevices(checkList, pod.Annotations)
	if err != nil {
		return nil, err
	}
	devices := []GPURelationPodDevice{}
	for _, podDevice := range podDevices {
		for _, containerDevices := range podDevice {
			for _, containerDevice := range containerDevices {
				devices = append(devices, GPURelationPodDevice{
					GPUID:           containerDevice.UUID,
					UsedCores:       containerDevice.Usedcores,
					UsedMemoryBytes: int64(containerDevice.Usedmem) * 1024, // gigabytes to bytes
				})
			}
		}
	}
	return devices, nil
}

func unixTimestampStringToTime(str string) (time.Time, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0), nil
}
