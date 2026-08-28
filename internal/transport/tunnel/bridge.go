package tunnel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"

	"github.com/otterscale/otterscale/internal/transport/pipe"
)

// Bridge listens on a localhost TCP port for chisel to forward to and relays
// every accepted connection into a pipe.Listener via net.Pipe. That keeps the
// HTTP server off the network entirely: it only ever sees in-memory pipe
// connections. Bridge implements transport.Listener.
type Bridge struct {
	pipeListener *pipe.Listener
	tcpListener  net.Listener
	log          *slog.Logger
	wg           sync.WaitGroup
}

// NewBridge binds an ephemeral localhost TCP port immediately, so Port() is
// available before Start is called.
func NewBridge(ctx context.Context, pl *pipe.Listener) (*Bridge, error) {
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("bridge listen: %w", err)
	}
	return &Bridge{
		pipeListener: pl,
		tcpListener:  ln,
		log:          slog.Default().With("component", "tunnel-bridge"),
	}, nil
}

// Port is where the tunnel client should forward.
func (b *Bridge) Port() int {
	return b.tcpListener.Addr().(*net.TCPAddr).Port
}

// Start blocks until ctx is canceled or an unrecoverable error occurs.
func (b *Bridge) Start(ctx context.Context) error {
	b.log.Info("starting", "address", b.tcpListener.Addr().String())

	// Closing the listener is what unblocks Accept. AfterFunc rather than a
	// goroutine parked on ctx.Done(): the stop function releases the watcher
	// when Start returns for its own reasons.
	stopWatch := context.AfterFunc(ctx, func() {
		b.tcpListener.Close()
	})
	defer stopWatch()

	var acceptErr error
	for {
		tcpConn, err := b.tcpListener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			// Deadline-exceeded is retried; anything else stops the loop.
			// errors.Is, not the deprecated net.Error.Timeout().
			if errors.Is(err, os.ErrDeadlineExceeded) {
				b.log.Warn("temporary accept error", "error", err)
				continue
			}
			// Fall through to the same wait as a clean shutdown: relays in
			// flight own a pipe connection each, and abandoning them here
			// would return while they still run.
			acceptErr = fmt.Errorf("bridge accept: %w", err)
			break
		}

		b.wg.Add(1)
		go b.relay(tcpConn)
	}

	b.wg.Wait()
	return acceptErr
}

// Stop closes both listeners, then waits for in-flight relays — but only as
// long as ctx allows.
//
// The deadline matters: a relay copies until its connection closes, and a
// request still streaming through the tunnel (a watch, a log follow, an exec)
// holds one open indefinitely. Waiting unconditionally would ignore the
// caller's shutdown budget and hang the process on a client that never
// disconnects.
func (b *Bridge) Stop(ctx context.Context) error {
	b.log.Info("shutting down")
	b.tcpListener.Close()
	b.pipeListener.Close()

	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		// The relays are left running: they hold connections this
		// bridge does not own, and the process is going away anyway.
		b.log.Warn("shutdown deadline reached with relays still in flight")
		return fmt.Errorf("bridge shutdown: %w", ctx.Err())
	}
}

// relay hands the server end of a net.Pipe pair to the pipe listener and copies
// bidirectionally against the client end. When either direction finishes —
// usually because the HTTP handler closed its end — both connections are closed
// so the other unwinds too.
func (b *Bridge) relay(tcpConn net.Conn) {
	defer b.wg.Done()

	pipeConn, err := b.pipeListener.Dial()
	if err != nil {
		b.log.Debug("pipe dial failed, listener likely closed", "error", err)
		tcpConn.Close()
		return
	}

	errc := make(chan error, 2)
	go func() {
		_, err := io.Copy(tcpConn, pipeConn) // pipe → TCP
		errc <- err
	}()
	go func() {
		_, err := io.Copy(pipeConn, tcpConn) // TCP → pipe
		errc <- err
	}()

	<-errc // first direction done
	pipeConn.Close()
	tcpConn.Close()
	<-errc // second direction done
}

var ErrBridgeRequired = errors.New("tunnel: bridge is required")
