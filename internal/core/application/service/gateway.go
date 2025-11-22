package service

import (
	"context"

	v1 "sigs.k8s.io/gateway-api/apis/v1"
)

// Gateway represents a Kubernetes Gateway resource.
type Gateway = v1.Gateway

type GatewayRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]Gateway, error)
	Get(ctx context.Context, scope, namespace, name string) (*Gateway, error)
	Update(ctx context.Context, scope, namespace string, g *Gateway) (*Gateway, error)
	Create(ctx context.Context, scope, namespace string, g *Gateway) (*Gateway, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}
