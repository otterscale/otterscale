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

	pl := pipe.NewListener()
	defer pl.Close()

	bridge, err := NewBridge(pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go bridge.Start(ctx)

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
		conn.Write([]byte(response))
	}()

	// Client side: connect to the bridge TCP port, send request, read response.
	tcpConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", bridge.Port()))
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

	pl := pipe.NewListener()
	defer pl.Close()

	bridge, err := NewBridge(pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go bridge.Start(ctx)

	const n = 5
	var wg sync.WaitGroup

	// Server side: accept n connections; each reads a request and
	// sends the same bytes back as a response, then closes.
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
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
		}()
	}

	// Client side: dial n connections concurrently.
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tcpConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", bridge.Port()))
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
		}()
	}

	wg.Wait()
	cancel()
}

func TestBridge_PortIsNonZero(t *testing.T) {
	t.Parallel()

	pl := pipe.NewListener()
	defer pl.Close()

	bridge, err := NewBridge(pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}
	defer bridge.Stop(context.Background())

	if bridge.Port() == 0 {
		t.Fatal("expected non-zero port")
	}
}

func TestBridge_StopClosesListener(t *testing.T) {
	t.Parallel()

	pl := pipe.NewListener()
	bridge, err := NewBridge(pl)
	if err != nil {
		t.Fatalf("NewBridge: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
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
