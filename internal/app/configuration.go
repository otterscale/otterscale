package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/otterscale/api/nexus/v1"
	"github.com/openhdc/otterscale/internal/domain/model"
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

func (a *NexusApp) CreateBootImage(ctx context.Context, req *connect.Request[pb.CreateBootImageRequest]) (*connect.Response[pb.Configuration_BootImage], error) {
	bi, err := a.svc.CreateBootImage(ctx, req.Msg.GetDistroSeries(), req.Msg.GetArchitectures())
	if err != nil {
		return nil, err
	}
	res := toProtoBootImage(bi)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) SetDefaultBootImage(ctx context.Context, req *connect.Request[pb.SetDefaultBootImageRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.SetDefaultBootImage(ctx, req.Msg.GetDistroSeries()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ImportBootImages(ctx context.Context, req *connect.Request[pb.ImportBootImagesRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.ImportBootImages(ctx); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) IsImportingBootImages(ctx context.Context, req *connect.Request[pb.IsImportingBootImagesRequest]) (*connect.Response[pb.IsImportingBootImagesResponse], error) {
	isImporting, err := a.svc.IsImportingBootImages(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.IsImportingBootImagesResponse{}
	res.SetImporting(isImporting)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ListBootImageSelections(ctx context.Context, req *connect.Request[pb.ListBootImageSelectionsRequest]) (*connect.Response[pb.ListBootImageSelectionsResponse], error) {
	bims, err := a.svc.ListBootImageSelections(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.ListBootImageSelectionsResponse{}
	res.SetBootImageSelections(toProtoBootImageSelections(bims))
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

func toProtoBootImages(bis []model.BootImage) []*pb.Configuration_BootImage {
	ret := []*pb.Configuration_BootImage{}
	for i := range bis {
		ret = append(ret, toProtoBootImage(&bis[i]))
	}
	return ret
}

func toProtoBootImage(bi *model.BootImage) *pb.Configuration_BootImage {
	ret := &pb.Configuration_BootImage{}
	ret.SetSource(bi.Source)
	ret.SetDistroSeries(bi.DistroSeries)
	ret.SetName(bi.Name)
	ret.SetArchitectureStatusMap(bi.ArchitectureStatusMap)
	ret.SetDefault(bi.Default)
	return ret
}

func toProtoBootImageSelections(biss []model.BootImageSelection) []*pb.Configuration_BootImageSelection {
	ret := []*pb.Configuration_BootImageSelection{}
	for i := range biss {
		ret = append(ret, toProtoBootImageSelection(&biss[i]))
	}
	return ret
}

func toProtoBootImageSelection(bis *model.BootImageSelection) *pb.Configuration_BootImageSelection {
	ret := &pb.Configuration_BootImageSelection{}
	ret.SetDistroSeries(bis.DistroSeries.String())
	ret.SetName(bis.Name)
	ret.SetArchitectures(bis.Architectures)
	return ret
}

func toProtoConfiguration(c *model.Configuration) *pb.Configuration {
	ret := &pb.Configuration{}
	ret.SetNtpServer(toProtoNTPServer(c.NTPServers))
	ret.SetPackageRepositories(toProtoPackageRepositories(c.PackageRepositories))
	ret.SetBootImages(toProtoBootImages(c.BootImages))
	return ret
}
