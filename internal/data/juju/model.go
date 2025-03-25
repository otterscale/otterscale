package juju

import (
	"context"
	"errors"
	"fmt"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/cloud"
	"github.com/juju/juju/api/client/modelmanager"
	"github.com/juju/names/v6"

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

func (r *model) List(ctx context.Context) ([]*base.UserModelSummary, error) {
	client := modelmanager.NewClient(r.juju)
	user := env.GetOrDefault(env.OPENHDC_JUJU_USERNAME, defaultUsername)

	umss, err := client.ListModelSummaries(ctx, user, true)
	if err != nil {
		return nil, err
	}

	ret := make([]*base.UserModelSummary, len(umss))
	for i, ums := range umss {
		ret[i] = &ums
	}

	return ret, nil
}

func (r *model) Create(ctx context.Context, name string) (*base.ModelInfo, error) {
	owner := env.GetOrDefault(env.OPENHDC_JUJU_USERNAME, defaultUsername)

	// Get cloud information
	cloudName, cloudRegion, err := r.getCloudInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Create the model
	mi, err := modelmanager.NewClient(r.juju).CreateModel(
		ctx,
		name,
		owner,
		cloudName,
		cloudRegion,
		names.CloudCredentialTag{},
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &mi, nil
}

// getCloudInfo fetches and returns the first available cloud name and region
func (r *model) getCloudInfo(ctx context.Context) (string, string, error) {
	clouds, err := cloud.NewClient(r.juju).Clouds(ctx)
	if err != nil {
		return "", "", err
	}

	if len(clouds) == 0 {
		return "", "", errors.New("no clouds found")
	}

	// Get the first cloud
	for cloudName, c := range clouds {
		if len(c.Regions) == 0 {
			return "", "", fmt.Errorf("no regions found for cloud: %s", cloudName)
		}

		// Return the first region of the first cloud
		return c.Name, c.Regions[0].Name, nil
	}

	return "", "", errors.New("unable to determine cloud and region")
}
