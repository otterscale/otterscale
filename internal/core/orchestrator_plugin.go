package core

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"golang.org/x/mod/semver"
	"golang.org/x/sync/errgroup"
)

type plugin struct {
	name        string
	namespace   string
	repoURL     string
	labels      map[string]string
	annotations map[string]string
	values      map[string]any
}

var pluginMap = map[string]plugin{
	"kube-prometheus-stack": {
		name:      "kube-prometheus-stack",
		namespace: "monitoring",
		repoURL:   "https://prometheus-community.github.io/helm-charts",
	},
	"gpu-operator": {
		name:      "gpu-operator",
		namespace: "nvidia-gpu-operator",
		repoURL:   "https://nvidia.github.io/gpu-operator",
	},
	"llm-d-infra": {
		name:      "llm-d-infra",
		namespace: "llm-d",
		repoURL:   "https://llm-d-incubation.github.io/llm-d-infra",
	},
	"kubevirt-infra": {
		name:      "kubevirt-infra",
		namespace: "kubevirt",
		repoURL:   "https://raw.githubusercontent.com/otterscale/otterscale-charts/refs/heads/main",
	},
	"samba-operator": {
		name:      "samba-operator",
		namespace: "samba-operator",
		repoURL:   "https://raw.githubusercontent.com/otterscale/otterscale-charts/refs/heads/main",
	},
	"nfs-operator": {
		name:      "nfs-operator",
		namespace: "nfs-operator",
		repoURL:   "https://raw.githubusercontent.com/otterscale/otterscale-charts/refs/heads/main",
	},
}

type Plugin struct {
	Release *Release
	Latest  *ChartVersion
}

func (uc *OrchestratorUseCase) ListGeneralPlugins(ctx context.Context, scope, facility string) ([]Plugin, error) {
	return uc.listPlugins(ctx, scope, facility, []string{
		"kube-prometheus-stack",
	})
}

func (uc *OrchestratorUseCase) ListModelPlugins(ctx context.Context, scope, facility string) ([]Plugin, error) {
	return uc.listPlugins(ctx, scope, facility, []string{
		"gpu-operator",
		"llm-d-infra",
	})
}

func (uc *OrchestratorUseCase) ListInstancePlugins(ctx context.Context, scope, facility string) ([]Plugin, error) {
	return uc.listPlugins(ctx, scope, facility, []string{
		"kubevirt-infra",
	})
}

func (uc *OrchestratorUseCase) ListStoragePlugins(ctx context.Context, scope, facility string) ([]Plugin, error) {
	return uc.listPlugins(ctx, scope, facility, []string{
		"samba-operator",
		"nfs-operator",
	})
}

func (uc *OrchestratorUseCase) InstallPlugins(ctx context.Context, scope, facility string, chartRefMap map[string]string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return err
	}

	names := slices.Collect(maps.Keys(chartRefMap))
	plugins := uc.filterInternalPlugins(names)

	eg, _ := errgroup.WithContext(ctx)
	for _, p := range plugins {
		eg.Go(func() error {
			ref := chartRefMap[p.name]
			_, err := uc.release.Install(config, p.namespace, p.name, false, ref, p.labels, p.annotations, p.values)
			return err
		})
	}
	return eg.Wait()
}

func (uc *OrchestratorUseCase) UpgradePlugins(ctx context.Context, scope, facility string, chartRefMap map[string]string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return err
	}

	names := slices.Collect(maps.Keys(chartRefMap))
	plugins := uc.filterInternalPlugins(names)

	eg, _ := errgroup.WithContext(ctx)
	for _, p := range plugins {
		eg.Go(func() error {
			ref := chartRefMap[p.name]
			_, err := uc.release.Upgrade(config, p.namespace, p.name, false, ref, p.values)
			return err
		})
	}
	return eg.Wait()
}

func (uc *OrchestratorUseCase) listPlugins(ctx context.Context, scope, facility string, names []string) ([]Plugin, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	plugins := uc.filterInternalPlugins(names)
	charts := make([][]Chart, len(plugins))
	releases := make([][]Release, len(plugins))

	eg, egctx := errgroup.WithContext(ctx)
	for i := range plugins {
		eg.Go(func() error {
			v, err := uc.chart.List(egctx, plugins[i].repoURL, true)
			if err == nil {
				charts[i] = v
			}
			return err
		})
		eg.Go(func() error {
			v, err := uc.release.List(config, plugins[i].namespace)
			if err == nil {
				releases[i] = v
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return uc.filterPlugins(flatten(charts), flatten(releases), names)
}

func (uc *OrchestratorUseCase) filterInternalPlugins(plugins []string) []plugin {
	result := []plugin{}
	for _, p := range plugins {
		if plg, ok := pluginMap[p]; ok {
			result = append(result, plg)
		}
	}
	return result
}

func (uc *OrchestratorUseCase) filterPlugins(charts []Chart, releases []Release, plugins []string) ([]Plugin, error) {
	result := []Plugin{}
	for _, plugin := range plugins {
		latest, err := uc.findLatestPluginChart(charts, plugin)
		if err != nil {
			return nil, err
		}
		result = append(result, Plugin{
			Release: uc.findPluginRelease(releases, plugin),
			Latest:  latest,
		})
	}
	return result, nil
}

func (uc *OrchestratorUseCase) findPluginRelease(releases []Release, plugin string) *Release {
	idx := slices.IndexFunc(releases, func(r Release) bool {
		return r.Chart != nil && r.Chart.Name() == plugin
	})
	if idx == -1 {
		return nil
	}
	return &releases[idx]
}

func (uc *OrchestratorUseCase) findLatestPluginChart(charts []Chart, plugin string) (*ChartVersion, error) {
	idx := slices.IndexFunc(charts, func(c Chart) bool {
		return c.Name == plugin
	})
	if idx == -1 {
		return nil, fmt.Errorf("chart not found for plugin: %s", plugin)
	}

	chart := charts[idx]
	if len(chart.Versions) == 0 {
		return nil, fmt.Errorf("no versions available for plugin: %s", plugin)
	}

	slices.SortFunc(chart.Versions, func(a, b *ChartVersion) int {
		return semver.Compare(b.Version, a.Version)
	})

	return chart.Versions[0], nil
}
