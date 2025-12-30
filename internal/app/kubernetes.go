package app

import (
	"context"

	pb "github.com/otterscale/otterscale/api/kubernetes/v1"
	"github.com/otterscale/otterscale/api/kubernetes/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
)

type KubernetesService struct {
	pbconnect.UnimplementedKubernetesServiceHandler

	kubernetes *kubernetes.Kubernetes
}

func NewKubernetesService(kubernetes *kubernetes.Kubernetes) *KubernetesService {
	return &KubernetesService{
		kubernetes: kubernetes,
	}
}

var _ pbconnect.KubernetesServiceHandler = (*KubernetesService)(nil)

func (s *KubernetesService) ValidateKubeConfig(ctx context.Context, req *pb.ValidateKubeConfigRequest) (*pb.ValidateKubeConfigResponse, error) {
	err := s.kubernetes.ValidateKubeConfig(ctx, req.GetKubeconfig())
	if err != nil {
		resp := &pb.ValidateKubeConfigResponse{}
		resp.SetMessage(err.Error())
		return resp, nil
	}

	resp := &pb.ValidateKubeConfigResponse{}
	resp.SetMessage("")
	return resp, nil
}
