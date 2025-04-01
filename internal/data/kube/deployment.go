package kube

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type deployment struct {
	kubes Kubes
}

func NewDeployment(kubes Kubes) service.KubeDeployment {
	return &deployment{
		kubes: kubes,
	}
}

var _ service.KubeDeployment = (*deployment)(nil)

func (r *deployment) List(ctx context.Context, cluster, namespace string) (*v1.DeploymentList, error) {
	client, err := r.kubes.Get(cluster)
	if err != nil {
		return nil, err
	}
	opts := metav1.ListOptions{}
	return client.AppsV1().Deployments(namespace).List(ctx, opts)
}
