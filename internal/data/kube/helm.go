package kube

import (
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type helm struct {
	kube *Kube
}

func NewHelm(kube *Kube) (service.KubeHelm, error) {
	return &helm{
		kube: kube,
	}, nil
}

var _ service.KubeHelm = (*helm)(nil)

func (r *helm) ListReleases(restConfig *rest.Config, namespace string) ([]release.Release, error) {
	config, err := r.config(restConfig, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewList(config)
	client.Deployed = true

	releases, err := client.Run()
	if err != nil {
		return nil, err
	}

	result := []release.Release{}
	for _, release := range releases {
		result = append(result, *release)
	}
	return result, nil
}

func (r *helm) InstallRelease(restConfig *rest.Config, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error) {
	if !action.ValidName.MatchString(name) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid release name %q", name)
	}

	config, err := r.config(restConfig, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewInstall(config)
	client.CreateNamespace = true
	client.Namespace = namespace
	client.DryRun = dryRun
	client.ReleaseName = name

	chartPath, err := client.ChartPathOptions.LocateChart(chartRef, r.kube.envSettings)
	if err != nil {
		return nil, err
	}

	chart, err := r.chartInstall(chartPath, client.DependencyUpdate, client.ChartPathOptions.Keyring)
	if err != nil {
		return nil, err
	}
	return client.Run(chart, values)
}

func (r *helm) UninstallRelease(restConfig *rest.Config, namespace, name string, dryRun bool) (*release.Release, error) {
	config, err := r.config(restConfig, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewUninstall(config)
	client.DeletionPropagation = "background"
	client.DryRun = dryRun

	res, err := client.Run(name)
	if err != nil {
		return nil, err
	}
	return res.Release, nil
}

func (r *helm) UpgradeRelease(restConfig *rest.Config, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error) {
	config, err := r.config(restConfig, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewUpgrade(config)
	client.Namespace = namespace
	client.DryRun = dryRun

	chartPath, err := client.ChartPathOptions.LocateChart(chartRef, r.kube.envSettings)
	if err != nil {
		return nil, err
	}

	chart, err := r.chartInstall(chartPath, client.DependencyUpdate, client.ChartPathOptions.Keyring)
	if err != nil {
		return nil, err
	}
	return client.Run(name, chart, values)
}

func (r *helm) RollbackRelease(restConfig *rest.Config, namespace, name string, dryRun bool) error {
	config, err := r.config(restConfig, namespace)
	if err != nil {
		return err
	}

	client := action.NewRollback(config)
	client.DryRun = dryRun
	return client.Run(name)
}

func (r *helm) GetValues(restConfig *rest.Config, namespace, name string) (map[string]any, error) {
	config, err := r.config(restConfig, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewGetValues(config)
	client.AllValues = true
	return client.Run(name)
}

func (r *helm) chartInstall(chartPath string, dependencyUpdate bool, keyring string) (*chart.Chart, error) {
	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	if chartDependencies := chart.Metadata.Dependencies; chartDependencies != nil {
		if err := action.CheckDependencies(chart, chartDependencies); err != nil {
			if !dependencyUpdate {
				return nil, fmt.Errorf("failed to check chart dependencies: %w", err)
			}

			manager := &downloader.Manager{
				ChartPath:        chartPath,
				Keyring:          keyring,
				SkipUpdate:       false,
				Getters:          getter.All(r.kube.envSettings),
				RepositoryConfig: r.kube.envSettings.RepositoryConfig,
				RepositoryCache:  r.kube.envSettings.RepositoryCache,
				Debug:            r.kube.envSettings.Debug,
				RegistryClient:   r.kube.registryClient,
			}
			if err := manager.Update(); err != nil {
				return nil, err
			}

			// Reload the chart with the updated Chart.lock file.
			if chart, err = loader.Load(chartPath); err != nil {
				return nil, fmt.Errorf("failed reloading chart after repo update: %w", err)
			}
		}
	}

	return chart, nil
}

func (r *helm) config(restConfig *rest.Config, namespace string) (*action.Configuration, error) {
	getter := genericclioptions.NewConfigFlags(true)
	getter.APIServer = &restConfig.Host
	getter.BearerToken = &restConfig.BearerToken
	getter.CAFile = &restConfig.CAFile
	getter.CertFile = &restConfig.TLSClientConfig.CertFile
	getter.KeyFile = &restConfig.TLSClientConfig.KeyFile
	getter.Insecure = &restConfig.Insecure
	getter.Namespace = &namespace

	config := new(action.Configuration)
	if err := config.Init(getter, namespace, "", log.Printf); err != nil {
		return nil, err
	}
	config.RegistryClient = r.kube.registryClient
	return config, nil
}
