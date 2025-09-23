package core

import (
	"context"
	"encoding/json"

	"golang.org/x/sync/errgroup"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/rest"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	clonev1beta1 "kubevirt.io/api/clone/v1beta1"
	virtv1 "kubevirt.io/api/core/v1"
	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"
	snapshotv1beta1 "kubevirt.io/api/snapshot/v1beta1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type SourceType int64

const (
	SourceTypeBlank SourceType = iota
	SourceTypeHTTP
	SourceTypePVC
)

type (
	VirtualMachine                    = virtv1.VirtualMachine
	VirtualMachineInstance            = virtv1.VirtualMachineInstance
	VirtualMachineInstanceMigration   = virtv1.VirtualMachineInstanceMigration
	VirtualMachineClone               = clonev1beta1.VirtualMachineClone
	VirtualMachineSnapshot            = snapshotv1beta1.VirtualMachineSnapshot
	VirtualMachineRestore             = snapshotv1beta1.VirtualMachineRestore
	VirtualMachineInstanceType        = instancetypev1beta1.VirtualMachineInstancetype
	VirtualMachineClusterInstanceType = instancetypev1beta1.VirtualMachineClusterInstancetype
	DataVolume                        = cdiv1beta1.DataVolume
)

type VirtualMachineDetails struct {
	*VirtualMachine
	*VirtualMachineInstance
	Clones    []VirtualMachineClone
	Snapshots []VirtualMachineSnapshot
	Restores  []VirtualMachineRestore
	MachineID string
}

type DataVolumeWithStorage struct {
	*DataVolume
	*Storage
}

type KubeVirtRepo interface {
	ListVirtualMachines(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachine, error)
	GetVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachine, error)
	CreateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name, instanceType, bootDataVolume, startupScript string) (*VirtualMachine, error)
	DeleteVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	StartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	StopVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	RestartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	GetInstance(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineInstance, error)
	PauseInstance(ctx context.Context, config *rest.Config, namespace, name string) error
	UnpauseInstance(ctx context.Context, config *rest.Config, namespace, name string) error
	MigrateInstance(ctx context.Context, config *rest.Config, namespace, name, hostname string) (*VirtualMachineInstanceMigration, error)
}

type KubeVirtCloneRepo interface {
	ListVirtualMachineClones(ctx context.Context, config *rest.Config, namespace, vmName string) ([]VirtualMachineClone, error)
	CreateVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name, source, target string) (*VirtualMachineClone, error)
	DeleteVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) error
}

type KubeVirtSnapshotRepo interface {
	ListVirtualMachineSnapshots(ctx context.Context, config *rest.Config, namespace, vmName string) ([]VirtualMachineSnapshot, error)
	CreateVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name, vmName string) (*VirtualMachineSnapshot, error)
	DeleteVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) error
	ListVirtualMachineRestores(ctx context.Context, config *rest.Config, namespace, vmName string) ([]VirtualMachineRestore, error)
	CreateVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name, vmName, snapshot string) (*VirtualMachineRestore, error)
	DeleteVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) error
}

type KubeCDIRepo interface {
	ListDataVolumes(ctx context.Context, config *rest.Config, namespace string, bootImage bool) ([]DataVolume, error)
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
	kubeVirt     KubeVirtRepo
	kubeClone    KubeVirtCloneRepo
	kubeSnapshot KubeVirtSnapshotRepo
	kubeCDI      KubeCDIRepo
	kubeIT       KubeInstanceTypeRepo
	kubeCore     KubeCoreRepo
	kubeStorage  KubeStorageRepo
	action       ActionRepo
	facility     FacilityRepo
	machine      MachineRepo
}

func NewVirtualMachineUseCase(kubeVirt KubeVirtRepo, kubeClone KubeVirtCloneRepo, kubeSnapshot KubeVirtSnapshotRepo, kubeCDI KubeCDIRepo, kubeIT KubeInstanceTypeRepo, kubeCore KubeCoreRepo, kubeStorage KubeStorageRepo, action ActionRepo, facility FacilityRepo, machine MachineRepo) *VirtualMachineUseCase {
	return &VirtualMachineUseCase{
		kubeVirt:     kubeVirt,
		kubeClone:    kubeClone,
		kubeSnapshot: kubeSnapshot,
		kubeCDI:      kubeCDI,
		kubeIT:       kubeIT,
		kubeCore:     kubeCore,
		kubeStorage:  kubeStorage,
		action:       action,
		facility:     facility,
		machine:      machine,
	}
}

