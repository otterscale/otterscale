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

const (
	healthOK           = 11
	healthNotInstalled = 21
)

const traefikShowAction = "show-proxied-endpoints"

type traefikProxiedEndpoints struct {
	Prometheus struct {
		URL string `json:"url"`
	} `json:"prometheus/0"`
}

type EnvironmentStatus struct {
	Started  bool
	Finished bool
	Phase    string
	Message  string
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

func (uc *EnvironmentUseCase) CheckHealth(ctx context.Context) (int32, error) {
	if !uc.isMAASConfigured() {
		return healthNotInstalled, nil
	}
	return healthOK, nil
}

func (uc *EnvironmentUseCase) LoadStatus(ctx context.Context) *EnvironmentStatus {
	v, ok := uc.statusMap.Load("")
	if ok {
		return v.(*EnvironmentStatus)
	}
	return &EnvironmentStatus{}
}

func (uc *EnvironmentUseCase) StoreStatus(ctx context.Context, phase, message string) {
	uc.statusMap.Store("", &EnvironmentStatus{
		Started:  true,
		Finished: uc.isMAASConfigured(),
		Phase:    phase,
		Message:  message,
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

func (uc *EnvironmentUseCase) isMAASConfigured() bool {
	return uc.conf.MAAS.Key != "::"
}
