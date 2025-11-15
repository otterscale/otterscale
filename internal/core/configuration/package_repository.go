package configuration

import (
	"context"

	"github.com/canonical/gomaasclient/entity"
)

const MAASUbuntuArchive = "Ubuntu archive"

// PackageRepository represents a MAAS PackageRepository resource.
type PackageRepository = entity.PackageRepository

type PackageRepositoryRepo interface {
	List(ctx context.Context) ([]PackageRepository, error)
	Update(ctx context.Context, id int, url string) (*PackageRepository, error)
}

type ScopeConfigRepo interface {
	Set(ctx context.Context, config map[string]any) error
}

func (uc *UseCase) UpdatePackageRepository(ctx context.Context, id int, url string) (*PackageRepository, error) {
	packageRepository, err := uc.packageRepository.Update(ctx, id, url)
	if err != nil {
		return nil, err
	}

	if packageRepository.Name == MAASUbuntuArchive {
		if err := uc.scopeConfig.Set(ctx, map[string]any{"apt-mirror": url}); err != nil {
			return nil, err
		}
	}

	return packageRepository, nil
}
