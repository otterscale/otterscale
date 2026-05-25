// Package cache provides TTL-based caching infrastructure for
// Kubernetes discovery data. It lives in the providers layer because
// caching is an infrastructure concern — the domain layer
// (internal/core) only defines the SchemaResolver interface.
package cache

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
	"k8s.io/apiserver/pkg/cel/openapi/resolver"
	"k8s.io/kube-openapi/pkg/validation/spec"

	"github.com/otterscale/otterscale/internal/core"
)

// DefaultTTL is the default TTL for cached OpenAPI schemas.
// Exported so that the DI layer can use it when constructing a
// DiscoveryCache.
const DefaultTTL = 10 * time.Minute

// defaultMaxGVEntries is the upper bound on the number of cached
// group/version entries. Each entry holds every kind in that
// group/version, so the bound on actual kinds cached is much higher
// than this number suggests.
const defaultMaxGVEntries = 1024

// ttlJitterFraction is the maximum random jitter applied to an entry's
// TTL, as a fraction of the configured TTL. Spreading expirations
// prevents a thundering-herd cache stampede when many entries are
// populated in the same burst.
const ttlJitterFraction = 0.2

// DiscoveryCache provides TTL-based caching with singleflight
// deduplication for OpenAPI schemas. It implements
// core.SchemaResolver and core.CacheEvictor, and reduces redundant
// discovery API calls when multiple concurrent requests target the
// same cluster.
//
// Entries are keyed at group/version granularity because the
// Kubernetes OpenAPI v3 endpoint returns one document per GV. Caching
// per-GVK would re-download the same document for every kind; the GV
// scheme amortizes a single fetch across all kinds in that GV and
// lets singleflight deduplicate concurrent misses for different kinds
// of the same GV.
type DiscoveryCache struct {
	discovery     core.DiscoveryClient
	ttl           time.Duration
	now           func() time.Time
	maxGVEntries  int
	jitterSampler func() float64

	mu        sync.RWMutex
	gvCache   map[string]*gvCacheEntry
	gvFlights singleflight.Group
}

// gvCacheEntry pairs the kind→schema map for a single group/version
// with its (jittered) expiration time.
type gvCacheEntry struct {
	schemas   map[string]*spec.Schema
	expiresAt time.Time
}

// singleflightFetchTimeout is the maximum time a cache-miss fetch is
// allowed to run. It uses context.WithoutCancel so that a single
// caller's cancellation does not fail all singleflight waiters.
const singleflightFetchTimeout = 30 * time.Second

// Option configures a DiscoveryCache at construction time.
type Option func(*DiscoveryCache)

// WithClock injects a custom time source for deterministic testing.
// When not set, time.Now is used.
func WithClock(now func() time.Time) Option {
	return func(c *DiscoveryCache) {
		c.now = now
	}
}

// WithMaxGVEntries overrides the default upper bound on cached
// group/version entries.
func WithMaxGVEntries(n int) Option {
	return func(c *DiscoveryCache) {
		if n > 0 {
			c.maxGVEntries = n
		}
	}
}

// WithJitterSampler injects a custom [0,1) sampler used to compute
// TTL jitter. Intended for deterministic testing; production callers
// can rely on the default math/rand/v2-backed sampler.
func WithJitterSampler(sample func() float64) Option {
	return func(c *DiscoveryCache) {
		c.jitterSampler = sample
	}
}

// NewDiscoveryCache returns a DiscoveryCache that wraps the given
// DiscoveryClient and caches results for the specified TTL.
func NewDiscoveryCache(discovery core.DiscoveryClient, ttl time.Duration, opts ...Option) *DiscoveryCache {
	c := &DiscoveryCache{
		discovery:     discovery,
		ttl:           ttl,
		now:           time.Now,
		maxGVEntries:  defaultMaxGVEntries,
		jitterSampler: rand.Float64,
		gvCache:       make(map[string]*gvCacheEntry),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// ResolveSchema fetches the OpenAPI schema for the given GVK. Results
// are cached per group/version for the configured TTL and concurrent
// requests for the same group/version are deduplicated via
// singleflight.
func (c *DiscoveryCache) ResolveSchema(
	ctx context.Context,
	cluster, group, version, kind string,
) (*spec.Schema, error) {
	key := c.gvCacheKey(cluster, group, version)

	c.mu.RLock()
	entry, ok := c.gvCache[key]
	c.mu.RUnlock()

	if ok && c.now().Before(entry.expiresAt) {
		return lookupKind(entry.schemas, group, version, kind)
	}

	v, err, _ := c.gvFlights.Do(key, func() (any, error) {
		// Use a non-cancellable context with its own timeout so that
		// a single caller's cancellation does not fail all waiters
		// sharing this singleflight key.
		fetchCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), singleflightFetchTimeout)
		defer cancel()

		schemas, err := c.discovery.ResolveGroupVersionSchemas(fetchCtx, cluster, group, version)
		if err != nil {
			return nil, err
		}

		c.mu.Lock()
		if len(c.gvCache) >= c.maxGVEntries {
			c.evictExpiredGVs()
		}
		if len(c.gvCache) < c.maxGVEntries {
			c.gvCache[key] = &gvCacheEntry{
				schemas:   schemas,
				expiresAt: c.now().Add(c.jitteredTTL()),
			}
		}
		c.mu.Unlock()

		return schemas, nil
	})
	if err != nil {
		return nil, err
	}

	return lookupKind(v.(map[string]*spec.Schema), group, version, kind)
}

// lookupKind returns the schema for kind from the given GV map, or a
// schema-not-found error wrapping resolver.ErrSchemaNotFound when the
// kind is absent. The wrapped sentinel preserves the error semantics
// of the upstream Kubernetes resolver.
func lookupKind(schemas map[string]*spec.Schema, group, version, kind string) (*spec.Schema, error) {
	if s, ok := schemas[kind]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("cannot resolve group version kind %q: %w",
		group+"/"+version+"/"+kind, resolver.ErrSchemaNotFound)
}

// gvCacheKey builds a cache key from the cluster/group/version tuple.
func (c *DiscoveryCache) gvCacheKey(cluster, group, version string) string {
	return strings.Join([]string{cluster, group, version}, "/")
}

// jitteredTTL returns the configured TTL with ±ttlJitterFraction
// random jitter applied. Spreading expirations prevents the
// thundering-herd stampede that occurs when a batch of entries
// populated at the same instant all expire together.
func (c *DiscoveryCache) jitteredTTL() time.Duration {
	jitter := (c.jitterSampler()*2 - 1) * ttlJitterFraction
	return c.ttl + time.Duration(float64(c.ttl)*jitter)
}

// StartEvictionLoop launches a background goroutine that periodically
// removes expired cache entries. This prevents memory leaks when
// clusters go offline or schemas are no longer queried. It blocks
// until ctx is canceled.
func (c *DiscoveryCache) StartEvictionLoop(ctx context.Context, interval time.Duration) {
	log := slog.Default().With("component", "discovery-cache-evictor")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			before := len(c.gvCache)
			c.evictExpiredGVs()
			after := len(c.gvCache)
			c.mu.Unlock()

			if evicted := before - after; evicted > 0 {
				log.Info("evicted expired cache entries", "count", evicted)
			}
		}
	}
}

// evictExpiredGVs removes expired entries from the group/version
// cache. Must be called with mu held for writing.
func (c *DiscoveryCache) evictExpiredGVs() {
	now := c.now()
	for key, entry := range c.gvCache {
		if now.After(entry.expiresAt) {
			delete(c.gvCache, key)
		}
	}
}
