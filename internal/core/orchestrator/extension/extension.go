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
	var (
		latest  *chart.Version
		release *release.Release
	)

	eg, egctx := errgroup.WithContext(ctx)

	for i := range bases {
		eg.Go(func() error {
			v, err := uc.chart.GetStableVersion(egctx, bases[i].RepoURL, bases[i].Name, true)
			if err == nil {
				latest = v
			}
			return err
		})

		eg.Go(func() error {
			v, err := uc.release.Get(egctx, scope, bases[i].Namespace, bases[i].Name)
			if err == nil {
				release = v
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

	for _, base := range bases {
		if len(latest.URLs) == 0 {
			return nil, fmt.Errorf("no chart URL found for %q", base.Name)
		}

		ret = append(ret, Extension{
			Release:   release,
			Latest:    latest,
			LatestURL: latest.URLs[0],
		})
	}

	return ret, nil
}
