package workload

import (
	"context"

	v1 "k8s.io/api/apps/v1"
)

// Deployment represents a Kubernetes Deployment resource.
type Deployment = v1.Deployment

type DeploymentRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]Deployment, error)
	Get(ctx context.Context, scope, namespace, name string) (*Deployment, error)
	Update(ctx context.Context, scope, namespace string, d *Deployment) (*Deployment, error)
	Create(ctx context.Context, scope, namespace string, d *Deployment) (*Deployment, error)
	Delete(ctx context.Context, scope, namespace, name string) error
	PublicAddress() string
}
