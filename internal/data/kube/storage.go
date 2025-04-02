package kube

import (
	"context"

	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type storage struct {
	kubes Kubes
}

func NewStorage(kubes Kubes) service.KubeStorage {
	return &storage{
		kubes: kubes,
	}
}

var _ service.KubeStorage = (*storage)(nil)

func (r *storage) ListStorageClasses(ctx context.Context, cluster string) ([]v1.StorageClass, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := client.StorageV1().StorageClasses().List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
