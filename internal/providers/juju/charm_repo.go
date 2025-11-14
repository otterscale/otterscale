package juju

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"connectrpc.com/connect"

	"github.com/otterscale/otterscale/internal/core/facility/charm"
)

type charmRepo struct {
	juju *Juju
}

func NewCharmRepo(juju *Juju) charm.CharmRepo {
	return &charmRepo{
		juju: juju,
	}
}

var _ charm.CharmRepo = (*charmRepo)(nil)

func (r *charmRepo) List(ctx context.Context) ([]charm.Charm, error) {
	return r.find(ctx, "")
}

func (r *charmRepo) Get(ctx context.Context, name string) (*charm.Charm, error) {
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

func (r *charmRepo) ListArtifacts(ctx context.Context, name string) ([]charm.CharmArtifact, error) {
	return r.info(ctx, name)
}

func (r *charmRepo) find(ctx context.Context, name string) ([]charm.Charm, error) {
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

	data, err := r.httpGet(ctx, queryURL.String())
	if err != nil {
		return nil, err
	}

	type response struct {
		Results []charm.Charm `json:"results"`
	}

	resp := new(response)

	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp.Results, nil
}

func (r *charmRepo) info(ctx context.Context, name string) ([]charm.CharmArtifact, error) {
	queryURL, err := url.ParseRequestURI(r.juju.charmhubAPIURL())
	if err != nil {
		return nil, err
	}
	queryURL = queryURL.JoinPath("v2", "charms", "info", name)

	queryParams := url.Values{}
	queryParams.Set("fields", "channel-map")
	queryURL.RawQuery = queryParams.Encode()

	data, err := r.httpGet(ctx, queryURL.String())
	if err != nil {
		return nil, err
	}

	resp := new(charm.Charm)

	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp.Artifacts, nil
}

func (r *charmRepo) httpGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("get %q failed: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get %q failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("get %q failed with code %d", url, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
