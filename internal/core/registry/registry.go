package registry

import (
	"context"
	"fmt"

	"github.com/otterscale/otterscale/internal/core/registry/chart"
)

type UseCase struct {
	manifest   ManifestRepo
	repository RepositoryRepo
}

func NewUseCase(manifest ManifestRepo, repository RepositoryRepo) *UseCase {
	return &UseCase{
		manifest:   manifest,
		repository: repository,
	}
}

func (uc *UseCase) GetRegistryURL(scope string) (string, error) {
	return uc.repository.GetRegistryURL(scope)
}

func (uc *UseCase) ListRepositories(ctx context.Context, scope string) ([]Repository, error) {
	return uc.repository.List(ctx, scope)
}

func (uc *UseCase) ListManifests(ctx context.Context, scope, repository string) ([]Manifest, error) {
	return uc.manifest.List(ctx, scope, repository)
}

func (uc *UseCase) DeleteManifest(ctx context.Context, scope, repository, digest string) error {
	return uc.manifest.Delete(ctx, scope, repository, digest)
}

func (uc *UseCase) ListCharts(ctx context.Context, scope string) ([]chart.Chart, error) {
	repos, err := uc.repository.List(ctx, scope)
	if err != nil {
		return nil, err
	}

	charts := []chart.Chart{}

	for _, repo := range repos {
		manifest, err := uc.manifest.Get(ctx, scope, repo.Name, repo.LatestTag)
		if err != nil {
			return nil, err
		}

		chart := manifest.Chart
		if chart == nil {
			continue
		}

		charts = append(charts, *chart)
	}

	return charts, nil
}

func (uc *UseCase) ListChartVersions(ctx context.Context, scope, chartName string) ([]chart.Version, error) {
	registryURL, err := uc.repository.GetRegistryURL(scope)
	if err != nil {
		return nil, err
	}

	manifests, err := uc.manifest.List(ctx, scope, chartName)
	if err != nil {
		return nil, err
	}

	versions := []chart.Version{}

	for i := range manifests {
		c := manifests[i].Chart
		if c != nil {
			versions = append(versions, chart.Version{
				ChartRef:           fmt.Sprintf("oci://%s/%s:%s", registryURL, c.Name, c.Version),
				ChartVersion:       c.Version,
				ApplicationVersion: c.AppVersion,
			})
		}
	}

	return versions, nil
}