func (uc *VirtualMachineUseCase) ListVirtualMachines(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachineDetails, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	vms, err := uc.kubeVirt.ListVirtualMachines(ctx, config, namespace)
	if err != nil {
		return nil, err
	}
	// FIX: Waited before sending request
	eg, egctx := errgroup.WithContext(ctx)
	details := make([]VirtualMachineDetails, len(vms))
	for i := range vms {
		eg.Go(func() error {
			vm, err := uc.GetVirtualMachine(egctx, uuid, facility, vms[i].Namespace, vms[i].Name)
			if err == nil {
				details[i] = *vm
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return details, nil
}

func (uc *VirtualMachineUseCase) GetVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) (*VirtualMachineDetails, error) {
	var (
		virtualMachine *VirtualMachine
		instance       *VirtualMachineInstance
		clones         []VirtualMachineClone
		snapshots      []VirtualMachineSnapshot
		restores       []VirtualMachineRestore
		machineID      string
	)

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		v, err := uc.kubeVirt.GetVirtualMachine(egctx, config, namespace, name)
		if err == nil {
			virtualMachine = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeVirt.GetInstance(egctx, config, namespace, name)
		if isKeyNotFoundError(err) {
			return nil
		}
		if err == nil {
			instance = v
			ms, err := uc.machine.List(egctx)
			if err == nil {
				for i := range ms {
					if ms[i].Hostname == instance.Status.NodeName {
						machineID = ms[i].SystemID
						break
					}
				}
			}
			return err
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeClone.ListVirtualMachineClones(egctx, config, namespace, name)
		if err == nil {
			clones = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeSnapshot.ListVirtualMachineSnapshots(egctx, config, namespace, name)
		if err == nil {
			snapshots = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeSnapshot.ListVirtualMachineRestores(egctx, config, namespace, name)
		if err == nil {
			restores = v
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return &VirtualMachineDetails{
		VirtualMachine:         virtualMachine,
		VirtualMachineInstance: instance,
		Clones:                 clones,
		Snapshots:              snapshots,
		Restores:               restores,
		MachineID:              machineID,
	}, nil
}

func (uc *VirtualMachineUseCase) CreateVirtualMachine(ctx context.Context, uuid, facility, namespace, name, instanceType, bootDataVolume, startupScript string) (*VirtualMachineDetails, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	vm, err := uc.kubeVirt.CreateVirtualMachine(ctx, config, namespace, name, instanceType, bootDataVolume, startupScript)
	if err != nil {
		return nil, err
	}
	return &VirtualMachineDetails{
		VirtualMachine: vm,
	}, nil
}

func (uc *VirtualMachineUseCase) DeleteVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirt.DeleteVirtualMachine(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) CreateVirtualMachineClone(ctx context.Context, uuid, facility, namespace, name, source, target string) (*VirtualMachineClone, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeClone.CreateVirtualMachineClone(ctx, config, namespace, name, source, target)
}

func (uc *VirtualMachineUseCase) DeleteVirtualMachineClone(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeClone.DeleteVirtualMachineClone(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) CreateVirtualMachineSnapshot(ctx context.Context, uuid, facility, namespace, name, vmName string) (*VirtualMachineSnapshot, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeSnapshot.CreateVirtualMachineSnapshot(ctx, config, namespace, name, vmName)
}

func (uc *VirtualMachineUseCase) DeleteVirtualMachineSnapshot(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeSnapshot.DeleteVirtualMachineSnapshot(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) CreateVirtualMachineRestore(ctx context.Context, uuid, facility, namespace, name, vmName, snapshot string) (*VirtualMachineRestore, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeSnapshot.CreateVirtualMachineRestore(ctx, config, namespace, name, vmName, snapshot)
}

func (uc *VirtualMachineUseCase) DeleteVirtualMachineRestore(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeSnapshot.DeleteVirtualMachineRestore(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) StartVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirt.StartVirtualMachine(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) StopVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirt.StopVirtualMachine(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) RestartVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirt.RestartVirtualMachine(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) PauseInstance(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirt.PauseInstance(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) ResumeInstance(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirt.UnpauseInstance(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) MigrateInstance(ctx context.Context, uuid, facility, namespace, name, hostname string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	if _, err := uc.kubeVirt.MigrateInstance(ctx, config, namespace, name, hostname); err != nil {
		return err
	}
	return nil
}

func (uc *VirtualMachineUseCase) ListDataVolumes(ctx context.Context, uuid, facility, namespace string, bootImage bool) ([]DataVolumeWithStorage, error) {
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
		v, err := uc.kubeCDI.ListDataVolumes(egctx, config, namespace, bootImage)
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

	ret := make([]DataVolumeWithStorage, len(dataVolumes))
	for i := range ret {
		ret[i] = DataVolumeWithStorage{
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

func (uc *VirtualMachineUseCase) GetDataVolume(ctx context.Context, uuid, facility, namespace, name string) (*DataVolumeWithStorage, error) {
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
	return &DataVolumeWithStorage{
		DataVolume: dataVolume,
		Storage: &Storage{
			PersistentVolumeClaim: persistentVolumeClaim,
			StorageClass:          storageClass,
		},
	}, nil
}

func (uc *VirtualMachineUseCase) CreateDataVolume(ctx context.Context, uuid, facility, namespace, name string, srcType SourceType, srcData string, size int64, bootImage bool) (*DataVolumeWithStorage, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	dataVolume, err := uc.kubeCDI.CreateDataVolume(ctx, config, namespace, name, srcType, srcData, size, bootImage)
	if err != nil {
		return nil, err
	}
	return &DataVolumeWithStorage{
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
