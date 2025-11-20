package extension

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/storage/driver"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/kyaml/filesys"

	"github.com/otterscale/otterscale/internal/core/application/chart"
	"github.com/otterscale/otterscale/internal/core/application/cluster"
	"github.com/otterscale/otterscale/internal/core/application/release"
)

type Manifest struct {
	ID      string
	Version string
}

type Extension struct {
	Name        string
	Description string
	Icon        string
	Status      string
	DeployedAt  *time.Time
	Current     *Manifest
	Latest      *Manifest
}

type UseCase struct {
	chart                    chart.ChartRepo
	customResourceDefinition cluster.CustomResourceDefinitionRepo
	release                  release.ReleaseRepo
}

func NewUseCase(chart chart.ChartRepo, customResourceDefinition cluster.CustomResourceDefinitionRepo, release release.ReleaseRepo) *UseCase {
	return &UseCase{
		chart:                    chart,
		customResourceDefinition: customResourceDefinition,
		release:                  release,
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

func (uc *UseCase) InstallExtensions(ctx context.Context, scope string, manifests []Manifest) error {
	eg, egctx := errgroup.WithContext(ctx)

	for _, manifest := range manifests {
		eg.Go(func() error {
			base, err := uc.base(manifest.ID)
			if err != nil {
				return err
			}

			for i := range base.Charts {
				chart := base.Charts[i]

				version, err := uc.chart.GetStableVersion(egctx, chart.RepoURL, base.ID, true)
				if err != nil {
					return err
				}

				if len(version.URLs) == 0 {
					return fmt.Errorf("no chart url found for %q version %q", base.ID, manifest.Version)
				}

				chartRef := version.URLs[0]

				if _, err := uc.release.Install(egctx, scope, chart.Namespace, base.Name, false, chartRef, chart.Labels, chart.Labels, chart.Annotations, "", chart.ValuesMap); err != nil {
					return err
				}
			}

			crd := base.CRD

			if crd != nil {
				fSys := filesys.MakeFsOnDisk()
				k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

				m, err := k.Run(fSys, uc.remoteCRDFormat(crd.RepoURL, manifest.Version))
				if err != nil {
					return err
				}

				for _, node := range m.Resources() {
					data, err := node.AsYAML()
					if err != nil {
						return err
					}

					var def *apiextensionsv1.CustomResourceDefinition
					if err := yaml.Unmarshal(data, &def); err != nil {
						return err
					}

					if _, err := uc.customResourceDefinition.Create(egctx, scope, def); err != nil {
						return err
					}
				}
			}

			return nil
		})
	}

	return eg.Wait()
}

func (uc *UseCase) UpgradeExtensions(ctx context.Context, scope string, manifests []Manifest) error {
	eg, egctx := errgroup.WithContext(ctx)

	for _, manifest := range manifests {
		eg.Go(func() error {
			base, err := uc.base(manifest.ID)
			if err != nil {
				return err
			}

			for i := range base.Charts {
				chart := base.Charts[i]

				version, err := uc.chart.GetStableVersion(egctx, chart.RepoURL, base.ID, true)
				if err != nil {
					return err
				}

				if len(version.URLs) == 0 {
					return fmt.Errorf("no chart url found for %q version %q", base.ID, manifest.Version)
				}

				chartRef := version.URLs[0]

				if _, err := uc.release.Upgrade(egctx, scope, chart.Namespace, base.Name, false, chartRef, "", chart.ValuesMap, true); err != nil {
					return err
				}
			}

			crd := base.CRD

			if crd != nil {
				fSys := filesys.MakeFsOnDisk()
				k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

				m, err := k.Run(fSys, uc.remoteCRDFormat(crd.RepoURL, manifest.Version))
				if err != nil {
					return err
				}

				for _, node := range m.Resources() {
					data, err := node.AsYAML()
					if err != nil {
						return err
					}

					var def *apiextensionsv1.CustomResourceDefinition
					if err := yaml.Unmarshal(data, &def); err != nil {
						return err
					}

					if _, err := uc.customResourceDefinition.Update(egctx, scope, def); err != nil {
						return err
					}
				}
			}

			return nil
		})
	}

	return eg.Wait()
}

func (uc *UseCase) listExtensions(ctx context.Context, scope string, bases []base) ([]Extension, error) {
	versions := make([]chart.Version, len(bases))
	releases := make([]release.Release, len(bases))
	resMaps := make([]resmap.ResMap, len(bases))
	crds := make([]cluster.CustomResourceDefinition, len(bases))

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.customResourceDefinition.List(egctx, scope, "")
		if err == nil {
			crds = v
		}
		return err
	})

	for i := range bases {
		base := bases[i]

		if len(base.Charts) > 0 {
			chart := base.Charts[0] // only get the parent chart

			eg.Go(func() error {
				v, err := uc.chart.GetStableVersion(egctx, chart.RepoURL, base.ID, true)
				if err == nil {
					versions[i] = *v
				}
				return err
			})

			eg.Go(func() error {
				v, err := uc.release.Get(egctx, scope, chart.Namespace, base.ID)
				if err == nil {
					releases[i] = *v
				}
				if errors.Is(err, driver.ErrReleaseNotFound) {
					return nil
				}
				return err
			})
		}

		crd := base.CRD

		if crd != nil {
			eg.Go(func() error {
				fSys := filesys.MakeFsOnDisk()
				k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

				resMap, err := k.Run(fSys, uc.remoteCRDFormat(crd.RepoURL, crd.Version))
				if err == nil {
					resMaps[i] = resMap
				}

				return err
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	ret := []Extension{}

	for i := range bases {
		base := bases[i]

		if len(base.Charts) > 0 {
			ret = append(ret, *uc.buildExtensionFromChart(&base, &releases[i], &versions[i]))

			continue
		}

		crd := base.CRD

		if crd != nil {
			if len(resMaps[i].Resources()) == 0 {
				return nil, fmt.Errorf("no crd resource found for %q", base.ID)
			}

			ret = append(ret, *uc.buildExtensionFromCRD(&base, resMaps[i], crds, crd.AnnotationVersionKey))

			continue
		}
	}

	return ret, nil
}

func (uc *UseCase) buildExtensionFromChart(base *base, release *release.Release, version *chart.Version) *Extension {
	var (
		status     string
		deployedAt *time.Time
		current    *Manifest
		latest     *Manifest
	)

	if release.Info != nil {
		status = release.Info.Status.String()
		deployedAt = &release.Info.LastDeployed.Time
		current = &Manifest{
			ID:      release.Chart.Metadata.Name,
			Version: release.Chart.Metadata.Version,
		}
	}

	if len(version.URLs) > 0 {
		latest = &Manifest{
			ID:      version.Name,
			Version: version.Version,
		}
	}

	return &Extension{
		Name:        base.Name,
		Description: base.Description,
		Icon:        base.Logo,
		Status:      status,
		DeployedAt:  deployedAt,
		Current:     current,
		Latest:      latest,
	}
}

func (uc *UseCase) buildExtensionFromCRD(base *base, resMap resmap.ResMap, crds []cluster.CustomResourceDefinition, annotationVersionKey string) *Extension {
	var (
		status     string
		deployedAt *time.Time
		current    *Manifest
		latest     *Manifest
	)

	for i := range crds {
		if version, ok := crds[i].Annotations[annotationVersionKey]; ok {
			status = "deployed"
			deployedAt = &crds[i].CreationTimestamp.Time
			current = &Manifest{
				ID:      base.ID,
				Version: version,
			}
			break
		}
	}

	annotations := resMap.Resources()[0].GetAnnotations()

	if version, ok := annotations[annotationVersionKey]; ok {
		latest = &Manifest{
			ID:      base.ID,
			Version: version,
		}
	}

	return &Extension{
		Name:        base.Name,
		Description: base.Description,
		Icon:        base.Logo,
		Status:      status,
		DeployedAt:  deployedAt,
		Current:     current,
		Latest:      latest,
	}
}

func (uc *UseCase) remoteCRDFormat(repoURL, version string) string {
	return fmt.Sprintf("%s?ref=%s", repoURL, version)
}
