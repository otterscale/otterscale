package kube

import (
	"context"

	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type storage struct {
	kube *Kube
}

func NewStorage(kube *Kube) service.KubeStorage {
	return &storage{
		kube: kube,
	}
}

var _ service.KubeStorage = (*storage)(nil)

func (r *storage) ListStorageClasses(ctx context.Context, config *rest.Config) ([]v1.StorageClass, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.StorageV1().StorageClasses().List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *storage) GetStorageClass(ctx context.Context, config *rest.Config, name string) (*v1.StorageClass, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return clientset.StorageV1().StorageClasses().Get(ctx, name, opts)
}
