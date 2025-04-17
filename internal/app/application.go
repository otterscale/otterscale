package app

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/openhdc/openhdc/api/nexus/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *NexusApp) ListApplications(ctx context.Context, req *connect.Request[pb.ListApplicationsRequest]) (*connect.Response[pb.ListApplicationsResponse], error) {
	as, err := a.svc.ListApplications(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	res := &pb.ListApplicationsResponse{}
	res.SetApplications(toProtoApplications(as))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) GetApplication(ctx context.Context, req *connect.Request[pb.GetApplicationRequest]) (*connect.Response[pb.Application], error) {
	app, err := a.svc.GetApplication(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	res := toProtoApplication(app)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ListReleases(ctx context.Context, req *connect.Request[pb.ListReleasesRequest]) (*connect.Response[pb.ListReleasesResponse], error) {
	rs, err := a.svc.ListReleases(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.ListReleasesResponse{}
	res.SetReleases(toProtoReleases(rs))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) CreateRelease(ctx context.Context, req *connect.Request[pb.CreateReleaseRequest]) (*connect.Response[pb.Application_Release], error) {
	r, err := a.svc.CreateRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun(), req.Msg.GetChartRef(), req.Msg.GetValuesYaml())
	if err != nil {
		return nil, err
	}
	res := toProtoRelease(r)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) UpdateRelease(ctx context.Context, req *connect.Request[pb.UpdateReleaseRequest]) (*connect.Response[pb.Application_Release], error) {
	r, err := a.svc.UpdateRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun(), req.Msg.GetChartRef(), req.Msg.GetValuesYaml())
	if err != nil {
		return nil, err
	}
	res := toProtoRelease(r)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) DeleteRelease(ctx context.Context, req *connect.Request[pb.DeleteReleaseRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.DeleteRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) RollbackRelease(ctx context.Context, req *connect.Request[pb.RollbackReleaseRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.RollbackRelease(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ListCharts(ctx context.Context, req *connect.Request[pb.ListChartsRequest]) (*connect.Response[pb.ListChartsResponse], error) {
	cs, err := a.svc.ListCharts(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.ListChartsResponse{}
	res.SetCharts(toProtoCharts(cs))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) GetChart(ctx context.Context, req *connect.Request[pb.GetChartRequest]) (*connect.Response[pb.Application_Release_Chart], error) {
	c, err := a.svc.GetChart(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	md := &chart.Metadata{}
	if len(c.Versions) > 0 {
		md = c.Versions[0].Metadata
	}
	res := toProtoChart(md, c.Versions...)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) GetChartMetadata(ctx context.Context, req *connect.Request[pb.GetChartMetadataRequest]) (*connect.Response[pb.Application_Release_Chart_Metadata], error) {
	md, err := a.svc.GetChartMetadata(ctx, req.Msg.GetChartRef())
	if err != nil {
		return nil, err
	}
	res := toProtoChartMetadata(md)
	return connect.NewResponse(res), nil
}

func toProtoApplications(as []model.Application) []*pb.Application {
	ret := []*pb.Application{}
	for i := range as {
		ret = append(ret, toProtoApplication(&as[i]))
	}
	return ret
}

func toProtoApplication(a *model.Application) *pb.Application {
	ret := &pb.Application{}
	return ret
}

func toProtoReleases(rs []model.Release) []*pb.Application_Release {
	ret := []*pb.Application_Release{}
	for i := range rs {
		ret = append(ret, toProtoRelease(&rs[i]))
	}
	return ret
}

func toProtoRelease(r *model.Release) *pb.Application_Release {
	ret := &pb.Application_Release{}
	return ret
}

func toProtoCharts(cs []model.Chart) []*pb.Application_Release_Chart {
	ret := []*pb.Application_Release_Chart{}
	for i := range cs {
		if len(cs[i].Versions) > 0 {
			ret = append(ret, toProtoChart(cs[i].Versions[0].Metadata, cs[i].Versions[0])) // latest only
		}
	}
	return ret
}

func toProtoChart(cmd *chart.Metadata, vs ...*repo.ChartVersion) *pb.Application_Release_Chart {
	ret := &pb.Application_Release_Chart{}
	ret.SetName(cmd.Name)
	ret.SetIcon(cmd.Icon)
	ret.SetDescription(cmd.Description)
	ret.SetDeprecated(cmd.Deprecated)
	ret.SetTags(cmd.Tags)
	ret.SetKeywords(cmd.Keywords)
	ret.SetHome(cmd.Home)
	ret.SetSources(cmd.Sources)
	ret.SetMaintainers(toProtoChartMaintainers(cmd.Maintainers))
	ret.SetDependencies(toProtoChartDependencies(cmd.Dependencies))
	ret.SetVersions(toProtoChartVersions(vs...))
	return ret
}

func toProtoChartMaintainers(ms []*chart.Maintainer) []*pb.Application_Release_Chart_Maintainer {
	ret := []*pb.Application_Release_Chart_Maintainer{}
	for _, m := range ms {
		maintainer := &pb.Application_Release_Chart_Maintainer{}
		maintainer.SetName(m.Name)
		maintainer.SetEmail(m.Email)
		maintainer.SetUrl(m.URL)
		ret = append(ret, maintainer)
	}
	return ret
}

func toProtoChartDependencies(ds []*chart.Dependency) []*pb.Application_Release_Chart_Dependency {
	ret := []*pb.Application_Release_Chart_Dependency{}
	for _, d := range ds {
		dependency := &pb.Application_Release_Chart_Dependency{}
		dependency.SetName(d.Name)
		dependency.SetVersion(d.Version)
		dependency.SetCondition(d.Condition)
		ret = append(ret, dependency)
	}
	return ret
}

func toProtoChartVersions(vs ...*repo.ChartVersion) []*pb.Application_Release_Chart_Version {
	ret := []*pb.Application_Release_Chart_Version{}
	for _, v := range vs {
		version := &pb.Application_Release_Chart_Version{}
		version.SetChartVersion(v.Version)
		version.SetApplicationVersion(v.AppVersion)
		if len(v.URLs) > 0 {
			version.SetChartRef(v.URLs[0])
		}
		ret = append(ret, version)
	}
	return ret
}

func toProtoChartMetadata(md *model.ChartMetadata) *pb.Application_Release_Chart_Metadata {
	ret := &pb.Application_Release_Chart_Metadata{}
	ret.SetValuesYaml(md.ValuesYAML)
	ret.SetReadmeMd(md.ReadmeMD)
	return ret
}
