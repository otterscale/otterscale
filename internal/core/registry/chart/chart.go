package chart

import (
	"context"

	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
)

type (
	// Chart represents a Helm Chart Metadata resource.
	Chart = chart.Metadata

	// Maintainer represents Helm Chart Maintainer resource.
	Maintainer = chart.Maintainer

	// Dependency represents Helm Chart Dependency resource.
	Dependency = chart.Dependency
)

type Version struct {
	ChartRef           string
	ChartVersion       string
	ApplicationVersion string
}

type Information struct {
	Readme string
	Values string
}

type ChartRepo interface {
	Show(ctx context.Context, chartRef string, format action.ShowOutputFormat) (string, error)
}

type UseCase struct {
	chart ChartRepo
}

func NewUseCase(chart ChartRepo) *UseCase {
	return &UseCase{
		chart: chart,
	}
}

func (uc *UseCase) GetChartInformation(ctx context.Context, chartRef string) (*Information, error) {
	info := &Information{}
	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.chart.Show(egctx, chartRef, action.ShowValues)
		if err == nil {
			info.Values = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.chart.Show(egctx, chartRef, action.ShowReadme)
		if err == nil {
			info.Readme = v
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return info, nil
}
