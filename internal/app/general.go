package app

import (
	"github.com/openhdc/otterscale/api/general/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/core"
)

type GeneralService struct {
	pbconnect.UnimplementedGeneralServiceHandler

	uc *core.GeneralUseCase
}

func NewGeneralService(uc *core.GeneralUseCase) *GeneralService {
	return &GeneralService{uc: uc}
}

var _ pbconnect.GeneralServiceHandler = (*GeneralService)(nil)

// func (s *GeneralService) VerifyEnvironment(ctx context.Context, req *connect.Request[pb.VerifyEnvironmentRequest]) (*connect.Response[pb.VerifyEnvironmentResponse], error) {
// 	es, err := s.uc.VerifyEnvironment(ctx, req.Msg.GetScopeUuid())
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := &pb.VerifyEnvironmentResponse{}
// 	resp.SetErrors(toProtoErrors(es))
// 	return connect.NewResponse(resp), nil
// }

// func (s *GeneralService) ListCephes(ctx context.Context, req *connect.Request[pb.ListCephesRequest]) (*connect.Response[pb.ListCephesResponse], error) {
// 	cephes, err := s.uc.ListCephes(ctx, req.Msg.GetScopeUuid())
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := &pb.ListCephesResponse{}
// 	resp.SetCephes(toProtoFacilityInfos(cephes))
// 	return connect.NewResponse(resp), nil
// }

// func (s *GeneralService) CreateCeph(ctx context.Context, req *connect.Request[pb.CreateCephRequest]) (*connect.Response[pb.Facility_Info], error) {
// 	ceph, err := s.uc.CreateCeph(ctx, req.Msg.GetScopeUuid(), req.Msg.GetMachineId(), req.Msg.GetPrefixName(), req.Msg.GetOsdDevices(), req.Msg.GetDevelopment())
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := toProtoFacilityInfo(ceph)
// 	return connect.NewResponse(resp), nil
// }

// func (s *GeneralService) AddCephUnits(ctx context.Context, req *connect.Request[pb.AddCephUnitsRequest]) (*connect.Response[emptypb.Empty], error) {
// 	if err := s.uc.AddCephUnits(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), int(req.Msg.GetNumber()), req.Msg.GetMachineIds()); err != nil {
// 		return nil, err
// 	}
// 	resp := &emptypb.Empty{}
// 	return connect.NewResponse(resp), nil
// }

// func (s *GeneralService) ListKuberneteses(ctx context.Context, req *connect.Request[pb.ListKubernetesesRequest]) (*connect.Response[pb.ListKubernetesesResponse], error) {
// 	kuberneteses, err := s.uc.ListKuberneteses(ctx, req.Msg.GetScopeUuid())
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := &pb.ListKubernetesesResponse{}
// 	resp.SetKuberneteses(toProtoFacilityInfos(kuberneteses))
// 	return connect.NewResponse(resp), nil
// }

// func (s *GeneralService) CreateKubernetes(ctx context.Context, req *connect.Request[pb.CreateKubernetesRequest]) (*connect.Response[pb.Facility_Info], error) {
// 	kubernetes, err := s.uc.CreateKubernetes(ctx, req.Msg.GetScopeUuid(), req.Msg.GetMachineId(), req.Msg.GetPrefixName(), req.Msg.GetVirtualIps(), req.Msg.GetCalicoCidr())
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := toProtoFacilityInfo(kubernetes)
// 	return connect.NewResponse(resp), nil
// }

// func (s *GeneralService) AddKubernetesUnits(ctx context.Context, req *connect.Request[pb.AddKubernetesUnitsRequest]) (*connect.Response[emptypb.Empty], error) {
// 	if err := s.uc.AddKubernetesUnits(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), int(req.Msg.GetNumber()), req.Msg.GetMachineIds(), req.Msg.GetForce()); err != nil {
// 		return nil, err
// 	}
// 	resp := &emptypb.Empty{}
// 	return connect.NewResponse(resp), nil
// }

// func (s *GeneralService) SetCephCSI(ctx context.Context, req *connect.Request[pb.SetCephCSIRequest]) (*connect.Response[emptypb.Empty], error) {
// 	k := req.Msg.GetKubernetes()
// 	if k == nil {
// 		return nil, status.Error(codes.InvalidArgument, "kubernetes is empty")
// 	}
// 	c := req.Msg.GetCeph()
// 	if c == nil {
// 		return nil, status.Error(codes.InvalidArgument, "ceph is empty")
// 	}
// 	kubernetes := &core.FacilityInfo{
// 		ScopeUUID:    k.GetScopeUuid(),
// 		ScopeName:    k.GetScopeName(),
// 		FacilityName: k.GetFacilityName(),
// 	}
// 	ceph := &core.FacilityInfo{
// 		ScopeUUID:    c.GetScopeUuid(),
// 		ScopeName:    c.GetScopeName(),
// 		FacilityName: c.GetFacilityName(),
// 	}
// 	if err := s.uc.SetCephCSI(ctx, kubernetes, ceph, req.Msg.GetPrefix(), req.Msg.GetDevelopment()); err != nil {
// 		return nil, err
// 	}
// 	resp := &emptypb.Empty{}
// 	return connect.NewResponse(resp), nil
// }

// func toProtoErrors(es []core.Error) []*pb.Error {
// 	ret := []*pb.Error{}
// 	for i := range es {
// 		ret = append(ret, toProtoError(&es[i]))
// 	}
// 	return ret
// }

// func toProtoError(e *core.Error) *pb.Error {
// 	ret := &pb.Error{}
// 	ret.SetCode(e.Code)
// 	ret.SetLevel(pb.ErrorLevel(e.Level)) //nolint:gosec
// 	ret.SetMessage(e.Message)
// 	ret.SetDetails(e.Details)
// 	ret.SetUrl(e.URL)
// 	return ret
// }

// func toProtoFacilityInfos(fis []core.FacilityInfo) []*pb.Facility_Info {
// 	ret := []*pb.Facility_Info{}
// 	for i := range fis {
// 		ret = append(ret, toProtoFacilityInfo(&fis[i]))
// 	}
// 	return ret
// }

// func toProtoFacilityInfo(fi *core.FacilityInfo) *pb.Facility_Info {
// 	ret := &pb.Facility_Info{}
// 	ret.SetScopeUuid(fi.ScopeUUID)
// 	ret.SetScopeName(fi.ScopeName)
// 	ret.SetFacilityName(fi.FacilityName)
// 	return ret
// }
