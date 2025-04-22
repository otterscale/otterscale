package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/openhdc/api/nexus/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *NexusApp) ListTags(ctx context.Context, req *connect.Request[pb.ListTagsRequest]) (*connect.Response[pb.ListTagsResponse], error) {
	ts, err := a.svc.ListTags(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.ListTagsResponse{}
	res.SetTags(toProtoTags(ts))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) GetTag(ctx context.Context, req *connect.Request[pb.GetTagRequest]) (*connect.Response[pb.Tag], error) {
	t, err := a.svc.GetTag(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	res := toProtoTag(t)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) CreateTag(ctx context.Context, req *connect.Request[pb.CreateTagRequest]) (*connect.Response[pb.Tag], error) {
	t, err := a.svc.CreateTag(ctx, req.Msg.GetName(), req.Msg.GetComment())
	if err != nil {
		return nil, err
	}
	res := toProtoTag(t)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) DeleteTag(ctx context.Context, req *connect.Request[pb.DeleteTagRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.DeleteTag(ctx, req.Msg.GetName()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func toProtoTags(ts []model.Tag) []*pb.Tag {
	ret := []*pb.Tag{}
	for i := range ts {
		ret = append(ret, toProtoTag(&ts[i]))
	}
	return ret
}

func toProtoTag(t *model.Tag) *pb.Tag {
	ret := &pb.Tag{}
	ret.SetName(t.Name)
	ret.SetComment(t.Comment)
	return ret
}
