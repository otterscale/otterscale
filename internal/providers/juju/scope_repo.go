package juju

import (
	"context"
	"fmt"
	"slices"

	"connectrpc.com/connect"
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/keymanager"
	"github.com/juju/juju/api/client/modelmanager"
	"github.com/juju/juju/core/status"
	"github.com/juju/names/v5"

	"github.com/otterscale/otterscale/internal/core/scope/scope"
)

type scopeRepo struct {
	juju *Juju
}

func NewScopeRepo(juju *Juju) scope.ScopeRepo {
	return &scopeRepo{
		juju: juju,
	}
}

var _ scope.ScopeRepo = (*scopeRepo)(nil)

func (r *scopeRepo) List(_ context.Context) ([]scope.Scope, error) {
	conn, err := r.juju.Connection("controller")
	if err != nil {
		return nil, err
	}

	summaries, err := modelmanager.NewClient(conn).ListModelSummaries(r.juju.username(), true)
	if err != nil {
		return nil, err
	}

	summaries = slices.DeleteFunc(summaries, func(model base.UserModelSummary) bool {
		return model.Name == "controller" || !status.ValidModelStatus(model.Status.Status)
	})

	return r.toScopesFromSummary(summaries), nil
}

func (r *scopeRepo) Get(ctx context.Context, name string) (*scope.Scope, error) {
	scopes, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	for i := range scopes {
		if scopes[i].Name == name {
			return &scopes[i], nil
		}
	}

	return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("scope %q not found", name))
}

func (r *scopeRepo) Create(ctx context.Context, name, sshKey string) (*scope.Scope, error) {
	info, err := r.create(ctx, name, sshKey)
	if err != nil {
		return nil, err
	}

	if err := r.addSSHKey(ctx, name, sshKey); err != nil {
		return nil, err
	}

	return r.toScopeFromInfo(info), nil
}

func (r *scopeRepo) create(_ context.Context, name, sshKey string) (base.ModelInfo, error) {
	conn, err := r.juju.Connection("controller")
	if err != nil {
		return base.ModelInfo{}, err
	}
	return modelmanager.NewClient(conn).CreateModel(name, r.juju.username(), r.juju.cloudName(), r.juju.cloudRegion(), names.CloudCredentialTag{}, nil)
}

func (r *scopeRepo) addSSHKey(_ context.Context, name, sshKey string) error {
	conn, err := r.juju.Connection(name)
	if err != nil {
		return err
	}

	_, err = keymanager.NewClient(conn).AddKeys(r.juju.username(), sshKey)
	return err
}

func (r *scopeRepo) toScopesFromSummary(ms []base.UserModelSummary) []scope.Scope {
	scopes := make([]scope.Scope, len(ms))
	for i, m := range ms {
		scopes[i] = r.toScopeFromSummary(m)
	}
	return scopes
}

func (r *scopeRepo) toScopeFromSummary(m base.UserModelSummary) scope.Scope {
	return scope.Scope{
		UUID: m.UUID,
		Name: m.Name,
	}
}

func (r *scopeRepo) toScopeFromInfo(m base.ModelInfo) *scope.Scope {
	return &scope.Scope{
		UUID: m.UUID,
		Name: m.Name,
	}
}
