package kubernetessigs

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/service"
)

type httpRouteRepo struct {
	kubernetesSigs *KubernetesSigs
}

func NewHTTPRouteRepo(kubernetesSigs *KubernetesSigs) service.HTTPRouteRepo {
	return &httpRouteRepo{
		kubernetesSigs: kubernetesSigs,
	}
}

var _ service.HTTPRouteRepo = (*httpRouteRepo)(nil)

func (r *httpRouteRepo) List(ctx context.Context, scope, namespace, selector string) ([]service.HTTPRoute, error) {
	clientset, err := r.kubernetesSigs.gaClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.GatewayV1().HTTPRoutes(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *httpRouteRepo) Get(ctx context.Context, scope, namespace, name string) (*service.HTTPRoute, error) {
	clientset, err := r.kubernetesSigs.gaClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.GatewayV1().HTTPRoutes(namespace).Get(ctx, name, opts)
}

func (r *httpRouteRepo) Create(ctx context.Context, scope, namespace string, hr *service.HTTPRoute) (*service.HTTPRoute, error) {
	clientset, err := r.kubernetesSigs.gaClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.GatewayV1().HTTPRoutes(namespace).Create(ctx, hr, opts)
}

func (r *httpRouteRepo) Update(ctx context.Context, scope, namespace string, hr *service.HTTPRoute) (*service.HTTPRoute, error) {
	clientset, err := r.kubernetesSigs.gaClientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.GatewayV1().HTTPRoutes(namespace).Update(ctx, hr, opts)
}

func (r *httpRouteRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetesSigs.gaClientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.GatewayV1().HTTPRoutes(namespace).Delete(ctx, name, opts)
}
