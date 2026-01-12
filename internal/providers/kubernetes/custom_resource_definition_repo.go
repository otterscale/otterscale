package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/cluster"
)

type customResourceDefinitionRepo struct {
	kubernetes *Kubernetes
}

func NewCustomResourceDefinitionRepo(kubernetes *Kubernetes) cluster.CustomResourceDefinitionRepo {
	return &customResourceDefinitionRepo{
		kubernetes: kubernetes,
	}
}

var _ cluster.CustomResourceDefinitionRepo = (*customResourceDefinitionRepo)(nil)

func (r *customResourceDefinitionRepo) List(ctx context.Context, scope, selector string) ([]cluster.CustomResourceDefinition, error) {
	clientset, err := r.kubernetes.apiClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *customResourceDefinitionRepo) Get(ctx context.Context, scope, name string) (*cluster.CustomResourceDefinition, error) {
	clientset, err := r.kubernetes.apiClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, name, opts)
}

func (r *customResourceDefinitionRepo) Create(ctx context.Context, scope string, crd *cluster.CustomResourceDefinition) (*cluster.CustomResourceDefinition, error) {
	clientset, err := r.kubernetes.apiClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, crd, opts)
}

func (r *customResourceDefinitionRepo) Update(ctx context.Context, scope string, crd *cluster.CustomResourceDefinition) (*cluster.CustomResourceDefinition, error) {
	clientset, err := r.kubernetes.apiClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.ApiextensionsV1().CustomResourceDefinitions().Update(ctx, crd, opts)
}

func (r *customResourceDefinitionRepo) Delete(ctx context.Context, scope, name string) error {
	clientset, err := r.kubernetes.apiClientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, name, opts)
}
