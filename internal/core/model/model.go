package model

import (
	"context"
	"fmt"
	"strconv"

	"github.com/otterscale/otterscale/internal/core/application/chart"
	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/application/storage"
	"github.com/otterscale/otterscale/internal/core/application/workload"
)

const (
	ModelNameAnnotation              = "otterscale.com/model.name"
	ModelArtifactModelNameAnnotation = "otterscale.com/model-artifact.model-name"
)

// Release represents a Helm Release resource.
type Model = release.Release

type ModelResource struct {
	VGPU       uint32
	VGPUMemory uint32
}

type ModelArtifact struct{}

type ModelUseCase struct {
	chart                 chart.ChartRepo
	deployment            workload.DeploymentRepo
	release               release.ReleaseRepo
	persistentVolumeClaim storage.PersistentVolumeClaimRepo
}

func NewModelUseCase(chart chart.ChartRepo, deployment workload.DeploymentRepo, release release.ReleaseRepo, persistentVolumeClaim storage.PersistentVolumeClaimRepo) *ModelUseCase {
	return &ModelUseCase{
		chart:                 chart,
		deployment:            deployment,
		release:               release,
		persistentVolumeClaim: persistentVolumeClaim,
	}
}

func (uc *ModelUseCase) ListModels(ctx context.Context, scope, namespace string) ([]Model, error) {
	selector := fmt.Sprintf("%s=%s", release.TypeLabel, "model")

	return uc.release.List(scope, namespace, selector)
}

func (uc *ModelUseCase) CreateModel(ctx context.Context, scope, namespace, name, modelName string) (*Model, error) {
	// find chart ref
	version, err := uc.chart.GetStableVersion(ctx, chart.RepoURL, "llm-d-modelservice", true)
	if err != nil {
		return nil, err
	}

	// check URLs
	if len(version.URLs) == 0 {
		return nil, fmt.Errorf("no URLs found for chart model-artifact")
	}

	chartRef := version.URLs[0]

	// labels
	labels := map[string]string{
		release.TypeLabel: "model",
	}

	// annotations
	annotations := map[string]string{
		ModelNameAnnotation: modelName,
	}

	// values
	valuesMap := map[string]string{} // TODO: waiting for v0.3.0

	return uc.release.Install(scope, namespace, name, false, chartRef, labels, labels, annotations, "", valuesMap)
}

func (uc *ModelUseCase) UpdateModel(ctx context.Context, scope, namespace, name string, requests, limits *ModelResource) (*Model, error) {
	rel, err := uc.release.Get(scope, namespace, name)
	if err != nil {
		return nil, err
	}

	chart := rel.Chart
	if chart == nil {
		return nil, fmt.Errorf("chart not found in release %s", name)
	}

	chartRef := ""
	if v, ok := rel.Config[release.ChartRefKey]; ok {
		if str, ok := v.(string); ok {
			chartRef = str
		}
	}

	if chartRef == "" {
		return nil, fmt.Errorf("chart ref not found in release %s", name)
	}

	// values
	_, _ = requests, limits          // TODO: waiting for v0.3.0
	valuesMap := map[string]string{} // TODO: waiting for v0.3.0

	return uc.release.Upgrade(scope, namespace, name, false, chartRef, "", valuesMap, true)
}

func (uc *ModelUseCase) DeleteModel(ctx context.Context, scope, namespace, name string) error {
	_, err := uc.release.Uninstall(scope, namespace, name, false)
	return err
}

func (uc *ModelUseCase) ListModelArtifacts(ctx context.Context, scope, namespace string) ([]ModelArtifact, error) {
	selector := fmt.Sprintf("%s=%s", release.TypeLabel, "model-artifact")

	pvcs, err := uc.persistentVolumeClaim.List(ctx, scope, namespace, selector)
	if err != nil {
		return nil, err
	}

	artifacts := make([]ModelArtifact, len(pvcs))
	for i := range pvcs {
		artifact := uc.toModelArtifact(&pvcs[i])
		artifacts[i] = *artifact
	}

	return artifacts, nil
}

func (uc *ModelUseCase) CreateModelArtifact(ctx context.Context, scope, namespace, name, modelName string, size int64) (*ModelArtifact, error) {
	// find chart ref
	version, err := uc.chart.GetStableVersion(ctx, chart.RepoURL, "model-artifact", true)
	if err != nil {
		return nil, err
	}

	// check URLs
	if len(version.URLs) == 0 {
		return nil, fmt.Errorf("no URLs found for chart model-artifact")
	}

	chartRef := version.URLs[0]

	// labels
	labels := map[string]string{
		release.TypeLabel: "model-artifact",
	}

	// annotations
	annotations := map[string]string{
		ModelArtifactModelNameAnnotation: modelName,
	}

	// values
	valuesMap := map[string]string{
		"model.name": modelName,
		"pvc.name":   name,
		"pvc.size":   strconv.FormatInt(size, 10),
	}

	// install
	if _, err := uc.release.Install(scope, namespace, (name), false, chartRef, labels, labels, annotations, "", valuesMap); err != nil {
		return nil, err
	}

	// get pvc
	pvc, err := uc.persistentVolumeClaim.Get(ctx, scope, namespace, name)
	if err != nil {
		return nil, err
	}

	// convert to model artifact
	return uc.toModelArtifact(pvc), nil
}

func (uc *ModelUseCase) DeleteModelArtifact(ctx context.Context, scope, namespace, name string) error {
	return uc.persistentVolumeClaim.Delete(ctx, scope, namespace, name)
}

func (uc *ModelUseCase) toModelArtifact(pvc *storage.PersistentVolumeClaim) *ModelArtifact {
	// size := int64(0)
	// capacity, ok := pvc.Status.Capacity[corev1.ResourceStorage]
	// if ok {
	// 	size = capacity.Value()
	// }
	return &ModelArtifact{
		// Name:       pvc.Name,
		// Namespace:  pvc.Namespace,
		// Modelname:  pvc.Annotations[ModelArtifactModelNameAnnotation],
		// Phase:      pvc.Status.Phase,
		// Size:       size,
		// VolumeName: pvc.Spec.VolumeName,
		// CreatedAt:  pvc.CreationTimestamp.Time,
	}
}
