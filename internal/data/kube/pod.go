package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type pod struct {
	kubes Kubes
}

func NewPod(kubes Kubes) service.KubePod {
	return &pod{
		kubes: kubes,
	}
}

var _ service.KubePod = (*pod)(nil)

func (r *pod) List(ctx context.Context, cluster, namespace string) (*v1.PodList, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	return client.CoreV1().Pods(namespace).List(ctx, opts)
}
