package samba

import (
	"sync"

	"github.com/otterscale/samba-operator-client-go/samba"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
)

type Samba struct {
	conf       *config.Config
	kubernetes *kubernetes.Kubernetes

	clientsets sync.Map
}

func New(conf *config.Config, kubernetes *kubernetes.Kubernetes) (*Samba, error) {
	return &Samba{
		conf:       conf,
		kubernetes: kubernetes,
	}, nil
}

func (m *Samba) clientset(scope string) (*samba.Clientset, error) {
	if v, ok := m.clientsets.Load(scope); ok {
		return v.(*samba.Clientset), nil
	}

	config, err := m.kubernetes.Config(scope)
	if err != nil {
		return nil, err
	}

	clientset, err := samba.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.clientsets.Store(scope, clientset)

	return clientset, nil
}
