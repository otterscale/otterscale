package vm

import (
	"context"

	snapshotv1beta1 "kubevirt.io/api/snapshot/v1beta1"
)

// VirtualMachineRestore represents a KubeVirt VirtualMachineRestore resource.
type VirtualMachineRestore = snapshotv1beta1.VirtualMachineRestore

type VirtualMachineRestoreRepo interface {
	List(ctx context.Context, scope, namespace, vmName string) ([]VirtualMachineRestore, error)
	Create(ctx context.Context, scope, namespace, name, vmName, snapshot string) (*VirtualMachineRestore, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

func (uc *VirtualMachineUseCase) CreateVirtualMachineRestore(ctx context.Context, scope, namespace, name, vmName, snapshot string) (*VirtualMachineRestore, error) {
	return uc.virtualMachineRestore.Create(ctx, scope, namespace, name, vmName, snapshot)
}

func (uc *VirtualMachineUseCase) DeleteVirtualMachineRestore(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineRestore.Delete(ctx, scope, namespace, name)
}
