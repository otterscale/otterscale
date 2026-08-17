package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// dialTTY starts a server that runs fn against an upgraded connection and
// returns a client connected to it.
func dialTTY(t *testing.T, fn func(*ttyConn)) *websocket.Conn {
	t.Helper()

	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(done)
		conn, err := ttyUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v", err)
			return
		}
		defer conn.Close()
		fn(&ttyConn{conn: conn})
	}))
	t.Cleanup(func() {
		server.Close()
		<-done
	})

	client, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/", nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	return client
}

func TestTTYConnReadStart(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		wantErr     bool
		wantCommand []string
		wantCols    uint16
	}{
		{
			name:        "full message",
			message:     `{"cluster":"c1","namespace":"ns","pod":"p","container":"app","command":["/bin/bash"],"cols":120,"rows":40}`,
			wantCommand: []string{"/bin/bash"},
			wantCols:    120,
		},
		{
			name:        "command defaults to a shell",
			message:     `{"pod":"p","cols":80,"rows":24}`,
			wantCommand: defaultTTYCommand,
			wantCols:    80,
		},
		{
			name:    "pod is required",
			message: `{"namespace":"ns"}`,
			wantErr: true,
		},
		{
			name:    "malformed json",
			message: `not json`,
			wantErr: true,
		},
		{
			name:    "dimensions out of range",
			message: `{"pod":"p","cols":99999,"rows":24}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type result struct {
				start *ttyStartMessage
				err   error
			}
			results := make(chan result, 1)

			client := dialTTY(t, func(c *ttyConn) {
				start, err := c.readStart()
				results <- result{start, err}
			})

			if err := client.WriteMessage(websocket.TextMessage, []byte(tt.message)); err != nil {
				t.Fatalf("write start: %v", err)
			}

			got := <-results
			if tt.wantErr {
				if got.err == nil {
					t.Fatalf("expected an error, got start %+v", got.start)
				}
				return
			}
			if got.err != nil {
				t.Fatalf("readStart: %v", got.err)
			}
			if got.start.Cols != tt.wantCols {
				t.Errorf("cols = %d, want %d", got.start.Cols, tt.wantCols)
			}
			if strings.Join(got.start.Command, " ") != strings.Join(tt.wantCommand, " ") {
				t.Errorf("command = %v, want %v", got.start.Command, tt.wantCommand)
			}
		})
	}
}

// TestTTYConnOutputFraming checks the wire contract the dashboard depends on:
// every output message is one binary frame prefixed with the output opcode.
func TestTTYConnOutputFraming(t *testing.T) {
	reader, writer := io.Pipe()
	defer writer.Close()

	client := dialTTY(t, func(c *ttyConn) {
		if err := c.pumpOutput(context.Background(), reader); err != nil {
			t.Errorf("pumpOutput: %v", err)
		}
	})

	if _, err := writer.Write([]byte("hello")); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := client.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatalf("set deadline: %v", err)
	}
	messageType, data, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if messageType != websocket.BinaryMessage {
		t.Errorf("message type = %d, want binary", messageType)
	}
	if len(data) == 0 || data[0] != ttyOpOutput {
		t.Fatalf("frame = %q, want the output opcode prefix", data)
	}
	if string(data[1:]) != "hello" {
		t.Errorf("payload = %q, want %q", data[1:], "hello")
	}
}

func TestTTYFlowPauseResume(t *testing.T) {
	var flow ttyFlow

	// Nothing to wait for until a pause is requested.
	if err := flow.wait(context.Background()); err != nil {
		t.Fatalf("wait while running: %v", err)
	}

	flow.pause()

	blocked := make(chan error, 1)
	go func() { blocked <- flow.wait(context.Background()) }()

	select {
	case err := <-blocked:
		t.Fatalf("wait returned %v while paused", err)
	case <-time.After(50 * time.Millisecond):
	}

	flow.resume()

	select {
	case err := <-blocked:
		if err != nil {
			t.Fatalf("wait after resume: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("wait did not return after resume")
	}
}

func TestTTYFlowWaitCancelled(t *testing.T) {
	var flow ttyFlow
	flow.pause()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := flow.wait(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("wait = %v, want context.Canceled", err)
	}
}

// TestTTYResizeMessageRoundTrip pins the resize payload shape, which the
// dashboard encodes by hand.
func TestTTYResizeMessageRoundTrip(t *testing.T) {
	var size ttyResizeMessage
	if err := json.Unmarshal([]byte(`{"cols":100,"rows":30}`), &size); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if size.Cols != 100 || size.Rows != 30 {
		t.Errorf("size = %+v, want cols 100 rows 30", size)
	}
}
