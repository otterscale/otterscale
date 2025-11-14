package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/config"
)

type secretRepo struct {
	kubernetes *Kubernetes
}

func NewSecretRepo(kubernetes *Kubernetes) config.SecretRepo {
	return &secretRepo{
		kubernetes: kubernetes,
	}
}

var _ config.SecretRepo = (*secretRepo)(nil)

func (r *secretRepo) List(ctx context.Context, scope, namespace, selector string) ([]config.Secret, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.CoreV1().Secrets(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *secretRepo) Get(ctx context.Context, scope, namespace, name string) (*config.Secret, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.CoreV1().Secrets(namespace).Get(ctx, name, opts)
}

func (r *secretRepo) Create(ctx context.Context, scope, namespace string, s *config.Secret) (*config.Secret, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.CoreV1().Secrets(namespace).Create(ctx, s, opts)
}

func (r *secretRepo) Update(ctx context.Context, scope, namespace string, s *config.Secret) (*config.Secret, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.CoreV1().Secrets(namespace).Update(ctx, s, opts)
}

func (r *secretRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.CoreV1().Secrets(namespace).Delete(ctx, name, opts)
}
