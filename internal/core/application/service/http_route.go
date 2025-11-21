package service

import (
	"context"

	v1 "sigs.k8s.io/gateway-api/apis/v1"
)

// HTTPRoute represents a Kubernetes HTTPRoute resource.
type HTTPRoute = v1.HTTPRoute

type HTTPRouteRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]HTTPRoute, error)
	Get(ctx context.Context, scope, namespace, name string) (*HTTPRoute, error)
	Update(ctx context.Context, scope, namespace string, hr *HTTPRoute) (*HTTPRoute, error)
	Create(ctx context.Context, scope, namespace string, hr *HTTPRoute) (*HTTPRoute, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}
