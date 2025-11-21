package kubernetessigs

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/model"
)

type modelRepo struct {
	kubernetesSigs *KubernetesSigs
}

func NewInferencePoolRepo(kubernetesSigs *KubernetesSigs) model.InferencePoolRepo {
	return &modelRepo{
		kubernetesSigs: kubernetesSigs,
	}
}

var _ model.InferencePoolRepo = (*modelRepo)(nil)

func (r *modelRepo) List(ctx context.Context, scope, namespace, selector string) ([]model.InferencePool, error) {
	clientset, err := r.kubernetesSigs.gaieClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.InferenceV1().InferencePools(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *modelRepo) Get(ctx context.Context, scope, namespace, name string) (*model.InferencePool, error) {
	clientset, err := r.kubernetesSigs.gaieClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.InferenceV1().InferencePools(namespace).Get(ctx, name, opts)
}

func (r *modelRepo) Create(ctx context.Context, scope, namespace string, ip *model.InferencePool) (*model.InferencePool, error) {
	clientset, err := r.kubernetesSigs.gaieClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.InferenceV1().InferencePools(namespace).Create(ctx, ip, opts)
}

func (r *modelRepo) Update(ctx context.Context, scope, namespace string, ip *model.InferencePool) (*model.InferencePool, error) {
	clientset, err := r.kubernetesSigs.gaieClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.InferenceV1().InferencePools(namespace).Update(ctx, ip, opts)
}

func (r *modelRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetesSigs.gaieClientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.InferenceV1().InferencePools(namespace).Delete(ctx, name, opts)
}
