package core

import (
	"context"

	"k8s.io/client-go/rest"
)

type KubeVirtVMRepo interface {
	// VirtualMachines
	CreateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *VirtualMachineSpec) (*VirtualMachine, error)
	GetVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachine, error)
	ListVirtualMachines(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachine, error)
	UpdateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string, spec *VirtualMachineSpec) (*VirtualMachine, error)
	DeleteVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	MigrateVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	StartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	RestartVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	StopVirtualMachine(ctx context.Context, config *rest.Config, namespace, name string) error
	CreateVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string, spec *VirtualMachineCloneSpec) (*VirtualMachineClone, error)
	GetVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineClone, error)
	ListVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) ([]VirtualMachineClone, error)
	DeleteVirtualMachineClone(ctx context.Context, config *rest.Config, namespace, name string) error
	CreateVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string, spec *VirtualMachineSnapshotSpec) (*VirtualMachineSnapshot, error)
	GetVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineSnapshot, error)
	ListVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) ([]VirtualMachineSnapshot, error)
	DeleteVirtualMachineSnapshot(ctx context.Context, config *rest.Config, namespace, name string) error
	CreateVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string, spec *VirtualMachineRestoreSpec) (*VirtualMachineRestore, error)
	GetVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineRestore, error)
	ListVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) ([]VirtualMachineRestore, error)
	DeleteVirtualMachineRestore(ctx context.Context, config *rest.Config, namespace, name string) error
	// VirtualMachine Instances
	GetVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineInstance, error)
	ListVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) ([]VirtualMachineInstance, error)
	UpdateVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string, spec *VirtualMachineInstanceSpec) (*VirtualMachineInstance, error)
	DeleteVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
	MigrateVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string, spec *VirtualMachineInstanceMigrationSpec) (*VirtualMachineInstanceMigration, error)
	PauseVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
	UnpauseVirtualMachineInstance(ctx context.Context, config *rest.Config, namespace, name string) error
}

func (uc *KubeVirtUseCase) CreateVirtualMachine(ctx context.Context, uuid, facility string, metadata Metadata, spec VirtualMachineSpec) (*VirtualMachine, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	return uc.kubeVirtVM.CreateVirtualMachine(ctx, config, metadata.Namespace, metadata.Name, &spec)
}

func (uc *KubeVirtUseCase) GetVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) (*VirtualMachine, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtVM.GetVirtualMachine(ctx, config, name, namespace)
}

func (uc *KubeVirtUseCase) ListVirtualMachines(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachine, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtVM.ListVirtualMachines(ctx, config, namespace)
}

func (uc *KubeVirtUseCase) UpdateVirtualMachine(ctx context.Context, uuid, facility, name, namespace, flavorName, networkName, startupScript string, labels, annotations map[string]string, dataVolumes []string, devices []Device) (*VirtualMachine, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	spec := &VirtualMachineSpec{}

	return uc.kubeVirtVM.UpdateVirtualMachine(ctx, config, namespace, name, spec)
}

func (uc *KubeVirtUseCase) DeleteVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtVM.DeleteVirtualMachine(ctx, config, namespace, name)
}

// Virtual Machine Control Operations
func (uc *KubeVirtUseCase) StartVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtVM.StartVirtualMachine(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) StopVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtVM.StopVirtualMachine(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) PauseVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtVM.PauseVirtualMachineInstance(ctx, config, name, namespace)
}

func (uc *KubeVirtUseCase) UnpauseVirtualMachine(ctx context.Context, uuid, facility, name, namespace string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtVM.UnpauseVirtualMachineInstance(ctx, config, namespace, name)
}

// Virtual Machine Advanced Operations
func (uc *KubeVirtUseCase) CloneVirtualMachine(ctx context.Context, uuid, facility, targetName, targetNamespace, sourceName, sourceNamespace string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	spec := &VirtualMachineCloneSpec{}
	_, err = uc.kubeVirtVM.CreateVirtualMachineClone(ctx, config, targetNamespace, targetName, spec)
	return err
}

func (uc *KubeVirtUseCase) SnapshotVirtualMachine(ctx context.Context, uuid, facility, name, namespace, snapshotName string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	spec := &VirtualMachineSnapshotSpec{}
	_, err = uc.kubeVirtVM.CreateVirtualMachineSnapshot(ctx, config, namespace, name, spec)

	return err
}

func (uc *KubeVirtUseCase) RestoreVirtualMachine(ctx context.Context, uuid, facility, name, namespace, snapshotName string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	spec := &VirtualMachineRestoreSpec{}
	_, err = uc.kubeVirtVM.CreateVirtualMachineRestore(ctx, config, namespace, name, spec)

	return err
}

func (uc *KubeVirtUseCase) MigrateVirtualMachine(ctx context.Context, uuid, facility, name, namespace, targetNode string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}

	return uc.kubeVirtVM.MigrateVirtualMachine(ctx, config, namespace, name)
}
