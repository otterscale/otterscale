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
	"sigs.k8s.io/kustomize/kyaml/filesys"

	"github.com/otterscale/otterscale/internal/core/application/chart"
	"github.com/otterscale/otterscale/internal/core/application/cluster"
	"github.com/otterscale/otterscale/internal/core/application/release"
)

type Manifest struct {
	Name    string
	Version string
}

type Extension struct {
	DisplayName string
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

func (uc *UseCase) ListExtensions(ctx context.Context, scope string, extType Type) ([]Extension, error) {
	switch extType {
	case TypeGeneral:
		return uc.listExtensions(ctx, scope, general)

	case TypeRegistry:
		return uc.listExtensions(ctx, scope, registry)

	case TypeModel:
		return uc.listExtensions(ctx, scope, model)

	case TypeInstance:
		return uc.listExtensions(ctx, scope, instance)

	case TypeStorage:
		return uc.listExtensions(ctx, scope, storage)

	default:
		return nil, fmt.Errorf("unknown extension type: %v", extType)
	}
}

func (uc *UseCase) InstallExtensions(ctx context.Context, scope string, manifests []Manifest) error {
	eg, egctx := errgroup.WithContext(ctx)

	for _, manifest := range manifests {
		eg.Go(func() error {
			base, err := uc.base(manifest.Name)
			if err != nil {
				return err
			}

			for i := range base.Charts {
				chart := base.Charts[i]

				if _, err := uc.release.Install(egctx, scope, base.Namespace, base.Name, false, chart.Ref, chart.Labels, chart.Labels, chart.Annotations, "", chart.ValuesMap); err != nil {
					return err
				}

				if chart.PostFunc != nil {
					if err := chart.PostFunc(scope); err != nil {
						return err
					}
				}
			}

			crd := base.CRD

			if crd != nil {
				fSys := filesys.MakeFsOnDisk()
				k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

				m, err := k.Run(fSys, crd.Ref)
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
			base, err := uc.base(manifest.Name)
			if err != nil {
				return err
			}

			for i := range base.Charts {
				chart := base.Charts[i]

				if _, err := uc.release.Upgrade(egctx, scope, base.Namespace, base.Name, false, chart.Ref, "", chart.ValuesMap, true); err != nil {
					return err
				}
			}

			crd := base.CRD

			if crd != nil {
				fSys := filesys.MakeFsOnDisk()
				k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

				m, err := k.Run(fSys, crd.Ref)
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
	releases := make([]release.Release, len(bases))
	crds := []cluster.CustomResourceDefinition{}

	eg, egctx := errgroup.WithContext(ctx)

	for i := range bases {
		base := bases[i]

		if len(base.Charts) > 0 {
			eg.Go(func() error {
				v, err := uc.release.Get(egctx, scope, base.Namespace, base.Name)
				if err == nil {
					releases[i] = *v
				}
				if errors.Is(err, driver.ErrReleaseNotFound) {
					return nil
				}
				return err
			})
		}
	}

	eg.Go(func() error {
		v, err := uc.customResourceDefinition.List(egctx, scope, "")
		if err == nil {
			crds = v
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	ret := []Extension{}

	for i := range bases {
		base := bases[i]

		if len(base.Charts) > 0 {
			ret = append(ret, *uc.buildExtensionFromChart(&base, &releases[i]))
		}

		crd := base.CRD

		if crd != nil {
			ret = append(ret, *uc.buildExtensionFromCRD(&base, crds, crd.AnnotationVersionKey))
		}
	}

	return ret, nil
}

func (uc *UseCase) buildExtensionFromChart(base *base, release *release.Release) *Extension {
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
			Name:    release.Chart.Metadata.Name,
			Version: release.Chart.Metadata.Version,
		}
	}

	charts := base.Charts

	if len(charts) > 0 {
		chart := charts[0]

		latest = &Manifest{
			Name:    base.Name,
			Version: chart.Version,
		}
	}

	return &Extension{
		DisplayName: base.DisplayName,
		Description: base.Description,
		Icon:        base.Logo,
		Status:      status,
		DeployedAt:  deployedAt,
		Current:     current,
		Latest:      latest,
	}
}

func (uc *UseCase) buildExtensionFromCRD(base *base, crds []cluster.CustomResourceDefinition, annotationVersionKey string) *Extension {
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
				Name:    base.Name,
				Version: version,
			}
			break
		}
	}

	crd := base.CRD

	if crd != nil {
		latest = &Manifest{
			Name:    base.Name,
			Version: crd.Version,
		}
	}

	return &Extension{
		DisplayName: base.DisplayName,
		Description: base.Description,
		Icon:        base.Logo,
		Status:      status,
		DeployedAt:  deployedAt,
		Current:     current,
		Latest:      latest,
	}
}
