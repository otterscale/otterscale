package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/facility/v1"
	"github.com/otterscale/otterscale/api/facility/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/facility"
	"github.com/otterscale/otterscale/internal/core/facility/action"
	"github.com/otterscale/otterscale/internal/core/facility/charm"
	"github.com/otterscale/otterscale/internal/core/machine"
)

type FacilityService struct {
	pbconnect.UnimplementedFacilityServiceHandler

	facility *facility.UseCase
	action   *action.UseCase
	charm    *charm.UseCase
}

func NewFacilityService(facility *facility.UseCase, action *action.UseCase, charm *charm.UseCase) *FacilityService {
	return &FacilityService{
		facility: facility,
		action:   action,
		charm:    charm,
	}
}

var _ pbconnect.FacilityServiceHandler = (*FacilityService)(nil)

func (s *FacilityService) ListFacilities(ctx context.Context, req *pb.ListFacilitiesRequest) (*pb.ListFacilitiesResponse, error) {
	facilities, err := s.facility.ListFacilities(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	machineMap, err := s.facility.JujuIDMachineMap(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListFacilitiesResponse{}
	resp.SetFacilities(toProtoFacilities(facilities, machineMap))
	return resp, nil
}

func (s *FacilityService) ResolveFacilityUnitErrors(ctx context.Context, req *pb.ResolveFacilityUnitErrorsRequest) (*emptypb.Empty, error) {
	if err := s.facility.ResolveFacilityUnitErrors(ctx, req.GetScope(), req.GetUnitName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func toProtoFacilityStatus(s *facility.DetailedStatus) *pb.Facility_Status {
	ret := &pb.Facility_Status{}
	ret.SetState(s.Status)
	ret.SetDetails(s.Info)

	since := s.Since
	if since != nil {
		ret.SetCreatedAt(timestamppb.New(*since))
	}

	return ret
}

func toProtoFacilityUnits(usm map[string]facility.UnitStatus, machineMap map[string]machine.Machine) []*pb.Facility_Unit {
	ret := []*pb.Facility_Unit{}

	for name := range usm {
		status := usm[name]
		ret = append(ret, toProtoFacilityUnit(name, &status, machineMap))
	}

	return ret
}

func toProtoFacilityUnit(name string, s *facility.UnitStatus, machineMap map[string]machine.Machine) *pb.Facility_Unit {
	ret := &pb.Facility_Unit{}
	ret.SetName(name)
	ret.SetAgentStatus(toProtoFacilityStatus(&s.AgentStatus))
	ret.SetWorkloadStatus(toProtoFacilityStatus(&s.WorkloadStatus))
	ret.SetLeader(s.Leader)
	ret.SetMachineId(machineMap[s.Machine].SystemID)
	ret.SetHostname(machineMap[s.Machine].Hostname)
	ret.SetIpAddress(s.Address + s.PublicAddress)
	ret.SetPorts(s.OpenedPorts)
	ret.SetCharmName(s.Charm)
	ret.SetVersion(s.WorkloadVersion)
	ret.SetSubordinates(toProtoFacilityUnits(s.Subordinates, machineMap))
	return ret
}

func toProtoFacilities(fs []facility.Facility, machineMap map[string]machine.Machine) []*pb.Facility {
	ret := []*pb.Facility{}

	for i := range fs {
		ret = append(ret, toProtoFacility(&fs[i], machineMap))
	}

	return ret
}

func toProtoFacility(f *facility.Facility, machineMap map[string]machine.Machine) *pb.Facility {
	ret := &pb.Facility{}
	ret.SetName(f.Name)
	ret.SetStatus(toProtoFacilityStatus(&f.Status.Status))
	ret.SetCharmName(f.Status.Charm)
	ret.SetVersion(f.Status.WorkloadVersion)
	ret.SetRevision(int64(f.Status.CharmRev))
	ret.SetChannel(f.Status.CharmChannel)
	ret.SetUnits(toProtoFacilityUnits(f.Status.Units, machineMap))
	return ret
}
