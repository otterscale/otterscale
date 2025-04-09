package kube

import (
	"k8s.io/client-go/rest"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type client struct {
	kubeMap KubeMap
}

func NewClient(kubeMap KubeMap) service.KubeClient {
	return &client{
		kubeMap: kubeMap,
	}
}

var _ service.KubeClient = (*client)(nil)

func (r *client) Exists(cluster string) bool {
	_, ok := r.kubeMap[cluster]
	return ok
}

func (r *client) Add(cluster string, cfg *rest.Config) error {
	return r.kubeMap.Add(cluster, cfg)
}
