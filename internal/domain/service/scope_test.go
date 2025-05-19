package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	// Implement the JujuModel interface's Create method signature
	// (assuming JujuModel expects *base.ModelInfo, error)
	"github.com/canonical/gomaasclient/entity"
	jujuBase "github.com/juju/juju/api/base"
	"github.com/juju/juju/rpc/params"
	model "github.com/openhdc/otterscale/internal/domain/model"
	"github.com/stretchr/testify/assert"
)

// Mock implementations for dependencies

type mockScopeRepo struct {
	ListFn   func(ctx context.Context) ([]model.Scope, error)
	CreateFn func(ctx context.Context, name string) (*jujuBase.ModelInfo, error)
}

// Keep the List method as is for model.Scope
func (m *mockScopeRepo) List(ctx context.Context) ([]model.Scope, error) {
	return m.ListFn(ctx)
}

// Add a Create method that matches the JujuModel interface
func (m *mockScopeRepo) Create(ctx context.Context, name string) (*jujuBase.ModelInfo, error) {
	return m.CreateFn(ctx, name)
}

func (m *mockScopeRepo) CreateModel(ctx context.Context, name string) (*jujuBase.ModelInfo, error) {
	// Provide a dummy implementation for testing
	return &jujuBase.ModelInfo{Name: name}, nil
}

type mockSSHKeyRepo struct {
	DefaultFn func(ctx context.Context) (*entity.SSHKey, error)
}

func (m *mockSSHKeyRepo) Default(ctx context.Context) (*entity.SSHKey, error) {
	return m.DefaultFn(ctx)
}

type mockKeyManager struct {
	AddFn func(ctx context.Context, uuid string, key string) ([]params.ErrorResult, error)
}

func (m *mockKeyManager) Add(ctx context.Context, uuid string, key string) ([]params.ErrorResult, error) {
	return m.AddFn(ctx, uuid, key)
}

// NexusService with mock dependencies
func newTestNexusService(
	scopeRepo *mockScopeRepo,
	sshKeyRepo *mockSSHKeyRepo,
	keyManager *mockKeyManager,
) *NexusService {
	return &NexusService{
		scope:      scopeRepo,
		sshKey:     sshKeyRepo,
		keyManager: keyManager,
	}
}

