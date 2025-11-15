package scope

import (
	"context"
	"fmt"

	"github.com/juju/juju/api/base"
)

const ReservedName = "otterscale"

// Scope represents a Juju UserModelSummary resource.
type Scope = base.UserModelSummary

//nolint:revive // allows this exported interface name for specific domain clarity.
type ScopeRepo interface {
	List(ctx context.Context) ([]Scope, error)
	Get(ctx context.Context, name string) (*Scope, error)
	Create(ctx context.Context, name, sshKey string) (*Scope, error)
}

type SSHKeyRepo interface {
	First(ctx context.Context) (string, error)
}

type UseCase struct {
	scope  ScopeRepo
	sshKey SSHKeyRepo
}

func NewUseCase(scope ScopeRepo, sshKey SSHKeyRepo) *UseCase {
	return &UseCase{
		scope:  scope,
		sshKey: sshKey,
	}
}

func (uc *UseCase) ListScopes(ctx context.Context) ([]Scope, error) {
	return uc.scope.List(ctx)
}

func (uc *UseCase) CreateScope(ctx context.Context, name string) (*Scope, error) {
	if name == ReservedName {
		return nil, fmt.Errorf("scope name %q is reserved", name)
	}

	sshKey, err := uc.sshKey.First(ctx)
	if err != nil {
		return nil, err
	}

	return uc.scope.Create(ctx, name, sshKey)
}
