package kubernetes

import (
	"runtime"
	"sync"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/otterscale/otterscale/internal/core"
)

// relayExitTimeout bounds how long a test waits for the relay
// goroutine to react. The bug these tests guard against is a leak, so
// a relay that never exits must fail the test rather than hang it.
const relayExitTimeout = 2 * time.Second

// fakeWatch implements watch.Interface over a channel so the adapter
// can be exercised without a live API server.
type fakeWatch struct {
	ch       chan watch.Event
	stopOnce sync.Once
}

func newFakeWatch(buf int) *fakeWatch {
	return &fakeWatch{ch: make(chan watch.Event, buf)}
}

func (f *fakeWatch) ResultChan() <-chan watch.Event { return f.ch }

func (f *fakeWatch) Stop() { f.stopOnce.Do(func() { close(f.ch) }) }

func podEvent(name string) watch.Event {
	return watch.Event{
		Type: watch.Added,
		Object: &unstructured.Unstructured{Object: map[string]any{
			"kind":     "Pod",
			"metadata": map[string]any{"name": name},
		}},
	}
}

// TestWatcherAdapterStopUnblocksParkedRelay is the regression test for
// the leak: a consumer that walks away mid-stream (canceled request
// context, failed stream send) leaves relay holding an event nobody
// will ever receive. Stopping only the upstream watch does not free it,
// because relay is parked on a send rather than on the range.
//
// The test deliberately never receives from the result channel — doing
// so would unpark relay and hide the bug — so relay's exit is observed
// through the goroutine count instead.
func TestWatcherAdapterStopUnblocksParkedRelay(t *testing.T) {
	before := goroutineFloor(t)

	// An unbuffered upstream channel makes the handoff a rendezvous:
	// this send returns only once relay has taken the event, so relay
	// is provably holding it by the time Stop is called.
	fw := newFakeWatch(0)
	w := newWatcherAdapter(fw)
	fw.ch <- podEvent("pod-a")

	w.Stop()

	waitForGoroutines(t, before, "relay goroutine still parked after Stop")
}

// TestWatcherAdapterStopIsIdempotent guards the stop channel against a
// double close: Stop is called from a handler defer and may also be
// reached on an error path.
func TestWatcherAdapterStopIsIdempotent(t *testing.T) {
	w := newWatcherAdapter(newFakeWatch(0))

	w.Stop()
	w.Stop() // must not panic

	assertClosed(t, w.ResultChan())
}

// TestWatcherAdapterRelaysEventsAndClosesOnUpstreamEnd covers the
// ordinary path: events are converted and forwarded, and the result
// channel closes once the upstream watch ends.
func TestWatcherAdapterRelaysEventsAndClosesOnUpstreamEnd(t *testing.T) {
	fw := newFakeWatch(2)
	fw.ch <- podEvent("pod-a")
	fw.ch <- watch.Event{Type: watch.Error, Object: &metav1.Status{Reason: metav1.StatusReasonExpired}}
	fw.Stop() // closes the upstream channel; buffered events still drain

	w := newWatcherAdapter(fw)

	got := receive(t, w.ResultChan())
	if got.Type != core.WatchEventAdded {
		t.Fatalf("first event type = %q, want %q", got.Type, core.WatchEventAdded)
	}
	metadata, ok := got.Object["metadata"].(map[string]any)
	if !ok || metadata["name"] != "pod-a" {
		t.Fatalf("first event object = %#v, want a pod named pod-a", got.Object)
	}

	got = receive(t, w.ResultChan())
	if got.Type != core.WatchEventError {
		t.Fatalf("second event type = %q, want %q", got.Type, core.WatchEventError)
	}
	if got.Object["reason"] != string(metav1.StatusReasonExpired) {
		t.Fatalf("second event reason = %v, want %q", got.Object["reason"], metav1.StatusReasonExpired)
	}

	assertClosed(t, w.ResultChan())
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// receive takes one event, failing the test rather than blocking the
// suite if none arrives.
func receive(t *testing.T, ch <-chan core.WatchEvent) core.WatchEvent {
	t.Helper()
	select {
	case ev, ok := <-ch:
		if !ok {
			t.Fatal("result channel closed, want an event")
		}
		return ev
	case <-time.After(relayExitTimeout):
		t.Fatal("timed out waiting for an event")
		return core.WatchEvent{}
	}
}

func assertClosed(t *testing.T, ch <-chan core.WatchEvent) {
	t.Helper()
	select {
	case ev, ok := <-ch:
		if ok {
			t.Fatalf("result channel delivered %#v, want it closed", ev)
		}
	case <-time.After(relayExitTimeout):
		t.Fatal("timed out waiting for the result channel to close")
	}
}

// goroutineFloor returns the current goroutine count once it has
// settled, so that goroutines left over from an earlier test are not
// counted against this one.
func goroutineFloor(t *testing.T) int {
	t.Helper()

	deadline := time.Now().Add(relayExitTimeout)
	last := runtime.NumGoroutine()
	for time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
		n := runtime.NumGoroutine()
		if n == last {
			return n
		}
		last = n
	}
	return last
}

// waitForGoroutines waits until the goroutine count drops to want,
// polling because a goroutine's exit is not synchronous with the call
// that released it.
func waitForGoroutines(t *testing.T, want int, msg string) {
	t.Helper()

	deadline := time.Now().Add(relayExitTimeout)
	var n int
	for time.Now().Before(deadline) {
		if n = runtime.NumGoroutine(); n <= want {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("%s: %d goroutines running, want at most %d", msg, n, want)
}
