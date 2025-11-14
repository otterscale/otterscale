package vmi

import (
	"context"
	"net/http"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	corev1 "kubevirt.io/api/core/v1"
)

// VirtualMachineInstance represents a KubeVirt VirtualMachineInstance resource.
type VirtualMachineInstance = corev1.VirtualMachineInstance

type VirtualMachineInstanceRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]VirtualMachineInstance, error)
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

func (uc *VirtualMachineInstanceUseCase) isKeyNotFoundError(err error) bool {
	statusErr, _ := err.(*k8serrors.StatusError)
	return statusErr != nil && statusErr.Status().Code == http.StatusNotFound
}
