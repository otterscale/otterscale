package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/openhdc/openhdc/internal/domain/model"
)

type NexusService struct {
	server            MAASServer
	packageRepository MAASPackageRepository
	bootResource      MAASBootResource
}

func NewNexusService(server MAASServer, packageRepository MAASPackageRepository, bootResource MAASBootResource) *NexusService {
	return &NexusService{
		server:            server,
		packageRepository: packageRepository,
		bootResource:      bootResource,
	}
}

func (s *NexusService) GetConfiguration(ctx context.Context) (*model.Configuration, error) {
	// ntp server
	cfg, err := s.server.Get(ctx, "ntp_servers")
	if err != nil {
		return nil, err
	}
	ntpServers := strings.Split(removeQuotes(cfg), " ")

	// boot resource
	defaultDistro, err := s.server.Get(ctx, "default_distro_series")
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
		if token[0] != "ubuntu" {
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
		}
	}
	brs := []*model.BootResource{}
	for _, br := range brm {
		brs = append(brs, br)
	}

	// package repository
	prs, err := s.packageRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return &model.Configuration{
		NTPServers:          ntpServers,
		BootResources:       brs,
		PackageRepositories: prs,
	}, nil
}

func ubuntuDistro(name, arch string) string {
	// version
	re := regexp.MustCompile(`-(\d{2}\.\d{2})`)
	match := re.FindStringSubmatch(arch)
	if len(match) == 0 {
		return name
	}

	version := match[1]
	token := strings.Split(version, ".")
	if len(token) == 1 {
		return name
	}
	minor := token[1]
	main, err := strconv.Atoi(token[0])
	if err != nil {
		return name
	}

	// code name
	token = strings.Split(name, "/")
	if len(token) == 1 {
		return name
	}
	codeName := token[1]

	isLTS := main%2 == 0 && minor == "04"
	if isLTS {
		return fmt.Sprintf("Ubuntu %s LTS (%s)", version, codeName)
	}
	return fmt.Sprintf("Ubuntu %s (%s)", version, codeName)
}

func removeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, ``)
}
