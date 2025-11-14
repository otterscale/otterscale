package facility

import (
	"context"

	"github.com/juju/juju/rpc/params"
)

type Facility struct {
	Name   string
	Status *params.ApplicationStatus
}

type FacilityRepo interface {
	List(ctx context.Context, scope, jujuID string) ([]Facility, error)
	Create(ctx context.Context, scope, name, configYAML, charmName, channel, placementScope string, subordinate bool, directive, series string) error
	Update(ctx context.Context, scope, name, configYAML string) error
	Delete(ctx context.Context, scope, name string, destroyStorage, force bool) error
	Resolve(ctx context.Context, scope, unitName string) error
	Config(ctx context.Context, scope string, name string) (map[string]any, error)
}

type RelationRepo interface {
	Create(ctx context.Context, scope string, endpoints []string) error
	Delete(ctx context.Context, scope string, id int) error
	Consume(ctx context.Context, scope, url string) error
}

type FacilityUseCase struct {
	facility FacilityRepo
}

func NewFacilityUseCase(facility FacilityRepo) *FacilityUseCase {
	return &FacilityUseCase{
		facility: facility,
	}
}

func (uc *FacilityUseCase) ListFacilities(ctx context.Context, scope string) ([]Facility, error) {
	return uc.facility.List(ctx, scope, "")
}

func (uc *FacilityUseCase) ResolveFacilityUnitErrors(ctx context.Context, scope, unitName string) error {
	return uc.facility.Resolve(ctx, scope, unitName)
}
