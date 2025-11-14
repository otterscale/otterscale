package workload

import (
	"context"
	"fmt"
	"io"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	PodPhaseRunning = v1.PodRunning

	PodPhaseSucceeded = v1.PodSucceeded
)

type ttySession struct {
	ID        string
	InReader  *io.PipeReader
	InWriter  *io.PipeWriter
	OutReader *io.PipeReader
	OutWriter *io.PipeWriter
}

type (
	// Pod represents a Kubernetes Pod resource.
	Pod = v1.Pod

	// PodCondition represents a Kubernetes PodCondition resource.
	PodCondition = v1.PodCondition

	// Container represents a Kubernetes Container resource.
	Container = v1.Container

	// ContainerStatus represents a Kubernetes ContainerStatus resource.
	ContainerStatus = v1.ContainerStatus
)

type PodRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]Pod, error)
	Create(ctx context.Context, scope, namespace string, p *Pod) (*Pod, error)
	Update(ctx context.Context, scope, namespace string, p *Pod) (*Pod, error)
	Delete(ctx context.Context, scope, namespace, name string) error
	Stream(ctx context.Context, scope, namespace, podName, containerName string, duration time.Duration, follow bool) (io.ReadCloser, error)
	Execute(scope, namespace, podName, containerName string, command []string) (remotecommand.Executor, error)
}

func (uc *WorkloadUseCase) DeletePod(ctx context.Context, scope, namespace, name string) error {
	return uc.pod.Delete(ctx, scope, namespace, name)
}

func (uc *WorkloadUseCase) StreamLogs(ctx context.Context, scope, namespace, podName, containerName string, duration time.Duration) (io.ReadCloser, error) {
	return uc.pod.Stream(ctx, scope, namespace, podName, containerName, duration, true)
}

func (uc *WorkloadUseCase) WriteToTTYSession(sessionID string, stdIn []byte) error {
	value, ok := uc.ttySessions.Load(sessionID)
	if !ok {
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("session %s not found", sessionID))
	}

	if _, err := value.(*ttySession).InWriter.Write(stdIn); err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to write to session: %w", err))
	}

	return nil
}

func (uc *WorkloadUseCase) CreateTTYSession() (string, error) {
	sessionID := uuid.New().String()

	inReader, inWriter := io.Pipe()
	outReader, outWriter := io.Pipe()

	uc.ttySessions.Store(sessionID, &ttySession{
		ID:        sessionID,
		InReader:  inReader,
		InWriter:  inWriter,
		OutReader: outReader,
		OutWriter: outWriter,
	})

	return sessionID, nil
}

func (uc *WorkloadUseCase) CleanupTTYSession(sessionID string) error {
	value, ok := uc.ttySessions.Load(sessionID)
	if !ok {
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("session %s not found", sessionID))
	}

	ttySession := value.(*ttySession)
	ttySession.InReader.Close()
	ttySession.InWriter.Close()
	ttySession.OutReader.Close()
	ttySession.OutWriter.Close()

	uc.ttySessions.Delete(sessionID)

	return nil
}

func (uc *WorkloadUseCase) ExecuteTTY(ctx context.Context, sessionID, scope, namespace, podName, containerName string, command []string, stdOut chan<- []byte) error {
	value, ok := uc.ttySessions.Load(sessionID)
	if !ok {
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("session %s not found", sessionID))
	}

	ttySession := value.(*ttySession)

	exec, err := uc.pod.Execute(scope, namespace, podName, containerName, command)
	if err != nil {
		return err
	}

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		buf := make([]byte, 1024) //nolint:mnd // 1KB buffer

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			default:
				n, err := ttySession.OutReader.Read(buf)
				if err != nil {
					if err == io.EOF {
						return nil
					}
					return err
				}

				if n > 0 {
					// write message to std out
					stdOut <- buf[:n]
				}
			}
		}
	})

	eg.Go(func() error {
		defer close(stdOut)

		return exec.StreamWithContext(egctx, remotecommand.StreamOptions{
			Stdin:  ttySession.InReader,
			Stdout: ttySession.OutWriter, // raw TTY manages stdout and stderr over the stdout stream
			Tty:    true,
		})
	})

	return eg.Wait()
}
