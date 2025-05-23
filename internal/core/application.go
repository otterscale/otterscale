package core

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"k8s.io/client-go/rest"
)

type ApplicationUseCase struct {
	kubeApps    KubeAppsRepo
	kubeCore    KubeCoreRepo
	kubeStorage KubeStorageRepo
	chart       ChartRepo
	release     ReleaseRepo
	facility    FacilityRepo
	scope       ScopeRepo
	client      ClientRepo

	configs sync.Map
}

func NewApplicationUseCase(kubeApps KubeAppsRepo, kubeCore KubeCoreRepo, kubeStorage KubeStorageRepo, chart ChartRepo, release ReleaseRepo, facility FacilityRepo, scope ScopeRepo, client ClientRepo) *ApplicationUseCase {
	return &ApplicationUseCase{
		kubeApps:    kubeApps,
		kubeCore:    kubeCore,
		kubeStorage: kubeStorage,
		chart:       chart,
		release:     release,
		facility:    facility,
		scope:       scope,
		client:      client,
	}
}

func (uc *ApplicationUseCase) GetPublicAddress(ctx context.Context, uuid, facility string) (string, error) {
	leader, err := uc.facility.GetLeader(ctx, uuid, facility)
	if err != nil {
		return "", err
	}
	unitInfo, err := uc.facility.GetUnitInfo(ctx, uuid, leader)
	if err != nil {
		return "", err
	}
	return unitInfo.PublicAddress, nil
}

func (uc *ApplicationUseCase) config(ctx context.Context, uuid, name string) (*rest.Config, error) {
	key := uuid + "/" + name

	if v, ok := uc.configs.Load(key); ok {
		return v.(*rest.Config), nil
	}

	config, err := uc.newConfig(ctx, uuid, name)
	if err != nil {
		return nil, err
	}

	uc.configs.Store(uuid, config)

	return config, nil
}

func (uc *ApplicationUseCase) newConfig(ctx context.Context, uuid, name string) (*rest.Config, error) {
	// kubernetes-control-plane
	leader, err := uc.facility.GetLeader(ctx, uuid, name)
	if err != nil {
		return nil, err
	}
	unitInfo, err := uc.facility.GetUnitInfo(ctx, uuid, leader)
	if err != nil {
		return nil, err
	}
	kubeControl, err := extractWorkerUnitName(unitInfo)
	if err != nil {
		return nil, err
	}

	// kubernetes-worker
	leader, err = uc.facility.GetLeader(ctx, uuid, kubeControl)
	if err != nil {
		return nil, err
	}
	unitInfo, err = uc.facility.GetUnitInfo(ctx, uuid, leader)
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

func extractWorkerUnitName(unitInfo *UnitInfo) (string, error) {
	for _, erd := range unitInfo.RelationData {
		if erd.Endpoint != "kube-control" {
			continue
		}
		for name := range erd.UnitRelationData {
			return name, nil
		}
	}
	return "", status.Error(codes.NotFound, "kube-control not found")
}

func extractEndpoint(unitInfo *UnitInfo) (string, error) {
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
	return "", status.Error(codes.NotFound, "endpoint not found")
}

func extractClientToken(unitInfo *UnitInfo) (string, error) {
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
