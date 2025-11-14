package vmi

import (
	"context"

	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"
)

type VirtualMachineInstanceType struct {
	*instancetypev1beta1.VirtualMachineInstancetype
	ClusterWide bool
}

type VirtualMachineInstanceTypeRepo interface {
	List(ctx context.Context, scope, namespace string, includeClusterWide bool) ([]VirtualMachineInstanceType, error)
	Get(ctx context.Context, scope, namespace, name string) (*VirtualMachineInstanceType, error)
	Create(ctx context.Context, scope, namespace, name string, cpu uint32, memory int64) (*VirtualMachineInstanceType, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

func (uc *VirtualMachineInstanceUseCase) ListInstanceTypes(ctx context.Context, scope, namespace string, includeClusterWide bool) ([]VirtualMachineInstanceType, error) {
	return uc.virtualMachineInstanceType.List(ctx, scope, namespace, includeClusterWide)
}

func (uc *VirtualMachineInstanceUseCase) GetInstanceType(ctx context.Context, scope, namespace, name string) (*VirtualMachineInstanceType, error) {
	return uc.virtualMachineInstanceType.Get(ctx, scope, namespace, name)
}

func (uc *VirtualMachineInstanceUseCase) CreateInstanceType(ctx context.Context, scope, namespace, name string, cpu uint32, memory int64) (*VirtualMachineInstanceType, error) {
	return uc.virtualMachineInstanceType.Create(ctx, scope, namespace, name, cpu, memory)
}

func (uc *VirtualMachineInstanceUseCase) DeleteInstanceType(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineInstanceType.Delete(ctx, scope, namespace, name)
}
