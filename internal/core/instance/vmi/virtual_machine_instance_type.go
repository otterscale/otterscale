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
	List(ctx context.Context, scope, selector string) ([]VirtualMachineClusterInstanceType, error)
	Get(ctx context.Context, scope, name string) (*VirtualMachineClusterInstanceType, error)
	Create(ctx context.Context, scope string, vmit *VirtualMachineClusterInstanceType) (*VirtualMachineClusterInstanceType, error)
	Delete(ctx context.Context, scope, name string) error
}

func (uc *UseCase) ListInstanceTypes(ctx context.Context, scope string) ([]VirtualMachineInstanceTypeData, error) {
	vmcits, err := uc.virtualMachineInstanceType.List(ctx, scope, "")
	if err != nil {
		return nil, err
	}

	return uc.toVirtualMachineInstanceTypes(vmcits), nil
}

func (uc *UseCase) GetInstanceType(ctx context.Context, scope, name string) (*VirtualMachineInstanceTypeData, error) {
	vmcit, err := uc.virtualMachineInstanceType.Get(ctx, scope, name)
	if err != nil {
		return nil, err
	}

	return uc.toVirtualMachineInstanceType(vmcit), nil
}

func (uc *UseCase) CreateInstanceType(ctx context.Context, scope, name string, cpu uint32, memory int64) (*VirtualMachineInstanceTypeData, error) {
	vmcit, err := uc.virtualMachineInstanceType.Create(ctx, scope, uc.buildVirtualMachineClusterInstanceType(name, cpu, memory))
	if err != nil {
		return nil, err
	}

	return uc.toVirtualMachineInstanceType(vmcit), nil
}

func (uc *UseCase) DeleteInstanceType(ctx context.Context, scope, name string) error {
	return uc.virtualMachineInstanceType.Delete(ctx, scope, name)
}

func (uc *UseCase) buildVirtualMachineClusterInstanceType(name string, cpu uint32, memory int64) *VirtualMachineClusterInstanceType {
	memoryQuantity := resource.NewQuantity(memory, resource.BinarySI)

	return &VirtualMachineClusterInstanceType{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
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

func (uc *UseCase) toVirtualMachineInstanceType(vmcit *VirtualMachineClusterInstanceType) *VirtualMachineInstanceTypeData {
	return &VirtualMachineInstanceTypeData{
		Type: &VirtualMachineInstanceType{
			TypeMeta:   vmcit.TypeMeta,
			ObjectMeta: vmcit.ObjectMeta,
			Spec:       vmcit.Spec,
		},
		ClusterWide: true,
	}
}

func (uc *UseCase) toVirtualMachineInstanceTypes(vmcits []VirtualMachineClusterInstanceType) []VirtualMachineInstanceTypeData {
	ret := make([]VirtualMachineInstanceTypeData, 0, len(vmcits))

	for i := range vmcits {
		ret = append(ret, *uc.toVirtualMachineInstanceType(&vmcits[i]))
	}

	return ret
}
