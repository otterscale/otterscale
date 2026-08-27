package kubernetes

import (
	"testing"
	"time"
)

// TestWatchListCache covers the TTL bookkeeping that keeps a /version
// round-trip off the critical path of every new watch.
func TestWatchListCache(t *testing.T) {
	now := time.Now()
	d := &discoveryClient{
		now:       func() time.Time { return now },
		watchList: make(map[string]watchListSupport),
	}

	if _, ok := d.cachedWatchList("prod"); ok {
		t.Fatal("expected a miss on an empty cache")
	}

	d.cacheWatchList("prod", true)

	supported, ok := d.cachedWatchList("prod")
	if !ok {
		t.Fatal("expected a hit right after caching")
	}
	if !supported {
		t.Error("cached capability = false, want true")
	}

	if _, ok := d.cachedWatchList("staging"); ok {
		t.Error("expected a miss for a cluster that was never cached")
	}

	// A negative answer is cached too, so an older cluster is not
	// re-probed on every watch either.
	d.cacheWatchList("staging", false)
	supported, ok = d.cachedWatchList("staging")
	if !ok {
		t.Fatal("expected a hit for a cached negative answer")
	}
	if supported {
		t.Error("cached capability = true, want false")
	}

	now = now.Add(watchListTTL + time.Second)

	if _, ok := d.cachedWatchList("prod"); ok {
		t.Error("expected a miss once the entry has expired")
	}
}
