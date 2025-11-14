package extension

import (
	"context"
	"fmt"
	"sync"

	"github.com/otterscale/otterscale/internal/core/application/chart"
	"github.com/otterscale/otterscale/internal/core/application/release"
	"golang.org/x/sync/errgroup"
)

type Extension struct {
	Release   *release.Release
	Latest    *chart.Version
	LatestURL string
}

type ExtensionUseCase struct {
	chart   chart.ChartRepo
	release release.ReleaseRepo
}

func NewExtensionUseCase(chart chart.ChartRepo, release release.ReleaseRepo) *ExtensionUseCase {
	return &ExtensionUseCase{
		chart:   chart,
		release: release,
	}
}

func (uc *ExtensionUseCase) ListGeneralExtensions(ctx context.Context, scope string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, general)
}

func (uc *ExtensionUseCase) ListModelExtensions(ctx context.Context, scope string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, model)
}

func (uc *ExtensionUseCase) ListInstanceExtensions(ctx context.Context, scope string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, instance)
}

func (uc *ExtensionUseCase) ListStorageExtensions(ctx context.Context, scope string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, storage)
}

func (uc *ExtensionUseCase) InstallExtensions(ctx context.Context, scope string, chartRefMap map[string]string) error {
	eg, _ := errgroup.WithContext(ctx)

	for name, chartRef := range chartRefMap {
		eg.Go(func() error {
			base, err := uc.base(name)
			if err != nil {
				return err
			}

			_, err = uc.release.Install(scope, base.Namespace, base.Name, false, chartRef, base.Labels, base.Labels, base.Annotations, "", base.ValuesMap)
			return err
		})
	}

	return eg.Wait()
}

func (uc *ExtensionUseCase) UpgradeExtensions(ctx context.Context, scope string, chartRefMap map[string]string) error {
	eg, _ := errgroup.WithContext(ctx)

	for name, chartRef := range chartRefMap {
		eg.Go(func() error {
			base, err := uc.base(name)
			if err != nil {
				return err
			}

			_, err = uc.release.Upgrade(scope, base.Namespace, base.Name, false, chartRef, "", base.ValuesMap, true)
			return err
		})
	}

	return eg.Wait()
}

func (uc *ExtensionUseCase) listExtensions(ctx context.Context, scope string, bases []base) ([]Extension, error) {

	versions := sync.Map{}
	releases := sync.Map{}

	eg, egctx := errgroup.WithContext(ctx)

	for i := range bases {
		eg.Go(func() error {
			v, err := uc.chart.GetStableVersion(egctx, bases[i].RepoURL, bases[i].Name, true)
			if err == nil {
				versions.Store(i, v)
			}
			return err
		})

		eg.Go(func() error {
			v, err := uc.release.Get(scope, bases[i].Namespace, bases[i].Name)
			if err == nil {
				releases.Store(i, v)
			}
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	ret := []Extension{}

	for _, base := range bases {
		r, ok := releases.Load(base.Name)
		if !ok {
			return nil, fmt.Errorf("failed to load release for %q", base.Name)
		}

		v, ok := versions.Load(base.Name)
		if !ok {
			return nil, fmt.Errorf("failed to load chart version for %q", base.Name)
		}

		latest := v.(*chart.Version)

		if len(latest.URLs) == 0 {
			return nil, fmt.Errorf("no chart URL found for %q", base.Name)
		}

		ret = append(ret, Extension{
			Release:   r.(*release.Release),
			Latest:    latest,
			LatestURL: latest.URLs[0],
		})
	}

	return ret, nil
}
