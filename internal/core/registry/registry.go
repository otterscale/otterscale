package registry

import (
	"context"
)

type UseCase struct {
	manifest   ManifestRepo
	repository RepositoryRepo
}

func NewUseCase(manifest ManifestRepo, repository RepositoryRepo) *UseCase {
	return &UseCase{
		manifest:   manifest,
		repository: repository,
	}
}

func (uc *UseCase) GetRegistryURL(ctx context.Context, scope string) (string, error) {
	return uc.repository.GetRegistryURL(ctx, scope)
}

func (uc *UseCase) ListRepositories(ctx context.Context, scope string) ([]Repository, error) {
	return uc.repository.List(ctx, scope)
}

func (uc *UseCase) ListManifests(ctx context.Context, scope, repository string) ([]Manifest, error) {
	return uc.manifest.List(ctx, scope, repository)
}

func (uc *UseCase) DeleteManifest(ctx context.Context, scope, repository, digest string) error {
	return uc.manifest.Delete(ctx, scope, repository, digest)
}
