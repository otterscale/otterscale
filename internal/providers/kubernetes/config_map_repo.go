package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/config"
)

type configMapRepo struct {
	kubernetes *Kubernetes
}

func NewConfigMapRepo(kubernetes *Kubernetes) config.ConfigMapRepo {
	return &configMapRepo{
		kubernetes: kubernetes,
	}
}

var _ config.ConfigMapRepo = (*configMapRepo)(nil)

func (r *configMapRepo) List(ctx context.Context, scope, namespace, selector string) ([]config.ConfigMap, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.CoreV1().ConfigMaps(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *configMapRepo) Get(ctx context.Context, scope, namespace, name string) (*config.ConfigMap, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, opts)
}

func (r *configMapRepo) Create(ctx context.Context, scope, namespace string, cm *config.ConfigMap) (*config.ConfigMap, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.CoreV1().ConfigMaps(namespace).Create(ctx, cm, opts)
}

func (r *configMapRepo) Update(ctx context.Context, scope, namespace string, cm *config.ConfigMap) (*config.ConfigMap, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.CoreV1().ConfigMaps(namespace).Update(ctx, cm, opts)
}

func (r *configMapRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, opts)
}
