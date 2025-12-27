package resource

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ClusterGroupVersionResource struct {
	Cluster string
	schema.GroupVersionResource
}

type DiscoveryRepo interface {
	List(cluster string) ([]*metav1.APIResourceList, error)
	Validate(cluster, group, version, resource string) (ClusterGroupVersionResource, error)
}
