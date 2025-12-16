package kubevirt

import (
	"context"
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

// filterOutInstanceType checks if an instance type should be filtered out
// Returns true only if:
// 1. The name starts with O, CX, M, N, or RT series, AND
// 2. It has the label "instancetype.kubevirt.io/vendor=kubevirt.io"
func filterOutInstanceType(name string, labels map[string]string) bool {
	excludedPrefixes := []string{"o", "cx", "m", "n", "rt"}
	lowerName := strings.ToLower(name)

	isTargetSeries := false
	for _, prefix := range excludedPrefixes {
		if strings.HasPrefix(lowerName, prefix) {
			if len(lowerName) > len(prefix) {
				nextChar := lowerName[len(prefix)]
				if (nextChar >= '0' && nextChar <= '9') || nextChar == '.' || nextChar == '-' {
					isTargetSeries = true
					break
				}
			}
		}
	}

	if !isTargetSeries {
		return false
	}

	if labels == nil {
		return false
	}
	vendor, exists := labels["instancetype.kubevirt.io/vendor"]
	return exists && vendor == "kubevirt.io"
}

func (r *virtualMachineInstanceTypeRepo) ListCluster(ctx context.Context, scope, selector string) ([]vmi.VirtualMachineClusterInstanceType, error) {
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

	filteredItems := make([]vmi.VirtualMachineClusterInstanceType, 0, len(list.Items))
	for _, item := range list.Items {
		if !filterOutInstanceType(item.Name, item.Labels) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, nil
}

func (r *virtualMachineInstanceTypeRepo) GetCluster(ctx context.Context, scope, name string) (*vmi.VirtualMachineClusterInstanceType, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.InstancetypeV1beta1().VirtualMachineClusterInstancetypes().Get(ctx, name, opts)
}

func (r *virtualMachineInstanceTypeRepo) List(ctx context.Context, scope, namespace, selector string) ([]vmi.VirtualMachineInstanceType, error) {
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

	filteredItems := make([]vmi.VirtualMachineInstanceType, 0, len(list.Items))
	for _, item := range list.Items {
		if !filterOutInstanceType(item.Name, item.Labels) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, nil
}

func (r *virtualMachineInstanceTypeRepo) Get(ctx context.Context, scope, namespace, name string) (*vmi.VirtualMachineInstanceType, error) {
	clientset, err := r.kubevirt.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.InstancetypeV1beta1().VirtualMachineInstancetypes(namespace).Get(ctx, name, opts)
}

func (r *virtualMachineInstanceTypeRepo) Create(ctx context.Context, scope, namespace string, vmit *vmi.VirtualMachineInstanceType) (*vmi.VirtualMachineInstanceType, error) {
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
