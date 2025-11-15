package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/workload"
)

type jobRepo struct {
	kubernetes *Kubernetes
}

func NewJobRepo(kubernetes *Kubernetes) workload.JobRepo {
	return &jobRepo{
		kubernetes: kubernetes,
	}
}

var _ workload.JobRepo = (*jobRepo)(nil)

func (r *jobRepo) List(ctx context.Context, scope, namespace, selector string) ([]workload.Job, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.BatchV1().Jobs(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *jobRepo) Get(ctx context.Context, scope, namespace, name string) (*workload.Job, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.BatchV1().Jobs(namespace).Get(ctx, name, opts)
}

func (r *jobRepo) Create(ctx context.Context, scope, namespace string, j *workload.Job) (*workload.Job, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.BatchV1().Jobs(namespace).Create(ctx, j, opts)
}

func (r *jobRepo) Update(ctx context.Context, scope, namespace string, j *workload.Job) (*workload.Job, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.BatchV1().Jobs(namespace).Update(ctx, j, opts)
}

func (r *jobRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.BatchV1().Jobs(namespace).Delete(ctx, name, opts)
}
