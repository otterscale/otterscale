package kube

import (
	"sync"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	cdi "github.com/otterscale/kubevirt-client-go/containerizeddataimporter"
	virt "github.com/otterscale/kubevirt-client-go/kubevirt"

	"github.com/otterscale/otterscale/internal/config"
)

type Kube struct {
	conf           *config.Config
	clientsets     sync.Map
	virtClientsets sync.Map
	cdiClientsets  sync.Map

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
	return m.conf.Kube.HelmRepositoryURLs
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

func (m *Kube) virtClientset(config *rest.Config) (*virt.Clientset, error) {
	if v, ok := m.virtClientsets.Load(config.Host); ok {
		return v.(*virt.Clientset), nil
	}

	clientset, err := virt.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.virtClientsets.Store(config.Host, clientset)

	return clientset, nil
}

func (m *Kube) cdiClientset(config *rest.Config) (*cdi.Clientset, error) {
	if v, ok := m.cdiClientsets.Load(config.Host); ok {
		return v.(*cdi.Clientset), nil
	}

	clientset, err := cdi.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.cdiClientsets.Store(config.Host, clientset)

	return clientset, nil
}
