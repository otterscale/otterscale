package chart

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/registry"
)

type (
	// Metadata represents Helm Chart Metadata resource.
	Metadata = chart.Metadata

	// Maintainer represents Helm Chart Maintainer resource.
	Maintainer = chart.Maintainer

	// Dependency represents Helm Chart Dependency resource.
	Dependency = chart.Dependency
)

type Chart struct {
	*chart.Metadata
	Repository string
}

type Version struct {
	ChartRef           string
	ChartVersion       string
	ApplicationVersion string
}

type Information struct {
	Readme string
	Values string
}

//nolint:revive // allows this exported interface name for specific domain clarity.
type ChartRepo interface {
	Show(ctx context.Context, chartRef string, format action.ShowOutputFormat) (string, error)
	Pull(ctx context.Context, chartRef, localDir string) error
	Push(ctx context.Context, chartPath, remoteOCI string) error
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

func (uc *UseCase) SyncArtifactHub(ctx context.Context, remoteOCI string) error {
	packages, err := fetchPackages(ctx, artifactHubSearchURL)
	if err != nil {
		return err
	}

	chartRefs, err := uc.fetchChartRefs(ctx, packages)
	if err != nil {
		return err
	}

	return uc.syncCharts(ctx, chartRefs, remoteOCI)
}

func (uc *UseCase) Import(ctx context.Context, chartRef, remoteOCI string) error {
	return uc.importChart(ctx, chartRef, remoteOCI)
}

func (uc *UseCase) fetchChartRefs(ctx context.Context, packages []Package) ([]string, error) {
	chartRefs := make([]string, len(packages))

	eg, egctx := errgroup.WithContext(ctx)

	for i := range packages {
		eg.Go(func() error {
			url := fmt.Sprintf(artifactHubPackageTemplate, packages[i].Repository.Name, packages[i].Name)
			chartRef, err := fetchContentURL(egctx, url)
			if err == nil {
				chartRefs[i] = chartRef
			}
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return chartRefs, nil
}

func (uc *UseCase) syncCharts(ctx context.Context, chartRefs []string, remoteOCI string) error {
	eg, egctx := errgroup.WithContext(ctx)

	for i := range chartRefs {
		eg.Go(func() error {
			return uc.syncChart(egctx, chartRefs[i], remoteOCI)
		})
	}

	return eg.Wait()
}

func (uc *UseCase) syncChart(ctx context.Context, chartRef, remoteOCI string) error {
	destDir, err := os.MkdirTemp("", "chart-sync-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(destDir)

	if err := uc.chart.Pull(ctx, chartRef, destDir); err != nil {
		slog.Error("failed to pull chart", "error", err, "chartRef", chartRef)
		return nil // continue on error
	}

	name := filepath.Base(chartRef)

	if strings.HasPrefix(chartRef, registry.OCIScheme) {
		idx := strings.LastIndexByte(name, ':')
		name = fmt.Sprintf("%s-%s.tgz", name[:idx], name[idx+1:])
	}

	chartPath := filepath.Join(destDir, name)

	return uc.chart.Push(ctx, chartPath, fmt.Sprintf("oci://%s/otterscale", remoteOCI))
}

// TODO: with repo name
func (uc *UseCase) importChart(ctx context.Context, chartRef, remoteOCI string) error {
	destDir, err := os.MkdirTemp("", "chart-import-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(destDir)

	if err := uc.chart.Pull(ctx, chartRef, destDir); err != nil {
		return fmt.Errorf("failed to pull chart: %s", err)
	}

	name := filepath.Base(chartRef)

	if strings.HasPrefix(chartRef, registry.OCIScheme) {
		idx := strings.LastIndexByte(name, ':')
		name = fmt.Sprintf("%s-%s.tgz", name[:idx], name[idx+1:])
	}

	chartPath := filepath.Join(destDir, name)

	return uc.chart.Push(ctx, chartPath, fmt.Sprintf("oci://%s", remoteOCI))
}
