package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/service"
)

type serviceRepo struct {
	kubernetes *Kubernetes
}

func NewServiceRepo(kubernetes *Kubernetes) service.ServiceRepo {
	return &serviceRepo{
		kubernetes: kubernetes,
	}
}

var _ service.ServiceRepo = (*serviceRepo)(nil)

func (r *serviceRepo) List(ctx context.Context, scope, namespace, selector string) ([]service.Service, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.CoreV1().Services(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *serviceRepo) Get(ctx context.Context, scope, namespace, name string) (*service.Service, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.CoreV1().Services(namespace).Get(ctx, name, opts)
}

func (r *serviceRepo) Create(ctx context.Context, scope, namespace string, s *service.Service) (*service.Service, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.CoreV1().Services(namespace).Create(ctx, s, opts)
}

func (r *serviceRepo) Update(ctx context.Context, scope, namespace string, s *service.Service) (*service.Service, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.CoreV1().Services(namespace).Update(ctx, s, opts)
}

func (r *serviceRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.CoreV1().Services(namespace).Delete(ctx, name, opts)
}
