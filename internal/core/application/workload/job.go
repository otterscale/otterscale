package workload

import (
	"context"

	"golang.org/x/sync/errgroup"
	v1 "k8s.io/api/batch/v1"
)

type (
	// Job represents a Kubernetes Job resource.
	Job = v1.Job

	// JobCondition represents a Kubernetes JobCondition resource.
	JobCondition = v1.JobCondition
)

type JobRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]Job, error)
	Get(ctx context.Context, scope, namespace, name string) (*Job, error)
	Create(ctx context.Context, scope, namespace string, j *Job) (*Job, error)
	Update(ctx context.Context, scope, namespace string, j *Job) (*Job, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

func (uc *UseCase) ListJobs(ctx context.Context, scope, namespace string) (jobs []Job, err error) {
	var ret []Job

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.job.List(egctx, scope, namespace, "")
		if err == nil {
			ret = v
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return ret, nil
}
