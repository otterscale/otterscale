package core

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

// sessionTimeout bounds the session lifecycle assertions: the bugs
// they guard against are deadlocks, so a stuck session must fail the
// test instead of hanging the suite.
const sessionTimeout = 2 * time.Second

// userContext returns a context carrying an authenticated subject, as
// the auth middleware provides for every session RPC.
func userContext(t *testing.T, subject string) context.Context {
	t.Helper()
	return WithUserInfo(t.Context(), UserInfo{Subject: subject, Groups: []string{"system:authenticated"}})
}

// mockRuntimeRepo implements RuntimeRepo for these tests.
type mockRuntimeRepo struct {
	subResourceResult map[string]any
	subResourceErr    error

	// drainStdin makes PortForward and VNC behave like the real
	// adapters: they only return once the session's stdin pipe is
	// closed, which is what makes cancellation handling observable.
	drainStdin bool

	// execErr is what Exec fails with, standing in for a session that
	// never starts (no such container, RBAC denial, upgrade failure).
	execErr error

	// blockExec makes Exec run until its context is canceled, so that a
	// session that has not finished can be observed.
	blockExec bool

	// portForwardErr is what PortForward fails with, standing in for a
	// forward the kubelet refused.
	portForwardErr error
}

func (m *mockRuntimeRepo) PodLogs(context.Context, string, string, string, PodLogOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (m *mockRuntimeRepo) Exec(ctx context.Context, _, _, _ string, _ *ExecOptions) error {
	if m.blockExec {
		<-ctx.Done()
		return ctx.Err()
	}
	return m.execErr
}

func (m *mockRuntimeRepo) UpdateScale(context.Context, string, schema.GroupVersionResource, string, string, int32) (int32, error) {
	return 0, nil
}

func (m *mockRuntimeRepo) Restart(context.Context, string, schema.GroupVersionResource, string, string) error {
	return nil
}

func (m *mockRuntimeRepo) PortForward(_ context.Context, _, _, _ string, opts PortForwardOptions) error {
	if m.portForwardErr != nil {
		return m.portForwardErr
	}
	if m.drainStdin {
		_, err := io.Copy(io.Discard, opts.Stdin)
		return err
	}
	return nil
}

func (m *mockRuntimeRepo) SubResourceAction(_ context.Context, _ string, _ schema.GroupVersionResource,
	_, _, _, _ string, _ []byte,
) (map[string]any, error) {
	return m.subResourceResult, m.subResourceErr
}

func (m *mockRuntimeRepo) VNC(_ context.Context, _, _, _ string, opts VNCOptions) error {
	if m.drainStdin {
		_, err := io.Copy(io.Discard, opts.Stdin)
		return err
	}
	return nil
}

// mockDiscovery implements DiscoveryClient for these tests.
type mockDiscovery struct {
	lookupErr    error
	watchList    bool
	watchListErr error
}

func (m *mockDiscovery) LookupResource(_ context.Context, _, group, ver, resource, _ string) (schema.GroupVersionResource, error) {
	if m.lookupErr != nil {
		return schema.GroupVersionResource{}, m.lookupErr
	}
	return schema.GroupVersionResource{Group: group, Version: ver, Resource: resource}, nil
}

func (m *mockDiscovery) ServerResources(context.Context, string) ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (m *mockDiscovery) ResolveGroupVersionSchemas(context.Context, string, string, string) (map[string]*spec.Schema, error) {
	return nil, nil
}

func (m *mockDiscovery) ServerVersion(context.Context, string) (*version.Info, error) {
	return nil, nil
}

func (m *mockDiscovery) SupportsWatchList(context.Context, string) (bool, error) {
	return m.watchList, m.watchListErr
}

// mockHelmRepoForRuntime implements HelmRepo for these tests.
type mockHelmRepoForRuntime struct{}

func (m *mockHelmRepoForRuntime) ShowChart(context.Context, string, string, string) (values, readme []byte, err error) {
	return nil, nil, nil
}

func newTestRuntimeUseCase(discovery DiscoveryClient, runtime RuntimeRepo) *RuntimeUseCase {
	return NewRuntimeUseCase(discovery, runtime, &mockHelmRepoForRuntime{}, NewSessionStore())
}

func TestRuntimeUseCase_SubResourceAction_Validation(t *testing.T) {
	disco := &mockDiscovery{}
	repo := &mockRuntimeRepo{}
	uc := newTestRuntimeUseCase(disco, repo)

	tests := []struct {
		name        string
		id          *ResourceIdentifier
		subresource string
		method      string
		wantField   string
	}{
		{
			name:        "empty name",
			id:          &ResourceIdentifier{Cluster: "c", Group: "kubevirt.io", Version: "v1", Resource: "virtualmachines", Namespace: "ns"},
			subresource: "start",
			method:      "PUT",
			wantField:   "name",
		},
		{
			name:        "empty subresource",
			id:          &ResourceIdentifier{Cluster: "c", Group: "kubevirt.io", Version: "v1", Resource: "virtualmachines", Namespace: "ns", Name: "vm1"},
			subresource: "",
			method:      "PUT",
			wantField:   "subresource",
		},
		{
			name:        "invalid method DELETE",
			id:          &ResourceIdentifier{Cluster: "c", Group: "kubevirt.io", Version: "v1", Resource: "virtualmachines", Namespace: "ns", Name: "vm1"},
			subresource: "start",
			method:      "DELETE",
			wantField:   "method",
		},
		{
			name:        "invalid method PATCH",
			id:          &ResourceIdentifier{Cluster: "c", Group: "kubevirt.io", Version: "v1", Resource: "virtualmachines", Namespace: "ns", Name: "vm1"},
			subresource: "start",
			method:      "PATCH",
			wantField:   "method",
		},
		{
			name:        "invalid method GET",
			id:          &ResourceIdentifier{Cluster: "c", Group: "kubevirt.io", Version: "v1", Resource: "virtualmachines", Namespace: "ns", Name: "vm1"},
			subresource: "start",
			method:      "GET",
			wantField:   "method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.id.SubResource = tt.subresource
			_, err := uc.SubResourceAction(t.Context(), tt.id, tt.method, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var invalidInput *ErrInvalidInput
			if !errors.As(err, &invalidInput) {
				t.Fatalf("expected ErrInvalidInput, got %T: %v", err, err)
			}
			if invalidInput.Field != tt.wantField {
				t.Errorf("field = %q, want %q", invalidInput.Field, tt.wantField)
			}
		})
	}
}

func TestRuntimeUseCase_SubResourceAction_Success(t *testing.T) {
	disco := &mockDiscovery{}
	repo := &mockRuntimeRepo{
		subResourceResult: map[string]any{"status": "started"},
	}
	uc := newTestRuntimeUseCase(disco, repo)

	result, err := uc.SubResourceAction(
		t.Context(),
		&ResourceIdentifier{
			Cluster:     "prod",
			Group:       "kubevirt.io",
			Version:     "v1",
			Resource:    "virtualmachines",
			SubResource: "start",
			Namespace:   "default",
			Name:        "my-vm",
		},
		"PUT",
		nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["status"] != "started" {
		t.Errorf("result[status] = %v, want %q", result["status"], "started")
	}
}

func TestRuntimeUseCase_SubResourceAction_POST(t *testing.T) {
	disco := &mockDiscovery{}
	repo := &mockRuntimeRepo{
		subResourceResult: map[string]any{"ok": true},
	}
	uc := newTestRuntimeUseCase(disco, repo)

	result, err := uc.SubResourceAction(
		t.Context(),
		&ResourceIdentifier{
			Cluster:     "prod",
			Group:       "kubevirt.io",
			Version:     "v1",
			Resource:    "virtualmachines",
			SubResource: "start",
			Namespace:   "default",
			Name:        "my-vm",
		},
		"POST",
		[]byte(`{"gracePeriod": 30}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["ok"] != true {
		t.Errorf("result[ok] = %v, want true", result["ok"])
	}
}

func TestRuntimeUseCase_SubResourceAction_LookupError(t *testing.T) {
	disco := &mockDiscovery{lookupErr: errors.New("resource not found")}
	repo := &mockRuntimeRepo{}
	uc := newTestRuntimeUseCase(disco, repo)

	_, err := uc.SubResourceAction(
		t.Context(),
		&ResourceIdentifier{
			Cluster:     "prod",
			Group:       "kubevirt.io",
			Version:     "v1",
			Resource:    "virtualmachines",
			SubResource: "start",
			Namespace:   "default",
			Name:        "my-vm",
		},
		"PUT",
		nil,
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestStartPortForward_CancellationReleasesSession is the regression
// test for the leaked port-forward session: the adapter is blocked
// reading the session's stdin pipe, so nothing but cancellation can
// release it, and the session must become reapable afterwards.
func TestStartPortForward_CancellationReleasesSession(t *testing.T) {
	store := NewSessionStore()
	uc := NewRuntimeUseCase(&mockDiscovery{}, &mockRuntimeRepo{drainStdin: true}, &mockHelmRepoForRuntime{}, store)

	ctx, cancel := context.WithCancel(userContext(t, "alice"))
	sess, out, err := uc.StartPortForward(ctx, "prod", "default", "my-pod", 8080)
	if err != nil {
		t.Fatalf("StartPortForward: %v", err)
	}
	defer out.Close()

	cancel()

	select {
	case <-sess.Done:
	case <-time.After(sessionTimeout):
		t.Fatalf("port-forward adapter still running %s after cancellation", sessionTimeout)
	}

	// Done must stay signaled after being read, or the reaper can no
	// longer tell that the session has finished.
	if n := store.ReapStaleSessions(); n != 1 {
		t.Fatalf("ReapStaleSessions() = %d, want 1", n)
	}
}

// TestStartVNC_CancellationReleasesSession is the same regression test
// for VNC sessions.
func TestStartVNC_CancellationReleasesSession(t *testing.T) {
	store := NewSessionStore()
	uc := NewRuntimeUseCase(&mockDiscovery{}, &mockRuntimeRepo{drainStdin: true}, &mockHelmRepoForRuntime{}, store)

	ctx, cancel := context.WithCancel(userContext(t, "alice"))
	sess, out, err := uc.StartVNC(ctx, "prod", "default", "my-vmi")
	if err != nil {
		t.Fatalf("StartVNC: %v", err)
	}
	defer out.Close()

	cancel()

	select {
	case <-sess.Done:
	case <-time.After(sessionTimeout):
		t.Fatalf("VNC adapter still running %s after cancellation", sessionTimeout)
	}

	if n := store.ReapStaleSessions(); n != 1 {
		t.Fatalf("ReapStaleSessions() = %d, want 1", n)
	}
}

// TestStartExec_DoneStaysSignaled guards the reaper against a Done
// channel whose single buffered value has already been consumed by
// another caller (for example WriteExec's fast path).
func TestStartExec_DoneStaysSignaled(t *testing.T) {
	store := NewSessionStore()
	uc := NewRuntimeUseCase(&mockDiscovery{}, &mockRuntimeRepo{}, &mockHelmRepoForRuntime{}, store)

	sess, stdout, stderr, err := uc.StartExec(userContext(t, "alice"), &StartExecParams{
		Cluster:   "prod",
		Namespace: "default",
		Name:      "my-pod",
		Command:   []string{"sh"},
	})
	if err != nil {
		t.Fatalf("StartExec: %v", err)
	}
	defer stdout.Close()
	defer stderr.Close()

	select {
	case <-sess.Done:
	case <-time.After(sessionTimeout):
		t.Fatalf("exec session did not finish within %s", sessionTimeout)
	}

	if n := store.ReapStaleSessions(); n != 1 {
		t.Fatalf("ReapStaleSessions() = %d, want 1", n)
	}
}

// TestStartSession_RequiresAuthenticatedUser checks that a session
// cannot be opened without an identity to bind it to.
func TestStartSession_RequiresAuthenticatedUser(t *testing.T) {
	uc := NewRuntimeUseCase(&mockDiscovery{}, &mockRuntimeRepo{}, &mockHelmRepoForRuntime{}, NewSessionStore())

	tests := map[string]func(context.Context) error{
		"exec": func(ctx context.Context) error {
			_, _, _, err := uc.StartExec(ctx, &StartExecParams{Name: "my-pod", Command: []string{"sh"}})
			return err
		},
		"port-forward": func(ctx context.Context) error {
			_, _, err := uc.StartPortForward(ctx, "prod", "default", "my-pod", 8080)
			return err
		},
		"vnc": func(ctx context.Context) error {
			_, _, err := uc.StartVNC(ctx, "prod", "default", "my-vmi")
			return err
		},
	}

	for name, start := range tests {
		t.Run(name, func(t *testing.T) {
			err := start(t.Context()) // no user info
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if code, ok := DomainErrorCode(err); !ok || code != ErrorCodeUnauthenticated {
				t.Errorf("code = %v (domain=%v), want ErrorCodeUnauthenticated", code, ok)
			}
		})
	}
}

// TestSessionsAreBoundToTheirOwner is the regression test for session
// identifiers being the only thing protecting an open shell: a second
// authenticated user must not be able to drive somebody else's session
// even when they know its identifier.
func TestSessionsAreBoundToTheirOwner(t *testing.T) {
	store := NewSessionStore()
	uc := NewRuntimeUseCase(&mockDiscovery{}, &mockRuntimeRepo{drainStdin: true}, &mockHelmRepoForRuntime{}, store)

	alice := userContext(t, "alice")
	mallory := userContext(t, "mallory")

	execSess, stdout, stderr, err := uc.StartExec(alice, &StartExecParams{Name: "my-pod", Command: []string{"sh"}})
	if err != nil {
		t.Fatalf("StartExec: %v", err)
	}
	defer stdout.Close()
	defer stderr.Close()

	pfSess, pfOut, err := uc.StartPortForward(alice, "prod", "default", "my-pod", 8080)
	if err != nil {
		t.Fatalf("StartPortForward: %v", err)
	}
	defer pfOut.Close()

	vncSess, vncOut, err := uc.StartVNC(alice, "prod", "default", "my-vmi")
	if err != nil {
		t.Fatalf("StartVNC: %v", err)
	}
	defer vncOut.Close()

	denied := map[string]func() error{
		"write exec":         func() error { return uc.WriteExec(mallory, execSess.ID, []byte("id\n")) },
		"resize exec":        func() error { return uc.ResizeExec(mallory, execSess.ID, 40, 120) },
		"write port-forward": func() error { return uc.WritePortForward(mallory, pfSess.ID, []byte("GET /")) },
		"write vnc":          func() error { return uc.WriteVNC(mallory, vncSess.ID, []byte{0x01}) },
	}

	for name, call := range denied {
		t.Run(name, func(t *testing.T) {
			err := call()
			var notFound *ErrSessionNotFound
			if !errors.As(err, &notFound) {
				t.Fatalf("got %v (%T), want ErrSessionNotFound", err, err)
			}
		})
	}

	t.Run("cleanup by another user leaves the session running", func(t *testing.T) {
		uc.CleanupExec(mallory, execSess.ID)
		uc.CleanupPortForward(mallory, pfSess.ID)
		uc.CleanupVNC(mallory, vncSess.ID)

		if _, ok := store.GetExec(execSess.ID); !ok {
			t.Error("exec session was removed by a caller that does not own it")
		}
		if _, ok := store.GetPortForward(pfSess.ID); !ok {
			t.Error("port-forward session was removed by a caller that does not own it")
		}
		if _, ok := store.GetVNC(vncSess.ID); !ok {
			t.Error("VNC session was removed by a caller that does not own it")
		}
	})

	t.Run("the owner can still resize and clean up", func(t *testing.T) {
		if err := uc.ResizeExec(alice, execSess.ID, 40, 120); err != nil {
			t.Errorf("ResizeExec by the owner: %v", err)
		}

		uc.CleanupExec(alice, execSess.ID)
		uc.CleanupPortForward(alice, pfSess.ID)
		uc.CleanupVNC(alice, vncSess.ID)

		if _, ok := store.GetExec(execSess.ID); ok {
			t.Error("exec session survived cleanup by its owner")
		}
		if _, ok := store.GetPortForward(pfSess.ID); ok {
			t.Error("port-forward session survived cleanup by its owner")
		}
		if _, ok := store.GetVNC(vncSess.ID); ok {
			t.Error("VNC session survived cleanup by its owner")
		}
	})
}

// TestWaitExecReportsSessionFailure is the regression test for exec
// failures that never reached the caller. The session goroutine records
// why the exec ended, but nothing ever read that field, so a session
// that failed to start was indistinguishable from one that produced no
// output.
func TestWaitExecReportsSessionFailure(t *testing.T) {
	wantErr := &DomainError{Code: ErrorCodeNotFound, Message: "container \"nope\" not found"}
	uc := newTestRuntimeUseCase(&mockDiscovery{}, &mockRuntimeRepo{execErr: wantErr})

	ctx := userContext(t, "alice")
	sess, stdout, stderr, err := uc.StartExec(ctx, &StartExecParams{Name: "pod-a", Command: []string{"sh"}})
	if err != nil {
		t.Fatalf("StartExec: %v", err)
	}
	defer stdout.Close()
	defer stderr.Close()

	waitCtx, cancel := context.WithTimeout(ctx, sessionTimeout)
	defer cancel()

	got := uc.WaitExec(waitCtx, sess)
	if !errors.Is(got, wantErr) {
		t.Fatalf("WaitExec = %v, want %v", got, wantErr)
	}
}

// TestWaitExecReportsCommandExit checks that a command which ran and
// exited non-zero is reported as such, so callers can tell it apart
// from a session that never started.
func TestWaitExecReportsCommandExit(t *testing.T) {
	uc := newTestRuntimeUseCase(&mockDiscovery{}, &mockRuntimeRepo{execErr: &ErrCommandExited{Code: 2}})

	ctx := userContext(t, "alice")
	sess, stdout, stderr, err := uc.StartExec(ctx, &StartExecParams{Name: "pod-a", Command: []string{"false"}})
	if err != nil {
		t.Fatalf("StartExec: %v", err)
	}
	defer stdout.Close()
	defer stderr.Close()

	waitCtx, cancel := context.WithTimeout(ctx, sessionTimeout)
	defer cancel()

	var exited *ErrCommandExited
	if got := uc.WaitExec(waitCtx, sess); !errors.As(got, &exited) {
		t.Fatalf("WaitExec = %v, want *ErrCommandExited", got)
	}
	if exited.Code != 2 {
		t.Fatalf("exit code = %d, want 2", exited.Code)
	}
}

// TestWaitExecGivesUpWithCaller checks that waiting for a session that
// has not finished is bounded by the caller's context, so that a wedged
// session cannot hold its RPC open indefinitely.
func TestWaitExecGivesUpWithCaller(t *testing.T) {
	uc := newTestRuntimeUseCase(&mockDiscovery{}, &mockRuntimeRepo{blockExec: true})

	ctx := userContext(t, "alice")
	sess, stdout, stderr, err := uc.StartExec(ctx, &StartExecParams{Name: "pod-a", Command: []string{"sleep"}})
	if err != nil {
		t.Fatalf("StartExec: %v", err)
	}
	defer stdout.Close()
	defer stderr.Close()
	defer uc.CleanupExec(ctx, sess.ID)

	waitCtx, cancel := context.WithCancel(ctx)
	cancel()

	if got := uc.WaitExec(waitCtx, sess); !errors.Is(got, context.Canceled) {
		t.Fatalf("WaitExec = %v, want context.Canceled", got)
	}
}

// TestWaitPortForwardReportsSessionFailure covers the port-forward half
// of the same bug.
func TestWaitPortForwardReportsSessionFailure(t *testing.T) {
	wantErr := &DomainError{Code: ErrorCodeFailedPrecondition, Message: "port-forward refused by kubelet: socat not found"}
	uc := newTestRuntimeUseCase(&mockDiscovery{}, &mockRuntimeRepo{portForwardErr: wantErr})

	ctx := userContext(t, "alice")
	sess, out, err := uc.StartPortForward(ctx, "c1", "default", "pod-a", 8080)
	if err != nil {
		t.Fatalf("StartPortForward: %v", err)
	}
	defer out.Close()

	waitCtx, cancel := context.WithTimeout(ctx, sessionTimeout)
	defer cancel()

	if got := uc.WaitPortForward(waitCtx, sess); !errors.Is(got, wantErr) {
		t.Fatalf("WaitPortForward = %v, want %v", got, wantErr)
	}
}
