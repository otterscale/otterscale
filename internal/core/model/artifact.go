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
	"github.com/otterscale/otterscale/internal/core/versions"
)

const ModelArtifactModelNameAnnotation = "otterscale.com/model-artifact.model-name"

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
	// find chart ref
	version, err := uc.chart.GetVersion(ctx, chart.RepoURL, "model-artifact", versions.ModelArtifact, true)
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
