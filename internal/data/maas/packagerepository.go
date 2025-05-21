package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type packageRepository struct {
	maas *MAAS
}

func NewPackageRepository(maas *MAAS) service.MAASPackageRepository {
	return &packageRepository{
		maas: maas,
	}
}

var _ service.MAASPackageRepository = (*packageRepository)(nil)

func (r *packageRepository) List(_ context.Context) ([]entity.PackageRepository, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.PackageRepositories.Get()
}

func (r *packageRepository) Update(_ context.Context, id int, params *entity.PackageRepositoryParams) (*entity.PackageRepository, error) {
	client, err := r.maas.client()
	if err != nil {
		return nil, err
	}
	return client.PackageRepository.Update(id, params)
}
