package app

import (
	"context"
	"sync"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/otterscale/otterscale/api/bootstrap/v1"
	"github.com/otterscale/otterscale/api/bootstrap/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type BootstrapService struct {
	pbconnect.UnimplementedBootstrapServiceHandler

	bootstrap *core.BootstrapUseCase

	statusChan chan *pb.WatchStatusResponse
	clients    sync.Map
}

func NewBootstrapService(bootstrap *core.BootstrapUseCase) *BootstrapService {
	size := 100
	svc := &BootstrapService{
		bootstrap:  bootstrap,
		statusChan: make(chan *pb.WatchStatusResponse, size),
	}
	go svc.broadcastStatus()
	return svc
}

var _ pbconnect.BootstrapServiceHandler = (*BootstrapService)(nil)

func (s *BootstrapService) WatchStatus(ctx context.Context, _ *pb.WatchStatusRequest, stream *connect.ServerStream[pb.WatchStatusResponse]) error {
	// Send initial status to the new client
	status := s.bootstrap.LoadStatus()
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

func (s *BootstrapService) UpdateStatus(_ context.Context, req *pb.UpdateStatusRequest) (*emptypb.Empty, error) {
	// Update the environment status in the use case layer
	s.bootstrap.StoreStatus(req.GetPhase(), req.GetMessage())

	status := s.bootstrap.LoadStatus()
	select {
	case s.statusChan <- toProtoWatchStatus(status):
	default:
		// Non-blocking send to avoid deadlock if channel is full
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *BootstrapService) broadcastStatus() {
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

func toProtoWatchStatus(status *core.BootstrapStatus) *pb.WatchStatusResponse {
	ret := &pb.WatchStatusResponse{}
	ret.SetStarted(status.Started)
	ret.SetFinished(status.Finished)
	ret.SetPhase(status.Phase)
	ret.SetMessage(status.Message)
	return ret
}
