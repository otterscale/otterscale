package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/moby/moby/pkg/namesgenerator"
	"golang.org/x/sync/errgroup"

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/strvals"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

type Release struct {
	ScopeName    string
	ScopeUUID    string
	FacilityName string
	*release.Release
}

type ReleaseRepo interface {
	List(config *rest.Config, namespace string) ([]release.Release, error)
	Install(config *rest.Config, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error)
	Uninstall(config *rest.Config, namespace, name string, dryRun bool) (*release.Release, error)
	Upgrade(config *rest.Config, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error)
	Rollback(config *rest.Config, namespace, name string, dryRun bool) error
	GetValues(config *rest.Config, namespace, name string) (map[string]any, error)
}

func (uc *ApplicationUseCase) ListReleases(ctx context.Context) ([]Release, error) {
	kuberneteses, err := listKuberneteses(ctx, uc.scope, uc.client, "")
	if err != nil {
		return nil, err
	}
	eg, ctx := errgroup.WithContext(ctx)
	result := make([][]Release, len(kuberneteses))
	for i := range kuberneteses {
		eg.Go(func() error {
			config, err := uc.config(ctx, kuberneteses[i].ScopeUUID, kuberneteses[i].Name)
			if err != nil {
				return err
			}
			releases, _ := uc.release.List(config, "")
			for _, release := range releases {
				result[i] = append(result[i], Release{
					ScopeName:    kuberneteses[i].ScopeName,
					ScopeUUID:    kuberneteses[i].ScopeUUID,
					FacilityName: kuberneteses[i].Name,
					Release:      &release,
				})
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	release := []Release{}
	for _, r := range result {
		release = append(release, r...)
	}
	return release, nil
}

func (uc *ApplicationUseCase) CreateRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool, chartRef, valuesYAML string, valuesMap map[string]string) (*Release, error) {
	values, err := toReleaseValues(valuesYAML, valuesMap)
	if err != nil {
		return nil, err
	}

	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	release, err := uc.release.Install(config, namespace, getReleaseName(name), dryRun, chartRef, values)
	if err != nil {
		return nil, err
	}
	return &Release{Release: release}, nil
}

func (uc *ApplicationUseCase) UpdateRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool, chartRef, valuesYAML string) (*Release, error) {
	values, err := toReleaseValues(valuesYAML, nil)
	if err != nil {
		return nil, err
	}

	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return nil, err
	}
	release, err := uc.release.Upgrade(config, namespace, name, dryRun, chartRef, values)
	if err != nil {
		return nil, err
	}
	return &Release{Release: release}, nil
}

func (uc *ApplicationUseCase) DeleteRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	_, err = uc.release.Uninstall(config, namespace, name, dryRun)
	return err
}

func (uc *ApplicationUseCase) RollbackRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool) error {
	config, err := uc.config(ctx, uuid, facility)
	if err != nil {
		return err
	}
	return uc.release.Rollback(config, namespace, name, dryRun)
}

func getReleaseName(name string) string {
	if name != "" {
		return name
	}
	return strings.ReplaceAll(namesgenerator.GetRandomName(0), "_", "-")
}

func toReleaseValues(valuesYAML string, valuesMap map[string]string) (map[string]any, error) {
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
