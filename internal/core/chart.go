package core

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/otterscale/otterscale/internal/config"
)

const (
	LocalChartRepoDir = "./charts"
	remoteOCIFormat   = "oci://%s:32000/charts"
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
	Push(chartRef, remoteOCI string) (string, error)
	Index(dir, url string) error
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

// TODO: multiple service on kubernetes
func (uc *ChartUseCase) UploadChart(chartContent []byte) error {
	if err := os.MkdirAll(LocalChartRepoDir, 0o700); err != nil { //nolint:mnd // default folder permission
		return err
	}

	name, version, err := extractChartMetadata(chartContent)
	if err != nil {
		return err
	}

	fileName := filepath.Join(LocalChartRepoDir, fmt.Sprintf("%s-%s.tgz", name, version))
	if err := os.WriteFile(fileName, chartContent, 0o600); err != nil { //nolint:mnd // default file permission
		return err
	}

	remoteOCI, err := uc.remoteOCI()
	if err != nil {
		return err
	}

	if _, err := uc.chart.Push(fileName, remoteOCI); err != nil {
		return err
	}

	return uc.chart.Index(LocalChartRepoDir, remoteOCI)
}

// TODO: replace with remote flag
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

func (uc *ChartUseCase) remoteOCI() (string, error) {
	config, err := uc.newMicroK8sConfig()
	if err != nil {
		return "", err
	}
	url, err := url.Parse(config.Host)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(remoteOCIFormat, url.Hostname()), nil
}

func extractChartMetadata(chartContent []byte) (name, version string, err error) {
	reader := bytes.NewReader(chartContent)
	chart, err := loader.LoadArchive(reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to load chart archive: %w", err)
	}

	metadata := chart.Metadata
	if metadata == nil {
		return "", "", fmt.Errorf("chart metadata not found")
	}
	if metadata.Name == "" {
		return "", "", fmt.Errorf("chart name not found")
	}
	if metadata.Version == "" {
		return "", "", fmt.Errorf("chart version not found")
	}
	return metadata.Name, metadata.Version, nil
}
