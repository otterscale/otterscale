package scope

import (
	"context"
)

type Scope struct {
	UUID string
	Name string
}

type ScopeRepo interface {
	List(ctx context.Context) ([]Scope, error)
	Get(ctx context.Context, name string) (*Scope, error)
	Create(ctx context.Context, name, sshKey string) (*Scope, error)
}

type SSHKeyRepo interface {
	First(ctx context.Context) (string, error)
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
	sshKey, err := uc.sshKey.First(ctx)
	if err != nil {
		return nil, err
	}

	return uc.scope.Create(ctx, name, sshKey)
}
