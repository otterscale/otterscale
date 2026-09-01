package server

import (
	"context"
	"errors"
	"time"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

// sessionReapInterval is how often stale sessions are swept.
const sessionReapInterval = 30 * time.Second

// cacheEvictionInterval is how often expired cache entries are removed.
const cacheEvictionInterval = 5 * time.Minute

// ProvideBackgroundListeners builds the session reaper and cache evictor that
// join the server's managed lifecycle. Taking evictor as an interface keeps the
// application layer off the concrete cache implementation.
func ProvideBackgroundListeners(runtime *core.RuntimeUseCase, evictor core.CacheEvictor) BackgroundListeners {
	return BackgroundListeners{
		&sessionReaperListener{runtime: runtime},
		&cacheEvictorListener{cache: evictor},
	}
}

// sessionReaperListener adapts RuntimeUseCase.StartSessionReaper to
// transport.Listener.
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

// cacheEvictorListener adapts a CacheEvictor to transport.Listener.
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

// ProvideJoinAuthority fails when no secret is configured: the registration
// endpoint is reachable without authentication, so a server without one would
// let any caller claim — and take over — any cluster.
func ProvideJoinAuthority(conf *config.Config) (*core.JoinAuthority, error) {
	secret, err := conf.ServerJoinSecret()
	if err != nil {
		return nil, err
	}
	if secret == "" {
		return nil, errors.New(
			"join secret is required but not configured; " +
				"set --join-secret, --join-secret-file or OTTERSCALE_SERVER_JOIN_SECRET",
		)
	}
	return core.NewJoinAuthority(secret)
}
