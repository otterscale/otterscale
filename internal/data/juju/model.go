package juju

import (
	"context"
	"errors"
	"fmt"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/cloud"
	"github.com/juju/juju/api/client/modelmanager"
	jujucloud "github.com/juju/juju/cloud"
	"github.com/juju/names/v6"

	"github.com/openhdc/openhdc/internal/domain/service"
	"github.com/openhdc/openhdc/internal/env"
)

type model struct {
	jujuMap JujuMap
	user    string
}

func NewModel(jujuMap JujuMap) service.JujuModel {
	return &model{
		jujuMap: jujuMap,
		user:    env.GetOrDefault(env.OPENHDC_JUJU_USERNAME, defaultUsername),
	}
}

var _ service.JujuModel = (*model)(nil)

func (r *model) List(ctx context.Context) ([]base.UserModelSummary, error) {
	conn, err := r.jujuMap.Get(ctx, "")
	if err != nil {
		return nil, err
	}
	return modelmanager.NewClient(conn).ListModelSummaries(ctx, r.user, true)
}

func (r *model) Create(ctx context.Context, name string) (*base.ModelInfo, error) {
	conn, err := r.jujuMap.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	clouds, err := cloud.NewClient(conn).Clouds(ctx)
	if err != nil {
		return nil, err
	}

	cloudName, cloudRegion, err := r.cloudInfo(clouds)
	if err != nil {
		return nil, err
	}

	mi, err := modelmanager.NewClient(conn).CreateModel(
		ctx,
		name,
		r.user,
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

func (r *model) cloudInfo(clouds map[names.CloudTag]jujucloud.Cloud) (name, region string, err error) {
	if len(clouds) == 0 {
		return "", "", errors.New("cloud not found")
	}

	for tag := range clouds {
		if len(clouds[tag].Regions) == 0 {
			return "", "", fmt.Errorf("cloud %q region not found ", tag)
		}
		return clouds[tag].Name, clouds[tag].Regions[0].Name, nil
	}

	return "", "", errors.New("unable to determine cloud and region")
}
