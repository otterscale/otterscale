package core

import (
	"context"

	"k8s.io/client-go/rest"
	"kubevirt.io/api/instancetype/v1beta1"
)

type (
	VirtualMachineInstanceType        = v1beta1.VirtualMachineInstancetype
	VirtualMachineClusterInstanceType = v1beta1.VirtualMachineClusterInstancetype
)

type InstanceTypeRepo interface {
	ListClusterWide(ctx context.Context, config *rest.Config) ([]VirtualMachineClusterInstanceType, error)
	List(ctx context.Context, config *rest.Config, namespace string) ([]VirtualMachineInstanceType, error)
	Get(ctx context.Context, config *rest.Config, namespace, name string) (*VirtualMachineInstanceType, error)
	Create(ctx context.Context, config *rest.Config, namespace, name string, cpu uint32, memory int64) (*VirtualMachineInstanceType, error)
	Delete(ctx context.Context, config *rest.Config, namespace, name string) error
}

type VirtualMachineUseCase struct {
	instanceType InstanceTypeRepo
	action       ActionRepo
	facility     FacilityRepo
}

func NewVirtualMachineUseCase(instanceType InstanceTypeRepo, action ActionRepo, facility FacilityRepo) *VirtualMachineUseCase {
	return &VirtualMachineUseCase{
		instanceType: instanceType,
		action:       action,
		facility:     facility,
	}
}

func (uc *VirtualMachineUseCase) ListClusterWideInstanceTypes(ctx context.Context, uuid, facility string) ([]VirtualMachineClusterInstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.instanceType.ListClusterWide(ctx, config)
}

func (uc *VirtualMachineUseCase) ListInstanceTypes(ctx context.Context, uuid, facility, namespace string) ([]VirtualMachineInstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.instanceType.List(ctx, config, namespace)
}

func (uc *VirtualMachineUseCase) GetInstanceType(ctx context.Context, uuid, facility, namespace, name string) (*VirtualMachineInstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.instanceType.Get(ctx, config, namespace, name)
}

func (uc *VirtualMachineUseCase) CreateInstanceType(ctx context.Context, uuid, facility, namespace, name string, cpu uint32, memory int64) (*VirtualMachineInstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.instanceType.Create(ctx, config, namespace, name, cpu, memory)
}

func (uc *VirtualMachineUseCase) DeleteInstanceType(ctx context.Context, uuid, facility, namespace, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.instanceType.Delete(ctx, config, namespace, name)
}
