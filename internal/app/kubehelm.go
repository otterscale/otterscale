package app

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"

	v1 "github.com/openhdc/openhdc/api/kube/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *KubeApp) ListReleases(ctx context.Context, req *connect.Request[v1.ListReleasesRequest]) (*connect.Response[v1.ListReleasesResponse], error) {
	rels, err := a.svc.ListReleases(ctx, req.Msg.GetModelUuid(), req.Msg.GetClusterName(), "")
	if err != nil {
		return nil, err
	}
	for _, rel := range rels {
		fmt.Println(rel.Name, rel.Namespace)
	}
	res := &v1.ListReleasesResponse{}
	return connect.NewResponse(res), nil
}

func (a *KubeApp) ListRepositories(ctx context.Context, req *connect.Request[v1.ListRepositoriesRequest]) (*connect.Response[v1.ListRepositoriesResponse], error) {
	repos, err := a.svc.ListRepositories(ctx)
	if err != nil {
		return nil, err
	}
	res := &v1.ListRepositoriesResponse{}
	res.SetRepositories(a.toRepositories(repos))
	return connect.NewResponse(res), nil
}

func (a *KubeApp) UpdateRepositoryCharts(ctx context.Context, req *connect.Request[v1.UpdateRepositoryChartsRequest]) (*connect.Response[v1.Repository], error) {
	repo, err := a.svc.UpdateRepositoryCharts(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(a.toRepository(repo)), nil
}

func (a *KubeApp) ListCharts(ctx context.Context, req *connect.Request[v1.ListChartsRequest]) (*connect.Response[v1.ListChartsResponse], error) {
	m, err := a.svc.ListCharts(ctx)
	if err != nil {
		return nil, err
	}
	res := &v1.ListChartsResponse{}
	res.SetCharts(a.toLatestCharts(m))
	return connect.NewResponse(res), nil
}

func (a *KubeApp) GetChart(ctx context.Context, req *connect.Request[v1.GetChartRequest]) (*connect.Response[v1.Chart], error) {
	cvs, err := a.svc.GetChart(ctx, req.Msg.GetName())
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
