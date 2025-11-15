package workload

import (
	"context"

	v1 "k8s.io/api/batch/v1"
)

// Job represents a Kubernetes Job resource.
type Job = v1.Job

type JobRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]Job, error)
	Get(ctx context.Context, scope, namespace, name string) (*Job, error)
	Create(ctx context.Context, scope, namespace string, j *Job) (*Job, error)
	Update(ctx context.Context, scope, namespace string, j *Job) (*Job, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}
