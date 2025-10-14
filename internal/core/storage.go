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
	action      ActionRepo
	facility    FacilityRepo
	cephCluster CephClusterRepo
	cephFS      CephFSRepo
	cephRBD     CephRBDRepo
	cephRGW     CephRGWRepo
	machine     MachineRepo
}

func NewStorageUseCase(action ActionRepo, facility FacilityRepo, cephCluster CephClusterRepo, cephFS CephFSRepo, cephRBD CephRBDRepo, cephRGW CephRGWRepo, machine MachineRepo) *StorageUseCase {
	return &StorageUseCase{
		action:      action,
		facility:    facility,
		cephCluster: cephCluster,
		cephFS:      cephFS,
		cephRBD:     cephRBD,
		cephRGW:     cephRGW,
		machine:     machine,
	}
}
