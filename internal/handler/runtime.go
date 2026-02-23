package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/otterscale/api/runtime/v1"
	"github.com/otterscale/otterscale/internal/core"
)

// streamChunkSize is the maximum bytes sent per streaming message.
const streamChunkSize = 32 * 1024

// RuntimeService implements the Runtime gRPC service. It proxies
// Kubernetes runtime operations (logs, exec, port-forward, scale,
// restart) through the tunnel.
type RuntimeService struct {
	pb.UnimplementedRuntimeServiceHandler

	runtime *core.RuntimeUseCase
}

// NewRuntimeService returns a RuntimeService backed by the given
// use-case.
func NewRuntimeService(runtime *core.RuntimeUseCase) *RuntimeService {
	return &RuntimeService{runtime: runtime}
}

var _ pb.RuntimeServiceHandler = (*RuntimeService)(nil)

// ---------------------------------------------------------------------------
// PodLog
// ---------------------------------------------------------------------------

// PodLog streams container log output to the client.
func (s *RuntimeService) PodLog(ctx context.Context, req *pb.PodLogRequest, stream *connect.ServerStream[pb.PodLogResponse]) error {
	opts := core.PodLogOptions{
		Container:  req.GetContainer(),
		Follow:     req.GetFollow(),
		Previous:   req.GetPrevious(),
		Timestamps: req.GetTimestamps(),
	}
	if req.HasTailLines() {
		v := req.GetTailLines()
		opts.TailLines = &v
	}
	if req.HasSinceSeconds() {
		v := req.GetSinceSeconds()
		opts.SinceSeconds = &v
	}
	if req.HasSinceTime() {
		t := req.GetSinceTime().AsTime()
		opts.SinceTime = &t
	}
	if req.HasLimitBytes() {
		v := req.GetLimitBytes()
		opts.LimitBytes = &v
	}

	reader, err := s.runtime.StartPodLogs(ctx, req.GetCluster(), req.GetNamespace(), req.GetName(), opts)
	if err != nil {
		return domainErrorToConnectError(err)
	}
	defer reader.Close()

	buf := make([]byte, streamChunkSize)
	for {
		n, readErr := reader.Read(buf)
		if n > 0 {
			msg := &pb.PodLogResponse{}
			msg.SetData(append([]byte(nil), buf[:n]...))
			if err := stream.Send(msg); err != nil {
				return err
			}
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return nil
			}
			return domainErrorToConnectError(readErr)
		}
	}
}

// ---------------------------------------------------------------------------
// ExecuteTTY / WriteTTY / ResizeTTY
// ---------------------------------------------------------------------------

// ExecuteTTY starts an interactive exec session and streams
// stdout/stderr back to the client. The first response message
// contains the session_id that the client must use for WriteTTY
// and ResizeTTY calls.
func (s *RuntimeService) ExecuteTTY(ctx context.Context, req *pb.ExecuteTTYRequest, stream *connect.ServerStream[pb.ExecuteTTYResponse]) error {
	rows := req.GetRows()
	cols := req.GetCols()
	if rows > math.MaxUint16 || cols > math.MaxUint16 {
		return connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("terminal dimensions out of range (max %d)", math.MaxUint16))
	}

	sess, stdoutR, stderrR, err := s.runtime.StartExec(ctx, &core.StartExecParams{
		Cluster:   req.GetCluster(),
		Namespace: req.GetNamespace(),
		Name:      req.GetName(),
		Container: req.GetContainer(),
		Command:   req.GetCommand(),
		TTY:       req.GetTty(),
		Rows:      uint16(rows),
		Cols:      uint16(cols),
	})
	if err != nil {
		return domainErrorToConnectError(err)
	}
	defer s.runtime.CleanupExec(ctx, sess.ID)

	// Send the session ID as the first message.
	first := &pb.ExecuteTTYResponse{}
	first.SetSessionId(sess.ID)
	if err := stream.Send(first); err != nil {
		return err
	}

	ch := mergeExecOutputs(ctx, stdoutR, stderrR)

	// Stream chunks to the client until all output is consumed.
	// The channel is closed by the readerWg goroutine once both
	// stdout and stderr readers exit (triggered by pipe closure
	// when the exec session ends or CleanupExec runs). This
	// guarantees all buffered data is delivered without relying on
	// a time-based heuristic.
	for {
		select {
		case <-ctx.Done():
			return nil

		case c, ok := <-ch:
			if !ok {
				return nil
			}
			msg := &pb.ExecuteTTYResponse{}
			if len(c.stdout) > 0 {
				msg.SetStdout(c.stdout)
			}
			if len(c.stderr) > 0 {
				msg.SetStderr(c.stderr)
			}
			if err := stream.Send(msg); err != nil {
				return err
			}
		}
	}
}

// execChunk holds a piece of stdout or stderr data from an exec session.
type execChunk struct {
	stdout []byte
	stderr []byte
}

// execChunkBufSize is the channel buffer size for merged exec output.
const execChunkBufSize = 8

