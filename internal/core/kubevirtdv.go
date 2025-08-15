package core

import (
	"context"

	"k8s.io/client-go/rest"
)

type KubeVirtDVRepo interface {
	CreateDataVolume(ctx context.Context, config *rest.Config, namespace, name string, source_type string, source string, sizeBytes int64) (*DataVolume, error)
	GetDataVolume(ctx context.Context, config *rest.Config, namespace, name string) (*DataVolume, error)
	ListDataVolume(ctx context.Context, config *rest.Config, namespace string) ([]DataVolume, error)
	DeleteDataVolume(ctx context.Context, config *rest.Config, namespace, name string) error
}

// Data Volume Operations
func (uc *KubeVirtUseCase) CreateDataVolume(ctx context.Context, uuid, facility, namespace, name string, source_type string, source string, sizeBytes int64) (*DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	return uc.kubeVirtDV.CreateDataVolume(ctx, config, namespace, name, source_type, source, sizeBytes)
}

func (uc *KubeVirtUseCase) GetDataVolume(ctx context.Context, uuid, facility, name, namespace string) (*DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtDV.GetDataVolume(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) ListDataVolumes(ctx context.Context, uuid, facility, namespace string) ([]DataVolume, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.kubeVirtDV.ListDataVolume(ctx, config, namespace)
}

func (uc *KubeVirtUseCase) DeleteDataVolume(ctx context.Context, uuid, facility, name, namespace string) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.kubeVirtDV.DeleteDataVolume(ctx, config, namespace, name)
}
