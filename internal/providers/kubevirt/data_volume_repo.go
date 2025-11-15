package kubevirt

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/instance/cdi"
)

type dataVolumeRepo struct {
	kubevirt *KubeVirt
}

func NewDataVolumeRepo(kubevirt *KubeVirt) cdi.DataVolumeRepo {
	return &dataVolumeRepo{
		kubevirt: kubevirt,
	}
}

var _ cdi.DataVolumeRepo = (*dataVolumeRepo)(nil)

func (r *dataVolumeRepo) List(ctx context.Context, scope, namespace, selector string) ([]cdi.DataVolume, error) {
	clientset, err := r.kubevirt.cdiClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.CdiV1beta1().DataVolumes(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *dataVolumeRepo) Get(ctx context.Context, scope, namespace, name string) (*cdi.DataVolume, error) {
	clientset, err := r.kubevirt.cdiClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.CdiV1beta1().DataVolumes(namespace).Get(ctx, name, opts)
}

func (r *dataVolumeRepo) Create(ctx context.Context, scope, namespace string, dv *cdi.DataVolume) (*cdi.DataVolume, error) {
	clientset, err := r.kubevirt.cdiClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.CdiV1beta1().DataVolumes(namespace).Create(ctx, dv, opts)
}

func (r *dataVolumeRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubevirt.cdiClientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.CdiV1beta1().DataVolumes(namespace).Delete(ctx, name, opts)
}
