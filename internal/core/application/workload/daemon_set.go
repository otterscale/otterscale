package workload

import (
	"context"

	v1 "k8s.io/api/apps/v1"
)

// DaemonSet represents a Kubernetes DaemonSet resource.
type DaemonSet = v1.DaemonSet

type DaemonSetRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]DaemonSet, error)
	Get(ctx context.Context, scope, namespace, name string) (*DaemonSet, error)
	Create(ctx context.Context, scope, namespace string, ds *DaemonSet) (*DaemonSet, error)
	Update(ctx context.Context, scope, namespace string, ds *DaemonSet) (*DaemonSet, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}
