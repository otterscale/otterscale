package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/otterscale/otterscale/internal/core"
)

// countingReviewer answers from a fixed table and counts how often it is asked.
type countingReviewer struct {
	calls   atomic.Int32
	allowed bool
	err     error
	gate    chan struct{} // when non-nil, every call blocks until it is closed
}

func (r *countingReviewer) Allowed(ctx context.Context, _ string, _ core.AccessAttributes) (bool, error) { //nolint:gocritic // hugeParam: signature fixed by core.AccessReviewer
	r.calls.Add(1)
	if r.gate != nil {
		select {
		case <-r.gate:
		case <-ctx.Done():
			return false, ctx.Err()
		}
	}
	return r.allowed, r.err
}

func userContext(t *testing.T, subject string, groups ...string) context.Context {
	return core.WithUserInfo(t.Context(), core.UserInfo{Subject: subject, Groups: groups})
}

func newTestAccessCache(r core.AccessReviewer, now *time.Time, opts ...AccessReviewOption) *AccessReviewCache {
	return NewAccessReviewCache(r, AccessReviewTTL, append([]AccessReviewOption{
		WithAccessReviewClock(func() time.Time { return *now }),
	}, opts...)...)
}

func TestAccessReviewCache_CachesAllowAndDenyForTTL(t *testing.T) {
	for _, allowed := range []bool{true, false} {
		t.Run(fmt.Sprint(allowed), func(t *testing.T) {
			now := time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)
			inner := &countingReviewer{allowed: allowed}
			c := newTestAccessCache(inner, &now)
			ctx := userContext(t, "alice", "system:authenticated")

			for range 5 {
				got, err := c.Allowed(ctx, "prod", core.ClusterMetricsAccess)
				if err != nil || got != allowed {
					t.Fatalf("Allowed = %v, %v; want %v, nil", got, err, allowed)
				}
			}
			if n := inner.calls.Load(); n != 1 {
				t.Errorf("reviewer asked %d times within the TTL, want 1", n)
			}

			now = now.Add(AccessReviewTTL + time.Second)
			if _, err := c.Allowed(ctx, "prod", core.ClusterMetricsAccess); err != nil {
				t.Fatal(err)
			}
			if n := inner.calls.Load(); n != 2 {
				t.Errorf("reviewer asked %d times after expiry, want 2", n)
			}
		})
	}
}

// TestAccessReviewCache_KeysOnEverythingThatChangesTheAnswer: a cached "yes"
// for one user, group set, attribute or cluster must never answer another.
func TestAccessReviewCache_KeysOnEverythingThatChangesTheAnswer(t *testing.T) {
	now := time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)
	inner := &countingReviewer{allowed: true}
	c := newTestAccessCache(inner, &now)

	ask := func(ctx context.Context, cluster string, attrs core.AccessAttributes) {
		t.Helper()
		if _, err := c.Allowed(ctx, cluster, attrs); err != nil {
			t.Fatal(err)
		}
	}

	ask(userContext(t, "alice", "a", "b"), "prod", core.ClusterMetricsAccess)
	ask(userContext(t, "bob", "a", "b"), "prod", core.ClusterMetricsAccess)             // other user
	ask(userContext(t, "alice", "a"), "prod", core.ClusterMetricsAccess)                // other groups
	ask(userContext(t, "alice", "a", "b"), "staging", core.ClusterMetricsAccess)        // other cluster
	ask(userContext(t, "alice", "a", "b"), "prod", core.NamespaceMetricsAccess("ws-a")) // other attributes
	ask(userContext(t, "alice", "a", "b"), "prod", core.NamespaceMetricsAccess("ws-b")) // other namespace
	if n := inner.calls.Load(); n != 6 {
		t.Errorf("reviewer asked %d times for 6 distinct questions, want 6", n)
	}

	ask(userContext(t, "alice", "b", "a"), "prod", core.ClusterMetricsAccess) // same groups, other order
	if n := inner.calls.Load(); n != 6 {
		t.Errorf("group order changed the cache key (reviewer asked %d times, want 6)", n)
	}
}

