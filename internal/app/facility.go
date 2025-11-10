package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/facility/v1"
	"github.com/otterscale/otterscale/api/facility/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type FacilityService struct {
	pbconnect.UnimplementedFacilityServiceHandler

	facility *core.FacilityUseCase
}

func NewFacilityService(facility *core.FacilityUseCase) *FacilityService {
	return &FacilityService{
		facility: facility,
	}
}

var _ pbconnect.FacilityServiceHandler = (*FacilityService)(nil)

func (s *FacilityService) ListFacilities(ctx context.Context, req *pb.ListFacilitiesRequest) (*pb.ListFacilitiesResponse, error) {
	facilities, err := s.facility.ListFacilities(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}
	machineMap, err := s.facility.JujuToMAASMachineMap(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListFacilitiesResponse{}
	resp.SetFacilities(toProtoFacilities(facilities, machineMap))
	return resp, nil
}

func (s *FacilityService) GetFacility(ctx context.Context, req *pb.GetFacilityRequest) (*pb.Facility, error) {
	facility, err := s.facility.GetFacility(ctx, req.GetScope(), req.GetName())
	if err != nil {
		return nil, err
	}
	machineMap, err := s.facility.JujuToMAASMachineMap(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}
	resp := toProtoFacility(facility, machineMap)
	return resp, nil
}

func (s *FacilityService) CreateFacility(ctx context.Context, req *pb.CreateFacilityRequest) (*pb.Facility, error) {
	facility, err := s.facility.CreateFacility(ctx, req.GetScope(), req.GetName(), req.GetConfigYaml(), req.GetCharmName(), req.GetChannel(), int(req.GetRevision()), int(req.GetNumber()), toModelPlacements(req.GetPlacements()), toModelConstraint(req.GetConstraint()), req.GetTrust())
	if err != nil {
		return nil, err
	}
	machineMap, err := s.facility.JujuToMAASMachineMap(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}
	resp := toProtoFacility(facility, machineMap)
	return resp, nil
}

func (s *FacilityService) UpdateFacility(ctx context.Context, req *pb.UpdateFacilityRequest) (*pb.Facility, error) {
	facility, err := s.facility.UpdateFacility(ctx, req.GetScope(), req.GetName(), req.GetConfigYaml())
	if err != nil {
		return nil, err
	}
	machineMap, err := s.facility.JujuToMAASMachineMap(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}
	resp := toProtoFacility(facility, machineMap)
	return resp, nil
}

func (s *FacilityService) DeleteFacility(ctx context.Context, req *pb.DeleteFacilityRequest) (*emptypb.Empty, error) {
	if err := s.facility.DeleteFacility(ctx, req.GetScope(), req.GetName(), req.GetDestroyStorage(), req.GetForce()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *FacilityService) ExposeFacility(ctx context.Context, req *pb.ExposeFacilityRequest) (*emptypb.Empty, error) {
	if err := s.facility.ExposeFacility(ctx, req.GetScope(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *FacilityService) AddFacilityUnits(ctx context.Context, req *pb.AddFacilityUnitsRequest) (*pb.AddFacilityUnitsResponse, error) {
	unitNames, err := s.facility.AddFacilityUnits(ctx, req.GetScope(), req.GetName(), int(req.GetNumber()), toModelPlacements(req.GetPlacements()))
	if err != nil {
		return nil, err
	}
	resp := &pb.AddFacilityUnitsResponse{}
	resp.SetUnitNames(unitNames)
	return resp, nil
}

func (s *FacilityService) ResolveFacilityUnitErrors(ctx context.Context, req *pb.ResolveFacilityUnitErrorsRequest) (*emptypb.Empty, error) {
	if err := s.facility.ResolveFacilityUnitErrors(ctx, req.GetScope(), req.GetUnitName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *FacilityService) ListActions(ctx context.Context, req *pb.ListActionsRequest) (*pb.ListActionsResponse, error) {
	actions, err := s.facility.ListActions(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListActionsResponse{}
	resp.SetActions(toProtoActions(actions))
	return resp, nil
}

func (s *FacilityService) ListCharms(ctx context.Context, _ *pb.ListCharmsRequest) (*pb.ListCharmsResponse, error) {
	charms, err := s.facility.ListCharms(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListCharmsResponse{}
	resp.SetCharms(toProtoCharms(charms))
	return resp, nil
}

func (s *FacilityService) GetCharm(ctx context.Context, req *pb.GetCharmRequest) (*pb.Facility_Charm, error) {
	charm, err := s.facility.GetCharm(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoCharm(charm)
	return resp, nil
}

func (s *FacilityService) ListCharmArtifacts(ctx context.Context, req *pb.ListCharmArtifactsRequest) (*pb.ListCharmArtifactsResponse, error) {
	artifacts, err := s.facility.ListArtifacts(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListCharmArtifactsResponse{}
	resp.SetArtifacts(toProtoCharmArtifacts(artifacts))
	return resp, nil
}

func toProtoFacilityStatus(s *core.FacilityDetailedStatus) *pb.Facility_Status {
	ret := &pb.Facility_Status{}
	ret.SetState(s.Status)
	ret.SetDetails(s.Info)
	since := s.Since
	if since != nil {
		ret.SetCreatedAt(timestamppb.New(*since))
	}
	return ret
}

func toProtoFacilityUnits(usm map[string]core.FacilityUnitStatus, machineMap map[string]core.FacilityMachineStatus) []*pb.Facility_Unit {
	ret := []*pb.Facility_Unit{}
	for name := range usm {
		status := usm[name]
		ret = append(ret, toProtoFacilityUnit(name, &status, machineMap))
	}
	return ret
}

func toProtoFacilityUnit(name string, s *core.FacilityUnitStatus, machineMap map[string]core.FacilityMachineStatus) *pb.Facility_Unit {
	ret := &pb.Facility_Unit{}
	ret.SetName(name)
	ret.SetAgentStatus(toProtoFacilityStatus(&s.AgentStatus))
	ret.SetWorkloadStatus(toProtoFacilityStatus(&s.WorkloadStatus))
	ret.SetLeader(s.Leader)
	ret.SetMachineId(string(machineMap[s.Machine].InstanceId))
	ret.SetHostname(machineMap[s.Machine].Hostname)
	ret.SetIpAddress(s.Address + s.PublicAddress)
	ret.SetPorts(s.OpenedPorts)
	ret.SetCharmName(s.Charm)
	ret.SetVersion(s.WorkloadVersion)
	ret.SetSubordinates(toProtoFacilityUnits(s.Subordinates, machineMap))
	return ret
}

func toProtoFacilities(fs []core.Facility, machineMap map[string]core.FacilityMachineStatus) []*pb.Facility {
	ret := []*pb.Facility{}
	for i := range fs {
		ret = append(ret, toProtoFacility(&fs[i], machineMap))
	}
	return ret
}

func toProtoFacility(f *core.Facility, machineMap map[string]core.FacilityMachineStatus) *pb.Facility {
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

func toProtoActions(as []core.FacilityAction) []*pb.Facility_Action {
	ret := []*pb.Facility_Action{}
	for i := range as {
		ret = append(ret, toProtoAction(&as[i]))
	}
	return ret
}

func toProtoAction(a *core.FacilityAction) *pb.Facility_Action {
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
