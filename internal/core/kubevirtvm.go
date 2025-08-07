package core

import (
	"context"

	"k8s.io/client-go/rest"
)

type KubeVirtVMRepo interface {
	// VirtualMachines
	CreateVirtualMachine(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string, spec *VirtualMachineSpec) (*VirtualMachine, error)
	GetVirtualMachine(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) (*VirtualMachine, error)
	ListVirtualMachines(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) ([]VirtualMachine, error)
	UpdateVirtualMachine(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string, spec *VirtualMachineSpec) (*VirtualMachine, error)
	DeleteVirtualMachine(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
	MigrateVirtualMachine(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
	StartVirtualMachine(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
	RestartVirtualMachine(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
	StopVirtualMachine(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
	CreateVirtualMachineClone(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string, spec *VirtualMachineCloneSpec) (*VirtualMachineClone, error)
	GetVirtualMachineClone(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) (*VirtualMachineClone, error)
	ListVirtualMachineClone(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) ([]VirtualMachineClone, error)
	DeleteVirtualMachineClone(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
	CreateVirtualMachineSnapshot(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string, spec *VirtualMachineSnapshotSpec) (*VirtualMachineSnapshot, error)
	GetVirtualMachineSnapshot(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) (*VirtualMachineSnapshot, error)
	ListVirtualMachineSnapshot(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) ([]VirtualMachineSnapshot, error)
	DeleteVirtualMachineSnapshot(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
	CreateVirtualMachineRestore(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string, spec *VirtualMachineRestoreSpec) (*VirtualMachineRestore, error)
	GetVirtualMachineRestore(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) (*VirtualMachineRestore, error)
	ListVirtualMachineRestore(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) ([]VirtualMachineRestore, error)
	DeleteVirtualMachineRestore(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
	// VirtualMachine Instances
	GetVirtualMachineInstance(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) (*VirtualMachineInstance, error)
	ListVirtualMachineInstance(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) ([]VirtualMachineInstance, error)
	UpdateVirtualMachineInstance(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string, spec *VirtualMachineInstanceSpec) (*VirtualMachineInstance, error)
	DeleteVirtualMachineInstance(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
	MigrateVirtualMachineInstance(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string, spec *VirtualMachineInstanceMigrationSpec) (*VirtualMachineInstanceMigration, error)
	PauseVirtualMachineInstance(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
	UnpauseVirtualMachineInstance(ctx context.Context, config *rest.Config, uuid, facility, namespace, name string) error
}

func (uc *KubeVirtUseCase) CreateVirtualMachine(ctx context.Context, uuid, facility string, metadata Metadata, spec VirtualMachineSpec) (*VirtualMachine, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	return uc.kubeVirtVM.CreateVirtualMachine(ctx, metadata, spec)
}

func (uc *KubeVirtUseCase) GetVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) (*VirtualMachine, error) {
	return uc.kubeVirtVM.GetVirtualMachine(ctx, name, namespace)
}

func (uc *KubeVirtUseCase) ListVirtualMachines(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachine, error) {
	return uc.kubeVirtVM.ListVirtualMachines(ctx, namespace)
}

func (uc *KubeVirtUseCase) UpdateVirtualMachine(ctx context.Context, uuid, facility, name, namespace, flavorName, networkName, startupScript string, labels, annotations map[string]string, dataVolumes []string, devices []Device) (*VirtualMachine, error) {
	return uc.kubeVirtVM.UpdateVirtualMachine(ctx, name, namespace, flavorName, networkName, startupScript, labels, annotations, dataVolumes, devices)
}

func (uc *KubeVirtUseCase) DeleteVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) error {
	return uc.kubeVirtVM.DeleteVirtualMachine(ctx, name, namespace)
}

// Virtual Machine Control Operations
func (uc *KubeVirtUseCase) StartVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) error {
	return uc.kubeVirtVM.StartVirtualMachine(ctx, name, namespace)
}

func (uc *KubeVirtUseCase) StopVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) error {
	return uc.kubeVirtVM.StopVirtualMachine(ctx, name, namespace)
}

func (uc *KubeVirtUseCase) PauseVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) error {
	return uc.kubeVirtVM.PauseVirtualMachineInstance(ctx, config*rest.Config, namespace, name)
}

func (uc *KubeVirtUseCase) UnpauseVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) error {
	return uc.kubeVirtVM.UnpauseVirtualMachine(ctx, name, namespace)
}

// Virtual Machine Advanced Operations
func (uc *KubeVirtUseCase) CloneVirtualMachine(ctx context.Context, uuid, facility, targetName, targetNamespace, sourceName, sourceNamespace string) error {
	return uc.kubeVirtVM.CloneVirtualMachine(ctx, targetName, targetNamespace, sourceName, sourceNamespace)
}

func (uc *KubeVirtUseCase) SnapshotVirtualMachine(ctx context.Context, uuid, facility, name, namespace, snapshotName string) error {
	return uc.kubeVirtVM.SnapshotVirtualMachine(ctx, name, namespace, snapshotName)
}

func (uc *KubeVirtUseCase) RestoreVirtualMachine(ctx context.Context, uuid, facility, name, namespace, snapshotName string) error {
	return uc.kubeVirtVM.RestoreVirtualMachine(ctx, name, namespace, snapshotName)
}

func (uc *KubeVirtUseCase) MigrateVirtualMachine(ctx context.Context, uuid, facility, name, namespace, targetNode string) error {
	return uc.kubeVirtVM.MigrateVirtualMachine(ctx, name, namespace, targetNode)
}
