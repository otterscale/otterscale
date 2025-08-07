package core

import (
	"context"

	"k8s.io/client-go/rest"
)

type KubeVirtNetRepo interface {
	CreateNetwork(ctx context.Context, config *rest.Config, namespace, name string, spec *DataVolumeSpec) (*DataVolume, error)
	GetNetwork(ctx context.Context, config *rest.Config, namespace, name string) (*DataVolume, error)
	ListNetworks(ctx context.Context, config *rest.Config, namespace, name string) ([]DataVolume, error)
	DeleteNetwork(ctx context.Context, config *rest.Config, namespace, name string) error
	UpdateNetwork(ctx context.Context, config *rest.Config, namespace, name string) error
}

func (uc *KubeVirtUseCase) CreateNetwork(ctx context.Context, uuid, facility, namespace, name string, network Network) (*Network, error) {
	return uc.kubeCore.CreateService(ctx, config, namespace, name)
}

func (uc *KubeVirtUseCase) GetNetwork(ctx context.Context, uuid, facility, name, namespace string) (*Network, error) {
	return uc.KubeCore.GetService()
}

func (uc *KubeVirtUseCase) ListNetworks(ctx context.Context, uuid, facility, namespace string) ([]Network, error) {
	return uc.kubeCore.ListServices(ctx, config, namespace)
}

func (uc *KubeVirtUseCase) UpdateNetwork(ctx context.Context, uuid, facility, network Network) (*Network, error) {
	return uc.KubeCore.UpdateService()
}

func (uc *KubeVirtUseCase) DeleteNetwork(ctx context.Context, uuid, facility, name, namespace string) error {
	return uc.KubeCore.DeleteService()
}
