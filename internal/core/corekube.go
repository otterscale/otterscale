package core

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"connectrpc.com/connect"
	"github.com/juju/juju/api/client/application"
	"k8s.io/client-go/rest"
)

var kubeConfigs sync.Map

func kubeConfig(ctx context.Context, facility FacilityRepo, uuid, name string) (*rest.Config, error) {
	key := uuid + "/" + name

	if v, ok := kubeConfigs.Load(key); ok {
		return v.(*rest.Config), nil
	}

	config, err := newKubeConfig(ctx, facility, uuid, name)
	if err != nil {
		return nil, err
	}

	kubeConfigs.Store(key, config)

	return config, nil
}

func newKubeConfig(ctx context.Context, facility FacilityRepo, uuid, name string) (*rest.Config, error) {
	// kubernetes-control-plane
	leader, err := facility.GetLeader(ctx, uuid, name)
	if err != nil {
		return nil, err
	}
	unitInfo, err := facility.GetUnitInfo(ctx, uuid, leader)
	if err != nil {
		return nil, err
	}

	// kubernetes-worker
	worker, err := extractWorkerUnitName(unitInfo)
	if err != nil {
		return nil, err
	}
	unitInfo, err = facility.GetUnitInfo(ctx, uuid, worker)
	if err != nil {
		return nil, err
	}

	// config
	endpoint, err := extractEndpoint(unitInfo)
	if err != nil {
		return nil, err
	}
	clientToken, err := extractClientToken(unitInfo)
	if err != nil {
		return nil, err
	}
	return &rest.Config{
		Host:        endpoint,
		BearerToken: clientToken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}, nil
}

func extractWorkerUnitName(unitInfo *application.UnitInfo) (string, error) {
	for _, erd := range unitInfo.RelationData {
		if erd.Endpoint != "kube-control" {
			continue
		}
		for name := range erd.UnitRelationData {
			return name, nil
		}
	}
	return "", connect.NewError(connect.CodeNotFound, errors.New("kube-control not found"))
}

func extractEndpoint(unitInfo *application.UnitInfo) (string, error) {
	var endpoints []string
	for _, erd := range unitInfo.RelationData {
		for _, rd := range erd.UnitRelationData {
			if endpointData, ok := rd.UnitData["api-endpoints"]; ok && endpointData != nil {
				if endpointStr, ok := endpointData.(string); ok {
					if err := json.Unmarshal([]byte(endpointStr), &endpoints); err != nil {
						return "", err
					}
				}
			}
		}
	}
	if len(endpoints) > 0 {
		return endpoints[0], nil
	}
	return "", connect.NewError(connect.CodeNotFound, errors.New("endpoint not found"))
}

func extractClientToken(unitInfo *application.UnitInfo) (string, error) {
	credentials := make(map[string]ControlPlaneCredential)
	for _, erd := range unitInfo.RelationData {
		for _, rd := range erd.UnitRelationData {
			if credsData, ok := rd.UnitData["creds"]; ok && credsData != nil {
				if credsStr, ok := credsData.(string); ok {
					if err := json.Unmarshal([]byte(credsStr), &credentials); err != nil {
						return "", err
					}
				}
			}
		}
	}
	for _, cred := range credentials {
		return cred.ClientToken, nil
	}
	return "", errors.New("token not found")
}
