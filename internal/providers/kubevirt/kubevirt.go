package kubevirt

import (
	"sync"

	"github.com/otterscale/kubevirt-client-go/containerizeddataimporter"
	"github.com/otterscale/kubevirt-client-go/kubevirt"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
)

type KubeVirt struct {
	conf       *config.Config
	kubernetes *kubernetes.Kubernetes

	kvClientsets  sync.Map
	cdiClientsets sync.Map
}

func New(conf *config.Config, kubernetes *kubernetes.Kubernetes) (*KubeVirt, error) {
	return &KubeVirt{
		conf:       conf,
		kubernetes: kubernetes,
	}, nil
}

func (m *KubeVirt) KVClientset(scope string) (*kubevirt.Clientset, error) {
	if v, ok := m.kvClientsets.Load(scope); ok {
		return v.(*kubevirt.Clientset), nil
	}

	config, err := m.kubernetes.Config(scope)
	if err != nil {
		return nil, err
	}

	clientset, err := kubevirt.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.kvClientsets.Store(scope, clientset)

	return clientset, nil
}

func (m *KubeVirt) CDIClientset(scope string) (*containerizeddataimporter.Clientset, error) {
	if v, ok := m.cdiClientsets.Load(scope); ok {
		return v.(*containerizeddataimporter.Clientset), nil
	}

	config, err := m.kubernetes.Config(scope)
	if err != nil {
		return nil, err
	}

	clientset, err := containerizeddataimporter.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.cdiClientsets.Store(scope, clientset)

	return clientset, nil
}
