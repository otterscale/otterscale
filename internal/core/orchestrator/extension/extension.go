package extension

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/storage/driver"

	"github.com/otterscale/otterscale/internal/core/application/chart"
	"github.com/otterscale/otterscale/internal/core/application/release"
)

type Extension struct {
	Release   *release.Release
	Latest    *chart.Version
	LatestURL string
}

type UseCase struct {
	chart   chart.ChartRepo
	release release.ReleaseRepo
}

func NewUseCase(chart chart.ChartRepo, release release.ReleaseRepo) *UseCase {
	return &UseCase{
		chart:   chart,
		release: release,
	}
}

func (uc *UseCase) ListGeneralExtensions(ctx context.Context, scope string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, general)
}

func (uc *UseCase) ListModelExtensions(ctx context.Context, scope string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, model)
}

func (uc *UseCase) ListInstanceExtensions(ctx context.Context, scope string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, instance)
}

func (uc *UseCase) ListStorageExtensions(ctx context.Context, scope string) ([]Extension, error) {
	return uc.listExtensions(ctx, scope, storage)
}

func (uc *UseCase) InstallExtensions(ctx context.Context, scope string, chartRefMap map[string]string) error {
	eg, egctx := errgroup.WithContext(ctx)

	for name, chartRef := range chartRefMap {
		eg.Go(func() error {
			base, err := uc.base(name)
			if err != nil {
				return err
			}

			_, err = uc.release.Install(egctx, scope, base.Namespace, base.Name, false, chartRef, base.Labels, base.Labels, base.Annotations, "", base.ValuesMap)
			return err
		})
	}

	return eg.Wait()
}

func (uc *UseCase) UpgradeExtensions(ctx context.Context, scope string, chartRefMap map[string]string) error {
	eg, egctx := errgroup.WithContext(ctx)

	for name, chartRef := range chartRefMap {
		eg.Go(func() error {
			base, err := uc.base(name)
			if err != nil {
				return err
			}

			_, err = uc.release.Upgrade(egctx, scope, base.Namespace, base.Name, false, chartRef, "", base.ValuesMap, true)
			return err
		})
	}

	return eg.Wait()
}

func (uc *UseCase) listExtensions(ctx context.Context, scope string, bases []base) ([]Extension, error) {
	versions := make([]chart.Version, len(bases))
	releases := make([]release.Release, len(bases))
	eg, egctx := errgroup.WithContext(ctx)

	for i := range bases {
		eg.Go(func() error {
			v, err := uc.chart.GetStableVersion(egctx, bases[i].RepoURL, bases[i].Name, true)
			if err == nil {
				versions[i] = *v
			}
			return err
		})

		eg.Go(func() error {
			v, err := uc.release.Get(egctx, scope, bases[i].Namespace, bases[i].Name)
			if err == nil {
				releases[i] = *v
			}
			if errors.Is(err, driver.ErrReleaseNotFound) {
				return nil
			}
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	ret := []Extension{}

	for i := range bases {
		if len(versions[i].URLs) == 0 {
			return nil, fmt.Errorf("no chart URL found for %q", bases[i].Name)
		}

		ret = append(ret, Extension{
			Release:   &releases[i],
			Latest:    &versions[i],
			LatestURL: versions[i].URLs[0],
		})
	}

	return ret, nil
}
