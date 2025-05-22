package kube

import (
	"sync"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/openhdc/otterscale/internal/config"
)

type Kube struct {
	conf       *config.Config
	clientsets *sync.Map

	envSettings    *cli.EnvSettings
	registryClient *registry.Client
}

func New(conf *config.Config) (*Kube, error) {
	opts := []registry.ClientOption{
		registry.ClientOptEnableCache(true),
	}
	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &Kube{
		conf:           conf,
		envSettings:    cli.New(),
		registryClient: registryClient,
	}, nil
}

func (m *Kube) helmRepoURLs() []string {
	kube := m.conf.GetKube()
	if kube != nil {
		return kube.GetHelmRepositoryUrls()
	}
	return nil
}

func (m *Kube) key(uuid, name string) string {
	return uuid + "/" + name
}

func (m *Kube) clientset(config *rest.Config) (*kubernetes.Clientset, error) {
	if v, ok := m.clientsets.Load(config.Host); ok {
		return v.(*kubernetes.Clientset), nil
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.clientsets.Store(config.Host, clientset)

	return clientset, nil
}
