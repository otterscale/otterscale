package core

type StorageCephConfig struct {
	FSID    string
	MONHost string
	Key     string
}

type StorageRGWConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

type StorageConfig struct {
	*StorageCephConfig
	*StorageRGWConfig
}

type StorageUseCase struct {
	action   ActionRepo
	facility FacilityRepo
	cluster  CephClusterRepo
	rbd      CephRBDRepo
	fs       CephFSRepo
	rgw      CephRGWRepo
	machine  MachineRepo
}

func NewStorageUseCase(action ActionRepo, facility FacilityRepo, cluster CephClusterRepo, rbd CephRBDRepo, fs CephFSRepo, rgw CephRGWRepo, machine MachineRepo) *StorageUseCase {
	return &StorageUseCase{
		action:   action,
		facility: facility,
		cluster:  cluster,
		rbd:      rbd,
		fs:       fs,
		rgw:      rgw,
		machine:  machine,
	}
}
