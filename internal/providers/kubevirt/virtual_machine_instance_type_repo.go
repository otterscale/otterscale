package kubevirt

import (
	"context"
	"slices"
	"strings"

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

func (r *virtualMachineInstanceTypeRepo) filterCluster(cits []vmi.VirtualMachineClusterInstanceType) []vmi.VirtualMachineClusterInstanceType {
	excludedPrefixes := []string{"o", "cx", "m", "n", "rt"}

	return slices.DeleteFunc(cits, func(cit vmi.VirtualMachineClusterInstanceType) bool {
		for _, excludedPrefix := range excludedPrefixes {
			if strings.HasPrefix(cit.Name, excludedPrefix) {
				return true
			}
		}
		return false
	})
}

func (r *virtualMachineInstanceTypeRepo) List(ctx context.Context, scope, selector string) ([]vmi.VirtualMachineClusterInstanceType, error) {
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

	return r.filterCluster(list.Items), nil
}

func (r *virtualMachineInstanceTypeRepo) Get(ctx context.Context, scope, name string) (*vmi.VirtualMachineClusterInstanceType, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.InstancetypeV1beta1().VirtualMachineClusterInstancetypes().Get(ctx, name, opts)
}

func (r *virtualMachineInstanceTypeRepo) Create(ctx context.Context, scope string, vmit *vmi.VirtualMachineClusterInstanceType) (*vmi.VirtualMachineClusterInstanceType, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.InstancetypeV1beta1().VirtualMachineClusterInstancetypes().Create(ctx, vmit, opts)
}

func (r *virtualMachineInstanceTypeRepo) Delete(ctx context.Context, scope, name string) error {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.InstancetypeV1beta1().VirtualMachineClusterInstancetypes().Delete(ctx, name, opts)
}
