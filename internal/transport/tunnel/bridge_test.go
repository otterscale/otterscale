package tunnel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/otterscale/otterscale/internal/transport/pipe"
)

// TestBridge_RelaysData exchanges data with a server behind the pipe listener.
func TestBridge_RelaysData(t *testing.T) {
	t.Parallel()

	pl := pipe.NewListener()
	defer pl.Close()

	bridge, err := NewBridge(t.Context(), pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	go func() {
		if err := bridge.Start(t.Context()); err != nil {
			t.Logf("bridge.Start: %v", err)
		}
	}()

	const request = "hello"
	const response = "world"

	// Server: read a fixed-size request, respond, close.
	go func() {
		conn, err := pl.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		buf := make([]byte, len(request))
		if _, err := io.ReadFull(conn, buf); err != nil {
			return
		}
		if _, err := conn.Write([]byte(response)); err != nil {
			return
		}
	}()

	// Client: connect to the bridge TCP port and round-trip.
	var d net.Dialer
	tcpConn, err := d.DialContext(t.Context(), "tcp", fmt.Sprintf("127.0.0.1:%d", bridge.Port()))
	if err != nil {
		t.Fatalf("tcp dial: %v", err)
	}
	defer tcpConn.Close()

	if _, err := tcpConn.Write([]byte(request)); err != nil {
		t.Fatalf("write: %v", err)
	}

	buf := make([]byte, len(response))
	if _, err := io.ReadFull(tcpConn, buf); err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(buf) != response {
		t.Fatalf("got %q, want %q", buf, response)
	}
}

// TestBridge_MultipleConnections relays several concurrent connections
// independently.
func TestBridge_MultipleConnections(t *testing.T) {
	t.Parallel()

	pl := pipe.NewListener()
	defer pl.Close()

	bridge, err := NewBridge(t.Context(), pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	go func() {
		if err := bridge.Start(t.Context()); err != nil {
			t.Logf("bridge.Start: %v", err)
		}
	}()

	const n = 5
	var wg sync.WaitGroup

	// Server: accept n connections and echo back.
	for i := range n {
		wg.Go(func() {
			echoConnection(t, pl, i)
		})
	}

	// Client: dial n connections concurrently.
	addr := fmt.Sprintf("127.0.0.1:%d", bridge.Port())
	for i := range n {
		wg.Go(func() {
			verifyRoundTrip(t.Context(), t, addr, i)
		})
	}

	wg.Wait()
}

// echoConnection reads one message off a pipe connection and sends it back.
func echoConnection(t *testing.T, pl *pipe.Listener, i int) {
	t.Helper()
	conn, err := pl.Accept()
	if err != nil {
		t.Errorf("pipe Accept #%d: %v", i, err)
		return
	}
	defer conn.Close()

	msg := fmt.Sprintf("msg-%d", i)
	buf := make([]byte, len(msg))
	if _, err := io.ReadFull(conn, buf); err != nil {
		t.Errorf("server read #%d: %v", i, err)
		return
	}
	if _, err := conn.Write(buf); err != nil {
		t.Errorf("server write #%d: %v", i, err)
	}
}

// verifyRoundTrip dials addr, sends a message, and checks what comes back.
func verifyRoundTrip(ctx context.Context, t *testing.T, addr string, i int) {
	t.Helper()
	var d net.Dialer
	tcpConn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		t.Errorf("tcp dial #%d: %v", i, err)
		return
	}
	defer tcpConn.Close()

	msg := fmt.Sprintf("msg-%d", i)
	if _, err := tcpConn.Write([]byte(msg)); err != nil {
		t.Errorf("write #%d: %v", i, err)
		return
	}

	buf := make([]byte, len(msg))
	if _, err := io.ReadFull(tcpConn, buf); err != nil {
		t.Errorf("read #%d: %v", i, err)
		return
	}
	if string(buf) != msg {
		t.Errorf("#%d: got %q, want %q", i, buf, msg)
	}
}

func TestBridge_PortIsNonZero(t *testing.T) {
	t.Parallel()

	pl := pipe.NewListener()
	defer pl.Close()

	bridge, err := NewBridge(t.Context(), pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}
	defer func() {
		if err := bridge.Stop(t.Context()); err != nil {
			t.Logf("bridge.Stop: %v", err)
		}
	}()

	if bridge.Port() == 0 {
		t.Fatal("expected non-zero port")
	}
}

func TestBridge_StopClosesListener(t *testing.T) {
	t.Parallel()

	pl := pipe.NewListener()
	bridge, err := NewBridge(t.Context(), pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- bridge.Start(ctx)
	}()

	time.Sleep(20 * time.Millisecond)

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Start returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return after cancel")
	}
}

// TestBridge_StopHonoursDeadline is the regression test for a shutdown
// that could not be bounded. A relay copies until its connection
// closes, and a request still streaming through the tunnel — a watch, a
// log follow, an exec — holds one open indefinitely. Stop used to wait
// on those unconditionally, ignoring the budget transport.Serve allots
// each listener and hanging the process on a client that never hangs up.
func TestBridge_StopHonoursDeadline(t *testing.T) {
	t.Parallel()

	pl := pipe.NewListener()
	defer pl.Close()

	bridge, err := NewBridge(t.Context(), pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	startBridge(t, bridge)

	// A server behind the pipe that accepts and then never says
	// anything, standing in for a long-lived stream.
	serverHolds := make(chan struct{})
	go func() {
		conn, acceptErr := pl.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()
		<-serverHolds
	}()
	defer close(serverHolds)

	var dialer net.Dialer
	conn, err := dialer.DialContext(t.Context(), "tcp", fmt.Sprintf("127.0.0.1:%d", bridge.Port()))
	if err != nil {
		t.Fatalf("dial bridge: %v", err)
	}
	defer conn.Close()

	time.Sleep(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	stopped := make(chan error, 1)
	go func() { stopped <- bridge.Stop(ctx) }()

	select {
	case err := <-stopped:
		if err == nil {
			t.Fatal("Stop reported success while a relay was still in flight")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Stop error = %v, want it to wrap context.DeadlineExceeded", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Stop ignored its deadline and blocked on an in-flight relay")
	}
}

// TestBridge_StopWaitsForIdleRelays checks the other half: when nothing
// is in flight, Stop still returns cleanly rather than reporting the
// deadline it never needed.
func TestBridge_StopWaitsForIdleRelays(t *testing.T) {
	t.Parallel()

	pl := pipe.NewListener()
	bridge, err := NewBridge(t.Context(), pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	startBridge(t, bridge)
	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	if err := bridge.Stop(ctx); err != nil {
		t.Fatalf("Stop on an idle bridge: %v", err)
	}
}

// startBridge runs Start in the background and makes the test wait for it on
// the way out. Start outlives Stop — which closes the listener without
// canceling Start's context — so a test that left it behind would log from a
// goroutine after finishing.
func startBridge(t *testing.T, bridge *Bridge) {
	t.Helper()

	ctx, cancel := context.WithCancel(t.Context())
	stopped := make(chan error, 1)

	go func() { stopped <- bridge.Start(ctx) }()

	t.Cleanup(func() {
		cancel()
		<-stopped
	})
}
