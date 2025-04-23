package app

import (
	"context"

	"connectrpc.com/connect"

	pb "github.com/openhdc/openhdc/api/nexus/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *NexusApp) VerifyEnvironment(ctx context.Context, req *connect.Request[pb.VerifyEnvironmentRequest]) (*connect.Response[pb.VerifyEnvironmentResponse], error) {
	es, err := a.svc.VerifyEnvironment(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.VerifyEnvironmentResponse{}
	res.SetErrors(toProtoErrors(es))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ListCephes(ctx context.Context, req *connect.Request[pb.ListCephesRequest]) (*connect.Response[pb.ListCephesResponse], error) {
	cephes, err := a.svc.ListCephes(ctx)
	if err != nil {
		return nil, err
	}
	a.svc.GetCephInfo(ctx, "1a675505-5618-4578-8350-b9fbc19ac78f", "ceph-osd")
	res := &pb.ListCephesResponse{}
	res.SetCephes(toProtoFacilityInfos(cephes))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ListKubernetes(ctx context.Context, req *connect.Request[pb.ListKubernetesRequest]) (*connect.Response[pb.ListKubernetesResponse], error) {
	kubernetes, err := a.svc.ListKubernetes(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.ListKubernetesResponse{}
	res.SetKubernetes(toProtoFacilityInfos(kubernetes))
	return connect.NewResponse(res), nil
}

func toProtoErrors(es []model.Error) []*pb.Error {
	ret := []*pb.Error{}
	for i := range es {
		ret = append(ret, toProtoError(&es[i]))
	}
	return ret
}

func toProtoError(e *model.Error) *pb.Error {
	ret := &pb.Error{}
	ret.SetCode(e.Code)
	ret.SetLevel(pb.ErrorLevel(e.Level))
	ret.SetMessage(e.Message)
	ret.SetDetails(e.Details)
	ret.SetUrl(e.URL)
	return ret
}

func toProtoFacilityInfos(fis []model.FacilityInfo) []*pb.Facility_Info {
	ret := []*pb.Facility_Info{}
	for i := range fis {
		ret = append(ret, toProtoFacilityInfo(&fis[i]))
	}
	return ret
}

func toProtoFacilityInfo(fi *model.FacilityInfo) *pb.Facility_Info {
	ret := &pb.Facility_Info{}
	ret.SetScopeUuid(fi.ScopeUUID)
	ret.SetScopeName(fi.ScopeName)
	ret.SetFacilityName(fi.FacilityName)
	return ret
}
