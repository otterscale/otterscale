package persistent

import v1 "k8s.io/api/core/v1"

// Volume represents a Kubernetes Volume resource.
type Volume = v1.Volume

type Persistent struct {
	*PersistentVolumeClaim
	*StorageClass
}

type PersistentUseCase struct {
	persistentVolumeClaim PersistentVolumeClaimRepo
	storageClass          StorageClassRepo
}

func NewPersistentUseCase(persistentVolumeClaim PersistentVolumeClaimRepo, storageClass StorageClassRepo) *PersistentUseCase {
	return &PersistentUseCase{
		persistentVolumeClaim: persistentVolumeClaim,
		storageClass:          storageClass,
	}
}
