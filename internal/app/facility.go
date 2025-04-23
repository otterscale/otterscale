package app

import (
	"context"
	"os"
	"slices"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/openhdc/openhdc/api/nexus/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
	"github.com/openhdc/openhdc/internal/env"
)

func (a *NexusApp) ListFacilities(ctx context.Context, req *connect.Request[pb.ListFacilitiesRequest]) (*connect.Response[pb.ListFacilitiesResponse], error) {
	fs, err := a.svc.ListFacilities(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	machineMap, err := a.svc.JujuToMAASMachineMap(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	res := &pb.ListFacilitiesResponse{}
	res.SetFacilities(toProtoFacilities(fs, machineMap))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) GetFacility(ctx context.Context, req *connect.Request[pb.GetFacilityRequest]) (*connect.Response[pb.Facility], error) {
	f, err := a.svc.GetFacility(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	machineMap, err := a.svc.JujuToMAASMachineMap(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	res := toProtoFacility(f, machineMap)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) CreateFacility(ctx context.Context, req *connect.Request[pb.CreateFacilityRequest]) (*connect.Response[pb.Facility], error) {
	f, err := a.svc.CreateFacility(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName(), req.Msg.GetConfigYaml(), req.Msg.GetCharmName(), req.Msg.GetChannel(), int(req.Msg.GetRevision()), int(req.Msg.GetNumber()), toModelPlacements(req.Msg.GetPlacements()), toModelConstraint(req.Msg.GetConstraint()), req.Msg.GetTrust())
	if err != nil {
		return nil, err
	}
	machineMap, err := a.svc.JujuToMAASMachineMap(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	res := toProtoFacility(f, machineMap)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) UpdateFacility(ctx context.Context, req *connect.Request[pb.UpdateFacilityRequest]) (*connect.Response[pb.Facility], error) {
	f, err := a.svc.UpdateFacility(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName(), req.Msg.GetConfigYaml())
	if err != nil {
		return nil, err
	}
	machineMap, err := a.svc.JujuToMAASMachineMap(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	res := toProtoFacility(f, machineMap)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) DeleteFacility(ctx context.Context, req *connect.Request[pb.DeleteFacilityRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.DeleteFacility(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName(), req.Msg.GetDestroyStorage(), req.Msg.GetForce()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ExposeFacility(ctx context.Context, req *connect.Request[pb.ExposeFacilityRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.ExposeFacility(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) AddFacilityUnits(ctx context.Context, req *connect.Request[pb.AddFacilityUnitsRequest]) (*connect.Response[pb.AddFacilityUnitsResponse], error) {
	units, err := a.svc.AddFacilityUnits(ctx, req.Msg.GetScopeUuid(), req.Msg.GetName(), int(req.Msg.GetNumber()), toModelPlacements(req.Msg.GetPlacements()))
	if err != nil {
		return nil, err
	}
	res := &pb.AddFacilityUnitsResponse{}
	res.SetUnits(units)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ListActions(ctx context.Context, req *connect.Request[pb.ListActionsRequest]) (*connect.Response[pb.ListActionsResponse], error) {
	as, err := a.svc.ListActions(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	res := &pb.ListActionsResponse{}
	res.SetActions(toProtoActions(as))
	return connect.NewResponse(res), nil
}

// TODO: Unimplemented
// func (a *NexusApp) DoAction(ctx context.Context, req *connect.Request[pb.DoActionRequest]) (*connect.Response[emptypb.Empty], error) {
// 	res := &emptypb.Empty{}
// 	return connect.NewResponse(res), nil
// }

func (a *NexusApp) ListCharms(ctx context.Context, req *connect.Request[pb.ListCharmsRequest]) (*connect.Response[pb.ListCharmsResponse], error) {
	cs, err := a.svc.ListCharms(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.ListCharmsResponse{}
	res.SetCharms(toProtoCharms(cs))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) GetCharm(ctx context.Context, req *connect.Request[pb.GetCharmRequest]) (*connect.Response[pb.Facility_Charm], error) {
	c, err := a.svc.GetCharm(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	res := toProtoCharm(c)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) ListCharmArtifacts(ctx context.Context, req *connect.Request[pb.ListCharmArtifactsRequest]) (*connect.Response[pb.ListCharmArtifactsResponse], error) {
	as, err := a.svc.ListArtifacts(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	res := &pb.ListCharmArtifactsResponse{}
	res.SetArtifacts(toProtoCharmArtifacts(as))
	return connect.NewResponse(res), nil
}

func toProtoFacilityStatus(s *model.DetailedStatus) *pb.Facility_Status {
	ret := &pb.Facility_Status{}
	ret.SetState(s.Status)
	ret.SetDetails(s.Info)
	since := s.Since
	if since != nil {
		ret.SetCreatedAt(timestamppb.New(*since))
	}
	return ret
}

func toProtoFacilityUnits(usm map[string]model.UnitStatus, machineMap map[string]string) []*pb.Facility_Unit {
	ret := []*pb.Facility_Unit{}
	for name := range usm {
		status := usm[name]
		ret = append(ret, toProtoFacilityUnit(name, &status, machineMap))
	}
	return ret
}

func toProtoFacilityUnit(name string, s *model.UnitStatus, machineMap map[string]string) *pb.Facility_Unit {
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

func toProtoFacilities(fs []model.Facility, machineMap map[string]string) []*pb.Facility {
	ret := []*pb.Facility{}
	for i := range fs {
		ret = append(ret, toProtoFacility(&fs[i], machineMap))
	}
	return ret
}

func toProtoFacility(f *model.Facility, machineMap map[string]string) *pb.Facility {
	ret := &pb.Facility{}
	ret.SetName(f.Name)
	ret.SetStatus(toProtoFacilityStatus(&f.Status.Status))
	ret.SetCharmName(f.Status.Charm)
	ret.SetVersion(f.Status.WorkloadVersion)
	ret.SetRevision(int64(f.Status.CharmRev))
	ret.SetUnits(toProtoFacilityUnits(f.Status.Units, machineMap))
	if f.FacilityMetadata != nil {
		ret.SetMetadata(toProtoFacilityMetadata(f.FacilityMetadata))
	}
	return ret
}

func toProtoFacilityMetadata(md *model.FacilityMetadata) *pb.Facility_Charm_Metadata {
	ret := &pb.Facility_Charm_Metadata{}
	ret.SetConfigYaml(md.ConfigYAML)
	return ret
}

func toProtoActions(as []model.Action) []*pb.Facility_Action {
	ret := []*pb.Facility_Action{}
	for i := range as {
		ret = append(ret, toProtoAction(&as[i]))
	}
	return ret
}

func toProtoAction(a *model.Action) *pb.Facility_Action {
	ret := &pb.Facility_Action{}
	ret.SetName(a.Name)
	ret.SetDescription(a.Spec.Description)
	return ret
}

func toProtoCharms(cs []model.Charm) []*pb.Facility_Charm {
	ret := []*pb.Facility_Charm{}
	for i := range cs {
		ret = append(ret, toProtoCharm(&cs[i]))
	}
	return ret
}

func toProtoCharm(c *model.Charm) *pb.Facility_Charm {
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
	ret.SetVerified(isVerified(c.Result.Publisher.DisplayName))
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

func toProtoCharmBases(bs []model.CharmBase) []*pb.Facility_Charm_Base {
	ret := []*pb.Facility_Charm_Base{}
	for i := range bs {
		ret = append(ret, toProtoCharmBase(&bs[i]))
	}
	return ret
}

func toProtoCharmBase(b *model.CharmBase) *pb.Facility_Charm_Base {
	ret := &pb.Facility_Charm_Base{}
	ret.SetName(b.Name)
	ret.SetChannel(b.Channel)
	ret.SetArchitecture(b.Architecture)
	return ret
}

func toProtoCharmArtifacts(as []model.CharmArtifact) []*pb.Facility_Charm_Artifact {
	ret := []*pb.Facility_Charm_Artifact{}
	for i := range as {
		ret = append(ret, toProtoCharmArtifact(&as[i]))
	}
	return ret
}

func toProtoCharmArtifact(r *model.CharmArtifact) *pb.Facility_Charm_Artifact {
	ret := &pb.Facility_Charm_Artifact{}
	ret.SetChannel(r.Channel.Name)
	ret.SetRevision(int64(r.Revision.Revision))
	ret.SetVersion(r.Revision.Version)
	ret.SetBases(toProtoCharmBases(r.Revision.Bases))
	ret.SetCreatedAt(timestamppb.New(r.Channel.ReleasedAt))
	return ret
}

func isVerified(values ...string) bool {
	keyword := os.Getenv(env.OPENHDC_COMPANY_NAME)
	if keyword == "" {
		return false
	}
	return slices.Contains(values, keyword)
}
