package core

import (
	"context"
	"testing"

	"github.com/juju/juju/api/base"
	"github.com/stretchr/testify/assert"
)

// Mock ScopeRepo
type mockScopeRepo struct {
	scopes []Scope
}

func (m *mockScopeRepo) List(ctx context.Context) ([]Scope, error) {
	return m.scopes, nil
}

func (m *mockScopeRepo) Create(ctx context.Context, name string) (*Scope, error) {
	return &base.UserModelSummary{Name: name, UUID: "uuid-123"}, nil
}

// Mock KeyRepo
type mockKeyRepo struct {
	added []string
}

func (m *mockKeyRepo) Add(ctx context.Context, uuid, key string) error {
	m.added = append(m.added, key)
	return nil
}

// Mock SSHKeyRepo
type mockSSHKeyRepo struct {
	keys []SSHKey
}

func (m *mockSSHKeyRepo) List(ctx context.Context) ([]SSHKey, error) {
	return m.keys, nil
}

func TestScopeUseCase_ListScopes(t *testing.T) {
	mockScope := &mockScopeRepo{
		scopes: []Scope{
			{Name: "scope1", UUID: "uuid1"},
			{Name: "scope2", UUID: "uuid2"},
		},
	}
	uc := NewScopeUseCase(mockScope, &mockKeyRepo{}, &mockSSHKeyRepo{})
	scopes, err := uc.ListScopes(context.Background())
	assert.NoError(t, err)
	assert.Len(t, scopes, 2)
	assert.Equal(t, "scope1", scopes[0].Name)
}

func TestScopeUseCase_CreateScope(t *testing.T) {
	mockScope := &mockScopeRepo{}
	mockKey := &mockKeyRepo{}
	mockSSH := &mockSSHKeyRepo{keys: []SSHKey{{Key: "ssh-rsa AAA"}}}
	uc := NewScopeUseCase(mockScope, mockKey, mockSSH)
	scope, err := uc.CreateScope(context.Background(), "testscope")
	assert.NoError(t, err)
	assert.Equal(t, "testscope", scope.Name)
	assert.Equal(t, "ssh-rsa AAA", mockKey.added[0])
}

func TestScopeUseCase_CreateScope_NoSSHKey(t *testing.T) {
	mockScope := &mockScopeRepo{}
	mockKey := &mockKeyRepo{}
	mockSSH := &mockSSHKeyRepo{keys: []SSHKey{}}
	uc := NewScopeUseCase(mockScope, mockKey, mockSSH)
	scope, err := uc.CreateScope(context.Background(), "testscope")
	assert.Error(t, err)
	assert.Nil(t, scope)
}
