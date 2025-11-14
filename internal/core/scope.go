package core

import (
	"context"
	"errors"

	"github.com/canonical/gomaasclient/entity"
	"github.com/juju/juju/api/base"
)

type (
	Scope  = base.UserModelSummary
	SSHKey = entity.SSHKey
)

type KeyRepo interface {
	Add(ctx context.Context, scope, key string) error
}

type ScopeRepo interface {
	List(ctx context.Context) ([]Scope, error)
	Get(ctx context.Context, name string) (*Scope, error)
	Create(ctx context.Context, name string, url string) (*Scope, error)
}

type ScopeConfigRepo interface {
	List(ctx context.Context, scope string) (map[string]any, error)
	Set(ctx context.Context, scope string, config map[string]any) error
	Unset(ctx context.Context, scope string, keys ...string) error
}

type SSHKeyRepo interface {
	List(ctx context.Context) ([]SSHKey, error)
}

type ScopeUseCase struct {
	key               KeyRepo
	scope             ScopeRepo
	sshKey            SSHKeyRepo
	packageRepository PackageRepositoryRepo
}

func NewScopeUseCase(key KeyRepo, scope ScopeRepo, sshKey SSHKeyRepo, packageRepository PackageRepositoryRepo) *ScopeUseCase {
	return &ScopeUseCase{
		key:               key,
		scope:             scope,
		sshKey:            sshKey,
		packageRepository: packageRepository,
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
	prs, err := uc.packageRepository.List(ctx)
	if err != nil {
		return nil, err
	}
	if len(prs) == 0 {
		return nil, errors.New("no package repositories configured; please configure at least one before creating a scope")
	}
	var url string
	for i := range prs {
		if prs[i].Name == "Ubuntu archive" {
			url = prs[i].URL
			break
		}
	}
	if url == "" {
		return nil, errors.New("ubuntu archive package repository not found; please ensure a package repository named 'Ubuntu archive' is configured")
	}

	scope, err := uc.scope.Create(ctx, name, url)
	if err != nil {
		return nil, err
	}
	if err := uc.key.Add(ctx, scope.Name, sshKey); err != nil {
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
