package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime/debug"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"

	"github.com/otterscale/otterscale/internal/core"
)

// resourceRepo implements core.ResourceRepo by delegating to the
// Kubernetes dynamic client, accessed through the tunnel.
type resourceRepo struct {
	kubernetes *Kubernetes
}

// NewResourceRepo returns a core.ResourceRepo backed by the Kubernetes
// dynamic API.
func NewResourceRepo(kubernetes *Kubernetes) core.ResourceRepo {
	return &resourceRepo{
		kubernetes: kubernetes,
	}
}

var _ core.ResourceRepo = (*resourceRepo)(nil)

// ---------------------------------------------------------------------------
// CRUD
// ---------------------------------------------------------------------------

// List returns a paged list of resources matching the given options.
func (r *resourceRepo) List(
	ctx context.Context,
	cluster string,
	gvr schema.GroupVersionResource,
	namespace string,
	opts core.ListOptions,
) (*unstructured.UnstructuredList, error) {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return nil, err
	}

	listOpts := metav1.ListOptions{
		LabelSelector: opts.LabelSelector,
		FieldSelector: opts.FieldSelector,
		Limit:         opts.Limit,
		Continue:      opts.Continue,
	}

	result, err := client.Resource(gvr).Namespace(namespace).List(ctx, listOpts)
	return result, wrapK8sError(err)
}

// Get returns a single resource by name.
func (r *resourceRepo) Get(
	ctx context.Context,
	cluster string,
	gvr schema.GroupVersionResource,
	namespace, name string,
) (*unstructured.Unstructured, error) {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return nil, err
	}

	result, err := client.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	return result, wrapK8sError(err)
}

// Create decodes a YAML manifest and creates the resource.
func (r *resourceRepo) Create(
	ctx context.Context,
	cluster string,
	gvr schema.GroupVersionResource,
	namespace string,
	manifest []byte,
) (*unstructured.Unstructured, error) {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return nil, err
	}

	obj, err := fromYAML(manifest)
	if err != nil {
		return nil, err
	}

	result, err := client.Resource(gvr).Namespace(namespace).Create(ctx, obj, metav1.CreateOptions{})
	return result, wrapK8sError(err)
}

// Apply decodes a YAML manifest, converts it to JSON, and performs a
// server-side apply (PATCH with ApplyPatchType). When force is true,
// conflicts are resolved in favor of the caller's field manager.
func (r *resourceRepo) Apply(
	ctx context.Context,
	cluster string,
	gvr schema.GroupVersionResource,
	namespace, name string,
	manifest []byte,
	opts core.ApplyOptions,
) (*unstructured.Unstructured, error) {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return nil, err
	}

	obj, err := fromYAML(manifest)
	if err != nil {
		return nil, err
	}

	data, err := obj.MarshalJSON()
	if err != nil {
		return nil, &core.DomainError{Code: core.ErrorCodeInternal, Message: "marshal manifest to JSON", Cause: err}
	}

	patchOpts := metav1.PatchOptions{
		Force:        &opts.Force,
		FieldManager: opts.FieldManager,
	}

	result, err := client.Resource(gvr).Namespace(namespace).Patch(ctx, name, types.ApplyPatchType, data, patchOpts)
	return result, wrapK8sError(err)
}

// Delete removes a resource.
func (r *resourceRepo) Delete(
	ctx context.Context,
	cluster string,
	gvr schema.GroupVersionResource,
	namespace, name string,
	opts core.DeleteOptions,
) error {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return err
	}

	deleteOpts := metav1.DeleteOptions{
		GracePeriodSeconds: opts.GracePeriodSeconds,
	}

	return wrapK8sError(client.Resource(gvr).Namespace(namespace).Delete(ctx, name, deleteOpts))
}

// ---------------------------------------------------------------------------
// Watch
// ---------------------------------------------------------------------------

// Watch opens a long-lived watch stream for resources matching the
// given options. When SendInitialEvents is true, the server streams
// the current state before switching to change notifications (requires
// Kubernetes >= 1.34).
//
// The returned core.Watcher adapts the Kubernetes watch.Interface to
// the domain-level event model, keeping the core layer free of
// client-go watch types.
func (r *resourceRepo) Watch(
	ctx context.Context,
	cluster string,
	gvr schema.GroupVersionResource,
	namespace string,
	opts core.WatchOptions,
) (core.Watcher, error) {
	client, err := r.watchDynamicClient(ctx, cluster)
	if err != nil {
		return nil, err
	}

	listOpts := metav1.ListOptions{
		LabelSelector:       opts.LabelSelector,
		FieldSelector:       opts.FieldSelector,
		Watch:               true,
		AllowWatchBookmarks: true,
		ResourceVersion:     opts.ResourceVersion,
	}

	if opts.SendInitialEvents {
		listOpts.ResourceVersionMatch = metav1.ResourceVersionMatchNotOlderThan
		listOpts.SendInitialEvents = &opts.SendInitialEvents
	}

	result, err := client.Resource(gvr).Namespace(namespace).Watch(ctx, listOpts)
	if err != nil {
		return nil, wrapK8sError(err)
	}

	return newWatcherAdapter(result), nil
}

