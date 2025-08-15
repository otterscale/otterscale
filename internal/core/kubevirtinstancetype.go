package core

import (
	"context"

	"k8s.io/client-go/rest"
)

type KubeVirtInstancetypeRepo interface {
	CreateInstancetype(ctx context.Context, config *rest.Config, instancetype Instancetype) (*Instancetype, error)
	GetInstancetype(ctx context.Context, config *rest.Config, name string) (*Instancetype, error)
	ListInstancetypes(ctx context.Context, config *rest.Config) ([]Instancetype, error)
	DeleteInstancetype(ctx context.Context, config *rest.Config, name string) error
}

// Instancetype Operations
func (uc *KubeVirtUseCase) CreateInstancetype(ctx context.Context, uuid, facility string, instancetype Instancetype) (*Instancetype, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtInstancetype.CreateInstancetype(ctx, config, instancetype)
}

func (uc *KubeVirtUseCase) GetInstancetype(ctx context.Context, uuid, facility, name string) (*Instancetype, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtInstancetype.GetInstancetype(ctx, config, name)
}

func (uc *KubeVirtUseCase) ListInstancetypes(ctx context.Context, uuid, facility string) ([]Instancetype, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtInstancetype.ListInstancetypes(ctx, config)
}

func (uc *KubeVirtUseCase) DeleteInstancetype(ctx context.Context, uuid, facility, name string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtInstancetype.DeleteInstancetype(ctx, config, name)
}
