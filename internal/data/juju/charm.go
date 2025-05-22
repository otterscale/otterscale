package juju

import (
	"context"
	"encoding/json"
	"net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	biz "github.com/openhdc/otterscale/internal/domain/model"
	"github.com/openhdc/otterscale/internal/domain/service"
	"github.com/openhdc/otterscale/internal/utils"
)

type charm struct {
	juju *Juju
}

func NewCharm(juju *Juju) service.JujuCharm {
	return &charm{
		juju: juju,
	}
}

var _ service.JujuCharm = (*charm)(nil)

func (r *charm) List(ctx context.Context) ([]biz.Charm, error) {
	return r.find(ctx, "")
}

func (r *charm) Get(ctx context.Context, name string) (*biz.Charm, error) {
	charms, err := r.find(ctx, name)
	if err != nil {
		return nil, err
	}
	for i := range charms {
		if charms[i].Name == name {
			return &charms[i], nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "charm name %q not found", name)
}

func (r *charm) ListArtifacts(ctx context.Context, name string) ([]biz.CharmArtifact, error) {
	return r.info(ctx, name)
}

func (r *charm) find(ctx context.Context, name string) ([]biz.Charm, error) {
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

	data, err := utils.Get(ctx, queryURL.String())
	if err != nil {
		return nil, err
	}

	type response struct {
		Results []biz.Charm `json:"results"`
	}
	resp := new(response)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (r *charm) info(ctx context.Context, name string) ([]biz.CharmArtifact, error) {
	queryURL, err := url.ParseRequestURI(r.juju.charmhubAPIURL())
	if err != nil {
		return nil, err
	}
	queryURL = queryURL.JoinPath("v2", "charms", "info", name)

	queryParams := url.Values{}
	queryParams.Set("fields", "channel-map")
	queryURL.RawQuery = queryParams.Encode()

	data, err := utils.Get(ctx, queryURL.String())
	if err != nil {
		return nil, err
	}

	resp := new(biz.Charm)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp.Artifacts, nil
}
