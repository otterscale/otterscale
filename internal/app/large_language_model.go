package app

import (
	"context"

	pb "github.com/otterscale/otterscale/api/large_language_model/v1"
	"github.com/otterscale/otterscale/api/large_language_model/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type LargeLanguageModelService struct {
	pbconnect.UnimplementedLargeLanguageModelServiceHandler

	uc *core.LargeLanguageModelUseCase
}

func NewLargeLanguageModelService(uc *core.LargeLanguageModelUseCase) *LargeLanguageModelService {
	return &LargeLanguageModelService{uc: uc}
}

var _ pbconnect.LargeLanguageModelServiceHandler = (*LargeLanguageModelService)(nil)

func (s *LargeLanguageModelService) CheckInfrastructureStatus(ctx context.Context, req *pb.CheckInfrastructureStatusRequest) (*pb.CheckInfrastructureStatusResponse, error) {
	result, err := s.uc.CheckInfrastructureStatus(ctx, req.GetScopeUuid(), req.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.CheckInfrastructureStatusResponse{}
	resp.SetResult(pb.CheckInfrastructureStatusResponse_Result(result))
	return resp, nil
}
