package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type persistentVolumeClaim struct {
	kubes Kubes
}

func NewPersistentVolumeClaim(kubes Kubes) service.KubePersistentVolumeClaim {
	return &persistentVolumeClaim{
		kubes: kubes,
	}
}

var _ service.KubePersistentVolumeClaim = (*persistentVolumeClaim)(nil)

func (r *persistentVolumeClaim) List(ctx context.Context, cluster, namespace string) (*v1.PersistentVolumeClaimList, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	return client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, opts)
}
