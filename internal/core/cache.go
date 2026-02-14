package core

import (
	"context"
	"time"
)

// CacheEvictor represents a cache that supports periodic eviction of
// expired entries. Implementations live in the infrastructure layer
// (e.g. providers/cache). Defining the interface here decouples the
// application layer from concrete cache implementations.
type CacheEvictor interface {
	StartEvictionLoop(ctx context.Context, interval time.Duration)
}
