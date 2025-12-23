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

func (r *chartRepo) Show(_ context.Context, chartRef string, format action.ShowOutputFormat) (string, error) {
	registryClient, err := newRegistryClient()
	if err != nil {
		return "", err
	}

	config := &action.Configuration{
		RegistryClient: registryClient,
	}

	client := action.NewShowWithConfig(format, config)

	chartPath, err := client.LocateChart(chartRef, newEnvSettings())
	if err != nil {
		return "", err
	}

	return client.Run(chartPath)
}

func (r *chartRepo) Pull(_ context.Context, chartRef, localDir string) error {
	registryClient, err := newRegistryClient()
	if err != nil {
		return err
	}

	config := &action.Configuration{
		RegistryClient: registryClient,
	}

	client := action.NewPullWithOpts(
		action.WithConfig(config),
	)

	client.Settings = newEnvSettings()
	client.DestDir = localDir

	_, err = client.Run(chartRef)
	return err
}

func (r *chartRepo) Push(_ context.Context, chartPath, remoteOCI string) error {
	registryClient, err := newRegistryClient()
	if err != nil {
		return err
	}

	config := &action.Configuration{
		RegistryClient: registryClient,
	}

	client := action.NewPushWithOpts(
		action.WithPushConfig(config),
	)

	client.Settings = newEnvSettings()

	_, err = client.Run(chartPath, remoteOCI)
	return err
}
