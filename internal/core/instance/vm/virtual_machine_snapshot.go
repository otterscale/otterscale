package vm

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	snapshotv1beta1 "kubevirt.io/api/snapshot/v1beta1"
)

type (
	// VirtualMachineSnapshot represents a KubeVirt VirtualMachineSnapshot resource.
	VirtualMachineSnapshot = snapshotv1beta1.VirtualMachineSnapshot

	// VirtualMachineSnapshotCondition represents a KubeVirt Condition resource.
	VirtualMachineSnapshotCondition = snapshotv1beta1.Condition
)

type VirtualMachineSnapshotRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]VirtualMachineSnapshot, error)
	Create(ctx context.Context, scope, namespace string, vms *VirtualMachineSnapshot) (*VirtualMachineSnapshot, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

func (uc *UseCase) CreateVirtualMachineSnapshot(ctx context.Context, scope, namespace, name, vmName string) (*VirtualMachineSnapshot, error) {
	return uc.virtualMachineSnapshot.Create(ctx, scope, namespace, uc.buildVirtualMachineSnapshot(namespace, name, vmName))
}

func (uc *UseCase) DeleteVirtualMachineSnapshot(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineSnapshot.Delete(ctx, scope, namespace, name)
}

func (uc *UseCase) buildVirtualMachineSnapshot(namespace, name, vmName string) *VirtualMachineSnapshot {
	apiGroup := groupName

	return &snapshotv1beta1.VirtualMachineSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				nameLabel: vmName,
			},
		},
		Spec: snapshotv1beta1.VirtualMachineSnapshotSpec{
			Source: corev1.TypedLocalObjectReference{
				APIGroup: &apiGroup,
				Kind:     kind,
				Name:     vmName,
			},
		},
	}
}
