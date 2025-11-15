package config

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

// Secret represents a Kubernetes Secret resource.
type Secret = v1.Secret

type SecretRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]Secret, error)
	Get(ctx context.Context, scope, namespace, name string) (*Secret, error)
	Create(ctx context.Context, scope, namespace string, s *Secret) (*Secret, error)
	Update(ctx context.Context, scope, namespace string, s *Secret) (*Secret, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}
