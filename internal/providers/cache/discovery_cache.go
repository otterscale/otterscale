// Package cache provides TTL-based caching infrastructure for
// Kubernetes discovery data. It lives in the providers layer because
// caching is an infrastructure concern — the domain layer
// (internal/core) only defines the SchemaResolver interface.
package cache

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
	"k8s.io/kube-openapi/pkg/validation/spec"

	"github.com/otterscale/otterscale/internal/core"
)

// DefaultTTL is the default TTL for cached OpenAPI schemas.
// Exported so that the DI layer can use it when constructing a
// DiscoveryCache.
const DefaultTTL = 10 * time.Minute

// defaultMaxSchemaEntries is the upper bound on the number of schema
// cache entries. When exceeded, expired entries are eagerly evicted
// before inserting new ones.
const defaultMaxSchemaEntries = 10000

// DiscoveryCache provides TTL-based caching with singleflight
// deduplication for OpenAPI schemas. It implements
// core.SchemaResolver and core.CacheEvictor, and reduces redundant
// discovery API calls when multiple concurrent requests target the
// same cluster.
type DiscoveryCache struct {
	discovery        core.DiscoveryClient
	ttl              time.Duration
	now              func() time.Time
	maxSchemaEntries int

	mu            sync.RWMutex
	schemaCache   map[string]*schemaCacheEntry
	schemaFlights singleflight.Group
}

// schemaCacheEntry pairs a cached schema with its expiration time.
type schemaCacheEntry struct {
	schema    *spec.Schema
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

// WithMaxSchemaEntries overrides the default upper bound on cached
// schema entries.
func WithMaxSchemaEntries(n int) Option {
	return func(c *DiscoveryCache) {
		if n > 0 {
			c.maxSchemaEntries = n
		}
	}
}

// NewDiscoveryCache returns a DiscoveryCache that wraps the given
// DiscoveryClient and caches results for the specified TTL.
func NewDiscoveryCache(discovery core.DiscoveryClient, ttl time.Duration, opts ...Option) *DiscoveryCache {
	c := &DiscoveryCache{
		discovery:        discovery,
		ttl:              ttl,
		now:              time.Now,
		maxSchemaEntries: defaultMaxSchemaEntries,
		schemaCache:      make(map[string]*schemaCacheEntry),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// ResolveSchema fetches the OpenAPI schema for the given GVK. Results
// are cached for the configured TTL and concurrent requests for the
// same key are deduplicated via singleflight.
func (c *DiscoveryCache) ResolveSchema(
	ctx context.Context,
	cluster, group, version, kind string,
) (*spec.Schema, error) {
	key := c.schemaCacheKey(cluster, group, version, kind)

	c.mu.RLock()
	entry, ok := c.schemaCache[key]
	c.mu.RUnlock()

	if ok && c.now().Before(entry.expiresAt) {
		return entry.schema, nil
	}

	v, err, _ := c.schemaFlights.Do(key, func() (any, error) {
		// Use a non-cancellable context with its own timeout so that
		// a single caller's cancellation does not fail all waiters
		// sharing this singleflight key.
		fetchCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), singleflightFetchTimeout)
		defer cancel()

		resolved, err := c.discovery.ResolveSchema(fetchCtx, cluster, group, version, kind)
		if err != nil {
			return nil, err
		}

		c.mu.Lock()
		// Enforce size limit: eagerly evict expired entries before
		// inserting a new one to stay within the bound.
		if len(c.schemaCache) >= c.maxSchemaEntries {
			c.evictExpiredSchemas()
		}
		// Only cache if eviction freed enough space; otherwise
		// return the result uncached to prevent unbounded growth.
		if len(c.schemaCache) < c.maxSchemaEntries {
			c.schemaCache[key] = &schemaCacheEntry{
				schema:    resolved,
				expiresAt: c.now().Add(c.ttl),
			}
		}
		c.mu.Unlock()

		return resolved, nil
	})
	if err != nil {
		return nil, err
	}

	return v.(*spec.Schema), nil
}

// schemaCacheKey builds a cache key from the cluster/group/version/kind tuple.
func (c *DiscoveryCache) schemaCacheKey(cluster, group, version, kind string) string {
	return strings.Join([]string{cluster, group, version, kind}, "/")
}

// StartEvictionLoop launches a background goroutine that periodically
// removes expired cache entries. This prevents memory leaks when
// clusters go offline or schemas are no longer queried. It blocks
// until ctx is cancelled.
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
			before := len(c.schemaCache)
			c.evictExpiredSchemas()
			after := len(c.schemaCache)
			c.mu.Unlock()

			if evicted := before - after; evicted > 0 {
				log.Info("evicted expired cache entries", "count", evicted)
			}
		}
	}
}

// evictExpiredSchemas removes expired entries from the schema cache.
// Must be called with mu held for writing.
func (c *DiscoveryCache) evictExpiredSchemas() {
	now := c.now()
	for key, entry := range c.schemaCache {
		if now.After(entry.expiresAt) {
			delete(c.schemaCache, key)
		}
	}
}
