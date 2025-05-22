package core

import (
	"context"

	v1 "k8s.io/api/storage/v1"
	"k8s.io/client-go/rest"
)

type StorageClass = v1.StorageClass

type StorageRepo interface {
	ListStorageClasses(ctx context.Context, config *rest.Config) ([]StorageClass, error)
	GetStorageClass(ctx context.Context, config *rest.Config, name string) (*StorageClass, error)
}

type StorageUseCase struct {
	storage StorageRepo
}

func NewStorageUseCase(storage StorageRepo) *StorageUseCase {
	return &StorageUseCase{
		storage: storage,
	}
}

func (uc *StorageUseCase) ListStorageClasses(ctx context.Context, config *rest.Config) ([]StorageClass, error) {
	return uc.storage.ListStorageClasses(ctx, config)
}
