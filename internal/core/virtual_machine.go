package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"

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
	VirtualMachineDisk                = virtv1.Disk
	VirtualMachineVolume              = virtv1.Volume
	VirtualMachineClone               = clonev1beta1.VirtualMachineClone
	VirtualMachineSnapshot            = snapshotv1beta1.VirtualMachineSnapshot
	VirtualMachineRestore             = snapshotv1beta1.VirtualMachineRestore
	VirtualMachineInstanceType        = instancetypev1beta1.VirtualMachineInstancetype
	VirtualMachineClusterInstanceType = instancetypev1beta1.VirtualMachineClusterInstancetype
	DataVolume                        = cdiv1beta1.DataVolume
)

type VirtualMachineData struct {
	*VirtualMachine
	*VirtualMachineInstance
	Clones    []VirtualMachineClone
	Snapshots []VirtualMachineSnapshot
	Restores  []VirtualMachineRestore
	Services  []Service
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
	UpdateVirtualMachine(ctx context.Context, config *rest.Config, namespace string, virtualMachine *VirtualMachine) (*VirtualMachine, error)
	DeleteVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	StartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	StopVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	RestartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	ListInstances(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineInstance, error)
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
	release      ReleaseRepo
	action       ActionRepo
	facility     FacilityRepo
	machine      MachineRepo
}

func NewVirtualMachineUseCase(kubeVirt KubeVirtRepo, kubeClone KubeVirtCloneRepo, kubeSnapshot KubeVirtSnapshotRepo, kubeCDI KubeCDIRepo, kubeIT KubeInstanceTypeRepo, kubeCore KubeCoreRepo, kubeStorage KubeStorageRepo, release ReleaseRepo, action ActionRepo, facility FacilityRepo, machine MachineRepo) *VirtualMachineUseCase {
	return &VirtualMachineUseCase{
		kubeVirt:     kubeVirt,
		kubeClone:    kubeClone,
		kubeSnapshot: kubeSnapshot,
		kubeCDI:      kubeCDI,
		kubeIT:       kubeIT,
		kubeCore:     kubeCore,
		kubeStorage:  kubeStorage,
		release:      release,
		action:       action,
		facility:     facility,
		machine:      machine,
	}
}

func (uc *VirtualMachineUseCase) CheckInfrastructureStatus(ctx context.Context, uuid, facility string) (int32, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return 0, err
	}
	rel, err := uc.release.Get(config, LLMd, LLMd)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			return kubevirtHealthNotInstalled, nil
		}
		return 0, err
	}
	switch {
	case rel.Info.Status.IsPending():
		return kubevirtHealthPending, nil
	case rel.Info.Status == release.StatusDeployed:
		return kubevirtHealthOK, nil
	case rel.Info.Status == release.StatusFailed:
		return kubevirtHealthFailed, nil
	}
	return 0, nil
}

func (uc *VirtualMachineUseCase) ListVirtualMachines(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachineData, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.fetchVirtualMachineData(ctx, config, namespace, "")
}

func (uc *VirtualMachineUseCase) GetVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) (*VirtualMachineData, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	vmds, err := uc.fetchVirtualMachineData(ctx, config, namespace, name)
	if err != nil {
		return nil, err
	}
	if len(vmds) == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("virtual machine %s/%s not found", namespace, name))
	}
	return &vmds[0], nil
}

func (uc *VirtualMachineUseCase) CreateVirtualMachine(ctx context.Context, uuid, facility, namespace, name, instanceType, bootDataVolume, startupScript string) (*VirtualMachineData, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	vm, err := uc.kubeVirt.CreateVirtualMachine(ctx, config, namespace, name, instanceType, bootDataVolume, startupScript)
	if err != nil {
		return nil, err
	}
	return &VirtualMachineData{
		VirtualMachine: vm,
	}, nil
}

func (uc *VirtualMachineUseCase) DeleteVirtualMachine(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	// Get related services before deleting the virtual machine
	services, err := uc.kubeCore.ListVirtualMachineServices(ctx, config, namespace, name)
	if err != nil {
		return err
	}

	// Get related snapshots before deleting the virtual machine
	snapshots, err := uc.kubeSnapshot.ListVirtualMachineSnapshots(ctx, config, namespace, name)
	if err != nil {
		return err
	}

	// Delete the virtual machine first
	err = uc.kubeVirt.DeleteVirtualMachine(ctx, config, namespace, name)
	if err != nil {
		return err
	}

	// Delete related services
	for i := range services {
		_ = uc.kubeCore.DeleteService(ctx, config, namespace, services[i].Name)
	}

	// Delete related snapshots
	for i := range snapshots {
		_ = uc.kubeSnapshot.DeleteVirtualMachineSnapshot(ctx, config, namespace, snapshots[i].Name)
	}

	return nil
}

