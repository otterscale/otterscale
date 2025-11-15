package workload

import (
	"context"

	v1 "k8s.io/api/apps/v1"
)

// StatefulSet represents a Kubernetes StatefulSet resource.
type StatefulSet = v1.StatefulSet

type StatefulSetRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]StatefulSet, error)
	Get(ctx context.Context, scope, namespace, name string) (*StatefulSet, error)
	Create(ctx context.Context, scope, namespace string, ss *StatefulSet) (*StatefulSet, error)
	Update(ctx context.Context, scope, namespace string, ss *StatefulSet) (*StatefulSet, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}
