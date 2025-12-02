package extension

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/storage/driver"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resource"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/yaml"

	"github.com/otterscale/otterscale/internal/core/application/cluster"
	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/facility"
	"github.com/otterscale/otterscale/internal/core/registry"
	"github.com/otterscale/otterscale/internal/core/scope"
)

type Extension struct {
	DisplayName string
	Description string
	Icon        string
	Status      string
	DeployedAt  *time.Time
	Current     *Manifest
	Latest      *Manifest
}

type Manifest struct {
	ID      string
	Version string
}

type UseCase struct {
	customResourceDefinition cluster.CustomResourceDefinitionRepo
	facility                 facility.FacilityRepo
	release                  release.ReleaseRepo
	repository               registry.RepositoryRepo
	scope                    scope.ScopeRepo
}

func NewUseCase(customResourceDefinition cluster.CustomResourceDefinitionRepo, facility facility.FacilityRepo, release release.ReleaseRepo, repository registry.RepositoryRepo, scope scope.ScopeRepo, service service.ServiceRepo) *UseCase {
	return &UseCase{
		customResourceDefinition: customResourceDefinition,
		facility:                 facility,
		release:                  release,
		repository:               repository,
		scope:                    scope,
	}
}

func (uc *UseCase) ListExtensions(ctx context.Context, scope string, extType Type) ([]Extension, error) {
	components, err := uc.getComponentsByType(extType)
	if err != nil {
		return nil, err
	}

	return uc.listExtensions(ctx, scope, components)
}

func (uc *UseCase) getComponentsByType(extType Type) ([]component, error) {
	switch extType {
	case TypeMetrics:
		return metricsComponents, nil

	case TypeServiceMesh:
		return serviceMeshComponents, nil

	case TypeRegistry:
		return registryComponents, nil

	case TypeModel:
		return modelComponents, nil

	case TypeInstance:
		return instanceComponents, nil

	case TypeStorage:
		return storageComponents, nil

	default:
		return nil, fmt.Errorf("unknown extension type: %v", extType)
	}
}

