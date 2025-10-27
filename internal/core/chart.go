package core

import (
	"context"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/otterscale/otterscale/internal/config"
)

type (
	ChartDependency = chart.Dependency
	ChartMaintainer = chart.Maintainer
	ChartMetadata   = chart.Metadata
	ChartVersion    = repo.ChartVersion
)

type Chart struct {
	Name     string
	Versions repo.ChartVersions
}

type ChartFile struct {
	ReadmeMarkdown string
	ValuesYAML     string
	Customization  map[string]any
}

type ChartRepo interface {
	List(ctx context.Context, url string) ([]Chart, error)
	Show(chartRef string, format action.ShowOutputFormat) (string, error)
}

type ChartUseCase struct {
	conf *config.Config

	action   ActionRepo
	chart    ChartRepo
	facility FacilityRepo
	release  ReleaseRepo
}

func NewChartUseCase(conf *config.Config, action ActionRepo, chart ChartRepo, facility FacilityRepo, release ReleaseRepo) *ChartUseCase {
	return &ChartUseCase{
		conf:     conf,
		action:   action,
		chart:    chart,
		facility: facility,
		release:  release,
	}
}

func (uc *ChartUseCase) ListCharts(ctx context.Context) ([]Chart, error) {
	urls := uc.conf.Kube.HelmRepositoryURLs
	eg, egctx := errgroup.WithContext(ctx)
	result := make([][]Chart, len(urls))
	for i := range urls {
		eg.Go(func() error {
			v, err := uc.chart.List(egctx, urls[i])
			if err == nil {
				result[i] = v
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return flatten(result), nil
}

func (uc *ChartUseCase) GetChartFile(chartRef string) (*ChartFile, error) {
	file := &ChartFile{}
	eg := errgroup.Group{}
	eg.Go(func() error {
		v, err := uc.chart.Show(chartRef, action.ShowValues)
		if err == nil {
			file.ValuesYAML = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.chart.Show(chartRef, action.ShowReadme)
		if err == nil {
			file.ReadmeMarkdown = v
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return file, nil
}

func (uc *ChartUseCase) GetChartFileFromApplication(ctx context.Context, scope, facility string, app *Application) (*ChartFile, error) {
	file := &ChartFile{}
	eg := errgroup.Group{}
	eg.Go(func() error {
		// FIXME: invalid label
		releaseName, ok := app.Labels["app.otterscale.com/release-name"]
		if ok {
			config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
			if err != nil {
				return err
			}
			v, err := uc.release.GetValues(config, app.Namespace, releaseName)
			if err != nil {
				return err
			}
			valuesYAML, _ := yaml.Marshal(v)
			file.ValuesYAML = string(valuesYAML)
		}
		return nil
	})
	eg.Go(func() error {
		// FIXME: invalid label format
		// if chartRef, ok := app.Labels["app.otterscale.com/chart-ref"]; ok {
		// 	v, err := uc.chart.Show(chartRef, action.ShowReadme)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	metadata.ReadmeMD = v
		// }
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return file, nil
}
