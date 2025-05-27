package app

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/otterscale/api/environment/v1"
	"github.com/openhdc/otterscale/api/environment/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/config"
	"github.com/openhdc/otterscale/internal/core"
)

type EnvironmentService struct {
	pbconnect.UnimplementedEnvironmentServiceHandler

	uc *core.EnvironmentUseCase
}

func NewEnvironmentService(uc *core.EnvironmentUseCase) *EnvironmentService {
	return &EnvironmentService{uc: uc}
}

var _ pbconnect.EnvironmentServiceHandler = (*EnvironmentService)(nil)

func (s *EnvironmentService) CheckHealthy(ctx context.Context, req *connect.Request[pb.CheckHealthyRequest]) (*connect.Response[pb.CheckHealthyResponse], error) {
	result, err := s.uc.CheckHealthy(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.CheckHealthyResponse{}
	resp.SetResult(pb.CheckHealthyResponse_Result(result))
	return connect.NewResponse(resp), nil
}

func (s *EnvironmentService) WatchStatuses(ctx context.Context, req *connect.Request[pb.WatchStatusesRequest], stream *connect.ServerStream[pb.WatchStatusesResponse]) error {
	ticker := time.NewTicker(10 * time.Second) //nolint:mnd
	defer ticker.Stop()

	// Send initial status immediately
	if err := s.sendStatus(ctx, stream); err != nil {
		return err
	}

	// Then send status every 10 seconds
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.sendStatus(ctx, stream); err != nil {
				return err
			}
		}
	}
}

func (s *EnvironmentService) UpdateStatus(ctx context.Context, req *connect.Request[pb.UpdateStatusRequest]) (*connect.Response[emptypb.Empty], error) {
	s.uc.StoreStatus(ctx, req.Msg.GetPhase(), req.Msg.GetMessage())
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *EnvironmentService) UpdateConfig(ctx context.Context, req *connect.Request[pb.UpdateConfigRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.UpdateConfig(ctx, toConfigSet(req.Msg)); err != nil {
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

func (s *EnvironmentService) sendStatus(ctx context.Context, stream *connect.ServerStream[pb.WatchStatusesResponse]) error {
	status, err := s.uc.LoadStatus(ctx)
	if err != nil {
		return err
	}
	resp := &pb.WatchStatusesResponse{}
	resp.SetPhase(status.Phase)
	resp.SetMessage(status.Message)
	return stream.Send(resp)
}

func toConfigSet(req *pb.UpdateConfigRequest) *config.ConfigSet {
	maas := &config.MAAS{}
	maas.SetUrl(req.GetMaasUrl())
	maas.SetKey(req.GetMaasKey())
	maas.SetVersion(req.GetMaasVersion())

	juju := &config.Juju{}
	juju.SetControllerAddresses(req.GetJujuControllerAddresses())
	juju.SetUsername(req.GetJujuUsername())
	juju.SetPassword(req.GetJujuPassword())
	juju.SetCaCert(req.GetJujuCaCert())
	juju.SetCloudName(req.GetJujuCloudName())
	juju.SetCloudRegion(req.GetJujuCloudRegion())

	set := &config.ConfigSet{}
	set.SetMaas(maas)
	set.SetJuju(juju)

	return set
}
