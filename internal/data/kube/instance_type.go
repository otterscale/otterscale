package kube

import (
	"context"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"kubevirt.io/api/instancetype/v1beta1"

	oscore "github.com/otterscale/otterscale/internal/core"
)

type instanceType struct {
	kube *Kube
}

func NewInstanceType(kube *Kube) oscore.KubeInstanceTypeRepo {
	return &instanceType{
		kube: kube,
	}
}

var _ oscore.KubeInstanceTypeRepo = (*instanceType)(nil)

func (r *instanceType) ListClusterWide(ctx context.Context, config *rest.Config) ([]oscore.VirtualMachineClusterInstanceType, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := clientset.InstancetypeV1beta1().VirtualMachineClusterInstancetypes().List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *instanceType) List(ctx context.Context, config *rest.Config, namespace string) ([]oscore.VirtualMachineInstanceType, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := clientset.InstancetypeV1beta1().VirtualMachineInstancetypes(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *instanceType) Get(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.VirtualMachineInstanceType, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	opts := metav1.GetOptions{}
	return clientset.InstancetypeV1beta1().VirtualMachineInstancetypes(namespace).Get(ctx, name, opts)
}

func (r *instanceType) Create(ctx context.Context, config *rest.Config, namespace, name string, cpu uint32, memory int64) (*oscore.VirtualMachineInstanceType, error) {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return nil, err
	}
	memoryQuantity := resource.NewQuantity(memory, resource.BinarySI)
	instanceType := &v1beta1.VirtualMachineInstancetype{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1beta1.VirtualMachineInstancetypeSpec{
			CPU: v1beta1.CPUInstancetype{
				Guest: cpu,
			},
			Memory: v1beta1.MemoryInstancetype{
				Guest: *memoryQuantity,
			},
		},
	}
	opts := metav1.CreateOptions{}
	return clientset.InstancetypeV1beta1().VirtualMachineInstancetypes(namespace).Create(ctx, instanceType, opts)
}

func (r *instanceType) Delete(ctx context.Context, config *rest.Config, namespace, name string) error {
	clientset, err := r.kube.virtClientset(config)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return clientset.InstancetypeV1beta1().VirtualMachineInstancetypes(namespace).Delete(ctx, name, opts)
}
