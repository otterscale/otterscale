package core

// ADR: Kubernetes types in the domain layer
//
// This file imports several k8s.io packages (apimachinery, kube-openapi)
// directly into the core (domain) layer. In a strict DDD interpretation,
// domain types should be infrastructure-agnostic. However, otterscale's
// core business *is* Kubernetes resource management: GVR, Unstructured,
// APIResourceList, and OpenAPI Schema are part of the domain's Ubiquitous
// Language, not incidental infrastructure details.
//
// Wrapping these types in custom DTOs would introduce a costly
// translation layer at every boundary with no material benefit — the
// domain would still be structurally identical to the K8s types.
//
// Trade-off accepted: we allow k8s.io/apimachinery and kube-openapi
// imports in core, treating them as domain-level vocabulary rather than
// infrastructure leakage. This decision should be revisited if the
// project ever needs to support non-Kubernetes backends.

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

// ---------------------------------------------------------------------------
// Interfaces
// ---------------------------------------------------------------------------

// DiscoveryClient abstracts Kubernetes API discovery so that the
// use-case layer can validate resources and fetch schemas without
// depending on a concrete client implementation.
type DiscoveryClient interface {
	// LookupResource validates that a group/version/resource triple
	// exists on the target cluster.
	LookupResource(ctx context.Context, cluster, group, version, resource string) (schema.GroupVersionResource, error)
	// ServerResources returns all API resources advertised by the cluster.
	ServerResources(ctx context.Context, cluster string) ([]*metav1.APIResourceList, error)
	// ResolveSchema fetches the OpenAPI schema for a given GVK.
	ResolveSchema(ctx context.Context, cluster, group, version, kind string) (*spec.Schema, error)
	// ServerVersion returns the Kubernetes version of the cluster.
	ServerVersion(ctx context.Context, cluster string) (*version.Info, error)
	// SupportsWatchList reports whether the target cluster supports
	// the WatchList streaming feature (Kubernetes >= 1.34).
	SupportsWatchList(ctx context.Context, cluster string) (bool, error)
}

// ResourceRepo abstracts Kubernetes resource CRUD and watch operations
// through the dynamic client. All methods accept a cluster name so
// that the underlying implementation can route requests through the
// correct tunnel.
type ResourceRepo interface {
	// List returns a paged list of resources matching the given options.
	List(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace string, opts ListOptions,
	) (*unstructured.UnstructuredList, error)

	// Get returns a single resource by name.
	Get(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace, name string,
	) (*unstructured.Unstructured, error)

	// Create decodes a YAML manifest and creates a new resource.
	Create(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace string, manifest []byte,
	) (*unstructured.Unstructured, error)

	// Apply decodes a YAML manifest and performs a server-side apply
	// (PATCH with ApplyPatchType) for the given resource.
	Apply(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace, name string, manifest []byte, opts ApplyOptions,
	) (*unstructured.Unstructured, error)

	// Delete removes a resource.
	Delete(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace, name string, opts DeleteOptions,
	) error

	// Watch opens a long-lived watch stream for resources matching the
	// given options.
	Watch(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace string, opts WatchOptions,
	) (Watcher, error)

	// ListEvents returns events matching the given options.
	// Used by DescribeResource to fetch events via involvedObject.uid.
	ListEvents(ctx context.Context, cluster, namespace string, opts ListOptions) (*unstructured.UnstructuredList, error)
}

// ---------------------------------------------------------------------------
// Options types
// ---------------------------------------------------------------------------

// ListOptions configures a resource list or event list query.
// Mirrors the commonly used fields of metav1.ListOptions.
type ListOptions struct {
	LabelSelector string
	FieldSelector string
	Limit         int64
	Continue      string
}

// ApplyOptions configures a server-side apply operation.
// Mirrors the commonly used fields of metav1.PatchOptions.
type ApplyOptions struct {
	Force        bool
	FieldManager string
}

// DeleteOptions configures a resource deletion.
// Mirrors the commonly used fields of metav1.DeleteOptions.
type DeleteOptions struct {
	GracePeriodSeconds *int64
}

// WatchOptions configures a watch stream.
// Mirrors the commonly used fields of metav1.ListOptions for watch.
type WatchOptions struct {
	LabelSelector     string
	FieldSelector     string
	ResourceVersion   string
	SendInitialEvents bool
}

// SchemaResolver resolves OpenAPI schemas for Kubernetes GVKs.
// Implementations may cache results and deduplicate concurrent
// requests. Defining this as an interface decouples the use-case
// layer from the caching infrastructure.
type SchemaResolver interface {
	ResolveSchema(ctx context.Context, cluster, group, version, kind string) (*spec.Schema, error)
}

// ---------------------------------------------------------------------------
// Identifiers
// ---------------------------------------------------------------------------

// ResourceIdentifier identifies a Kubernetes resource type (and
// optionally a specific instance) across clusters. It replaces long
// positional parameter lists in use-case methods with a single,
// self-documenting value object.
type ResourceIdentifier struct {
	Cluster   string
	Group     string
	Version   string
	Resource  string
	Namespace string
	Name      string
}

