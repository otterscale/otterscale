package juju

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	md "github.com/openhdc/openhdc/internal/domain/model"
	"github.com/openhdc/openhdc/internal/domain/service"
	"github.com/openhdc/openhdc/internal/env"
	"github.com/openhdc/openhdc/internal/utils"
)

const defaultCharmHubAPIURL = "https://api.charmhub.io"

type charmHub struct {
	apiURL          string
	charms          []md.Charm
	charmsCacheTime time.Time
}

func NewCharmHub() service.JujuCharmHub {
	return &charmHub{
		apiURL: env.GetOrDefault(env.OPENHDC_CHARMHUB_API_URL, defaultCharmHubAPIURL),
	}
}

var _ service.JujuCharmHub = (*charmHub)(nil)

func (r *charmHub) List(ctx context.Context) ([]md.Charm, error) {
	if r.charms != nil && time.Since(r.charmsCacheTime) < time.Hour*24 {
		return r.charms, nil
	}
	cs, err := r.find(ctx, "")
	if err != nil {
		return nil, err
	}
	r.charms = cs
	r.charmsCacheTime = time.Now()
	return cs, nil
}

func (r *charmHub) Get(ctx context.Context, name string) (*md.Charm, error) {
	if r.charms != nil && time.Since(r.charmsCacheTime) < time.Hour*24 {
		for i := range r.charms {
			if r.charms[i].Name != name {
				continue
			}
			return &r.charms[i], nil
		}
	}
	cs, err := r.find(ctx, name)
	if err != nil {
		return nil, err
	}
	for i := range cs {
		if cs[i].Name != name {
			continue
		}
		return &cs[i], nil
	}
	return nil, status.Errorf(codes.NotFound, "charm name %q not found", name)
}

func (r *charmHub) ListArtifacts(ctx context.Context, name string) ([]md.CharmArtifact, error) {
	return r.info(ctx, name)
}

func (r *charmHub) find(ctx context.Context, name string) ([]md.Charm, error) {
	queryURL, err := url.ParseRequestURI(r.apiURL)
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
		Results []md.Charm `json:"results"`
	}
	resp := new(response)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (r *charmHub) info(ctx context.Context, name string) ([]md.CharmArtifact, error) {
	queryURL, err := url.ParseRequestURI(r.apiURL)
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

	resp := new(md.Charm)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp.Artifacts, nil
}
