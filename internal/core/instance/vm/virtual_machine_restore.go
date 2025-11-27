package vm

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	snapshotv1beta1 "kubevirt.io/api/snapshot/v1beta1"
)

// VirtualMachineRestore represents a KubeVirt VirtualMachineRestore resource.
type VirtualMachineRestore = snapshotv1beta1.VirtualMachineRestore

type VirtualMachineRestoreRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]VirtualMachineRestore, error)
	Create(ctx context.Context, scope, namespace string, vmr *VirtualMachineRestore) (*VirtualMachineRestore, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

func (uc *UseCase) CreateVirtualMachineRestore(ctx context.Context, scope, namespace, name, vmName, snapshot string) (*VirtualMachineRestore, error) {
	return uc.virtualMachineRestore.Create(ctx, scope, namespace, uc.buildVirtualMachineRestore(namespace, name, vmName, snapshot))
}

func (uc *UseCase) DeleteVirtualMachineRestore(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineRestore.Delete(ctx, scope, namespace, name)
}

func (uc *UseCase) buildVirtualMachineRestore(namespace, name, vmName, snapshot string) *VirtualMachineRestore {
	apiGroup := groupName

	return &snapshotv1beta1.VirtualMachineRestore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				nameLabel: vmName,
			},
		},
		Spec: snapshotv1beta1.VirtualMachineRestoreSpec{
			Target: corev1.TypedLocalObjectReference{
				APIGroup: &apiGroup,
				Kind:     kind,
				Name:     vmName,
			},
			VirtualMachineSnapshotName: snapshot,
		},
	}
}
