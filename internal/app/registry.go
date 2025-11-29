package app

import (
	"context"

	pb "github.com/otterscale/otterscale/api/registry/v1"
	"github.com/otterscale/otterscale/api/registry/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/registry"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RegistryService struct {
	pbconnect.UnimplementedRegistryServiceHandler

	registry *registry.UseCase
}

func NewRegistryService(registry *registry.UseCase) *RegistryService {
	return &RegistryService{
		registry: registry,
	}
}

var _ pbconnect.RegistryServiceHandler = (*RegistryService)(nil)

func (s *RegistryService) GetRegistryURL(ctx context.Context, req *pb.GetRegistryURLRequest) (*pb.GetRegistryURLResponse, error) {
	url, err := s.registry.GetRegistryURL(ctx, req.GetScope())
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
		ret.SetChart(toProtoRegistryChart(chart))
	}

	return ret
}

func toProtoRegistryImage(i *registry.Image) *pb.Image {
	ret := &pb.Image{}
	// ret.SetArchitecture(i.Architecture)
	// ret.SetOs(i.Os)
	// ret.SetCreatedAt(i.CreatedAt.String())
	// ret.SetAuthor(i.Author)
	// ret.SetDockerVersion(i.DockerVersion)
	// ret.SetConfigJson(string(i.ConfigJson))
	return ret
}

func toProtoRegistryChart(c *registry.Chart) *pb.Chart {
	ret := &pb.Chart{}
	// ret.SetName(c.Name)
	// ret.SetVersion(c.Version)
	// ret.SetDescription(c.Description)
	// ret.SetApiVersion(c.ApiVersion)
	// ret.SetAppVersion(c.AppVersion)
	// ret.SetKeywords(c.Keywords)
	// ret.SetHome(c.Home)
	// ret.SetSources(c.Sources)
	// ret.SetMaintainers(toProtoChartMaintainers(c.Maintainers))
	// ret.SetDependencies(toProtoChartDependencies(c.Dependencies))
	return ret
}
