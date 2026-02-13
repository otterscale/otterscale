package app

import (
	"context"
	"maps"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/registry/v1"
	"github.com/otterscale/otterscale/api/registry/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/registry"
	"github.com/otterscale/otterscale/internal/core/registry/chart"
	"github.com/otterscale/otterscale/internal/core/registry/img"
)

type RegistryService struct {
	pbconnect.UnimplementedRegistryServiceHandler

	chart    *chart.UseCase
	registry *registry.UseCase
}

func NewRegistryService(chart *chart.UseCase, registry *registry.UseCase) *RegistryService {
	return &RegistryService{
		chart:    chart,
		registry: registry,
	}
}

var _ pbconnect.RegistryServiceHandler = (*RegistryService)(nil)

func (s *RegistryService) GetRegistryURL(_ context.Context, req *pb.GetRegistryURLRequest) (*pb.GetRegistryURLResponse, error) {
	url, err := s.registry.GetRegistryURL(req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.GetRegistryURLResponse{}
	resp.SetRegistryUrl(url)
	return resp, nil
}

func (s *RegistryService) ListRepositories(ctx context.Context, req *pb.ListRepositoriesRequest) (*pb.ListRepositoriesResponse, error) {
	repositories, err := s.registry.ListRepositories(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListRepositoriesResponse{}
	resp.SetRepositories(toProtoRepositories(repositories))
	return resp, nil
}

func (s *RegistryService) ListManifests(ctx context.Context, req *pb.ListManifestsRequest) (*pb.ListManifestsResponse, error) {
	manifests, err := s.registry.ListManifests(ctx, req.GetScope(), req.GetRepositoryName())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListManifestsResponse{}
	resp.SetManifests(toProtoManifests(manifests))
	return resp, nil
}

func (s *RegistryService) DeleteManifest(ctx context.Context, req *pb.DeleteManifestRequest) (*emptypb.Empty, error) {
	if err := s.registry.DeleteManifest(ctx, req.GetScope(), req.GetRepositoryName(), req.GetDigest()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *RegistryService) ListCharts(ctx context.Context, req *pb.ListChartsRequest) (*pb.ListChartsResponse, error) {
	charts, err := s.registry.ListCharts(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListChartsResponse{}
	resp.SetCharts(toProtoCharts(charts))
	return resp, nil
}

func (s *RegistryService) ListChartVersions(ctx context.Context, req *pb.ListChartVersionsRequest) (*pb.ListChartVersionsResponse, error) {
	versions, err := s.registry.ListChartVersions(ctx, req.GetScope(), req.GetRepositoryName())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListChartVersionsResponse{}
	resp.SetVersions(toProtoChartVersions(versions))
	return resp, nil
}

func (s *RegistryService) GetChartInformation(ctx context.Context, req *pb.GetChartInformationRequest) (*pb.Chart_Information, error) {
	info, err := s.chart.GetChartInformation(ctx, req.GetChartRef())
	if err != nil {
		return nil, err
	}

	resp := toProtoChartInformation(info)
	return resp, nil
}

func (s *RegistryService) SyncArtifactHub(ctx context.Context, req *pb.SyncArtifactHubRequest) (*emptypb.Empty, error) {
	if err := s.chart.SyncArtifactHub(ctx, req.GetRegistryUrl()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *RegistryService) ImportChart(ctx context.Context, req *pb.ImportChartRequest) (*emptypb.Empty, error) {
	if err := s.chart.Import(ctx, req.GetChartRef(), req.GetRegistryUrl()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func toProtoRepositories(rs []registry.Repository) []*pb.Repository {
	ret := []*pb.Repository{}

	for i := range rs {
		ret = append(ret, toProtoRepository(&rs[i]))
	}

	return ret
}

func toProtoRepository(r *registry.Repository) *pb.Repository {
	ret := &pb.Repository{}
	ret.SetName(r.Name)
	ret.SetManifestCount(r.ManifestCount)
	ret.SetSizeBytes(r.SizeBytes)
	ret.SetLatestTag(r.LatestTag)
	return ret
}

func toProtoManifests(ms []registry.Manifest) []*pb.Manifest {
	ret := []*pb.Manifest{}

	for i := range ms {
		ret = append(ret, toProtoManifest(&ms[i]))
	}

	return ret
}

func toProtoManifest(m *registry.Manifest) *pb.Manifest {
	ret := &pb.Manifest{}
	ret.SetRepositoryName(m.Repository)
	ret.SetTag(m.Tag)
	ret.SetDigest(m.Digest)
	ret.SetSizeBytes(m.SizeBytes)

	image := m.Image
	if image != nil {
		ret.SetImage(toProtoRegistryImage(image))
	}

	chart := m.Chart
	if chart != nil {
		ret.SetChart(toProtoChart(chart.Metadata, chart.Repository))
	}

	return ret
}

func toProtoRegistryImage(i *img.Image) *pb.Image {
	ret := &pb.Image{}

	createdAt := i.Created
	if createdAt != nil {
		ret.SetCreatedAt(timestamppb.New(*createdAt))
	}

	ret.SetAuthor(i.Author)
	ret.SetPlatform(toProtoRegistryImagePlatform(&i.Platform))
	ret.SetConfig(toProtoRegistryImageImageConfig(&i.Config))
	ret.SetRootFs(toProtoRegistryImageRootFS(&i.RootFS))
	return ret
}

func toProtoRegistryImagePlatform(p *img.Platform) *pb.Image_Platform {
	ret := &pb.Image_Platform{}
	ret.SetArchitecture(p.Architecture)
	ret.SetOs(p.OS)
	ret.SetOsVersion(p.OSVersion)
	ret.SetOsFeatures(p.OSFeatures)
	ret.SetVariant(p.Variant)
	return ret
}

func toProtoRegistryImageImageConfig(ic *img.ImageConfig) *pb.Image_Config {
	ret := &pb.Image_Config{}
	ret.SetUser(ic.User)

	exposedPorts := []string{}
	for k := range maps.Keys(ic.ExposedPorts) {
		exposedPorts = append(exposedPorts, k)
	}
	ret.SetExposedPorts(exposedPorts)

	ret.SetEnvironments(ic.Env)
	ret.SetEntrypoint(ic.Entrypoint)
	ret.SetCmd(ic.Cmd)

	volumes := []string{}
	for k := range maps.Keys(ic.Volumes) {
		volumes = append(volumes, k)
	}
	ret.SetVolumes(volumes)

	ret.SetWorkingDir(ic.WorkingDir)
	ret.SetLabels(ic.Labels)
	ret.SetStopSignal(ic.StopSignal)
	return ret
}

func toProtoRegistryImageRootFS(r *img.RootFS) *pb.Image_RootFS {
	ret := &pb.Image_RootFS{}
	ret.SetType(r.Type)

	diffIDs := []string{}
	for _, d := range r.DiffIDs {
		diffIDs = append(diffIDs, d.String())
	}
	ret.SetDiffIds(diffIDs)

	return ret
}

func toProtoCharts(cs []chart.Chart) []*pb.Chart {
	ret := []*pb.Chart{}

	for i := range cs {
		ret = append(ret, toProtoChart(cs[i].Metadata, cs[i].Repository))
	}

	return ret
}

func toProtoChart(c *chart.Metadata, repository string) *pb.Chart {
	ret := &pb.Chart{}
	ret.SetName(c.Name)
	ret.SetHome(c.Home)
	ret.SetSources(c.Sources)
	ret.SetVersion(c.Version)
	ret.SetDescription(c.Description)
	ret.SetKeywords(c.Keywords)
	ret.SetMaintainers(toProtoChartMaintainers(c.Maintainers))
	ret.SetIcon(c.Icon)
	ret.SetApiVersion(c.APIVersion)
	ret.SetCondition(c.Condition)
	ret.SetTags(c.Tags)
	ret.SetAppVersion(c.AppVersion)
	ret.SetDeprecated(c.Deprecated)
	ret.SetAnnotations(c.Annotations)
	ret.SetKubeVersion(c.KubeVersion)
	ret.SetDependencies(toProtoChartDependencies(c.Dependencies))
	ret.SetType(c.Type)
	ret.SetRepositoryName(repository)
	return ret
}

func toProtoChartMaintainers(ms []*chart.Maintainer) []*pb.Chart_Maintainer {
	ret := []*pb.Chart_Maintainer{}

	for i := range ms {
		ret = append(ret, toProtoChartMaintainer(ms[i]))
	}

	return ret
}

func toProtoChartMaintainer(m *chart.Maintainer) *pb.Chart_Maintainer {
	ret := &pb.Chart_Maintainer{}
	ret.SetName(m.Name)
	ret.SetEmail(m.Email)
	ret.SetUrl(m.URL)
	return ret
}

func toProtoChartDependencies(ds []*chart.Dependency) []*pb.Chart_Dependency {
	ret := []*pb.Chart_Dependency{}

	for i := range ds {
		ret = append(ret, toProtoChartDependency(ds[i]))
	}

	return ret
}

func toProtoChartDependency(d *chart.Dependency) *pb.Chart_Dependency {
	ret := &pb.Chart_Dependency{}
	ret.SetName(d.Name)
	ret.SetVersion(d.Version)
	ret.SetRepository(d.Repository)
	ret.SetCondition(d.Condition)
	ret.SetTags(d.Tags)
	ret.SetEnabled(d.Enabled)
	ret.SetAlias(d.Alias)
	return ret
}

func toProtoChartVersions(vs []chart.Version) []*pb.Chart_Version {
	ret := []*pb.Chart_Version{}

	for _, v := range vs {
		ret = append(ret, toProtoChartVersion(&v))
	}

	return ret
}

func toProtoChartVersion(v *chart.Version) *pb.Chart_Version {
	ret := &pb.Chart_Version{}
	ret.SetChartRef(v.ChartRef)
	ret.SetChartVersion(v.ChartVersion)
	ret.SetApplicationVersion(v.ApplicationVersion)
	return ret
}

func toProtoChartInformation(info *chart.Information) *pb.Chart_Information {
	ret := &pb.Chart_Information{}
	ret.SetReadme(info.Readme)
	ret.SetValues(info.Values)
	return ret
}
