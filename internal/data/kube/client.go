package kube

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type client struct {
	kubes Kubes
}

func NewClient(kubes Kubes) service.KubeClient {
	return &client{
		kubes: kubes,
	}
}

var _ service.KubeClient = (*client)(nil)

func (r *client) Get(cluster string) (*kubernetes.Clientset, error) {
	return r.kubes.Get(cluster)
}

func (r *client) Add(cluster string, cfg *rest.Config) error {
	return r.kubes.Add(cluster, cfg)
}
