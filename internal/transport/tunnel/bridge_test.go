package tunnel

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/otterscale/otterscale/internal/transport/pipe"
)

// TestBridge_RelaysData verifies that a TCP client can exchange data
// with a server behind the pipe listener through the bridge.
func TestBridge_RelaysData(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pl := pipe.NewListener()
	defer pl.Close()

	bridge, err := NewBridge(ctx, pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	go func() {
		if err := bridge.Start(ctx); err != nil {
			t.Logf("bridge.Start: %v", err)
		}
	}()

	const request = "hello"
	const response = "world"

	// Server side: read a fixed-size request, send a response, close.
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

	// Client side: connect to the bridge TCP port, send request, read response.
	var d net.Dialer
	tcpConn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("127.0.0.1:%d", bridge.Port()))
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

// TestBridge_MultipleConnections verifies that the bridge can handle
// several concurrent connections, each independently relaying data.
func TestBridge_MultipleConnections(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pl := pipe.NewListener()
	defer pl.Close()

	bridge, err := NewBridge(ctx, pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	go func() {
		if err := bridge.Start(ctx); err != nil {
			t.Logf("bridge.Start: %v", err)
		}
	}()

	const n = 5
	var wg sync.WaitGroup

	// Server side: accept n connections and echo back.
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			echoConnection(t, pl, i)
		}()
	}

	// Client side: dial n connections concurrently.
	addr := fmt.Sprintf("127.0.0.1:%d", bridge.Port())
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			verifyRoundTrip(ctx, t, addr, i)
		}()
	}

	wg.Wait()
	cancel()
}

// echoConnection accepts a pipe connection, reads a message, and
// sends it back unchanged. Used by TestBridge_MultipleConnections.
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

// verifyRoundTrip dials a TCP address, sends a message, reads it
// back, and verifies it matches. Used by TestBridge_MultipleConnections.
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pl := pipe.NewListener()
	defer pl.Close()

	bridge, err := NewBridge(ctx, pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}
	defer func() {
		if err := bridge.Stop(context.Background()); err != nil {
			t.Logf("bridge.Stop: %v", err)
		}
	}()

	if bridge.Port() == 0 {
		t.Fatal("expected non-zero port")
	}
}

func TestBridge_StopClosesListener(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	pl := pipe.NewListener()
	bridge, err := NewBridge(ctx, pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- bridge.Start(ctx)
	}()

	// Give Start time to begin accepting.
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
