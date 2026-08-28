package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"

	"github.com/otterscale/otterscale/internal/core"
)

// resourceRepo implements core.ResourceRepo through the Kubernetes dynamic
// client, reached through the tunnel.
type resourceRepo struct {
	kubernetes *Kubernetes
}

func NewResourceRepo(kubernetes *Kubernetes) core.ResourceRepo {
	return &resourceRepo{
		kubernetes: kubernetes,
	}
}

var _ core.ResourceRepo = (*resourceRepo)(nil)

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

// Apply server-side applies the manifest. With opts.Force, conflicts are
// resolved in favor of the caller's field manager.
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

// Update fully replaces the resource (PUT). The manifest must carry
// metadata.name matching the request path, so the dynamic client cannot be
// pointed at a different resource than the one addressed.
func (r *resourceRepo) Update(
	ctx context.Context,
	cluster string,
	gvr schema.GroupVersionResource,
	namespace, name string,
	manifest []byte,
	opts core.UpdateOptions,
) (*unstructured.Unstructured, error) {
	client, err := r.dynamicClient(ctx, cluster)
	if err != nil {
		return nil, err
	}

	obj, err := fromYAML(manifest)
	if err != nil {
		return nil, err
	}
	if obj.GetName() == "" {
		return nil, &core.ErrInvalidInput{Field: "manifest.metadata.name", Message: "must not be empty for update"}
	}
	if obj.GetName() != name {
		return nil, &core.ErrInvalidInput{Field: "manifest.metadata.name", Message: fmt.Sprintf("must match request name %q", name)}
	}

	updateOpts := metav1.UpdateOptions{
		FieldManager: opts.FieldManager,
	}

	result, err := client.Resource(gvr).Namespace(namespace).Update(ctx, obj, updateOpts)
	return result, wrapK8sError(err)
}

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

// Watch opens a long-lived stream. With SendInitialEvents the server streams
// current state first (Kubernetes >= 1.34). The returned core.Watcher adapts
// watch.Interface to the domain event model, keeping core free of client-go
// watch types.
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

	// The API server accepts sendInitialEvents only together with
	// resourceVersionMatch=NotOlderThan, and only when the resource version is
	// unset or "0" — the use case guarantees the latter before setting this.
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

// watcherAdapter converts watch.Event objects into core.WatchEvent values with
// generic map[string]any payloads.
type watcherAdapter struct {
	inner watch.Interface
	ch    chan core.WatchEvent

	// stop is closed by Stop to tell relay the consumer is gone. Stopping the
	// upstream watch is not enough on its own: relay may be parked on a send to
	// ch, which no longer has a receiver, and a closed upstream channel is only
	// observed from the range.
	stop     chan struct{}
	stopOnce sync.Once
}

func newWatcherAdapter(inner watch.Interface) *watcherAdapter {
	w := &watcherAdapter{
		inner: inner,
		ch:    make(chan core.WatchEvent),
		stop:  make(chan struct{}),
	}
	go w.relay()
	return w
}

func (w *watcherAdapter) ResultChan() <-chan core.WatchEvent {
	return w.ch
}

// Stop is safe to call more than once. relay observes the signal and closes the
// result channel on its way out, so the channel is closed whether the watch
// ended upstream or the consumer gave up.
func (w *watcherAdapter) Stop() {
	w.stopOnce.Do(func() { close(w.stop) })
	w.inner.Stop()
}

// relay converts upstream events until the upstream channel closes or Stop is
// called. The panic recovery keeps a malformed event from killing the goroutine
// silently — the deferred close still runs, so the caller sees "watch closed"
// rather than hanging.
//
// Every send selects on stop as well. A consumer that abandons the watch — a
// canceled request context, a failed stream send — stops receiving without
// draining what is in flight, and a send with no escape hatch would park this
// goroutine for the lifetime of the process.
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
			domainEvent.Object = statusToGenericMap(obj)
		}

		select {
		case w.ch <- domainEvent:
		case <-w.stop:
			return
		}
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

func statusToGenericMap(status *metav1.Status) map[string]any {
	// A JSON round-trip is the simplest accurate conversion.
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

var eventsGVR = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "events"}

// ListEvents backs DescribeResource, which filters by involvedObject.uid.
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

// watchDynamicClient drops the HTTP timeout, since a long-lived Watch relies on
// context cancellation instead.
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

// dynamicClient builds a fresh impersonated client per request, because each
// request may carry different impersonation credentials. The underlying HTTP
// transport is cached per cluster in Kubernetes.roundTripper, so only the
// Go-level wrapper is allocated — negligible against the API call latency.
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

// fromYAML returns a domain validation error if the manifest is invalid.
func fromYAML(manifest []byte) (*unstructured.Unstructured, error) {
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}

	if _, _, err := dec.Decode(manifest, nil, obj); err != nil {
		return nil, &core.ErrInvalidInput{Field: "manifest", Message: fmt.Sprintf("invalid YAML: %s", err)}
	}

	return obj, nil
}
