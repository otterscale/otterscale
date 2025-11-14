package vm

import (
	"context"

	snapshotv1beta1 "kubevirt.io/api/snapshot/v1beta1"
)

// VirtualMachineSnapshot represents a KubeVirt VirtualMachineSnapshot resource.
type VirtualMachineSnapshot = snapshotv1beta1.VirtualMachineSnapshot

type VirtualMachineSnapshotRepo interface {
	List(ctx context.Context, scope, namespace, vmName string) ([]VirtualMachineSnapshot, error)
	Create(ctx context.Context, scope, namespace, name, vmName string) (*VirtualMachineSnapshot, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

func (uc *VirtualMachineUseCase) CreateVirtualMachineSnapshot(ctx context.Context, scope, namespace, name, vmName string) (*VirtualMachineSnapshot, error) {
	return uc.virtualMachineSnapshot.Create(ctx, scope, namespace, name, vmName)
}

func (uc *VirtualMachineUseCase) DeleteVirtualMachineSnapshot(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineSnapshot.Delete(ctx, scope, namespace, name)
}
