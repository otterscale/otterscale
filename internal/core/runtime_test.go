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

// mockRuntimeRepo implements RuntimeRepo for SubResourceAction testing.
type mockRuntimeRepo struct {
	subResourceResult map[string]any
	subResourceErr    error

	// drainStdin makes PortForward and VNC behave like the real
	// adapters: they only return once the session's stdin pipe is
	// closed, which is what makes cancellation handling observable.
	drainStdin bool
}

func (m *mockRuntimeRepo) PodLogs(context.Context, string, string, string, PodLogOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (m *mockRuntimeRepo) Exec(context.Context, string, string, string, *ExecOptions) error {
	return nil
}

func (m *mockRuntimeRepo) UpdateScale(context.Context, string, schema.GroupVersionResource, string, string, int32) (int32, error) {
	return 0, nil
}

func (m *mockRuntimeRepo) Restart(context.Context, string, schema.GroupVersionResource, string, string) error {
	return nil
}

func (m *mockRuntimeRepo) PortForward(_ context.Context, _, _, _ string, opts PortForwardOptions) error {
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

// mockDiscovery implements DiscoveryClient for use-case tests.
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

// mockHelmRepoForRuntime implements HelmRepo for runtime tests.
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

	ctx, cancel := context.WithCancel(t.Context())
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

	ctx, cancel := context.WithCancel(t.Context())
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

	sess, stdout, stderr, err := uc.StartExec(t.Context(), &StartExecParams{
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
