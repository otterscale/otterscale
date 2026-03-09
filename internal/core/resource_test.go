package core

import (
	"context"
	"errors"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

// mockDiscoveryForColumns implements DiscoveryClient for columns tests.
type mockDiscoveryForColumns struct {
	lookupErr  error
	columns    []ColumnDefinition
	columnsErr error
}

func (m *mockDiscoveryForColumns) LookupResource(_ context.Context, _, group, ver, resource string) (schema.GroupVersionResource, error) {
	if m.lookupErr != nil {
		return schema.GroupVersionResource{}, m.lookupErr
	}
	return schema.GroupVersionResource{Group: group, Version: ver, Resource: resource}, nil
}

func (m *mockDiscoveryForColumns) ServerResources(context.Context, string) ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (m *mockDiscoveryForColumns) ResolveSchema(context.Context, string, string, string, string) (*spec.Schema, error) {
	return nil, nil
}

func (m *mockDiscoveryForColumns) ServerVersion(context.Context, string) (*version.Info, error) {
	return nil, nil
}

func (m *mockDiscoveryForColumns) SupportsWatchList(context.Context, string) (bool, error) {
	return false, nil
}

func (m *mockDiscoveryForColumns) Columns(_ context.Context, _ string, _ schema.GroupVersionResource, _ string) ([]ColumnDefinition, error) {
	return m.columns, m.columnsErr
}

// mockColumnsProvider implements ColumnsProvider for resource use-case tests.
type mockColumnsProvider struct {
	columns []ColumnDefinition
	err     error
}

func (m *mockColumnsProvider) Columns(context.Context, string, schema.GroupVersionResource, string) ([]ColumnDefinition, error) {
	return m.columns, m.err
}

// mockResourceRepo implements ResourceRepo (minimal, only for Columns testing).
type mockResourceRepo struct{}

func (m *mockResourceRepo) List(context.Context, string, schema.GroupVersionResource, string, ListOptions) (*unstructured.UnstructuredList, error) {
	return nil, nil
}

func (m *mockResourceRepo) Get(context.Context, string, schema.GroupVersionResource, string, string) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (m *mockResourceRepo) Create(context.Context, string, schema.GroupVersionResource, string, []byte) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (m *mockResourceRepo) Apply(context.Context, string, schema.GroupVersionResource, string, string, []byte, ApplyOptions) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (m *mockResourceRepo) Delete(context.Context, string, schema.GroupVersionResource, string, string, DeleteOptions) error {
	return nil
}

func (m *mockResourceRepo) Watch(context.Context, string, schema.GroupVersionResource, string, WatchOptions) (Watcher, error) {
	return nil, nil
}

func (m *mockResourceRepo) ListEvents(context.Context, string, string, ListOptions) (*unstructured.UnstructuredList, error) {
	return nil, nil
}

func TestResourceUseCase_Columns_Success(t *testing.T) {
	cols := []ColumnDefinition{
		{Name: "Name", Type: "string", Priority: 0},
		{Name: "Ready", Type: "string", Priority: 0},
		{Name: "Node", Type: "string", Priority: 1},
	}

	disco := &mockDiscoveryForColumns{}
	cp := &mockColumnsProvider{columns: cols}
	uc := NewResourceUseCase(disco, &mockResourceRepo{}, nil, cp)

	result, err := uc.Columns(t.Context(), &ResourceIdentifier{
		Cluster:  "prod",
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(result))
	}
	if result[0].Name != "Name" {
		t.Errorf("result[0].Name = %q, want %q", result[0].Name, "Name")
	}
	if result[2].Priority != 1 {
		t.Errorf("result[2].Priority = %d, want 1", result[2].Priority)
	}
}

func TestResourceUseCase_Columns_LookupError(t *testing.T) {
	disco := &mockDiscoveryForColumns{lookupErr: errors.New("not found")}
	cp := &mockColumnsProvider{}
	uc := NewResourceUseCase(disco, &mockResourceRepo{}, nil, cp)

	_, err := uc.Columns(t.Context(), &ResourceIdentifier{
		Cluster:  "prod",
		Group:    "apps",
		Version:  "v1",
		Resource: "nonexistent",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResourceUseCase_Columns_ProviderError(t *testing.T) {
	disco := &mockDiscoveryForColumns{}
	cp := &mockColumnsProvider{err: errors.New("table API error")}
	uc := NewResourceUseCase(disco, &mockResourceRepo{}, nil, cp)

	_, err := uc.Columns(t.Context(), &ResourceIdentifier{
		Cluster:  "prod",
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResourceUseCase_Columns_Empty(t *testing.T) {
	disco := &mockDiscoveryForColumns{}
	cp := &mockColumnsProvider{columns: []ColumnDefinition{}}
	uc := NewResourceUseCase(disco, &mockResourceRepo{}, nil, cp)

	result, err := uc.Columns(t.Context(), &ResourceIdentifier{
		Cluster:  "prod",
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 columns, got %d", len(result))
	}
}
