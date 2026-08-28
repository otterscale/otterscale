package core

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

// TerminalSize is the domain-level terminal dimension type, keeping core free
// of k8s.io/client-go/tools/remotecommand. The adapter layer converts it to
// remotecommand.TerminalSize for SPDY executors.
type TerminalSize struct {
	Width  uint16
	Height uint16
}

// TerminalSizer provides the next terminal size event. Implementations block
// until one is available, or return nil once no more will be produced.
type TerminalSizer interface {
	Next() *TerminalSize
}

// TerminalSizeQueue is a concurrency-safe, single-slot mailbox implementing
// TerminalSizer: Set publishes, Next consumes.
//
// It holds one size because a terminal has one size: anything still waiting
// when a newer size arrives is already stale, and delivering it first would
// only make the consumer apply a dimension the terminal no longer has.
type TerminalSizeQueue struct {
	mu     sync.Mutex
	ch     chan TerminalSize
	closed bool
}

// terminalSizeQueueBuf is the mailbox capacity: the latest size only.
const terminalSizeQueueBuf = 1

func NewTerminalSizeQueue() *TerminalSizeQueue {
	return &TerminalSizeQueue{ch: make(chan TerminalSize, terminalSizeQueueBuf)}
}

// Next blocks until an event is available, and returns nil once the channel is
// closed.
func (q *TerminalSizeQueue) Next() *TerminalSize {
	size, ok := <-q.ch
	if !ok {
		return nil
	}
	return &size
}

// Set publishes a resize event, replacing any size the consumer has not taken
// yet. Calls after Close are ignored, to prevent a send-on-closed-channel panic.
//
// Every channel operation here is non-blocking, and that is the point. Set runs
// under q.mu, which Close also needs; an earlier version used a blocking
// receive to drop the oldest of a four-deep buffer, and a consumer that drained
// the channel in between left Set parked on that receive holding the lock. Close
// then blocked forever, and because it is the first deferred call in the exec
// goroutine, so did every close after it — the session's pipes stayed open, its
// Done channel never fired, and the reaper could not collect it.
func (q *TerminalSizeQueue) Set(width, height uint16) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return
	}

	// Discard a stale pending size, then publish. The send cannot fail: only
	// Set fills the slot and it holds the lock, so nothing can refill what was
	// just drained. The default guards the invariant, not a reachable case.
	select {
	case <-q.ch:
	default:
	}
	select {
	case q.ch <- TerminalSize{Width: width, Height: height}:
	default:
	}
}

// Close makes Next return nil. It is safe to call more than once.
func (q *TerminalSizeQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.closed {
		q.closed = true
		close(q.ch)
	}
}

type ExecSession struct {
	ID string
	// Owner is the authenticated subject that opened the session, and the only
	// one that may write to, resize, or close it.
	Owner     string
	Stdin     io.WriteCloser
	SizeQueue *TerminalSizeQueue
	Cancel    context.CancelFunc
	// Done is closed when the exec goroutine finishes. Closing, rather than
	// sending, is what makes the signal permanent: a value would be consumed by
	// whichever observer read it first, hiding the finished state from the reaper.
	Done chan struct{}
	// Err is safe to read once Done is closed.
	Err error
}

type PortForwardSession struct {
	ID     string
	Owner  string
	Writer io.WriteCloser
	Cancel context.CancelFunc
	// Done: see ExecSession.Done for why it is closed rather than sent to.
	Done chan struct{}
	// Err is safe to read once Done is closed.
	Err error
}

// VNCSession is an active VNC session to a KubeVirt VMI.
type VNCSession struct {
	ID     string
	Owner  string
	Writer io.WriteCloser
	Cancel context.CancelFunc
	Done   chan struct{}
	// Err is safe to read once Done is closed.
	Err error
}

// Session caps prevent resource exhaustion from clients that create sessions
// without cleaning them up.
const maxExecSessions = 100

const maxPortForwardSessions = 100

const maxVNCSessions = 100

// SessionStore manages active exec, port-forward, and VNC sessions.
type SessionStore struct {
	mu       sync.RWMutex
	execSess map[string]*ExecSession
	pfSess   map[string]*PortForwardSession
	vncSess  map[string]*VNCSession
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		execSess: make(map[string]*ExecSession),
		pfSess:   make(map[string]*PortForwardSession),
		vncSess:  make(map[string]*VNCSession),
	}
}

