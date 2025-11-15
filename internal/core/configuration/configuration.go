package configuration

import (
	"context"

	"github.com/otterscale/otterscale/internal/config"
)

type Configuration struct {
	NTPServers          []string
	PackageRepositories []PackageRepository
	BootImages          []BootImage
	HelmRepositorys     []string
}

type UseCase struct {
	conf *config.Config

	bootResource        BootResourceRepo
	bootSource          BootSourceRepo
	bootSourceSelection BootSourceSelectionRepo
	packageRepository   PackageRepositoryRepo
	provisioner         ProvisionerRepo
	scopeConfig         ScopeConfigRepo
}

func NewUseCase(conf *config.Config, bootResource BootResourceRepo, bootSource BootSourceRepo, bootSourceSelection BootSourceSelectionRepo, packageRepository PackageRepositoryRepo, provisioner ProvisionerRepo, scopeConfig ScopeConfigRepo) *UseCase {
	return &UseCase{
		conf:                conf,
		bootResource:        bootResource,
		bootSource:          bootSource,
		bootSourceSelection: bootSourceSelection,
		packageRepository:   packageRepository,
		provisioner:         provisioner,
		scopeConfig:         scopeConfig,
	}
}

func (uc *UseCase) GetConfiguration(ctx context.Context) (*Configuration, error) {
	ntpServers, err := uc.listNTPServers(ctx)
	if err != nil {
		return nil, err
	}

	packageRepositories, err := uc.packageRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	bootImages, err := uc.listBootImages(ctx)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		NTPServers:          ntpServers,
		PackageRepositories: packageRepositories,
		BootImages:          bootImages,
		HelmRepositorys:     uc.conf.Kube.HelmRepositoryURLs,
	}, nil
}

func (uc *UseCase) UpdateHelmRepository(urls []string) ([]string, error) {
	uc.conf.Kube.HelmRepositoryURLs = urls

	if err := uc.conf.Override(uc.conf); err != nil {
		return nil, err
	}

	return uc.conf.Kube.HelmRepositoryURLs, nil
}
