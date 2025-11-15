package network

import (
	"context"

	"github.com/canonical/gomaasclient/entity"
)

// Fabric represents a MAAS Fabric resource.
type Fabric = entity.Fabric

type FabricRepo interface {
	List(ctx context.Context) ([]Fabric, error)
	Get(ctx context.Context, id int) (*Fabric, error)
	Create(ctx context.Context) (*Fabric, error)
	Update(ctx context.Context, id int, name string) (*Fabric, error)
	Delete(ctx context.Context, id int) error
}

func (uc *UseCase) UpdateFabric(ctx context.Context, id int, name string) (*Fabric, error) {
	return uc.fabric.Update(ctx, id, name)
}
