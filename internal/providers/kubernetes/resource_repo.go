package kubernetes

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/utils/ptr"

	"github.com/otterscale/otterscale/internal/core/resource"
)

type resourceRepo struct {
	kubernetes *Kubernetes
}

func NewresourceRepo(kubernetes *Kubernetes) resource.ResourceRepo {
	return &resourceRepo{
		kubernetes: kubernetes,
	}
}

var _ resource.ResourceRepo = (*resourceRepo)(nil)

func (r *resourceRepo) List(ctx context.Context, cluster, group, version, resource, namespace, labelSelector, fieldSelector string, limit int64, continueToken string) (*unstructured.UnstructuredList, error) {
	dynamic, err := r.kubernetes.dynamic(cluster, "", nil) // from context
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	opts := metav1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
		Limit:         limit,
		Continue:      continueToken,
	}

	return dynamic.Resource(gvr).Namespace(namespace).List(ctx, opts)
}

func (r *resourceRepo) Get(ctx context.Context, cluster, group, version, resource, namespace, name string) (*unstructured.Unstructured, error) {
	dynamic, err := r.kubernetes.dynamic(cluster, "", nil) // from context
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	opts := metav1.GetOptions{}

	return dynamic.Resource(gvr).Namespace(namespace).Get(ctx, name, opts)
}

func (r *resourceRepo) Create(ctx context.Context, cluster, group, version, resource, namespace string, manifest []byte) (*unstructured.Unstructured, error) {
	dynamic, err := r.kubernetes.dynamic(cluster, "", nil) // from context
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	obj := &unstructured.Unstructured{}
	if err := obj.UnmarshalJSON(manifest); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid manifest: %v", err))
	}

	opts := metav1.CreateOptions{}

	return dynamic.Resource(gvr).Namespace(namespace).Create(ctx, obj, opts)
}

func (r *resourceRepo) Apply(ctx context.Context, cluster, group, version, resource, namespace, name string, manifest []byte, force bool, fieldManager string) (*unstructured.Unstructured, error) {
	dynamic, err := r.kubernetes.dynamic(cluster, "", nil) // from context
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	opts := metav1.PatchOptions{
		Force:        ptr.To(force),
		FieldManager: fieldManager,
	}

	return dynamic.Resource(gvr).Namespace(namespace).Patch(ctx, name, types.ApplyPatchType, manifest, opts)
}

func (r *resourceRepo) Delete(ctx context.Context, cluster, group, version, resource, namespace, name string, gracePeriodSeconds *int64) error {
	dynamic, err := r.kubernetes.dynamic(cluster, "", nil) // from context
	if err != nil {
		return err
	}

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	opts := metav1.DeleteOptions{
		GracePeriodSeconds: gracePeriodSeconds,
	}

	return dynamic.Resource(gvr).Namespace(namespace).Delete(ctx, name, opts)
}

func (r *resourceRepo) Watch(ctx context.Context, cluster, group, version, resource, namespace, labelSelector, fieldSelector, resourceVersion string) (watch.Interface, error) {
	dynamic, err := r.kubernetes.dynamic(cluster, "", nil) // from context
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	opts := metav1.ListOptions{
		LabelSelector:   labelSelector,
		FieldSelector:   fieldSelector,
		ResourceVersion: resourceVersion,
	}

	return dynamic.Resource(gvr).Namespace(namespace).Watch(ctx, opts)
}
