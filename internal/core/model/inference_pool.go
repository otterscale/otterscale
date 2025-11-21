package model

import (
	"context"

	v1 "sigs.k8s.io/gateway-api-inference-extension/api/v1"
)

// InferencePool represents a Kubernetes InferencePool resource.
type InferencePool = v1.InferencePool

type InferencePoolRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]InferencePool, error)
	Get(ctx context.Context, scope, namespace, name string) (*InferencePool, error)
	Update(ctx context.Context, scope, namespace string, ip *InferencePool) (*InferencePool, error)
	Create(ctx context.Context, scope, namespace string, ip *InferencePool) (*InferencePool, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}
