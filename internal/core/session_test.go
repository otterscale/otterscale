package core

import (
	"testing"
	"time"
)

func TestTerminalSizeQueue_SetAndNext(t *testing.T) {
	q := NewTerminalSizeQueue()

	q.Set(80, 24)
	size := q.Next()
	if size == nil {
		t.Fatal("expected non-nil size")
	}
	if size.Width != 80 || size.Height != 24 {
		t.Errorf("got %dx%d, want 80x24", size.Width, size.Height)
	}
}

// TestTerminalSizeQueue_KeepsOnlyLatest pins the coalescing rule: a
// terminal has one size, so a size the consumer has not taken yet is
// stale the moment a newer one arrives.
func TestTerminalSizeQueue_KeepsOnlyLatest(t *testing.T) {
	q := NewTerminalSizeQueue()

	for i := range uint16(5) {
		q.Set(i, i)
	}
	q.Set(99, 99)

	size := q.Next()
	if size == nil {
		t.Fatal("expected non-nil size")
	}
	if size.Width != 99 || size.Height != 99 {
		t.Errorf("got %dx%d, want 99x99", size.Width, size.Height)
	}
}

// TestTerminalSizeQueue_SetNeverBlocks is the regression test for the deadlock.
// Set runs under the lock Close needs, so a Set that blocks takes Close with it
// — and Close is the first deferred call in the exec goroutine, wedging the
// whole session: pipes stay open, Done never fires, the reaper cannot collect.
//
// The old implementation dropped the oldest entry with a blocking receive,
// which parked whenever a consumer emptied the channel between the failed send
// and that receive. Driving Set and Next concurrently hits that interleaving.
func TestTerminalSizeQueue_SetNeverBlocks(t *testing.T) {
	q := NewTerminalSizeQueue()

	const rounds = 2000

	consumerDone := make(chan struct{})
	go func() {
		defer close(consumerDone)
		for q.Next() != nil { //nolint:revive // draining until the queue closes is the point
		}
	}()

	producerDone := make(chan struct{})
	go func() {
		defer close(producerDone)
		for i := range uint16(rounds) {
			q.Set(i, i)
		}
	}()

	select {
	case <-producerDone:
	case <-time.After(sessionTimeout):
		t.Fatal("Set blocked: a resize parked while holding the queue lock")
	}

	closed := make(chan struct{})
	go func() {
		defer close(closed)
		q.Close()
	}()

	select {
	case <-closed:
	case <-time.After(sessionTimeout):
		t.Fatal("Close blocked: Set is still holding the queue lock")
	}

	select {
	case <-consumerDone:
	case <-time.After(sessionTimeout):
		t.Fatal("Next did not return nil after Close")
	}
}

func TestTerminalSizeQueue_Close(t *testing.T) {
	q := NewTerminalSizeQueue()

	q.Close()

	size := q.Next()
	if size != nil {
		t.Errorf("expected nil after close, got %v", size)
	}

	q.Close()
}

func TestTerminalSizeQueue_SetAfterClose(_ *testing.T) {
	q := NewTerminalSizeQueue()
	q.Close()

	q.Set(80, 24)
}

func TestSessionStore_ExecCRUD(t *testing.T) {
	store := NewSessionStore()
	done := make(chan struct{})
	close(done)

	sess := &ExecSession{
		ID:   "exec-1",
		Done: done,
	}

	if err := store.PutExec(sess); err != nil {
		t.Fatalf("PutExec: %v", err)
	}

	got, ok := store.GetExec("exec-1")
	if !ok {
		t.Fatal("expected to find exec session")
	}
	if got.ID != "exec-1" {
		t.Errorf("got ID %q, want %q", got.ID, "exec-1")
	}

	_, ok = store.GetExec("nonexistent")
	if ok {
		t.Error("expected not to find nonexistent session")
	}

	removed := store.RemoveExec("exec-1")
	if removed == nil {
		t.Fatal("expected RemoveExec to return the session")
	}
	_, ok = store.GetExec("exec-1")
	if ok {
		t.Error("expected session to be removed")
	}

	if store.RemoveExec("exec-1") != nil {
		t.Error("expected nil for already-removed session")
	}
}

func TestSessionStore_PortForwardCRUD(t *testing.T) {
	store := NewSessionStore()
	done := make(chan struct{})
	close(done)

	sess := &PortForwardSession{
		ID:   "pf-1",
		Done: done,
	}

	if err := store.PutPortForward(sess); err != nil {
		t.Fatalf("PutPortForward: %v", err)
	}

	got, ok := store.GetPortForward("pf-1")
	if !ok {
		t.Fatal("expected to find port-forward session")
	}
	if got.ID != "pf-1" {
		t.Errorf("got ID %q, want %q", got.ID, "pf-1")
	}

	removed := store.RemovePortForward("pf-1")
	if removed == nil {
		t.Fatal("expected RemovePortForward to return the session")
	}
	_, ok = store.GetPortForward("pf-1")
	if ok {
		t.Error("expected session to be removed")
	}

	if store.RemovePortForward("pf-1") != nil {
		t.Error("expected nil for already-removed session")
	}
}

func TestSessionStore_ReapStaleSessions(t *testing.T) {
	store := NewSessionStore()

	// Stale: Done already closed.
	execDone := make(chan struct{})
	close(execDone)

	if err := store.PutExec(&ExecSession{
		ID:     "stale-exec",
		Done:   execDone,
		Cancel: func() {},
		Stdin:  &nopCloser{},
	}); err != nil {
		t.Fatalf("PutExec stale: %v", err)
	}

	// Live: Done still open.
	liveDone := make(chan struct{})
	if err := store.PutExec(&ExecSession{
		ID:     "live-exec",
		Done:   liveDone,
		Cancel: func() {},
		Stdin:  &nopCloser{},
	}); err != nil {
		t.Fatalf("PutExec live: %v", err)
	}

	pfDone := make(chan struct{})
	close(pfDone)

	if err := store.PutPortForward(&PortForwardSession{
		ID:     "stale-pf",
		Done:   pfDone,
		Cancel: func() {},
		Writer: &nopCloser{},
	}); err != nil {
		t.Fatalf("PutPortForward stale: %v", err)
	}

	reaped := store.ReapStaleSessions()
	if reaped != 2 {
		t.Errorf("expected 2 reaped sessions, got %d", reaped)
	}

	if _, ok := store.GetExec("stale-exec"); ok {
		t.Error("stale exec session should have been reaped")
	}
	if _, ok := store.GetPortForward("stale-pf"); ok {
		t.Error("stale port-forward session should have been reaped")
	}

	if _, ok := store.GetExec("live-exec"); !ok {
		t.Error("live exec session should still exist")
	}
}

type nopCloser struct{}

func (n *nopCloser) Write(p []byte) (int, error) { return len(p), nil }
func (n *nopCloser) Close() error                { return nil }
