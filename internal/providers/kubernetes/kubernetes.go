package kubernetes

import (
	"context"
	"encoding/base64"
	"os"
	"sync"
	"time"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core/scope"
	"github.com/otterscale/otterscale/internal/providers/juju"
)

type Kubernetes struct {
	conf           *config.Config
	juju           *juju.Juju
	envSettings    *cli.EnvSettings
	registryClient *registry.Client

	configs    sync.Map
	clientsets sync.Map
}

func New(conf *config.Config, juju *juju.Juju) (*Kubernetes, error) {
	opts := []registry.ClientOption{
		registry.ClientOptEnableCache(true),
	}
	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &Kubernetes{
		conf:           conf,
		juju:           juju,
		envSettings:    cli.New(),
		registryClient: registryClient,
	}, nil
}

func (m *Kubernetes) newMicroK8sConfig() (*rest.Config, error) {
	kubeConfig, err := base64.StdEncoding.DecodeString(m.conf.MicroK8s.Config)
	if err != nil {
		return nil, err
	}

	configAPI, err := clientcmd.Load(kubeConfig)
	if err != nil {
		return nil, err
	}

	return clientcmd.NewDefaultClientConfig(*configAPI, &clientcmd.ConfigOverrides{}).ClientConfig()
}

func (m *Kubernetes) getKubeConfig(ctx context.Context, scope, name string) (*api.Config, error) {
	result, err := m.juju.Run(ctx, scope, name, "get-kubeconfig", nil)
	if err != nil {
		return nil, err
	}
	return clientcmd.Load([]byte(result["kubeconfig"].(string)))
}

func (m *Kubernetes) writeCAToFile(caData []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", "otterscale-ca-*.crt")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(caData); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func (m *Kubernetes) newConfig(scope string) (*rest.Config, error) {
	name := scope + "-kubernetes-control-plane"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	kubeConfig, err := m.getKubeConfig(ctx, scope, name)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.NewDefaultClientConfig(*kubeConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}

	// Write CA data to temp file for helm
	if config.CAData != nil {
		fileName, err := m.writeCAToFile(config.CAData)
		if err != nil {
			return nil, err
		}
		config.CAFile = fileName
	}

	return config, nil
}

func (m *Kubernetes) getConfig(scopeName string) (*rest.Config, error) {
	if scopeName == scope.ReservedName {
		return m.newMicroK8sConfig()
	}
	return m.newConfig(scopeName)
}

func (m *Kubernetes) Config(scope string) (*rest.Config, error) {
	if v, ok := m.configs.Load(scope); ok {
		return v.(*rest.Config), nil
	}

	config, err := m.getConfig(scope)
	if err != nil {
		return nil, err
	}

	m.configs.Store(scope, config)

	return config, nil
}

func (m *Kubernetes) clientset(scope string) (*kubernetes.Clientset, error) {
	if v, ok := m.clientsets.Load(scope); ok {
		return v.(*kubernetes.Clientset), nil
	}

	config, err := m.Config(scope)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.clientsets.Store(scope, clientset)

	return clientset, nil
}
