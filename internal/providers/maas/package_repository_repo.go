package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/configuration"
)

type packageRepositoryRepo struct {
	maas *MAAS
}

func NewPackageRepositoryRepo(maas *MAAS) configuration.PackageRepositoryRepo {
	return &packageRepositoryRepo{
		maas: maas,
	}
}

var _ configuration.PackageRepositoryRepo = (*packageRepositoryRepo)(nil)

func (r *packageRepositoryRepo) List(_ context.Context) ([]configuration.PackageRepository, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	repos, err := client.PackageRepositories.Get()
	if err != nil {
		return nil, err
	}

	return r.toPackageRepositories(repos), nil
}

func (r *packageRepositoryRepo) Update(_ context.Context, id int, url string) (*configuration.PackageRepository, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.PackageRepositoryParams{
		URL: url,
	}

	repo, err := client.PackageRepository.Update(id, params)
	if err != nil {
		return nil, err
	}

	return r.toPackageRepository(repo), nil
}

func (r *packageRepositoryRepo) toPackageRepository(pr *entity.PackageRepository) *configuration.PackageRepository {
	return &configuration.PackageRepository{}
}

func (r *packageRepositoryRepo) toPackageRepositories(prs []entity.PackageRepository) []configuration.PackageRepository {
	ret := make([]configuration.PackageRepository, 0, len(prs))

	return ret
}
