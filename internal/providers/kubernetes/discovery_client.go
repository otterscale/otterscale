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

// extGVK is the OpenAPI v3 extension key naming the GroupVersionKind tuple(s)
// a schema definition represents.
const extGVK = "x-kubernetes-group-version-kind"

// refPrefix is the OpenAPI v3 prefix on JSON Pointer refs into the schemas
// component.
const refPrefix = "#/components/schemas/"

// minWatchListVersion is the first Kubernetes version with WatchList streaming
// (beta, default-on since 1.34).
var minWatchListVersion = semver.MustParse("v1.34.0")

// watchListTTL bounds how long a cluster's watch-list capability is reused.
// The answer only changes on a cluster upgrade, and every Watch call needs it,
// so caching keeps a /version round-trip off the critical path of each stream.
const watchListTTL = 10 * time.Minute

type watchListSupport struct {
	supported bool
	expiresAt time.Time
}

// resourceListTTL bounds how long a group/version's resource list is reused.
// Every resource operation validates its GVR first, so without this each List,
// Get, Create, Watch or Scale costs two round-trips through the tunnel.
//
// The TTL can be generous because it is not what keeps the list fresh: a lookup
// that misses forces a refresh before reporting the resource as unknown, so a
// newly installed CRD is usable immediately. What the TTL bounds is how long a
// *removed* resource keeps being accepted — and the request that follows fails
// against the API server anyway.
const resourceListTTL = 10 * time.Minute

// maxResourceListEntries bounds the resource cache. Entries are keyed per
// cluster and group/version, so the count grows with clusters registered over
// the process's lifetime rather than with traffic.
const maxResourceListEntries = 2048

// resourceListFetchTimeout bounds a cache-miss fetch, which runs detached from
// the caller that triggered it.
const resourceListFetchTimeout = 30 * time.Second

type resourceListEntry struct {
	resources []metav1.APIResource
	expiresAt time.Time
}

// discoveryClient implements core.DiscoveryClient against the target cluster's
// discovery API, reached through the tunnel.
//
// Discovery answers are cached across users. They are impersonated like every
// other call, so in principle a user who may not read discovery can be served
// an answer another user's request populated. In practice the built-in
// system:discovery ClusterRole grants exactly these reads to
// system:authenticated, so the API resource list is universally readable; this
// is the same trade-off the OpenAPI schema cache makes, and the two should be
// revisited together if a deployment ever restricts discovery.
type discoveryClient struct {
	kubernetes *Kubernetes
	now        func() time.Time

	mu        sync.RWMutex
	watchList map[string]watchListSupport  // keyed by cluster name
	resources map[string]resourceListEntry // keyed by cluster and group/version

	// resourceFlights collapses concurrent misses for the same group/version
	// into one fetch.
	resourceFlights singleflight.Group
}

func NewDiscoveryClient(kubernetes *Kubernetes) core.DiscoveryClient {
	return &discoveryClient{
		kubernetes: kubernetes,
		now:        time.Now,
		watchList:  make(map[string]watchListSupport),
		resources:  make(map[string]resourceListEntry),
	}
}

var _ core.DiscoveryClient = (*discoveryClient)(nil)

// LookupResource returns the validated GVR, or a BadRequest error if the
// cluster does not recognize the resource.
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

	// A miss against a cached list may only mean the list is stale: a CRD
	// installed after it was fetched would not appear in it. Drop the entry and
	// look once more before declaring the resource unknown, so a freshly
	// installed type works immediately instead of after the TTL. That costs one
	// discovery round-trip — precisely what every lookup cost before this cache
	// existed — so the miss path is never slower than the old behavior.
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

// containsResource expects subresources already joined as
// "resource/subresource", which is how discovery lists them.
func containsResource(resources []metav1.APIResource, name string) bool {
	for i := range resources {
		if resources[i].Name == name {
			return true
		}
	}
	return false
}

// groupVersionResources serves the resource list from the cache when possible.
// fromCache is what tells the caller a miss is worth re-checking against a
// fresh list.
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
		// Detach from the caller's cancellation: waiters sharing this flight
		// must not all fail because whoever triggered the fetch gave up. Values
		// are preserved, so the request still carries its impersonated identity.
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

func resourceCacheKey(cluster string, gv schema.GroupVersion) string {
	return cluster + "/" + gv.String()
}

func (d *discoveryClient) cachedResources(key string) ([]metav1.APIResource, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entry, found := d.resources[key]
	if !found || !d.now().Before(entry.expiresAt) {
		return nil, false
	}
	return entry.resources, true
}

// cacheResources drops expired entries first when the cache is full; if that
// frees nothing, the result goes uncached rather than evicting an entry another
// cluster is actively using.
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

// invalidateGroupVersion also forgets any in-flight fetch, so the next lookup
// starts a new one rather than joining a flight returning the same stale list.
func (d *discoveryClient) invalidateGroupVersion(cluster string, gv schema.GroupVersion) {
	key := resourceCacheKey(cluster, gv)

	d.mu.Lock()
	delete(d.resources, key)
	d.mu.Unlock()

	d.resourceFlights.Forget(key)
}

