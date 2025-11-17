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

	"github.com/otterscale/otterscale/internal/core/scope"
)

// Note: Juju API do not support context.
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
	conn, err := r.juju.connection("controller")
	if err != nil {
		return nil, err
	}

	summaries, err := modelmanager.NewClient(conn).ListModelSummaries(r.juju.conf.JujuUsername(), true)
	if err != nil {
		return nil, err
	}

	return slices.DeleteFunc(summaries, func(model base.UserModelSummary) bool {
		return model.Name == "controller" || !status.ValidModelStatus(model.Status.Status)
	}), nil
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

func (r *scopeRepo) Create(ctx context.Context, name, aptMirrorURL, sshKey string) (*scope.Scope, error) {
	info, err := r.create(ctx, name, aptMirrorURL)
	if err != nil {
		return nil, err
	}

	if err := r.addSSHKey(ctx, name, sshKey); err != nil {
		return nil, err
	}

	return r.toScopeFromInfo(&info), nil
}

func (r *scopeRepo) create(_ context.Context, name, aptMirrorURL string) (base.ModelInfo, error) {
	conn, err := r.juju.connection("controller")
	if err != nil {
		return base.ModelInfo{}, err
	}

	userName := r.juju.conf.JujuUsername()
	cloudName := r.juju.conf.JujuCloudName()
	cloudRegion := r.juju.conf.JujuCloudRegion()
	cloudCredential := names.CloudCredentialTag{}
	config := map[string]any{"apt-mirror": aptMirrorURL}

	return modelmanager.NewClient(conn).CreateModel(name, userName, cloudName, cloudRegion, cloudCredential, config)
}

func (r *scopeRepo) addSSHKey(_ context.Context, name, sshKey string) error {
	conn, err := r.juju.connection(name)
	if err != nil {
		return err
	}

	_, err = keymanager.NewClient(conn).AddKeys(r.juju.conf.JujuUsername(), sshKey)
	return err
}

func (r *scopeRepo) toScopeFromInfo(m *base.ModelInfo) *scope.Scope {
	return &scope.Scope{
		UUID:         m.UUID,
		Name:         m.Name,
		Type:         m.Type,
		ProviderType: m.ProviderType,
		Life:         m.Life,
		Status:       m.Status,
		AgentVersion: m.AgentVersion,
		IsController: m.IsController,
	}
}
