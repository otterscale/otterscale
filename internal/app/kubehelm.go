package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"

	v1 "github.com/openhdc/openhdc/api/kube/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *KubeApp) InstallRelease(ctx context.Context, req *connect.Request[v1.InstallReleaseRequest]) (*connect.Response[v1.Release], error) {
	values := map[string]any{}
	if err := yaml.Unmarshal([]byte(req.Msg.GetValuesYaml()), &values); err != nil {
		return nil, err
	}
	rel, err := a.svc.InstallRelease(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun(), req.Msg.GetChartRef(), values)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(a.toRelease(rel)), nil
}

func (a *KubeApp) UninstallRelease(ctx context.Context, req *connect.Request[v1.UninstallReleaseRequest]) (*connect.Response[v1.Release], error) {
	rel, err := a.svc.UninstallRelease(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(a.toRelease(rel)), nil
}

func (a *KubeApp) UpgradeRelease(ctx context.Context, req *connect.Request[v1.UpgradeReleaseRequest]) (*connect.Response[v1.Release], error) {
	values := map[string]any{}
	if err := yaml.Unmarshal([]byte(req.Msg.GetValuesYaml()), &values); err != nil {
		return nil, err
	}
	rel, err := a.svc.UpgradeRelease(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun(), req.Msg.GetChartRef(), values)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(a.toRelease(rel)), nil
}

func (a *KubeApp) RollbackRelease(ctx context.Context, req *connect.Request[v1.RollbackReleaseRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.RollbackRelease(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetDryRun()); err != nil {
		return nil, err
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

// func (a *KubeApp) ListReleases(ctx context.Context, req *connect.Request[v1.ListReleasesRequest]) (*connect.Response[v1.ListReleasesResponse], error) {
// 	rels, err := a.svc.ListReleases(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName(), "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, rel := range rels {
// 		fmt.Println(rel.Name, rel.Namespace)
// 	}
// 	res := &v1.ListReleasesResponse{}
// 	return connect.NewResponse(res), nil
// }

func (a *KubeApp) ListRepositories(_ context.Context, req *connect.Request[v1.ListRepositoriesRequest]) (*connect.Response[v1.ListRepositoriesResponse], error) {
	repos, err := a.svc.ListRepositories()
	if err != nil {
		return nil, err
	}
	res := &v1.ListRepositoriesResponse{}
	res.SetRepositories(a.toRepositories(repos))
	return connect.NewResponse(res), nil
}

func (a *KubeApp) UpdateRepositoryCharts(_ context.Context, req *connect.Request[v1.UpdateRepositoryChartsRequest]) (*connect.Response[v1.Repository], error) {
	repo, err := a.svc.UpdateRepositoryCharts(req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(a.toRepository(repo)), nil
}

func (a *KubeApp) ListCharts(_ context.Context, req *connect.Request[v1.ListChartsRequest]) (*connect.Response[v1.ListChartsResponse], error) {
	m, err := a.svc.ListCharts()
	if err != nil {
		return nil, err
	}
	res := &v1.ListChartsResponse{}
	res.SetCharts(a.toLatestCharts(m))
	return connect.NewResponse(res), nil
}

func (a *KubeApp) GetChart(_ context.Context, req *connect.Request[v1.GetChartRequest]) (*connect.Response[v1.Chart], error) {
	cvs, err := a.svc.GetChart(req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(a.toChart(cvs)), nil
}

func (a *KubeApp) toRepositories(hrs []*model.HelmRepo) []*v1.Repository {
	ret := []*v1.Repository{}
	for idx := range hrs {
		ret = append(ret, a.toRepository(hrs[idx]))
	}
	return ret
}

func (a *KubeApp) toRepository(hr *model.HelmRepo) *v1.Repository {
	ret := &v1.Repository{}
	ret.SetName(hr.Name)
	ret.SetUrl(hr.URL)
	ret.SetUsername(hr.Username)
	ret.SetPassword(hr.Password)
	ret.SetCertFile(hr.CertFile)
	ret.SetKeyFile(hr.KeyFile)
	ret.SetCaFile(hr.CAFile)
	ret.SetInsecureSkipTlsVerify(hr.InsecureSkipTLSverify)
	ret.SetPassCredentialsAll(hr.PassCredentialsAll)
	ret.SetChartNames(hr.ChartNames)
	return ret
}

func (a *KubeApp) toLatestCharts(m map[string]repo.ChartVersions) []*v1.Chart {
	ret := []*v1.Chart{}
	for _, cvs := range m {
		if len(cvs) == 0 {
			continue
		}
		latest := cvs[0]
		chart := a.metadataToChart(latest.Metadata)
		chart.SetVersions(a.toChartVersions(latest)) // only latest
		ret = append(ret, chart)
	}
	return ret
}

func (a *KubeApp) toChart(cvs repo.ChartVersions) *v1.Chart {
	ret := &v1.Chart{}
	for _, cv := range cvs {
		ret = a.metadataToChart(cv.Metadata)
		break
	}
	ret.SetVersions(a.toChartVersions(cvs...)) // all
	return ret
}

func (a *KubeApp) metadataToChart(cv *chart.Metadata) *v1.Chart {
	ret := &v1.Chart{}
	ret.SetName(cv.Name)
	ret.SetIcon(cv.Icon)
	ret.SetDescription(cv.Description)
	ret.SetDeprecated(cv.Deprecated)
	ret.SetTags(cv.Tags)
	ret.SetKeywords(cv.Keywords)
	ret.SetHome(cv.Home)
	ret.SetSources(cv.Sources)
	ret.SetMaintainers(a.toChartMaintainers(cv.Maintainers))
	ret.SetDependencies(a.toChartDependencies(cv.Dependencies))
	return ret
}

func (a *KubeApp) toChartMaintainers(ms []*chart.Maintainer) []*v1.Chart_Maintainer {
	ret := []*v1.Chart_Maintainer{}
	for _, m := range ms {
		maintainer := &v1.Chart_Maintainer{}
		maintainer.SetName(m.Name)
		maintainer.SetEmail(m.Email)
		maintainer.SetUrl(m.URL)
		ret = append(ret, maintainer)
	}
	return ret
}

func (a *KubeApp) toChartDependencies(ds []*chart.Dependency) []*v1.Chart_Dependency {
	ret := []*v1.Chart_Dependency{}
	for _, d := range ds {
		dependency := &v1.Chart_Dependency{}
		dependency.SetName(d.Name)
		dependency.SetVersion(d.Version)
		dependency.SetCondition(d.Condition)
		ret = append(ret, dependency)
	}
	return ret
}

func (a *KubeApp) toChartVersions(cvs ...*repo.ChartVersion) []*v1.Chart_Version {
	ret := []*v1.Chart_Version{}
	for _, cv := range cvs {
		ret = append(ret, a.toChartVersion(cv))
	}
	return ret
}

func (a *KubeApp) toChartVersion(cv *repo.ChartVersion) *v1.Chart_Version {
	ret := &v1.Chart_Version{}
	ret.SetChartVersion(cv.Version)
	ret.SetApplicationVersion(cv.AppVersion)
	if len(cv.URLs) > 0 {
		ret.SetChartRef(cv.URLs[0])
	}
	return ret
}

func (a *KubeApp) toRelease(rel *release.Release) *v1.Release {
	ret := &v1.Release{}
	ret.SetName(rel.Name)
	ret.SetNamespace(rel.Namespace)
	ret.SetRevision(int32(rel.Version)) //nolint:gosec
	return ret
}
