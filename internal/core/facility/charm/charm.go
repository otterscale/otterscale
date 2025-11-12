package charm

import (
	"context"
	"slices"
)

type Charm struct {
	Name         string
	Type         string
	DeployableOn []string

	Artifacts []CharmArtifact
}

type CharmArtifact struct{}

type CharmRepo interface {
	List(ctx context.Context) ([]Charm, error)
	Get(ctx context.Context, name string) (*Charm, error)
	ListArtifacts(ctx context.Context, name string) ([]CharmArtifact, error)
}

type CharmUseCase struct {
	charm CharmRepo
}

func NewCharmUseCase(charm CharmRepo) *CharmUseCase {
	return &CharmUseCase{
		charm: charm,
	}
}

func (uc *CharmUseCase) ListCharms(ctx context.Context) ([]Charm, error) {
	charms, err := uc.charm.List(ctx)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(charms, func(charm Charm) bool {
		return slices.Contains(charm.DeployableOn, "kubernetes") || charm.Type != "charm"
	}), nil
}

func (uc *CharmUseCase) GetCharm(ctx context.Context, name string) (*Charm, error) {
	return uc.charm.Get(ctx, name)
}

func (uc *CharmUseCase) ListArtifacts(ctx context.Context, name string) ([]CharmArtifact, error) {
	return uc.charm.ListArtifacts(ctx, name)
}
