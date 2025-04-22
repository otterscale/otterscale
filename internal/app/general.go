package app

import (
	"context"

	"connectrpc.com/connect"

	pb "github.com/openhdc/openhdc/api/nexus/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *NexusApp) VerifyEnvironment(ctx context.Context, req *connect.Request[pb.VerifyEnvironmentRequest]) (*connect.Response[pb.VerifyEnvironmentResponse], error) {
	es, err := a.svc.VerifyEnvironment(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.VerifyEnvironmentResponse{}
	res.SetErrors(toProtoErrors(es))
	return connect.NewResponse(res), nil
}

func toProtoErrors(es []model.Error) []*pb.Error {
	ret := []*pb.Error{}
	for i := range es {
		ret = append(ret, toProtoError(&es[i]))
	}
	return ret
}

func toProtoError(e *model.Error) *pb.Error {
	ret := &pb.Error{}
	ret.SetCode(e.Code)
	ret.SetLevel(pb.ErrorLevel(e.Level))
	ret.SetMessage(e.Message)
	ret.SetDetails(e.Details)
	ret.SetUrl(e.URL)
	return ret
}