func (uc *VirtualMachineUseCase) AttachVirtualMachineDisk(ctx context.Context, uuid, facility, namespace, name, dvName string) (disk *VirtualMachineDisk, volume *VirtualMachineVolume, err error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, nil, err
	}
	vm, err := uc.kubeVirt.GetVirtualMachine(ctx, config, namespace, name)
	if err != nil {
		return nil, nil, err
	}

	// check if disk already attached
	volumes := vm.Spec.Template.Spec.Volumes
	for i := range volumes {
		if volumes[i].Name == dvName {
			return nil, nil, fmt.Errorf("disk %s is already attached to virtual machine %s/%s", dvName, namespace, name)
		}
	}

	// attach disk
	vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, virtv1.Disk{
		Name: dvName,
		DiskDevice: virtv1.DiskDevice{
			Disk: &virtv1.DiskTarget{
				Bus: virtv1.DiskBusVirtio,
			},
		},
	})

	// add volume
	vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, virtv1.Volume{
		Name: dvName,
		VolumeSource: virtv1.VolumeSource{
			DataVolume: &virtv1.DataVolumeSource{
				Name: dvName,
			},
		},
	})

	// update virtual machine
	newVM, err := uc.kubeVirt.UpdateVirtualMachine(ctx, config, namespace, vm)
	if err != nil {
		return nil, nil, err
	}

	// get disk and volume to verify
	disks := newVM.Spec.Template.Spec.Domain.Devices.Disks
	for i := range disks {
		if disks[i].Name == dvName {
			disk = &disks[i]
			break
		}
	}
	volumes = newVM.Spec.Template.Spec.Volumes
	for i := range volumes {
		if volumes[i].Name == dvName {
			volume = &volumes[i]
			break
		}
	}
	return disk, volume, nil
}

