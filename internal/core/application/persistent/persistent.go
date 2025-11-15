package persistent

import v1 "k8s.io/api/core/v1"

// Volume represents a Kubernetes Volume resource.
type Volume = v1.Volume

type Persistent struct {
	*PersistentVolumeClaim
	*StorageClass
}

type UseCase struct {
	persistentVolumeClaim PersistentVolumeClaimRepo
	storageClass          StorageClassRepo
}

func NewUseCase(persistentVolumeClaim PersistentVolumeClaimRepo, storageClass StorageClassRepo) *UseCase {
	return &UseCase{
		persistentVolumeClaim: persistentVolumeClaim,
		storageClass:          storageClass,
	}
}
