package kube

import (
	"k8s.io/client-go/rest"

	"github.com/openhdc/openhdc/internal/domain/service"
)

type client struct {
	kubeMap KubeMap
	helmMap HelmMap
}

func NewClient(kubeMap KubeMap, helmMap HelmMap) service.KubeClient {
	return &client{
		kubeMap: kubeMap,
		helmMap: helmMap,
	}
}

var _ service.KubeClient = (*client)(nil)

func (r *client) Exists(uuid, facility string) bool {
	return r.kubeMap.exists(uuid, facility) && r.helmMap.exists(uuid, facility)
}

func (r *client) Set(uuid, facility string, config *rest.Config) error {
	if err := r.kubeMap.set(uuid, facility, config); err != nil {
		return err
	}
	return r.helmMap.set(uuid, facility, config)
}
