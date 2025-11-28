package model

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/versions"
)

type Artifact struct {
	Name       string
	Namespace  string
	Phase      persistent.PersistentVolumeClaimPhase
	Size       int64
	VolumeName string
	Modelname  string
	CreatedAt  time.Time
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
	if err := uc.installModelArtifact(ctx, scope, namespace, name, modelName, size); err != nil {
		return nil, err
	}

	pvc, err := uc.persistentVolumeClaim.Get(ctx, scope, namespace, name)
	if err != nil {
		return nil, err
	}

	return uc.toArtifact(pvc), nil
}

func (uc *UseCase) DeleteModelArtifact(ctx context.Context, scope, namespace, name string) error {
	return uc.persistentVolumeClaim.Delete(ctx, scope, namespace, name)
}

func (uc *UseCase) installModelArtifact(ctx context.Context, scope, namespace, name, modelName string, size int64) error {
	// chart ref
	chartRef := fmt.Sprintf("https://github.com/otterscale/charts/releases/download/model-artifact-%[1]s/model-artifact-%[1]s.tgz", versions.ModelArtifact)

	// labels
	labels := map[string]string{
		release.TypeLabel: "model-artifact",
	}

	// annotations
	annotations := map[string]string{
		ModelNameAnnotation: modelName,
	}

	// values
	valuesMap := map[string]string{
		"model.name": modelName,
		"pvc.name":   name,
		"pvc.size":   strconv.FormatInt(size, 10),
	}

	// install
	_, err := uc.release.Install(ctx, scope, namespace, name, false, chartRef, labels, labels, annotations, "", valuesMap)
	return err
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
		Modelname:  pvc.Annotations[ModelNameAnnotation],
		Phase:      pvc.Status.Phase,
		Size:       size,
		VolumeName: pvc.Spec.VolumeName,
		CreatedAt:  pvc.CreationTimestamp.Time,
	}
}
