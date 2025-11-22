package chart

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core/scope"
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

	// Maintainer represents Helm Chart Maintainer.
	Maintainer = chart.Maintainer

	// Dependency represents Helm Chart Dependency.
	Dependency = chart.Dependency
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

//nolint:revive // allows this exported interface name for specific domain clarity.
type ChartRepo interface {
	List(ctx context.Context, url string, useCache bool) ([]Chart, error)
	Show(ctx context.Context, chartRef string, format action.ShowOutputFormat) (string, error)
	Push(ctx context.Context, chartRef, remoteOCI string) (string, error)
	Index(ctx context.Context, dir, url string) error
	LocalOCI(scope string) (string, error)
}

type UseCase struct {
	conf  *config.Config
	chart ChartRepo
}

func NewUseCase(conf *config.Config, chart ChartRepo) *UseCase {
	return &UseCase{
		conf:  conf,
		chart: chart,
	}
}

func (uc *UseCase) ListCharts(ctx context.Context) ([]Chart, error) {
	urls := uc.conf.KubeHelmRepositoryURLs()

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

func (uc *UseCase) GetChartFile(ctx context.Context, chartRef string) (*File, error) {
	file := &File{}
	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.chart.Show(egctx, chartRef, action.ShowValues)
		if err == nil {
			file.ValuesYAML = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.chart.Show(egctx, chartRef, action.ShowReadme)
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
func (uc *UseCase) UploadChart(ctx context.Context, chartContent []byte) error {
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

	if _, err := uc.chart.Push(ctx, fileName, localOCI); err != nil {
		return err
	}

	return uc.chart.Index(ctx, localRepoDir, localOCI)
}

func (uc *UseCase) localOCI() (string, error) {
	hostname, err := uc.chart.LocalOCI(scope.ReservedName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(localOCIFormat, hostname), nil
}

func (uc *UseCase) extractMetadata(chartContent []byte) (name, version string, err error) {
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