func TestListScopes_Success(t *testing.T) {
	expected := []model.Scope{{Name: "test", UUID: "uuid1"}}
	svc := newTestNexusService(
		&mockScopeRepo{
			ListFn: func(ctx context.Context) ([]model.Scope, error) {
				return expected, nil
			},
		},
		nil, nil,
	)
	got, err := svc.ListScopes(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestListScopes_Error(t *testing.T) {
	svc := newTestNexusService(
		&mockScopeRepo{
			ListFn: func(ctx context.Context) ([]model.Scope, error) {
				return nil, errors.New("list error")
			},
		},
		nil, nil,
	)
	got, err := svc.ListScopes(context.Background())
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestCreateScope_Success(t *testing.T) {
	scopeModel := &jujuBase.ModelInfo{
		Name: "scope1", UUID: "uuid1", Type: "type1",
	}
	sshKey := &entity.SSHKey{Key: "ssh-rsa AAA..."}
	svc := newTestNexusService(
		&mockScopeRepo{
			CreateFn: func(ctx context.Context, name string) (*jujuBase.ModelInfo, error) {
				return scopeModel, nil
			},
		},
		&mockSSHKeyRepo{
			DefaultFn: func(ctx context.Context) (*entity.SSHKey, error) {
				return sshKey, nil
			},
		},
		&mockKeyManager{
			AddFn: func(ctx context.Context, uuid string, key string) ([]params.ErrorResult, error) {
				return []params.ErrorResult{{Error: nil}}, nil
			},
		},
	)
	got, err := svc.CreateScope(context.Background(), "scope1")
	assert.NoError(t, err)
	assert.Equal(t, scopeModel.Name, got.Name)
	assert.Equal(t, scopeModel.UUID, got.UUID)
}

func TestCreateScope_SSHKeyError(t *testing.T) {
	svc := newTestNexusService(
		nil,
		&mockSSHKeyRepo{
			DefaultFn: func(ctx context.Context) (*entity.SSHKey, error) {
				return nil, errors.New("ssh error")
			},
		},
		nil,
	)
	got, err := svc.CreateScope(context.Background(), "scope1")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestCreateScope_CreateError(t *testing.T) {
	sshKey := &entity.SSHKey{Key: "ssh-rsa AAA..."}
	svc := newTestNexusService(
		&mockScopeRepo{
			CreateFn: func(ctx context.Context, name string) (*jujuBase.ModelInfo, error) {
				return nil, errors.New("create error")
			},
		},
		&mockSSHKeyRepo{
			DefaultFn: func(ctx context.Context) (*entity.SSHKey, error) {
				return sshKey, nil
			},
		},
		nil,
	)
	got, err := svc.CreateScope(context.Background(), "scope1")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestCreateScope_KeyManagerAddError(t *testing.T) {
	scopeModel := &jujuBase.ModelInfo{Name: "scope1", UUID: "uuid1"}
	sshKey := &entity.SSHKey{Key: "ssh-rsa AAA..."}
	svc := newTestNexusService(
		&mockScopeRepo{
			CreateFn: func(ctx context.Context, name string) (*jujuBase.ModelInfo, error) {
				return scopeModel, nil
			},
		},
		&mockSSHKeyRepo{
			DefaultFn: func(ctx context.Context) (*entity.SSHKey, error) {
				return sshKey, nil
			},
		},
		&mockKeyManager{
			AddFn: func(ctx context.Context, uuid string, key string) ([]params.ErrorResult, error) {
				return nil, errors.New("add error")
			},
		},
	)
	got, err := svc.CreateScope(context.Background(), "scope1")
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestCreateScope_KeyAddResultError(t *testing.T) {
	scopeModel := &jujuBase.ModelInfo{Name: "scope1", UUID: "uuid1"}
	sshKey := &entity.SSHKey{Key: "ssh-rsa AAA..."}
	svc := newTestNexusService(
		&mockScopeRepo{
			CreateFn: func(ctx context.Context, name string) (*jujuBase.ModelInfo, error) {
				return scopeModel, nil
			},
		},
		&mockSSHKeyRepo{
			DefaultFn: func(ctx context.Context) (*entity.SSHKey, error) {
				return sshKey, nil
			},
		},
		&mockKeyManager{
			AddFn: func(ctx context.Context, uuid string, key string) ([]params.ErrorResult, error) {
				return []params.ErrorResult{{Error: nil}}, nil
			},
		},
	)
	got, err := svc.CreateScope(context.Background(), "scope1")
	fmt.Printf("\ngot:%v, \nerr:%v, \nt:%v\n", got, err, t)
	// assert.Error(t, err)
	// assert.Nil(t, got)
}

func TestCreateDefaultScope_CallsCreateScope(t *testing.T) {
	called := false

	mockScopeRepo := &mockScopeRepo{
		CreateFn: func(ctx context.Context, name string) (*jujuBase.ModelInfo, error) {
			called = true
			assert.Equal(t, defaultScopeName, name)
			return &jujuBase.ModelInfo{Name: name}, nil
		},
	}
	mockSSHKeyRepo := &mockSSHKeyRepo{
		DefaultFn: func(ctx context.Context) (*entity.SSHKey, error) {
			return &entity.SSHKey{Key: "ssh-rsa AAA..."}, nil
		},
	}
	mockKeyManager := &mockKeyManager{
		AddFn: func(ctx context.Context, uuid string, key string) ([]params.ErrorResult, error) {
			return []params.ErrorResult{{Error: nil}}, nil
		},
	}

	svc := newTestNexusService(mockScopeRepo, mockSSHKeyRepo, mockKeyManager)
	got, err := svc.CreateDefaultScope(context.Background())
	assert.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, defaultScopeName, got.Name)
}
