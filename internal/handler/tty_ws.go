package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/otterscale/otterscale/internal/core"
)

// ---------------------------------------------------------------------------
// TTY over WebSocket
// ---------------------------------------------------------------------------
//
// ExecuteTTY streams output but has no way to carry stdin, so its callers must
// send every keystroke as a separate unary WriteTTY call. That costs a full
// authenticated round trip per character and, because the calls are
// independent, lets concurrent writes reach the stdin pipe out of order.
//
// This endpoint carries the whole session on one connection instead. The frame
// format follows ttyd (https://github.com/tsl0922/ttyd): the first message is
// JSON describing the session, and every message after it is binary with a
// one-byte opcode prefix.
//
//	client -> server   '0' stdin bytes
//	                   '1' resize, payload is {"cols":N,"rows":N}
//	                   '2' pause output
//	                   '3' resume output
//	server -> client   '0' stdout bytes
//	                   '2' error text
//
// The ExecuteTTY/WriteTTY/ResizeTTY RPCs remain for dashboards that predate
// this endpoint, since the two components are deployed independently.

const (
	ttyOpInput  byte = '0'
	ttyOpResize byte = '1'
	ttyOpPause  byte = '2'
	ttyOpResume byte = '3'

	ttyOpOutput byte = '0'
	ttyOpError  byte = '2'
)

const (
	// ttyHandshakeTimeout bounds how long we wait for the opening JSON message.
	ttyHandshakeTimeout = 10 * time.Second
	// ttyPingInterval and ttyPongTimeout detect connections dropped by an
	// intermediary without a close frame. Browsers answer pings automatically.
	ttyPingInterval = 30 * time.Second
	ttyPongTimeout  = 60 * time.Second
	ttyWriteTimeout = 10 * time.Second
	// ttyMaxMessageSize caps a single client frame. Stdin is tiny; the headroom
	// is for large pastes.
	ttyMaxMessageSize = 1 << 20 // 1 MiB
)

// defaultTTYCommand is used when the client does not specify one.
var defaultTTYCommand = []string{"/bin/sh"}

var ttyUpgrader = websocket.Upgrader{
	HandshakeTimeout: ttyHandshakeTimeout,
	ReadBufferSize:   4 * 1024,
	WriteBufferSize:  streamChunkSize,
	// Requests reach this endpoint through the dashboard's own server-side
	// proxy, which terminates the browser connection and opens a fresh one
	// here, so there is no meaningful browser Origin to match against.
	// Rejecting cross-origin handshakes would only guard against
	// cross-site WebSocket hijacking, which needs ambient credentials to
	// work; this endpoint authenticates with a bearer token that a foreign
	// page cannot attach.
	CheckOrigin: func(_ *http.Request) bool { return true },
}

// ttyStartMessage is the opening JSON frame that describes the session.
type ttyStartMessage struct {
	Cluster   string   `json:"cluster"`
	Namespace string   `json:"namespace"`
	Pod       string   `json:"pod"`
	Container string   `json:"container"`
	Command   []string `json:"command"`
	Cols      uint16   `json:"cols"`
	Rows      uint16   `json:"rows"`
}

// ttyResizeMessage is the payload of a resize frame.
type ttyResizeMessage struct {
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

// ServeTTYWebSocket upgrades the request and runs an interactive exec session
// for the lifetime of the connection.
func (s *RuntimeService) ServeTTYWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ttyUpgrader.Upgrade(w, r, nil)
	if err != nil {
		// Upgrade has already written an error response.
		slog.Debug("tty websocket upgrade failed", "error", err)
		return
	}

	c := &ttyConn{conn: conn}
	defer conn.Close()

	if err := c.serve(r.Context(), s.runtime); err != nil {
		slog.Debug("tty websocket session ended", "error", err)
	}
}

// ttyConn owns one WebSocket connection. Gorilla allows a single concurrent
// writer, so every write goes through the mutex.
type ttyConn struct {
	conn *websocket.Conn
	mu   sync.Mutex
	flow ttyFlow
}

// serve runs the session until either side goes away.
func (c *ttyConn) serve(ctx context.Context, runtime *core.RuntimeUseCase) error {
	start, err := c.readStart()
	if err != nil {
		// Close cleanly even on failure: the client reads a normal close as
		// "this session is over" and anything else as "the connection dropped,
		// retry". Retrying a bad request would never succeed.
		c.shutdown(err)
		return err
	}

	sess, stdoutReader, stderrReader, err := runtime.StartExec(ctx, &core.StartExecParams{
		Cluster:   start.Cluster,
		Namespace: start.Namespace,
		Name:      start.Pod,
		Container: start.Container,
		Command:   start.Command,
		TTY:       true,
		Rows:      start.Rows,
		Cols:      start.Cols,
	})
	if err != nil {
		c.shutdown(err)
		return err
	}
	defer runtime.CleanupExec(ctx, sess.ID)
	// A TTY session merges stderr into stdout, so StartExec never writes to
	// this pipe. Release it immediately rather than pumping it.
	stderrReader.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan error, 2)
	go func() { done <- c.pumpOutput(ctx, stdoutReader) }()
	go func() { done <- c.pumpInput(sess) }()
	go c.pumpPing(ctx)

	// Whichever side finishes first ends the session.
	err = <-done
	cancel()
	// Unblock the peer: the output pump may be parked in Read, and the input
	// pump in ReadMessage.
	stdoutReader.Close()
	c.shutdown(err)
	<-done

	return err
}

