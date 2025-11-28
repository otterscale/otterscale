package release

import (
	"context"

	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"sigs.k8s.io/yaml"

	"github.com/otterscale/otterscale/internal/core/application/chart"
)

const (
	TypeLabel          = "otterscale.com/type"
	ReleaseNameLabel   = "otterscale.com/release-name"
	ChartRefAnnotation = "otterscale.com/chart-ref"
)

// Release represents a Helm Release resource.
type Release = release.Release

//nolint:revive // allows this exported interface name for specific domain clarity.
type ReleaseRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]Release, error)
	Get(ctx context.Context, scope, namespace, name string) (*Release, error)
	Install(ctx context.Context, scope, namespace, name string, dryRun bool, chartRef string, labelsInSecrets, labels, annotations map[string]string, valuesYAML string, valuesMap map[string]string) (*Release, error)
	Uninstall(ctx context.Context, scope, namespace, name string, dryRun bool) (*Release, error)
	Upgrade(ctx context.Context, scope, namespace, name string, dryRun bool, chartRef string, valuesYAML string, valuesMap map[string]string, reuseValues bool) (*Release, error)
	Rollback(ctx context.Context, scope, namespace, name string, dryRun bool) error
	GetValues(ctx context.Context, scope, namespace, name string) (map[string]any, error)
}

type UseCase struct {
	release ReleaseRepo

	chart chart.ChartRepo
}

func NewUseCase(release ReleaseRepo, chart chart.ChartRepo) *UseCase {
	return &UseCase{
		release: release,
		chart:   chart,
	}
}

func (uc *UseCase) ListReleases(ctx context.Context, scope string) ([]Release, error) {
	selector := "!" + TypeLabel

	return uc.release.List(ctx, scope, "", selector)
}

func (uc *UseCase) CreateRelease(ctx context.Context, scope, namespace, name string, dryRun bool, chartRef, valuesYAML string, valuesMap map[string]string) (*Release, error) {
	// labels
	labels := map[string]string{
		ReleaseNameLabel: name,
	}

	return uc.release.Install(ctx, scope, namespace, name, dryRun, chartRef, nil, labels, nil, valuesYAML, valuesMap)
}

func (uc *UseCase) UpdateRelease(ctx context.Context, scope, namespace, name string, dryRun bool, chartRef, valuesYAML string) (*Release, error) {
	return uc.release.Upgrade(ctx, scope, namespace, name, dryRun, chartRef, valuesYAML, nil, false)
}

func (uc *UseCase) DeleteRelease(ctx context.Context, scope, namespace, name string, dryRun bool) error {
	_, err := uc.release.Uninstall(ctx, scope, namespace, name, dryRun)
	return err
}

func (uc *UseCase) RollbackRelease(ctx context.Context, scope, namespace, name string, dryRun bool) error {
	return uc.release.Rollback(ctx, scope, namespace, name, dryRun)
}

func (uc *UseCase) GetChartFileFromApplication(ctx context.Context, scope, namespace string, labels, annotations map[string]string) (*chart.File, error) {
	file := &chart.File{}
	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		releaseName, ok := labels[ReleaseNameLabel]
		if ok {
			v, err := uc.release.GetValues(egctx, scope, namespace, releaseName)
			if err != nil {
				return err
			}

			valuesYAML, _ := yaml.Marshal(v) // ignore error
			file.ValuesYAML = string(valuesYAML)
		}

		return nil
	})

	eg.Go(func() error {
		chartRef, ok := annotations[ChartRefAnnotation]
		if ok {
			v, err := uc.chart.Show(egctx, chartRef, action.ShowReadme)
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
