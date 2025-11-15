package vnc

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	kvcorev1 "github.com/otterscale/kubevirt-client-go/kubevirt/typed/core/v1"
)

//nolint:revive // allows this exported interface name for specific domain clarity.
type VNCRepo interface {
	Streamer(scope, namespace, name string) (kvcorev1.StreamInterface, error)
}

type UseCase struct {
	vnc VNCRepo

	vncSessions sync.Map
}

func NewUseCase(vnc VNCRepo) *UseCase {
	return &UseCase{
		vnc: vnc,
	}
}

func (uc *UseCase) CreateVNCSession(scope, namespace, name string) (string, error) {
	sessionID := uuid.New().String()

	vnc, err := uc.vnc.Streamer(scope, namespace, name)
	if err != nil {
		return "", err
	}

	uc.vncSessions.Store(sessionID, vnc)

	return sessionID, nil
}

func (uc *UseCase) VNCPathPrefix() string {
	return "/vnc/"
}

func (uc *UseCase) VNCHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade to websocket
	upgrader := kvcorev1.NewUpgrader()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// get vnc session
	sessionID := strings.TrimPrefix(r.URL.Path, uc.VNCPathPrefix())

	if sessionID == "" {
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "missing VNC session ID"))
		return
	}

	// load vnc session
	value, ok := uc.vncSessions.Load(sessionID)
	if !ok {
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "VNC session not found"))
		return
	}
	defer uc.vncSessions.Delete(sessionID)

	// configure websocket connection
	const wait = 5 * time.Minute

	_ = conn.SetReadDeadline(time.Now().Add(wait))

	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(wait))
		return nil
	})

	// create context for ping handler
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start ping handler
	go uc.pingHandler(ctx, conn)

	pipeInReader, pipeInWriter := io.Pipe()
	pipeOutReader, pipeOutWriter := io.Pipe()
	defer pipeInWriter.Close()
	defer pipeOutWriter.Close()

	// start stream handler
	stream := value.(kvcorev1.StreamInterface)
	errChan := make(chan error, 3)
	go uc.streamHandler(stream, pipeInReader, pipeOutWriter, conn, pipeOutReader, pipeInWriter, errChan)

	// wait for any handler to complete
	finalErr := <-errChan

	if finalErr != nil && !errors.Is(finalErr, context.Canceled) && finalErr != io.EOF && !websocket.IsCloseError(finalErr, websocket.CloseNoStatusReceived, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, finalErr.Error()))
	}
}

func (uc *UseCase) pingHandler(ctx context.Context, conn *websocket.Conn) {
	const (
		period   = 1 * time.Minute
		deadline = 10 * time.Second
	)

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(deadline)); err != nil {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (uc *UseCase) streamHandler(stream kvcorev1.StreamInterface, inReader io.Reader, outWriter io.Writer, conn *websocket.Conn, outReader io.Reader, inWriter io.Writer, errChan chan error) {
	go func() {
		errChan <- stream.Stream(kvcorev1.StreamOptions{
			In:  inReader,
			Out: outWriter,
		})
	}()

	go func() {
		_, err := kvcorev1.CopyTo(conn, outReader)
		errChan <- err
	}()

	go func() {
		_, err := kvcorev1.CopyFrom(inWriter, conn)
		errChan <- err
	}()
}
