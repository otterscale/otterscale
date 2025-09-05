package core

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	clonev1 "kubevirt.io/api/clone/v1beta1"
	virtCorev1 "kubevirt.io/api/core/v1"
	snapshotv1 "kubevirt.io/api/snapshot/v1beta1"
	v1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type (
	VirtualMachineSpec                  = virtCorev1.VirtualMachineSpec
	VirtualMachine                      = virtCorev1.VirtualMachine
	VirtualMachineInstance              = virtCorev1.VirtualMachineInstance
	VirtualMachineInstanceMigration     = virtCorev1.VirtualMachineInstanceMigration
	VirtualMachineInstanceMigrationSpec = virtCorev1.VirtualMachineInstanceMigrationSpec
	VirtualMachineInstanceSpec          = virtCorev1.VirtualMachineInstanceSpec
	VirtualMachineCloneSpec             = clonev1.VirtualMachineCloneSpec
	VirtualMachineClone                 = clonev1.VirtualMachineClone
	VirtualMachineSnapshotSpec          = snapshotv1.VirtualMachineSnapshotSpec
	VirtualMachineSnapshot              = snapshotv1.VirtualMachineSnapshot
	VirtualMachineRestoreSpec           = snapshotv1.VirtualMachineRestoreSpec
	VirtualMachineRestore               = snapshotv1.VirtualMachineRestore
	DataVolumeSpec                      = v1beta1.DataVolumeSpec
	DataVolume                          = v1beta1.DataVolume
	KubeVirtVolume                      = virtCorev1.Volume
	KubeVirtDisk                        = virtCorev1.Disk
	KubeVirtDevice                      = virtCorev1.Devices
	VirtualMachineService               = corev1.Service
	VirtualMachineServiceSpec           = corev1.ServiceSpec
	VirtualMachineServicePort           = corev1.ServicePort
)

type Metadata struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
	CreatedAt   *timestamppb.Timestamp
	UpdatedAt   *timestamppb.Timestamp
}

// VirtualMachineSpec defines the specification for a virtual machine
type KubeVirtVirtualMachineSpec struct {
	InstanceTypeName string
	NetworkName      string
	StartupScript    string
	DataVolumes      []string
	Devices          []Device
}

// VirtualMachine represents a virtual machine resource
type KubeVirtVirtualMachine struct {
	Metadata  Metadata
	Spec      KubeVirtVirtualMachineSpec
	Status    string // Maps to pb.VirtualMachine_Status enum (e.g., "RUNNING", "STOPPED")
	Snapshots []Operation
	Clones    []Operation
	Migrates  []Operation
	Restores  []Operation
}

type VirtualMachineResources struct {
	InstanceName string
	CPUcores     uint32
	MemoryBytes  int64
}

type VirtualMachineOperation struct {
	Name        string
	Type        string
	Description string
	CreatedAt   *time.Time
	Status      VirtualMachineOperationStatus
}

// VirtualMachineOperationStatus represents the status of an operation
type VirtualMachineOperationStatus struct {
	Status  int32
	Message string
	Reason  string
}

type DiskDevice struct {
	Name     string
	DiskType string
	Bus      string
	Data     string
}

// Operation represents an operation on a virtual machine (snapshot, clone, etc.)
type Operation struct {
	Name        string
	Type        string
	Description string
	CreatedAt   *timestamppb.Timestamp
	Status      OperationResult
}

// OperationResult represents the result of an operation
type OperationResult struct {
	Status  string // Maps to pb.VirtualMachine_Operation_Result_Status enum (e.g., "IN_PROGRESS", "SUCCEEDED")
	Message string
	Reason  string
}

// Device represents a device attached to a virtual machine
type Device struct {
	Name string
	Type string
}

// InstanceType represents a flavor resource
type InstanceType struct {
	Metadata    Metadata
	CPUCores    uint32
	MemoryBytes int64
}

type KubeVirtUseCase struct {
	kubeCore             KubeCoreRepo
	kubeApps             KubeAppsRepo
	kubeVirtVM           KubeVirtVMRepo
	kubeVirtDV           KubeVirtDVRepo
	kubeVirtInstanceType KubeVirtInstanceTypeRepo
	action               ActionRepo
	facility             FacilityRepo
}

func NewKubeVirtUseCase(kubeCore KubeCoreRepo, kubeApps KubeAppsRepo, kubeVirtVM KubeVirtVMRepo, kubeVirtDV KubeVirtDVRepo, kubeVirtInstanceType KubeVirtInstanceTypeRepo, action ActionRepo, facility FacilityRepo) *KubeVirtUseCase {
	return &KubeVirtUseCase{
		kubeCore:             kubeCore,
		kubeApps:             kubeApps,
		kubeVirtVM:           kubeVirtVM,
		kubeVirtDV:           kubeVirtDV,
		kubeVirtInstanceType: kubeVirtInstanceType,
		action:               action,
		facility:             facility,
	}
}
