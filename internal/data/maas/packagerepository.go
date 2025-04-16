package maas

import (
	"context"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/openhdc/internal/domain/service"
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
	return r.maas.PackageRepositories.Get()
}

func (r *packageRepository) Update(_ context.Context, id int, params *entity.PackageRepositoryParams) (*entity.PackageRepository, error) {
	return r.maas.PackageRepository.Update(id, params)
}