// readStart reads and validates the opening JSON frame.
func (c *ttyConn) readStart() (*ttyStartMessage, error) {
	c.conn.SetReadLimit(ttyMaxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(ttyHandshakeTimeout)); err != nil {
		return nil, err
	}

	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var start ttyStartMessage
	if err := json.Unmarshal(data, &start); err != nil {
		return nil, fmt.Errorf("invalid start message: %w", err)
	}
	if start.Pod == "" {
		return nil, errors.New("pod name is required")
	}
	if len(start.Command) == 0 {
		start.Command = defaultTTYCommand
	}

	// From here on the read deadline is driven by the ping/pong exchange.
	if err := c.conn.SetReadDeadline(time.Now().Add(ttyPongTimeout)); err != nil {
		return nil, err
	}
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(ttyPongTimeout))
	})

	return &start, nil
}

// pumpInput relays client frames into the exec session. Running in a single
// goroutine is what keeps stdin ordered — the RPC path writes each request
// from its own goroutine and can interleave them.
func (c *ttyConn) pumpInput(sess *core.ExecSession) error {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return nil
			}
			return err
		}
		if len(data) == 0 {
			continue
		}

		payload := data[1:]
		switch data[0] {
		case ttyOpInput:
			if _, err := sess.Stdin.Write(payload); err != nil {
				return err
			}
		case ttyOpResize:
			var size ttyResizeMessage
			if err := json.Unmarshal(payload, &size); err != nil {
				// A malformed resize is not worth dropping the session over;
				// the next one corrects the grid.
				slog.Debug("tty websocket bad resize frame", "error", err)
				continue
			}
			sess.SizeQueue.Set(size.Cols, size.Rows)
		case ttyOpPause:
			c.flow.pause()
		case ttyOpResume:
			c.flow.resume()
		default:
			slog.Debug("tty websocket unknown opcode", "opcode", data[0])
		}
	}
}

// pumpOutput relays container output to the client, honouring pause requests.
func (c *ttyConn) pumpOutput(ctx context.Context, reader io.Reader) error {
	buf := make([]byte, streamChunkSize)
	frame := make([]byte, 0, streamChunkSize+1)

	for {
		if err := c.flow.wait(ctx); err != nil {
			return err
		}

		n, readErr := reader.Read(buf)
		if n > 0 {
			frame = append(frame[:0], ttyOpOutput)
			frame = append(frame, buf[:n]...)
			if err := c.write(websocket.BinaryMessage, frame); err != nil {
				return err
			}
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) || errors.Is(readErr, io.ErrClosedPipe) {
				return nil
			}
			return readErr
		}
	}
}

// pumpPing keeps the connection alive through idle-timeout proxies.
func (c *ttyConn) pumpPing(ctx context.Context) {
	ticker := time.NewTicker(ttyPingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *ttyConn) write(messageType int, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.conn.SetWriteDeadline(time.Now().Add(ttyWriteTimeout)); err != nil {
		return err
	}
	return c.conn.WriteMessage(messageType, data)
}

func (c *ttyConn) writeError(message string) {
	frame := append([]byte{ttyOpError}, message...)
	if err := c.write(websocket.BinaryMessage, frame); err != nil {
		slog.Debug("tty websocket error frame not delivered", "error", err)
	}
}

// shutdown reports the outcome and closes the connection. Detail travels in an
// error frame rather than the close reason, which is capped at 125 bytes; the
// close code itself only has to tell the client whether to reconnect.
func (c *ttyConn) shutdown(cause error) {
	if cause != nil {
		c.writeError(cause.Error())
	}

	c.mu.Lock()
	deadline := time.Now().Add(ttyWriteTimeout)
	message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	if err := c.conn.WriteControl(websocket.CloseMessage, message, deadline); err != nil {
		slog.Debug("tty websocket close frame not delivered", "error", err)
	}
	c.mu.Unlock()

	c.conn.Close()
}

// ttyFlow implements the pause/resume half of ttyd's flow control. The client
// asks for a pause when its renderer falls behind a flooding process, and the
// output pump parks until it catches up.
type ttyFlow struct {
	mu     sync.Mutex
	paused chan struct{}
}

func (f *ttyFlow) pause() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.paused == nil {
		f.paused = make(chan struct{})
	}
}

func (f *ttyFlow) resume() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.paused != nil {
		close(f.paused)
		f.paused = nil
	}
}

// wait blocks while output is paused.
func (f *ttyFlow) wait(ctx context.Context) error {
	f.mu.Lock()
	paused := f.paused
	f.mu.Unlock()

	if paused == nil {
		return nil
	}

	select {
	case <-paused:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
