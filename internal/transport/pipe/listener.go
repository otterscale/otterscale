// Package pipe provides an in-memory net.Listener backed by
// net.Pipe. Only code with a reference to the Listener can inject
// connections via Dial, making it suitable for isolating an HTTP
// server so that it has no TCP presence.
package pipe

import (
	"net"
	"sync"
)

// Listener has no network presence: the only way to create a connection is
// Dial, which returns the client side of a net.Pipe pair and hands the server
// side to Accept.
type Listener struct {
	connCh chan net.Conn
	once   sync.Once
	done   chan struct{}
}

func NewListener() *Listener {
	return &Listener{
		connCh: make(chan net.Conn),
		done:   make(chan struct{}),
	}
}

// Accept blocks until Dial creates a connection or the listener closes.
func (l *Listener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.connCh:
		return conn, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

// Close unblocks any pending Accept with net.ErrClosed. Safe to call more than
// once.
func (l *Listener) Close() error {
	l.once.Do(func() { close(l.done) })
	return nil
}

func (l *Listener) Addr() net.Addr {
	return pipeAddr{}
}

// Dial hands the server side of a net.Pipe pair to Accept and returns the
// client side. Once the listener is closed, both ends are cleaned up and
// net.ErrClosed returned.
func (l *Listener) Dial() (net.Conn, error) {
	server, client := net.Pipe()
	select {
	case l.connCh <- server:
		return client, nil
	case <-l.done:
		server.Close()
		client.Close()
		return nil, net.ErrClosed
	}
}

type pipeAddr struct{}

func (pipeAddr) Network() string { return "pipe" }
func (pipeAddr) String() string  { return "pipe" }
