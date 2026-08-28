package core

// WatchEventType is the domain-level stand-in for
// k8s.io/apimachinery/pkg/watch.EventType.
type WatchEventType string

const (
	WatchEventAdded    WatchEventType = "ADDED"
	WatchEventModified WatchEventType = "MODIFIED"
	WatchEventDeleted  WatchEventType = "DELETED"
	WatchEventBookmark WatchEventType = "BOOKMARK"
	WatchEventError    WatchEventType = "ERROR"
)

// WatchEvent carries the raw resource in Object as a generic map, so the domain
// layer does not depend on unstructured.Unstructured.
type WatchEvent struct {
	Type   WatchEventType
	Object map[string]any
}

// Watcher is the domain-level stand-in for
// k8s.io/apimachinery/pkg/watch.Interface.
//
// Stop is the only signal an implementation gets that the consumer has gone
// away, so it must be enough to release one mid-delivery. A consumer stops
// receiving as soon as it gives up — a canceled request context, a failed
// stream send — without draining what is already in flight.
type Watcher interface {
	// ResultChan is closed when the watch ends or Stop is called.
	ResultChan() <-chan WatchEvent
	// Stop must unblock a producer parked on a send nobody will receive, and
	// must tolerate being called more than once.
	Stop()
}
