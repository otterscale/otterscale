package core

import (
	"context"
	"errors"
	"io"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

// mockRuntimeRepo implements RuntimeRepo for SubResourceAction testing.
type mockRuntimeRepo struct {
	subResourceResult map[string]any
	subResourceErr    error
}

func (m *mockRuntimeRepo) PodLogs(context.Context, string, string, string, PodLogOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (m *mockRuntimeRepo) Exec(context.Context, string, string, string, *ExecOptions) error {
	return nil
}

func (m *mockRuntimeRepo) GetScale(context.Context, string, schema.GroupVersionResource, string, string) (int32, error) {
	return 0, nil
}

func (m *mockRuntimeRepo) UpdateScale(context.Context, string, schema.GroupVersionResource, string, string, int32) (int32, error) {
	return 0, nil
}

func (m *mockRuntimeRepo) Restart(context.Context, string, schema.GroupVersionResource, string, string) error {
	return nil
}

func (m *mockRuntimeRepo) PortForward(context.Context, string, string, string, PortForwardOptions) error {
	return nil
}

func (m *mockRuntimeRepo) SubResourceAction(_ context.Context, _ string, _ schema.GroupVersionResource,
	_, _, _, _ string, _ []byte,
) (map[string]any, error) {
	return m.subResourceResult, m.subResourceErr
}

// mockDiscoveryForRuntime implements DiscoveryClient for runtime tests.
type mockDiscoveryForRuntime struct {
	lookupErr error
}

func (m *mockDiscoveryForRuntime) LookupResource(_ context.Context, _, group, ver, resource, _ string) (schema.GroupVersionResource, error) {
	if m.lookupErr != nil {
		return schema.GroupVersionResource{}, m.lookupErr
	}
	return schema.GroupVersionResource{Group: group, Version: ver, Resource: resource}, nil
}

func (m *mockDiscoveryForRuntime) ServerResources(context.Context, string) ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (m *mockDiscoveryForRuntime) ResolveSchema(context.Context, string, string, string, string) (*spec.Schema, error) {
	return nil, nil
}

func (m *mockDiscoveryForRuntime) ServerVersion(context.Context, string) (*version.Info, error) {
	return nil, nil
}

func (m *mockDiscoveryForRuntime) SupportsWatchList(context.Context, string) (bool, error) {
	return false, nil
}

func newTestRuntimeUseCase(discovery DiscoveryClient, runtime RuntimeRepo) *RuntimeUseCase {
	return NewRuntimeUseCase(discovery, runtime, NewSessionStore())
}

func TestRuntimeUseCase_SubResourceAction_Validation(t *testing.T) {
	disco := &mockDiscoveryForRuntime{}
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
	disco := &mockDiscoveryForRuntime{}
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
	disco := &mockDiscoveryForRuntime{}
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
	disco := &mockDiscoveryForRuntime{lookupErr: errors.New("resource not found")}
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
