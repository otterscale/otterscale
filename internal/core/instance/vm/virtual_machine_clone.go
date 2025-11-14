package vm

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clonev1beta1 "kubevirt.io/api/clone/v1beta1"
)

// VirtualMachineClone represents a KubeVirt VirtualMachineClone resource.
type VirtualMachineClone = clonev1beta1.VirtualMachineClone

type VirtualMachineCloneRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]VirtualMachineClone, error)
	Create(ctx context.Context, scope, namespace string, vmc *VirtualMachineClone) (*VirtualMachineClone, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

func (uc *VirtualMachineUseCase) CreateVirtualMachineClone(ctx context.Context, scope, namespace, name, source, target string) (*VirtualMachineClone, error) {
	return uc.virtualMachineClone.Create(ctx, scope, namespace, uc.buildVirtualMachineClone(namespace, name, source, target))
}

func (uc *VirtualMachineUseCase) DeleteVirtualMachineClone(ctx context.Context, scope, namespace, name string) error {
	return uc.virtualMachineClone.Delete(ctx, scope, namespace, name)
}

func (uc *VirtualMachineUseCase) buildVirtualMachineClone(namespace, name, source, target string) *VirtualMachineClone {
	apiGroup := groupName

	return &clonev1beta1.VirtualMachineClone{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				VirtualMachineNameLabel: source,
			},
		},
		Spec: clonev1beta1.VirtualMachineCloneSpec{
			Source: &corev1.TypedLocalObjectReference{
				APIGroup: &apiGroup,
				Kind:     kind,
				Name:     source,
			},
			Target: &corev1.TypedLocalObjectReference{
				APIGroup: &apiGroup,
				Kind:     kind,
				Name:     target,
			},
		},
	}
}
