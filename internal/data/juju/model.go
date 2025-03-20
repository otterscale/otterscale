package juju

import (
	"context"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/modelmanager"
	"github.com/openhdc/openhdc/internal/domain/service"
	"github.com/openhdc/openhdc/internal/env"
)

type model struct {
	juju Juju
}

func NewModel(juju Juju) service.JujuModel {
	return &model{
		juju: juju,
	}
}

var _ service.JujuModel = (*model)(nil)

func (r *model) List(ctx context.Context) ([]base.UserModel, error) {
	user := env.GetOrDefault(env.OPENHDC_JUJU_USERNAME, defaultUsername)
	return modelmanager.NewClient(r.juju).ListModels(ctx, user)
}
