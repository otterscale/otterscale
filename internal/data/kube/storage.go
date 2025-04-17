package kube

import (
	"context"

	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type storage struct {
	kubeMap KubeMap
}

func NewStorage(kubeMap KubeMap) service.KubeStorage {
	return &storage{
		kubeMap: kubeMap,
	}
}

var _ service.KubeStorage = (*storage)(nil)

func (r *storage) ListStorageClasses(ctx context.Context, uuid, facility string) ([]v1.StorageClass, error) {
	clientset, err := r.kubeMap.get(uuid, facility)
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
