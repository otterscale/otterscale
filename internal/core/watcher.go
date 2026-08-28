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
//
// Stop is the only signal an implementation gets that the consumer has
// gone away, so it must be enough to release one that is mid-delivery.
// A consumer stops receiving as soon as it gives up — a canceled
// request context, a failed stream send — without draining what is
// already in flight.
type Watcher interface {
	// ResultChan returns a channel that receives watch events.
	// The channel is closed when the watch ends or Stop is called.
	ResultChan() <-chan WatchEvent
	// Stop terminates the watch and closes the result channel.
	// Implementations must unblock a producer parked on a send that
	// nobody will receive, and must tolerate being called more than
	// once.
	Stop()
}
