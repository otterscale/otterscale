// Package transport coordinates the lifecycle of multiple server
// components (HTTP, tunnel, health-check, etc.) using an errgroup.
package transport

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"
)

// shutdownTimeout bounds each listener's graceful shutdown.
const shutdownTimeout = 15 * time.Second

// Listener is a component in the server lifecycle. Start blocks until the
// component finishes or ctx is canceled; Stop shuts down gracefully within its
// context deadline.
type Listener interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// TunnelService supplies the tunnel infrastructure the server needs for
// transport setup and health monitoring. It lives here because its methods
// return transport.Listener values, which avoids a reverse dependency from the
// providers layer back into cmd/server.
type TunnelService interface {
	// BuildTunnelListener configures a tunnel server for the given address and
	// host SAN.
	BuildTunnelListener(address, host string) (Listener, error)
	// BuildHealthListener periodically probes registered tunnel endpoints.
	BuildHealthListener() Listener
}

// Serve runs all listeners concurrently and coordinates graceful shutdown.
// Every listener is started first, and only then does a single goroutine wait
// on the derived context and call Stop on each — so Stop never runs before
// Start had a chance to.
func Serve(ctx context.Context, lis ...Listener) error {
	eg, egCtx := errgroup.WithContext(ctx)

	for _, li := range lis {
		eg.Go(func() error {
			return li.Start(egCtx)
		})
	}

	// Cancellation comes from the parent ctx or a listener failure. Each
	// listener gets its own timeout, so a slow one cannot starve the rest.
	eg.Go(func() error {
		<-egCtx.Done()

		var errs []error
		for _, li := range lis {
			stopCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			if err := li.Stop(stopCtx); err != nil {
				errs = append(errs, err)
			}
			cancel()
		}
		return errors.Join(errs...)
	})

	return eg.Wait()
}
