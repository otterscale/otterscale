package vmi

import (
	"context"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"
)

type (
	// VirtualMachineInstanceType represents a KubeVirt VirtualMachineInstancetype resource.
	VirtualMachineInstanceType = instancetypev1beta1.VirtualMachineInstancetype

	// VirtualMachineClusterInstancetype represents a KubeVirt VirtualMachineClusterInstancetype resource.
	VirtualMachineClusterInstanceType = instancetypev1beta1.VirtualMachineClusterInstancetype
)

type VirtualMachineInstanceTypeData struct {
	Type        *VirtualMachineInstanceType
	ClusterWide bool
}

type VirtualMachineInstanceTypeRepo interface {
	ListCluster(ctx context.Context, scope, selector string) ([]VirtualMachineClusterInstanceType, error)
	GetCluster(ctx context.Context, scope, name string) (*VirtualMachineClusterInstanceType, error)
	List(ctx context.Context, scope, namespace, selector string) ([]VirtualMachineInstanceType, error)
	Get(ctx context.Context, scope, namespace, name string) (*VirtualMachineInstanceType, error)
	Create(ctx context.Context, scope, namespace string, vmit *VirtualMachineInstanceType) (*VirtualMachineInstanceType, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

func (uc *VirtualMachineInstanceUseCase) ListInstanceTypes(ctx context.Context, scope, namespace string, includeClusterWide bool) ([]VirtualMachineInstanceTypeData, error) {
	vmits, err := uc.virtualMachineInstanceType.List(ctx, scope, namespace, "")
	if err != nil {
		return nil, err
	}

	ret := make([]VirtualMachineInstanceTypeData, 0, len(vmits))

	for i := range vmits {
		ret = append(ret, VirtualMachineInstanceTypeData{
			Type:        &vmits[i],
			ClusterWide: false,
		})
	}

	if includeClusterWide {
		vmcits, err := uc.virtualMachineInstanceType.ListCluster(ctx, scope, "")
		if err != nil {
			return nil, err
		}

		ret = append(ret, uc.toVirtualMachineInstanceTypes(vmcits)...)
	}

	return ret, nil
}

func (uc *VirtualMachineInstanceUseCase) GetInstanceType(ctx context.Context, scope, namespace, name string) (*VirtualMachineInstanceTypeData, error) {
	vmcit, err := uc.virtualMachineInstanceType.GetCluster(ctx, scope, name)
	if uc.isKeyNotFoundError(err) {
		vmit, err := uc.virtualMachineInstanceType.Get(ctx, scope, namespace, name)
		if err != nil {
			return nil, err
		}

		return &VirtualMachineInstanceTypeData{
			Type:        vmit,
			ClusterWide: false,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return uc.toVirtualMachineInstanceType(vmcit), nil
}

func (uc *VirtualMachineInstanceUseCase) CreateInstanceType(ctx context.Context, scope, namespace, name string, cpu uint32, memory int64) (*VirtualMachineInstanceTypeData, error) {
	vmit, err := uc.virtualMachineInstanceType.Create(ctx, scope, namespace, uc.buildVirtualMachineInstanceType(namespace, name, cpu, memory))
	if err != nil {
		return nil, err
	}

	return &VirtualMachineInstanceTypeData{
		Type:        vmit,
		ClusterWide: false,
	}, nil
}

func (uc *VirtualMachineInstanceUseCase) DeleteInstanceType(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineInstanceType.Delete(ctx, scope, namespace, name)
}

func (uc *VirtualMachineInstanceUseCase) buildVirtualMachineInstanceType(namespace, name string, cpu uint32, memory int64) *VirtualMachineInstanceType {
	memoryQuantity := resource.NewQuantity(memory, resource.BinarySI)

	return &VirtualMachineInstanceType{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: instancetypev1beta1.VirtualMachineInstancetypeSpec{
			CPU: instancetypev1beta1.CPUInstancetype{
				Guest: cpu,
			},
			Memory: instancetypev1beta1.MemoryInstancetype{
				Guest: *memoryQuantity,
			},
		},
	}
}

func (uc *VirtualMachineInstanceUseCase) toVirtualMachineInstanceType(vmcit *VirtualMachineClusterInstanceType) *VirtualMachineInstanceTypeData {
	return &VirtualMachineInstanceTypeData{
		Type: &VirtualMachineInstanceType{
			TypeMeta:   vmcit.TypeMeta,
			ObjectMeta: vmcit.ObjectMeta,
			Spec:       vmcit.Spec,
		},
		ClusterWide: true,
	}
}

func (uc *VirtualMachineInstanceUseCase) toVirtualMachineInstanceTypes(vmcits []VirtualMachineClusterInstanceType) []VirtualMachineInstanceTypeData {
	ret := make([]VirtualMachineInstanceTypeData, 0, len(vmcits))

	for i := range vmcits {
		ret = append(ret, *uc.toVirtualMachineInstanceType(&vmcits[i]))
	}

	return ret
}
