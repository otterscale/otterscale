package core

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"sync"

	"connectrpc.com/connect"

	"github.com/openhdc/otterscale/internal/config"
)

const traefikShowAction = "show-proxied-endpoints"

type traefikProxiedEndpoints struct {
	Prometheus struct {
		URL string `json:"url"`
	} `json:"prometheus/0"`
}

type EnvironmentStatus struct {
	Phase   string
	Message string
}

type EnvironmentUseCase struct {
	scope    ScopeRepo
	action   ActionRepo
	facility FacilityRepo

	prometheusEndpoint string
	prometheusBaseURL  string

	conf      *config.Config
	statusMap sync.Map
}

func NewEnvironmentUseCase(scope ScopeRepo, action ActionRepo, facility FacilityRepo, conf *config.Config) *EnvironmentUseCase {
	return &EnvironmentUseCase{
		scope:             scope,
		action:            action,
		facility:          facility,
		conf:              conf,
		prometheusBaseURL: "/api/v1",
	}
}

//nolint:mnd
func (uc *EnvironmentUseCase) CheckHealth(ctx context.Context) (int32, error) {
	if uc.conf.MAAS.Key == "::" {
		return 21, nil // NOT_INSTALLED
	}
	return 11, nil // OK
}

func (uc *EnvironmentUseCase) LoadStatus(ctx context.Context) (*EnvironmentStatus, error) {
	v, ok := uc.statusMap.Load("")
	if ok {
		return v.(*EnvironmentStatus), nil
	}
	return nil, connect.NewError(connect.CodeNotFound, errors.New("status not found"))
}

func (uc *EnvironmentUseCase) StoreStatus(ctx context.Context, phase, message string) {
	uc.statusMap.Store("", &EnvironmentStatus{
		Phase:   phase,
		Message: message,
	})
}

func (uc *EnvironmentUseCase) UpdateConfig(ctx context.Context, conf *config.Config) error {
	uc.conf.MAAS = conf.MAAS
	uc.conf.Juju = conf.Juju
	uc.conf.MicroK8s = conf.MicroK8s
	return uc.conf.Override(uc.conf)
}

func (uc *EnvironmentUseCase) UpdateConfigHelmRepos(ctx context.Context, urls []string) error {
	uc.conf.Kube.HelmRepositoryURLs = urls
	return uc.conf.Override(uc.conf)
}

func (uc *EnvironmentUseCase) GetPrometheusInfo(ctx context.Context) (endpoint, baseURL string, err error) {
	if uc.prometheusEndpoint != "" {
		return uc.prometheusEndpoint, uc.prometheusBaseURL, nil
	}

	scopes, err := uc.scope.List(ctx)
	if err != nil {
		return "", "", err
	}

	cosScopes := []string{"cos", "cos-lite", "cos-dev"}
	scopes = slices.DeleteFunc(scopes, func(s Scope) bool {
		return !slices.Contains(cosScopes, s.Name)
	})

	for i := range scopes {
		leader, err := uc.facility.GetLeader(ctx, scopes[i].UUID, "traefik")
		if err != nil {
			continue
		}
		result, err := runAction(ctx, uc.action, scopes[i].UUID, leader, traefikShowAction, nil)
		if err != nil {
			continue
		}
		var endpoints traefikProxiedEndpoints
		if err := json.Unmarshal([]byte(result.Output["proxied-endpoints"].(string)), &endpoints); err != nil {
			continue
		}
		uc.prometheusEndpoint = endpoints.Prometheus.URL
		return uc.prometheusEndpoint, uc.prometheusBaseURL, nil
	}
	return "", "", connect.NewError(connect.CodeNotFound, errors.New("prometheus info not found"))
}
