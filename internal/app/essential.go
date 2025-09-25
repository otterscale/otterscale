package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/otterscale/otterscale/api/essential/v1"
	"github.com/otterscale/otterscale/api/essential/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type EssentialService struct {
	pbconnect.UnimplementedEssentialServiceHandler

	uc *core.EssentialUseCase
}

func NewEssentialService(uc *core.EssentialUseCase) *EssentialService {
	return &EssentialService{uc: uc}
}

var _ pbconnect.EssentialServiceHandler = (*EssentialService)(nil)

func (s *EssentialService) GetCapabilities(ctx context.Context, req *connect.Request[pb.GetCapabilitiesRequest]) (*connect.Response[pb.GetCapabilitiesResponse], error) {
	capabilities, err := s.uc.GetCapabilities(ctx, req.Msg.GetLanguage())
	if err != nil {
		return nil, err
	}
	
	resp := &pb.GetCapabilitiesResponse{}
	resp.SetPlatformName(capabilities.PlatformName)
	resp.SetPlatformDescription(capabilities.PlatformDescription)
	resp.SetCapabilities(toProtoCapabilities(capabilities.Capabilities))
	resp.SetUseCases(capabilities.UseCases)
	resp.SetDocumentationUrl(capabilities.DocumentationURL)
	resp.SetVersion(capabilities.Version)
	
	return connect.NewResponse(resp), nil
}

func (s *EssentialService) IsMachineDeployed(ctx context.Context, req *connect.Request[pb.IsMachineDeployedRequest]) (*connect.Response[pb.IsMachineDeployedResponse], error) {
	message, deployed, err := s.uc.IsMachineDeployed(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.IsMachineDeployedResponse{}
	resp.SetMessage(message)
	resp.SetDeployed(deployed)
	return connect.NewResponse(resp), nil
}

func (s *EssentialService) ListStatuses(ctx context.Context, req *connect.Request[pb.ListStatusesRequest]) (*connect.Response[pb.ListStatusesResponse], error) {
	statuses, err := s.uc.ListStatuses(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListStatusesResponse{}
	resp.SetStatuses(toProtoStatuses(statuses))
	return connect.NewResponse(resp), nil
}

func (s *EssentialService) ListEssentials(ctx context.Context, req *connect.Request[pb.ListEssentialsRequest]) (*connect.Response[pb.ListEssentialsResponse], error) {
	essentials, err := s.uc.ListEssentials(ctx, int32(req.Msg.GetType()), req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListEssentialsResponse{}
	resp.SetEssentials(toProtoEssentials(essentials))
	return connect.NewResponse(resp), nil
}

func (s *EssentialService) CreateSingleNode(ctx context.Context, req *connect.Request[pb.CreateSingleNodeRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.CreateSingleNode(ctx,
		req.Msg.GetScopeUuid(),
		req.Msg.GetMachineId(),
		req.Msg.GetPrefixName(),
		req.Msg.GetVirtualIps(),
		req.Msg.GetCalicoCidr(),
		req.Msg.GetOsdDevices(),
	); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func toProtoStatuses(ess []core.EssentialStatus) []*pb.Status {
	ret := []*pb.Status{}
	for i := range ess {
		ret = append(ret, toProtoStatus(&ess[i]))
	}
	return ret
}

func toProtoStatus(es *core.EssentialStatus) *pb.Status {
	ret := &pb.Status{}
	ret.SetLevel(pb.Status_Level(es.Level))
	ret.SetMessage(es.Message)
	ret.SetDetails(es.Details)
	return ret
}

func toProtoEssentials(es []core.Essential) []*pb.Essential {
	ret := []*pb.Essential{}
	for i := range es {
		ret = append(ret, toProtoEssential(&es[i]))
	}
	return ret
}

func toProtoEssential(e *core.Essential) *pb.Essential {
	ret := &pb.Essential{}
	ret.SetType(pb.Essential_Type(e.Type))
	ret.SetName(e.Name)
	ret.SetScopeUuid(e.ScopeUUID)
	ret.SetScopeName(e.ScopeName)
	ret.SetUnits(toProtoEssentialUnits(e.Units))
	return ret
}

func toProtoEssentialUnits(eus []core.EssentialUnit) []*pb.Essential_Unit {
	ret := []*pb.Essential_Unit{}
	for i := range eus {
		ret = append(ret, toProtoEssentialUnit(&eus[i]))
	}
	return ret
}

func toProtoEssentialUnit(eu *core.EssentialUnit) *pb.Essential_Unit {
	ret := &pb.Essential_Unit{}
	ret.SetName(eu.Name)
	ret.SetDirective(eu.Directive)
	return ret
}

func toProtoCapabilities(caps []core.Capability) []*pb.GetCapabilitiesResponse_Capability {
	ret := []*pb.GetCapabilitiesResponse_Capability{}
	for i := range caps {
		ret = append(ret, toProtoCapability(&caps[i]))
	}
	return ret
}

func toProtoCapability(cap *core.Capability) *pb.GetCapabilitiesResponse_Capability {
	ret := &pb.GetCapabilitiesResponse_Capability{}
	ret.SetCategory(cap.Category)
	ret.SetName(cap.Name)
	ret.SetDescription(cap.Description)
	ret.SetFeatures(cap.Features)
	ret.SetAvailable(cap.Available)
	return ret
}
