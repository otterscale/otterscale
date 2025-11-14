package persistent

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

type (
	// PersistentVolumeClaim represents a Kubernetes PersistentVolumeClaim resource.
	PersistentVolumeClaim = v1.PersistentVolumeClaim

	// PersistentVolumeClaimPhase represents a Kubernetes PersistentVolumeClaimPhase resource.
	PersistentVolumeClaimPhase = v1.PersistentVolumeClaimPhase

	// PersistentVolumeAccessMode represents a Kubernetes PersistentVolumeAccessMode resource.
	PersistentVolumeAccessMode = v1.PersistentVolumeAccessMode
)

type PersistentVolumeClaimRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]PersistentVolumeClaim, error)
	Get(ctx context.Context, scope, namespace, name string) (*PersistentVolumeClaim, error)
	Create(ctx context.Context, scope, namespace string, pvc *PersistentVolumeClaim) (*PersistentVolumeClaim, error)
	Update(ctx context.Context, scope, namespace string, pvc *PersistentVolumeClaim) (*PersistentVolumeClaim, error)
	Patch(ctx context.Context, scope, namespace, name string, data []byte) (*PersistentVolumeClaim, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}
