package kube

import (
	"bytes"
	"fmt"
	"log/slog"
	"maps"

	"connectrpc.com/connect"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/kustomize/kyaml/kio"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"

	oscore "github.com/otterscale/otterscale/internal/core"
)

type postRenderer struct {
	extraLabels      map[string]string
	extraAnnotations map[string]string
}

func newPostRenderer(extraLabels, extraAnnotaions map[string]string) *postRenderer {
	return &postRenderer{
		extraLabels:      extraLabels,
		extraAnnotations: extraAnnotaions,
	}
}

func (p *postRenderer) Run(renderedManifests *bytes.Buffer) (*bytes.Buffer, error) {
	if len(p.extraLabels) == 0 {
		return renderedManifests, nil
	}
	nodes, err := kio.FromBytes(renderedManifests.Bytes())
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		// labels
		labels := node.GetLabels()
		maps.Copy(labels, p.extraLabels)
		if err := node.SetLabels(labels); err != nil {
			return nil, err
		}
		// annotations
		annotations := node.GetAnnotations()
		maps.Copy(annotations, p.extraAnnotations)
		if err := node.SetAnnotations(annotations); err != nil {
			return nil, err
		}
	}
	str, err := kio.StringAll(nodes)
	if err != nil {
		return nil, err
	}
	return bytes.NewBufferString(str), nil
}

type helmRelease struct {
	kube *Kube
}

func NewHelmRelease(kube *Kube) (oscore.ReleaseRepo, error) {
	return &helmRelease{
		kube: kube,
	}, nil
}

var _ oscore.ReleaseRepo = (*helmRelease)(nil)

func (r *helmRelease) List(restConfig *rest.Config, namespace string) ([]release.Release, error) {
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

func (r *helmRelease) Get(restConfig *rest.Config, namespace, name string) (*release.Release, error) {
	config, err := r.config(restConfig, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewGet(config)
	return client.Run(name)
}

func (r *helmRelease) Install(restConfig *rest.Config, namespace, name string, dryRun bool, chartRef string, labels, annotations map[string]string, values map[string]any) (*release.Release, error) {
	if !action.ValidName.MatchString(name) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid release name %q", name))
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
	client.PostRenderer = newPostRenderer(labels, annotations)

	chartPath, err := client.LocateChart(chartRef, r.kube.envSettings)
	if err != nil {
		return nil, err
	}

	chart, err := r.chartInstall(chartPath, client.DependencyUpdate, client.Keyring)
	if err != nil {
		return nil, err
	}
	return client.Run(chart, values)
}

func (r *helmRelease) Uninstall(restConfig *rest.Config, namespace, name string, dryRun bool) (*release.Release, error) {
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

func (r *helmRelease) Upgrade(restConfig *rest.Config, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error) {
	config, err := r.config(restConfig, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewUpgrade(config)
	client.Namespace = namespace
	client.DryRun = dryRun

	chartPath, err := client.LocateChart(chartRef, r.kube.envSettings)
	if err != nil {
		return nil, err
	}

	chart, err := r.chartInstall(chartPath, client.DependencyUpdate, client.Keyring)
	if err != nil {
		return nil, err
	}
	return client.Run(name, chart, values)
}

func (r *helmRelease) Rollback(restConfig *rest.Config, namespace, name string, dryRun bool) error {
	config, err := r.config(restConfig, namespace)
	if err != nil {
		return err
	}

	client := action.NewRollback(config)
	client.DryRun = dryRun
	return client.Run(name)
}

func (r *helmRelease) GetValues(restConfig *rest.Config, namespace, name string) (map[string]any, error) {
	config, err := r.config(restConfig, namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewGetValues(config)
	client.AllValues = true
	return client.Run(name)
}

func (r *helmRelease) chartInstall(chartPath string, dependencyUpdate bool, keyring string) (*chart.Chart, error) {
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

func (r *helmRelease) config(restConfig *rest.Config, namespace string) (*action.Configuration, error) {
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
	config.RegistryClient = r.kube.registryClient
	return config, nil
}
