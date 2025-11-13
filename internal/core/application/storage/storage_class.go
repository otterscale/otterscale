package storage

import (
	"context"

	v1 "k8s.io/api/storage/v1"
)

// StorageClass represents a Kubernetes StorageClass resource.
type StorageClass = v1.StorageClass

type StorageClassRepo interface {
	List(ctx context.Context, scope, selector string) ([]StorageClass, error)
	Get(ctx context.Context, scope, name string) (*StorageClass, error)
}

func (uc *StorageUseCase) ListStorageClasses(ctx context.Context, scope string) ([]StorageClass, error) {
	return uc.storageClass.List(ctx, scope, "")
}
