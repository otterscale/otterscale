package facility

import (
	"context"
)

type Facility struct {
	Name      string
	Charm     string
	UnitNames []string
}

type FacilityRepo interface {
	List(ctx context.Context, scope, jujuID string) ([]Facility, error)
	Create(ctx context.Context, scope, name, configYAML, charmName, channel string, revision int, series, directive, placementScope string) error
	Update(ctx context.Context, scope, name, configYAML string) error
	Delete(ctx context.Context, scope, name string, destroyStorage, force bool) error
	Resolve(ctx context.Context, scope, unitName string) error
	Config(ctx context.Context, scope string, name string) (map[string]any, error)
	PublicAddress(ctx context.Context, scope, name string) (string, error)
}

type RelationRepo interface {
	Create(ctx context.Context, scope string, endpoints []string) error
	Delete(ctx context.Context, scope string, id int) error
	// Consume(ctx context.Context, scope string, args *crossmodel.ConsumeApplicationArgs) error
	// ConsumeInfo(ctx context.Context, url string) (params.ConsumeOfferDetails, error)
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
