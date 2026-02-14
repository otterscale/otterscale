package core

// WatchEventType represents the type of a resource watch event.
// This is a domain-level type that decouples the core layer from
// k8s.io/apimachinery/pkg/watch.EventType.
type WatchEventType string

const (
	WatchEventAdded    WatchEventType = "ADDED"
	WatchEventModified WatchEventType = "MODIFIED"
	WatchEventDeleted  WatchEventType = "DELETED"
	WatchEventBookmark WatchEventType = "BOOKMARK"
	WatchEventError    WatchEventType = "ERROR"
)

// WatchEvent represents a single event from a resource watch stream.
// Object carries the raw Kubernetes resource as a generic map so that
// the domain layer does not depend on unstructured.Unstructured.
type WatchEvent struct {
	Type   WatchEventType
	Object map[string]any
}

// Watcher provides a channel of WatchEvents and a way to stop the
// underlying watch. This replaces the direct use of
// k8s.io/apimachinery/pkg/watch.Interface in the domain layer,
// keeping the core package free of client-go dependencies for watch
// operations.
type Watcher interface {
	// ResultChan returns a channel that receives watch events.
	// The channel is closed when the watch ends or Stop is called.
	ResultChan() <-chan WatchEvent
	// Stop terminates the watch and closes the result channel.
	Stop()
}
