package core

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	kvcorev1 "github.com/otterscale/kubevirt-client-go/kubevirt/typed/core/v1"
)

func (uc *InstanceUseCase) CreateVNCSession(ctx context.Context, scopeUUID, facility, namespace, name string) (string, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scopeUUID, facility)
	if err != nil {
		return "", err
	}
	sessionID := uuid.New().String()
	vnc, err := uc.kubeVirt.VNCInstance(config, namespace, name)
	if err != nil {
		return "", err
	}
	uc.vncSessionMap.Store(sessionID, vnc)
	return sessionID, nil
}

func (uc *InstanceUseCase) VNCPathPrefix() string {
	return "/vnc/"
}

func (uc *InstanceUseCase) VNCHandler(w http.ResponseWriter, r *http.Request) {
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
	value, ok := uc.vncSessionMap.Load(sessionID)
	if !ok {
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "VNC session not found"))
		return
	}
	defer uc.vncSessionMap.Delete(sessionID)

	// configure websocket connection
	wait := 5 * time.Minute
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
	stream := value.(VirtualMachineStream)
	errChan := make(chan error, 3)
	go uc.streamHandler(stream, pipeInReader, pipeOutWriter, conn, pipeOutReader, pipeInWriter, errChan)

	// wait for any handler to complete
	finalErr := <-errChan

	if finalErr != nil && !errors.Is(finalErr, context.Canceled) && finalErr != io.EOF && !websocket.IsCloseError(finalErr, websocket.CloseNoStatusReceived, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, finalErr.Error()))
	}
}

func (uc *InstanceUseCase) pingHandler(ctx context.Context, conn *websocket.Conn) {
	var (
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

func (uc *InstanceUseCase) streamHandler(vnc VirtualMachineStream, inReader io.Reader, outWriter io.Writer, conn *websocket.Conn, outReader io.Reader, inWriter io.Writer, errChan chan error) {
	go func() {
		errChan <- vnc.Stream(VirtualMachineStreamOptions{
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
