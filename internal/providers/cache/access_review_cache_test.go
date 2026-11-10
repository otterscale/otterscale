package cache

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
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

// TestAccessReviewCache_NoSecondReviewAfterAFlightLands: a caller that missed
// just before a flight stored its answer, and reaches singleflight only after
// that flight is gone, must take the stored answer rather than start a second
// review inside the TTL: at most one review per user per attribute per 30 s.
func TestAccessReviewCache_NoSecondReviewAfterAFlightLands(t *testing.T) {
	start := time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)
	inner := &countingReviewer{allowed: true}
	clock := &steppingClock{now: start}
	c := NewAccessReviewCache(inner, AccessReviewTTL, WithAccessReviewClock(clock.Now))
	ctx := userContext(t, "alice")

	// An answer that has lapsed: the next lookups find an entry but miss.
	if _, err := c.Allowed(ctx, "prod", core.ClusterMetricsAccess); err != nil {
		t.Fatal(err)
	}
	clock.set(start.Add(AccessReviewTTL + time.Second))

	// Caller B reads the lapsed entry, then stalls in the clock before it
	// decides; meanwhile caller A misses, reviews, stores and finishes.
	clock.pauseNext()
	bDone := make(chan error, 1)
	go func() {
		_, err := c.Allowed(ctx, "prod", core.ClusterMetricsAccess)
		bDone <- err
	}()
	clock.waitPaused()
	if _, err := c.Allowed(ctx, "prod", core.ClusterMetricsAccess); err != nil {
		t.Fatal(err)
	}
	clock.resume()
	if err := <-bDone; err != nil {
		t.Fatal(err)
	}

	if n := inner.calls.Load(); n != 2 {
		t.Errorf("reviewer asked %d times, want 2: one before expiry, one after — B must reuse A's answer", n)
	}
}

// steppingClock is a settable clock that can hold the next caller of Now until
// resumed, to stop a goroutine between reading the cache and deciding.
type steppingClock struct {
	mu      sync.Mutex
	now     time.Time
	pause   bool
	paused  chan struct{}
	release chan struct{}
}

func (c *steppingClock) Now() time.Time {
	c.mu.Lock()
	if !c.pause {
		defer c.mu.Unlock()
		return c.now
	}
	c.pause = false
	paused, release, now := c.paused, c.release, c.now
	c.mu.Unlock()
	close(paused)
	<-release
	return now
}

func (c *steppingClock) set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t
}

func (c *steppingClock) pauseNext() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pause, c.paused, c.release = true, make(chan struct{}), make(chan struct{})
}

func (c *steppingClock) waitPaused() {
	c.mu.Lock()
	paused := c.paused
	c.mu.Unlock()
	<-paused
}

func (c *steppingClock) resume() {
	c.mu.Lock()
	release := c.release
	c.mu.Unlock()
	close(release)
}

// TestAccessReviewCache_LogsAFailureOncePerReview: while a cluster's apiserver
// is down, a page load is dozens of queries; the log gets one line per review
// that failed, not one per query.
func TestAccessReviewCache_LogsAFailureOncePerReview(t *testing.T) {
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })

	now := time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)
	inner := &countingReviewer{err: errors.New("apiserver down"), gate: make(chan struct{})}
	c := newTestAccessCache(inner, &now)
	ctx := userContext(t, "alice")

	const callers = 16
	var wg sync.WaitGroup
	wg.Add(callers)
	for range callers {
		go func() {
			defer wg.Done()
			if _, err := c.Allowed(ctx, "prod", core.NamespaceMetricsAccess("ws-a")); err == nil {
				t.Error("Allowed returned no error")
			}
		}()
	}
	for inner.calls.Load() == 0 {
		time.Sleep(time.Millisecond)
	}
	time.Sleep(20 * time.Millisecond) // let the rest pile up on the flight
	close(inner.gate)
	wg.Wait()

	lines := strings.FieldsFunc(logs.String(), func(r rune) bool { return r == '\n' })
	if n := int(inner.calls.Load()); len(lines) != n || n >= callers {
		t.Fatalf("%d log lines for %d reviews and %d callers, want one per review:\n%s", len(lines), n, callers, logs.String())
	}
	for _, want := range []string{"level=WARN", "cluster=prod", "resource=pods", "namespace=ws-a", `error="apiserver down"`} {
		if !strings.Contains(lines[0], want) {
			t.Errorf("log line %q lacks %s", lines[0], want)
		}
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
