package pipe

import (
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

func TestListener_DialAndAccept(t *testing.T) {
	t.Parallel()

	ln := NewListener()
	defer ln.Close()

	const msg = "hello pipe"

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := ln.Accept()
		if err != nil {
			t.Errorf("Accept: %v", err)
			return
		}
		defer conn.Close()

		buf := make([]byte, len(msg))
		if _, err := io.ReadFull(conn, buf); err != nil {
			t.Errorf("server read: %v", err)
			return
		}
		if string(buf) != msg {
			t.Errorf("server got %q, want %q", buf, msg)
		}
	}()

	client, err := ln.Dial()
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	if _, err := client.Write([]byte(msg)); err != nil {
		t.Fatalf("client write: %v", err)
	}

	client.Close()
	wg.Wait()
}

func TestListener_BidirectionalData(t *testing.T) {
	t.Parallel()

	ln := NewListener()
	defer ln.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := ln.Accept()
		if err != nil {
			t.Errorf("Accept: %v", err)
			return
		}
		defer conn.Close()

		buf := make([]byte, 4)
		if _, err := io.ReadFull(conn, buf); err != nil {
			t.Errorf("server read: %v", err)
			return
		}
		if _, err := conn.Write([]byte("pong")); err != nil {
			t.Errorf("server write: %v", err)
		}
	}()

	client, err := ln.Dial()
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	if _, err := client.Write([]byte("ping")); err != nil {
		t.Fatalf("client write: %v", err)
	}
	buf := make([]byte, 4)
	if _, err := io.ReadFull(client, buf); err != nil {
		t.Fatalf("client read: %v", err)
	}
	if string(buf) != "pong" {
		t.Fatalf("got %q, want %q", buf, "pong")
	}

	wg.Wait()
}

func TestListener_CloseUnblocksAccept(t *testing.T) {
	t.Parallel()

	ln := NewListener()

	done := make(chan error, 1)
	go func() {
		_, err := ln.Accept()
		done <- err
	}()

	// Give Accept time to block.
	time.Sleep(20 * time.Millisecond)
	ln.Close()

	select {
	case err := <-done:
		if err != net.ErrClosed {
			t.Fatalf("expected net.ErrClosed, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Accept did not unblock after Close")
	}
}

func TestListener_CloseUnblocksDial(t *testing.T) {
	t.Parallel()

	ln := NewListener()
	// No goroutine calling Accept, so Dial should block.
	ln.Close()

	_, err := ln.Dial()
	if err != net.ErrClosed {
		t.Fatalf("expected net.ErrClosed, got %v", err)
	}
}

func TestListener_MultipleCloseIsSafe(t *testing.T) {
	t.Parallel()

	ln := NewListener()
	ln.Close()
	ln.Close() // must not panic
}

func TestListener_Addr(t *testing.T) {
	t.Parallel()

	ln := NewListener()
	defer ln.Close()

	addr := ln.Addr()
	if addr.Network() != "pipe" {
		t.Fatalf("Network() = %q, want %q", addr.Network(), "pipe")
	}
	if addr.String() != "pipe" {
		t.Fatalf("String() = %q, want %q", addr.String(), "pipe")
	}
}

func TestListener_MultipleConcurrentDials(t *testing.T) {
	t.Parallel()

	ln := NewListener()
	defer ln.Close()

	const n = 10
	var wg sync.WaitGroup

	// Server: accept n connections.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range n {
			conn, err := ln.Accept()
			if err != nil {
				t.Errorf("Accept #%d: %v", i, err)
				return
			}
			conn.Close()
		}
	}()

	// Clients: dial n times concurrently.
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := ln.Dial()
			if err != nil {
				t.Errorf("Dial: %v", err)
				return
			}
			conn.Close()
		}()
	}

	wg.Wait()
}
