package core

import (
	"context"
	"fmt"
	"slices"

	"golang.org/x/mod/semver"
	"golang.org/x/sync/errgroup"
)

type Plugin struct {
	Release *Release
	Latest  *ChartVersion
}

func (uc *OrchestratorUseCase) ListGeneralPlugins(ctx context.Context, scope, facility string) ([]Plugin, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	repo := "https://prometheus-community.github.io/helm-charts"
	charts, err := uc.chart.List(ctx, repo)
	if err != nil {
		return nil, err
	}

	namespace := "monitoring"
	releases, err := uc.release.List(config, namespace)
	if err != nil {
		return nil, err
	}

	return uc.filterPlugins(charts, releases, []string{
		"kube-prometheus-stack",
	})
}

func (uc *OrchestratorUseCase) ListModelPlugins(ctx context.Context, scope, facility string) ([]Plugin, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	repos := []string{
		"https://llm-d-incubation.github.io/llm-d-infra",
		"https://nvidia.github.io/gpu-operator",
	}
	charts := make([][]Chart, len(repos))

	namespaces := []string{
		"llm-d",
		"nvidia-gpu-operator",
	}
	releases := make([][]Release, len(namespaces))

	eg, egctx := errgroup.WithContext(ctx)
	for i := range repos {
		eg.Go(func() error {
			v, err := uc.chart.List(egctx, repos[i])
			if err == nil {
				charts[i] = v
			}
			return err
		})
	}
	for i := range namespaces {
		eg.Go(func() error {
			v, err := uc.release.List(config, namespaces[i])
			if err == nil {
				releases[i] = v
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return uc.filterPlugins(flatten(charts), flatten(releases), []string{
		"gpu-operator",
		"llm-d-infra",
	})
}

func (uc *OrchestratorUseCase) ListInstancePlugins(ctx context.Context, scope, facility string) ([]Plugin, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	repo := "https://raw.githubusercontent.com/otterscale/otterscale-charts/refs/heads/main"
	charts, err := uc.chart.List(ctx, repo)
	if err != nil {
		return nil, err
	}

	releases, err := uc.release.List(config, "kubevirt")
	if err != nil {
		return nil, err
	}

	return uc.filterPlugins(charts, releases, []string{
		"kubevirt-infra",
	})
}

func (uc *OrchestratorUseCase) ListStoragePlugins(ctx context.Context, scope, facility string) ([]Plugin, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	repo := "https://raw.githubusercontent.com/otterscale/otterscale-charts/refs/heads/main"
	charts, err := uc.chart.List(ctx, repo)
	if err != nil {
		return nil, err
	}

	namespaces := []string{
		"samba-operator",
		"nfs-operator",
	}
	releases := make([][]Release, len(namespaces))

	eg, _ := errgroup.WithContext(ctx)
	for i := range namespaces {
		eg.Go(func() error {
			v, err := uc.release.List(config, namespaces[i])
			if err == nil {
				releases[i] = v
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return uc.filterPlugins(charts, flatten(releases), []string{
		"samba-operator",
		"nfs-operator",
	})
}

func (uc *OrchestratorUseCase) filterPlugins(charts []Chart, releases []Release, plugins []string) ([]Plugin, error) {
	result := []Plugin{}
	for _, plugin := range plugins {
		release, err := uc.findPluginRelease(releases, plugin)
		if err != nil {
			return nil, err
		}
		latest, err := uc.findLatestPluginChart(charts, plugin)
		if err != nil {
			return nil, err
		}
		result = append(result, Plugin{
			Release: release,
			Latest:  latest,
		})
	}
	return result, nil
}

func (uc *OrchestratorUseCase) findPluginRelease(releases []Release, plugin string) (*Release, error) {
	idx := slices.IndexFunc(releases, func(r Release) bool {
		return r.Chart != nil && r.Chart.Name() == plugin
	})
	if idx == -1 {
		return nil, nil
	}
	return &releases[idx], nil
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
