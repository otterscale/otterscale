package helm

import (
	"context"

	"helm.sh/helm/v3/pkg/action"

	"github.com/otterscale/otterscale/internal/core/registry/chart"
)

// Note: Helm API do not support context.
type chartRepo struct {
	helm *Helm
}

func NewChartRepo(helm *Helm) chart.ChartRepo {
	return &chartRepo{
		helm: helm,
	}
}

var _ chart.ChartRepo = (*chartRepo)(nil)

func (r *chartRepo) Show(ctx context.Context, chartRef string, format action.ShowOutputFormat) (string, error) {
	client := action.NewShow(format)
	client.SetRegistryClient(r.helm.registryClient)

	chartPath, err := client.LocateChart(chartRef, r.helm.envSettings)
	if err != nil {
		return "", err
	}

	return client.Run(chartPath)
}
