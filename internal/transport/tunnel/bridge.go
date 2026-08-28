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

// Bridge is a TCP-to-pipe relay. It listens on a localhost TCP port
// (for chisel to forward to) and bridges every accepted connection
// into a pipe.Listener via net.Pipe. This keeps the HTTP server
// completely off the network: it only sees in-memory pipe
// connections supplied by this bridge.
//
// Bridge implements transport.Listener.
type Bridge struct {
	pipeListener *pipe.Listener
	tcpListener  net.Listener
	log          *slog.Logger
	wg           sync.WaitGroup
}

// NewBridge creates a Bridge that feeds connections into pl.
// It binds to an ephemeral localhost TCP port immediately so that
// Port() is available before Start is called.
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

// Port returns the TCP port the bridge is listening on. The tunnel
// client should forward to this port.
func (b *Bridge) Port() int {
	return b.tcpListener.Addr().(*net.TCPAddr).Port
}

// Start accepts TCP connections and bridges them into the pipe
// listener. It blocks until ctx is canceled or an unrecoverable
// error occurs.
func (b *Bridge) Start(ctx context.Context) error {
	b.log.Info("starting", "address", b.tcpListener.Addr().String())

	// Close the TCP listener when the context is done so that Accept
	// unblocks. AfterFunc rather than a goroutine parked on ctx.Done():
	// the stop function releases the watcher when Start returns for its
	// own reasons, instead of leaving it alive until the process ends.
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
			// Deadline-exceeded errors (e.g. Accept with a
			// timeout set) are retried; other errors stop the
			// loop. We use errors.Is instead of the deprecated
			// net.Error.Timeout() method.
			if errors.Is(err, os.ErrDeadlineExceeded) {
				b.log.Warn("temporary accept error", "error", err)
				continue
			}
			// Fall through to the same wait as a clean shutdown: relays
			// already in flight own a pipe connection each, and
			// abandoning them here would return while they still run.
			acceptErr = fmt.Errorf("bridge accept: %w", err)
			break
		}

		b.wg.Add(1)
		go b.relay(tcpConn)
	}

	b.wg.Wait()
	return acceptErr
}

// Stop gracefully shuts down the bridge. It closes the TCP listener and
// the pipe listener, then waits for in-flight relays to finish — but
// only for as long as ctx allows.
//
// The deadline matters: a relay copies until its connection closes, and
// a request still streaming through the tunnel (a watch, a log follow,
// an exec) holds one open indefinitely. Waiting unconditionally would
// ignore the shutdown budget the caller allotted and hang the whole
// process on a client that never disconnects.
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

// relay bridges a single TCP connection to the pipe listener. It
// creates a net.Pipe pair, hands the server end to the pipe listener,
// and copies data bidirectionally between the TCP connection and the
// client end of the pipe.
//
// When either copy direction finishes (typically because the HTTP
// handler closed its end of the pipe), both connections are closed so
// the other direction terminates as well.
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

// ErrBridgeRequired is returned when a Bridge is expected but nil.
var ErrBridgeRequired = errors.New("tunnel: bridge is required")
