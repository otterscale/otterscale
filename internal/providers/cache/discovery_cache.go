// Package cache provides TTL-based caching for Kubernetes discovery data. It
// lives in the providers layer because caching is an infrastructure concern —
// internal/core only defines the SchemaResolver interface.
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

// DefaultTTL is exported for the DI layer to use when constructing a
// DiscoveryCache.
const DefaultTTL = 10 * time.Minute

// defaultMaxGVEntries bounds cached group/version entries. Each entry holds
// every kind in that group/version, so the bound on kinds cached is much higher
// than this number suggests.
const defaultMaxGVEntries = 1024

// ttlJitterFraction is the maximum random jitter applied to an entry's TTL, as
// a fraction of it. Spreading expirations prevents a cache stampede when many
// entries are populated in the same burst.
const ttlJitterFraction = 0.2

// DiscoveryCache adds TTL caching and singleflight deduplication to OpenAPI
// schema resolution, implementing core.SchemaResolver and core.CacheEvictor.
//
// Entries are keyed at group/version granularity because the Kubernetes OpenAPI
// v3 endpoint returns one document per GV. Caching per GVK would re-download
// the same document for every kind; the GV scheme amortizes one fetch across
// all kinds in that GV and lets singleflight deduplicate concurrent misses for
// different kinds of the same GV.
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

// gvCacheEntry pairs one group/version's kind→schema map with its jittered
// expiration.
type gvCacheEntry struct {
	schemas   map[string]*spec.Schema
	expiresAt time.Time
}

// singleflightFetchTimeout bounds a cache-miss fetch, which runs under
// context.WithoutCancel so one caller's cancellation does not fail all waiters.
const singleflightFetchTimeout = 30 * time.Second

type Option func(*DiscoveryCache)

// WithClock injects a time source for deterministic testing; time.Now otherwise.
func WithClock(now func() time.Time) Option {
	return func(c *DiscoveryCache) {
		c.now = now
	}
}

func WithMaxGVEntries(n int) Option {
	return func(c *DiscoveryCache) {
		if n > 0 {
			c.maxGVEntries = n
		}
	}
}

// WithJitterSampler injects a [0,1) sampler for deterministic testing;
// production relies on the math/rand/v2-backed default.
func WithJitterSampler(sample func() float64) Option {
	return func(c *DiscoveryCache) {
		c.jitterSampler = sample
	}
}

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

// ResolveSchema caches per group/version for the configured TTL and
// deduplicates concurrent requests for the same group/version.
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
		// Non-cancellable with its own timeout, so one caller's cancellation
		// does not fail all waiters sharing this singleflight key.
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

// lookupKind reports an absent kind as a domain not-found error: asking for a
// kind the cluster does not define is a bad request, not a server fault, and
// without a domain code it would surface as an internal error.
// resolver.ErrSchemaNotFound stays in the chain so errors.Is keeps working for
// callers that check the upstream sentinel.
func lookupKind(schemas map[string]*spec.Schema, group, version, kind string) (*spec.Schema, error) {
	if s, ok := schemas[kind]; ok {
		return s, nil
	}
	return nil, &core.DomainError{
		Code:    core.ErrorCodeNotFound,
		Message: fmt.Sprintf("cannot resolve group version kind %q", group+"/"+version+"/"+kind),
		Cause:   resolver.ErrSchemaNotFound,
	}
}

func (c *DiscoveryCache) gvCacheKey(cluster, group, version string) string {
	return strings.Join([]string{cluster, group, version}, "/")
}

// jitteredTTL applies ±ttlJitterFraction, so a batch of entries populated at
// the same instant does not expire together.
func (c *DiscoveryCache) jitteredTTL() time.Duration {
	jitter := (c.jitterSampler()*2 - 1) * ttlJitterFraction
	return c.ttl + time.Duration(float64(c.ttl)*jitter)
}

// StartEvictionLoop periodically removes expired entries, which is what keeps
// clusters going offline from leaking memory. It blocks until ctx is canceled.
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

// evictExpiredGVs must be called with mu held for writing.
func (c *DiscoveryCache) evictExpiredGVs() {
	now := c.now()
	for key, entry := range c.gvCache {
		if now.After(entry.expiresAt) {
			delete(c.gvCache, key)
		}
	}
}
