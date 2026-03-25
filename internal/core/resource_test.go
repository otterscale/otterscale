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

// mockDiscoveryForResource implements DiscoveryClient for resource tests.
type mockDiscoveryForResource struct {
	lookupGVR         schema.GroupVersionResource
	lookupNamespaced  bool
	lookupErr         error
	supportsWatchList bool
}

func (m *mockDiscoveryForResource) LookupResource(_ context.Context, _, group, ver, resource, _ string) (schema.GroupVersionResource, bool, error) {
	if m.lookupErr != nil {
		return schema.GroupVersionResource{}, false, m.lookupErr
	}
	if m.lookupGVR.Empty() {
		return schema.GroupVersionResource{Group: group, Version: ver, Resource: resource}, m.lookupNamespaced, nil
	}
	return m.lookupGVR, m.lookupNamespaced, nil
}

func (m *mockDiscoveryForResource) ServerResources(context.Context, string) ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (m *mockDiscoveryForResource) ResolveSchema(context.Context, string, string, string, string) (*spec.Schema, error) {
	return nil, nil
}

func (m *mockDiscoveryForResource) ServerVersion(context.Context, string) (*version.Info, error) {
	return nil, nil
}

func (m *mockDiscoveryForResource) SupportsWatchList(context.Context, string) (bool, error) {
	return m.supportsWatchList, nil
}

