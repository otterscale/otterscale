package core

import (
	"context"
	"errors"

	"github.com/canonical/gomaasclient/entity"
	"github.com/juju/juju/api/base"
)

type Scope = base.UserModelSummary

type ScopeRepo interface {
	List(ctx context.Context) ([]Scope, error)
	Create(ctx context.Context, name string) (*Scope, error)
	AddSSHKey(ctx context.Context, uuid, sshKey string) error
}

type ScopeConfigRepo interface {
	List(ctx context.Context, uuid string) (map[string]any, error)
	Set(ctx context.Context, uuid string, config map[string]any) error
	Unset(ctx context.Context, uuid string, keys ...string) error
}

type SSHKey = entity.SSHKey

type SSHKeyRepo interface {
	List(ctx context.Context) ([]entity.SSHKey, error)
}

type ScopeUseCase struct {
	scope  ScopeRepo
	sshKey SSHKeyRepo
}

func NewScopeUseCase(scope ScopeRepo, sshKey SSHKeyRepo) *ScopeUseCase {
	return &ScopeUseCase{
		scope:  scope,
		sshKey: sshKey,
	}
}

func (uc *ScopeUseCase) ListScopes(ctx context.Context) ([]Scope, error) {
	return uc.scope.List(ctx)
}

func (uc *ScopeUseCase) CreateScope(ctx context.Context, name string) (*Scope, error) {
	sshKey, err := uc.firstSSHKey(ctx)
	if err != nil {
		return nil, err
	}
	scope, err := uc.scope.Create(ctx, name)
	if err != nil {
		return nil, err
	}
	if err := uc.scope.AddSSHKey(ctx, scope.UUID, sshKey); err != nil {
		return nil, err
	}
	return scope, nil
}

func (uc *ScopeUseCase) firstSSHKey(ctx context.Context) (string, error) {
	sshKeys, err := uc.sshKey.List(ctx)
	if err != nil {
		return "", err
	}
	if len(sshKeys) == 0 {
		return "", errors.New("ssh key not found")
	}
	return sshKeys[0].Key, nil
}
