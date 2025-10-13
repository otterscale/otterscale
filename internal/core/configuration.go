package core

import (
	"context"
	"strings"

	"github.com/canonical/gomaasclient/entity"
	"github.com/juju/juju/core/base"

	"github.com/otterscale/otterscale/internal/config"
)

var ubuntuDistroSeriesMap = map[base.SeriesName]BootImageSelection{
	base.Xenial: {
		DistroSeries:  base.Xenial,
		Name:          "Ubuntu 16.04 LTS Xenial Xerus",
		Architectures: []string{"amd64", "arm64", "armhf", "i386", "ppc64el", "s390x"},
	},
	base.Bionic: {
		DistroSeries:  base.Bionic,
		Name:          "Ubuntu 18.04 LTS Bionic Beaver",
		Architectures: []string{"amd64", "arm64", "armhf", "i386", "ppc64el", "s390x"},
	},
	base.Focal: {
		DistroSeries:  base.Focal,
		Name:          "Ubuntu 20.04 LTS Focal Fossa",
		Architectures: []string{"amd64", "arm64", "armhf", "ppc64el", "s390x"},
	},
	base.Jammy: {
		DistroSeries:  base.Jammy,
		Name:          "Ubuntu 22.04 LTS Jammy Jellyfish",
		Architectures: []string{"amd64", "arm64", "armhf", "ppc64el", "s390x"},
	},
	base.Noble: {
		DistroSeries:  base.Noble,
		Name:          "Ubuntu 24.04 LTS Noble Numbat",
		Architectures: []string{"amd64", "arm64", "armhf", "ppc64el", "s390x"},
	},
}

type Configuration struct {
	NTPServers          []string
	PackageRepositories []PackageRepository
	BootImages          []BootImage
	HelmRepositorys     []string
}

type BootImage struct {
	Source                string
	DistroSeries          string
	Name                  string
	ArchitectureStatusMap map[string]string
	Default               bool
}

type BootImageSelection struct {
	DistroSeries  base.SeriesName
	Name          string
	Architectures []string
}

type PackageRepository = entity.PackageRepository

type BootResourceRepo interface {
	List(ctx context.Context) ([]entity.BootResource, error)
	Import(ctx context.Context) error
	IsImporting(ctx context.Context) (bool, error)
}

type BootSourceRepo interface {
	List(ctx context.Context) ([]entity.BootSource, error)
}

type BootSourceSelectionRepo interface {
	List(ctx context.Context, id int) ([]entity.BootSourceSelection, error)
	Create(ctx context.Context, bootSourceID int, params *entity.BootSourceSelectionParams) (*entity.BootSourceSelection, error)
}

type PackageRepositoryRepo interface {
	List(ctx context.Context) ([]PackageRepository, error)
	Update(ctx context.Context, id int, params *entity.PackageRepositoryParams) (*PackageRepository, error)
}

type ConfigurationUseCase struct {
	conf                *config.Config
	bootResource        BootResourceRepo
	bootSource          BootSourceRepo
	bootSourceSelection BootSourceSelectionRepo
	packageRepository   PackageRepositoryRepo
	scope               ScopeRepo
	scopeConfig         ScopeConfigRepo
	server              ServerRepo
}

func NewConfigurationUseCase(conf *config.Config, bootResource BootResourceRepo, bootSource BootSourceRepo, bootSourceSelection BootSourceSelectionRepo, packageRepository PackageRepositoryRepo, scope ScopeRepo, scopeConfig ScopeConfigRepo, server ServerRepo) *ConfigurationUseCase {
	return &ConfigurationUseCase{
		conf:                conf,
		bootResource:        bootResource,
		bootSource:          bootSource,
		bootSourceSelection: bootSourceSelection,
		packageRepository:   packageRepository,
		scope:               scope,
		scopeConfig:         scopeConfig,
		server:              server,
	}
}

func (uc *ConfigurationUseCase) GetConfiguration(ctx context.Context) (*Configuration, error) {
	ntpServers, err := uc.listNTPServers(ctx)
	if err != nil {
		return nil, err
	}
	packageRepositories, err := uc.listPackageRepositories(ctx)
	if err != nil {
		return nil, err
	}
	bootImages, err := uc.listBootImages(ctx)
	if err != nil {
		return nil, err
	}
	helmRepositories := uc.listHelmRepositories()
	return &Configuration{
		NTPServers:          ntpServers,
		PackageRepositories: packageRepositories,
		BootImages:          bootImages,
		HelmRepositorys:     helmRepositories,
	}, nil
}

func (uc *ConfigurationUseCase) UpdateNTPServer(ctx context.Context, addresses []string) ([]string, error) {
	if err := uc.server.Update(ctx, "ntp_servers", strings.Join(addresses, " ")); err != nil {
		return nil, err
	}
	return uc.listNTPServers(ctx)
}

