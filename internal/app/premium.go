package app

import (
	"context"

	"connectrpc.com/connect"

	pb "github.com/otterscale/otterscale/api/premium/v1"
	"github.com/otterscale/otterscale/api/premium/v1/pbconnect"
)

type PremiumService struct {
	pbconnect.UnimplementedPremiumServiceHandler
}

func NewPremiumService() *PremiumService {
	return &PremiumService{}
}

var _ pbconnect.PremiumServiceHandler = (*PremiumService)(nil)

func (s *PremiumService) GetTier(ctx context.Context, req *connect.Request[pb.GetTierRequest]) (*connect.Response[pb.GetTierResponse], error) {
	resp := &pb.GetTierResponse{}
	return connect.NewResponse(resp), nil
}
