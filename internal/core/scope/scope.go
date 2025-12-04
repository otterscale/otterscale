package scope

import (
	"context"
	"fmt"
	"strings"

	"github.com/juju/juju/api/base"

	"github.com/otterscale/otterscale/internal/core/configuration"
)

const ReservedName = "otterscale"

// Scope represents a Juju UserModelSummary resource.
type Scope = base.UserModelSummary

//nolint:revive // allows this exported interface name for specific domain clarity.
type ScopeRepo interface {
	List(ctx context.Context) ([]Scope, error)
	Get(ctx context.Context, name string) (*Scope, error)
	Create(ctx context.Context, name string, config map[string]any, sshKey string) (*Scope, error)
}

type SSHKeyRepo interface {
	First(ctx context.Context) (string, error)
}

type UseCase struct {
	scope  ScopeRepo
	sshKey SSHKeyRepo

	packageRepository configuration.PackageRepositoryRepo
}

func NewUseCase(scope ScopeRepo, sshKey SSHKeyRepo, packageRepository configuration.PackageRepositoryRepo) *UseCase {
	return &UseCase{
		scope:             scope,
		sshKey:            sshKey,
		packageRepository: packageRepository,
	}
}

func (uc *UseCase) ListScopes(ctx context.Context) ([]Scope, error) {
	return uc.scope.List(ctx)
}

func (uc *UseCase) CreateScope(ctx context.Context, name string) (*Scope, error) {
	if strings.EqualFold(name, ReservedName) {
		return nil, fmt.Errorf("scope name %q is reserved", name)
	}

	aptMirrorURL, err := uc.aptMirrorURL(ctx)
	if err != nil {
		return nil, err
	}

	config := map[string]any{
		"apt-mirror":                      aptMirrorURL,
		"num-container-provision-workers": 8,
	}

	sshKey, err := uc.sshKey.First(ctx)
	if err != nil {
		return nil, err
	}

	return uc.scope.Create(ctx, name, config, sshKey)
}

func (uc *UseCase) aptMirrorURL(ctx context.Context) (string, error) {
	packageRepositories, err := uc.packageRepository.List(ctx)
	if err != nil {
		return "", err
	}

	for i := range packageRepositories {
		if packageRepositories[i].Name == configuration.MAASUbuntuArchive {
			return packageRepositories[i].URL, nil
		}
	}

	return "", nil
}
