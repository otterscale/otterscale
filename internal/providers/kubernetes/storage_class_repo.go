package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/persistent"
)

type storageClassRepo struct {
	kubernetes *Kubernetes
}

func NewStorageClassRepo(kubernetes *Kubernetes) persistent.StorageClassRepo {
	return &storageClassRepo{
		kubernetes: kubernetes,
	}
}

var _ persistent.StorageClassRepo = (*storageClassRepo)(nil)

func (r *storageClassRepo) List(ctx context.Context, scope, selector string) ([]persistent.StorageClass, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.StorageV1().StorageClasses().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *storageClassRepo) Get(ctx context.Context, scope, name string) (*persistent.StorageClass, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.StorageV1().StorageClasses().Get(ctx, name, opts)
}
