package model

import "github.com/canonical/gomaasclient/entity"

type PackageRepository = entity.PackageRepository

type Configuration struct {
	NTPServers          []string
	PackageRepositories []*PackageRepository
	BootResources       []*BootResource
}

type BootResource struct {
	Name         string
	Architecture string
	Status       string
	Default      bool
}
