package kubernetes

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/utils/ptr"

	"github.com/otterscale/otterscale/internal/core/resource"
)

type resourceRepo struct {
	kubernetes *Kubernetes
}

func NewResourceRepo(kubernetes *Kubernetes) resource.ResourceRepo {
	return &resourceRepo{
		kubernetes: kubernetes,
	}
}

var _ resource.ResourceRepo = (*resourceRepo)(nil)

func (r *resourceRepo) List(ctx context.Context, cgvr resource.ClusterGroupVersionResource, namespace, labelSelector, fieldSelector string, limit int64, continueToken string) (*unstructured.UnstructuredList, error) {
	client, err := r.kubernetes.dynamic(cgvr.Cluster, "", nil) // from context
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
		Limit:         limit,
		Continue:      continueToken,
	}

	return client.Resource(cgvr.GroupVersionResource).Namespace(namespace).List(ctx, opts)
}

func (r *resourceRepo) Get(ctx context.Context, cgvr resource.ClusterGroupVersionResource, namespace, name string) (*unstructured.Unstructured, error) {
	client, err := r.kubernetes.dynamic(cgvr.Cluster, "", nil) // from context
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return client.Resource(cgvr.GroupVersionResource).Namespace(namespace).Get(ctx, name, opts)
}

func (r *resourceRepo) Create(ctx context.Context, cgvr resource.ClusterGroupVersionResource, namespace string, manifest []byte) (*unstructured.Unstructured, error) {
	client, err := r.kubernetes.dynamic(cgvr.Cluster, "", nil) // from context
	if err != nil {
		return nil, err
	}

	obj := &unstructured.Unstructured{}
	if err := obj.UnmarshalJSON(manifest); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid manifest: %v", err))
	}

	opts := metav1.CreateOptions{}

	return client.Resource(cgvr.GroupVersionResource).Namespace(namespace).Create(ctx, obj, opts)
}

func (r *resourceRepo) Apply(ctx context.Context, cgvr resource.ClusterGroupVersionResource, namespace, name string, manifest []byte, force bool, fieldManager string) (*unstructured.Unstructured, error) {
	client, err := r.kubernetes.dynamic(cgvr.Cluster, "", nil) // from context
	if err != nil {
		return nil, err
	}

	opts := metav1.PatchOptions{
		Force:        ptr.To(force),
		FieldManager: fieldManager,
	}

	return client.Resource(cgvr.GroupVersionResource).Namespace(namespace).Patch(ctx, name, types.ApplyPatchType, manifest, opts)
}

func (r *resourceRepo) Delete(ctx context.Context, cgvr resource.ClusterGroupVersionResource, namespace, name string, gracePeriodSeconds *int64) error {
	client, err := r.kubernetes.dynamic(cgvr.Cluster, "", nil) // from context
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{
		GracePeriodSeconds: gracePeriodSeconds,
	}

	return client.Resource(cgvr.GroupVersionResource).Namespace(namespace).Delete(ctx, name, opts)
}

func (r *resourceRepo) Watch(ctx context.Context, cgvr resource.ClusterGroupVersionResource, namespace, labelSelector, fieldSelector, resourceVersion string) (watch.Interface, error) {
	client, err := r.kubernetes.dynamic(cgvr.Cluster, "", nil) // from context
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector:   labelSelector,
		FieldSelector:   fieldSelector,
		ResourceVersion: resourceVersion,
	}

	return client.Resource(cgvr.GroupVersionResource).Namespace(namespace).Watch(ctx, opts)
}