// mergeExecOutputs reads stdout and stderr concurrently and merges
// them into a single channel. The returned channel is closed once
// both readers finish.
func mergeExecOutputs(ctx context.Context, stdoutR, stderrR io.ReadCloser) <-chan execChunk {
	ch := make(chan execChunk, execChunkBufSize)
	var readerWg sync.WaitGroup
	readerWg.Add(2)

	go func() {
		defer readerWg.Done()
		defer stdoutR.Close()
		buf := make([]byte, streamChunkSize)
		for {
			n, readErr := stdoutR.Read(buf)
			if n > 0 {
				select {
				case ch <- execChunk{stdout: append([]byte(nil), buf[:n]...)}:
				case <-ctx.Done():
					return
				}
			}
			if readErr != nil {
				return
			}
		}
	}()

	go func() {
		defer readerWg.Done()
		defer stderrR.Close()
		buf := make([]byte, streamChunkSize)
		for {
			n, readErr := stderrR.Read(buf)
			if n > 0 {
				select {
				case ch <- execChunk{stderr: append([]byte(nil), buf[:n]...)}:
				case <-ctx.Done():
					return
				}
			}
			if readErr != nil {
				return
			}
		}
	}()

	go func() {
		readerWg.Wait()
		close(ch)
	}()

	return ch
}

// WriteTTY sends stdin data to an active exec session.
func (s *RuntimeService) WriteTTY(ctx context.Context, req *pb.WriteTTYRequest) (*emptypb.Empty, error) {
	if err := s.runtime.WriteExec(ctx, req.GetSessionId(), req.GetStdin()); err != nil {
		return nil, domainErrorToConnectError(err)
	}
	return &emptypb.Empty{}, nil
}

// ResizeTTY updates the terminal dimensions of an active exec session.
func (s *RuntimeService) ResizeTTY(ctx context.Context, req *pb.ResizeTTYRequest) (*emptypb.Empty, error) {
	rows := req.GetRows()
	cols := req.GetCols()
	if rows > math.MaxUint16 || cols > math.MaxUint16 {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("terminal dimensions out of range (max %d)", math.MaxUint16))
	}
	if err := s.runtime.ResizeExec(ctx, req.GetSessionId(), uint16(rows), uint16(cols)); err != nil {
		return nil, domainErrorToConnectError(err)
	}
	return &emptypb.Empty{}, nil
}

// ---------------------------------------------------------------------------
// PortForward / WritePortForward
// ---------------------------------------------------------------------------

// PortForward opens a port-forward session and streams data from the
// pod back to the client. The first response message contains the
// session_id that the client must use for WritePortForward calls.
func (s *RuntimeService) PortForward(ctx context.Context, req *pb.PortForwardRequest, stream *connect.ServerStream[pb.PortForwardResponse]) error {
	sess, dataOutR, err := s.runtime.StartPortForward(
		ctx,
		req.GetCluster(),
		req.GetNamespace(),
		req.GetName(),
		req.GetPort(),
	)
	if err != nil {
		return domainErrorToConnectError(err)
	}
	defer s.runtime.CleanupPortForward(ctx, sess.ID)

	// Send the session ID as the first message.
	first := &pb.PortForwardResponse{}
	first.SetSessionId(sess.ID)
	if err := stream.Send(first); err != nil {
		return err
	}

	// Stream data from the pod.
	buf := make([]byte, streamChunkSize)
	for {
		n, readErr := dataOutR.Read(buf)
		if n > 0 {
			msg := &pb.PortForwardResponse{}
			msg.SetData(append([]byte(nil), buf[:n]...))
			if err := stream.Send(msg); err != nil {
				return err
			}
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return nil
			}
			return domainErrorToConnectError(readErr)
		}
	}
}

// WritePortForward sends data to an active port-forward session.
func (s *RuntimeService) WritePortForward(ctx context.Context, req *pb.WritePortForwardRequest) (*emptypb.Empty, error) {
	if err := s.runtime.WritePortForward(ctx, req.GetSessionId(), req.GetData()); err != nil {
		return nil, domainErrorToConnectError(err)
	}
	return &emptypb.Empty{}, nil
}

// ---------------------------------------------------------------------------
// Scale
// ---------------------------------------------------------------------------

// Scale updates the replica count and returns the new value.
func (s *RuntimeService) Scale(ctx context.Context, req *pb.ScaleRequest) (*pb.ScaleResponse, error) {
	replicas, err := s.runtime.Scale(
		ctx,
		&core.ResourceIdentifier{
			Cluster:   req.GetCluster(),
			Group:     req.GetGroup(),
			Version:   req.GetVersion(),
			Resource:  req.GetResource(),
			Namespace: req.GetNamespace(),
			Name:      req.GetName(),
		},
		req.GetReplicas(),
	)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	resp := &pb.ScaleResponse{}
	resp.SetReplicas(replicas)
	return resp, nil
}

// ---------------------------------------------------------------------------
// Restart
// ---------------------------------------------------------------------------

// Restart triggers a rolling restart of a workload.
func (s *RuntimeService) Restart(ctx context.Context, req *pb.RestartRequest) (*emptypb.Empty, error) {
	if err := s.runtime.Restart(
		ctx,
		&core.ResourceIdentifier{
			Cluster:   req.GetCluster(),
			Group:     req.GetGroup(),
			Version:   req.GetVersion(),
			Resource:  req.GetResource(),
			Namespace: req.GetNamespace(),
			Name:      req.GetName(),
		},
	); err != nil {
		return nil, domainErrorToConnectError(err)
	}
	return &emptypb.Empty{}, nil
}
