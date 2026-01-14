package vault

import (
	"context"
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/otterscale/otterscale/internal/core/secret"
)

type kubeConfigRepo struct {
	secretRepo secret.Repository
}

func NewKubeConfigRepo(secretRepo secret.Repository) secret.KubeConfigRepo {
	return &kubeConfigRepo{
		secretRepo: secretRepo,
	}
}

var _ secret.KubeConfigRepo = (*kubeConfigRepo)(nil)

func (r *kubeConfigRepo) GetKubeConfig(ctx context.Context, scope string) (*api.Config, error) {
	s, err := r.secretRepo.Get(ctx, scope, secret.SecretTypeKubeConfig, secret.KubeConfigName)
	if err != nil {
		return nil, err
	}

	data, ok := s.Data[secret.KubeConfigDataKey]
	if !ok {
		return nil, fmt.Errorf("kubeconfig data not found in secret")
	}

	return clientcmd.Load(data)
}

func (r *kubeConfigRepo) StoreKubeConfig(ctx context.Context, scope string, config *api.Config) error {
	data, err := clientcmd.Write(*config)
	if err != nil {
		return fmt.Errorf("failed to serialize kubeconfig: %w", err)
	}

	s := &secret.Secret{
		Scope: scope,
		Type:  secret.SecretTypeKubeConfig,
		Name:  secret.KubeConfigName,
		Data: map[string][]byte{
			secret.KubeConfigDataKey: data,
		},
		Metadata: map[string]string{
			"managed-by": "otterscale",
		},
	}

	return r.secretRepo.Store(ctx, s)
}

func (r *kubeConfigRepo) DeleteKubeConfig(ctx context.Context, scope string) error {
	return r.secretRepo.Delete(ctx, scope, secret.SecretTypeKubeConfig, secret.KubeConfigName)
}

func (r *kubeConfigRepo) ExistsKubeConfig(ctx context.Context, scope string) (bool, error) {
	return r.secretRepo.Exists(ctx, scope, secret.SecretTypeKubeConfig, secret.KubeConfigName)
}
