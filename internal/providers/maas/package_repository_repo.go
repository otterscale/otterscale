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

	return client.PackageRepositories.Get()
}

func (r *packageRepositoryRepo) Update(_ context.Context, id int, url string) (*configuration.PackageRepository, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.PackageRepositoryParams{
		URL: url,
	}

	return client.PackageRepository.Update(id, params)
}
