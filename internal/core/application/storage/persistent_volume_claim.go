package storage

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

// PersistentVolumeClaim represents a Kubernetes PersistentVolumeClaim resource.
type PersistentVolumeClaim = v1.PersistentVolumeClaim

type PersistentVolumeClaimRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]PersistentVolumeClaim, error)
	Get(ctx context.Context, scope, namespace, name string) (*PersistentVolumeClaim, error)
	Create(ctx context.Context, scope, namespace string, pvc *PersistentVolumeClaim) (*PersistentVolumeClaim, error)
	Update(ctx context.Context, scope, namespace string, pvc *PersistentVolumeClaim) (*PersistentVolumeClaim, error)
	Patch(ctx context.Context, scope, namespace, name string, data []byte) (*PersistentVolumeClaim, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}
