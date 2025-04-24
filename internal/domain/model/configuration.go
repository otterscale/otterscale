package model

import (
	"github.com/canonical/gomaasclient/entity"
	"github.com/juju/juju/core/base"
)

type (
	PackageRepository       = entity.PackageRepository
	PackageRepositoryParams = entity.PackageRepositoryParams
)

type Configuration struct {
	NTPServers          []string
	PackageRepositories []PackageRepository
	BootImages          []BootImage
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
