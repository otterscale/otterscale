package core

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
)

type Model = release.Release

type ModelResource struct {
	VGPU       uint32
	VGPUMemory uint32
}

type ModelArtifact struct {
	Name       string
	Namespace  string
	Phase      corev1.PersistentVolumeClaimPhase
	Size       int64
	VolumeName string
	Modelname  string
	CreatedAt  time.Time
}

type ModelUseCase struct {
	action   ActionRepo
	chart    ChartRepo
	facility FacilityRepo
	kubeCore KubeCoreRepo
	release  ReleaseRepo
}

func NewModelUseCase(action ActionRepo, chart ChartRepo, facility FacilityRepo, kubeCore KubeCoreRepo, release ReleaseRepo) *ModelUseCase {
	return &ModelUseCase{
		action:   action,
		chart:    chart,
		facility: facility,
		kubeCore: kubeCore,
		release:  release,
	}
}

func (uc *ModelUseCase) ListModels(ctx context.Context, scope, facility, namespace string) ([]Release, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	selector := fmt.Sprintf("%s=%s", TypeLabel, "model")
	return uc.release.List(config, namespace, selector)
}

func (uc *ModelUseCase) CreateModel(ctx context.Context, scope, facility, namespace, name, modelName string) (*Release, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	// find chart ref
	charts, err := uc.chart.List(ctx, chartRepoURL, true)
	if err != nil {
		return nil, err
	}
	chartRef, err := findLatestChartRef(charts, "llm-d-modelservice")
	if err != nil {
		return nil, err
	}

	// labels
	labels := map[string]string{
		TypeLabel: "model",
	}

	// annotations
	annotations := map[string]string{
		ModelNameAnnotation: modelName,
	}

	// values
	values, err := toReleaseValues("", map[string]string{ // TODO: waiting for v0.3.0
		ChartRefKey: chartRef,
	})
	if err != nil {
		return nil, err
	}

	return uc.release.Install(config, namespace, getReleaseName(name), false, chartRef, labels, labels, annotations, values)
}

func (uc *ModelUseCase) UpdateModel(ctx context.Context, scope, facility, namespace, name string, requests, limits *ModelResource) (*Release, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	release, err := uc.release.Get(config, namespace, name)
	if err != nil {
		return nil, err
	}

	chart := release.Chart
	if chart == nil {
		return nil, fmt.Errorf("chart not found in release %s", name)
	}

	chartRef := ""
	if v, ok := release.Config[ChartRefKey]; ok {
		if str, ok := v.(string); ok {
			chartRef = str
		}
	}
	if chartRef == "" {
		return nil, fmt.Errorf("chart ref not found in release %s", name)
	}

	// values
	_, _ = requests, limits                                 // TODO: waiting for v0.3.0
	values, err := toReleaseValues("", map[string]string{}) // TODO: waiting for v0.3.0
	if err != nil {
		return nil, err
	}

	return uc.release.Upgrade(config, namespace, getReleaseName(name), false, chartRef, values, true)
}

func (uc *ModelUseCase) DeleteModel(ctx context.Context, scope, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return err
	}
	if _, err := uc.release.Uninstall(config, namespace, name, false); err != nil {
		return err
	}
	return nil
}

func (uc *ModelUseCase) ListModelArtifacts(ctx context.Context, scope, facility, namespace string) ([]ModelArtifact, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	selector := fmt.Sprintf("%s=%s", TypeLabel, "model-artifact")
	pvcs, err := uc.kubeCore.ListPersistentVolumeClaimsByLabel(ctx, config, namespace, selector)
	if err != nil {
		return nil, err
	}

	artifacts := make([]ModelArtifact, len(pvcs))
	for i := range pvcs {
		artifact := toModelArtifact(&pvcs[i])
		artifacts[i] = *artifact
	}
	return artifacts, nil
}

func (uc *ModelUseCase) CreateModelArtifact(ctx context.Context, scope, facility, namespace, name, modelName string, size int64) (*ModelArtifact, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	// find chart ref
	charts, err := uc.chart.List(ctx, chartRepoURL, true)
	if err != nil {
		return nil, err
	}
	chartRef, err := findLatestChartRef(charts, "model-artifact")
	if err != nil {
		return nil, err
	}

	// labels
	labels := map[string]string{
		TypeLabel: "model-artifact",
	}

	// annotations
	annotations := map[string]string{
		ModelArtifactModelNameAnnotation: modelName,
	}

	// values
	values, err := toReleaseValues("", map[string]string{
		"model.name": modelName,
		"pvc.name":   name,
		"pvc.size":   strconv.FormatInt(size, 10),
	})
	if err != nil {
		return nil, err
	}

	// install
	if _, err := uc.release.Install(config, namespace, getReleaseName(name), false, chartRef, labels, labels, annotations, values); err != nil {
		return nil, err
	}

	// get pvc
	pvc, err := uc.kubeCore.GetPersistentVolumeClaim(ctx, config, namespace, name)
	if err != nil {
		return nil, err
	}

	// convert to model artifact
	return toModelArtifact(pvc), nil
}

func (uc *ModelUseCase) DeleteModelArtifact(ctx context.Context, scope, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return err
	}
	return uc.kubeCore.DeletePersistentVolumeClaim(ctx, config, namespace, name)
}

func toModelArtifact(pvc *corev1.PersistentVolumeClaim) *ModelArtifact {
	size := int64(0)
	capacity, ok := pvc.Status.Capacity[corev1.ResourceStorage]
	if ok {
		size = capacity.Value()
	}
	return &ModelArtifact{
		Name:       pvc.Name,
		Namespace:  pvc.Namespace,
		Modelname:  pvc.Annotations[ModelArtifactModelNameAnnotation],
		Phase:      pvc.Status.Phase,
		Size:       size,
		VolumeName: pvc.Spec.VolumeName,
		CreatedAt:  pvc.CreationTimestamp.Time,
	}
}
