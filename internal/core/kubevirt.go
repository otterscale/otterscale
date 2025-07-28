package core

import (
	"context"

	"k8s.io/client-go/rest"
	clonev1 "kubevirt.io/api/clone/v1beta1"
	v1 "kubevirt.io/api/core/v1"
	snapshotv1 "kubevirt.io/api/snapshot/v1beta1"
	v1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type (
	VMSpec                          = v1.VirtualMachineSpec
	VirtualMachine                  = v1.VirtualMachine
	VirtualMachineInstance          = v1.VirtualMachineInstance
	VirtualMachineInstanceMigration = v1.VirtualMachineInstanceMigration
	VMIMigrateSpec                  = v1.VirtualMachineInstanceMigrationSpec
	VMISpec                         = v1.VirtualMachineInstanceSpec
	VMCloneSpec                     = clonev1.VirtualMachineCloneSpec
	VirtualMachineClone             = clonev1.VirtualMachineClone
	VMSnapshotSpec                  = snapshotv1.VirtualMachineSnapshotSpec
	VirtualMachineSnapshot          = snapshotv1.VirtualMachineSnapshot
	VMRestoreSpec                   = snapshotv1.VirtualMachineRestoreSpec
	VirtualMachineRestore           = snapshotv1.VirtualMachineRestore
	DVSpec                          = v1beta1.DataVolumeSpec
	DataVolume                      = v1beta1.DataVolume
)

type KubeVirtVMRepo interface {
	// VirtualMachines
	CreateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *VMSpec) (*VirtualMachine, error)
	GetVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachine, error)
	ListVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) ([]VirtualMachine, error)
	UpdateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *VMSpec) (*VirtualMachine, error)
	DeleteVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	MigrateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *VMSpec) error
	StartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *VMSpec) error
	RestartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *VMSpec) error
	StopVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *VMSpec) error
	CreateVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string, spec *VMCloneSpec) (*VirtualMachineClone, error)
	GetVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineClone, error)
	ListVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) ([]VirtualMachineClone, error)
	DeleteVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) error
	CreateVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string, spec *VMSnapshotSpec) (*VirtualMachineSnapshot, error)
	GetVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineSnapshot, error)
	ListVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string, spec *VMSnapshotSpec) ([]VirtualMachineSnapshot, error)
	DeleteVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string, spec *VMSnapshotSpec) error
	CreateVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string, spec *VMRestoreSpec) (*VirtualMachineRestore, error)
	GetVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineRestore, error)
	ListVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string, spec *VMSnapshotSpec) ([]VirtualMachineRestore, error)
	DeleteVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string, spec *VMRestoreSpec) error
	// VirtualMachine Instances
	GetVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineInstance, error)
	ListVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) ([]VirtualMachineInstance, error)
	UpdateVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string, spec *VMISpec) (*VirtualMachineInstance, error)
	DeleteVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
	MigrateVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string, spec *VMIMigrateSpec) (*VirtualMachineInstanceMigration, error)
	PauseVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
	UnpauseVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
}

type KubeVirtDVRepo interface {
	CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string, spec *DVSpec) (*DataVolume, error)
	GetDataVolume(ctx context.Context, config *rest.Config, namespace, name string, spec *DVSpec) (*DataVolume, error)
	ListDataVolume(ctx context.Context, config *rest.Config, namespace, name string) ([]DataVolume, error)
	DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error
}
