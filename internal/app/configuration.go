package app

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/openhdc/openhdc/api/nexus/v1"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *NexusApp) GetConfiguration(ctx context.Context, req *connect.Request[pb.GetConfigurationRequest]) (*connect.Response[pb.Configuration], error) {
	cfg, err := a.svc.GetConfiguration(ctx)
	if err != nil {
		return nil, err
	}
	res := toProtoConfiguration(cfg)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) UpdateNTPServer(ctx context.Context, req *connect.Request[pb.UpdateNTPServerRequest]) (*connect.Response[pb.Configuration_NTPServer], error) {
	ntpServers, err := a.svc.UpdateNTPServer(ctx, req.Msg.GetAddresses())
	if err != nil {
		return nil, err
	}
	res := toProtoNTPServer(ntpServers)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) UpdatePackageRepository(ctx context.Context, req *connect.Request[pb.UpdatePackageRepositoryRequest]) (*connect.Response[pb.Configuration_PackageRepository], error) {
	pr, err := a.svc.UpdatePackageRepository(ctx, int(req.Msg.GetId()), req.Msg.GetUrl(), req.Msg.GetSkipJuju())
	if err != nil {
		return nil, err
	}
	res := toProtoPackageRepository(pr)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) UpdateDefaultBootResource(ctx context.Context, req *connect.Request[pb.UpdateDefaultBootResourceRequest]) (*connect.Response[pb.Configuration_BootResource], error) {
	br, err := a.svc.UpdateDefaultBootResource(ctx, req.Msg.GetDistroSeries())
	if err != nil {
		return nil, err
	}
	res := toProtoBootResource(br)
	return connect.NewResponse(res), nil
}

func toProtoNTPServer(addresses []string) *pb.Configuration_NTPServer {
	ret := &pb.Configuration_NTPServer{}
	ret.SetAddresses(addresses)
	return ret
}

func toProtoPackageRepositories(prs []model.PackageRepository) []*pb.Configuration_PackageRepository {
	ret := []*pb.Configuration_PackageRepository{}
	for i := range prs {
		ret = append(ret, toProtoPackageRepository(&prs[i]))
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

func toProtoBootResources(brs []model.BootResource) []*pb.Configuration_BootResource {
	ret := []*pb.Configuration_BootResource{}
	for i := range brs {
		ret = append(ret, toProtoBootResource(&brs[i]))
	}
	return ret
}

func toProtoBootResource(br *model.BootResource) *pb.Configuration_BootResource {
	ret := &pb.Configuration_BootResource{}
	ret.SetName(br.Name)
	ret.SetArchitecture(br.Architecture)
	ret.SetStatus(br.Status)
	ret.SetDefault(br.Default)
	ret.SetDistroSeries(br.DistroSeries)
	return ret
}

func toProtoConfiguration(c *model.Configuration) *pb.Configuration {
	ret := &pb.Configuration{}
	ret.SetNtpServer(toProtoNTPServer(c.NTPServers))
	ret.SetPackageRepositories(toProtoPackageRepositories(c.PackageRepositories))
	ret.SetBootResources(toProtoBootResources(c.BootResources))
	return ret
}
