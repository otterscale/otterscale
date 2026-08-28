package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/sync/singleflight"
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

// watchListTTL is how long a cluster's watch-list capability is reused
// before it is probed again. The answer only changes when the cluster
// is upgraded, and every Watch call needs it, so a short TTL keeps a
// /version round-trip off the critical path of each new stream.
const watchListTTL = 10 * time.Minute

// watchListSupport is a cached capability answer with its expiry.
type watchListSupport struct {
	supported bool
	expiresAt time.Time
}

// resourceListTTL is how long a group/version's resource list is reused
// before it is fetched again. Every resource operation validates its
// GVR first, so without this each List, Get, Create, Watch or Scale
// costs two round-trips through the tunnel instead of one.
//
// The TTL can be generous because it is not what keeps the list fresh:
// a lookup that misses forces a refresh before reporting the resource
// as unknown, so a newly installed CRD is usable immediately. What the
// TTL bounds is how long a *removed* resource keeps being accepted —
// and the request that follows fails against the API server anyway.
const resourceListTTL = 10 * time.Minute

// maxResourceListEntries bounds the resource cache. Entries are keyed
// per cluster and group/version, so the count grows with clusters
// registered over the process's lifetime rather than with traffic.
const maxResourceListEntries = 2048

// resourceListFetchTimeout bounds a cache-miss fetch, which runs
// detached from the caller that triggered it.
const resourceListFetchTimeout = 30 * time.Second

// resourceListEntry is a cached group/version resource list with its
// expiry.
type resourceListEntry struct {
	resources []metav1.APIResource
	expiresAt time.Time
}

// discoveryClient implements core.DiscoveryClient by delegating to the
// Kubernetes discovery API of the target cluster, accessed through the
// tunnel.
//
// Discovery answers are cached across users. They are impersonated like
// every other call, so in principle a user who may not read discovery
// can be served an answer another user's request populated. In practice
// the built-in system:discovery ClusterRole grants exactly these reads
// to system:authenticated, so the API resource list is universally
// readable; this is the same trade-off the OpenAPI schema cache makes,
// and it should be revisited together with that one if a deployment
// ever restricts discovery.
type discoveryClient struct {
	kubernetes *Kubernetes
	now        func() time.Time

	mu        sync.RWMutex
	watchList map[string]watchListSupport  // keyed by cluster name
	resources map[string]resourceListEntry // keyed by cluster and group/version

	// resourceFlights collapses concurrent misses for the same
	// group/version into one fetch.
	resourceFlights singleflight.Group
}

// NewDiscoveryClient returns a core.DiscoveryClient backed by the
// Kubernetes discovery API.
func NewDiscoveryClient(kubernetes *Kubernetes) core.DiscoveryClient {
	return &discoveryClient{
		kubernetes: kubernetes,
		now:        time.Now,
		watchList:  make(map[string]watchListSupport),
		resources:  make(map[string]resourceListEntry),
	}
}

var _ core.DiscoveryClient = (*discoveryClient)(nil)

// LookupResource verifies that the given group/version/resource triple
// exists on the target cluster. It returns the validated GVR or a
// BadRequest error if the resource is not recognized.
func (d *discoveryClient) LookupResource(ctx context.Context, cluster, group, version, resource, subresource string) (schema.GroupVersionResource, error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	gv := gvr.GroupVersion()

	target := resource
	if subresource != "" {
		target += "/" + subresource
	}

	resources, fromCache, err := d.groupVersionResources(ctx, cluster, gv)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	if containsResource(resources, target) {
		return gvr, nil
	}

	// A miss against a cached list may only mean the list is stale: a
	// CRD installed after it was fetched would not appear in it. Drop
	// the entry and look once more before declaring the resource
	// unknown, so a freshly installed type works immediately instead of
	// after the TTL.
	//
	// This costs one discovery round-trip — precisely what every lookup
	// cost before this cache existed — so the miss path is never slower
	// than the old behavior.
	if fromCache {
		d.invalidateGroupVersion(cluster, gv)

		resources, _, err = d.groupVersionResources(ctx, cluster, gv)
		if err != nil {
			return schema.GroupVersionResource{}, err
		}
		if containsResource(resources, target) {
			return gvr, nil
		}
	}

	return schema.GroupVersionResource{}, wrapK8sError(apierrors.NewBadRequest(fmt.Sprintf("unable to recognize resource %s", gvr)))
}

// containsResource reports whether name appears in the group/version's
// resource list. Subresources are listed under "resource/subresource",
// which is why callers join them before looking up.
func containsResource(resources []metav1.APIResource, name string) bool {
	for i := range resources {
		if resources[i].Name == name {
			return true
		}
	}
	return false
}

