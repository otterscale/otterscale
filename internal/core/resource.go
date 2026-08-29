package core

// ADR: Kubernetes types in the domain layer
//
// This file imports k8s.io packages directly into core. Strict DDD would wrap
// them in DTOs, but otterscale's core business *is* Kubernetes resource
// management: GVR, Unstructured, APIResourceList, and OpenAPI Schema belong to
// the domain's ubiquitous language, and wrapping them would add a translation
// layer at every boundary for a domain that stays structurally identical.
// Revisit if the project ever supports non-Kubernetes backends.

import (
	"context"
	"fmt"
	"log/slog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

// DiscoveryClient abstracts Kubernetes API discovery so the use-case layer can
// validate resources and fetch schemas without a concrete client.
type DiscoveryClient interface {
	// LookupResource validates that a group/version/resource triple exists.
	LookupResource(ctx context.Context, cluster, group, version, resource, subresource string) (schema.GroupVersionResource, error)
	ServerResources(ctx context.Context, cluster string) ([]*metav1.APIResourceList, error)
	// ResolveGroupVersionSchemas returns the schemas for every kind in a
	// group/version. Kubernetes serves one OpenAPI document per group/version,
	// so callers should cache the result rather than refetch per GVK.
	ResolveGroupVersionSchemas(ctx context.Context, cluster, group, version string) (map[string]*spec.Schema, error)
	ServerVersion(ctx context.Context, cluster string) (*version.Info, error)
	// SupportsWatchList reports whether the cluster has the WatchList
	// streaming feature (Kubernetes >= 1.34).
	SupportsWatchList(ctx context.Context, cluster string) (bool, error)
}

// ResourceRepo abstracts resource CRUD and watch through the dynamic client.
// Every method takes a cluster name so the implementation can route through
// the right tunnel.
type ResourceRepo interface {
	List(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace string, opts ListOptions,
	) (*unstructured.UnstructuredList, error)

	Get(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace, name string,
	) (*unstructured.Unstructured, error)

	// Create decodes a YAML manifest.
	Create(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace string, manifest []byte,
	) (*unstructured.Unstructured, error)

	// Apply decodes a YAML manifest and server-side applies it (PATCH with
	// ApplyPatchType).
	Apply(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace, name string, manifest []byte, opts ApplyOptions,
	) (*unstructured.Unstructured, error)

	// Update decodes a YAML manifest and fully replaces the resource (PUT).
	Update(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace, name string, manifest []byte, opts UpdateOptions,
	) (*unstructured.Unstructured, error)

	Delete(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace, name string, opts DeleteOptions,
	) error

	Watch(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace string, opts WatchOptions,
	) (Watcher, error)

	// ListEvents backs DescribeResource, which filters by involvedObject.uid.
	ListEvents(ctx context.Context, cluster, namespace string, opts ListOptions) (*unstructured.UnstructuredList, error)
}

// ListOptions mirrors the commonly used fields of metav1.ListOptions.
type ListOptions struct {
	LabelSelector string
	FieldSelector string
	Limit         int64
	Continue      string
}

// ApplyOptions mirrors the commonly used fields of metav1.PatchOptions.
type ApplyOptions struct {
	Force        bool
	FieldManager string
}

// UpdateOptions mirrors the commonly used fields of metav1.UpdateOptions.
type UpdateOptions struct {
	FieldManager string
}

// DeleteOptions mirrors the commonly used fields of metav1.DeleteOptions.
type DeleteOptions struct {
	GracePeriodSeconds *int64
}

// WatchOptions mirrors the metav1.ListOptions fields relevant to a watch.
type WatchOptions struct {
	LabelSelector     string
	FieldSelector     string
	ResourceVersion   string
	SendInitialEvents bool
}

// SchemaResolver resolves OpenAPI schemas for GVKs, decoupling the use-case
// layer from the caching infrastructure. Implementations may cache results and
// deduplicate concurrent requests.
type SchemaResolver interface {
	ResolveSchema(ctx context.Context, cluster, group, version, kind string) (*spec.Schema, error)
}

// ResourceIdentifier identifies a resource type, and optionally one instance,
// across clusters. It replaces long positional parameter lists in use-case
// methods with a single value object.
type ResourceIdentifier struct {
	Cluster     string
	Group       string
	Version     string
	Resource    string
	SubResource string
	Namespace   string
	Name        string
}

func (id *ResourceIdentifier) lookupGVR(ctx context.Context, dc DiscoveryClient) (schema.GroupVersionResource, error) {
	return dc.LookupResource(ctx, id.Cluster, id.Group, id.Version, id.Resource, id.SubResource)
}

// ResourceUseCase manages Kubernetes resources across clusters. It validates
// GVRs via the DiscoveryClient and resolves schemas through the injected
// SchemaResolver.
type ResourceUseCase struct {
	discovery      DiscoveryClient
	resource       ResourceRepo
	schemaResolver SchemaResolver
}

func NewResourceUseCase(discovery DiscoveryClient, resource ResourceRepo, schemaResolver SchemaResolver) *ResourceUseCase {
	return &ResourceUseCase{
		discovery:      discovery,
		resource:       resource,
		schemaResolver: schemaResolver,
	}
}

func (uc *ResourceUseCase) ServerResources(ctx context.Context, cluster string) ([]*metav1.APIResourceList, error) {
	return uc.discovery.ServerResources(ctx, cluster)
}

func (uc *ResourceUseCase) ResolveSchema(
	ctx context.Context,
	cluster, group, version, kind string,
) (*spec.Schema, error) {
	return uc.schemaResolver.ResolveSchema(ctx, cluster, group, version, kind)
}

func (uc *ResourceUseCase) ListResources(
	ctx context.Context,
	id *ResourceIdentifier,
	opts ListOptions,
) (*unstructured.UnstructuredList, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	return uc.resource.List(ctx, id.Cluster, gvr, id.Namespace, opts)
}

func (uc *ResourceUseCase) GetResource(
	ctx context.Context,
	id *ResourceIdentifier,
) (*unstructured.Unstructured, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	return uc.resource.Get(ctx, id.Cluster, gvr, id.Namespace, id.Name)
}

// DescribeResource is the backend equivalent of `kubectl describe`: it fetches
// the resource and the events referencing its UID.
func (uc *ResourceUseCase) DescribeResource(
	ctx context.Context,
	id *ResourceIdentifier,
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
		// Events are supplementary, so return the resource anyway. Log it, or a
		// permission problem is indistinguishable from having no events.
		slog.Warn("failed to list events for describe",
			fieldCluster, id.Cluster,
			"namespace", id.Namespace,
			fieldName, id.Name,
			"error", err,
		)
		return obj, &unstructured.UnstructuredList{}, nil
	}

	return obj, events, nil
}

