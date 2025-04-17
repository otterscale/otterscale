package kube

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	"sigs.k8s.io/yaml"

	"github.com/openhdc/openhdc/internal/domain/service"
	"github.com/openhdc/openhdc/internal/env"
)

const (
	defaultRepositoryURL = "http://chartmuseum:8080"
)

type helm struct {
	helmMap            HelmMap
	settings           *cli.EnvSettings
	providers          getter.Providers
	registryClient     *registry.Client
	repoIndex          *repo.IndexFile
	repoIndexCacheTime time.Time
}

func NewHelm(helmMap HelmMap) (service.KubeHelm, error) {
	settings := cli.New()

	opts := []registry.ClientOption{
		registry.ClientOptEnableCache(true),
	}
	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &helm{
		helmMap:        helmMap,
		settings:       settings,
		providers:      getter.All(settings),
		registryClient: registryClient,
	}, nil
}

var _ service.KubeHelm = (*helm)(nil)

func (r *helm) ListReleases(uuid, facility, namespace string) ([]release.Release, error) {
	config, err := r.helmMap.get(uuid, facility, namespace, r.registryClient)
	if err != nil {
		return nil, err
	}

	client := action.NewList(config)
	client.Deployed = true
	rels, err := client.Run()
	if err != nil {
		return nil, err
	}

	rs := []release.Release{}
	for _, rel := range rels {
		rs = append(rs, *rel)
	}
	return rs, nil
}

func (r *helm) InstallRelease(uuid, facility, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error) {
	if !action.ValidName.MatchString(name) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid release name %q", name)
	}

	config, err := r.helmMap.get(uuid, facility, namespace, r.registryClient)
	if err != nil {
		return nil, err
	}

	client := action.NewInstall(config)
	client.CreateNamespace = true
	client.Namespace = namespace
	client.DryRun = dryRun
	client.ReleaseName = name

	chartPath, err := client.ChartPathOptions.LocateChart(chartRef, r.settings)
	if err != nil {
		return nil, err
	}

	chart, err := r.chartInstall(chartPath, client.DependencyUpdate, client.ChartPathOptions.Keyring, config.RegistryClient)
	if err != nil {
		return nil, err
	}
	return client.Run(chart, values)
}

func (r *helm) UninstallRelease(uuid, facility, namespace, name string, dryRun bool) (*release.Release, error) {
	config, err := r.helmMap.get(uuid, facility, namespace, r.registryClient)
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

func (r *helm) UpgradeRelease(uuid, facility, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error) {
	config, err := r.helmMap.get(uuid, facility, namespace, r.registryClient)
	if err != nil {
		return nil, err
	}

	client := action.NewUpgrade(config)
	client.Namespace = namespace
	client.DryRun = dryRun

	chartPath, err := client.ChartPathOptions.LocateChart(chartRef, r.settings)
	if err != nil {
		return nil, err
	}

	chart, err := r.chartInstall(chartPath, client.DependencyUpdate, client.ChartPathOptions.Keyring, config.RegistryClient)
	if err != nil {
		return nil, err
	}
	return client.Run(name, chart, values)
}

func (r *helm) RollbackRelease(uuid, facility, namespace, name string, dryRun bool) error {
	config, err := r.helmMap.get(uuid, facility, namespace, r.registryClient)
	if err != nil {
		return err
	}

	client := action.NewRollback(config)
	client.DryRun = dryRun
	return client.Run(name)
}

func (r *helm) ShowChart(chartRef string, format action.ShowOutputFormat) (string, error) {
	client := action.NewShow(format)
	client.SetRegistryClient(r.registryClient)

	chartPath, err := client.ChartPathOptions.LocateChart(chartRef, r.settings)
	if err != nil {
		return "", err
	}
	return client.Run(chartPath)
}

func (r *helm) ListChartVersions(ctx context.Context) (map[string]repo.ChartVersions, error) {
	if err := r.fetchRepositoryIndex(ctx); err != nil {
		return nil, err
	}
	r.repoIndex.SortEntries()
	return r.repoIndex.Entries, nil
}

// FIXME: WORKAROUND
func (r *helm) fetchRepositoryIndex(ctx context.Context) error {
	if r.repoIndex != nil && time.Since(r.repoIndexCacheTime) < time.Hour*24 {
		return nil
	}

	queryURL, err := url.ParseRequestURI(env.GetOrDefault(env.OPENHDC_HELM_REPOSITORY_URL, defaultRepositoryURL))
	if err != nil {
		return err
	}
	queryURL = queryURL.JoinPath("index.yaml")

	url := queryURL.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return fmt.Errorf("fetch repository index from %q failed: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch repository index from %q failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("fetch repository index from %q failed: %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("fetch repository index from %q failed: %w", url, err)
	}

	r.repoIndex = new(repo.IndexFile)
	if err := yaml.Unmarshal(data, r.repoIndex); err != nil {
		return err
	}
	r.repoIndexCacheTime = time.Now()

	return nil
}

func (r *helm) chartInstall(chartPath string, dependencyUpdate bool, keyring string, rc *registry.Client) (*chart.Chart, error) {
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
				Getters:          r.providers,
				RepositoryConfig: r.settings.RepositoryConfig,
				RepositoryCache:  r.settings.RepositoryCache,
				Debug:            r.settings.Debug,
				RegistryClient:   rc,
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
