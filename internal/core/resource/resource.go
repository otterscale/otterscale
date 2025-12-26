package resource

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

type ResourceRepo interface {
	List(ctx context.Context, cluster, group, version, resource, namespace, labelSelector, fieldSelector string, limit int64, continueToken string) (*unstructured.UnstructuredList, error)
	Get(ctx context.Context, cluster, group, version, resource, namespace, name string) (*unstructured.Unstructured, error)
	Create(ctx context.Context, cluster, group, version, resource, namespace string, manifest []byte) (*unstructured.Unstructured, error)
	Apply(ctx context.Context, cluster, group, version, resource, namespace, name string, manifest []byte, force bool, fieldManager string) (*unstructured.Unstructured, error)
	Delete(ctx context.Context, cluster, group, version, resource, namespace, name string, gracePeriodSeconds *int64) error
	Watch(ctx context.Context, cluster, group, version, resource, namespace, labelSelector, fieldSelector, resourceVersion string) (watch.Interface, error)
}

type UseCase struct {
	resource ResourceRepo
}

func NewUseCase(resource ResourceRepo) *UseCase {
	return &UseCase{
		resource: resource,
	}
}

func (uc *UseCase) ListResources(ctx context.Context, cluster, group, version, resource, namespace, labelSelector, fieldSelector string, limit int64, continueToken string) (*unstructured.UnstructuredList, error) {
	return uc.resource.List(ctx, cluster, group, version, resource, namespace, labelSelector, fieldSelector, limit, continueToken)
}

func (uc *UseCase) GetResource(ctx context.Context, cluster, group, version, resource, namespace, name string) (*unstructured.Unstructured, error) {
	return uc.resource.Get(ctx, cluster, group, version, resource, namespace, name)
}

func (uc *UseCase) CreateResource(ctx context.Context, cluster, group, version, resource, namespace string, manifest []byte) (*unstructured.Unstructured, error) {
	return uc.resource.Create(ctx, cluster, group, version, resource, namespace, manifest)
}

func (uc *UseCase) ApplyResource(ctx context.Context, cluster, group, version, resource, namespace, name string, manifest []byte, force bool, fieldManager string) (*unstructured.Unstructured, error) {
	return uc.resource.Apply(ctx, cluster, group, version, resource, namespace, name, manifest, force, fieldManager)
}

func (uc *UseCase) DeleteResource(ctx context.Context, cluster, group, version, resource, namespace, name string, gracePeriodSeconds int64) error {
	return uc.resource.Delete(ctx, cluster, group, version, resource, namespace, name, &gracePeriodSeconds)
}

func (uc *UseCase) WatchResource(ctx context.Context, cluster, group, version, resource, namespace, labelSelector, fieldSelector, resourceVersion string) (watch.Interface, error) {
	return uc.resource.Watch(ctx, cluster, group, version, resource, namespace, labelSelector, fieldSelector, resourceVersion)
}
