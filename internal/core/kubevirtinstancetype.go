package core

import (
	"context"

	"k8s.io/client-go/rest"
)

type KubeVirtInstanceTypeRepo interface {
	CreateInstanceType(ctx context.Context, config *rest.Config, InstanceType InstanceType) (*InstanceType, error)
	GetInstanceType(ctx context.Context, config *rest.Config, name string) (*InstanceType, error)
	ListInstanceTypes(ctx context.Context, config *rest.Config) ([]InstanceType, error)
	DeleteInstanceType(ctx context.Context, config *rest.Config, name string) error
}

// InstanceType Operations
func (uc *KubeVirtUseCase) CreateInstanceType(ctx context.Context, uuid, facility string, InstanceType InstanceType) (*InstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtInstanceType.CreateInstanceType(ctx, config, InstanceType)
}

func (uc *KubeVirtUseCase) GetInstanceType(ctx context.Context, uuid, facility, name string) (*InstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtInstanceType.GetInstanceType(ctx, config, name)
}

func (uc *KubeVirtUseCase) ListInstanceTypes(ctx context.Context, uuid, facility string) ([]InstanceType, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtInstanceType.ListInstanceTypes(ctx, config)
}

func (uc *KubeVirtUseCase) DeleteInstanceType(ctx context.Context, uuid, facility, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtInstanceType.DeleteInstanceType(ctx, config, name)
}
