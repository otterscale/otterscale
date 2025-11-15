package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/workload"
)

type daemonSetRepo struct {
	kubernetes *Kubernetes
}

func NewDaemonSetRepo(kubernetes *Kubernetes) workload.DaemonSetRepo {
	return &daemonSetRepo{
		kubernetes: kubernetes,
	}
}

var _ workload.DaemonSetRepo = (*daemonSetRepo)(nil)

func (r *daemonSetRepo) List(ctx context.Context, scope, namespace, selector string) ([]workload.DaemonSet, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *daemonSetRepo) Get(ctx context.Context, scope, namespace, name string) (*workload.DaemonSet, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, opts)
}

func (r *daemonSetRepo) Create(ctx context.Context, scope, namespace string, ds *workload.DaemonSet) (*workload.DaemonSet, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.AppsV1().DaemonSets(namespace).Create(ctx, ds, opts)
}

func (r *daemonSetRepo) Update(ctx context.Context, scope, namespace string, ds *workload.DaemonSet) (*workload.DaemonSet, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.AppsV1().DaemonSets(namespace).Update(ctx, ds, opts)
}

func (r *daemonSetRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.AppsV1().DaemonSets(namespace).Delete(ctx, name, opts)
}