// lookupGVR validates the resource triple via the DiscoveryClient.
func (id ResourceIdentifier) lookupGVR(ctx context.Context, dc DiscoveryClient) (schema.GroupVersionResource, error) {
	return dc.LookupResource(ctx, id.Cluster, id.Group, id.Version, id.Resource)
}

// ---------------------------------------------------------------------------
// Use case
// ---------------------------------------------------------------------------

// ResourceUseCase provides the application-level API for managing
// Kubernetes resources across multiple clusters. It validates GVRs
// via the DiscoveryClient and resolves OpenAPI schemas through the
// injected SchemaResolver.
type ResourceUseCase struct {
	discovery      DiscoveryClient
	resource       ResourceRepo
	schemaResolver SchemaResolver
}

// NewResourceUseCase returns a ResourceUseCase wired to the given
// discovery, resource, and schema resolver backends. The
// SchemaResolver is injected to decouple caching infrastructure
// from the domain use-case.
func NewResourceUseCase(discovery DiscoveryClient, resource ResourceRepo, schemaResolver SchemaResolver) *ResourceUseCase {
	return &ResourceUseCase{
		discovery:      discovery,
		resource:       resource,
		schemaResolver: schemaResolver,
	}
}

// ServerResources returns all API resource lists from the target cluster.
func (uc *ResourceUseCase) ServerResources(ctx context.Context, cluster string) ([]*metav1.APIResourceList, error) {
	return uc.discovery.ServerResources(ctx, cluster)
}

// ResolveSchema fetches the OpenAPI schema for the given GVK via the
// injected SchemaResolver.
func (uc *ResourceUseCase) ResolveSchema(
	ctx context.Context,
	cluster, group, version, kind string,
) (*spec.Schema, error) {
	return uc.schemaResolver.ResolveSchema(ctx, cluster, group, version, kind)
}

// ListResources validates the GVR and fetches a paged resource list.
func (uc *ResourceUseCase) ListResources(
	ctx context.Context,
	id ResourceIdentifier,
	opts ListOptions,
) (*unstructured.UnstructuredList, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	return uc.resource.List(ctx, id.Cluster, gvr, id.Namespace, opts)
}

// GetResource validates the GVR and fetches a single resource.
func (uc *ResourceUseCase) GetResource(
	ctx context.Context,
	id ResourceIdentifier,
) (*unstructured.Unstructured, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	return uc.resource.Get(ctx, id.Cluster, gvr, id.Namespace, id.Name)
}

// DescribeResource validates the GVR, fetches the resource, extracts
// its UID, then queries related Kubernetes events filtered by
// involvedObject.uid. This is the backend equivalent of
// `kubectl describe`.
func (uc *ResourceUseCase) DescribeResource(
	ctx context.Context,
	id ResourceIdentifier,
) (*unstructured.Unstructured, *unstructured.UnstructuredList, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, nil, err
	}

	obj, err := uc.resource.Get(ctx, id.Cluster, gvr, id.Namespace, id.Name)
	if err != nil {
		return nil, nil, err
	}

	uid := string(obj.GetUID())

	events, err := uc.resource.ListEvents(ctx, id.Cluster, id.Namespace, ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.uid=%s", uid),
	})
	if err != nil {
		// Events are supplementary; return the resource even if event
		// listing fails (e.g. RBAC restrictions on events).
		return obj, &unstructured.UnstructuredList{}, nil
	}

	return obj, events, nil
}

// CreateResource validates the GVR and creates the resource on the
// target cluster from the given YAML manifest.
func (uc *ResourceUseCase) CreateResource(
	ctx context.Context,
	id ResourceIdentifier,
	manifest []byte,
) (*unstructured.Unstructured, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	return uc.resource.Create(ctx, id.Cluster, gvr, id.Namespace, manifest)
}

// ApplyResource validates the GVR and performs a server-side apply on
// the target cluster from the given YAML manifest.
func (uc *ResourceUseCase) ApplyResource(
	ctx context.Context,
	id ResourceIdentifier,
	manifest []byte,
	opts ApplyOptions,
) (*unstructured.Unstructured, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	return uc.resource.Apply(ctx, id.Cluster, gvr, id.Namespace, id.Name, manifest, opts)
}

// DeleteResource validates the GVR and deletes the named resource.
func (uc *ResourceUseCase) DeleteResource(
	ctx context.Context,
	id ResourceIdentifier,
	opts DeleteOptions,
) error {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return err
	}

	return uc.resource.Delete(ctx, id.Cluster, gvr, id.Namespace, id.Name, opts)
}

// WatchResource validates the GVR and opens a long-lived watch stream.
// If the cluster supports the WatchList feature (Kubernetes >= 1.34),
// initial events are streamed before switching to change notifications.
func (uc *ResourceUseCase) WatchResource(
	ctx context.Context,
	id ResourceIdentifier,
	opts WatchOptions,
) (Watcher, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	watchList, err := uc.discovery.SupportsWatchList(ctx, id.Cluster)
	if err != nil {
		return nil, err
	}

	opts.SendInitialEvents = watchList
	return uc.resource.Watch(ctx, id.Cluster, gvr, id.Namespace, opts)
}
