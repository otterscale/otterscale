package kube

import (
	"context"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type helm struct {
	kubeMap KubeMap
}

func NewHelm(kubeMap KubeMap) service.KubeHelm {
	return &helm{
		kubeMap: kubeMap,
	}
}

var _ service.KubeHelm = (*helm)(nil)

func (r *helm) ListReleases(ctx context.Context, cluster, namespace string) ([]*release.Release, error) {
	config, err := r.kubeMap.GetHelmConfig(cluster, namespace)
	if err != nil {
		return nil, err
	}
	client := action.NewList(config)
	client.Deployed = true
	return client.Run()
}

func (r *helm) ListRepositories(ctx context.Context, cluster, namespace string) ([]*release.Release, error) {
	config, err := r.kubeMap.GetHelmConfig(cluster, namespace)
	if err != nil {
		return nil, err
	}
	// repo.LoadFile()
	client := action.NewList(config)
	client.Deployed = true
	return client.Run()
}
