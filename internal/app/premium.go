package app

import (
	"context"

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

func (s *PremiumService) GetTier(_ context.Context, _ *pb.GetTierRequest) (*pb.GetTierResponse, error) {
	resp := &pb.GetTierResponse{}
	return resp, nil
}
