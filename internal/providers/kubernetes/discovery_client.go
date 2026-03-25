package kubernetes

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/cel/openapi/resolver"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/kube-openapi/pkg/validation/spec"

	"github.com/otterscale/otterscale/internal/core"
)

// minWatchListVersion is the minimum Kubernetes version that supports
// the WatchList streaming feature (beta, default-on since 1.34).
var minWatchListVersion = semver.MustParse("v1.34.0")

// discoveryClient implements core.DiscoveryClient by delegating to the
// Kubernetes discovery API of the target cluster, accessed through the
// tunnel.
type discoveryClient struct {
	kubernetes *Kubernetes
}

// NewDiscoveryClient returns a core.DiscoveryClient backed by the
// Kubernetes discovery API.
func NewDiscoveryClient(kubernetes *Kubernetes) core.DiscoveryClient {
	return &discoveryClient{
		kubernetes: kubernetes,
	}
}

var _ core.DiscoveryClient = (*discoveryClient)(nil)

// LookupResource verifies that the given group/version/resource triple
// exists on the target cluster. It returns the validated GVR or a
// BadRequest error if the resource is not recognized.
func (d *discoveryClient) LookupResource(ctx context.Context, cluster, group, version, resource, subresource string) (schema.GroupVersionResource, error) {
	client, err := d.client(ctx, cluster)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	resources, err := client.ServerResourcesForGroupVersion(gvr.GroupVersion().String())
	if err != nil {
		return schema.GroupVersionResource{}, wrapK8sError(err)
	}

	target := resource
	if subresource != "" {
		target += "/" + subresource
	}

	for i := range resources.APIResources {
		if resources.APIResources[i].Name == target {
			return gvr, nil
		}
	}
	return schema.GroupVersionResource{}, wrapK8sError(apierrors.NewBadRequest(fmt.Sprintf("unable to recognize resource %s", gvr)))
}

// ServerResources returns the full list of API resources available on
// the target cluster.
func (d *discoveryClient) ServerResources(ctx context.Context, cluster string) ([]*metav1.APIResourceList, error) {
	client, err := d.client(ctx, cluster)
	if err != nil {
		return nil, err
	}

	_, resources, err := client.ServerGroupsAndResources()
	return resources, wrapK8sError(err)
}

// ResolveSchema fetches the OpenAPI schema for the given GVK from the
// target cluster's discovery endpoint.
func (d *discoveryClient) ResolveSchema(ctx context.Context, cluster, group, version, kind string) (*spec.Schema, error) {
	client, err := d.client(ctx, cluster)
	if err != nil {
		return nil, err
	}

	schemaResolver := &resolver.ClientDiscoveryResolver{
		Discovery: client,
	}
	gvk := schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}
	resolved, err := schemaResolver.ResolveSchema(gvk)
	return resolved, wrapK8sError(err)
}

// ServerVersion returns the Kubernetes version of the target cluster.
func (d *discoveryClient) ServerVersion(ctx context.Context, cluster string) (*version.Info, error) {
	client, err := d.client(ctx, cluster)
	if err != nil {
		return nil, err
	}
	info, err := client.ServerVersion()
	return info, wrapK8sError(err)
}

// SupportsWatchList reports whether the target cluster supports the
// WatchList streaming feature (Kubernetes >= 1.34).
// See https://kubernetes.io/docs/reference/using-api/api-concepts/#streaming-lists
func (d *discoveryClient) SupportsWatchList(ctx context.Context, cluster string) (bool, error) {
	info, err := d.ServerVersion(ctx, cluster)
	if err != nil {
		return false, err
	}

	kubeVersion, err := semver.NewVersion(info.String())
	if err != nil {
		return false, err
	}

	return kubeVersion.GreaterThanEqual(minWatchListVersion), nil
}

// IsNamespacedResource reports whether the given resource is
// namespace-scoped (true) or cluster-scoped (false).
func (d *discoveryClient) IsNamespacedResource(ctx context.Context, cluster string, gvr schema.GroupVersionResource) (bool, error) {
	client, err := d.client(ctx, cluster)
	if err != nil {
		return false, err
	}

	resources, err := client.ServerResourcesForGroupVersion(gvr.GroupVersion().String())
	if err != nil {
		return false, wrapK8sError(err)
	}

	for i := range resources.APIResources {
		if resources.APIResources[i].Name == gvr.Resource {
			return resources.APIResources[i].Namespaced, nil
		}
	}

	return false, wrapK8sError(apierrors.NewNotFound(schema.GroupResource{
		Group:    gvr.Group,
		Resource: gvr.Resource,
	}, ""))
}

// client returns a fresh discovery client for the given cluster with
// impersonation headers set for the calling user. A new client is
// created per request because each request may carry different
// impersonation credentials (user subject + groups). The underlying
// HTTP transport is cached per-cluster in Kubernetes.roundTripper, so
// only the Go-level wrapper is allocated per call.
func (d *discoveryClient) client(ctx context.Context, cluster string) (*discovery.DiscoveryClient, error) {
	config, err := d.kubernetes.impersonationConfig(ctx, cluster)
	if err != nil {
		return nil, err
	}

	// Build a discovery client that reuses the cached transport but
	// applies per-request impersonation via a WrapTransport layer.
	dc, err := discovery.NewDiscoveryClientForConfig(rest.CopyConfig(config))
	if err != nil {
		return nil, wrapK8sError(err)
	}
	return dc, nil
}
