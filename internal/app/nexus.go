package app

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/openhdc/openhdc/api/nexus/v1"
	"github.com/openhdc/openhdc/api/nexus/v1/pbconnect"

	"github.com/openhdc/openhdc/internal/domain/model"
	"github.com/openhdc/openhdc/internal/domain/service"
)

type NexusApp struct {
	pbconnect.UnimplementedNexusHandler
	svc *service.NexusService
}

func NewNexusApp(svc *service.NexusService) *NexusApp {
	return &NexusApp{svc: svc}
}

var _ pbconnect.NexusHandler = (*NexusApp)(nil)

func (a *NexusApp) GetConfiguration(ctx context.Context, req *connect.Request[pb.GetConfigurationRequest]) (*connect.Response[pb.Configuration], error) {
	cfg, err := a.svc.GetConfiguration(ctx)
	if err != nil {
		return nil, err
	}
	res := toProtoConfiguration(cfg)
	return connect.NewResponse(res), nil
}

func toProtoNTPServer(addresses []string) *pb.Configuration_NTPServer {
	ret := &pb.Configuration_NTPServer{}
	ret.SetAddresses(addresses)
	return ret
}

func toProtoPackageRepositories(prs []*model.PackageRepository) []*pb.Configuration_PackageRepository {
	ret := []*pb.Configuration_PackageRepository{}
	for _, pr := range prs {
		ret = append(ret, toProtoPackageRepository(pr))
	}
	return ret
}

func toProtoPackageRepository(pr *model.PackageRepository) *pb.Configuration_PackageRepository {
	ret := &pb.Configuration_PackageRepository{}
	ret.SetId(int64(pr.ID))
	ret.SetName(pr.Name)
	ret.SetUrl(pr.URL)
	ret.SetEnabled(pr.Enabled)
	return ret
}

func toProtoBootResources(brs []*model.BootResource) []*pb.Configuration_BootResource {
	ret := []*pb.Configuration_BootResource{}
	for _, br := range brs {
		ret = append(ret, toProtoBootResource(br))
	}
	return ret
}

func toProtoBootResource(br *model.BootResource) *pb.Configuration_BootResource {
	ret := &pb.Configuration_BootResource{}
	ret.SetName(br.Name)
	ret.SetArchitecture(br.Architecture)
	ret.SetStatus(br.Status)
	ret.SetDefault(br.Default)
	return ret
}

func toProtoConfiguration(c *model.Configuration) *pb.Configuration {
	ret := &pb.Configuration{}
	ret.SetNtpServer(toProtoNTPServer(c.NTPServers))
	ret.SetPackageRepositories(toProtoPackageRepositories(c.PackageRepositories))
	ret.SetBootResources(toProtoBootResources(c.BootResources))
	return ret
}
