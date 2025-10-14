package core

import (
	"context"
	"slices"
	"strings"
)

func (uc *OrchestratorUseCase) ListPlugins(ctx context.Context, uuid, facility string) ([]Release, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	releases, err := uc.release.List(config, "")
	if err != nil {
		return nil, err
	}

	return uc.filterPluginReleases(releases), nil
}

func (uc *OrchestratorUseCase) filterPluginReleases(releases []Release) []Release {
	plugins := []string{"llm-d-infra", "kubevirt-infra"}

	return slices.DeleteFunc(releases, func(r Release) bool {
		if r.Chart == nil {
			return true
		}
		return !uc.isPluginChart(r.Chart.Name(), plugins)
	})
}

func (uc *OrchestratorUseCase) isPluginChart(chartName string, plugins []string) bool {
	for _, plugin := range plugins {
		if strings.Contains(chartName, plugin) {
			return true
		}
	}
	return false
}
