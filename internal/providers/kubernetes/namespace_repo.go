package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/cluster"
)

type namespaceRepo struct {
	kubernetes *Kubernetes
}

func NewNamespaceRepo(kubernetes *Kubernetes) cluster.NamespaceRepo {
	return &namespaceRepo{
		kubernetes: kubernetes,
	}
}

var _ cluster.NamespaceRepo = (*namespaceRepo)(nil)

func (r *namespaceRepo) List(ctx context.Context, scope, selector string) ([]cluster.Namespace, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.CoreV1().Namespaces().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *namespaceRepo) Get(ctx context.Context, scope, name string) (*cluster.Namespace, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.CoreV1().Namespaces().Get(ctx, name, opts)
}

func (r *namespaceRepo) Create(ctx context.Context, scope string, ns *cluster.Namespace) (*cluster.Namespace, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.CoreV1().Namespaces().Create(ctx, ns, opts)
}

func (r *namespaceRepo) Update(ctx context.Context, scope string, ns *cluster.Namespace) (*cluster.Namespace, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.CoreV1().Namespaces().Update(ctx, ns, opts)
}

func (r *namespaceRepo) Delete(ctx context.Context, scope, name string) error {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.CoreV1().Namespaces().Delete(ctx, name, opts)
}
