package kubernetes

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/cel/openapi/resolver"
	"k8s.io/kube-openapi/pkg/validation/spec"

	"github.com/otterscale/otterscale/internal/core/resource"
)

type discoveryRepo struct {
	kubernetes *Kubernetes
}

func NewDiscoveryRepo(kubernetes *Kubernetes) resource.DiscoveryRepo {
	return &discoveryRepo{
		kubernetes: kubernetes,
	}
}

var _ resource.DiscoveryRepo = (*discoveryRepo)(nil)

func (r *discoveryRepo) List(cluster string) ([]*metav1.APIResourceList, error) {
	client, err := r.kubernetes.discovery(cluster)
	if err != nil {
		return nil, err
	}

	_, resources, err := client.ServerGroupsAndResources()
	return resources, err
}

func (r *discoveryRepo) Schema(cluster, group, version, kind string) (*spec.Schema, error) {
	client, err := r.kubernetes.discovery(cluster)
	if err != nil {
		return nil, err
	}

	resolver := &resolver.ClientDiscoveryResolver{
		Discovery: client,
	}

	gvk := schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}

	return resolver.ResolveSchema(gvk)
}

func (r *discoveryRepo) Validate(cluster, group, version, res string) (resource.ClusterGroupVersionResource, error) {
	client, err := r.kubernetes.discovery(cluster)
	if err != nil {
		return resource.ClusterGroupVersionResource{}, err
	}

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: res,
	}

	resources, err := client.ServerResourcesForGroupVersion(gvr.GroupVersion().String())
	if err != nil {
		return resource.ClusterGroupVersionResource{}, err
	}

	for i := range resources.APIResources {
		if resources.APIResources[i].Name == gvr.Resource {
			return resource.ClusterGroupVersionResource{
				Cluster:              cluster,
				GroupVersionResource: gvr,
			}, nil
		}
	}

	return resource.ClusterGroupVersionResource{}, apierrors.NewNotFound(schema.GroupResource{Group: gvr.Group, Resource: gvr.Resource}, gvr.Resource)
}

func (r *discoveryRepo) Version(cluster string) (*version.Info, error) {
	client, err := r.kubernetes.discovery(cluster)
	if err != nil {
		return nil, err
	}

	return client.ServerVersion()
}
