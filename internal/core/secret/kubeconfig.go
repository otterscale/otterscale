package secret

import (
	"context"

	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	KubeConfigDataKey = "kubeconfig"
	KubeConfigName    = "cluster-kubeconfig"
)

// KubeConfigRepo defines the interface for kubeconfig storage operations.
type KubeConfigRepo interface {
	GetKubeConfig(ctx context.Context, scope string) (*api.Config, error)
	StoreKubeConfig(ctx context.Context, scope string, config *api.Config) error
	DeleteKubeConfig(ctx context.Context, scope string) error
	ExistsKubeConfig(ctx context.Context, scope string) (bool, error)
}
