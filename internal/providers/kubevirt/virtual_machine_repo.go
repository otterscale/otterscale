package kubevirt

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "kubevirt.io/api/core/v1"

	"github.com/otterscale/otterscale/internal/core/instance/vm"
)

type virtualMachineRepo struct {
	kubevirt *KubeVirt
}

func NewVirtualMachineRepo(kubevirt *KubeVirt) vm.VirtualMachineRepo {
	return &virtualMachineRepo{
		kubevirt: kubevirt,
	}
}

var _ vm.VirtualMachineRepo = (*virtualMachineRepo)(nil)

func (r *virtualMachineRepo) List(ctx context.Context, scope, namespace, selector string) ([]vm.VirtualMachine, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.KubevirtV1().VirtualMachines(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *virtualMachineRepo) Get(ctx context.Context, scope, namespace, name string) (*vm.VirtualMachine, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.KubevirtV1().VirtualMachines(namespace).Get(ctx, name, opts)
}

func (r *virtualMachineRepo) Create(ctx context.Context, scope, namespace string, vm *vm.VirtualMachine) (*vm.VirtualMachine, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.KubevirtV1().VirtualMachines(namespace).Create(ctx, vm, opts)
}

func (r *virtualMachineRepo) Update(ctx context.Context, scope, namespace string, vm *vm.VirtualMachine) (*vm.VirtualMachine, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.KubevirtV1().VirtualMachines(namespace).Update(ctx, vm, opts)
}

func (r *virtualMachineRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.KubevirtV1().VirtualMachines(namespace).Delete(ctx, name, opts)
}

func (r *virtualMachineRepo) Start(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := &corev1.StartOptions{}

	return clientset.KubevirtV1().VirtualMachines(namespace).Start(ctx, name, opts)
}

func (r *virtualMachineRepo) Stop(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := &corev1.StopOptions{}

	return clientset.KubevirtV1().VirtualMachines(namespace).Stop(ctx, name, opts)
}

func (r *virtualMachineRepo) Restart(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := &corev1.RestartOptions{}

	return clientset.KubevirtV1().VirtualMachines(namespace).Restart(ctx, name, opts)
}
