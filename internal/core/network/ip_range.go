package network

import (
	"context"

	entity "github.com/otterscale/otterscale/internal/core/_entity"
)

// IPRange represents a MAAS IPRange resource.
type IPRange = entity.IPRange

type IPRangeRepo interface {
	List(ctx context.Context) ([]IPRange, error)
	Create(ctx context.Context, subnetID int, startIP, endIP, comment string) (*IPRange, error)
	Update(ctx context.Context, id int, startIP, endIP, comment string) (*IPRange, error)
	Delete(ctx context.Context, id int) error
}

func (uc *NetworkUseCase) CreateIPRange(ctx context.Context, subnetID int, startIP, endIP, comment string) (*IPRange, error) {
	return uc.ipRange.Create(ctx, subnetID, startIP, endIP, comment)
}

func (uc *NetworkUseCase) UpdateIPRange(ctx context.Context, id int, startIP, endIP, comment string) (*IPRange, error) {
	return uc.ipRange.Update(ctx, id, startIP, endIP, comment)
}

func (uc *NetworkUseCase) DeleteIPRange(ctx context.Context, id int) error {
	return uc.ipRange.Delete(ctx, id)
}
