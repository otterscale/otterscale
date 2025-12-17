package kubevirt

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kvcorev1 "kubevirt.io/api/core/v1"

	"github.com/otterscale/otterscale/internal/core/instance/vmi"
)

type virtualMachineInstanceRepo struct {
	kubevirt *KubeVirt
}

func NewVirtualMachineInstanceRepo(kubevirt *KubeVirt) vmi.VirtualMachineInstanceRepo {
	return &virtualMachineInstanceRepo{
		kubevirt: kubevirt,
	}
}

var _ vmi.VirtualMachineInstanceRepo = (*virtualMachineInstanceRepo)(nil)

func (r *virtualMachineInstanceRepo) List(ctx context.Context, scope, namespace, selector string) ([]vmi.VirtualMachineInstance, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.KubevirtV1().VirtualMachineInstances(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *virtualMachineInstanceRepo) Get(ctx context.Context, scope, namespace, name string) (*vmi.VirtualMachineInstance, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.KubevirtV1().VirtualMachineInstances(namespace).Get(ctx, name, opts)
}

func (r *virtualMachineInstanceRepo) Pause(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := &kvcorev1.PauseOptions{}

	return clientset.KubevirtV1().VirtualMachineInstances(namespace).Pause(ctx, name, opts)
}

func (r *virtualMachineInstanceRepo) Resume(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := &kvcorev1.UnpauseOptions{}

	return clientset.KubevirtV1().VirtualMachineInstances(namespace).Unpause(ctx, name, opts)
}
