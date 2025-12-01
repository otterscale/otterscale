package helm

import (
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
)

type Helm struct {
	conf       *config.Config
	kubernetes *kubernetes.Kubernetes

	envSettings    *cli.EnvSettings
	registryClient *registry.Client
}

func New(conf *config.Config, kubernetes *kubernetes.Kubernetes) (*Helm, error) {
	opts := []registry.ClientOption{
		registry.ClientOptEnableCache(true),
	}

	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	envSettings := cli.New()
	envSettings.RepositoryConfig = filepath.Join(os.TempDir(), "helm", "repositories.yaml")
	envSettings.RepositoryCache = filepath.Join(os.TempDir(), "helm", "repository")

	return &Helm{
		conf:           conf,
		kubernetes:     kubernetes,
		envSettings:    envSettings,
		registryClient: registryClient,
	}, nil
}