func (uc *ConfigurationUseCase) UpdatePackageRepository(ctx context.Context, id int, url string, skipJuju bool) (*PackageRepository, error) {
	params := &entity.PackageRepositoryParams{
		URL: url,
	}
	packageRepository, err := uc.packageRepository.Update(ctx, id, params)
	if err != nil {
		return nil, err
	}

	if !skipJuju {
		scopes, err := uc.scope.List(ctx)
		if err != nil {
			return nil, err
		}
		for i := range scopes {
			if err := uc.scopeConfig.Set(ctx, scopes[i].UUID, map[string]any{"apt-mirror": url}); err != nil {
				return nil, err
			}
		}
	}

	return packageRepository, nil
}

func (uc *ConfigurationUseCase) UpdateHelmRepository(urls []string) ([]string, error) {
	uc.conf.Kube.HelmRepositoryURLs = urls
	if err := uc.conf.Override(uc.conf); err != nil {
		return nil, err
	}
	return uc.listHelmRepositories(), nil
}

func (uc *ConfigurationUseCase) CreateBootImage(ctx context.Context, distroSeries string, architectures []string) (*BootImage, error) {
	if len(architectures) == 0 {
		architectures = []string{"amd64"} // set default
	}

	maasIO := 1
	params := &entity.BootSourceSelectionParams{
		OS:        "ubuntu",
		Release:   distroSeries,
		Arches:    architectures,
		Subarches: []string{"*"},
		Labels:    []string{"*"},
	}
	selections, err := uc.bootSourceSelection.Create(ctx, maasIO, params)
	if err != nil {
		return nil, err
	}

	statusMap := map[string]string{}
	for _, arch := range selections.Arches {
		statusMap[arch] = ""
	}

	return &BootImage{
		DistroSeries:          selections.Release,
		Name:                  selections.OS,
		ArchitectureStatusMap: statusMap,
	}, nil
}

func (uc *ConfigurationUseCase) SetDefaultBootImage(ctx context.Context, distroSeries string) error {
	if err := uc.server.Update(ctx, "default_osystem", "ubuntu"); err != nil {
		return err
	}
	if err := uc.server.Update(ctx, "default_distro_series", distroSeries); err != nil {
		return err
	}
	return uc.server.Update(ctx, "commissioning_distro_series", distroSeries)
}

func (uc *ConfigurationUseCase) ImportBootImages(ctx context.Context) error {
	return uc.bootResource.Import(ctx)
}

func (uc *ConfigurationUseCase) IsImportingBootImages(ctx context.Context) (bool, error) {
	return uc.bootResource.IsImporting(ctx)
}

func (uc *ConfigurationUseCase) ListBootImageSelections() ([]BootImageSelection, error) {
	selections := []BootImageSelection{}
	for distro := range ubuntuDistroSeriesMap {
		selections = append(selections, ubuntuDistroSeriesMap[distro])
	}
	return selections, nil
}

func (uc *ConfigurationUseCase) listNTPServers(ctx context.Context) ([]string, error) {
	ntpServers, err := uc.server.Get(ctx, "ntp_servers")
	if err != nil {
		return nil, err
	}
	return strings.Split(ntpServers, " "), nil
}

func (uc *ConfigurationUseCase) listPackageRepositories(ctx context.Context) ([]PackageRepository, error) {
	return uc.packageRepository.List(ctx)
}

func (uc *ConfigurationUseCase) listBootImages(ctx context.Context) ([]BootImage, error) {
	defaultDistro, err := uc.server.Get(ctx, "default_distro_series")
	if err != nil {
		return nil, err
	}

	resources, err := uc.bootResource.List(ctx)
	if err != nil {
		return nil, err
	}
	statusMaps := map[string]map[string]string{}
	for i := range resources {
		token := strings.Split(resources[i].Name, "/")
		if len(token) > 1 {
			distro := token[1]
			if _, ok := statusMaps[distro]; !ok {
				statusMaps[distro] = map[string]string{}
			}
			statusMaps[distro][resources[i].Architecture] = resources[i].Type
		}
	}

	sources, err := uc.bootSource.List(ctx)
	if err != nil {
		return nil, err
	}
	bootImages := []BootImage{}
	for i := range sources {
		selections, err := uc.bootSourceSelection.List(ctx, sources[i].ID)
		if err != nil {
			return nil, err
		}
		for j := range selections {
			distro := selections[j].Release
			name := selections[j].Release
			if ds, ok := ubuntuDistroSeriesMap[base.SeriesName(selections[j].Release)]; ok {
				name = ds.Name
			}
			bootImages = append(bootImages, BootImage{
				Source:                sources[i].URL,
				DistroSeries:          distro,
				Name:                  name,
				ArchitectureStatusMap: statusMaps[distro],
				Default:               distro == defaultDistro,
			})
		}
	}
	return bootImages, nil
}

func (uc *ConfigurationUseCase) listHelmRepositories() []string {
	return uc.conf.Kube.HelmRepositoryURLs
}
