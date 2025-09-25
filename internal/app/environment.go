package app

import (
	"context"
	"net/url"
	"sync"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/otterscale/otterscale/api/environment/v1"
	"github.com/otterscale/otterscale/api/environment/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

const statusChanSize = 100

type EnvironmentService struct {
	pbconnect.UnimplementedEnvironmentServiceHandler

	uc         *core.EnvironmentUseCase
	clients    sync.Map
	statusChan chan *pb.WatchStatusResponse
}

func NewEnvironmentService(uc *core.EnvironmentUseCase) *EnvironmentService {
	s := &EnvironmentService{
		uc:         uc,
		statusChan: make(chan *pb.WatchStatusResponse, statusChanSize),
	}
	go s.broadcastStatus()
	return s
}

var _ pbconnect.EnvironmentServiceHandler = (*EnvironmentService)(nil)

func (s *EnvironmentService) CheckHealth(ctx context.Context, _ *pb.CheckHealthRequest) (*pb.CheckHealthResponse, error) {
	result, err := s.uc.CheckHealth(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.CheckHealthResponse{}
	resp.SetResult(pb.CheckHealthResponse_Result(result))
	return resp, nil
}

func (s *EnvironmentService) WatchStatus(ctx context.Context, _ *pb.WatchStatusRequest, stream *connect.ServerStream[pb.WatchStatusResponse]) error {
	// Send initial status to the new client
	status := s.uc.LoadStatus(ctx)
	if err := stream.Send(toProtoWatchStatus(status)); err != nil {
		return err
	}

	// Register client for status updates
	s.clients.Store(stream, struct{}{})
	defer s.clients.Delete(stream)

	// Wait for context cancellation
	<-ctx.Done()
	return ctx.Err()
}

func (s *EnvironmentService) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*emptypb.Empty, error) {
	// Update the environment status in the use case layer
	s.uc.StoreStatus(ctx, req.GetPhase(), req.GetMessage())

	status := s.uc.LoadStatus(ctx)
	select {
	case s.statusChan <- toProtoWatchStatus(status):
	default:
		// Non-blocking send to avoid deadlock if channel is full
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *EnvironmentService) broadcastStatus() {
	for status := range s.statusChan {
		s.clients.Range(func(key, _ any) bool {
			stream := key.(*connect.ServerStream[pb.WatchStatusResponse])
			if err := stream.Send(status); err != nil {
				s.clients.Delete(stream)
			}
			return true
		})
	}
}

func (s *EnvironmentService) UpdateConfig(ctx context.Context, req *pb.UpdateConfigRequest) (*emptypb.Empty, error) {
	if err := s.uc.UpdateConfig(ctx, toConfig(req)); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *EnvironmentService) GetConfigHelmRepositories(_ context.Context, _ *pb.GetConfigHelmRepositoriesRequest) (*pb.GetConfigHelmRepositoriesResponse, error) {
	helmRepos := s.uc.GetConfigHelmRepos()
	resp := &pb.GetConfigHelmRepositoriesResponse{}
	resp.SetUrls(helmRepos)
	return resp, nil
}

func (s *EnvironmentService) UpdateConfigHelmRepositories(_ context.Context, req *pb.UpdateConfigHelmRepositoriesRequest) (*emptypb.Empty, error) {
	if err := s.uc.UpdateConfigHelmRepos(req.GetUrls()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *EnvironmentService) GetPrometheus(ctx context.Context, _ *pb.GetPrometheusRequest) (*pb.Prometheus, error) {
	_, err := s.uc.FetchPrometheusInfo(ctx)
	if err != nil {
		return nil, err
	}
	resp := toProtoPrometheus()
	return resp, nil
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
			Config: req.GetMircoK8SConfig(),
		},
	}
}

func (s *EnvironmentService) GetPrometheusURL() *url.URL {
	if s.uc != nil {
		return s.uc.GetPrometheusURL()
	}
	return nil
}

func toProtoWatchStatus(status *core.EnvironmentStatus) *pb.WatchStatusResponse {
	ret := &pb.WatchStatusResponse{}
	ret.SetStarted(status.Started)
	ret.SetFinished(status.Finished)
	ret.SetPhase(status.Phase)
	ret.SetMessage(status.Message)
	return ret
}

func toProtoPrometheus() *pb.Prometheus {
	ret := &pb.Prometheus{}
	ret.SetBaseUrl("/api/v1")
	return ret
}
