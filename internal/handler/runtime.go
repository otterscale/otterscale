package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"sync"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/otterscale/otterscale/api/runtime/v1"

	"github.com/otterscale/otterscale/internal/core"
)

// streamChunkSize is the maximum bytes sent per streaming message.
const streamChunkSize = 32 * 1024

// sessionOutcomeTimeout bounds the wait for a finished session to
// publish its result. Once a session's output pipes are drained the
// goroutine behind it has only a few deferred closes left to run, so
// this is a backstop against a wedged adapter rather than a normal
// path — without it, a session that never signals completion would hold
// its RPC open until the client gave up.
const sessionOutcomeTimeout = 5 * time.Second

// sessionOutcome converts the error a finished session ended with into
// the error its RPC should report.
//
// Two endings are not RPC failures. A caller whose context was canceled
// is already gone, and the session was torn down on its behalf. And a
// command that ran and exited non-zero delivered exactly what was asked
// for: its output has already been streamed, so reporting the exit
// status as a transport failure would make every unsuccessful shell
// command look like a server error. The status itself has no field in
// the current response schema, so it is logged instead.
func sessionOutcome(kind, id string, err error) error {
	if err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil
	}

	var exited *core.ErrCommandExited
	if errors.As(err, &exited) {
		slog.Debug("exec command exited non-zero", "session", id, "exit_code", exited.Code)
		return nil
	}

	// Log as well as return: before this, a session that failed to even
	// start left no trace on either side.
	slog.Warn("session ended with an error", "kind", kind, "session", id, "error", err)
	return domainErrorToConnectError(err)
}

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

	for c := range readChunks(ctx, reader) {
		if c.err != nil {
			return domainErrorToConnectError(c.err)
		}
		if len(c.data) == 0 {
			continue
		}
		msg := &pb.PodLogResponse{}
		msg.SetData(c.data)
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	return nil
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
				// Output exhausted; the session's own result is what
				// says whether it ever ran. Without this, a session
				// that failed to start looks identical to one that
				// produced no output.
				return s.execOutcome(ctx, sess)
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

// execOutcome waits for the exec session to finish and reports why it
// ended.
func (s *RuntimeService) execOutcome(ctx context.Context, sess *core.ExecSession) error {
	waitCtx, cancel := context.WithTimeout(ctx, sessionOutcomeTimeout)
	defer cancel()
	return sessionOutcome("exec", sess.ID, s.runtime.WaitExec(waitCtx, sess))
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
	for c := range readChunks(ctx, dataOutR) {
		if c.err != nil {
			return domainErrorToConnectError(c.err)
		}
		if len(c.data) == 0 {
			continue
		}
		msg := &pb.PortForwardResponse{}
		msg.SetData(c.data)
		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	// The pipe reaching EOF says the session is over, not that it
	// succeeded — a forward kubelet refused outright ends the same way.
	waitCtx, cancel := context.WithTimeout(ctx, sessionOutcomeTimeout)
	defer cancel()
	return sessionOutcome("port-forward", sess.ID, s.runtime.WaitPortForward(waitCtx, sess))
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

// ---------------------------------------------------------------------------
// VNC / WriteVNC
// ---------------------------------------------------------------------------

// VNC opens a VNC session to a KubeVirt VMI and streams raw VNC
// protocol data back to the client. The first response message
// contains the session_id that the client must use for WriteVNC calls.
func (s *RuntimeService) VNC(ctx context.Context, req *pb.VNCRequest, stream *connect.ServerStream[pb.VNCResponse]) error {
	sess, dataOutR, err := s.runtime.StartVNC(
		ctx,
		req.GetCluster(),
		req.GetNamespace(),
		req.GetName(),
	)
	if err != nil {
		return domainErrorToConnectError(err)
	}
	defer s.runtime.CleanupVNC(ctx, sess.ID)

	// Send the session ID as the first message.
	first := &pb.VNCResponse{}
	first.SetSessionId(sess.ID)
	if err := stream.Send(first); err != nil {
		return err
	}

	// Stream VNC data from the VMI.
	for c := range readChunks(ctx, dataOutR) {
		if c.err != nil {
			return domainErrorToConnectError(c.err)
		}
		if len(c.data) == 0 {
			continue
		}
		msg := &pb.VNCResponse{}
		msg.SetData(c.data)
		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	waitCtx, cancel := context.WithTimeout(ctx, sessionOutcomeTimeout)
	defer cancel()
	return sessionOutcome("vnc", sess.ID, s.runtime.WaitVNC(waitCtx, sess))
}

// chunk carries one streamed payload, or the terminal read error once
// the reader is exhausted. A clean EOF is reported as a chunk with no
// data and no error.
type chunk struct {
	data []byte
	err  error
}

// readChunks streams r in fixed-size chunks on the returned channel and
// closes it once the reader is exhausted or ctx is canceled. Reading in
// a goroutine is what makes cancellation observable: Read on an
// io.Pipe blocks until the writer side is closed, so a handler that
// called it directly would keep running after its client vanished —
// and never run the cleanup that releases the session.
func readChunks(ctx context.Context, r io.Reader) <-chan chunk {
	ch := make(chan chunk, 1)

	go func() {
		defer close(ch)

		buf := make([]byte, streamChunkSize)
		for {
			n, readErr := r.Read(buf)
			if n > 0 {
				select {
				case ch <- chunk{data: append([]byte(nil), buf[:n]...)}:
				case <-ctx.Done():
					return
				}
			}
			if readErr != nil {
				if errors.Is(readErr, io.EOF) {
					readErr = nil
				}
				select {
				case ch <- chunk{err: readErr}:
				case <-ctx.Done():
				}
				return
			}
		}
	}()

	return ch
}

// WriteVNC sends VNC data to an active VNC session.
func (s *RuntimeService) WriteVNC(ctx context.Context, req *pb.WriteVNCRequest) (*emptypb.Empty, error) {
	if err := s.runtime.WriteVNC(ctx, req.GetSessionId(), req.GetData()); err != nil {
		return nil, domainErrorToConnectError(err)
	}
	return &emptypb.Empty{}, nil
}

// ---------------------------------------------------------------------------
// SubResourceAction
// ---------------------------------------------------------------------------

// SubResourceAction performs a generic subresource action (e.g.
// KubeVirt VM start/stop/restart). The request is forwarded to
// kube-apiserver via impersonation, so Kubernetes RBAC is enforced.
func (s *RuntimeService) SubResourceAction(ctx context.Context, req *pb.SubResourceActionRequest) (*pb.SubResourceActionResponse, error) {
	result, err := s.runtime.SubResourceAction(
		ctx,
		&core.ResourceIdentifier{
			Cluster:     req.GetCluster(),
			Group:       req.GetGroup(),
			Version:     req.GetVersion(),
			Resource:    req.GetResource(),
			SubResource: req.GetSubresource(),
			Namespace:   req.GetNamespace(),
			Name:        req.GetName(),
		},
		req.GetMethod(),
		req.GetBody(),
	)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	resp := &pb.SubResourceActionResponse{}
	if result != nil {
		pbStruct, err := structpb.NewStruct(result)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("marshal subresource result: %w", err))
		}
		resp.SetResult(pbStruct)
	}
	return resp, nil
}

// ---------------------------------------------------------------------------
// ShowChart
// ---------------------------------------------------------------------------

// ShowChart retrieves chart values and readme from a remote
// Helm repository.
func (s *RuntimeService) ShowChart(ctx context.Context, req *pb.ShowChartRequest) (*pb.ShowChartResponse, error) {
	values, readme, err := s.runtime.ShowChart(ctx, req.GetRepoUrl(), req.GetChartName(), req.GetVersion())
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	resp := &pb.ShowChartResponse{}
	resp.SetValues(values)
	resp.SetReadme(readme)
	return resp, nil
}
