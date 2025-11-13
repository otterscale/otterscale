package storage

import v1 "k8s.io/api/core/v1"

// Volume represents a Kubernetes Volume resource.
type Volume = v1.Volume

type Storage struct {
	*PersistentVolumeClaim
	*StorageClass
}

type StorageUseCase struct {
	persistentVolumeClaim PersistentVolumeClaimRepo
	storageClass          StorageClassRepo
}

func NewStorageUseCase(persistentVolumeClaim PersistentVolumeClaimRepo, storageClass StorageClassRepo) *StorageUseCase {
	return &StorageUseCase{
		persistentVolumeClaim: persistentVolumeClaim,
		storageClass:          storageClass,
	}
}
