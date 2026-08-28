package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/otterscale/otterscale/internal/core"
)

// copyTimeout bounds every copy-loop assertion. The bug these tests
// guard against is a deadlock, so a hung copy loop must fail the test
// rather than hang the suite.
const copyTimeout = 2 * time.Second

var errFakeConnClosed = errors.New("fake vnc conn closed")

// fakeVNCConn implements vncConn over channels so the copy loop can be
// exercised without a live WebSocket server.
type fakeVNCConn struct {
	incoming chan []byte // messages handed to ReadMessage
	readErr  error       // returned once incoming is closed; set before closing
	written  chan []byte // messages captured from WriteMessage

	closed    chan struct{}
	closeOnce sync.Once
}

func newFakeVNCConn() *fakeVNCConn {
	return &fakeVNCConn{
		incoming: make(chan []byte, 4),
		written:  make(chan []byte, 4),
		closed:   make(chan struct{}),
	}
}

func (c *fakeVNCConn) ReadMessage() (messageType int, p []byte, err error) {
	select {
	case msg, ok := <-c.incoming:
		if !ok {
			if c.readErr != nil {
				return 0, nil, c.readErr
			}
			return 0, nil, errFakeConnClosed
		}
		return websocket.BinaryMessage, msg, nil
	case <-c.closed:
		return 0, nil, errFakeConnClosed
	}
}

func (c *fakeVNCConn) WriteMessage(_ int, data []byte) error {
	select {
	case c.written <- append([]byte(nil), data...):
		return nil
	case <-c.closed:
		return errFakeConnClosed
	}
}

func (c *fakeVNCConn) Close() error {
	c.closeOnce.Do(func() { close(c.closed) })
	return nil
}

func (c *fakeVNCConn) isClosed() bool {
	select {
	case <-c.closed:
		return true
	default:
		return false
	}
}

// chanWriter captures writes so a test can observe the VMI → client
// direction without blocking the copy loop.
type chanWriter struct {
	ch chan []byte
}

func (w chanWriter) Write(p []byte) (int, error) {
	w.ch <- append([]byte(nil), p...)
	return len(p), nil
}

// vncFixture wires a fake connection to a pair of pipes standing in for
// the session's stdin/stdout, mirroring what core.StartVNC provides.
type vncFixture struct {
	conn   *fakeVNCConn
	stdinW *io.PipeWriter
	stdout chan []byte
	opts   core.VNCOptions
}

func newVNCFixture(t *testing.T) *vncFixture {
	t.Helper()

	conn := newFakeVNCConn()
	stdinR, stdinW := io.Pipe()
	stdout := make(chan []byte, 4)

	f := &vncFixture{
		conn:   conn,
		stdinW: stdinW,
		stdout: stdout,
		opts: core.VNCOptions{
			Stdin:  stdinR,
			Stdout: chanWriter{ch: stdout},
		},
	}
	// Closing the read end releases the stdin copy goroutine, which is
	// what core.StartVNC does once the session ends.
	t.Cleanup(func() {
		stdinW.Close()
		stdinR.Close()
		conn.Close()
	})
	return f
}

// run executes the copy loop and fails the test if it does not return
// within copyTimeout.
func (f *vncFixture) run(ctx context.Context, t *testing.T) error {
	t.Helper()

	errc := make(chan error, 1)
	go func() {
		errc <- (&runtimeRepo{}).copyVNCBidirectional(ctx, f.conn, f.opts)
	}()

	select {
	case err := <-errc:
		return err
	case <-time.After(copyTimeout):
		t.Fatalf("copyVNCBidirectional did not return within %s", copyTimeout)
		return nil
	}
}

// TestCopyVNCBidirectional_CancelWhileIdle is the regression test for
// the session leak: with no traffic in either direction, cancellation
// alone must end the copy loop.
func TestCopyVNCBidirectional_CancelWhileIdle(t *testing.T) {
	f := newVNCFixture(t)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if err := f.run(ctx, t); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
	if !f.conn.isClosed() {
		t.Error("expected the WebSocket connection to be closed on return")
	}
}

// TestCopyVNCBidirectional_PeerGracefulClose covers the case that used
// to enqueue two results from one direction and wedge the loop.
func TestCopyVNCBidirectional_PeerGracefulClose(t *testing.T) {
	f := newVNCFixture(t)

	f.conn.readErr = &websocket.CloseError{Code: websocket.CloseNormalClosure}
	close(f.conn.incoming)

	if err := f.run(t.Context(), t); err != nil {
		t.Fatalf("got %v, want nil for a graceful peer close", err)
	}
}

func TestCopyVNCBidirectional_ReadErrorSurfaces(t *testing.T) {
	f := newVNCFixture(t)

	boom := errors.New("boom")
	f.conn.readErr = boom
	close(f.conn.incoming)

	if err := f.run(t.Context(), t); !errors.Is(err, boom) {
		t.Fatalf("got %v, want %v", err, boom)
	}
}

