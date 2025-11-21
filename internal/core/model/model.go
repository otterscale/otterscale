package model

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/otterscale/otterscale/internal/core/application/chart"
	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/application/workload"
)

const (
	ModelNameAnnotation              = "otterscale.com/model.name"
	ModelArtifactModelNameAnnotation = "otterscale.com/model-artifact.model-name"
)

// Release represents a Helm Release resource.
type Model = release.Release

type Artifact struct {
	Name       string
	Namespace  string
	Phase      persistent.PersistentVolumeClaimPhase
	Size       int64
	VolumeName string
	Modelname  string
	CreatedAt  time.Time
}

type Resource struct {
	VGPU       uint32
	VGPUMemory uint32
}

type UseCase struct {
	inferencePoolRepo InferencePoolRepo

	chart                 chart.ChartRepo
	deployment            workload.DeploymentRepo
	httpRoute             service.HTTPRouteRepo
	persistentVolumeClaim persistent.PersistentVolumeClaimRepo
	release               release.ReleaseRepo
	service               service.ServiceRepo
}

func NewUseCase(inferencePoolRepo InferencePoolRepo, chart chart.ChartRepo, deployment workload.DeploymentRepo, httpRoute service.HTTPRouteRepo, persistentVolumeClaim persistent.PersistentVolumeClaimRepo, release release.ReleaseRepo, service service.ServiceRepo) *UseCase {
	return &UseCase{
		inferencePoolRepo:     inferencePoolRepo,
		chart:                 chart,
		deployment:            deployment,
		httpRoute:             httpRoute,
		persistentVolumeClaim: persistentVolumeClaim,
		release:               release,
		service:               service,
	}
}

func (uc *UseCase) ListModels(ctx context.Context, scope, namespace string) (models []Model, uri string, err error) {
	selector := release.TypeLabel + "=" + "model"

	models, err = uc.release.List(ctx, scope, namespace, selector)
	if err != nil {
		return nil, "", err
	}

	url, err := uc.service.URL(scope)
	if err != nil {
		return nil, "", err
	}

	// TODO: get port from service
	uri = fmt.Sprintf("http://%s:8080", url.Hostname())

	return models, uri, nil
}

func (uc *UseCase) CreateModel(ctx context.Context, scope, namespace, name, modelName string) (*Model, error) {
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

	return uc.release.Install(ctx, scope, namespace, name, false, chartRef, labels, labels, annotations, "", valuesMap)
}

func (uc *UseCase) UpdateModel(ctx context.Context, scope, namespace, name string, requests, limits *Resource) (*Model, error) {
	rel, err := uc.release.Get(ctx, scope, namespace, name)
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

	return uc.release.Upgrade(ctx, scope, namespace, name, false, chartRef, "", valuesMap, true)
}

func (uc *UseCase) DeleteModel(ctx context.Context, scope, namespace, name string) error {
	_, err := uc.release.Uninstall(ctx, scope, namespace, name, false)
	return err
}

func (uc *UseCase) ListModelArtifacts(ctx context.Context, scope, namespace string) ([]Artifact, error) {
	selector := release.TypeLabel + "=" + "model-artifact"

	pvcs, err := uc.persistentVolumeClaim.List(ctx, scope, namespace, selector)
	if err != nil {
		return nil, err
	}

	artifacts := make([]Artifact, len(pvcs))
	for i := range pvcs {
		artifact := uc.toArtifact(&pvcs[i])
		artifacts[i] = *artifact
	}

	return artifacts, nil
}

func (uc *UseCase) CreateModelArtifact(ctx context.Context, scope, namespace, name, modelName string, size int64) (*Artifact, error) {
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
	if _, err := uc.release.Install(ctx, scope, namespace, (name), false, chartRef, labels, labels, annotations, "", valuesMap); err != nil {
		return nil, err
	}

	// get pvc
	pvc, err := uc.persistentVolumeClaim.Get(ctx, scope, namespace, name)
	if err != nil {
		return nil, err
	}

	// convert to model artifact
	return uc.toArtifact(pvc), nil
}

func (uc *UseCase) DeleteModelArtifact(ctx context.Context, scope, namespace, name string) error {
	return uc.persistentVolumeClaim.Delete(ctx, scope, namespace, name)
}

func (uc *UseCase) toArtifact(pvc *persistent.PersistentVolumeClaim) *Artifact {
	size := int64(0)
	capacity, ok := pvc.Status.Capacity[v1.ResourceStorage]
	if ok {
		size = capacity.Value()
	}
	return &Artifact{
		Name:       pvc.Name,
		Namespace:  pvc.Namespace,
		Modelname:  pvc.Annotations[ModelArtifactModelNameAnnotation],
		Phase:      pvc.Status.Phase,
		Size:       size,
		VolumeName: pvc.Spec.VolumeName,
		CreatedAt:  pvc.CreationTimestamp.Time,
	}
}
