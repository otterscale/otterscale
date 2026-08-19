package server

import (
	"context"
	"time"

	"github.com/otterscale/otterscale/internal/core"
)

// sessionReapInterval is the interval at which the session reaper
// scans for and removes stale sessions.
const sessionReapInterval = 30 * time.Second

// cacheEvictionInterval is the interval at which the discovery cache
// evictor removes expired schema and version entries.
const cacheEvictionInterval = 5 * time.Minute

// ProvideBackgroundListeners constructs the background transport
// listeners (session reaper, cache evictor) that participate in the
// server's managed lifecycle. The CacheEvictor interface decouples
// this function from the concrete cache implementation, keeping the
// application layer free of infrastructure dependencies.
func ProvideBackgroundListeners(runtime *core.RuntimeUseCase, evictor core.CacheEvictor) BackgroundListeners {
	return BackgroundListeners{
		&sessionReaperListener{runtime: runtime},
		&cacheEvictorListener{cache: evictor},
	}
}

// sessionReaperListener adapts RuntimeUseCase.StartSessionReaper to
// the transport.Listener interface so it participates in the managed
// lifecycle alongside other servers.
type sessionReaperListener struct {
	runtime *core.RuntimeUseCase
}

func (l *sessionReaperListener) Start(ctx context.Context) error {
	l.runtime.StartSessionReaper(ctx, sessionReapInterval)
	return nil
}

func (l *sessionReaperListener) Stop(_ context.Context) error {
	return nil // reaper stops when its context is canceled
}

// cacheEvictorListener adapts a CacheEvictor to the
// transport.Listener interface so it participates in the managed
// lifecycle alongside other servers.
type cacheEvictorListener struct {
	cache core.CacheEvictor
}

func (l *cacheEvictorListener) Start(ctx context.Context) error {
	l.cache.StartEvictionLoop(ctx, cacheEvictionInterval)
	return nil
}

func (l *cacheEvictorListener) Stop(_ context.Context) error {
	return nil // evictor stops when its context is canceled
}
