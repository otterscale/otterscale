package kubernetes

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

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

	for _, res := range resources.APIResources {
		if res.Name == gvr.Resource {
			return resource.ClusterGroupVersionResource{
				Cluster:              cluster,
				GroupVersionResource: gvr,
			}, nil
		}
	}

	return resource.ClusterGroupVersionResource{}, fmt.Errorf("resource %s not found in group version %s", gvr.Resource, gvr.GroupVersion().String())
}
