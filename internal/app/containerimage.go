package app

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/otterscale/otterscale/api/containerimage/v1"
	"github.com/otterscale/otterscale/api/containerimage/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/containerimage"
)

type ContainerImageService struct {
	pbconnect.UnimplementedContainerImageServiceHandler

	containerimage *containerimage.UseCase
}

func NewContainerImageService(containerimage *containerimage.UseCase) *ContainerImageService {
	return &ContainerImageService{
		containerimage: containerimage,
	}
}

var _ pbconnect.ContainerImageServiceHandler = (*ContainerImageService)(nil)

func (s *ContainerImageService) ListContainerImages(ctx context.Context, req *pb.ListContainerImagesRequest) (*pb.ListContainerImagesResponse, error) {
	images, err := s.containerimage.ListContainerImages(ctx, req.GetScope(), req.GetNamespace(), req.GetEndpoint())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListContainerImagesResponse{}
	resp.SetImages(toProtoContainerImages(images))
	return resp, nil
}

func (s *ContainerImageService) UploadContainerImage(ctx context.Context, req *pb.UploadContainerImageRequest) (*emptypb.Empty, error) {
	err := s.containerimage.UploadContainerImage(ctx, req.GetScope(), req.GetNamespace(), req.GetEndpoint(), req.GetName(), req.GetTag(), req.GetImageTar())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ContainerImageService) DeleteContainerImage(ctx context.Context, req *pb.DeleteContainerImageRequest) (*emptypb.Empty, error) {
	err := s.containerimage.DeleteContainerImage(ctx, req.GetScope(), req.GetNamespace(), req.GetEndpoint(), req.GetName(), req.GetTag())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func toProtoContainerImages(images []containerimage.ContainerImage) []*pb.ContainerImage {
	ret := []*pb.ContainerImage{}

	for i := range images {
		ret = append(ret, toProtoContainerImage(&images[i]))
	}

	return ret
}

func toProtoContainerImage(i *containerimage.ContainerImage) *pb.ContainerImage {
	ret := &pb.ContainerImage{}
	ret.SetNamespace(i.Namespace)
	ret.SetName(i.Name)
	ret.SetRef(i.Ref)
	ret.SetSize(i.Size)
	ret.SetTag(i.Tag)
	ret.SetCreateAt(timestamppb.New(i.CreateAt))
	return ret
}