func (s *SessionStore) PutExec(sess *ExecSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.execSess) >= maxExecSessions {
		return &DomainError{
			Code:    ErrorCodeResourceExhausted,
			Message: fmt.Sprintf("max concurrent exec sessions (%d) reached", maxExecSessions),
		}
	}
	s.execSess[sess.ID] = sess
	return nil
}

func (s *SessionStore) GetExec(id string) (*ExecSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.execSess[id]
	return sess, ok
}

// RemoveExec atomically gets and removes a session, returning nil if it does
// not exist. Claiming ownership in one step is what prevents the double-close
// race between CleanupExec and ReapStaleSessions.
func (s *SessionStore) RemoveExec(id string) *ExecSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.execSess[id]
	if !ok {
		return nil
	}
	delete(s.execSess, id)
	return sess
}

func (s *SessionStore) PutPortForward(sess *PortForwardSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.pfSess) >= maxPortForwardSessions {
		return &DomainError{
			Code:    ErrorCodeResourceExhausted,
			Message: fmt.Sprintf("max concurrent port-forward sessions (%d) reached", maxPortForwardSessions),
		}
	}
	s.pfSess[sess.ID] = sess
	return nil
}

func (s *SessionStore) GetPortForward(id string) (*PortForwardSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.pfSess[id]
	return sess, ok
}

// RemovePortForward: see RemoveExec.
func (s *SessionStore) RemovePortForward(id string) *PortForwardSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.pfSess[id]
	if !ok {
		return nil
	}
	delete(s.pfSess, id)
	return sess
}

func (s *SessionStore) PutVNC(sess *VNCSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.vncSess) >= maxVNCSessions {
		return &DomainError{
			Code:    ErrorCodeResourceExhausted,
			Message: fmt.Sprintf("max concurrent VNC sessions (%d) reached", maxVNCSessions),
		}
	}
	s.vncSess[sess.ID] = sess
	return nil
}

func (s *SessionStore) GetVNC(id string) (*VNCSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.vncSess[id]
	return sess, ok
}

// RemoveVNC: see RemoveExec.
func (s *SessionStore) RemoveVNC(id string) *VNCSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.vncSess[id]
	if !ok {
		return nil
	}
	delete(s.vncSess, id)
	return sess
}

// ReapStaleSessions removes sessions whose goroutine has already finished,
// which is what keeps clients that disconnect without calling Cleanup from
// leaking sessions.
//
// Map mutations happen under the write lock, but the potentially blocking
// Cancel/Close calls happen after it is released — otherwise a goroutine
// holding a pipe write and waiting for a read lock would deadlock the reaper.
func (s *SessionStore) ReapStaleSessions() int {
	s.mu.Lock()

	var staleExec []*ExecSession
	for id, sess := range s.execSess {
		select {
		case <-sess.Done:
			staleExec = append(staleExec, sess)
			delete(s.execSess, id)
		default:
		}
	}

	var stalePF []*PortForwardSession
	for id, sess := range s.pfSess {
		select {
		case <-sess.Done:
			stalePF = append(stalePF, sess)
			delete(s.pfSess, id)
		default:
		}
	}

	var staleVNC []*VNCSession
	for id, sess := range s.vncSess {
		select {
		case <-sess.Done:
			staleVNC = append(staleVNC, sess)
			delete(s.vncSess, id)
		default:
		}
	}

	s.mu.Unlock()

	for _, sess := range staleExec {
		sess.Cancel()
		if err := sess.Stdin.Close(); err != nil {
			slog.Warn("failed to close exec stdin", "session", sess.ID, "error", err)
		}
	}
	for _, sess := range stalePF {
		sess.Cancel()
		if err := sess.Writer.Close(); err != nil {
			slog.Warn("failed to close port-forward writer", "session", sess.ID, "error", err)
		}
	}
	for _, sess := range staleVNC {
		sess.Cancel()
		if err := sess.Writer.Close(); err != nil {
			slog.Warn("failed to close VNC writer", "session", sess.ID, "error", err)
		}
	}

	return len(staleExec) + len(stalePF) + len(staleVNC)
}
