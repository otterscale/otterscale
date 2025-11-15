package app

import (
	"context"
	"net/url"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/otterscale/otterscale/api/environment/v1"
	"github.com/otterscale/otterscale/api/environment/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core/environment"
)

type EnvironmentService struct {
	pbconnect.UnimplementedEnvironmentServiceHandler

	environment *environment.EnvironmentUseCase
}

func NewEnvironmentService(environment *environment.EnvironmentUseCase) *EnvironmentService {
	return &EnvironmentService{
		environment: environment,
	}
}

var _ pbconnect.EnvironmentServiceHandler = (*EnvironmentService)(nil)

func (s *EnvironmentService) CheckHealth(_ context.Context, _ *pb.CheckHealthRequest) (*pb.CheckHealthResponse, error) {
	result, err := s.environment.CheckHealth()
	if err != nil {
		return nil, err
	}

	resp := &pb.CheckHealthResponse{}
	resp.SetResult(pb.CheckHealthResponse_Result(result))
	return resp, nil
}

func (s *EnvironmentService) UpdateConfig(_ context.Context, req *pb.UpdateConfigRequest) (*emptypb.Empty, error) {
	if err := s.environment.UpdateConfig(toConfig(req)); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *EnvironmentService) GetPrometheus(ctx context.Context, _ *pb.GetPrometheusRequest) (*pb.Prometheus, error) {
	_, err := s.environment.DiscoverPrometheusURL(ctx)
	if err != nil {
		return nil, err
	}

	resp := toProtoPrometheus()
	return resp, nil
}

func (s *EnvironmentService) GetPremiumTier(_ context.Context, _ *pb.GetPremiumTierRequest) (*pb.PremiumTier, error) {
	resp := &pb.PremiumTier{}
	return resp, nil
}

func (s *EnvironmentService) GetPrometheusURL() *url.URL {
	return s.environment.GetPrometheusURL()
}

func toConfig(req *pb.UpdateConfigRequest) *config.Config {
	return &config.Config{
		MAAS: config.MAAS{
			URL:     req.GetMaasUrl(),
			Key:     req.GetMaasKey(),
			Version: req.GetMaasVersion(),
		},
		Juju: config.Juju{
			Controller:          req.GetJujuController(),
			ControllerAddresses: req.GetJujuControllerAddresses(),
			Username:            req.GetJujuUsername(),
			Password:            req.GetJujuPassword(),
			CACert:              req.GetJujuCaCert(),
			CloudName:           req.GetJujuCloudName(),
			CloudRegion:         req.GetJujuCloudRegion(),
			CharmhubAPIURL:      req.GetJujuCharmhubApiUrl(),
		},
		MicroK8s: config.MicroK8s{
			Config: req.GetMicroK8SConfig(),
		},
	}
}

func toProtoPrometheus() *pb.Prometheus {
	ret := &pb.Prometheus{}
	ret.SetBaseUrl("/api/v1")
	return ret
}
