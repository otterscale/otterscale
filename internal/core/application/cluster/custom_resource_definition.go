package cluster

import (
	"context"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type CustomResourceDefinition = apiextensionsv1.CustomResourceDefinition

type CustomResourceDefinitionRepo interface {
	List(ctx context.Context, scope, selector string) ([]CustomResourceDefinition, error)
	Get(ctx context.Context, scope, name string) (*CustomResourceDefinition, error)
	Update(ctx context.Context, scope string, crd *CustomResourceDefinition) (*CustomResourceDefinition, error)
	Create(ctx context.Context, scope string, crd *CustomResourceDefinition) (*CustomResourceDefinition, error)
	Delete(ctx context.Context, scope, name string) error
}