func (uc *VirtualMachineUseCase) DetachVirtualMachineDisk(ctx context.Context, uuid, facility, namespace, name, dvName string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	vm, err := uc.kubeVirt.GetVirtualMachine(ctx, config, namespace, name)
	if err != nil {
		return err
	}

	// check if disk is attached
	found := false
	volumes := vm.Spec.Template.Spec.Volumes
	for i := range vm.Spec.Template.Spec.Volumes {
		if volumes[i].Name == dvName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("disk %s is not attached to virtual machine %s/%s", dvName, namespace, name)
	}

	// detach disk
	newDisks := make([]virtv1.Disk, 0, len(vm.Spec.Template.Spec.Domain.Devices.Disks)-1)
	disks := vm.Spec.Template.Spec.Domain.Devices.Disks
	for i := range disks {
		if disks[i].Name != dvName {
			newDisks = append(newDisks, disks[i])
		}
	}
	vm.Spec.Template.Spec.Domain.Devices.Disks = newDisks

	// remove volume
	newVolumes := make([]virtv1.Volume, 0, len(vm.Spec.Template.Spec.Volumes)-1)
	volumes = vm.Spec.Template.Spec.Volumes
	for i := range volumes {
		if volumes[i].Name != dvName {
			newVolumes = append(newVolumes, volumes[i])
		}
	}
	vm.Spec.Template.Spec.Volumes = newVolumes

	// update virtual machine
	if _, err := uc.kubeVirt.UpdateVirtualMachine(ctx, config, namespace, vm); err != nil {
		return err
	}
	return nil
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

func (uc *VirtualMachineUseCase) CreateVirtualMachineService(ctx context.Context, uuid, facility, namespace, name, vmName string, ports []ServicePort) (*corev1.Service, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeCore.CreateVirtualMachineService(ctx, config, namespace, name, vmName, ports)
}

func (uc *VirtualMachineUseCase) UpdateVirtualMachineService(ctx context.Context, uuid, facility, namespace, name string, ports []ServicePort) (*corev1.Service, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	svc, err := uc.kubeCore.GetService(ctx, config, namespace, name)
	if err != nil {
		return nil, err
	}
	svc.Spec.Ports = ports
	return uc.kubeCore.UpdateService(ctx, config, namespace, svc)
}

func (uc *VirtualMachineUseCase) DeleteVirtualMachineService(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeCore.DeleteService(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) fetchVirtualMachineData(ctx context.Context, config *rest.Config, namespace, vmName string) ([]VirtualMachineData, error) {
	var (
		virtualMachines []VirtualMachine
		instances       []VirtualMachineInstance
		clones          []VirtualMachineClone
		snapshots       []VirtualMachineSnapshot
		restores        []VirtualMachineRestore
		services        []Service
		machines        []Machine
	)

	eg, egctx := errgroup.WithContext(ctx)

	if vmName == "" {
		eg.Go(func() error {
			v, err := uc.kubeVirt.ListVirtualMachines(egctx, config, namespace)
			if err == nil {
				virtualMachines = v
			}
			return err
		})
		eg.Go(func() error {
			v, err := uc.kubeVirt.ListInstances(egctx, config, namespace)
			if err == nil {
				instances = v
			}
			return err
		})
	} else {
		eg.Go(func() error {
			v, err := uc.kubeVirt.GetVirtualMachine(egctx, config, namespace, vmName)
			if err == nil {
				virtualMachines = []VirtualMachine{*v}
			}
			return err
		})
		eg.Go(func() error {
			v, err := uc.kubeVirt.GetInstance(egctx, config, namespace, vmName)
			if isKeyNotFoundError(err) {
				return nil
			}
			if err == nil {
				instances = []VirtualMachineInstance{*v}
			}
			return err
		})
	}

	eg.Go(func() error {
		v, err := uc.kubeClone.ListVirtualMachineClones(egctx, config, namespace, vmName)
		if err == nil {
			clones = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeSnapshot.ListVirtualMachineSnapshots(egctx, config, namespace, vmName)
		if err == nil {
			snapshots = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeSnapshot.ListVirtualMachineRestores(egctx, config, namespace, vmName)
		if err == nil {
			restores = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.kubeCore.ListVirtualMachineServices(egctx, config, namespace, vmName)
		if err == nil {
			services = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := uc.machine.List(egctx)
		if err == nil {
			machines = v
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return uc.assembleVMData(virtualMachines, instances, clones, snapshots, restores, services, machines), nil
}

func (uc *VirtualMachineUseCase) assembleVMData(virtualMachines []VirtualMachine, instances []VirtualMachineInstance, clones []VirtualMachineClone, snapshots []VirtualMachineSnapshot, restores []VirtualMachineRestore, services []Service, machines []Machine) []VirtualMachineData {
	machineMap := make(map[string]string, len(machines))
	for _, m := range machines {
		machineMap[m.Hostname] = m.SystemID
	}

	result := make([]VirtualMachineData, len(virtualMachines))
	for i := range virtualMachines {
		var instance *VirtualMachineInstance
		var machineID string

		// find matching instance
		for j := range instances {
			if virtualMachines[i].Namespace == instances[j].Namespace && virtualMachines[i].Name == instances[j].Name {
				instance = &instances[j]
				if nodeName := instance.Status.NodeName; nodeName != "" {
					machineID = machineMap[nodeName]
				}
				break
			}
		}

		result[i] = VirtualMachineData{
			VirtualMachine:         &virtualMachines[i],
			VirtualMachineInstance: instance,
			Clones:                 uc.filterByVM(clones, virtualMachines[i].Namespace, virtualMachines[i].Name).([]VirtualMachineClone),
			Snapshots:              uc.filterByVM(snapshots, virtualMachines[i].Namespace, virtualMachines[i].Name).([]VirtualMachineSnapshot),
			Restores:               uc.filterByVM(restores, virtualMachines[i].Namespace, virtualMachines[i].Name).([]VirtualMachineRestore),
			Services:               uc.filterByVM(services, virtualMachines[i].Namespace, virtualMachines[i].Name).([]Service),
			MachineID:              machineID,
		}
	}
	return result
}

func (uc *VirtualMachineUseCase) filterByVM(items any, namespace, vmName string) any {
	switch v := items.(type) {
	case []VirtualMachineClone:
		var result []VirtualMachineClone
		for i := range v {
			if v[i].Namespace == namespace && v[i].Labels[VirtualMachineNameLabel] == vmName {
				result = append(result, v[i])
			}
		}
		return result
	case []VirtualMachineSnapshot:
		var result []VirtualMachineSnapshot
		for i := range v {
			if v[i].Namespace == namespace && v[i].Labels[VirtualMachineNameLabel] == vmName {
				result = append(result, v[i])
			}
		}
		return result
	case []VirtualMachineRestore:
		var result []VirtualMachineRestore
		for i := range v {
			if v[i].Namespace == namespace && v[i].Labels[VirtualMachineNameLabel] == vmName {
				result = append(result, v[i])
			}
		}
		return result
	case []Service:
		var result []Service
		for i := range v {
			if v[i].Namespace == namespace && v[i].Labels[VirtualMachineNameLabel] == vmName {
				result = append(result, v[i])
			}
		}
		return result
	default:
		return items
	}
}
