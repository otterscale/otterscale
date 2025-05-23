package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	oscore "github.com/openhdc/otterscale/internal/core"
)

type storage struct {
	kube *Kube
}

func NewStorage(kube *Kube) oscore.KubeStorageRepo {
	return &storage{
		kube: kube,
	}
}

var _ oscore.KubeStorageRepo = (*storage)(nil)

func (r *storage) ListStorageClasses(ctx context.Context, config *rest.Config) ([]oscore.StorageClass, error) {
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

func (r *storage) GetStorageClass(ctx context.Context, config *rest.Config, name string) (*oscore.StorageClass, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return clientset.StorageV1().StorageClasses().Get(ctx, name, opts)
}