func TestCopyVNCBidirectional_StdinEOF(t *testing.T) {
	f := newVNCFixture(t)

	// The client closing its side of the session is a clean shutdown.
	f.stdinW.Close()

	if err := f.run(t.Context(), t); err != nil {
		t.Fatalf("got %v, want nil for stdin EOF", err)
	}
}

func TestCopyVNCBidirectional_CopiesBothDirections(t *testing.T) {
	f := newVNCFixture(t)

	ctx, cancel := context.WithCancel(t.Context())

	errc := make(chan error, 1)
	go func() {
		errc <- (&runtimeRepo{}).copyVNCBidirectional(ctx, f.conn, f.opts)
	}()

	// io.Pipe is synchronous, so stdin can only be written once the
	// copy loop is running to consume it.
	f.conn.incoming <- []byte("from-vmi")
	if _, err := f.stdinW.Write([]byte("from-client")); err != nil {
		t.Fatalf("write stdin: %v", err)
	}

	select {
	case got := <-f.stdout:
		if string(got) != "from-vmi" {
			t.Errorf("stdout = %q, want %q", got, "from-vmi")
		}
	case <-time.After(copyTimeout):
		t.Fatal("VMI → client data never reached stdout")
	}

	select {
	case got := <-f.conn.written:
		if string(got) != "from-client" {
			t.Errorf("written = %q, want %q", got, "from-client")
		}
	case <-time.After(copyTimeout):
		t.Fatal("client → VMI data never reached the connection")
	}

	cancel()
	select {
	case <-errc:
	case <-time.After(copyTimeout):
		t.Fatalf("copyVNCBidirectional did not return within %s of cancellation", copyTimeout)
	}
}

func TestVNCURL(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		vmi       string
		want      string
	}{
		{
			name:      "plain names",
			namespace: "default",
			vmi:       "my-vmi",
			want:      "https://127.1.1.1:16598/apis/subresources.kubevirt.io/v1/namespaces/default/virtualmachineinstances/my-vmi/vnc",
		},
		{
			name:      "a name carrying path separators cannot escape its segment",
			namespace: "default",
			vmi:       "../../../api/v1/nodes",
			want:      "https://127.1.1.1:16598/apis/subresources.kubevirt.io/v1/namespaces/default/virtualmachineinstances/..%2F..%2F..%2Fapi%2Fv1%2Fnodes/vnc",
		},
		{
			name:      "a namespace carrying a query cannot append parameters",
			namespace: "ns?watch=true",
			vmi:       "my-vmi",
			want:      "https://127.1.1.1:16598/apis/subresources.kubevirt.io/v1/namespaces/ns%3Fwatch=true/virtualmachineinstances/my-vmi/vnc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vncURL("https://127.1.1.1:16598", tt.namespace, tt.vmi); got != tt.want {
				t.Errorf("vncURL() =\n\t%q\nwant\n\t%q", got, tt.want)
			}
		})
	}
}

// stubExitError implements utilexec.ExitError, which is what
// remotecommand returns when a command runs and exits non-zero.
type stubExitError struct {
	code   int
	exited bool
}

func (e stubExitError) Error() string {
	return fmt.Sprintf("command terminated with exit code %d", e.code)
}
func (e stubExitError) String() string  { return e.Error() }
func (e stubExitError) Exited() bool    { return e.exited }
func (e stubExitError) ExitStatus() int { return e.code }

// TestWrapExecError checks the distinction the handler layer depends
// on: a command that ran and exited non-zero is not the same kind of
// failure as a session that never started, and only the latter should
// reach the caller as an RPC error.
func TestWrapExecError(t *testing.T) {
	t.Run("nil stays nil", func(t *testing.T) {
		if err := wrapExecError(nil); err != nil {
			t.Fatalf("wrapExecError(nil) = %v, want nil", err)
		}
	})

	t.Run("non-zero exit becomes ErrCommandExited", func(t *testing.T) {
		err := wrapExecError(stubExitError{code: 3, exited: true})

		var exited *core.ErrCommandExited
		if !errors.As(err, &exited) {
			t.Fatalf("wrapExecError = %T (%v), want *core.ErrCommandExited", err, err)
		}
		if exited.Code != 3 {
			t.Fatalf("exit code = %d, want 3", exited.Code)
		}
	})

	t.Run("signaled command is not an exit status", func(t *testing.T) {
		// Exited() reports false when the process was signaled rather
		// than exiting on its own; that is a session failure, not a
		// command result.
		err := wrapExecError(stubExitError{code: -1, exited: false})

		var exited *core.ErrCommandExited
		if errors.As(err, &exited) {
			t.Fatal("wrapExecError classified a signaled command as an exit status")
		}
	})

	t.Run("api status maps to a domain code", func(t *testing.T) {
		err := wrapExecError(apierrors.NewForbidden(
			schema.GroupResource{Resource: "pods"}, "pod-a", errors.New("denied")))

		code, ok := core.DomainErrorCode(err)
		if !ok {
			t.Fatalf("wrapExecError = %T (%v), want a *core.DomainError", err, err)
		}
		if code != core.ErrorCodePermissionDenied {
			t.Fatalf("code = %v, want ErrorCodePermissionDenied", code)
		}
	})
}
