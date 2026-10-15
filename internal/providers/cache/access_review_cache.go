package cache

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/otterscale/otterscale/internal/core"
)

// AccessReviewTTL is how long one answer about one caller stands. A dashboard
// page sends dozens of metric queries; this turns them into one review per
// (user, attribute) per half minute, and bounds how long a revoked
// RoleBinding keeps working through the hub.
const AccessReviewTTL = 30 * time.Second

// defaultMaxAccessReviewEntries bounds the cache. Keys include the namespace a
// caller asks about, so a caller cycling through made-up namespaces would
// otherwise grow it without limit.
const defaultMaxAccessReviewEntries = 10000

// accessReviewTimeout bounds a review, which runs under context.WithoutCancel
// so one caller giving up does not fail everyone waiting on the same flight.
const accessReviewTimeout = 10 * time.Second

// AccessReviewCache remembers SelfSubjectAccessReview answers, allowed and
// denied alike, for a TTL, and collapses concurrent identical questions into
// one. Failures are not remembered: no answer is not a "no" to keep. They are
// logged here, once per review, so an outage costs a line per review rather
// than one per query that waited on it.
//
// Expiry needs no background sweeper: a lookup ignores lapsed entries and an
// insert at the ceiling clears them first (the AgentValuesStore pattern).
type AccessReviewCache struct {
	reviewer   core.AccessReviewer
	ttl        time.Duration
	now        func() time.Time
	maxEntries int

	mu      sync.Mutex
	entries map[string]accessReviewEntry
	flights singleflight.Group
}

type accessReviewEntry struct {
	allowed   bool
	expiresAt time.Time
}

var _ core.AccessReviewer = (*AccessReviewCache)(nil)

type AccessReviewOption func(*AccessReviewCache)

// WithAccessReviewClock injects a time source for deterministic tests.
func WithAccessReviewClock(now func() time.Time) AccessReviewOption {
	return func(c *AccessReviewCache) { c.now = now }
}

func WithMaxAccessReviewEntries(n int) AccessReviewOption {
	return func(c *AccessReviewCache) {
		if n > 0 {
			c.maxEntries = n
		}
	}
}

func NewAccessReviewCache(reviewer core.AccessReviewer, ttl time.Duration, opts ...AccessReviewOption) *AccessReviewCache {
	c := &AccessReviewCache{
		reviewer:   reviewer,
		ttl:        ttl,
		now:        time.Now,
		maxEntries: defaultMaxAccessReviewEntries,
		entries:    make(map[string]accessReviewEntry),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *AccessReviewCache) Allowed(ctx context.Context, cluster string, attrs core.AccessAttributes) (bool, error) { //nolint:gocritic // hugeParam: signature fixed by core.AccessReviewer
	user, ok := core.UserInfoFromContext(ctx)
	if !ok || user.Subject == "" {
		// Nothing to key on. The reviewer refuses such a context itself.
		return c.reviewer.Allowed(ctx, cluster, attrs)
	}
	key := accessReviewKey(cluster, user, attrs)

	if allowed, ok := c.lookup(key); ok {
		return allowed, nil
	}

	v, err, _ := c.flights.Do(key, func() (any, error) {
		// A flight that landed between the lookup above and this one being
		// started has already stored the answer; asking again would be a
		// second review of the same question inside the TTL.
		if allowed, ok := c.lookup(key); ok {
			return allowed, nil
		}

		reviewCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), accessReviewTimeout)
		defer cancel()

		allowed, err := c.reviewer.Allowed(reviewCtx, cluster, attrs)
		if err != nil {
			slog.Warn("access review failed", "cluster", cluster, "resource", attrs.Resource, "namespace", attrs.Namespace, "error", err)
			return false, err
		}
		c.store(key, allowed)
		return allowed, nil
	})
	if err != nil {
		return false, err
	}
	allowed, _ := v.(bool)
	return allowed, nil
}

// lookup reports a stored answer that has not lapsed.
func (c *AccessReviewCache) lookup(key string) (allowed, ok bool) {
	c.mu.Lock()
	entry, hit := c.entries[key]
	c.mu.Unlock()
	return entry.allowed, hit && c.now().Before(entry.expiresAt)
}

func (c *AccessReviewCache) store(key string, allowed bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.entries) >= c.maxEntries {
		now := c.now()
		for k, e := range c.entries {
			if !now.Before(e.expiresAt) {
				delete(c.entries, k)
			}
		}
	}
	// Still full: answer without remembering rather than evict live entries.
	if len(c.entries) >= c.maxEntries {
		return
	}
	c.entries[key] = accessReviewEntry{allowed: allowed, expiresAt: c.now().Add(c.ttl)}
}

// len is for tests.
func (c *AccessReviewCache) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}

// accessReviewKey quotes every part, so no subject, group or namespace can be
// crafted to collide with another key. Groups are sorted: the same set in
// another order is the same identity.
func accessReviewKey(cluster string, user core.UserInfo, attrs core.AccessAttributes) string { //nolint:gocritic // hugeParam: attrs is a small value type, copied once per call
	groups := slices.Clone(user.Groups)
	slices.Sort(groups)
	return fmt.Sprintf("%q %q %q %q", cluster, user.Subject, groups, attrs)
}
