package core

import (
	"context"

	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/repo"

	"sigs.k8s.io/yaml"
)

type Chart struct {
	Name     string
	Versions repo.ChartVersions
}

type ChartMetadata struct {
	ReadmeMD      string
	ValuesYAML    string
	Customization map[string]any
}

type ChartRepo interface {
	List(ctx context.Context) ([]repo.IndexFile, error)
	Show(chartRef string, format action.ShowOutputFormat) (string, error)
}

func (uc *ApplicationUseCase) GetChartMetadataFromApplication(ctx context.Context, uuid, facility string, app *Application) (*ChartMetadata, error) {
	metadata := &ChartMetadata{}
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if releaseName, ok := app.Labels["app.otterscale.io/release-name"]; ok {
			config, err := uc.config(ctx, uuid, facility)
			if err != nil {
				return err
			}
			v, err := uc.release.GetValues(config, app.Namespace, releaseName)
			if err != nil {
				return err
			}
			valuesYAML, _ := yaml.Marshal(v)
			metadata.ValuesYAML = string(valuesYAML)
		}
		return nil
	})
	eg.Go(func() error {
		if chartRef, ok := app.Labels["app.otterscale.io/chart-ref"]; ok {
			v, err := uc.chart.Show(chartRef, action.ShowReadme)
			if err != nil {
				return err
			}
			metadata.ReadmeMD = v
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return metadata, nil
}

func (uc *ApplicationUseCase) ListCharts(ctx context.Context) ([]Chart, error) {
	indice, err := uc.chart.List(ctx)
	if err != nil {
		return nil, err
	}
	charts := []Chart{}
	for i := range indice {
		for name := range indice[i].Entries {
			charts = append(charts, Chart{
				Name:     name,
				Versions: indice[i].Entries[name],
			})
		}
	}
	return charts, nil
}

func (uc *ApplicationUseCase) GetChart(ctx context.Context, chartName string) (*Chart, error) {
	indice, err := uc.chart.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range indice {
		for name := range indice[i].Entries {
			if name != chartName {
				continue
			}
			return &Chart{
				Name:     name,
				Versions: indice[i].Entries[name],
			}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "chart %q not found", chartName)
}

func (uc *ApplicationUseCase) GetChartMetadata(ctx context.Context, chartRef string) (*ChartMetadata, error) {
	metadata := &ChartMetadata{}
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		v, err := uc.chart.Show(chartRef, action.ShowValues)
		if err == nil {
			metadata.ValuesYAML = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.chart.Show(chartRef, action.ShowReadme)
		if err == nil {
			metadata.ReadmeMD = v
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return metadata, nil
}
