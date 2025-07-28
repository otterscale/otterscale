package kube

import (
	"sync"

	"github.com/openhdc/otterscale/internal/config"
	"k8s.io/client-go/rest"
	"kubevirt.io/client-go/kubecli"
)

type kubevirt struct {
	conf        *config.Config
	virtClients sync.Map
	kube        *Kube
}

func NewKubeVirt(kube *Kube) *kubevirt {
	return &kubevirt{
		kube: kube,
	}
}

func (r *kubevirt) virtClient(config *rest.Config) (kubecli.KubevirtClient, error) {
	if v, ok := r.virtClients.Load(config.Host); ok {
		return v.(kubecli.KubevirtClient), nil
	}

	virtClient, err := kubecli.GetKubevirtClientFromRESTConfig(config)
	if err != nil {
		return nil, err
	}

	r.virtClients.Store(config.Host, virtClient)

	return virtClient, nil
}