func (uc *UseCase) InstallOrUpgradeExtensions(ctx context.Context, scope string, manifests []Manifest) error {
	steps := []func(context.Context, string, []Manifest) error{
		uc.precondition,
		uc.processCRDs,
		uc.processReleases,
		uc.processPostFuncs,
	}

	for _, step := range steps {
		if err := step(ctx, scope, manifests); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) listExtensions(ctx context.Context, scope string, components []component) ([]Extension, error) {
	crds, err := uc.customResourceDefinition.List(ctx, scope, "")
	if err != nil {
		return nil, err
	}

	var extensions []Extension

	for i := range components {
		ext, err := uc.buildExtension(ctx, scope, &components[i], crds)
		if err != nil {
			return nil, err
		}

		extensions = append(extensions, *ext)
	}

	return extensions, nil
}

func (uc *UseCase) buildExtension(ctx context.Context, scope string, comp *component, crds []cluster.CustomResourceDefinition) (*Extension, error) {
	chart := comp.Chart
	if chart != nil {
		return uc.buildExtensionFromChart(ctx, scope, comp, chart)
	}

	crd := comp.CRD
	if crd != nil {
		return uc.buildExtensionFromCRD(comp, crd, crds)
	}

	return nil, fmt.Errorf("component %s has neither chart nor CRD", comp.ID)
}

func (uc *UseCase) buildExtensionFromChart(ctx context.Context, scope string, comp *component, chart *chartComponent) (*Extension, error) {
	ext := &Extension{
		DisplayName: comp.DisplayName,
		Description: comp.Description,
		Icon:        comp.Logo,
		Latest: &Manifest{
			ID:      comp.ID,
			Version: chart.Version,
		},
	}

	rel, err := uc.release.Get(ctx, scope, chart.Namespace, chart.Name)
	if err != nil && !errors.Is(err, driver.ErrReleaseNotFound) {
		return nil, err
	}

	if rel != nil && rel.Info != nil {
		ext.Status = rel.Info.Status.String()
		ext.DeployedAt = &rel.Info.LastDeployed.Time
		ext.Current = &Manifest{
			ID:      comp.ID,
			Version: rel.Chart.Metadata.Version,
		}
	}

	return ext, nil
}

func (uc *UseCase) buildExtensionFromCRD(comp *component, crd *crdComponent, crds []cluster.CustomResourceDefinition) (*Extension, error) {
	ext := &Extension{
		DisplayName: comp.DisplayName,
		Description: comp.Description,
		Icon:        comp.Logo,
		Latest: &Manifest{
			ID:      comp.ID,
			Version: crd.Version,
		},
	}

	for i := range crds {
		if version, ok := crds[i].Annotations[crd.VersionAnnotation]; ok {
			ext.Status = "deployed"
			ext.DeployedAt = &crds[i].CreationTimestamp.Time
			ext.Current = &Manifest{
				ID:      comp.ID,
				Version: strings.TrimPrefix(version, "v"),
			}
			break
		}
	}

	return ext, nil
}

func (uc *UseCase) precondition(ctx context.Context, scope string, manifests []Manifest) error {
	return uc.processManifests(ctx, scope, manifests, func(ctx context.Context, scope string, comp *component) error {
		for _, dep := range comp.Dependencies {
			if err := uc.verifyDependency(ctx, scope, dep); err != nil {
				return err
			}
		}
		return nil
	})
}

func (uc *UseCase) verifyDependency(ctx context.Context, scope, depID string) error {
	depComponent, err := uc.whichComponent(depID)
	if err != nil {
		return err
	}

	if depComponent.Chart != nil {
		_, err := uc.release.Get(ctx, scope, depComponent.Chart.Namespace, depComponent.Chart.Name)
		if err != nil {
			return fmt.Errorf("dependency %s not satisfied: %w", depID, err)
		}
	}

	return nil
}

func (uc *UseCase) processCRDs(ctx context.Context, scope string, manifests []Manifest) error {
	return uc.processManifests(ctx, scope, manifests, func(ctx context.Context, scope string, comp *component) error {
		if comp.CRD != nil {
			return uc.createOrUpdateCRDsFromRef(ctx, scope, comp.CRD)
		}
		return nil
	})
}

func (uc *UseCase) processReleases(ctx context.Context, scope string, manifests []Manifest) error {
	return uc.processManifests(ctx, scope, manifests, func(ctx context.Context, scope string, comp *component) error {
		if comp.Chart != nil {
			return uc.installOrUpgradeRelease(ctx, scope, comp.Chart)
		}
		return nil
	})
}

func (uc *UseCase) processPostFuncs(ctx context.Context, scope string, manifests []Manifest) error {
	return uc.processManifests(ctx, scope, manifests, func(ctx context.Context, scope string, comp *component) error {
		if comp.PostFunc != nil {
			return comp.PostFunc(uc, ctx, scope)
		}
		return nil
	})
}

func (uc *UseCase) processManifests(ctx context.Context, scope string, manifests []Manifest, fn func(context.Context, string, *component) error) error {
	eg, egctx := errgroup.WithContext(ctx)

	for _, manifest := range manifests {
		eg.Go(func() error {
			component, err := uc.whichComponent(manifest.ID)
			if err != nil {
				return err
			}
			return fn(egctx, scope, component)
		})
	}

	return eg.Wait()
}

func (uc *UseCase) createOrUpdateCRDsFromRef(ctx context.Context, scope string, crd *crdComponent) error {
	fSys := filesys.MakeFsOnDisk()
	k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

	m, err := k.Run(fSys, crd.Ref)
	if err != nil {
		return err
	}

	for _, node := range m.Resources() {
		if err := uc.processCRDNode(ctx, scope, node, crd.VersionAnnotation); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) processCRDNode(ctx context.Context, scope string, node *resource.Resource, versionAnnotation string) error {
	data, err := node.AsYAML()
	if err != nil {
		return err
	}

	var def *apiextensionsv1.CustomResourceDefinition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return err
	}

	return uc.createOrUpdateCRD(ctx, scope, def, versionAnnotation)
}

func (uc *UseCase) createOrUpdateCRD(ctx context.Context, scope string, crd *apiextensionsv1.CustomResourceDefinition, versionAnnotation string) error {
	existing, err := uc.customResourceDefinition.Get(ctx, scope, crd.Name)
	if k8serrors.IsNotFound(err) {
		_, err := uc.customResourceDefinition.Create(ctx, scope, crd)
		return err
	}
	if err != nil {
		return err
	}

	if uc.shouldUpdateCRD(existing, crd, versionAnnotation) {
		crd.ObjectMeta = existing.ObjectMeta
		_, err = uc.customResourceDefinition.Update(ctx, scope, crd)
		return err
	}

	return nil
}

func (uc *UseCase) shouldUpdateCRD(existing *cluster.CustomResourceDefinition, new *apiextensionsv1.CustomResourceDefinition, versionAnnotation string) bool {
	currentVersion := existing.Annotations[versionAnnotation]
	newVersion := new.Annotations[versionAnnotation]
	return currentVersion != newVersion
}

func (uc *UseCase) installOrUpgradeRelease(ctx context.Context, scope string, chart *chartComponent) error {
	_, err := uc.release.Get(ctx, scope, chart.Namespace, chart.Name)
	if errors.Is(err, driver.ErrReleaseNotFound) {
		return uc.installRelease(ctx, scope, chart)
	}
	if err != nil {
		return err
	}

	return uc.upgradeRelease(ctx, scope, chart)
}

func (uc *UseCase) installRelease(ctx context.Context, scope string, chart *chartComponent) error {
	labels := map[string]string{
		release.TypeLabel: "extension",
	}

	if err := uc.patchValuesMap(ctx, scope, chart.ValuesMap); err != nil {
		return err
	}

	_, err := uc.release.Install(ctx, scope, chart.Namespace, chart.Name, false, chart.Ref, labels, labels, nil, "", chart.ValuesMap)
	return err
}

func (uc *UseCase) upgradeRelease(ctx context.Context, scope string, chart *chartComponent) error {
	_, err := uc.release.Upgrade(ctx, scope, chart.Namespace, chart.Name, false, chart.Ref, "", chart.ValuesMap, false)
	return err
}

func (uc *UseCase) patchValuesMap(ctx context.Context, scopeName string, valuesMap map[string]string) error {
	scope, err := uc.scope.Get(ctx, scopeName)
	if err != nil {
		return err
	}

	for key, value := range valuesMap {
		switch value {
		case "{{ .Scope }}":
			valuesMap[key] = scope.Name
		case "{{ .Scope.UUID }}":
			valuesMap[key] = scope.UUID
		}
	}

	return nil
}
