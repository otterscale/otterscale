package kubernetes

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/go-faker/faker/v4"
	"github.com/goccy/go-yaml"
	"helm.sh/helm/v3/pkg/action"
	helmchart "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/otterscale/otterscale/internal/core/application/release"
)

// Note: Helm API do not support context.
type releaseRepo struct {
	kubernetes *Kubernetes
}

func NewReleaseRepo(kubernetes *Kubernetes) (release.ReleaseRepo, error) {
	return &releaseRepo{
		kubernetes: kubernetes,
	}, nil
}

func (r *releaseRepo) List(_ context.Context, scope, namespace, selector string) ([]release.Release, error) {
	config, err := r.config(scope, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewList(config)
	client.Selector = selector
	client.Deployed = true

	releases, err := client.Run()
	if err != nil {
		return nil, err
	}

	return pointersToValues(releases), nil
}

func (r *releaseRepo) Get(_ context.Context, scope, namespace, name string) (*release.Release, error) {
	config, err := r.config(scope, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewGet(config)

	return client.Run(name)
}

func (r *releaseRepo) Install(_ context.Context, scope, namespace, name string, dryRun bool, chartRef string, labelsInSecrets, labels, annotations map[string]string, valuesYAML string, valuesMap map[string]string) (*release.Release, error) {
	name = r.newName(name)

	if !action.ValidName.MatchString(name) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid release name %q", name))
	}

	config, err := r.config(scope, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewInstall(config)
	client.CreateNamespace = true
	client.Namespace = namespace
	client.DryRun = dryRun
	client.ReleaseName = name
	client.Labels = labelsInSecrets
	client.PostRenderer = newPostRenderer(labels, annotations)

	chartPath, err := client.LocateChart(chartRef, r.kubernetes.envSettings)
	if err != nil {
		return nil, err
	}

	chart, err := r.chartInstall(chartPath, client.DependencyUpdate, client.Keyring)
	if err != nil {
		return nil, err
	}

	valuesMap[release.ChartRefKey] = chartRef

	values, err := r.toValues(valuesYAML, valuesMap)
	if err != nil {
		return nil, err
	}

	return client.Run(chart, values)
}

func (r *releaseRepo) Uninstall(_ context.Context, scope, namespace, name string, dryRun bool) (*release.Release, error) {
	config, err := r.config(scope, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewUninstall(config)
	client.DeletionPropagation = "background"
	client.DryRun = dryRun

	resp, err := client.Run(name)
	if err != nil {
		return nil, err
	}

	return resp.Release, nil
}

func (r *releaseRepo) Upgrade(_ context.Context, scope, namespace, name string, dryRun bool, chartRef, valuesYAML string, valuesMap map[string]string, reuseValues bool) (*release.Release, error) {
	config, err := r.config(scope, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewUpgrade(config)
	client.Namespace = namespace
	client.DryRun = dryRun
	client.ReuseValues = reuseValues

	chartPath, err := client.LocateChart(chartRef, r.kubernetes.envSettings)
	if err != nil {
		return nil, err
	}

	chart, err := r.chartInstall(chartPath, client.DependencyUpdate, client.Keyring)
	if err != nil {
		return nil, err
	}

	valuesMap[release.ChartRefKey] = chartRef

	values, err := r.toValues(valuesYAML, valuesMap)
	if err != nil {
		return nil, err
	}

	return client.Run(name, chart, values)
}

func (r *releaseRepo) Rollback(_ context.Context, scope, namespace, name string, dryRun bool) error {
	config, err := r.config(scope, namespace)
	if err != nil {
		return err
	}

	client := action.NewRollback(config)
	client.DryRun = dryRun

	return client.Run(name)
}

func (r *releaseRepo) GetValues(_ context.Context, scope, namespace, name string) (map[string]any, error) {
	config, err := r.config(scope, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewGetValues(config)
	client.AllValues = true

	return client.Run(name)
}

func (r *releaseRepo) chartInstall(chartPath string, dependencyUpdate bool, keyring string) (*helmchart.Chart, error) {
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
				Getters:          getter.All(r.kubernetes.envSettings),
				RepositoryConfig: r.kubernetes.envSettings.RepositoryConfig,
				RepositoryCache:  r.kubernetes.envSettings.RepositoryCache,
				Debug:            r.kubernetes.envSettings.Debug,
				RegistryClient:   r.kubernetes.registryClient,
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

func (r *releaseRepo) config(scope, namespace string) (*action.Configuration, error) {
	restConfig, err := r.kubernetes.Config(scope)
	if err != nil {
		return nil, err
	}

	getter := genericclioptions.NewConfigFlags(true)
	getter.APIServer = &restConfig.Host
	getter.BearerToken = &restConfig.BearerToken
	getter.Insecure = &restConfig.Insecure
	getter.CAFile = &restConfig.CAFile
	getter.Namespace = &namespace

	config := new(action.Configuration)
	if err := config.Init(getter, namespace, "", slog.Debug); err != nil {
		return nil, err
	}

	config.RegistryClient = r.kubernetes.registryClient

	return config, nil
}

func (uc *releaseRepo) newName(name string) string {
	if name != "" {
		return name
	}
	return strings.ToLower(faker.FirstName() + "-" + faker.Username())
}

func (r *releaseRepo) toValues(valuesYAML string, valuesMap map[string]string) (map[string]any, error) {
	// advanced
	values := map[string]any{}
	if err := yaml.Unmarshal([]byte(valuesYAML), &values); err != nil {
		return nil, err
	}

	// basic
	vals := []string{}
	for k, v := range valuesMap {
		if v != "" {
			vals = append(vals, fmt.Sprintf("%s=%s", k, v))
		}
	}

	if err := strvals.ParseInto(strings.Join(vals, ","), values); err != nil {
		return nil, err
	}

	return values, nil
}

func pointersToValues[T any](ptrs []*T) []T {
	ret := make([]T, 0, len(ptrs))

	for _, ptr := range ptrs {
		ret = append(ret, *ptr)
	}

	return ret
}
