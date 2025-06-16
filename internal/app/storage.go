package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/otterscale/api/storage/v1"
	"github.com/openhdc/otterscale/api/storage/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/core"
)

type StorageService struct {
	pbconnect.UnimplementedStorageServiceHandler

	uc *core.StorageUseCase
}

func NewStorageService(uc *core.StorageUseCase) *StorageService {
	return &StorageService{uc: uc}
}

var _ pbconnect.StorageServiceHandler = (*StorageService)(nil)

func (s *StorageService) ListPools(ctx context.Context, req *connect.Request[pb.ListPoolsRequest]) (*connect.Response[pb.ListPoolsResponse], error) {
	pools, err := s.uc.ListPools(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListPoolsResponse{}
	resp.SetPools(toProtoPools(pools))
	return connect.NewResponse(resp), nil
}

func (s *StorageService) CreatePool(ctx context.Context, req *connect.Request[pb.CreatePoolRequest]) (*connect.Response[pb.Pool], error) {
	pool, err := s.uc.CreatePool(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoPool(pool)
	return connect.NewResponse(resp), nil
}

func (s *StorageService) DeletePool(ctx context.Context, req *connect.Request[pb.DeletePoolRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeletePool(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func toProtoPools(ps []core.Pool) []*pb.Pool {
	ret := []*pb.Pool{}
	for i := range ps {
		ret = append(ret, toProtoPool(&ps[i]))
	}
	return ret
}

func toProtoPool(p *core.Pool) *pb.Pool {
	ret := &pb.Pool{}
	ret.SetName(p.Name)
	return ret
}
