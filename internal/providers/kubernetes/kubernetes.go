package kubernetes

import (
	"context"
	"os"
	"sync"

	"github.com/otterscale/kubevirt-client-go/containerizeddataimporter"
	"github.com/otterscale/kubevirt-client-go/kubevirt"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/providers/juju"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Kube struct {
	conf           *config.Config
	juju           *juju.Juju
	envSettings    *cli.EnvSettings
	registryClient *registry.Client

	configs       sync.Map
	clientsets    sync.Map
	kvClientsets  sync.Map
	cdiClientsets sync.Map
}

func New(conf *config.Config, juju *juju.Juju) (*Kube, error) {
	opts := []registry.ClientOption{
		registry.ClientOptEnableCache(true),
	}
	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &Kube{
		conf:           conf,
		juju:           juju,
		envSettings:    cli.New(),
		registryClient: registryClient,
	}, nil
}

func (m *Kube) getKubeConfig(ctx context.Context, scope, name string) (*api.Config, error) {
	result, err := m.juju.RunAction(ctx, scope, name, "get-kubeconfig", nil)
	if err != nil {
		return nil, err
	}
	return clientcmd.Load([]byte(result.Output["kubeconfig"].(string)))
}

func (m *Kube) writeCAToFile(caData []byte) (string, error) {
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

func (m *Kube) newConfig(scope string) (*rest.Config, error) {
	name := scope + "-kubernetes-control-plane"

	kubeConfig, err := m.getKubeConfig(context.Background(), scope, name)
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

func (m *Kube) config(scope string) (*rest.Config, error) {
	if v, ok := m.configs.Load(scope); ok {
		return v.(*rest.Config), nil
	}

	config, err := m.newConfig(scope)
	if err != nil {
		return nil, err
	}

	m.configs.Store(scope, config)

	return config, nil
}

func (m *Kube) Clientset(scope string) (*kubernetes.Clientset, error) {
	if v, ok := m.clientsets.Load(scope); ok {
		return v.(*kubernetes.Clientset), nil
	}

	config, err := m.config(scope)
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

func (m *Kube) KVClientset(scope string) (*kubevirt.Clientset, error) {
	if v, ok := m.kvClientsets.Load(scope); ok {
		return v.(*kubevirt.Clientset), nil
	}

	config, err := m.config(scope)
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

func (m *Kube) CDIClientset(scope string) (*containerizeddataimporter.Clientset, error) {
	if v, ok := m.cdiClientsets.Load(scope); ok {
		return v.(*containerizeddataimporter.Clientset), nil
	}

	config, err := m.config(scope)
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
