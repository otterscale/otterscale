package kubevirt

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/instance/vm"
)

type virtualMachineCloneRepo struct {
	kubevirt *KubeVirt
}

func NewVirtualMachineCloneRepo(kubevirt *KubeVirt) vm.VirtualMachineCloneRepo {
	return &virtualMachineCloneRepo{
		kubevirt: kubevirt,
	}
}

var _ vm.VirtualMachineCloneRepo = (*virtualMachineCloneRepo)(nil)

func (r *virtualMachineCloneRepo) List(ctx context.Context, scope, namespace, selector string) ([]vm.VirtualMachineClone, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.CloneV1beta1().VirtualMachineClones(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *virtualMachineCloneRepo) Create(ctx context.Context, scope, namespace string, vmc *vm.VirtualMachineClone) (*vm.VirtualMachineClone, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.CloneV1beta1().VirtualMachineClones(namespace).Create(ctx, vmc, opts)
}

func (r *virtualMachineCloneRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.CloneV1beta1().VirtualMachineClones(namespace).Delete(ctx, name, opts)
}
