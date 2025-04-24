package service

import (
	"context"
	"strings"

	"github.com/juju/juju/core/base"

	"github.com/openhdc/openhdc/internal/domain/model"
)

const (
	maasConfigNTPServers                = "ntp_servers"
	maasConfigDefaultOSSystem           = "default_osystem"
	maasConfigDefaultDistroSeries       = "default_distro_series"
	maasConfigCommissioningDistroSeries = "commissioning_distro_series"
	jujuConfigAPTMirror                 = "apt-mirror"
)

const defaultOSSystem = "ubuntu"

var ubuntuDistroSeriesMap = map[base.SeriesName]model.BootImageSelection{
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

func (s *NexusService) GetConfiguration(ctx context.Context) (*model.Configuration, error) {
	ntpServers, err := s.listNTPServers(ctx)
	if err != nil {
		return nil, err
	}
	prs, err := s.listPackageRepositories(ctx)
	if err != nil {
		return nil, err
	}
	brs, err := s.listBootImages(ctx)
	if err != nil {
		return nil, err
	}
	return &model.Configuration{
		NTPServers:          ntpServers,
		PackageRepositories: prs,
		BootImages:          brs,
	}, nil
}

func (s *NexusService) UpdateNTPServer(ctx context.Context, addresses []string) ([]string, error) {
	if err := s.server.Update(ctx, maasConfigNTPServers, strings.Join(addresses, " ")); err != nil {
		return nil, err
	}
	return s.listNTPServers(ctx)
}

func (s *NexusService) UpdatePackageRepository(ctx context.Context, id int, url string, skipJuju bool) (*model.PackageRepository, error) {
	params := &model.PackageRepositoryParams{
		URL: url,
	}
	pr, err := s.packageRepository.Update(ctx, id, params)
	if err != nil {
		return nil, err
	}
	if !skipJuju {
		cfg := map[string]any{jujuConfigAPTMirror: url}
		umss, err := s.scope.List(ctx)
		if err != nil {
			return nil, err
		}
		for i := range umss {
			if err := s.scopeConfig.Set(ctx, umss[i].UUID, cfg); err != nil {
				return nil, err
			}
		}
	}
	return pr, nil
}

func (s *NexusService) CreateBootImage(ctx context.Context, distroSeries string, architectures []string) (*model.BootImage, error) {
	if len(architectures) == 0 {
		architectures = []string{"amd64"} // default
	}
	bss, err := s.bootSourceSelection.CreateFromMAASIO(ctx, distroSeries, architectures)
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	for _, arch := range bss.Arches {
		m[arch] = ""
	}
	return &model.BootImage{
		DistroSeries:          bss.Release,
		Name:                  bss.OS,
		ArchitectureStatusMap: m,
	}, nil
}

func (s *NexusService) SetDefaultBootImage(ctx context.Context, distroSeries string) error {
	if err := s.server.Update(ctx, maasConfigDefaultOSSystem, defaultOSSystem); err != nil {
		return err
	}
	if err := s.server.Update(ctx, maasConfigDefaultDistroSeries, distroSeries); err != nil {
		return err
	}
	return s.server.Update(ctx, maasConfigCommissioningDistroSeries, distroSeries)
}

func (s *NexusService) ImportBootImages(ctx context.Context) error {
	return s.bootResource.Import(ctx)
}

func (s *NexusService) IsImportingBootImages(ctx context.Context) (bool, error) {
	return s.bootResource.IsImporting(ctx)
}

func (s *NexusService) ListBootImageSelections(ctx context.Context) ([]model.BootImageSelection, error) {
	biss := []model.BootImageSelection{}
	for distro := range ubuntuDistroSeriesMap {
		biss = append(biss, ubuntuDistroSeriesMap[distro])
	}
	return biss, nil
}

func (s *NexusService) listNTPServers(ctx context.Context) ([]string, error) {
	val, err := s.server.Get(ctx, maasConfigNTPServers)
	if err != nil {
		return nil, err
	}
	return strings.Split(removeQuotes(string(val)), " "), nil
}

func (s *NexusService) listPackageRepositories(ctx context.Context) ([]model.PackageRepository, error) {
	return s.packageRepository.List(ctx)
}

func (s *NexusService) listBootImages(ctx context.Context) ([]model.BootImage, error) {
	val, err := s.server.Get(ctx, maasConfigDefaultDistroSeries)
	if err != nil {
		return nil, err
	}
	brs, err := s.bootResource.List(ctx)
	if err != nil {
		return nil, err
	}
	brm := map[string]map[string]string{}
	for i := range brs {
		token := strings.Split(brs[i].Name, "/")
		if len(token) > 1 {
			distro := token[1]
			if _, ok := brm[distro]; !ok {
				brm[distro] = map[string]string{}
			}
			brm[distro][brs[i].Architecture] = brs[i].Type
		}
	}
	bss, err := s.bootSource.List(ctx)
	if err != nil {
		return nil, err
	}
	bis := []model.BootImage{}
	for i := range bss {
		brss, err := s.bootSourceSelection.List(ctx, bss[i].ID)
		if err != nil {
			return nil, err
		}
		for j := range brss {
			distro := brss[j].Release
			name := brss[j].Release
			if ds, ok := ubuntuDistroSeriesMap[base.SeriesName(brss[j].Release)]; ok {
				name = ds.Name
			}
			bis = append(bis, model.BootImage{
				Source:                bss[i].URL,
				DistroSeries:          distro,
				Name:                  name,
				ArchitectureStatusMap: brm[distro],
				Default:               distro == removeQuotes(string(val)),
			})
		}
	}
	return bis, nil
}

func (s *NexusService) imageBase(ctx context.Context) (*base.Base, error) {
	val, err := s.server.Get(ctx, maasConfigDefaultDistroSeries)
	if err != nil {
		return nil, err
	}
	b, err := base.GetBaseFromSeries(removeQuotes(string(val)))
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func removeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, ``)
}