// TestResourceIdentifier_validateNamespaceScope tests the namespace scope validation logic.
func TestResourceIdentifier_validateNamespaceScope(t *testing.T) {
	tests := []struct {
		name         string
		id           *ResourceIdentifier
		gvr          schema.GroupVersionResource
		isNamespaced bool
		wantErr      bool
		wantErrField string
		wantErrMsg   string
	}{
		{
			name: "cluster-scoped resource with namespace should error",
			id: &ResourceIdentifier{
				Namespace: "default",
			},
			gvr: schema.GroupVersionResource{
				Group:    "apiextensions.k8s.io",
				Version:  "v1",
				Resource: "customresourcedefinitions",
			},
			isNamespaced: false,
			wantErr:      true,
			wantErrField: "namespace",
			wantErrMsg:   "cluster-scoped and cannot have a namespace",
		},
		{
			name: "cluster-scoped resource without namespace should pass",
			id: &ResourceIdentifier{
				Namespace: "",
			},
			gvr: schema.GroupVersionResource{
				Group:    "apiextensions.k8s.io",
				Version:  "v1",
				Resource: "customresourcedefinitions",
			},
			isNamespaced: false,
			wantErr:      false,
		},
		{
			name: "namespace-scoped resource without namespace should error",
			id: &ResourceIdentifier{
				Namespace: "",
			},
			gvr: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "pods",
			},
			isNamespaced: true,
			wantErr:      true,
			wantErrField: "namespace",
			wantErrMsg:   "namespace-scoped and requires a namespace",
		},
		{
			name: "namespace-scoped resource with namespace should pass",
			id: &ResourceIdentifier{
				Namespace: "default",
			},
			gvr: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "pods",
			},
			isNamespaced: true,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.id.validateNamespaceScope(tt.gvr, tt.isNamespaced)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateNamespaceScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				var invalidInput *ErrInvalidInput
				if !errors.As(err, &invalidInput) {
					t.Fatalf("expected ErrInvalidInput, got %T: %v", err, err)
				}
				if invalidInput.Field != tt.wantErrField {
					t.Errorf("field = %q, want %q", invalidInput.Field, tt.wantErrField)
				}
				if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("error message = %q, want to contain %q", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}

// TestResourceUseCase_CreateResource_NamespaceValidation tests namespace scope validation in CreateResource.
func TestResourceUseCase_CreateResource_NamespaceValidation(t *testing.T) {
	tests := []struct {
		name         string
		id           *ResourceIdentifier
		isNamespaced bool
		wantErr      bool
		wantErrMsg   string
	}{
		{
			name: "create cluster-scoped resource with namespace should fail",
			id: &ResourceIdentifier{
				Cluster:   "test-cluster",
				Group:     "apiextensions.k8s.io",
				Version:   "v1",
				Resource:  "customresourcedefinitions",
				Namespace: "default",
			},
			isNamespaced: false,
			wantErr:      true,
			wantErrMsg:   "cluster-scoped and cannot have a namespace",
		},
		{
			name: "create namespace-scoped resource without namespace should fail",
			id: &ResourceIdentifier{
				Cluster:   "test-cluster",
				Group:     "",
				Version:   "v1",
				Resource:  "pods",
				Namespace: "",
			},
			isNamespaced: true,
			wantErr:      true,
			wantErrMsg:   "namespace-scoped and requires a namespace",
		},
		{
			name: "create cluster-scoped resource without namespace should succeed",
			id: &ResourceIdentifier{
				Cluster:   "test-cluster",
				Group:     "apiextensions.k8s.io",
				Version:   "v1",
				Resource:  "customresourcedefinitions",
				Namespace: "",
			},
			isNamespaced: false,
			wantErr:      false,
		},
		{
			name: "create namespace-scoped resource with namespace should succeed",
			id: &ResourceIdentifier{
				Cluster:   "test-cluster",
				Group:     "",
				Version:   "v1",
				Resource:  "pods",
				Namespace: "default",
			},
			isNamespaced: true,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disco := &mockDiscoveryForResource{
				lookupNamespaced: tt.isNamespaced,
			}
			repo := &mockResourceRepo{}
			uc := NewResourceUseCase(disco, repo, nil)

			_, err := uc.CreateResource(context.Background(), tt.id, []byte("manifest"))
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("error message = %q, want to contain %q", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

// TestResourceUseCase_GetResource_NamespaceValidation tests namespace scope validation in GetResource.
func TestResourceUseCase_GetResource_NamespaceValidation(t *testing.T) {
	disco := &mockDiscoveryForResource{
		lookupNamespaced: false, // cluster-scoped
	}
	repo := &mockResourceRepo{}
	uc := NewResourceUseCase(disco, repo, nil)

	id := &ResourceIdentifier{
		Cluster:   "test-cluster",
		Group:     "apiextensions.k8s.io",
		Version:   "v1",
		Resource:  "customresourcedefinitions",
		Namespace: "default", // should fail
		Name:      "test-crd",
	}

	_, err := uc.GetResource(context.Background(), id)
	if err == nil {
		t.Fatal("expected error for cluster-scoped resource with namespace, got nil")
	}

	var invalidInput *ErrInvalidInput
	if !errors.As(err, &invalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %T: %v", err, err)
	}
	if !contains(err.Error(), "cluster-scoped") {
		t.Errorf("error message = %q, want to contain 'cluster-scoped'", err.Error())
	}
}

// TestResourceUseCase_ListResources_NamespaceValidation tests namespace scope validation in ListResources.
func TestResourceUseCase_ListResources_NamespaceValidation(t *testing.T) {
	disco := &mockDiscoveryForResource{
		lookupNamespaced: true, // namespace-scoped
	}
	repo := &mockResourceRepo{}
	uc := NewResourceUseCase(disco, repo, nil)

	id := &ResourceIdentifier{
		Cluster:   "test-cluster",
		Group:     "",
		Version:   "v1",
		Resource:  "pods",
		Namespace: "", // should fail
	}

	_, err := uc.ListResources(context.Background(), id, ListOptions{})
	if err == nil {
		t.Fatal("expected error for namespace-scoped resource without namespace, got nil")
	}

	var invalidInput *ErrInvalidInput
	if !errors.As(err, &invalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %T: %v", err, err)
	}
	if !contains(err.Error(), "requires a namespace") {
		t.Errorf("error message = %q, want to contain 'requires a namespace'", err.Error())
	}
}

// TestResourceUseCase_ApplyResource_NamespaceValidation tests namespace scope validation in ApplyResource.
func TestResourceUseCase_ApplyResource_NamespaceValidation(t *testing.T) {
	disco := &mockDiscoveryForResource{
		lookupNamespaced: false, // cluster-scoped
	}
	repo := &mockResourceRepo{}
	uc := NewResourceUseCase(disco, repo, nil)

	id := &ResourceIdentifier{
		Cluster:   "test-cluster",
		Group:     "apiextensions.k8s.io",
		Version:   "v1",
		Resource:  "customresourcedefinitions",
		Namespace: "kube-system",
		Name:      "test-crd",
	}

	_, err := uc.ApplyResource(context.Background(), id, []byte("manifest"), ApplyOptions{})
	if err == nil {
		t.Fatal("expected error for cluster-scoped resource with namespace, got nil")
	}

	var invalidInput *ErrInvalidInput
	if !errors.As(err, &invalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %T: %v", err, err)
	}
}

// TestResourceUseCase_DeleteResource_NamespaceValidation tests namespace scope validation in DeleteResource.
func TestResourceUseCase_DeleteResource_NamespaceValidation(t *testing.T) {
	disco := &mockDiscoveryForResource{
		lookupNamespaced: true, // namespace-scoped
	}
	repo := &mockResourceRepo{}
	uc := NewResourceUseCase(disco, repo, nil)

	id := &ResourceIdentifier{
		Cluster:  "test-cluster",
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
		Name:     "test-deploy",
		// Missing namespace
	}

	err := uc.DeleteResource(context.Background(), id, DeleteOptions{})
	if err == nil {
		t.Fatal("expected error for namespace-scoped resource without namespace, got nil")
	}

	var invalidInput *ErrInvalidInput
	if !errors.As(err, &invalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %T: %v", err, err)
	}
}

// TestResourceUseCase_WatchResource_NamespaceValidation tests namespace scope validation in WatchResource.
func TestResourceUseCase_WatchResource_NamespaceValidation(t *testing.T) {
	disco := &mockDiscoveryForResource{
		lookupNamespaced: false, // cluster-scoped
	}
	repo := &mockResourceRepo{}
	uc := NewResourceUseCase(disco, repo, nil)

	id := &ResourceIdentifier{
		Cluster:   "test-cluster",
		Group:     "",
		Version:   "v1",
		Resource:  "nodes",
		Namespace: "default", // should fail
	}

	_, err := uc.WatchResource(context.Background(), id, WatchOptions{})
	if err == nil {
		t.Fatal("expected error for cluster-scoped resource with namespace, got nil")
	}

	var invalidInput *ErrInvalidInput
	if !errors.As(err, &invalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %T: %v", err, err)
	}
}

// mockResourceRepo is a minimal mock for ResourceRepo.
type mockResourceRepo struct{}

func (m *mockResourceRepo) List(context.Context, string, schema.GroupVersionResource, string, ListOptions) (*unstructured.UnstructuredList, error) {
	return &unstructured.UnstructuredList{}, nil
}

func (m *mockResourceRepo) Get(context.Context, string, schema.GroupVersionResource, string, string) (*unstructured.Unstructured, error) {
	return &unstructured.Unstructured{}, nil
}

func (m *mockResourceRepo) Create(context.Context, string, schema.GroupVersionResource, string, []byte) (*unstructured.Unstructured, error) {
	return &unstructured.Unstructured{}, nil
}

func (m *mockResourceRepo) Apply(context.Context, string, schema.GroupVersionResource, string, string, []byte, ApplyOptions) (*unstructured.Unstructured, error) {
	return &unstructured.Unstructured{}, nil
}

func (m *mockResourceRepo) Delete(context.Context, string, schema.GroupVersionResource, string, string, DeleteOptions) error {
	return nil
}

func (m *mockResourceRepo) Watch(context.Context, string, schema.GroupVersionResource, string, WatchOptions) (Watcher, error) {
	return nil, nil
}

func (m *mockResourceRepo) ListEvents(context.Context, string, string, ListOptions) (*unstructured.UnstructuredList, error) {
	return &unstructured.UnstructuredList{}, nil
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
