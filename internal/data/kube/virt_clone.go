package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	clonev1beta1 "kubevirt.io/api/clone/v1beta1"
	virtcore "kubevirt.io/api/core"

	oscore "github.com/otterscale/otterscale/internal/core"
)

type virtClone struct {
	kube *Kube
}

func NewVirtClone(kube *Kube) oscore.KubeVirtCloneRepo {
	return &virtClone{
		kube: kube,
	}
}

var _ oscore.KubeVirtCloneRepo = (*virtClone)(nil)

func (r *virtClone) ListVirtualMachineClones(ctx context.Context, config *rest.Config, namespace, vmName string) ([]oscore.VirtualMachineClone, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{
		LabelSelector: oscore.VirtualMachineNameLabel + "=" + vmName,
	}
	list, err := clientset.CloneV1beta1().VirtualMachineClones(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *virtClone) CreateVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name, source, target string) (*oscore.VirtualMachineClone, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	apiGroup := virtcore.GroupName
	kind := "VirtualMachine"
	clone := &clonev1beta1.VirtualMachineClone{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				oscore.VirtualMachineNameLabel: source,
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
	opts := metav1.CreateOptions{}
	return clientset.CloneV1beta1().VirtualMachineClones(namespace).Create(ctx, clone, opts)
}

func (r *virtClone) DeleteVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return clientset.CloneV1beta1().VirtualMachineClones(namespace).Delete(ctx, name, opts)
}
