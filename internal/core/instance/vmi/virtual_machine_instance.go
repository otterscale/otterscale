package vmi

import (
	"context"

	corev1 "kubevirt.io/api/core/v1"
)

// VirtualMachineInstance represents a KubeVirt VirtualMachineInstance resource.
type VirtualMachineInstance = corev1.VirtualMachineInstance

type VirtualMachineInstanceRepo interface {
	List(ctx context.Context, scope, namespace string) ([]VirtualMachineInstance, error)
	Get(ctx context.Context, scope, namespace, name string) (*VirtualMachineInstance, error)
	Pause(ctx context.Context, scope, namespace, name string) error
	Resume(ctx context.Context, scope, namespace, name string) error
}

type VirtualMachineInstanceUseCase struct {
	virtualMachineInstance          VirtualMachineInstanceRepo
	virtualMachineInstanceType      VirtualMachineInstanceTypeRepo
	virtualMachineInstanceMigration VirtualMachineInstanceMigrationRepo
}

func NewVirtualMachineInstanceUseCase(virtualMachineInstance VirtualMachineInstanceRepo, virtualMachineInstanceType VirtualMachineInstanceTypeRepo, virtualMachineInstanceMigration VirtualMachineInstanceMigrationRepo) *VirtualMachineInstanceUseCase {
	return &VirtualMachineInstanceUseCase{
		virtualMachineInstance:          virtualMachineInstance,
		virtualMachineInstanceType:      virtualMachineInstanceType,
		virtualMachineInstanceMigration: virtualMachineInstanceMigration,
	}
}

func (uc *VirtualMachineInstanceUseCase) PauseInstance(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineInstance.Pause(ctx, scope, namespace, name)
}

func (uc *VirtualMachineInstanceUseCase) ResumeInstance(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineInstance.Resume(ctx, scope, namespace, name)
}