// watcherAdapter bridges a Kubernetes watch.Interface to the domain
// core.Watcher interface by converting watch.Event objects into
// core.WatchEvent values with generic map[string]any payloads.
type watcherAdapter struct {
	inner watch.Interface
	ch    chan core.WatchEvent
}

func newWatcherAdapter(inner watch.Interface) *watcherAdapter {
	w := &watcherAdapter{
		inner: inner,
		ch:    make(chan core.WatchEvent),
	}
	go w.relay()
	return w
}

func (w *watcherAdapter) ResultChan() <-chan core.WatchEvent {
	return w.ch
}

func (w *watcherAdapter) Stop() {
	w.inner.Stop()
}

// relay reads from the Kubernetes watch channel and converts events
// to domain WatchEvents. It closes the output channel when the
// upstream channel is closed. A panic recovery is installed to
// prevent a malformed event from crashing the goroutine silently —
// the output channel is still closed via defer so the caller sees
// "watch closed" instead of hanging indefinitely.
func (w *watcherAdapter) relay() {
	defer close(w.ch)
	defer func() {
		if r := recover(); r != nil {
			slog.Error("watch relay panic recovered",
				"panic", r,
				"stack", string(debug.Stack()),
			)
		}
	}()

	for event := range w.inner.ResultChan() {
		domainEvent := core.WatchEvent{
			Type: toCorEventType(event.Type),
		}

		switch obj := event.Object.(type) {
		case *unstructured.Unstructured:
			domainEvent.Object = obj.Object
		case *metav1.Status:
			// Convert Status to a generic map for error events.
			domainEvent.Object = statusToGenericMap(obj)
		}

		w.ch <- domainEvent
	}
}

func toCorEventType(t watch.EventType) core.WatchEventType {
	switch t {
	case watch.Added:
		return core.WatchEventAdded
	case watch.Modified:
		return core.WatchEventModified
	case watch.Deleted:
		return core.WatchEventDeleted
	case watch.Bookmark:
		return core.WatchEventBookmark
	case watch.Error:
		return core.WatchEventError
	default:
		return core.WatchEventError
	}
}

// statusToGenericMap converts a Kubernetes Status object to a generic
// map for inclusion in domain watch events.
func statusToGenericMap(status *metav1.Status) map[string]any {
	// Use JSON round-trip for a simple, accurate conversion.
	data, err := json.Marshal(status)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}
	return m
}

// ---------------------------------------------------------------------------
// Events
// ---------------------------------------------------------------------------

// eventsGVR is the GroupVersionResource for core/v1 Events.
var eventsGVR = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}

// ListEvents returns events matching the given options. This is
// used by DescribeResource to fetch events related to a specific
// resource via involvedObject.uid.
func (r *resourceRepo) ListEvents(
	ctx context.Context,
	cluster, namespace string,
	opts core.ListOptions,
) (*unstructured.UnstructuredList, error) {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return nil, err
	}

	listOpts := metav1.ListOptions{
		LabelSelector: opts.LabelSelector,
		FieldSelector: opts.FieldSelector,
		Limit:         opts.Limit,
		Continue:      opts.Continue,
	}

	result, err := client.Resource(eventsGVR).Namespace(namespace).List(ctx, listOpts)
	return result, wrapK8sError(err)
}

// ---------------------------------------------------------------------------
// Client helpers
// ---------------------------------------------------------------------------

// watchDynamicClient builds a dynamic client with no HTTP timeout, suitable
// for long-lived Watch streams that rely on context cancellation instead.
func (r *resourceRepo) watchDynamicClient(ctx context.Context, cluster string) (*dynamic.DynamicClient, error) {
	config, err := r.kubernetes.impersonationConfig(ctx, cluster)
	if err != nil {
		return nil, err
	}
	config.Timeout = 0
	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, &core.DomainError{Code: core.ErrorCodeInternal, Message: "create dynamic client", Cause: err}
	}
	return dc, nil
}

// dynamicClient builds a fresh impersonated dynamic client for the given cluster.
// A new client is created per request because each request may carry
// different impersonation credentials (user subject + groups). The
// underlying HTTP transport (TCP connections) is cached per-cluster in
// Kubernetes.roundTripper, so the per-request cost is limited to
// allocating the Go-level client wrapper — negligible compared to the
// actual API call latency.
func (r *resourceRepo) dynamicClient(ctx context.Context, cluster string) (*dynamic.DynamicClient, error) {
	config, err := r.kubernetes.impersonationConfig(ctx, cluster)
	if err != nil {
		return nil, err
	}
	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, &core.DomainError{Code: core.ErrorCodeInternal, Message: "create dynamic client", Cause: err}
	}
	return dc, nil
}

// fromYAML decodes a YAML manifest into an Unstructured object.
// Returns a domain validation error if the manifest is invalid.
func fromYAML(manifest []byte) (*unstructured.Unstructured, error) {
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}

	if _, _, err := dec.Decode(manifest, nil, obj); err != nil {
		return nil, &core.ErrInvalidInput{Field: "manifest", Message: fmt.Sprintf("invalid YAML: %s", err)}
	}

	return obj, nil
}
