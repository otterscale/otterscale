package core

import (
	"context"

	"k8s.io/client-go/rest"
)

type KubeVirtInstancetypeRepo interface {
	CreateInstancetype(ctx context.Context, config *rest.Config, name string, obj *Instancetype) (*Instancetype, error)
	GetInstancetype(ctx context.Context, config *rest.Config, name string) (*Instancetype, error)
	ListInstancetypes(ctx context.Context, config *rest.Config) ([]Instancetype, error)
	DeleteInstancetype(ctx context.Context, config *rest.Config, name string) error
}

// Data Volume Operations
func (uc *KubeVirtUseCase) CreateInstancetype(ctx context.Context, uuid, facility, obj Instancetype) (*Instancetype, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	return uc.kubeVirtDV.CreateDataVolume(ctx, config, obj)
}

func (uc *KubeVirtUseCase) GetInstancetype(ctx context.Context, obj Instancetype) (*Instancetype, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtDV.GetDataVolume(ctx, config, obj)
}

func (uc *KubeVirtUseCase) ListInstancetypes(ctx context.Context, obj Instancetype) ([]Instancetype, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtDV.ListDataVolume(ctx, config)
}

func (uc *KubeVirtUseCase) DeleteInstancetype(ctx context.Context, uuid, fa) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtDV.DeleteDataVolume(ctx, config, name)
}
