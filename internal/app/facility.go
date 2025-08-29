package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/facility/v1"
	"github.com/otterscale/otterscale/api/facility/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type FacilityService struct {
	pbconnect.UnimplementedFacilityServiceHandler

	uc *core.FacilityUseCase
}

func NewFacilityService(uc *core.FacilityUseCase) *FacilityService {
	return &FacilityService{uc: uc}
}

var _ pbconnect.FacilityServiceHandler = (*FacilityService)(nil)

func (s *FacilityService) ListFacilities(ctx context.Context, req *connect.Request[pb.ListFacilitiesRequest]) (*connect.Response[pb.ListFacilitiesResponse], error) {
	facilities, err := s.uc.ListFacilities(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	machineMap, err := s.uc.JujuToMAASMachineMap(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListFacilitiesResponse{}
	resp.SetFacilities(toProtoFacilities(facilities, machineMap))
	return connect.NewResponse(resp), nil
}

func (s *FacilityService) GetFacility(ctx context.Context, req *connect.Request[pb.GetFacilityRequest]) (*connect.Response[pb.Facility], error) {
	facility, err := s.uc.GetFacility(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	machineMap, err := s.uc.JujuToMAASMachineMap(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := toProtoFacility(facility, machineMap)
	return connect.NewResponse(resp), nil
}

func (s *FacilityService) CreateFacility(ctx context.Context, req *connect.Request[pb.CreateFacilityRequest]) (*connect.Response[pb.Facility], error) {
	facility, err := s.uc.CreateFacility(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName(), req.Msg.GetConfigYaml(), req.Msg.GetCharmName(), req.Msg.GetChannel(), int(req.Msg.GetRevision()), int(req.Msg.GetNumber()), toModelPlacements(req.Msg.GetPlacements()), toModelConstraint(req.Msg.GetConstraint()), req.Msg.GetTrust())
	if err != nil {
		return nil, err
	}
	machineMap, err := s.uc.JujuToMAASMachineMap(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := toProtoFacility(facility, machineMap)
	return connect.NewResponse(resp), nil
}

func (s *FacilityService) UpdateFacility(ctx context.Context, req *connect.Request[pb.UpdateFacilityRequest]) (*connect.Response[pb.Facility], error) {
	facility, err := s.uc.UpdateFacility(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName(), req.Msg.GetConfigYaml())
	if err != nil {
		return nil, err
	}
	machineMap, err := s.uc.JujuToMAASMachineMap(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := toProtoFacility(facility, machineMap)
	return connect.NewResponse(resp), nil
}

func (s *FacilityService) DeleteFacility(ctx context.Context, req *connect.Request[pb.DeleteFacilityRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteFacility(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName(), req.Msg.GetDestroyStorage(), req.Msg.GetForce()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *FacilityService) ExposeFacility(ctx context.Context, req *connect.Request[pb.ExposeFacilityRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.ExposeFacility(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *FacilityService) AddFacilityUnits(ctx context.Context, req *connect.Request[pb.AddFacilityUnitsRequest]) (*connect.Response[pb.AddFacilityUnitsResponse], error) {
	units, err := s.uc.AddFacilityUnits(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName(), int(req.Msg.GetNumber()), toModelPlacements(req.Msg.GetPlacements()))
	if err != nil {
		return nil, err
	}
	resp := &pb.AddFacilityUnitsResponse{}
	resp.SetUnits(units)
	return connect.NewResponse(resp), nil
}

func (s *FacilityService) ListActions(ctx context.Context, req *connect.Request[pb.ListActionsRequest]) (*connect.Response[pb.ListActionsResponse], error) {
	actions, err := s.uc.ListActions(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListActionsResponse{}
	resp.SetActions(toProtoActions(actions))
	return connect.NewResponse(resp), nil
}

func (s *FacilityService) ListCharms(ctx context.Context, req *connect.Request[pb.ListCharmsRequest]) (*connect.Response[pb.ListCharmsResponse], error) {
	charms, err := s.uc.ListCharms(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListCharmsResponse{}
	resp.SetCharms(toProtoCharms(charms))
	return connect.NewResponse(resp), nil
}

func (s *FacilityService) GetCharm(ctx context.Context, req *connect.Request[pb.GetCharmRequest]) (*connect.Response[pb.Facility_Charm], error) {
	charm, err := s.uc.GetCharm(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoCharm(charm)
	return connect.NewResponse(resp), nil
}

func (s *FacilityService) ListCharmArtifacts(ctx context.Context, req *connect.Request[pb.ListCharmArtifactsRequest]) (*connect.Response[pb.ListCharmArtifactsResponse], error) {
	artifacts, err := s.uc.ListArtifacts(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListCharmArtifactsResponse{}
	resp.SetArtifacts(toProtoCharmArtifacts(artifacts))
	return connect.NewResponse(resp), nil
}

func toProtoFacilityStatus(s *core.DetailedStatus) *pb.Facility_Status {
	ret := &pb.Facility_Status{}
	ret.SetState(s.Status)
	ret.SetDetails(s.Info)
	since := s.Since
	if since != nil {
		ret.SetCreatedAt(timestamppb.New(*since))
	}
	return ret
}

func toProtoFacilityUnits(usm map[string]core.UnitStatus, machineMap map[string]string) []*pb.Facility_Unit {
	ret := []*pb.Facility_Unit{}
	for name := range usm {
		status := usm[name]
		ret = append(ret, toProtoFacilityUnit(name, &status, machineMap))
	}
	return ret
}

func toProtoFacilityUnit(name string, s *core.UnitStatus, machineMap map[string]string) *pb.Facility_Unit {
	ret := &pb.Facility_Unit{}
	ret.SetName(name)
	ret.SetAgentStatus(toProtoFacilityStatus(&s.AgentStatus))
	ret.SetWorkloadStatus(toProtoFacilityStatus(&s.WorkloadStatus))
	ret.SetLeader(s.Leader)
	ret.SetMachineId(machineMap[s.Machine])
	ret.SetIpAddress(s.Address + s.PublicAddress)
	ret.SetPorts(s.OpenedPorts)
	ret.SetCharmName(s.Charm)
	ret.SetVersion(s.WorkloadVersion)
	ret.SetSubordinates(toProtoFacilityUnits(s.Subordinates, machineMap))
	return ret
}

func toProtoFacilities(fs []core.Facility, machineMap map[string]string) []*pb.Facility {
	ret := []*pb.Facility{}
	for i := range fs {
		ret = append(ret, toProtoFacility(&fs[i], machineMap))
	}
	return ret
}

func toProtoFacility(f *core.Facility, machineMap map[string]string) *pb.Facility {
	ret := &pb.Facility{}
	ret.SetName(f.Name)
	ret.SetStatus(toProtoFacilityStatus(&f.Status.Status))
	ret.SetCharmName(f.Status.Charm)
	ret.SetVersion(f.Status.WorkloadVersion)
	ret.SetRevision(int64(f.Status.CharmRev))
	ret.SetChannel(f.Status.CharmChannel)
	ret.SetUnits(toProtoFacilityUnits(f.Status.Units, machineMap))
	if f.Metadata != nil {
		ret.SetMetadata(toProtoFacilityMetadata(f.Metadata))
	}
	return ret
}

func toProtoFacilityMetadata(md *core.FacilityMetadata) *pb.Facility_Charm_Metadata {
	ret := &pb.Facility_Charm_Metadata{}
	ret.SetConfigYaml(md.ConfigYAML)
	return ret
}

func toProtoActions(as []core.Action) []*pb.Facility_Action {
	ret := []*pb.Facility_Action{}
	for i := range as {
		ret = append(ret, toProtoAction(&as[i]))
	}
	return ret
}

func toProtoAction(a *core.Action) *pb.Facility_Action {
	ret := &pb.Facility_Action{}
	ret.SetName(a.Name)
	ret.SetDescription(a.Spec.Description)
	return ret
}

func toProtoCharms(cs []core.Charm) []*pb.Facility_Charm {
	ret := []*pb.Facility_Charm{}
	for i := range cs {
		ret = append(ret, toProtoCharm(&cs[i]))
	}
	return ret
}

func toProtoCharm(c *core.Charm) *pb.Facility_Charm {
	categories := []string{}
	for _, ca := range c.Result.Categories {
		categories = append(categories, ca.Name)
	}
	icon := ""
	for _, m := range c.Result.Media {
		icon = m.URL
		break
	}
	ret := &pb.Facility_Charm{}
	ret.SetId(c.ID)
	ret.SetType(c.Type)
	ret.SetName(c.Name)
	ret.SetVerified(false) // TODO: custom
	ret.SetTitle(c.Result.Title)
	ret.SetSummary(c.Result.Summary)
	ret.SetIcon(icon)
	ret.SetDescription(c.Result.Description)
	ret.SetCategories(categories)
	ret.SetDeployableOn(c.Result.DeployableOn)
	ret.SetPublisher(c.Result.Publisher.DisplayName)
	ret.SetLicense(c.Result.License)
	ret.SetStoreUrl(c.Result.StoreURL)
	ret.SetWebsite(c.Result.Website)
	ret.SetDefaultArtifact(toProtoCharmArtifact(&c.DefaultArtifact))
	return ret
}

func toProtoCharmBases(bs []core.CharmBase) []*pb.Facility_Charm_Base {
	ret := []*pb.Facility_Charm_Base{}
	for i := range bs {
		ret = append(ret, toProtoCharmBase(&bs[i]))
	}
	return ret
}

func toProtoCharmBase(b *core.CharmBase) *pb.Facility_Charm_Base {
	ret := &pb.Facility_Charm_Base{}
	ret.SetName(b.Name)
	ret.SetChannel(b.Channel)
	ret.SetArchitecture(b.Architecture)
	return ret
}

func toProtoCharmArtifacts(as []core.CharmArtifact) []*pb.Facility_Charm_Artifact {
	ret := []*pb.Facility_Charm_Artifact{}
	for i := range as {
		ret = append(ret, toProtoCharmArtifact(&as[i]))
	}
	return ret
}

func toProtoCharmArtifact(r *core.CharmArtifact) *pb.Facility_Charm_Artifact {
	ret := &pb.Facility_Charm_Artifact{}
	ret.SetChannel(r.Channel.Name)
	ret.SetRevision(int64(r.Revision.Revision))
	ret.SetVersion(r.Revision.Version)
	ret.SetBases(toProtoCharmBases(r.Revision.Bases))
	ret.SetCreatedAt(timestamppb.New(r.Channel.ReleasedAt))
	return ret
}

func toModelPlacements(ps []*pb.Placement) []core.MachinePlacement {
	ret := []core.MachinePlacement{}
	for i := range ps {
		ret = append(ret, *toModelPlacement(ps[i]))
	}
	return ret
}

func toModelPlacement(p *pb.Placement) *core.MachinePlacement {
	return &core.MachinePlacement{
		LXD:       p.GetLxd(),
		KVM:       p.GetKvm(),
		Machine:   p.GetMachine(),
		MachineID: p.GetMachineId(),
	}
}

func toModelConstraint(c *pb.Constraint) *core.MachineConstraint {
	return &core.MachineConstraint{
		Architecture: c.GetArchitecture(),
		CPUCores:     c.GetCpuCores(),
		MemoryMB:     c.GetMemoryMb(),
		Tags:         c.GetTags(),
	}
}
