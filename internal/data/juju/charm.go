package juju

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"connectrpc.com/connect"
	"github.com/openhdc/otterscale/internal/core"
	"github.com/openhdc/otterscale/internal/utils"
)

type charm struct {
	juju *Juju
}

func NewCharm(juju *Juju) core.CharmRepo {
	return &charm{
		juju: juju,
	}
}

var _ core.CharmRepo = (*charm)(nil)

func (r *charm) List(ctx context.Context) ([]core.Charm, error) {
	return r.find(ctx, "")
}

func (r *charm) Get(ctx context.Context, name string) (*core.Charm, error) {
	charms, err := r.find(ctx, name)
	if err != nil {
		return nil, err
	}
	for i := range charms {
		if charms[i].Name == name {
			return &charms[i], nil
		}
	}
	return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("charm name %q not found", name))
}

func (r *charm) ListArtifacts(ctx context.Context, name string) ([]core.CharmArtifact, error) {
	return r.info(ctx, name)
}

func (r *charm) find(ctx context.Context, name string) ([]core.Charm, error) {
	queryURL, err := url.ParseRequestURI(r.juju.charmhubAPIURL())
	if err != nil {
		return nil, err
	}
	queryURL = queryURL.JoinPath("v2", "charms", "find")

	queryParams := url.Values{}
	queryParams.Set("fields", "default-release,result")
	if name != "" {
		queryParams.Set("q", name)
	}
	queryURL.RawQuery = queryParams.Encode()

	data, err := utils.HTTPGet(ctx, queryURL.String())
	if err != nil {
		return nil, err
	}

	type response struct {
		Results []core.Charm `json:"results"`
	}
	resp := new(response)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (r *charm) info(ctx context.Context, name string) ([]core.CharmArtifact, error) {
	queryURL, err := url.ParseRequestURI(r.juju.charmhubAPIURL())
	if err != nil {
		return nil, err
	}
	queryURL = queryURL.JoinPath("v2", "charms", "info", name)

	queryParams := url.Values{}
	queryParams.Set("fields", "channel-map")
	queryURL.RawQuery = queryParams.Encode()

	data, err := utils.HTTPGet(ctx, queryURL.String())
	if err != nil {
		return nil, err
	}

	resp := new(core.Charm)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp.Artifacts, nil
}
