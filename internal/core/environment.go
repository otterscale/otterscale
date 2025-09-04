package core

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"slices"
	"sync"

	"connectrpc.com/connect"

	"github.com/otterscale/otterscale/internal/config"
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

	conf          *config.Config
	statusMap     sync.Map
	prometheusURL *url.URL
}

func NewEnvironmentUseCase(scope ScopeRepo, action ActionRepo, facility FacilityRepo, conf *config.Config) *EnvironmentUseCase {
	return &EnvironmentUseCase{
		scope:         scope,
		action:        action,
		facility:      facility,
		conf:          conf,
		prometheusURL: &url.URL{},
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
	return &EnvironmentStatus{
		Finished: uc.isMAASConfigured(),
	}
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

func (uc *EnvironmentUseCase) GetConfigHelmRepos() []string {
	return uc.conf.Kube.HelmRepositoryURLs
}

func (uc *EnvironmentUseCase) UpdateConfigHelmRepos(urls []string) error {
	uc.conf.Kube.HelmRepositoryURLs = urls
	return uc.conf.Override(uc.conf)
}

func (uc *EnvironmentUseCase) GetPrometheusURL() *url.URL {
	return uc.prometheusURL
}

func (uc *EnvironmentUseCase) FetchPrometheusInfo(ctx context.Context) (*url.URL, error) {
	if uc.prometheusURL.Scheme != "" {
		return uc.prometheusURL, nil
	}

	scopes, err := uc.scope.List(ctx)
	if err != nil {
		return nil, err
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
		url, err := url.Parse(endpoints.Prometheus.URL)
		if err != nil {
			continue
		}
		*uc.prometheusURL = *url
		return uc.prometheusURL, nil
	}
	return nil, connect.NewError(connect.CodeNotFound, errors.New("prometheus info not found"))
}

func (uc *EnvironmentUseCase) isMAASConfigured() bool {
	return uc.conf.MAAS.Key != "::"
}
