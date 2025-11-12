package maas

import (
	"context"

	entity "github.com/otterscale/otterscale/internal/core/_entity"
	"github.com/otterscale/otterscale/internal/core/configuration"
)

type bootSourceRepo struct {
	maas *MAAS
}

func NewBootSourceRepo(maas *MAAS) configuration.BootSourceRepo {
	return &bootSourceRepo{
		maas: maas,
	}
}

var _ configuration.BootSourceRepo = (*bootSourceRepo)(nil)

func (r *bootSourceRepo) List(_ context.Context) ([]configuration.BootSource, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	sources, err := client.BootSources.Get()
	if err != nil {
		return nil, err
	}

	return r.toBootSources(sources), nil
}

func (r *bootSourceRepo) toBootSources(bss []entity.BootSource) []configuration.BootSource {
	ret := []configuration.BootSource{}
	// for _, bs := range bss {
	// 	ret = append(ret, configuration.BootSource{
	// 		// ID: bs.ID,
	// 	})
	// }
	return ret
}
