package kubevirt

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/instance/vmi"
)

type virtualMachineInstanceTypeRepo struct {
	kubevirt *KubeVirt
}

func NewVirtualMachineInstanceTypeRepo(kubevirt *KubeVirt) vmi.VirtualMachineInstanceTypeRepo {
	return &virtualMachineInstanceTypeRepo{
		kubevirt: kubevirt,
	}
}

var _ vmi.VirtualMachineInstanceTypeRepo = (*virtualMachineInstanceTypeRepo)(nil)

func (r *virtualMachineInstanceTypeRepo) ListCluster(ctx context.Context, scope, selector string) ([]vmi.VirtualMachineClusterInstancetype, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.InstancetypeV1beta1().VirtualMachineClusterInstancetypes().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *virtualMachineInstanceTypeRepo) GetCluster(ctx context.Context, scope, name string) (*vmi.VirtualMachineClusterInstancetype, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.InstancetypeV1beta1().VirtualMachineClusterInstancetypes().Get(ctx, name, opts)
}

func (r *virtualMachineInstanceTypeRepo) List(ctx context.Context, scope, namespace, selector string) ([]vmi.VirtualMachineInstancetype, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.InstancetypeV1beta1().VirtualMachineInstancetypes(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *virtualMachineInstanceTypeRepo) Get(ctx context.Context, scope, namespace, name string) (*vmi.VirtualMachineInstancetype, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.InstancetypeV1beta1().VirtualMachineInstancetypes(namespace).Get(ctx, name, opts)
}

func (r *virtualMachineInstanceTypeRepo) Create(ctx context.Context, scope, namespace string, vmit *vmi.VirtualMachineInstancetype) (*vmi.VirtualMachineInstancetype, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.InstancetypeV1beta1().VirtualMachineInstancetypes(namespace).Create(ctx, vmit, opts)
}

func (r *virtualMachineInstanceTypeRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.InstancetypeV1beta1().VirtualMachineInstancetypes(namespace).Delete(ctx, name, opts)
}