// groupVersionResources returns the resource list for gv on cluster,
// serving it from the cache when possible. fromCache reports whether
// the answer came from a stored entry, which is what tells the caller a
// miss is worth re-checking against a fresh list.
func (d *discoveryClient) groupVersionResources(
	ctx context.Context,
	cluster string,
	gv schema.GroupVersion,
) (resources []metav1.APIResource, fromCache bool, err error) {
	key := resourceCacheKey(cluster, gv)

	if cached, ok := d.cachedResources(key); ok {
		return cached, true, nil
	}

	v, err, _ := d.resourceFlights.Do(key, func() (any, error) {
		// Detach from the caller's cancellation: waiters sharing this
		// flight must not all fail because whoever triggered the fetch
		// gave up. Values are preserved, so the request still carries
		// its impersonated identity.
		fetchCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), resourceListFetchTimeout)
		defer cancel()

		client, err := d.client(fetchCtx, cluster)
		if err != nil {
			return nil, err
		}

		list, err := client.ServerResourcesForGroupVersion(gv.String())
		if err != nil {
			return nil, wrapK8sError(err)
		}

		d.cacheResources(key, list.APIResources)
		return list.APIResources, nil
	})
	if err != nil {
		return nil, false, err
	}

	return v.([]metav1.APIResource), false, nil
}

// resourceCacheKey builds the cache key for a cluster's group/version.
func resourceCacheKey(cluster string, gv schema.GroupVersion) string {
	return cluster + "/" + gv.String()
}

// cachedResources returns the resource list for key if it has not
// expired.
func (d *discoveryClient) cachedResources(key string) ([]metav1.APIResource, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entry, found := d.resources[key]
	if !found || !d.now().Before(entry.expiresAt) {
		return nil, false
	}
	return entry.resources, true
}

// cacheResources stores a resource list until resourceListTTL has
// elapsed. When the cache is full, expired entries are dropped first;
// if that frees nothing, the result simply goes uncached rather than
// evicting an entry another cluster is actively using.
func (d *discoveryClient) cacheResources(key string, resources []metav1.APIResource) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.resources) >= maxResourceListEntries {
		d.evictExpiredResources()
	}
	if len(d.resources) >= maxResourceListEntries {
		return
	}

	d.resources[key] = resourceListEntry{
		resources: resources,
		expiresAt: d.now().Add(resourceListTTL),
	}
}

// invalidateGroupVersion drops a cached resource list and any in-flight
// fetch for it, so the next lookup starts a new one rather than joining
// a flight that would return the same stale list.
func (d *discoveryClient) invalidateGroupVersion(cluster string, gv schema.GroupVersion) {
	key := resourceCacheKey(cluster, gv)

	d.mu.Lock()
	delete(d.resources, key)
	d.mu.Unlock()

	d.resourceFlights.Forget(key)
}

// evictExpiredResources removes expired entries. Must be called with mu
// held for writing.
func (d *discoveryClient) evictExpiredResources() {
	now := d.now()
	for key, entry := range d.resources {
		if now.After(entry.expiresAt) {
			delete(d.resources, key)
		}
	}
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
// WatchList streaming feature (Kubernetes >= 1.34). The answer is
// cached per cluster for watchListTTL.
// See https://kubernetes.io/docs/reference/using-api/api-concepts/#streaming-lists
func (d *discoveryClient) SupportsWatchList(ctx context.Context, cluster string) (bool, error) {
	if supported, ok := d.cachedWatchList(cluster); ok {
		return supported, nil
	}

	info, err := d.ServerVersion(ctx, cluster)
	if err != nil {
		return false, err
	}

	// Distributions decorate the version with build metadata
	// (v1.34.1+k3s1, v1.34.1-gke.100), which semver parses as
	// prerelease or build identifiers rather than rejecting.
	kubeVersion, err := semver.NewVersion(info.String())
	if err != nil {
		return false, fmt.Errorf("parse cluster version %q: %w", info.String(), err)
	}

	supported := kubeVersion.GreaterThanEqual(minWatchListVersion)
	d.cacheWatchList(cluster, supported)
	return supported, nil
}

// cachedWatchList returns the cached capability for cluster, if it has
// not expired.
func (d *discoveryClient) cachedWatchList(cluster string) (supported, ok bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entry, found := d.watchList[cluster]
	if !found || !d.now().Before(entry.expiresAt) {
		return false, false
	}
	return entry.supported, true
}

// cacheWatchList stores the capability for cluster until watchListTTL
// has elapsed.
func (d *discoveryClient) cacheWatchList(cluster string, supported bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.watchList[cluster] = watchListSupport{
		supported: supported,
		expiresAt: d.now().Add(watchListTTL),
	}
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

	// The config already carries this request's impersonation headers
	// and the cluster's cached transport; it is copied so the discovery
	// client cannot mutate the caller's copy.
	dc, err := discovery.NewDiscoveryClientForConfig(rest.CopyConfig(config))
	if err != nil {
		return nil, wrapK8sError(err)
	}
	return dc, nil
}