func (uc *ResourceUseCase) CreateResource(
	ctx context.Context,
	id *ResourceIdentifier,
	manifest []byte,
) (*unstructured.Unstructured, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	return uc.resource.Create(ctx, id.Cluster, gvr, id.Namespace, manifest)
}

func (uc *ResourceUseCase) ApplyResource(
	ctx context.Context,
	id *ResourceIdentifier,
	manifest []byte,
	opts ApplyOptions,
) (*unstructured.Unstructured, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	return uc.resource.Apply(ctx, id.Cluster, gvr, id.Namespace, id.Name, manifest, opts)
}

func (uc *ResourceUseCase) UpdateResource(
	ctx context.Context,
	id *ResourceIdentifier,
	manifest []byte,
	opts UpdateOptions,
) (*unstructured.Unstructured, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	return uc.resource.Update(ctx, id.Cluster, gvr, id.Namespace, id.Name, manifest, opts)
}

func (uc *ResourceUseCase) DeleteResource(
	ctx context.Context,
	id *ResourceIdentifier,
	opts DeleteOptions,
) error {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return err
	}

	return uc.resource.Delete(ctx, id.Cluster, gvr, id.Namespace, id.Name, opts)
}

// WatchResource opens a long-lived watch stream. A fresh watch on a cluster
// with the WatchList feature streams current state before switching to change
// notifications; see wantsInitialEvents for why a resumed watch does not.
func (uc *ResourceUseCase) WatchResource(
	ctx context.Context,
	id *ResourceIdentifier,
	opts WatchOptions,
) (Watcher, error) {
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	opts.SendInitialEvents = uc.wantsInitialEvents(ctx, id.Cluster, opts.ResourceVersion)

	return uc.resource.Watch(ctx, id.Cluster, gvr, id.Namespace, opts)
}

// wantsInitialEvents reports whether the watch should ask for current state
// before change notifications.
//
// A caller supplying a resource version is resuming an earlier watch —
// typically from a BOOKMARK event — and wants only what changed since. Asking
// for initial events there is not merely redundant: the API server rejects
// sendInitialEvents unless the resource version is unset or "0".
//
// Discovery failures degrade to a plain watch rather than failing the call.
// Starting from "now" is still correct, and it keeps watches working when the
// version endpoint is briefly unavailable or reports an unparseable version.
func (uc *ResourceUseCase) wantsInitialEvents(ctx context.Context, cluster, resourceVersion string) bool {
	if resourceVersion != "" && resourceVersion != "0" {
		return false
	}

	supported, err := uc.discovery.SupportsWatchList(ctx, cluster)
	if err != nil {
		slog.Warn("watch-list support unknown, falling back to a plain watch",
			fieldCluster, cluster,
			"error", err,
		)
		return false
	}
	return supported
}
