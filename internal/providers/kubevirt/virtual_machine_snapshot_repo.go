package kubevirt

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/instance/vm"
)

type virtualMachineSnapshotRepo struct {
	kubevirt *KubeVirt
}

func NewVirtualMachineSnapshotRepo(kubevirt *KubeVirt) vm.VirtualMachineSnapshotRepo {
	return &virtualMachineSnapshotRepo{
		kubevirt: kubevirt,
	}
}

var _ vm.VirtualMachineSnapshotRepo = (*virtualMachineSnapshotRepo)(nil)

func (r *virtualMachineSnapshotRepo) List(ctx context.Context, scope, namespace, selector string) ([]vm.VirtualMachineSnapshot, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.SnapshotV1beta1().VirtualMachineSnapshots(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *virtualMachineSnapshotRepo) Create(ctx context.Context, scope, namespace string, vms *vm.VirtualMachineSnapshot) (*vm.VirtualMachineSnapshot, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.SnapshotV1beta1().VirtualMachineSnapshots(namespace).Create(ctx, vms, opts)
}

func (r *virtualMachineSnapshotRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.SnapshotV1beta1().VirtualMachineSnapshots(namespace).Delete(ctx, name, opts)
}
