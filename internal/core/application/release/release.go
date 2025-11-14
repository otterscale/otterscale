package release

import (
	"context"

	"github.com/goccy/go-yaml"
	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"

	"github.com/otterscale/otterscale/internal/core/application/chart"
)

const (
	TypeLabel        = "otterscale.com/type"
	ReleaseNameLabel = "otterscale.com/release-name"
	ChartRefKey      = "chart-ref"
)

// Release represents a Helm Release resource.
type Release = release.Release

type ReleaseRepo interface {
	List(scope, namespace, selector string) ([]Release, error)
	Get(scope, namespace, name string) (*Release, error)
	Install(scope, namespace, name string, dryRun bool, chartRef string, labelsInSecrets, labels, annotations map[string]string, valuesYAML string, valuesMap map[string]string) (*Release, error)
	Uninstall(scope, namespace, name string, dryRun bool) (*Release, error)
	Upgrade(scope, namespace, name string, dryRun bool, chartRef string, valuesYAML string, valuesMap map[string]string, reuseValues bool) (*Release, error)
	Rollback(scope, namespace, name string, dryRun bool) error
	GetValues(scope, namespace, name string) (map[string]any, error)
}

type ReleaseUseCase struct {
	release ReleaseRepo

	chart chart.ChartRepo
}

func NewReleaseUseCase(release ReleaseRepo, chart chart.ChartRepo) *ReleaseUseCase {
	return &ReleaseUseCase{
		release: release,
		chart:   chart,
	}
}

func (uc *ReleaseUseCase) ListReleases(ctx context.Context, scope string) ([]Release, error) {
	selector := "!" + TypeLabel

	return uc.release.List(scope, "", selector)
}

func (uc *ReleaseUseCase) CreateRelease(ctx context.Context, scope, namespace, name string, dryRun bool, chartRef, valuesYAML string, valuesMap map[string]string) (*Release, error) {
	// labels
	labels := map[string]string{
		ReleaseNameLabel: name,
	}

	return uc.release.Install(scope, namespace, name, dryRun, chartRef, nil, labels, nil, valuesYAML, valuesMap)
}

func (uc *ReleaseUseCase) UpdateRelease(ctx context.Context, scope, namespace, name string, dryRun bool, chartRef, valuesYAML string) (*Release, error) {
	return uc.release.Upgrade(scope, namespace, name, dryRun, chartRef, valuesYAML, nil, false)
}

func (uc *ReleaseUseCase) DeleteRelease(ctx context.Context, scope, namespace, name string, dryRun bool) error {
	_, err := uc.release.Uninstall(scope, namespace, name, dryRun)
	return err
}

func (uc *ReleaseUseCase) RollbackRelease(ctx context.Context, scope, namespace, name string, dryRun bool) error {
	return uc.release.Rollback(scope, namespace, name, dryRun)
}

func (uc *ReleaseUseCase) GetChartFileFromApplication(ctx context.Context, scope, namespace string, labels map[string]string) (*chart.File, error) {
	file := &chart.File{}
	eg := errgroup.Group{}

	eg.Go(func() error {
		releaseName, ok := labels[ReleaseNameLabel]
		if ok {
			v, err := uc.release.GetValues(scope, namespace, releaseName)
			if err != nil {
				return err
			}

			valuesYAML, _ := yaml.Marshal(v) // ignore error
			file.ValuesYAML = string(valuesYAML)
		}

		return nil
	})

	eg.Go(func() error {
		releaseName, ok := labels[ReleaseNameLabel]
		if ok {
			rel, err := uc.release.Get(scope, namespace, releaseName)
			if err != nil {
				return err
			}

			chart := rel.Chart
			if chart == nil {
				return nil // skip if chart is nil
			}

			chartRef := ""
			if v, ok := rel.Config[ChartRefKey]; ok {
				if str, ok := v.(string); ok {
					chartRef = str
				}
			}
			if chartRef == "" {
				return nil // skip if chartRef is empty
			}

			v, err := uc.chart.Show(chartRef, action.ShowReadme)
			if err != nil {
				return err
			}

			file.ReadmeMarkdown = v
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return file, nil
}
