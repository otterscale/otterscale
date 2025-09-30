package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/otterscale/otterscale/api/configuration/v1"
	"github.com/otterscale/otterscale/api/configuration/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type ConfigurationService struct {
	pbconnect.UnimplementedConfigurationServiceHandler

	uc *core.ConfigurationUseCase
}

func NewConfigurationService(uc *core.ConfigurationUseCase) *ConfigurationService {
	return &ConfigurationService{uc: uc}
}

var _ pbconnect.ConfigurationServiceHandler = (*ConfigurationService)(nil)

func (s *ConfigurationService) GetConfiguration(ctx context.Context, _ *pb.GetConfigurationRequest) (*pb.Configuration, error) {
	config, err := s.uc.GetConfiguration(ctx)
	if err != nil {
		return nil, err
	}
	resp := toProtoConfiguration(config)
	return resp, nil
}

func (s *ConfigurationService) UpdateNTPServer(ctx context.Context, req *pb.UpdateNTPServerRequest) (*pb.Configuration_NTPServer, error) {
	ntpServers, err := s.uc.UpdateNTPServer(ctx, req.GetAddresses())
	if err != nil {
		return nil, err
	}
	resp := toProtoNTPServer(ntpServers)
	return resp, nil
}

func (s *ConfigurationService) UpdatePackageRepository(ctx context.Context, req *pb.UpdatePackageRepositoryRequest) (*pb.Configuration_PackageRepository, error) {
	repo, err := s.uc.UpdatePackageRepository(ctx, int(req.GetId()), req.GetUrl(), req.GetSkipJuju())
	if err != nil {
		return nil, err
	}
	resp := toProtoPackageRepository(repo)
	return resp, nil
}

func (s *ConfigurationService) CreateBootImage(ctx context.Context, req *pb.CreateBootImageRequest) (*pb.Configuration_BootImage, error) {
	image, err := s.uc.CreateBootImage(ctx, req.GetDistroSeries(), req.GetArchitectures())
	if err != nil {
		return nil, err
	}
	resp := toProtoBootImage(image)
	return resp, nil
}

func (s *ConfigurationService) SetDefaultBootImage(ctx context.Context, req *pb.SetDefaultBootImageRequest) (*emptypb.Empty, error) {
	if err := s.uc.SetDefaultBootImage(ctx, req.GetDistroSeries()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ConfigurationService) ImportBootImages(ctx context.Context, _ *pb.ImportBootImagesRequest) (*emptypb.Empty, error) {
	if err := s.uc.ImportBootImages(ctx); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *ConfigurationService) IsImportingBootImages(ctx context.Context, _ *pb.IsImportingBootImagesRequest) (*pb.IsImportingBootImagesResponse, error) {
	isImporting, err := s.uc.IsImportingBootImages(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.IsImportingBootImagesResponse{}
	resp.SetImporting(isImporting)
	return resp, nil
}

func (s *ConfigurationService) ListBootImageSelections(_ context.Context, _ *pb.ListBootImageSelectionsRequest) (*pb.ListBootImageSelectionsResponse, error) {
	selections, err := s.uc.ListBootImageSelections()
	if err != nil {
		return nil, err
	}
	resp := &pb.ListBootImageSelectionsResponse{}
	resp.SetBootImageSelections(toProtoBootImageSelections(selections))
	return resp, nil
}

func toProtoNTPServer(addresses []string) *pb.Configuration_NTPServer {
	ret := &pb.Configuration_NTPServer{}
	ret.SetAddresses(addresses)
	return ret
}

func toProtoPackageRepositories(prs []core.PackageRepository) []*pb.Configuration_PackageRepository {
	ret := []*pb.Configuration_PackageRepository{}
	for i := range prs {
		ret = append(ret, toProtoPackageRepository(&prs[i]))
	}
	return ret
}

func toProtoPackageRepository(pr *core.PackageRepository) *pb.Configuration_PackageRepository {
	ret := &pb.Configuration_PackageRepository{}
	ret.SetId(int64(pr.ID))
	ret.SetName(pr.Name)
	ret.SetUrl(pr.URL)
	ret.SetEnabled(pr.Enabled)
	return ret
}

func toProtoBootImages(bis []core.BootImage) []*pb.Configuration_BootImage {
	ret := []*pb.Configuration_BootImage{}
	for i := range bis {
		ret = append(ret, toProtoBootImage(&bis[i]))
	}
	return ret
}

func toProtoBootImage(bi *core.BootImage) *pb.Configuration_BootImage {
	ret := &pb.Configuration_BootImage{}
	ret.SetSource(bi.Source)
	ret.SetDistroSeries(bi.DistroSeries)
	ret.SetName(bi.Name)
	ret.SetArchitectureStatusMap(bi.ArchitectureStatusMap)
	ret.SetDefault(bi.Default)
	return ret
}

func toProtoBootImageSelections(biss []core.BootImageSelection) []*pb.Configuration_BootImageSelection {
	ret := []*pb.Configuration_BootImageSelection{}
	for i := range biss {
		ret = append(ret, toProtoBootImageSelection(&biss[i]))
	}
	return ret
}

func toProtoBootImageSelection(bis *core.BootImageSelection) *pb.Configuration_BootImageSelection {
	ret := &pb.Configuration_BootImageSelection{}
	ret.SetDistroSeries(bis.DistroSeries.String())
	ret.SetName(bis.Name)
	ret.SetArchitectures(bis.Architectures)
	return ret
}

func toProtoConfiguration(c *core.Configuration) *pb.Configuration {
	ret := &pb.Configuration{}
	ret.SetNtpServer(toProtoNTPServer(c.NTPServers))
	ret.SetPackageRepositories(toProtoPackageRepositories(c.PackageRepositories))
	ret.SetBootImages(toProtoBootImages(c.BootImages))
	return ret
}
