package core

const (
	cephConfigCommand        = "ceph config generate-minimal-conf && ceph auth get client.admin"
	cephRGWUserListCommand   = "radosgw-admin user list"
	cephRGWUserCreateCommand = "radosgw-admin user create --system --uid=otterscale --display-name=OtterScale --format json"
	cephRGWUserInfoCommand   = "radosgw-admin user info --uid=otterscale --format=json"
)

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
