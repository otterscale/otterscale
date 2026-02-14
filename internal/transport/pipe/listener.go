// Package pipe provides an in-memory net.Listener backed by
// net.Pipe. Only code with a reference to the Listener can inject
// connections via Dial, making it suitable for isolating an HTTP
// server so that it has no TCP presence.
package pipe

import (
	"net"
	"sync"
)

// Listener implements net.Listener using in-memory net.Pipe
// connections. It has no network presence; the only way to create a
// connection is to call Dial, which returns the client side of a
// net.Pipe pair while handing the server side to Accept.
type Listener struct {
	connCh chan net.Conn
	once   sync.Once
	done   chan struct{}
}

// NewListener returns a ready-to-use Listener.
func NewListener() *Listener {
	return &Listener{
		connCh: make(chan net.Conn),
		done:   make(chan struct{}),
	}
}

// Accept blocks until a new connection is available (created by Dial)
// or the listener is closed.
func (l *Listener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.connCh:
		return conn, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

// Close shuts down the listener. Any blocked Accept calls will
// return net.ErrClosed. Close is safe to call multiple times.
func (l *Listener) Close() error {
	l.once.Do(func() { close(l.done) })
	return nil
}

// Addr returns a synthetic address identifying this as a pipe listener.
func (l *Listener) Addr() net.Addr {
	return pipeAddr{}
}

// Dial creates a new in-memory connection pair via net.Pipe and
// hands the server side to Accept. It returns the client side to
// the caller. If the listener has been closed, both ends are
// cleaned up and net.ErrClosed is returned.
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

// pipeAddr is the net.Addr returned by Listener.Addr.
type pipeAddr struct{}

func (pipeAddr) Network() string { return "pipe" }
func (pipeAddr) String() string  { return "pipe" }
