package kube

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
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

	"github.com/openhdc/otterscale/internal/domain/service"
	"github.com/openhdc/otterscale/internal/env"
	"github.com/openhdc/otterscale/internal/utils"
)

const defaultRepositoryURL = "http://chartmuseum:8080"

type helm struct {
	helmMap            HelmMap
	settings           *cli.EnvSettings
	providers          getter.Providers
	registryClient     *registry.Client
	repoURLs           []string
	repoIndexFiles     []*repo.IndexFile
	repoIndexCacheTime time.Time
}

func NewHelm(helmMap HelmMap) (service.KubeHelm, error) {
	settings := cli.New()

	opts := []registry.ClientOption{
		registry.ClientOptEnableCache(true),
	}
	username := os.Getenv(env.OPENHDC_REGISTRY_USERNAME)
	password := os.Getenv(env.OPENHDC_REGISTRY_PASSWORD)
	if username != "" && password != "" {
		opts = append(opts, registry.ClientOptBasicAuth(username, password))
	}

	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	repoURLs := []string{}
	repoURLs = append(repoURLs, strings.Split(env.GetOrDefault(env.OPENHDC_HELM_REPOSITORY_URLS, defaultRepositoryURL), ",")...)

	return &helm{
		helmMap:        helmMap,
		settings:       settings,
		providers:      getter.All(settings),
		registryClient: registryClient,
		repoURLs:       repoURLs,
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

func (r *helm) GetValues(uuid, facility, namespace, name string) (map[string]any, error) {
	config, err := r.helmMap.get(uuid, facility, namespace, r.registryClient)
	if err != nil {
		return nil, err
	}

	client := action.NewGetValues(config)
	client.AllValues = true
	return client.Run(name)
}

func (r *helm) ListChartVersions(ctx context.Context) ([]*repo.IndexFile, error) {
	if r.repoIndexFiles != nil && time.Since(r.repoIndexCacheTime) < time.Hour*24 {
		return r.repoIndexFiles, nil
	}
	eg, ctx := errgroup.WithContext(ctx)
	result := make([]*repo.IndexFile, len(r.repoURLs))
	for i := range r.repoURLs {
		url := r.repoURLs[i]
		eg.Go(func() error {
			f, err := r.fetchRepositoryIndex(ctx, url)
			if err == nil {
				f.SortEntries()
				result[i] = f
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	r.repoIndexFiles = result
	r.repoIndexCacheTime = time.Now()
	return r.repoIndexFiles, nil
}

func (r *helm) fetchRepositoryIndex(ctx context.Context, repoURL string) (*repo.IndexFile, error) {
	var data []byte

	if strings.HasPrefix(repoURL, "http") {
		queryURL, err := url.ParseRequestURI(repoURL)
		if err != nil {
			return nil, err
		}
		queryURL = queryURL.JoinPath("index.yaml")

		data, err = utils.Get(ctx, queryURL.String())
		if err != nil {
			return nil, err
		}
	} else {
		path, err := filepath.Abs(repoURL)
		if err != nil {
			return nil, err
		}

		path = filepath.Join(path, "index.yaml")

		data, err = os.ReadFile(path)
		if err != nil {
			return nil, err
		}
	}

	f := new(repo.IndexFile)
	if err := yaml.Unmarshal(data, f); err != nil {
		return nil, err
	}
	return f, nil
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
