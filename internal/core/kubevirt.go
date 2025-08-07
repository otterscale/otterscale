package core

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	clonev1 "kubevirt.io/api/clone/v1beta1"
	v1 "kubevirt.io/api/core/v1"
	snapshotv1 "kubevirt.io/api/snapshot/v1beta1"
	v1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type (
	VirtualMachineSpec                  = v1.VirtualMachineSpec
	VirtualMachine                      = v1.VirtualMachine
	VirtualMachineInstance              = v1.VirtualMachineInstance
	VirtualMachineInstanceMigration     = v1.VirtualMachineInstanceMigration
	VirtualMachineInstanceMigrationSpec = v1.VirtualMachineInstanceMigrationSpec
	VirtualMachineInstanceSpec          = v1.VirtualMachineInstanceSpec
	VirtualMachineCloneSpec             = clonev1.VirtualMachineCloneSpec
	VirtualMachineClone                 = clonev1.VirtualMachineClone
	VirtualMachineSnapshotSpec          = snapshotv1.VirtualMachineSnapshotSpec
	VirtualMachineSnapshot              = snapshotv1.VirtualMachineSnapshot
	VirtualMachineRestoreSpec           = snapshotv1.VirtualMachineRestoreSpec
	VirtualMachineRestore               = snapshotv1.VirtualMachineRestore
	DataVolumeSpec                      = v1beta1.DataVolumeSpec
	DataVolume                          = v1beta1.DataVolume
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
	FlavorName    string
	NetworkName   string
	StartupScript string
	DataVolumes   []string
	Devices       []Device
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

// DataVolume represents a data volume resource
type KubeVirtDataVolume struct {
	Metadata  Metadata
	Source    string
	Type      string
	SizeBytes int64
}

// Network represents a network resource
type KubeVirtNetwork struct {
	Metadata      Metadata
	ServiceType   string
	Port          int32
	NodePort      int32
	ContainerPort int32
}

// Flavor represents a flavor resource
type Flavor struct {
	Metadata    Metadata
	CpuCores    float32
	MemoryBytes int64
}

type KubeVirtUseCase struct {
	kubeCore    KubeCoreRepo
	KubeApps    KubeAppsRepo
	kubeVirtVM  KubeVirtVMRepo
	KubeVirtDV  KubeVirtDVRepo
	KubeVirtNet KubeVirtNetRepo
	action      ActionRepo
	facility    FacilityRepo
}
