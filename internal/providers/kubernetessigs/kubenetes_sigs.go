package kubernetessigs

import (
	"sync"

	gaie "sigs.k8s.io/gateway-api-inference-extension/client-go/clientset/versioned"
	ga "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
)

type KubernetesSigs struct {
	conf       *config.Config
	kubernetes *kubernetes.Kubernetes

	gaClientsets   sync.Map
	gaieClientsets sync.Map
}

func New(conf *config.Config, kubernetes *kubernetes.Kubernetes) (*KubernetesSigs, error) {
	return &KubernetesSigs{
		conf:       conf,
		kubernetes: kubernetes,
	}, nil
}

func (m *KubernetesSigs) gaClientset(scope string) (*ga.Clientset, error) {
	if v, ok := m.gaClientsets.Load(scope); ok {
		return v.(*ga.Clientset), nil
	}

	config, err := m.kubernetes.Config(scope)
	if err != nil {
		return nil, err
	}

	clientset, err := ga.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.gaClientsets.Store(scope, clientset)

	return clientset, nil
}

func (m *KubernetesSigs) gaieClientset(scope string) (*gaie.Clientset, error) {
	if v, ok := m.gaieClientsets.Load(scope); ok {
		return v.(*gaie.Clientset), nil
	}

	config, err := m.kubernetes.Config(scope)
	if err != nil {
		return nil, err
	}

	clientset, err := gaie.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.gaieClientsets.Store(scope, clientset)

	return clientset, nil
}
