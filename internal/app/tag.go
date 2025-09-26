package app

import (
	"context"

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

func (s *TagService) ListTags(ctx context.Context, _ *pb.ListTagsRequest) (*pb.ListTagsResponse, error) {
	tags, err := s.uc.ListTags(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListTagsResponse{}
	resp.SetTags(toProtoTags(tags))
	return resp, nil
}

func (s *TagService) GetTag(ctx context.Context, req *pb.GetTagRequest) (*pb.Tag, error) {
	tag, err := s.uc.GetTag(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoTag(tag)
	return resp, nil
}

func (s *TagService) CreateTag(ctx context.Context, req *pb.CreateTagRequest) (*pb.Tag, error) {
	tag, err := s.uc.CreateTag(ctx, req.GetName(), req.GetComment())
	if err != nil {
		return nil, err
	}
	resp := toProtoTag(tag)
	return resp, nil
}

func (s *TagService) DeleteTag(ctx context.Context, req *pb.DeleteTagRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteTag(ctx, req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
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
