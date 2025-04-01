package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type svc struct {
	kubes Kubes
}

func NewService(kubes Kubes) service.KubeSVC {
	return &svc{
		kubes: kubes,
	}
}

var _ service.KubeSVC = (*svc)(nil)

func (r *svc) List(ctx context.Context, cluster, namespace string) (*v1.ServiceList, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	return client.CoreV1().Services(namespace).List(ctx, opts)
}
