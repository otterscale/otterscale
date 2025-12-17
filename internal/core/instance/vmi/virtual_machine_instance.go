package vmi

import (
	"context"

	kvcorev1 "kubevirt.io/api/core/v1"
)

// VirtualMachineInstance represents a KubeVirt VirtualMachineInstance resource.
type VirtualMachineInstance = kvcorev1.VirtualMachineInstance

type VirtualMachineInstanceRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]VirtualMachineInstance, error)
	Get(ctx context.Context, scope, namespace, name string) (*VirtualMachineInstance, error)
	Pause(ctx context.Context, scope, namespace, name string) error
	Resume(ctx context.Context, scope, namespace, name string) error
}

type UseCase struct {
	virtualMachineInstance          VirtualMachineInstanceRepo
	virtualMachineInstanceType      VirtualMachineInstanceTypeRepo
	virtualMachineInstanceMigration VirtualMachineInstanceMigrationRepo
}

func NewUseCase(virtualMachineInstance VirtualMachineInstanceRepo, virtualMachineInstanceType VirtualMachineInstanceTypeRepo, virtualMachineInstanceMigration VirtualMachineInstanceMigrationRepo) *UseCase {
	return &UseCase{
		virtualMachineInstance:          virtualMachineInstance,
		virtualMachineInstanceType:      virtualMachineInstanceType,
		virtualMachineInstanceMigration: virtualMachineInstanceMigration,
	}
}

func (uc *UseCase) PauseInstance(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineInstance.Pause(ctx, scope, namespace, name)
}

func (uc *UseCase) ResumeInstance(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineInstance.Resume(ctx, scope, namespace, name)
}
