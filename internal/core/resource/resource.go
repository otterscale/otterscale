package resource

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

type ResourceRepo interface {
	List(ctx context.Context, cgvr ClusterGroupVersionResource, namespace, labelSelector, fieldSelector string, limit int64, continueToken string) (*unstructured.UnstructuredList, error)
	Get(ctx context.Context, cgvr ClusterGroupVersionResource, namespace, name string) (*unstructured.Unstructured, error)
	Create(ctx context.Context, cgvr ClusterGroupVersionResource, namespace string, manifest []byte) (*unstructured.Unstructured, error)
	Apply(ctx context.Context, cgvr ClusterGroupVersionResource, namespace, name string, manifest []byte, force bool, fieldManager string) (*unstructured.Unstructured, error)
	Delete(ctx context.Context, cgvr ClusterGroupVersionResource, namespace, name string, gracePeriodSeconds *int64) error
	Watch(ctx context.Context, cgvr ClusterGroupVersionResource, namespace, labelSelector, fieldSelector, resourceVersion string) (watch.Interface, error)
}

type UseCase struct {
	discovery DiscoveryRepo
	resource  ResourceRepo
}

func NewUseCase(discovery DiscoveryRepo, resource ResourceRepo) *UseCase {
	return &UseCase{
		discovery: discovery,
		resource:  resource,
	}
}

func (uc *UseCase) Validate(ctx context.Context, cluster, group, version, resource string) (ClusterGroupVersionResource, error) {
	return uc.discovery.Validate(cluster, group, version, resource)
}

func (uc *UseCase) ListAPIResources(ctx context.Context, cluster string) ([]*metav1.APIResourceList, error) {
	return uc.discovery.List(cluster)
}

func (uc *UseCase) ListResources(ctx context.Context, cgvr ClusterGroupVersionResource, namespace, labelSelector, fieldSelector string, limit int64, continueToken string) (*unstructured.UnstructuredList, error) {
	return uc.resource.List(ctx, cgvr, namespace, labelSelector, fieldSelector, limit, continueToken)
}

func (uc *UseCase) GetResource(ctx context.Context, cgvr ClusterGroupVersionResource, namespace, name string) (*unstructured.Unstructured, error) {
	return uc.resource.Get(ctx, cgvr, namespace, name)
}

func (uc *UseCase) CreateResource(ctx context.Context, cgvr ClusterGroupVersionResource, namespace string, manifest []byte) (*unstructured.Unstructured, error) {
	return uc.resource.Create(ctx, cgvr, namespace, manifest)
}

func (uc *UseCase) ApplyResource(ctx context.Context, cgvr ClusterGroupVersionResource, namespace, name string, manifest []byte, force bool, fieldManager string) (*unstructured.Unstructured, error) {
	return uc.resource.Apply(ctx, cgvr, namespace, name, manifest, force, fieldManager)
}

func (uc *UseCase) DeleteResource(ctx context.Context, cgvr ClusterGroupVersionResource, namespace, name string, gracePeriodSeconds int64) error {
	return uc.resource.Delete(ctx, cgvr, namespace, name, &gracePeriodSeconds)
}

func (uc *UseCase) WatchResource(ctx context.Context, cgvr ClusterGroupVersionResource, namespace, labelSelector, fieldSelector, resourceVersion string) (watch.Interface, error) {
	return uc.resource.Watch(ctx, cgvr, namespace, labelSelector, fieldSelector, resourceVersion)
}

// func (uc *UseCase) cleanObject(obj *unstructured.Unstructured) {
// 	if obj == nil {
// 		return
// 	}

// 	unstructured.RemoveNestedField(obj.Object, "metadata", "managedFields")

// 	annotations, found, _ := unstructured.NestedStringMap(obj.Object, "metadata", "annotations")
// 	if found {
// 		delete(annotations, "kubectl.kubernetes.io/last-applied-configuration")
// 		if len(annotations) == 0 {
// 			unstructured.RemoveNestedField(obj.Object, "metadata", "annotations")
// 		} else {
// 			unstructured.SetNestedStringMap(obj.Object, annotations, "metadata", "annotations")
// 		}
// 	}
// }
