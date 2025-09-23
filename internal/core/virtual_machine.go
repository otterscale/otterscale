package core

import (
	"context"
	"encoding/json"

	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/rest"
	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type SourceType int64

const (
	SourceTypeBlank SourceType = iota
	SourceTypeHTTP
	SourceTypePVC
)

type (
	VirtualMachineInstanceType        = instancetypev1beta1.VirtualMachineInstancetype
	VirtualMachineClusterInstanceType = instancetypev1beta1.VirtualMachineClusterInstancetype
	DataVolume                        = cdiv1beta1.DataVolume
)

type DataVolumePVC struct {
	*DataVolume
	*Storage
}

type KubeCDIRepo interface {
	ListDataVolumes(ctx context.Context, config *rest.Config, namespace string) ([]DataVolume, error)
	GetDataVolume(ctx context.Context, config *rest.Config, namespace, name string) (*DataVolume, error)
	CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string, srcType SourceType, srcData string, size int64, bootImage bool) (*DataVolume, error)
	DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error
}

type KubeInstanceTypeRepo interface {
	ListClusterWide(ctx context.Context, config *rest.Config) ([]VirtualMachineClusterInstanceType, error)
	List(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineInstanceType, error)
	Get(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineInstanceType, error)
	Create(ctx context.Context, config *rest.Config, namespace, name string, cpu uint32, memory int64) (*VirtualMachineInstanceType, error)
	Delete(ctx context.Context, config *rest.Config, namespace, name string) error
}

type VirtualMachineUseCase struct {
	kubeCDI     KubeCDIRepo
	kubeIT      KubeInstanceTypeRepo
	kubeCore    KubeCoreRepo
	kubeStorage KubeStorageRepo
	action      ActionRepo
	facility    FacilityRepo
}

func NewVirtualMachineUseCase(kubeCDI KubeCDIRepo, kubeIT KubeInstanceTypeRepo, kubeCore KubeCoreRepo, kubeStorage KubeStorageRepo, action ActionRepo, facility FacilityRepo) *VirtualMachineUseCase {
	return &VirtualMachineUseCase{
		kubeCDI:     kubeCDI,
		kubeIT:      kubeIT,
		kubeCore:    kubeCore,
		kubeStorage: kubeStorage,
		action:      action,
		facility:    facility,
	}
}

func (uc *VirtualMachineUseCase) ListDataVolumes(ctx context.Context, uuid, facility, namespace string) ([]DataVolumePVC, error) {
	var (
		dataVolumes            []DataVolume
		persistentVolumeClaims []corev1.PersistentVolumeClaim
		storageClasses         []storagev1.StorageClass
	)

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		v, err := uc.kubeCDI.ListDataVolumes(egctx, config, namespace)
		if err == nil {
			dataVolumes = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeCore.ListPersistentVolumeClaims(egctx, config, namespace)
		if err == nil {
			persistentVolumeClaims = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeStorage.ListStorageClasses(egctx, config)
		if err == nil {
			storageClasses = v
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	storageClassMap := toStorageClassMap(storageClasses)

	ret := make([]DataVolumePVC, len(dataVolumes))
	for i := range ret {
		ret[i] = DataVolumePVC{
			DataVolume: &dataVolumes[i],
		}
		for j := range persistentVolumeClaims {
			if dataVolumes[i].Name != persistentVolumeClaims[j].Name {
				continue
			}
			if dataVolumes[i].Namespace != persistentVolumeClaims[j].Namespace {
				continue
			}
			if name := persistentVolumeClaims[j].Spec.StorageClassName; name != nil {
				if sc, ok := storageClassMap[*name]; ok {
					ret[i].Storage = &Storage{
						PersistentVolumeClaim: &persistentVolumeClaims[j],
						StorageClass:          &sc,
					}
				}
			}
			break
		}
	}
	return ret, nil
}

func (uc *VirtualMachineUseCase) GetDataVolume(ctx context.Context, uuid, facility, namespace, name string) (*DataVolumePVC, error) {
	var (
		dataVolume            *DataVolume
		persistentVolumeClaim *corev1.PersistentVolumeClaim
		storageClass          *storagev1.StorageClass
	)
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		v, err := uc.kubeCDI.GetDataVolume(egctx, config, namespace, name)
		if err == nil {
			dataVolume = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeCore.GetPersistentVolumeClaim(egctx, config, namespace, name)
		if err == nil {
			persistentVolumeClaim = v
			scName := persistentVolumeClaim.Spec.StorageClassName
			if scName != nil {
				s, err := uc.kubeStorage.GetStorageClass(egctx, config, *scName)
				if err == nil {
					storageClass = s
				}
				return err
			}
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return &DataVolumePVC{
		DataVolume: dataVolume,
		Storage: &Storage{
			PersistentVolumeClaim: persistentVolumeClaim,
			StorageClass:          storageClass,
		},
	}, nil
}

func (uc *VirtualMachineUseCase) CreateDataVolume(ctx context.Context, uuid, facility, namespace, name string, srcType SourceType, srcData string, size int64, bootImage bool) (*DataVolumePVC, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	dataVolume, err := uc.kubeCDI.CreateDataVolume(ctx, config, namespace, name, srcType, srcData, size, bootImage)
	if err != nil {
		return nil, err
	}
	return &DataVolumePVC{
		DataVolume: dataVolume,
	}, nil
}

func (uc *VirtualMachineUseCase) DeleteDataVolume(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeCDI.DeleteDataVolume(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) ExtendDataVolume(ctx context.Context, uuid, facility, namespace, name string, newSize int64) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	data, err := json.Marshal([]map[string]any{
		{
			"op":    "replace",
			"path":  "/spec/resources/requests/storage",
			"value": resource.NewQuantity(newSize, resource.BinarySI).String(),
		},
	})
	if err != nil {
		return err
	}
	if _, err := uc.kubeCore.PatchPersistentVolumeClaim(ctx, config, namespace, name, data); err != nil {
		return err
	}
	return nil
}

func (uc *VirtualMachineUseCase) ListClusterWideInstanceTypes(ctx context.Context, uuid, facility string) ([]VirtualMachineClusterInstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeIT.ListClusterWide(ctx, config)
}

func (uc *VirtualMachineUseCase) ListInstanceTypes(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachineInstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeIT.List(ctx, config, namespace)
}

func (uc *VirtualMachineUseCase) GetInstanceType(ctx context.Context, uuid, facility, namespace, name string) (*VirtualMachineInstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeIT.Get(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) CreateInstanceType(ctx context.Context, uuid, facility, namespace, name string, cpu uint32, memory int64) (*VirtualMachineInstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeIT.Create(ctx, config, namespace, name, cpu, memory)
}

func (uc *VirtualMachineUseCase) DeleteInstanceType(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeIT.Delete(ctx, config, namespace, name)
}
