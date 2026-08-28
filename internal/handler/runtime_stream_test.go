package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/kube-openapi/pkg/validation/spec"

	pb "github.com/otterscale/otterscale/api/runtime/v1"

	"github.com/otterscale/otterscale/internal/core"
)

// These tests exercise the streaming handlers through a real Connect
// server, because the bug they guard against lives in what the handler
// returns once its output channel closes — which is exactly what a
// client observes and nothing below the handler can see.

// stubRuntimeRepo fails the operation under test the way the real
// adapter does when a session cannot be established.
type stubRuntimeRepo struct {
	execErr        error
	portForwardErr error
	vncErr         error
}

func (s *stubRuntimeRepo) PodLogs(context.Context, string, string, string, core.PodLogOptions) (io.ReadCloser, error) {
	return io.NopCloser(nil), nil
}

func (s *stubRuntimeRepo) Exec(context.Context, string, string, string, *core.ExecOptions) error {
	return s.execErr
}

func (s *stubRuntimeRepo) UpdateScale(context.Context, string, schema.GroupVersionResource, string, string, int32) (int32, error) {
	return 0, nil
}

func (s *stubRuntimeRepo) Restart(context.Context, string, schema.GroupVersionResource, string, string) error {
	return nil
}

func (s *stubRuntimeRepo) PortForward(context.Context, string, string, string, core.PortForwardOptions) error {
	return s.portForwardErr
}

func (s *stubRuntimeRepo) SubResourceAction(context.Context, string, schema.GroupVersionResource,
	string, string, string, string, []byte,
) (map[string]any, error) {
	return nil, nil
}

func (s *stubRuntimeRepo) VNC(context.Context, string, string, string, core.VNCOptions) error {
	return s.vncErr
}

type stubDiscovery struct{}

func (stubDiscovery) LookupResource(_ context.Context, _, group, ver, resource, _ string) (schema.GroupVersionResource, error) {
	return schema.GroupVersionResource{Group: group, Version: ver, Resource: resource}, nil
}

func (stubDiscovery) ServerResources(context.Context, string) ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (stubDiscovery) ResolveGroupVersionSchemas(context.Context, string, string, string) (map[string]*spec.Schema, error) {
	return nil, nil
}

func (stubDiscovery) ServerVersion(context.Context, string) (*version.Info, error) { return nil, nil }

func (stubDiscovery) SupportsWatchList(context.Context, string) (bool, error) { return false, nil }

type stubHelmRepo struct{}

func (stubHelmRepo) ShowChart(context.Context, string, string, string) (values, readme []byte, err error) {
	return nil, nil, nil
}

// newRuntimeTestClient mounts the RuntimeService on a real Connect
// server and returns a client for it. The auth middleware is stubbed by
// a handler that injects UserInfo, as the real one does for every
// session RPC.
func newRuntimeTestClient(t *testing.T, repo core.RuntimeRepo) pb.RuntimeServiceClient {
	t.Helper()

	uc := core.NewRuntimeUseCase(stubDiscovery{}, repo, stubHelmRepo{}, core.NewSessionStore())

	mux := http.NewServeMux()
	mux.Handle(pb.NewRuntimeServiceHandler(NewRuntimeService(uc)))

	withUser := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := core.WithUserInfo(r.Context(), core.UserInfo{
			Subject: "alice",
			Groups:  []string{"system:authenticated"},
		})
		mux.ServeHTTP(w, r.WithContext(ctx))
	})

	srv := httptest.NewServer(withUser)
	t.Cleanup(srv.Close)

	return pb.NewRuntimeServiceClient(srv.Client(), srv.URL)
}

// drain consumes a server stream to completion and returns the error it
// ended with, along with how many messages arrived.
func drain[T any](stream *connect.ServerStreamForClient[T]) (int, error) {
	n := 0
	for stream.Receive() {
		n++
	}
	if err := stream.Err(); err != nil {
		return n, err
	}
	return n, stream.Close()
}

