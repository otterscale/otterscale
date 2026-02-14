package core

import (
	"testing"
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

func TestTerminalSizeQueue_OverflowDropsOldest(t *testing.T) {
	q := NewTerminalSizeQueue()

	// Fill the buffer (capacity 4).
	for i := uint16(0); i < 4; i++ {
		q.Set(i, i)
	}

	// Push one more, which should drop the oldest (0x0).
	q.Set(99, 99)

	// The first dequeued element should be 1x1 (0x0 was dropped).
	size := q.Next()
	if size == nil {
		t.Fatal("expected non-nil size")
	}
	if size.Width != 1 || size.Height != 1 {
		t.Errorf("got %dx%d, want 1x1", size.Width, size.Height)
	}
}

func TestTerminalSizeQueue_Close(t *testing.T) {
	q := NewTerminalSizeQueue()

	q.Close()

	// Next should return nil after close.
	size := q.Next()
	if size != nil {
		t.Errorf("expected nil after close, got %v", size)
	}

	// Double close should not panic.
	q.Close()
}

func TestTerminalSizeQueue_SetAfterClose(t *testing.T) {
	q := NewTerminalSizeQueue()
	q.Close()

	// Should not panic.
	q.Set(80, 24)
}

func TestSessionStore_ExecCRUD(t *testing.T) {
	store := NewSessionStore()
	done := make(chan error, 1)
	done <- nil

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

	// RemoveExec on a non-existent session should return nil.
	if store.RemoveExec("exec-1") != nil {
		t.Error("expected nil for already-removed session")
	}
}

func TestSessionStore_PortForwardCRUD(t *testing.T) {
	store := NewSessionStore()
	done := make(chan error, 1)
	done <- nil

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

	// RemovePortForward on a non-existent session should return nil.
	if store.RemovePortForward("pf-1") != nil {
		t.Error("expected nil for already-removed session")
	}
}

func TestSessionStore_ReapStaleSessions(t *testing.T) {
	store := NewSessionStore()

	// Create a "stale" exec session (Done already received a value).
	execDone := make(chan error, 1)
	execDone <- nil
	close(execDone)

	if err := store.PutExec(&ExecSession{
		ID:     "stale-exec",
		Done:   execDone,
		Cancel: func() {},
		Stdin:  &nopCloser{},
	}); err != nil {
		t.Fatalf("PutExec stale: %v", err)
	}

	// Create a "live" exec session (Done has no value yet).
	liveDone := make(chan error, 1)
	if err := store.PutExec(&ExecSession{
		ID:     "live-exec",
		Done:   liveDone,
		Cancel: func() {},
		Stdin:  &nopCloser{},
	}); err != nil {
		t.Fatalf("PutExec live: %v", err)
	}

	// Create a "stale" port-forward session.
	pfDone := make(chan error, 1)
	pfDone <- nil
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

	// Stale sessions should be gone.
	if _, ok := store.GetExec("stale-exec"); ok {
		t.Error("stale exec session should have been reaped")
	}
	if _, ok := store.GetPortForward("stale-pf"); ok {
		t.Error("stale port-forward session should have been reaped")
	}

	// Live session should remain.
	if _, ok := store.GetExec("live-exec"); !ok {
		t.Error("live exec session should still exist")
	}
}

// nopCloser is a no-op io.WriteCloser for tests.
type nopCloser struct{}

func (n *nopCloser) Write(p []byte) (int, error) { return len(p), nil }
func (n *nopCloser) Close() error                { return nil }
