package juju

import (
	"context"
	"slices"

	"github.com/juju/juju/api/base"
	api "github.com/juju/juju/api/client/modelmanager"
	"github.com/juju/juju/core/status"
	"github.com/juju/names/v5"

	"github.com/openhdc/otterscale/internal/domain/service"
)

type model struct {
	juju *Juju
}

func NewModel(juju *Juju) service.JujuModel {
	return &model{
		juju: juju,
	}
}

var _ service.JujuModel = (*model)(nil)

func (r *model) List(_ context.Context) ([]base.UserModelSummary, error) {
	conn, err := r.juju.connection("")
	if err != nil {
		return nil, err
	}

	models, err := api.NewClient(conn).ListModelSummaries(r.juju.username(), true)
	if err != nil {
		return nil, err
	}
	return r.filterValidModels(models), nil
}

func (r *model) Create(_ context.Context, name string) (*base.ModelInfo, error) {
	conn, err := r.juju.connection("")
	if err != nil {
		return nil, err
	}

	cloudCredential := names.CloudCredentialTag{}
	model, err := api.NewClient(conn).CreateModel(name, r.juju.username(), r.juju.cloudName(), r.juju.cloudRegion(), cloudCredential, nil)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *model) filterValidModels(models []base.UserModelSummary) []base.UserModelSummary {
	return slices.DeleteFunc(models, func(model base.UserModelSummary) bool {
		return model.Name == "controller" || !status.ValidModelStatus(model.Status.Status)
	})
}
