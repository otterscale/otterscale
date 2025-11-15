package config

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

// ConfigMap represents a Kubernetes ConfigMap resource.
type ConfigMap = v1.ConfigMap

//nolint:revive // allows this exported interface name for specific domain clarity.
type ConfigMapRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]ConfigMap, error)
	Get(ctx context.Context, scope, namespace, name string) (*ConfigMap, error)
	Create(ctx context.Context, scope, namespace string, cm *ConfigMap) (*ConfigMap, error)
	Update(ctx context.Context, scope, namespace string, cm *ConfigMap) (*ConfigMap, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}
