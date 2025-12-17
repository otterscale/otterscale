package vmi

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kvcorev1 "kubevirt.io/api/core/v1"
)

// VirtualMachineInstanceMigration represents a KubeVirt VirtualMachineInstanceMigration resource.
type VirtualMachineInstanceMigration = kvcorev1.VirtualMachineInstanceMigration

type VirtualMachineInstanceMigrationRepo interface {
	Create(ctx context.Context, scope, namespace string, vmim *VirtualMachineInstanceMigration) (*VirtualMachineInstanceMigration, error)
}

func (uc *UseCase) MigrateInstance(ctx context.Context, scope, namespace, name, hostname string) (*VirtualMachineInstanceMigration, error) {
	return uc.virtualMachineInstanceMigration.Create(ctx, scope, namespace, uc.buildVirtualMachineInstanceMigration(namespace, name, hostname))
}

func (uc *UseCase) buildVirtualMachineInstanceMigration(namespace, name, hostname string) *VirtualMachineInstanceMigration {
	return &kvcorev1.VirtualMachineInstanceMigration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-migration-%s", name, uuid.New().String()),
			Namespace: namespace,
		},
		Spec: kvcorev1.VirtualMachineInstanceMigrationSpec{
			VMIName: name,
			AddedNodeSelector: map[string]string{
				"kubernetes.io/hostname": hostname,
			},
		},
	}
}
