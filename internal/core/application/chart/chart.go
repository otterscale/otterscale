package chart

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"github.com/otterscale/otterscale/internal/config"
	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const RepoURL = "https://otterscale.github.io/charts"

const (
	localRepoDir   = "./charts"
	localOCIFormat = "oci://%s:32000/charts"
)

type (
	// Version represents a Helm ChartVersion resource.
	Version = repo.ChartVersion

	// Metadata represents Helm Chart Metadata.
	Metadata = chart.Metadata
)

type Chart struct {
	Name     string
	Versions []*Version
}

type File struct {
	ReadmeMarkdown string
	ValuesYAML     string
	Customization  map[string]any
}

type ChartRepo interface {
	List(ctx context.Context, url string, useCache bool) ([]Chart, error)
	Show(chartRef string, format action.ShowOutputFormat) (string, error)
	Push(chartRef, remoteOCI string) (string, error)
	Index(dir, url string) error
	GetStableVersion(ctx context.Context, url, name string, useCache bool) (*Version, error)
}

type ChartUseCase struct {
	conf  *config.Config
	chart ChartRepo
}

func NewChartUseCase(conf *config.Config, chart ChartRepo) *ChartUseCase {
	return &ChartUseCase{
		conf:  conf,
		chart: chart,
	}
}

func (uc *ChartUseCase) ListCharts(ctx context.Context) ([]Chart, error) {
	urls := slices.Clone(uc.conf.Kube.HelmRepositoryURLs)

	exists, err := checkDirExists(localRepoDir)
	if err != nil {
		return nil, err
	}

	if exists {
		urls = append(urls, localRepoDir)
	}

	ret := make([][]Chart, len(urls))
	eg, egctx := errgroup.WithContext(ctx)

	for i := range urls {
		eg.Go(func() error {
			useCache := urls[i] != localRepoDir

			v, err := uc.chart.List(egctx, urls[i], useCache)
			if err == nil {
				ret[i] = v
			}
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return flatten(ret), nil
}

func (uc *ChartUseCase) GetChartFile(chartRef string) (*File, error) {
	file := &File{}
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

// TODO: multiple service on kubernetes
func (uc *ChartUseCase) UploadChart(chartContent []byte) error {
	if err := os.MkdirAll(localRepoDir, 0o700); err != nil { //nolint:mnd // default folder permission
		return err
	}

	name, version, err := uc.extractMetadata(chartContent)
	if err != nil {
		return err
	}

	fileName := filepath.Join(localRepoDir, fmt.Sprintf("%s-%s.tgz", name, version))
	if err := os.WriteFile(fileName, chartContent, 0o600); err != nil { //nolint:mnd // default file permission
		return err
	}

	localOCI, err := uc.localOCI()
	if err != nil {
		return err
	}

	if _, err := uc.chart.Push(fileName, localOCI); err != nil {
		return err
	}

	return uc.chart.Index(localRepoDir, localOCI)
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

func (uc *ChartUseCase) localOCI() (string, error) {
	config, err := uc.newMicroK8sConfig()
	if err != nil {
		return "", err
	}

	url, err := url.Parse(config.Host)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(localOCIFormat, url.Hostname()), nil
}

func (uc *ChartUseCase) extractMetadata(chartContent []byte) (name, version string, err error) {
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

func checkDirExists(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

func flatten[T any](data [][]T) []T {
	totalLen := 0

	for _, innerSlice := range data {
		totalLen += len(innerSlice)
	}

	ret := make([]T, 0, totalLen)

	for _, innerSlice := range data {
		ret = append(ret, innerSlice...)
	}

	return ret
}
