package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/cel/openapi/resolver"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/kube-openapi/pkg/validation/spec"

	"github.com/otterscale/otterscale/internal/core"
)

// extGVK is the OpenAPI v3 extension key identifying the
// GroupVersionKind tuple(s) that a schema definition represents.
const extGVK = "x-kubernetes-group-version-kind"

// refPrefix is the OpenAPI v3 prefix on JSON Pointer refs into the
// schemas component.
const refPrefix = "#/components/schemas/"

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

// ResolveGroupVersionSchemas fetches the OpenAPI schemas for every
// kind in the given group/version from the target cluster's discovery
// endpoint. The Kubernetes OpenAPI v3 endpoint serves one document per
// group/version, so fetching kinds in bulk amortizes a single HTTP
// roundtrip across all kinds the caller subsequently looks up.
func (d *discoveryClient) ResolveGroupVersionSchemas(ctx context.Context, cluster, group, version string) (map[string]*spec.Schema, error) {
	client, err := d.client(ctx, cluster)
	if err != nil {
		return nil, err
	}

	paths, err := client.OpenAPIV3().Paths()
	if err != nil {
		return nil, wrapK8sError(err)
	}

	gv := schema.GroupVersion{Group: group, Version: version}
	gvPath, ok := paths[resourcePathFromGV(gv)]
	if !ok {
		return nil, fmt.Errorf("cannot resolve group version %q: %w", gv, resolver.ErrSchemaNotFound)
	}

	raw, err := gvPath.Schema(runtime.ContentTypeJSON)
	if err != nil {
		return nil, wrapK8sError(err)
	}

	var doc struct {
		Components struct {
			Schemas map[string]*spec.Schema `json:"schemas"`
		} `json:"components"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	schemaOf := func(ref string) (*spec.Schema, bool) {
		s, ok := doc.Components.Schemas[strings.TrimPrefix(ref, refPrefix)]
		return s, ok
	}

	result := make(map[string]*spec.Schema)
	for ref, s := range doc.Components.Schemas {
		var gvks []schema.GroupVersionKind
		if err := s.Extensions.GetObject(extGVK, &gvks); err != nil {
			return nil, err
		}
		for _, g := range gvks {
			if g.Group != group || g.Version != version {
				continue
			}
			populated, err := resolver.PopulateRefs(schemaOf, refPrefix+ref)
			if err != nil {
				return nil, err
			}
			result[g.Kind] = populated
		}
	}
	return result, nil
}

// resourcePathFromGV mirrors the path scheme served by the Kubernetes
// OpenAPI v3 endpoint: core types under api/<version>, named groups
// under apis/<group>/<version>.
func resourcePathFromGV(gv schema.GroupVersion) string {
	if gv.Group == "" {
		return fmt.Sprintf("api/%s", gv.Version)
	}
	return fmt.Sprintf("apis/%s/%s", gv.Group, gv.Version)
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
