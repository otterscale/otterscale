package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core"
)

type packageRepository struct {
	maas *MAAS
}

func NewPackageRepository(maas *MAAS) core.PackageRepositoryRepo {
	return &packageRepository{
		maas: maas,
	}
}

var _ core.PackageRepositoryRepo = (*packageRepository)(nil)

func (r *packageRepository) List(_ context.Context) ([]core.PackageRepository, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.PackageRepositories.Get()
}

func (r *packageRepository) Update(_ context.Context, id int, params *entity.PackageRepositoryParams) (*core.PackageRepository, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.PackageRepository.Update(id, params)
}
