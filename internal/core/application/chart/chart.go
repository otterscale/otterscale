package chart

import (
	"context"

	"helm.sh/helm/v3/pkg/repo"
)

type Chart struct {
	Name     string
	Versions repo.ChartVersions
}

type File struct {
	ReadmeMarkdown string
	ValuesYAML     string
	Customization  map[string]any
}

type ChartRepo interface {
	List(ctx context.Context, url string, useCache bool) ([]Chart, error)
	Show(chartRef string, format string) (string, error)
	Push(chartRef, remoteOCI string) (string, error)
	Index(dir, url string) error
}

type ChartUseCase struct {
	chart ChartRepo
}

func NewChartUseCase(chart ChartRepo) *ChartUseCase {
	return &ChartUseCase{
		chart: chart,
	}
}