func TestAccessReviewCache_DoesNotCacheErrors(t *testing.T) {
	now := time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)
	inner := &countingReviewer{err: errors.New("apiserver down")}
	c := newTestAccessCache(inner, &now)
	ctx := userContext(t, "alice")

	for range 2 {
		if _, err := c.Allowed(ctx, "prod", core.ClusterMetricsAccess); err == nil {
			t.Fatal("Allowed returned no error")
		}
	}
	if n := inner.calls.Load(); n != 2 {
		t.Errorf("reviewer asked %d times, want 2: a failure must not be remembered", n)
	}
}

// TestAccessReviewCache_CollapsesConcurrentMisses: a dashboard page fires
// dozens of queries at once; they must cost one review, not dozens.
func TestAccessReviewCache_CollapsesConcurrentMisses(t *testing.T) {
	now := time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)
	inner := &countingReviewer{allowed: true, gate: make(chan struct{})}
	c := newTestAccessCache(inner, &now)
	ctx := userContext(t, "alice")

	const callers = 32
	var wg sync.WaitGroup
	wg.Add(callers)
	for range callers {
		go func() {
			defer wg.Done()
			if _, err := c.Allowed(ctx, "prod", core.ClusterMetricsAccess); err != nil {
				t.Errorf("Allowed: %v", err)
			}
		}()
	}
	for inner.calls.Load() == 0 {
		time.Sleep(time.Millisecond)
	}
	time.Sleep(20 * time.Millisecond) // let the rest pile up on the flight
	close(inner.gate)
	wg.Wait()

	if n := inner.calls.Load(); n != 1 {
		t.Errorf("reviewer asked %d times for %d concurrent callers, want 1", n, callers)
	}
}

func TestAccessReviewCache_CallerCancelDoesNotFailOthers(t *testing.T) {
	now := time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)
	inner := &countingReviewer{allowed: true, gate: make(chan struct{})}
	c := newTestAccessCache(inner, &now)

	first, cancel := context.WithCancel(userContext(t, "alice"))
	firstDone := make(chan error, 1)
	go func() {
		_, err := c.Allowed(first, "prod", core.ClusterMetricsAccess)
		firstDone <- err
	}()
	for inner.calls.Load() == 0 {
		time.Sleep(time.Millisecond)
	}

	secondDone := make(chan error, 1)
	go func() {
		_, err := c.Allowed(userContext(t, "alice"), "prod", core.ClusterMetricsAccess)
		secondDone <- err
	}()
	time.Sleep(20 * time.Millisecond)
	cancel()
	close(inner.gate)

	<-firstDone
	if err := <-secondDone; err != nil {
		t.Errorf("second caller failed because the first gave up: %v", err)
	}
}

func TestAccessReviewCache_Bounded(t *testing.T) {
	now := time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)
	inner := &countingReviewer{allowed: true}
	c := newTestAccessCache(inner, &now, WithMaxAccessReviewEntries(2))
	ctx := userContext(t, "alice")

	for _, ns := range []string{"a", "b", "c"} {
		if _, err := c.Allowed(ctx, "prod", core.NamespaceMetricsAccess(ns)); err != nil {
			t.Fatal(err)
		}
	}
	if got := c.len(); got != 2 {
		t.Errorf("cache holds %d entries, want the ceiling 2", got)
	}

	// Once the old entries lapse, room is made for new ones.
	now = now.Add(AccessReviewTTL + time.Second)
	if _, err := c.Allowed(ctx, "prod", core.NamespaceMetricsAccess("d")); err != nil {
		t.Fatal(err)
	}
	if got := c.len(); got != 1 {
		t.Errorf("cache holds %d entries after expiry, want 1", got)
	}
}
