package core

import (
	"context"
	"errors"
	"regexp"
	"slices"
	"sync"

	"connectrpc.com/connect"
	"github.com/juju/juju/api/client/application"

	"github.com/openhdc/otterscale/internal/config"
)

type EnvironmentStatus struct {
	Phase   string
	Message string
}

type EnvironmentUseCase struct {
	scope    ScopeRepo
	facility FacilityRepo

	conf      *config.Config
	statusMap sync.Map
}

func NewEnvironmentUseCase(scope ScopeRepo, facility FacilityRepo, conf *config.Config) *EnvironmentUseCase {
	return &EnvironmentUseCase{
		scope:    scope,
		facility: facility,
		conf:     conf,
	}
}

//nolint:mnd
func (uc *EnvironmentUseCase) CheckHealthy(ctx context.Context) (int32, error) {
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
	return uc.conf.Override(uc.conf)
}

func (uc *EnvironmentUseCase) UpdateConfigHelmRepos(ctx context.Context, urls []string) error {
	uc.conf.Kube.HelmRepositoryURLs = urls
	return uc.conf.Override(uc.conf)
}

func (uc *EnvironmentUseCase) GetPrometheusInfo(ctx context.Context) (endpoint, baseURL string, err error) {
	scopes, err := uc.scope.List(ctx)
	if err != nil {
		return "", "", err
	}

	cosScopes := []string{"cos", "cos-lite", "cos-dev"}
	scopes = slices.DeleteFunc(scopes, func(s Scope) bool {
		return !slices.Contains(cosScopes, s.Name)
	})

	for i := range scopes {
		leader, err := uc.facility.GetLeader(ctx, scopes[i].UUID, "prometheus")
		if err != nil {
			return "", "", err
		}
		unitInfo, err := uc.facility.GetUnitInfo(ctx, scopes[i].UUID, leader)
		if err != nil {
			return "", "", err
		}
		endpoint, _ := extractPrometheusIngress(unitInfo)
		if endpoint != "" {
			return endpoint, "/api/v1", nil
		}
	}
	return "", "", connect.NewError(connect.CodeNotFound, errors.New("prometheus info not found"))
}

func extractPrometheusIngress(unitInfo *application.UnitInfo) (string, error) {
	for _, erd := range unitInfo.RelationData {
		ingress, ok := erd.ApplicationData["ingress"]
		if !ok || ingress == nil {
			continue
		}
		val, ok := ingress.(string)
		if !ok {
			return "", errors.New("prometheus ingress value is not a string")
		}
		re := regexp.MustCompile(`url: (.*)`)
		matches := re.FindStringSubmatch(val)
		if len(matches) < 2 {
			return "", errors.New("prometheus ingress url not found")
		}
		return matches[1], nil
	}
	return "", errors.New("prometheus ingress not found")
}
