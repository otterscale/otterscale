package core

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/otterscale/otterscale/internal/config"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
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
	UploadChart(ctx context.Context, ociRegistryURL, chartName, chartVersion string, chartContent []byte) error
}

type ChartUseCase struct {
	action   ActionRepo
	chart    ChartRepo
	facility FacilityRepo
	release  ReleaseRepo
	conf     *config.Config
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

func (uc *ChartUseCase) newMicroK8sConfig() (*rest.Config, error) {
	kubeConfig, err := base64.StdEncoding.DecodeString(uc.conf.MicroK8s.Config)
	if err != nil {
		return nil, err
	}
	configAPI, err := clientcmd.Load(kubeConfig)
	if err != nil {
		return nil, err
	}
	return clientcmd.NewDefaultClientConfig(*configAPI, &clientcmd.ConfigOverrides{}).ClientConfig()
}

func (uc *ChartUseCase) UploadChart(ctx context.Context, chartContent []byte) error {
	// Parse chart metadata from the chart content to extract name and version
	chartName, chartVersion, err := extractChartMetadata(chartContent)
	if err != nil {
		return fmt.Errorf("failed to extract chart metadata: %w", err)
	}

	microK8sConfig, err := uc.newMicroK8sConfig()
	if err != nil {
		return err
	}

	// Convert https://IP:PORT to oci://IP:32000/charts
	host := microK8sConfig.Host
	host = strings.TrimPrefix(host, "https://")
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}
	ociRegistryURL := fmt.Sprintf("oci://%s:32000/charts", host)

	return uc.chart.UploadChart(ctx, ociRegistryURL, chartName, chartVersion, chartContent)
}

func extractChartMetadata(chartContent []byte) (string, string, error) {
	reader := bytes.NewReader(chartContent)
	chart, err := loader.LoadArchive(reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to load chart archive: %w", err)
	}

	if chart.Metadata == nil {
		return "", "", fmt.Errorf("chart metadata is nil")
	}

	if chart.Metadata.Name == "" {
		return "", "", fmt.Errorf("chart name not found in Chart.yaml")
	}
	if chart.Metadata.Version == "" {
		return "", "", fmt.Errorf("chart version not found in Chart.yaml")
	}

	return chart.Metadata.Name, chart.Metadata.Version, nil
}
