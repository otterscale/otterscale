package app

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

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
