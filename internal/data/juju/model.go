package juju

import (
	"context"
	"errors"
	"slices"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/cloud"
	"github.com/juju/juju/api/client/modelmanager"
	jujucloud "github.com/juju/juju/cloud"
	jujustatus "github.com/juju/juju/core/status"
	"github.com/juju/names/v5"

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

func (r *model) List(_ context.Context) ([]base.UserModelSummary, error) {
	conn, err := r.jujuMap.Get("")
	if err != nil {
		return nil, err
	}
	ms, err := modelmanager.NewClient(conn).ListModelSummaries(r.user, true)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(ms, func(ums base.UserModelSummary) bool {
		return ums.Name == "controller" || !jujustatus.ValidModelStatus(ums.Status.Status)
	}), nil
}

func (r *model) Create(_ context.Context, name string) (*base.ModelInfo, error) {
	conn, err := r.jujuMap.Get("")
	if err != nil {
		return nil, err
	}

	clouds, err := cloud.NewClient(conn).Clouds()
	if err != nil {
		return nil, err
	}

	cloudName, cloudRegion, err := r.cloudInfo(clouds)
	if err != nil {
		return nil, err
	}

	mi, err := modelmanager.NewClient(conn).CreateModel(
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
		return "", "", status.Error(codes.NotFound, "cloud not found")
	}

	for tag := range clouds {
		if len(clouds[tag].Regions) == 0 {
			return "", "", status.Errorf(codes.NotFound, "cloud %q not found", tag)
		}
		return clouds[tag].Name, clouds[tag].Regions[0].Name, nil
	}

	return "", "", errors.New("unable to determine cloud and region")
}
