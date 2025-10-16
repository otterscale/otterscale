package core

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"
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
	List(ctx context.Context) ([]Chart, error)
	Show(chartRef string, format action.ShowOutputFormat) (string, error)
}

type ChartUseCase struct {
	action   ActionRepo
	chart    ChartRepo
	facility FacilityRepo
	release  ReleaseRepo
}

func NewChartUseCase(action ActionRepo, chart ChartRepo, facility FacilityRepo, release ReleaseRepo) *ChartUseCase {
	return &ChartUseCase{
		action:   action,
		chart:    chart,
		facility: facility,
		release:  release,
	}
}

func (uc *ChartUseCase) ListCharts(ctx context.Context) ([]Chart, error) {
	return uc.chart.List(ctx)
}

func (uc *ChartUseCase) GetChart(ctx context.Context, chartName string) (*Chart, error) {
	charts, err := uc.chart.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range charts {
		if charts[i].Name != chartName {
			continue
		}
		return &Chart{
			Name:     charts[i].Name,
			Versions: charts[i].Versions,
		}, nil
	}
	return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("chart %q not found", chartName))
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
