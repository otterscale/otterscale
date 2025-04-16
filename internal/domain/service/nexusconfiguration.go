package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/openhdc/openhdc/internal/domain/model"
)

const (
	maasConfigNTPServers                = "ntp_servers"
	maasConfigDefaultOSSystem           = "default_osystem"
	maasConfigDefaultDistroSeries       = "default_distro_series"
	maasConfigCommissioningDistroSeries = "commissioning_distro_series"
	jujuConfigAPTMirror                 = "apt-mirror"
)

const (
	defaultOSSystem = "ubuntu"
)

func (s *NexusService) GetConfiguration(ctx context.Context) (*model.Configuration, error) {
	ntpServers, err := s.listNTPServers(ctx)
	if err != nil {
		return nil, err
	}
	brs, err := s.listBootResources(ctx)
	if err != nil {
		return nil, err
	}
	prs, err := s.listPackageRepositories(ctx)
	if err != nil {
		return nil, err
	}
	return &model.Configuration{
		NTPServers:          ntpServers,
		BootResources:       brs,
		PackageRepositories: prs,
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
		for _, ums := range umss {
			if err := s.scopeConfig.Set(ctx, ums.UUID, cfg); err != nil {
				return nil, err
			}
		}
	}
	return pr, nil
}

func (s *NexusService) UpdateDefaultBootResource(ctx context.Context, distroSeries string) (*model.BootResource, error) {
	if err := s.server.Update(ctx, maasConfigDefaultOSSystem, defaultOSSystem); err != nil {
		return nil, err
	}
	if err := s.server.Update(ctx, maasConfigDefaultDistroSeries, distroSeries); err != nil {
		return nil, err
	}
	if err := s.server.Update(ctx, maasConfigCommissioningDistroSeries, distroSeries); err != nil {
		return nil, err
	}
	brs, err := s.listBootResources(ctx)
	if err != nil {
		return nil, err
	}
	for _, br := range brs {
		if br.DistroSeries != distroSeries {
			continue
		}
		return br, nil
	}
	return nil, fmt.Errorf("distro series %q not found", distroSeries)
}

func (s *NexusService) SyncBootResources(ctx context.Context) error {
	return s.bootResource.Import(ctx)
}

func (s *NexusService) listNTPServers(ctx context.Context) ([]string, error) {
	cfg, err := s.server.Get(ctx, maasConfigNTPServers)
	if err != nil {
		return nil, err
	}
	return strings.Split(removeQuotes(cfg), " "), nil
}

func (s *NexusService) listBootResources(ctx context.Context) ([]*model.BootResource, error) {
	defaultDistro, err := s.server.Get(ctx, maasConfigDefaultDistroSeries)
	if err != nil {
		return nil, err
	}
	ebrs, err := s.bootResource.List(ctx)
	if err != nil {
		return nil, err
	}
	brm := map[string]*model.BootResource{}
	for _, br := range ebrs {
		token := strings.Split(br.Name, "/")
		if token[0] != defaultOSSystem {
			continue
		}
		isDefault := false
		if len(token) > 1 {
			isDefault = token[1] == removeQuotes(defaultDistro)
		}
		arch := strings.Split(br.Architecture, "/")[0]
		group := br.Name + arch
		brm[group] = &model.BootResource{
			Name:         ubuntuDistro(br.Name, br.Architecture),
			Architecture: arch,
			Status:       br.Type,
			Default:      isDefault,
			DistroSeries: ubuntuDistroSeries(br.Name),
		}
	}
	brs := []*model.BootResource{}
	for _, br := range brm {
		brs = append(brs, br)
	}
	return brs, nil
}

func (s *NexusService) listPackageRepositories(ctx context.Context) ([]*model.PackageRepository, error) {
	return s.packageRepository.List(ctx)
}

func ubuntuDistro(name, arch string) string {
	distroSeries := ubuntuDistroSeries(name)
	version := ubuntuVersion(arch)

	suffix := ""
	if isLTS(arch) {
		suffix = " LTS"
	}

	return fmt.Sprintf("Ubuntu %s%s (%s)", version, suffix, distroSeries)
}

func ubuntuVersion(arch string) string {
	re := regexp.MustCompile(`-(\d{2}\.\d{2})`)
	match := re.FindStringSubmatch(arch)
	if len(match) == 0 {
		return ""
	}
	return match[1]
}

func ubuntuDistroSeries(name string) string {
	token := strings.Split(name, "/")
	if len(token) == 1 {
		return name
	}
	return token[1]
}

func isLTS(arch string) bool {
	token := strings.Split(ubuntuVersion(arch), ".")
	if len(token) == 1 {
		return false
	}
	minor := token[1]
	main, err := strconv.Atoi(token[0])
	if err != nil {
		return false
	}
	return main%2 == 0 && minor == "04"
}

func removeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, ``)
}
