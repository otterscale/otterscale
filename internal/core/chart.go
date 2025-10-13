package core

import (
	"context"
	"fmt"
	"sync"

	"connectrpc.com/connect"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"
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

type (
	ChartDependency = chart.Dependency
	ChartMaintainer = chart.Maintainer
	ChartMetadata   = chart.Metadata
	ChartVersion    = repo.ChartVersion
)

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

func (uc *ChartUseCase) GetChartFile(ctx context.Context, chartRef string) (*ChartFile, error) {
	file := &ChartFile{}
	wg := sync.WaitGroup{}
	wg.Go(func() {
		file.ValuesYAML, _ = uc.chart.Show(chartRef, action.ShowValues)
	})
	wg.Go(func() {
		file.ReadmeMarkdown, _ = uc.chart.Show(chartRef, action.ShowReadme)
	})
	wg.Wait()
	return file, nil
}

func (uc *ChartUseCase) GetChartFileFromApplication(ctx context.Context, uuid, facility string, app *Application) (*ChartFile, error) {
	file := &ChartFile{}
	wg := sync.WaitGroup{}
	wg.Go(func() {
		// FIXME: invalid label
		// if releaseName, ok := app.Labels["app.otterscale.com/release-name"]; ok {
		// 	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
		// 	if err != nil {
		// 		return
		// 	}
		// 	v, err := uc.release.GetValues(config, app.Namespace, releaseName)
		// 	if err != nil {
		// 		return
		// 	}
		// 	valuesYAML, _ := yaml.Marshal(v)
		// 	metadata.ValuesYAML = string(valuesYAML)
		// }
	})
	wg.Go(func() {
		// FIXME: invalid label format
		// if chartRef, ok := app.Labels["app.otterscale.com/chart-ref"]; ok {
		// 	v, err := uc.chart.Show(chartRef, action.ShowReadme)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	metadata.ReadmeMD = v
		// }
	})
	wg.Wait()
	return file, nil
}
