package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

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
	cephes, err := a.svc.ListCephes(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	res := &pb.ListCephesResponse{}
	res.SetCephes(toProtoFacilityInfos(cephes))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) CreateCeph(ctx context.Context, req *connect.Request[pb.CreateCephRequest]) (*connect.Response[pb.Facility_Info], error) {
	ceph, err := a.svc.CreateCeph(ctx, req.Msg.GetScopeUuid(), req.Msg.GetMachineId(), req.Msg.GetPrefixName(), req.Msg.GetOsdDevices(), req.Msg.GetDevelopment())
	if err != nil {
		return nil, err
	}
	res := toProtoFacilityInfo(ceph)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) AddCephUnits(ctx context.Context, req *connect.Request[pb.AddCephUnitsRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.AddCephUnits(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), int(req.Msg.GetNumber()), req.Msg.GetMachineIds()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ListKuberneteses(ctx context.Context, req *connect.Request[pb.ListKubernetesesRequest]) (*connect.Response[pb.ListKubernetesesResponse], error) {
	kuberneteses, err := a.svc.ListKuberneteses(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	res := &pb.ListKubernetesesResponse{}
	res.SetKuberneteses(toProtoFacilityInfos(kuberneteses))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) CreateKubernetes(ctx context.Context, req *connect.Request[pb.CreateKubernetesRequest]) (*connect.Response[pb.Facility_Info], error) {
	kubernetes, err := a.svc.CreateKubernetes(ctx, req.Msg.GetScopeUuid(), req.Msg.GetMachineId(), req.Msg.GetPrefixName(), req.Msg.GetVirtualIps(), req.Msg.GetCalicoCidr())
	if err != nil {
		return nil, err
	}
	res := toProtoFacilityInfo(kubernetes)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) AddKubernetesUnits(ctx context.Context, req *connect.Request[pb.AddKubernetesUnitsRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.AddKubernetesUnits(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), int(req.Msg.GetNumber()), req.Msg.GetMachineIds(), req.Msg.GetForce()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
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
	ret.SetLevel(pb.ErrorLevel(e.Level)) //nolint:gosec
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
