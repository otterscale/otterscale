package core

import (
	"context"

	"k8s.io/client-go/rest"
)

type KubeVirtDVRepo interface {
	CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string, spec *DataVolumeSpec) (*DataVolume, error)
	GetDataVolume(ctx context.Context, config *rest.Config, namespace, name string) (*DataVolume, error)
	ListDataVolume(ctx context.Context, config *rest.Config, namespace string) ([]DataVolume, error)
	DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error
}

// Data Volume Operations
func (uc *KubeVirtUseCase) CreateDataVolume(ctx context.Context, uuid, facility, namespace, name string, dataVolume DataVolume) (*DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	return uc.KubeVirtDV.CreateDataVolume(ctx, config, namespace, name, &dataVolume.Spec)
}

func (uc *KubeVirtUseCase) GetDataVolume(ctx context.Context, uuid, facility, name, namespace string) (*DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.KubeVirtDV.GetDataVolume(ctx, config, name, namespace)
}

func (uc *KubeVirtUseCase) ListDataVolumes(ctx context.Context, uuid, facility, namespace string) ([]DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.KubeVirtDV.ListDataVolume(ctx, config, namespace)
}

func (uc *KubeVirtUseCase) DeleteDataVolume(ctx context.Context, uuid, facility, name, namespace string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.KubeVirtDV.DeleteDataVolume(ctx, config, name, namespace)
}
