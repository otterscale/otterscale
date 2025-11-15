package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/workload"
)

type deploymentRepo struct {
	kubernetes *Kubernetes
}

func NewDeploymentRepo(kubernetes *Kubernetes) workload.DeploymentRepo {
	return &deploymentRepo{
		kubernetes: kubernetes,
	}
}

var _ workload.DeploymentRepo = (*deploymentRepo)(nil)

func (r *deploymentRepo) List(ctx context.Context, scope, namespace, selector string) ([]workload.Deployment, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.AppsV1().Deployments(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *deploymentRepo) Get(ctx context.Context, scope, namespace, name string) (*workload.Deployment, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.AppsV1().Deployments(namespace).Get(ctx, name, opts)
}

func (r *deploymentRepo) Create(ctx context.Context, scope, namespace string, d *workload.Deployment) (*workload.Deployment, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.AppsV1().Deployments(namespace).Create(ctx, d, opts)
}

func (r *deploymentRepo) Update(ctx context.Context, scope, namespace string, d *workload.Deployment) (*workload.Deployment, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.AppsV1().Deployments(namespace).Update(ctx, d, opts)
}

func (r *deploymentRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.AppsV1().Deployments(namespace).Delete(ctx, name, opts)
}
