package kubernetessigs

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/service"
)

type GatewayRepo struct {
	kubernetesSigs *KubernetesSigs
}

func NewGatewayRepo(kubernetesSigs *KubernetesSigs) service.GatewayRepo {
	return &GatewayRepo{
		kubernetesSigs: kubernetesSigs,
	}
}

var _ service.GatewayRepo = (*GatewayRepo)(nil)

func (r *GatewayRepo) List(ctx context.Context, scope, namespace, selector string) ([]service.Gateway, error) {
	clientset, err := r.kubernetesSigs.gaClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.GatewayV1().Gateways(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *GatewayRepo) Get(ctx context.Context, scope, namespace, name string) (*service.Gateway, error) {
	clientset, err := r.kubernetesSigs.gaClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.GatewayV1().Gateways(namespace).Get(ctx, name, opts)
}

func (r *GatewayRepo) Create(ctx context.Context, scope, namespace string, g *service.Gateway) (*service.Gateway, error) {
	clientset, err := r.kubernetesSigs.gaClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.GatewayV1().Gateways(namespace).Create(ctx, g, opts)
}

func (r *GatewayRepo) Update(ctx context.Context, scope, namespace string, g *service.Gateway) (*service.Gateway, error) {
	clientset, err := r.kubernetesSigs.gaClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.GatewayV1().Gateways(namespace).Update(ctx, g, opts)
}

func (r *GatewayRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetesSigs.gaClientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.GatewayV1().Gateways(namespace).Delete(ctx, name, opts)
}
