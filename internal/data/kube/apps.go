package kube

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type apps struct {
	kubes Kubes
}

func NewApps(kubes Kubes) service.KubeApps {
	return &apps{
		kubes: kubes,
	}
}

var _ service.KubeApps = (*apps)(nil)

func (r *apps) ListDeployments(ctx context.Context, cluster, namespace string) ([]v1.Deployment, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	list, err := client.AppsV1().Deployments(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
