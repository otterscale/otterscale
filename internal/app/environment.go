package app

import (
	"context"
	"sync"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/otterscale/api/environment/v1"
	"github.com/openhdc/otterscale/api/environment/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/config"
	"github.com/openhdc/otterscale/internal/core"
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

func (s *EnvironmentService) CheckHealth(ctx context.Context, req *connect.Request[pb.CheckHealthRequest]) (*connect.Response[pb.CheckHealthResponse], error) {
	result, err := s.uc.CheckHealth(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.CheckHealthResponse{}
	resp.SetResult(pb.CheckHealthResponse_Result(result))
	return connect.NewResponse(resp), nil
}

func (s *EnvironmentService) WatchStatus(ctx context.Context, req *connect.Request[pb.WatchStatusRequest], stream *connect.ServerStream[pb.WatchStatusResponse]) error {
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

func (s *EnvironmentService) UpdateStatus(ctx context.Context, req *connect.Request[pb.UpdateStatusRequest]) (*connect.Response[emptypb.Empty], error) {
	// Update the environment status in the use case layer
	s.uc.StoreStatus(ctx, req.Msg.GetPhase(), req.Msg.GetMessage())

	status := s.uc.LoadStatus(ctx)
	select {
	case s.statusChan <- toProtoWatchStatus(status):
	default:
		// Non-blocking send to avoid deadlock if channel is full
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
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

func (s *EnvironmentService) UpdateConfig(ctx context.Context, req *connect.Request[pb.UpdateConfigRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.UpdateConfig(ctx, toConfig(req.Msg)); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *EnvironmentService) UpdateConfigHelmRepositories(ctx context.Context, req *connect.Request[pb.UpdateConfigHelmRepositoriesRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.UpdateConfigHelmRepos(ctx, req.Msg.GetUrls()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *EnvironmentService) GetPrometheus(ctx context.Context, req *connect.Request[pb.GetPrometheusRequest]) (*connect.Response[pb.Prometheus], error) {
	endpoint, baseURL, err := s.uc.GetPrometheusInfo(ctx)
	if err != nil {
		return nil, err
	}
	resp := toProtoPrometheus(endpoint, baseURL)
	return connect.NewResponse(resp), nil
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
			Host:  req.GetMicroK8SHost(),
			Token: req.GetMicroK8SToken(),
		},
	}
}

func toProtoWatchStatus(status *core.EnvironmentStatus) *pb.WatchStatusResponse {
	ret := &pb.WatchStatusResponse{}
	ret.SetStarted(status.Started)
	ret.SetFinished(status.Finished)
	ret.SetPhase(status.Phase)
	ret.SetMessage(status.Message)
	return ret
}

func toProtoPrometheus(endpoint, baseURL string) *pb.Prometheus {
	ret := &pb.Prometheus{}
	ret.SetEndpoint(endpoint)
	ret.SetBaseUrl(baseURL)
	return ret
}
