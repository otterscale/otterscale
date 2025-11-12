package network

import (
	"context"
)

type Fabric struct {
	ID    int
	VLANs []VLAN
}

type FabricRepo interface {
	List(ctx context.Context) ([]Fabric, error)
	Get(ctx context.Context, id int) (*Fabric, error)
	Create(ctx context.Context) (*Fabric, error)
	Update(ctx context.Context, id int, name string) (*Fabric, error)
	Delete(ctx context.Context, id int) error
}

func (uc *NetworkUseCase) UpdateFabric(ctx context.Context, id int, name string) (*Fabric, error) {
	return uc.fabric.Update(ctx, id, name)
}
