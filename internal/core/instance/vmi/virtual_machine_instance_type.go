package vmi

import (
	"context"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"
)

type (
	// VirtualMachineInstancetype represents a KubeVirt VirtualMachineInstancetype resource.
	VirtualMachineInstancetype = instancetypev1beta1.VirtualMachineInstancetype

	// VirtualMachineClusterInstancetype represents a KubeVirt VirtualMachineClusterInstancetype resource.
	VirtualMachineClusterInstancetype = instancetypev1beta1.VirtualMachineClusterInstancetype
)

type VirtualMachineInstanceType struct {
	*instancetypev1beta1.VirtualMachineInstancetype
	ClusterWide bool
}

type VirtualMachineInstanceTypeRepo interface {
	ListCluster(ctx context.Context, scope, selector string) ([]VirtualMachineClusterInstancetype, error)
	GetCluster(ctx context.Context, scope, name string) (*VirtualMachineClusterInstancetype, error)
	List(ctx context.Context, scope, namespace, selector string) ([]VirtualMachineInstancetype, error)
	Get(ctx context.Context, scope, namespace, name string) (*VirtualMachineInstancetype, error)
	Create(ctx context.Context, scope, namespace string, vmit *VirtualMachineInstancetype) (*VirtualMachineInstancetype, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

func (uc *VirtualMachineInstanceUseCase) ListInstanceTypes(ctx context.Context, scope, namespace string, includeClusterWide bool) ([]VirtualMachineInstanceType, error) {
	vmits, err := uc.virtualMachineInstanceType.List(ctx, scope, namespace, "")
	if err != nil {
		return nil, err
	}

	result := make([]VirtualMachineInstanceType, 0, len(vmits))

	for i := range vmits {
		result = append(result, VirtualMachineInstanceType{
			VirtualMachineInstancetype: &vmits[i],
			ClusterWide:                false,
		})
	}

	if includeClusterWide {
		vmcits, err := uc.virtualMachineInstanceType.ListCluster(ctx, scope, "")
		if err != nil {
			return nil, err
		}

		result = append(result, uc.toVirtualMachineInstanceTypes(vmcits)...)
	}

	return result, nil
}

func (uc *VirtualMachineInstanceUseCase) GetInstanceType(ctx context.Context, scope, namespace, name string) (*VirtualMachineInstanceType, error) {
	vmcit, err := uc.virtualMachineInstanceType.GetCluster(ctx, scope, name)
	if uc.isKeyNotFoundError(err) {
		vmit, err := uc.virtualMachineInstanceType.Get(ctx, scope, namespace, name)
		if err != nil {
			return nil, err
		}

		return &VirtualMachineInstanceType{
			VirtualMachineInstancetype: vmit,
			ClusterWide:                false,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return uc.toVirtualMachineInstanceType(vmcit), nil
}

func (uc *VirtualMachineInstanceUseCase) CreateInstanceType(ctx context.Context, scope, namespace, name string, cpu uint32, memory int64) (*VirtualMachineInstanceType, error) {
	vmit, err := uc.virtualMachineInstanceType.Create(ctx, scope, namespace, uc.buildVirtualMachineInstanceType(namespace, name, cpu, memory))
	if err != nil {
		return nil, err
	}

	return &VirtualMachineInstanceType{
		VirtualMachineInstancetype: vmit,
		ClusterWide:                false,
	}, nil
}

func (uc *VirtualMachineInstanceUseCase) DeleteInstanceType(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineInstanceType.Delete(ctx, scope, namespace, name)
}

func (uc *VirtualMachineInstanceUseCase) buildVirtualMachineInstanceType(namespace, name string, cpu uint32, memory int64) *VirtualMachineInstancetype {
	memoryQuantity := resource.NewQuantity(memory, resource.BinarySI)

	return &instancetypev1beta1.VirtualMachineInstancetype{
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

func (uc *VirtualMachineInstanceUseCase) toVirtualMachineInstanceType(vmcit *VirtualMachineClusterInstancetype) *VirtualMachineInstanceType {
	return &VirtualMachineInstanceType{
		VirtualMachineInstancetype: &VirtualMachineInstancetype{
			TypeMeta:   vmcit.TypeMeta,
			ObjectMeta: vmcit.ObjectMeta,
			Spec:       vmcit.Spec,
		},
		ClusterWide: true,
	}
}

func (uc *VirtualMachineInstanceUseCase) toVirtualMachineInstanceTypes(vmcits []VirtualMachineClusterInstancetype) []VirtualMachineInstanceType {
	result := make([]VirtualMachineInstanceType, 0, len(vmcits))

	for _, vmcit := range vmcits {
		result = append(result, *uc.toVirtualMachineInstanceType(&vmcit))
	}

	return result
}