// TestExecuteTTYReportsSessionFailure is the end-to-end regression test:
// an exec that never starts used to end as a clean stream carrying only
// the session id, so the client saw success.
func TestExecuteTTYReportsSessionFailure(t *testing.T) {
	client := newRuntimeTestClient(t, &stubRuntimeRepo{
		execErr: &core.DomainError{Code: core.ErrorCodeNotFound, Message: `container "nope" not found`},
	})

	req := &pb.ExecuteTTYRequest{}
	req.SetCluster("c1")
	req.SetNamespace("default")
	req.SetName("pod-a")
	req.SetCommand([]string{"sh"})

	stream, err := client.ExecuteTTY(t.Context(), req)
	if err != nil {
		t.Fatalf("ExecuteTTY: %v", err)
	}

	n, err := drain(stream)
	if err == nil {
		t.Fatalf("stream ended cleanly after %d message(s), want a NotFound error", n)
	}
	if code := connect.CodeOf(err); code != connect.CodeNotFound {
		t.Fatalf("code = %v (%v), want %v", code, err, connect.CodeNotFound)
	}
}

// TestExecuteTTYDoesNotFailOnNonZeroExit guards the other half of the
// distinction: a command that ran and exited non-zero is a result, not
// an RPC failure.
func TestExecuteTTYDoesNotFailOnNonZeroExit(t *testing.T) {
	client := newRuntimeTestClient(t, &stubRuntimeRepo{
		execErr: &core.ErrCommandExited{Code: 1, Reason: "command terminated with exit code 1"},
	})

	req := &pb.ExecuteTTYRequest{}
	req.SetCluster("c1")
	req.SetNamespace("default")
	req.SetName("pod-a")
	req.SetCommand([]string{"false"})

	stream, err := client.ExecuteTTY(t.Context(), req)
	if err != nil {
		t.Fatalf("ExecuteTTY: %v", err)
	}

	if _, err := drain(stream); err != nil {
		t.Fatalf("stream failed with %v, want a clean end", err)
	}
}

// TestPortForwardReportsSessionFailure covers the port-forward half.
func TestPortForwardReportsSessionFailure(t *testing.T) {
	client := newRuntimeTestClient(t, &stubRuntimeRepo{
		portForwardErr: &core.DomainError{
			Code:    core.ErrorCodeFailedPrecondition,
			Message: "port-forward refused by kubelet: socat not found",
		},
	})

	req := &pb.PortForwardRequest{}
	req.SetCluster("c1")
	req.SetNamespace("default")
	req.SetName("pod-a")
	req.SetPort(8080)

	stream, err := client.PortForward(t.Context(), req)
	if err != nil {
		t.Fatalf("PortForward: %v", err)
	}

	n, err := drain(stream)
	if err == nil {
		t.Fatalf("stream ended cleanly after %d message(s), want a FailedPrecondition error", n)
	}
	if code := connect.CodeOf(err); code != connect.CodeFailedPrecondition {
		t.Fatalf("code = %v (%v), want %v", code, err, connect.CodeFailedPrecondition)
	}
}

// TestVNCReportsSessionFailure covers the VNC half.
func TestVNCReportsSessionFailure(t *testing.T) {
	client := newRuntimeTestClient(t, &stubRuntimeRepo{
		vncErr: &core.DomainError{Code: core.ErrorCodeNotFound, Message: "VMI not found"},
	})

	req := &pb.VNCRequest{}
	req.SetCluster("c1")
	req.SetNamespace("default")
	req.SetName("vmi-a")

	stream, err := client.VNC(t.Context(), req)
	if err != nil {
		t.Fatalf("VNC: %v", err)
	}

	n, err := drain(stream)
	if err == nil {
		t.Fatalf("stream ended cleanly after %d message(s), want a NotFound error", n)
	}
	if code := connect.CodeOf(err); code != connect.CodeNotFound {
		t.Fatalf("code = %v (%v), want %v", code, err, connect.CodeNotFound)
	}
}
