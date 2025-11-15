package kubevirt

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/instance/vm"
)

type virtualMachineRestoreRepo struct {
	kubevirt *KubeVirt
}

func NewVirtualMachineRestoreRepo(kubevirt *KubeVirt) vm.VirtualMachineRestoreRepo {
	return &virtualMachineRestoreRepo{
		kubevirt: kubevirt,
	}
}

var _ vm.VirtualMachineRestoreRepo = (*virtualMachineRestoreRepo)(nil)

func (r *virtualMachineRestoreRepo) List(ctx context.Context, scope, namespace, selector string) ([]vm.VirtualMachineRestore, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.SnapshotV1beta1().VirtualMachineRestores(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *virtualMachineRestoreRepo) Create(ctx context.Context, scope, namespace string, vmr *vm.VirtualMachineRestore) (*vm.VirtualMachineRestore, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.SnapshotV1beta1().VirtualMachineRestores(namespace).Create(ctx, vmr, opts)
}

func (r *virtualMachineRestoreRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.SnapshotV1beta1().VirtualMachineRestores(namespace).Delete(ctx, name, opts)
}
