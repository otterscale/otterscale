package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/otterscale/otterscale/api/tag/v1"
	"github.com/otterscale/otterscale/api/tag/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type TagService struct {
	pbconnect.UnimplementedTagServiceHandler

	uc *core.TagUseCase
}

func NewTagService(uc *core.TagUseCase) *TagService {
	return &TagService{uc: uc}
}

var _ pbconnect.TagServiceHandler = (*TagService)(nil)

func (s *TagService) ListTags(ctx context.Context, req *connect.Request[pb.ListTagsRequest]) (*connect.Response[pb.ListTagsResponse], error) {
	tags, err := s.uc.ListTags(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListTagsResponse{}
	resp.SetTags(toProtoTags(tags))
	return connect.NewResponse(resp), nil
}

func (s *TagService) GetTag(ctx context.Context, req *connect.Request[pb.GetTagRequest]) (*connect.Response[pb.Tag], error) {
	tag, err := s.uc.GetTag(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoTag(tag)
	return connect.NewResponse(resp), nil
}

func (s *TagService) CreateTag(ctx context.Context, req *connect.Request[pb.CreateTagRequest]) (*connect.Response[pb.Tag], error) {
	tag, err := s.uc.CreateTag(ctx, req.Msg.GetName(), req.Msg.GetComment())
	if err != nil {
		return nil, err
	}
	resp := toProtoTag(tag)
	return connect.NewResponse(resp), nil
}

func (s *TagService) DeleteTag(ctx context.Context, req *connect.Request[pb.DeleteTagRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteTag(ctx, req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func toProtoTags(ts []core.Tag) []*pb.Tag {
	ret := []*pb.Tag{}
	for i := range ts {
		ret = append(ret, toProtoTag(&ts[i]))
	}
	return ret
}

func toProtoTag(t *core.Tag) *pb.Tag {
	ret := &pb.Tag{}
	ret.SetName(t.Name)
	ret.SetComment(t.Comment)
	return ret
}
