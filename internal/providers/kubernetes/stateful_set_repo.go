package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/workload"
)

type statefulSetRepo struct {
	kubernetes *Kubernetes
}

func NewStatefulSetRepo(kubernetes *Kubernetes) workload.StatefulSetRepo {
	return &statefulSetRepo{
		kubernetes: kubernetes,
	}
}

var _ workload.StatefulSetRepo = (*statefulSetRepo)(nil)

func (r *statefulSetRepo) List(ctx context.Context, scope, namespace, selector string) ([]workload.StatefulSet, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *statefulSetRepo) Get(ctx context.Context, scope, namespace, name string) (*workload.StatefulSet, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.AppsV1().StatefulSets(namespace).Get(ctx, name, opts)
}

func (r *statefulSetRepo) Create(ctx context.Context, scope, namespace string, ss *workload.StatefulSet) (*workload.StatefulSet, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.AppsV1().StatefulSets(namespace).Create(ctx, ss, opts)
}

func (r *statefulSetRepo) Update(ctx context.Context, scope, namespace string, ss *workload.StatefulSet) (*workload.StatefulSet, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.AppsV1().StatefulSets(namespace).Update(ctx, ss, opts)
}

func (r *statefulSetRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.AppsV1().StatefulSets(namespace).Delete(ctx, name, opts)
}
