package core

import (
	"context"

	"k8s.io/client-go/rest"
)

type CephClusterConfig struct {
	FSID    string
	MONHost string
	Key     string
}

type CephObjectConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

type CephConfig struct {
	*CephClusterConfig
	*CephObjectConfig
}

// KubeSMBRepo defines the interface for SMB CRD operations (data layer only)
type KubeSMBRepo interface {
	ListSMBShares(ctx context.Context, config *rest.Config, namespace string) ([]SMBShare, error)
	GetSMBShare(ctx context.Context, config *rest.Config, namespace, name string) (*SMBShare, error)
	CreateSMBShare(ctx context.Context, config *rest.Config, namespace, name string, sizeBytes uint64, browseable, readOnly, guestOk bool, validUsers, commonConfig, securityConfig, userSecret string) error
	CreateSMBCommonConfig(ctx context.Context, config *rest.Config, namespace, name string, shareConfig *SMBShareConfig) error
	CreateSMBSecurityConfig(ctx context.Context, config *rest.Config, namespace, name, mode, realm, secretName string) error
	UpdateSMBShare(ctx context.Context, config *rest.Config, namespace, name string, sizeBytes uint64, browseable, readOnly, guestOk bool, validUsers string) error
	UpdateSMBCommonConfig(ctx context.Context, config *rest.Config, namespace, name string, shareConfig *SMBShareConfig) error
	UpdateSMBSecurityConfig(ctx context.Context, config *rest.Config, namespace, name, mode, realm, secretName string) error
	DeleteSMBCommonConfig(ctx context.Context, config *rest.Config, namespace, name string) error
	DeleteSMBSecurityConfig(ctx context.Context, config *rest.Config, namespace, name string) error
}

type StorageUseCase struct {
	action      ActionRepo
	facility    FacilityRepo
	cephCluster CephClusterRepo
	cephFS      CephFSRepo
	cephRBD     CephRBDRepo
	cephRGW     CephRGWRepo
	machine     MachineRepo
	kubeCore    KubeCoreRepo
	kubeApps    KubeAppsRepo
	kubeSMB     KubeSMBRepo
}

func NewStorageUseCase(action ActionRepo, facility FacilityRepo, cephCluster CephClusterRepo, cephFS CephFSRepo, cephRBD CephRBDRepo, cephRGW CephRGWRepo, machine MachineRepo, kubeCore KubeCoreRepo, kubeApps KubeAppsRepo, kubeSMB KubeSMBRepo) *StorageUseCase {
	return &StorageUseCase{
		action:      action,
		facility:    facility,
		cephCluster: cephCluster,
		cephFS:      cephFS,
		cephRBD:     cephRBD,
		cephRGW:     cephRGW,
		machine:     machine,
		kubeCore:    kubeCore,
		kubeApps:    kubeApps,
		kubeSMB:     kubeSMB,
	}
}
