package configuration

import (
	"context"

	"github.com/canonical/gomaasclient/entity"
)

// PackageRepository represents a MAAS PackageRepository resource.
type PackageRepository = entity.PackageRepository

type PackageRepositoryRepo interface {
	List(ctx context.Context) ([]PackageRepository, error)
	Update(ctx context.Context, id int, url string) (*PackageRepository, error)
}

type ScopeConfigRepo interface {
	Set(ctx context.Context, config map[string]any) error
}

func (uc *ConfigurationUseCase) UpdatePackageRepository(ctx context.Context, id int, url string, skipJuju bool) (*PackageRepository, error) {
	packageRepository, err := uc.packageRepository.Update(ctx, id, url)
	if err != nil {
		return nil, err
	}

	if !skipJuju {
		if err := uc.scopeConfig.Set(ctx, map[string]any{"apt-mirror": url}); err != nil {
			return nil, err
		}
	}

	return packageRepository, nil
}
