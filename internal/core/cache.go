package core

import (
	"context"
	"time"
)

// CacheEvictor is a cache that periodically evicts expired entries.
// Implementations live in the infrastructure layer (providers/cache).
type CacheEvictor interface {
	StartEvictionLoop(ctx context.Context, interval time.Duration)
}