// evictExpiredResources must be called with mu held for writing.
func (d *discoveryClient) evictExpiredResources() {
	now := d.now()
	for key, entry := range d.resources {
		if now.After(entry.expiresAt) {
			delete(d.resources, key)
		}
	}
}

func (d *discoveryClient) ServerResources(ctx context.Context, cluster string) ([]*metav1.APIResourceList, error) {
	client, err := d.client(ctx, cluster)
	if err != nil {
		return nil, err
	}

	_, resources, err := client.ServerGroupsAndResources()
	return resources, wrapK8sError(err)
}

// ResolveGroupVersionSchemas fetches every kind in one call: the Kubernetes
// OpenAPI v3 endpoint serves one document per group/version, so the single
// round-trip is amortized across all kinds the caller looks up afterwards.
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
		// A group/version the cluster does not serve is a bad request, not a
		// server fault. The sentinel stays in the chain so errors.Is keeps
		// working for callers that check it.
		return nil, &core.DomainError{
			Code:    core.ErrorCodeNotFound,
			Message: fmt.Sprintf("cannot resolve group version %q", gv),
			Cause:   resolver.ErrSchemaNotFound,
		}
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
		return nil, &core.DomainError{
			Code:    core.ErrorCodeInternal,
			Message: "decode OpenAPI document",
			Cause:   err,
		}
	}

	schemaOf := func(ref string) (*spec.Schema, bool) {
		s, ok := doc.Components.Schemas[strings.TrimPrefix(ref, refPrefix)]
		return s, ok
	}

	result := make(map[string]*spec.Schema)
	for ref, s := range doc.Components.Schemas {
		var gvks []schema.GroupVersionKind
		if err := s.Extensions.GetObject(extGVK, &gvks); err != nil {
			return nil, &core.DomainError{
				Code:    core.ErrorCodeInternal,
				Message: fmt.Sprintf("read %s extension of %q", extGVK, ref),
				Cause:   err,
			}
		}
		for _, g := range gvks {
			if g.Group != group || g.Version != version {
				continue
			}
			populated, err := resolver.PopulateRefs(schemaOf, refPrefix+ref)
			if err != nil {
				return nil, &core.DomainError{
					Code:    core.ErrorCodeInternal,
					Message: fmt.Sprintf("resolve schema refs for %q", ref),
					Cause:   err,
				}
			}
			result[g.Kind] = populated
		}
	}
	return result, nil
}

// resourcePathFromGV mirrors the OpenAPI v3 path scheme: core types under
// api/<version>, named groups under apis/<group>/<version>.
func resourcePathFromGV(gv schema.GroupVersion) string {
	if gv.Group == "" {
		return fmt.Sprintf("api/%s", gv.Version)
	}
	return fmt.Sprintf("apis/%s/%s", gv.Group, gv.Version)
}

func (d *discoveryClient) ServerVersion(ctx context.Context, cluster string) (*version.Info, error) {
	client, err := d.client(ctx, cluster)
	if err != nil {
		return nil, err
	}
	info, err := client.ServerVersion()
	return info, wrapK8sError(err)
}

// SupportsWatchList caches its answer per cluster for watchListTTL.
// See https://kubernetes.io/docs/reference/using-api/api-concepts/#streaming-lists
func (d *discoveryClient) SupportsWatchList(ctx context.Context, cluster string) (bool, error) {
	if supported, ok := d.cachedWatchList(cluster); ok {
		return supported, nil
	}

	info, err := d.ServerVersion(ctx, cluster)
	if err != nil {
		return false, err
	}

	// Distributions decorate the version with build metadata (v1.34.1+k3s1,
	// v1.34.1-gke.100), which semver parses as prerelease or build identifiers
	// rather than rejecting.
	kubeVersion, err := semver.NewVersion(info.String())
	if err != nil {
		return false, fmt.Errorf("parse cluster version %q: %w", info.String(), err)
	}

	supported := kubeVersion.GreaterThanEqual(minWatchListVersion)
	d.cacheWatchList(cluster, supported)
	return supported, nil
}

func (d *discoveryClient) cachedWatchList(cluster string) (supported, ok bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entry, found := d.watchList[cluster]
	if !found || !d.now().Before(entry.expiresAt) {
		return false, false
	}
	return entry.supported, true
}

func (d *discoveryClient) cacheWatchList(cluster string, supported bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.watchList[cluster] = watchListSupport{
		supported: supported,
		expiresAt: d.now().Add(watchListTTL),
	}
}

// client builds a fresh discovery client per request, because each request may
// carry different impersonation credentials. The underlying HTTP transport is
// cached per cluster in Kubernetes.roundTripper, so only the Go-level wrapper
// is allocated per call.
func (d *discoveryClient) client(ctx context.Context, cluster string) (*discovery.DiscoveryClient, error) {
	config, err := d.kubernetes.impersonationConfig(ctx, cluster)
	if err != nil {
		return nil, err
	}

	// The config already carries this request's impersonation headers and the
	// cluster's cached transport; it is copied so the discovery client cannot
	// mutate the caller's copy.
	dc, err := discovery.NewDiscoveryClientForConfig(rest.CopyConfig(config))
	if err != nil {
		return nil, wrapK8sError(err)
	}
	return dc, nil
}
