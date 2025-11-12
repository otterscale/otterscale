package configuration

import (
	"context"
	"strings"

	"github.com/juju/juju/core/base"
)

type BootImage struct {
	Source                string
	DistroSeries          string
	Name                  string
	ID                    int
	Architectures         []string
	ArchitectureStatusMap map[string]string
	Default               bool
}

type BootResource struct {
	ID           int
	Type         string
	Name         string
	Architecture string
}

type BootSource struct {
	ID  int
	URL string
}

type BootSourceSelection struct {
	ID           int
	BootSourceID int
	Release      string
	OS           string
	Arches       []string
}

type BootResourceRepo interface {
	List(ctx context.Context) ([]BootResource, error)
	Import(ctx context.Context) error
	IsImporting(ctx context.Context) (bool, error)
}

type BootSourceRepo interface {
	List(ctx context.Context) ([]BootSource, error)
}

type BootSourceSelectionRepo interface {
	List(ctx context.Context, id int) ([]BootSourceSelection, error)
	Create(ctx context.Context, bootSourceID int, distroSeries string, architectures []string) (*BootSourceSelection, error)
	Update(ctx context.Context, bootSourceID int, id int, distroSeries string, architectures []string) (*BootSourceSelection, error)
}

func (uc *ConfigurationUseCase) CreateBootImage(ctx context.Context, distroSeries string, architectures []string) (*BootImage, error) {
	if len(architectures) == 0 {
		architectures = []string{"amd64"} // set default
	}

	maasIO := 1
	selections, err := uc.bootSourceSelection.Create(ctx, maasIO, distroSeries, architectures)
	if err != nil {
		return nil, err
	}

	statusMap := map[string]string{}
	for _, arch := range selections.Arches {
		statusMap[arch] = ""
	}

	return &BootImage{
		ID:                    selections.ID,
		DistroSeries:          selections.Release,
		Name:                  selections.OS,
		ArchitectureStatusMap: statusMap,
	}, nil
}

func (uc *ConfigurationUseCase) UpdateBootImage(ctx context.Context, id int, distroSeries string, architectures []string) (*BootImage, error) {
	if len(architectures) == 0 {
		architectures = []string{"amd64"} // set default
	}

	maasIO := 1
	selections, err := uc.bootSourceSelection.Update(ctx, maasIO, id, distroSeries, architectures)
	if err != nil {
		return nil, err
	}

	statusMap := map[string]string{}
	for _, arch := range selections.Arches {
		statusMap[arch] = ""
	}

	return &BootImage{
		ID:                    selections.ID,
		DistroSeries:          selections.Release,
		Name:                  selections.OS,
		ArchitectureStatusMap: statusMap,
	}, nil
}

func (uc *ConfigurationUseCase) ImportBootImages(ctx context.Context) error {
	return uc.bootResource.Import(ctx)
}

func (uc *ConfigurationUseCase) IsImportingBootImages(ctx context.Context) (bool, error) {
	return uc.bootResource.IsImporting(ctx)
}

func (uc *ConfigurationUseCase) ListBootImageSelections() ([]BootImageSelection, error) {
	selections := []BootImageSelection{}

	for distro := range DistroSeriesMap {
		selections = append(selections, DistroSeriesMap[distro])
	}

	return selections, nil
}

func (uc *ConfigurationUseCase) listBootImages(ctx context.Context) ([]BootImage, error) {
	defaultDistro, err := uc.provisioner.Get(ctx, "default_distro_series")
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

	bootImages := []BootImage{}

	sources, err := uc.bootSource.List(ctx)
	if err != nil {
		return nil, err
	}

	for i := range sources {
		selections, err := uc.bootSourceSelection.List(ctx, sources[i].ID)
		if err != nil {
			return nil, err
		}

		for j := range selections {
			distro := selections[j].Release

			name := selections[j].Release
			if ds, ok := DistroSeriesMap[base.SeriesName(selections[j].Release)]; ok {
				name = ds.DisplayName
			}

			bootImages = append(bootImages, BootImage{
				Source:                sources[i].URL,
				DistroSeries:          distro,
				Name:                  name,
				ID:                    selections[j].ID,
				Architectures:         selections[j].Arches,
				ArchitectureStatusMap: statusMaps[distro],
				Default:               distro == defaultDistro,
			})
		}
	}

	return bootImages, nil
}
