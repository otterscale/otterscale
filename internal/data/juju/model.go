package juju

import (
	"context"
	"slices"

	"github.com/juju/juju/api/base"
	api "github.com/juju/juju/api/client/modelmanager"
	"github.com/juju/juju/core/status"
	"github.com/juju/names/v5"

	"github.com/otterscale/otterscale/internal/core"
)

type model struct {
	juju *Juju
}

func NewModel(juju *Juju) core.ScopeRepo {
	return &model{
		juju: juju,
	}
}

var _ core.ScopeRepo = (*model)(nil)

func (r *model) List(_ context.Context) ([]core.Scope, error) {
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

func (r *model) Create(_ context.Context, name string) (*core.Scope, error) {
	conn, err := r.juju.connection("")
	if err != nil {
		return nil, err
	}

	cloudCredential := names.CloudCredentialTag{}
	model, err := api.NewClient(conn).CreateModel(name, r.juju.username(), r.juju.cloudName(), r.juju.cloudRegion(), cloudCredential, nil)
	if err != nil {
		return nil, err
	}
	return &core.Scope{
		UUID:         model.UUID,
		Name:         model.Name,
		Type:         model.Type,
		ProviderType: model.ProviderType,
		Life:         model.Life,
		Status:       model.Status,
		AgentVersion: model.AgentVersion,
		IsController: model.IsController,
	}, nil
}

func (r *model) filterValidModels(models []base.UserModelSummary) []core.Scope {
	return slices.DeleteFunc(models, func(model base.UserModelSummary) bool {
		return model.Name == "controller" || !status.ValidModelStatus(model.Status.Status)
	})
}
