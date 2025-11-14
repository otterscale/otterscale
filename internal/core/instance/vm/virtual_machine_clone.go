package vm

import (
	"context"

	clonev1beta1 "kubevirt.io/api/clone/v1beta1"
)

// VirtualMachineClone represents a KubeVirt VirtualMachineClone resource.
type VirtualMachineClone = clonev1beta1.VirtualMachineClone

type VirtualMachineCloneRepo interface {
	List(ctx context.Context, scope, namespace, vmName string) ([]VirtualMachineClone, error)
	Create(ctx context.Context, scope, namespace, name, source, target string) (*VirtualMachineClone, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

func (uc *VirtualMachineUseCase) CreateVirtualMachineClone(ctx context.Context, scope, namespace, name, source, target string) (*VirtualMachineClone, error) {
	return uc.virtualMachineClone.Create(ctx, scope, namespace, name, source, target)
}

func (uc *VirtualMachineUseCase) DeleteVirtualMachineClone(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineClone.Delete(ctx, scope, namespace, name)
}
