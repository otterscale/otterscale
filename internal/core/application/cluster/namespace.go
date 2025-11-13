package cluster

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

// Namespace represents a Kubernetes Namespace resource.
type Namespace = v1.Namespace

type NamespaceRepo interface {
	List(ctx context.Context, scope, selector string) ([]Namespace, error)
	Get(ctx context.Context, scope, name string) (*Namespace, error)
	Create(ctx context.Context, scope string, ns *Namespace) (*Namespace, error)
	Update(ctx context.Context, scope string, ns *Namespace) (*Namespace, error)
	Delete(ctx context.Context, scope, name string) error
}

func (uc *ClusterUseCase) ListNamespaces(ctx context.Context, scope string) ([]Namespace, error) {
	return uc.namespace.List(ctx, scope, "")
}
