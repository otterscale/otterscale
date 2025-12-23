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
}

func New(conf *config.Config, kubernetes *kubernetes.Kubernetes) (*Helm, error) {
	return &Helm{
		conf:       conf,
		kubernetes: kubernetes,
	}, nil
}

func newEnvSettings() *cli.EnvSettings {
	envSettings := cli.New()
	envSettings.RepositoryConfig = filepath.Join(os.TempDir(), "helm", "repositories.yaml")
	envSettings.RepositoryCache = filepath.Join(os.TempDir(), "helm", "repository")
	return envSettings
}

func newRegistryClient() (*registry.Client, error) {
	return registry.NewClient(registry.ClientOptEnableCache(true))
}
