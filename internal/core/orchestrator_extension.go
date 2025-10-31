package core

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"golang.org/x/mod/semver"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
)

const chartRepoURL = "https://otterscale.github.io/charts"

type extension struct {
	name        string
	namespace   string
	repoURL     string
	labels      map[string]string
	annotations map[string]string
	valuesMap   map[string]string
}

var extensionMap = map[string]extension{
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
		valuesMap: map[string]string{
			"nameOverride": "llm-d-infra",
			"gateway.gatewayParameters.resources.limits.cpu":    "4",
			"gateway.gatewayParameters.resources.limits.memory": "2Gi",
			"gateway.service.type":                              string(corev1.ServiceTypeNodePort),
		},
	},
	"kubevirt-infra": {
		name:      "kubevirt-infra",
		namespace: "kubevirt",
		repoURL:   chartRepoURL,
	},
	"samba-operator": {
		name:      "samba-operator",
		namespace: "samba-operator",
		repoURL:   chartRepoURL,
	},
	"nfs-operator": {
		name:      "nfs-operator",
		namespace: "nfs-operator",
		repoURL:   chartRepoURL,
	},
}

type Extension struct {
	Release   *Release
	Latest    *ChartVersion
	LatestURL string
}

func (uc *OrchestratorUseCase) ListGeneralExtensions(ctx context.Context, scope, facility string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, facility, []string{
		"kube-prometheus-stack",
	})
}

func (uc *OrchestratorUseCase) ListModelExtensions(ctx context.Context, scope, facility string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, facility, []string{
		"gpu-operator",
		"llm-d-infra",
	})
}

func (uc *OrchestratorUseCase) ListInstanceExtensions(ctx context.Context, scope, facility string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, facility, []string{
		"kubevirt-infra",
	})
}

func (uc *OrchestratorUseCase) ListStorageExtensions(ctx context.Context, scope, facility string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, facility, []string{
		"samba-operator",
	})
}

func (uc *OrchestratorUseCase) InstallExtensions(ctx context.Context, scope, facility string, chartRefMap map[string]string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return err
	}

	labels := map[string]string{
		TypeLabel: "extension",
	}
	names := slices.Collect(maps.Keys(chartRefMap))
	extensions := uc.filterInternalExtensions(names)

	eg, _ := errgroup.WithContext(ctx)
	for _, p := range extensions {
		eg.Go(func() error {
			ref := chartRefMap[p.name]
			values, err := toReleaseValues("", p.valuesMap)
			if err != nil {
				return err
			}
			_, err = uc.release.Install(config, p.namespace, p.name, false, ref, labels, p.labels, p.annotations, values)
			return err
		})
	}
	return eg.Wait()
}

func (uc *OrchestratorUseCase) UpgradeExtensions(ctx context.Context, scope, facility string, chartRefMap map[string]string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return err
	}

	names := slices.Collect(maps.Keys(chartRefMap))
	extensions := uc.filterInternalExtensions(names)

	eg, _ := errgroup.WithContext(ctx)
	for _, p := range extensions {
		eg.Go(func() error {
			ref := chartRefMap[p.name]
			values, err := toReleaseValues("", p.valuesMap)
			if err != nil {
				return err
			}
			_, err = uc.release.Upgrade(config, p.namespace, p.name, false, ref, values, true)
			return err
		})
	}
	return eg.Wait()
}

func (uc *OrchestratorUseCase) listExtensions(ctx context.Context, scope, facility string, names []string) ([]Extension, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	extensions := uc.filterInternalExtensions(names)
	charts := make([][]Chart, len(extensions))
	releases := make([][]Release, len(extensions))

	eg, egctx := errgroup.WithContext(ctx)
	for i := range extensions {
		eg.Go(func() error {
			v, err := uc.chart.List(egctx, extensions[i].repoURL, true)
			if err == nil {
				charts[i] = v
			}
			return err
		})
		eg.Go(func() error {
			selector := fmt.Sprintf("%s=%s", TypeLabel, "extension")
			v, err := uc.release.List(config, extensions[i].namespace, selector)
			if err == nil {
				releases[i] = v
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return uc.filterExtensions(flatten(charts), flatten(releases), names)
}

func (uc *OrchestratorUseCase) filterInternalExtensions(extensions []string) []extension {
	result := []extension{}
	for _, p := range extensions {
		if plg, ok := extensionMap[p]; ok {
			result = append(result, plg)
		}
	}
	return result
}

func (uc *OrchestratorUseCase) filterExtensions(charts []Chart, releases []Release, extensions []string) ([]Extension, error) {
	result := []Extension{}
	for _, extension := range extensions {
		latest, url, err := findLatestChart(charts, extension)
		if err != nil {
			return nil, err
		}
		result = append(result, Extension{
			Release:   findRelease(releases, extension),
			Latest:    latest,
			LatestURL: url,
		})
	}
	return result, nil
}

func findRelease(releases []Release, name string) *Release {
	idx := slices.IndexFunc(releases, func(r Release) bool {
		return r.Chart != nil && r.Chart.Name() == name
	})
	if idx == -1 {
		return nil
	}
	return &releases[idx]
}

func findLatestChart(charts []Chart, name string) (latest *ChartVersion, ref string, err error) {
	idx := slices.IndexFunc(charts, func(c Chart) bool {
		return c.Name == name
	})
	if idx == -1 {
		return nil, "", fmt.Errorf("chart not found for %q", name)
	}

	chart := charts[idx]
	if len(chart.Versions) == 0 {
		return nil, "", fmt.Errorf("no versions available for %q", name)
	}

	slices.SortFunc(chart.Versions, func(a, b *ChartVersion) int {
		return semver.Compare(b.Version, a.Version)
	})

	latest = chart.Versions[0]
	if len(latest.URLs) == 0 {
		return nil, "", fmt.Errorf("no chart URL found for %q", name)
	}
	return latest, latest.URLs[0], nil
}
