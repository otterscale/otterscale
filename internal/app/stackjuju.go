package app

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	jujuyaml "gopkg.in/yaml.v2"

	v1 "github.com/openhdc/openhdc/api/stack/v1"
)

func (a *StackApp) AddMachines(ctx context.Context, req *connect.Request[v1.AddMachinesRequest]) (*connect.Response[v1.AddMachinesResponse], error) {
	uuid := req.Msg.GetModelUuid()

	factors := req.Msg.GetFactors()
	machineParams := make([]params.AddMachineParams, 0, len(factors))
	for i, f := range factors {
		machineParam := params.AddMachineParams{}
		if c := f.GetConstraint(); c != nil {
			machineParam.Constraints = buildConstraints(c)
		}
		if p := f.GetPlacement(); p != nil {
			placement, err := a.buildPlacement(ctx, uuid, p)
			if err != nil {
				return nil, err
			}
			machineParam.Placement = placement
		}
		machineParams[i] = machineParam
	}

	machines, err := a.svc.AddMachine(ctx, uuid, machineParams)
	if err != nil {
		return nil, err
	}

	res := &v1.AddMachinesResponse{}
	res.SetMachines(machines)
	return connect.NewResponse(res), nil
}

func (a *StackApp) ListModels(ctx context.Context, req *connect.Request[v1.ListModelsRequest]) (*connect.Response[v1.ListModelsResponse], error) {
	mds, err := a.svc.ListModels(ctx)
	if err != nil {
		return nil, err
	}
	res := &v1.ListModelsResponse{}
	res.SetModels(toModels(mds))
	return connect.NewResponse(res), nil
}

func (a *StackApp) CreateModel(ctx context.Context, req *connect.Request[v1.CreateModelRequest]) (*connect.Response[v1.Model], error) {
	mi, err := a.svc.CreateModel(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(modelInfoToModel(mi)), nil
}

func (a *StackApp) GetModelConfig(ctx context.Context, req *connect.Request[v1.GetModelConfigRequest]) (*connect.Response[v1.GetModelConfigResponse], error) {
	mc, err := a.svc.GetModelConfig(ctx, req.Msg.GetUuid())
	if err != nil {
		return nil, err
	}
	configYAML, err := jujuyaml.Marshal(mc)
	if err != nil {
		return nil, err
	}
	res := &v1.GetModelConfigResponse{}
	res.SetConfigYaml(string(configYAML))
	return connect.NewResponse(res), nil
}

func (a *StackApp) ListApplications(ctx context.Context, req *connect.Request[v1.ListApplicationsRequest]) (*connect.Response[v1.ListApplicationsResponse], error) {
	as, err := a.svc.ListApplications(ctx, req.Msg.GetModelUuid())
	if err != nil {
		return nil, err
	}
	ms, err := a.svc.ListJujuMachines(ctx, req.Msg.GetModelUuid())
	if err != nil {
		return nil, err
	}
	cs, err := a.svc.ListApplicationConfigs(ctx, req.Msg.GetModelUuid(), as)
	if err != nil {
		return nil, err
	}
	res := &v1.ListApplicationsResponse{}
	res.SetApplications(toApplications(as, ms, cs))
	return connect.NewResponse(res), nil
}

func (a *StackApp) GetApplication(ctx context.Context, req *connect.Request[v1.GetApplicationRequest]) (*connect.Response[v1.Application], error) {
	as, err := a.svc.ListApplications(ctx, req.Msg.GetModelUuid(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	ms, err := a.svc.ListJujuMachines(ctx, req.Msg.GetModelUuid())
	if err != nil {
		return nil, err
	}
	cs, err := a.svc.ListApplicationConfigs(ctx, req.Msg.GetModelUuid(), as)
	if err != nil {
		return nil, err
	}
	apps := toApplications(as, ms, cs)
	if len(apps) == 1 {
		return connect.NewResponse(apps[0]), nil
	}
	return nil, status.Error(codes.Internal, "something went wrong")
}

func (a *StackApp) CreateApplication(ctx context.Context, req *connect.Request[v1.CreateApplicationRequest]) (*connect.Response[v1.Application], error) {
	placements, err := a.buildPlacements(ctx, req.Msg.GetModelUuid(), req.Msg.GetPlacements())
	if err != nil {
		return nil, err
	}
	as, err := a.svc.CreateApplication(ctx, req.Msg.GetModelUuid(), req.Msg.GetCharmName(), req.Msg.GetName(), req.Msg.GetChannel(), int(req.Msg.GetRevision()), int(req.Msg.GetNumber()), req.Msg.GetOverrideConfig(), buildConstraints(req.Msg.GetConstraint()), placements, req.Msg.GetTrust())
	if err != nil {
		return nil, err
	}
	ms, err := a.svc.ListJujuMachines(ctx, req.Msg.GetModelUuid())
	if err != nil {
		return nil, err
	}
	cs, err := a.svc.ListApplicationConfigs(ctx, req.Msg.GetModelUuid(), as)
	if err != nil {
		return nil, err
	}
	apps := toApplications(as, ms, cs)
	if len(apps) == 1 {
		return connect.NewResponse(apps[0]), nil
	}
	return nil, status.Error(codes.Internal, "something went wrong")
}

func (a *StackApp) UpdateApplication(ctx context.Context, req *connect.Request[v1.UpdateApplicationRequest]) (*connect.Response[v1.Application], error) {
	as, err := a.svc.UpdateApplication(ctx, req.Msg.GetModelUuid(), req.Msg.GetName(), req.Msg.GetOverrideConfig())
	if err != nil {
		return nil, err
	}
	ms, err := a.svc.ListJujuMachines(ctx, req.Msg.GetModelUuid())
	if err != nil {
		return nil, err
	}
	cs, err := a.svc.ListApplicationConfigs(ctx, req.Msg.GetModelUuid(), as)
	if err != nil {
		return nil, err
	}
	apps := toApplications(as, ms, cs)
	if len(apps) == 1 {
		return connect.NewResponse(apps[0]), nil
	}
	return nil, status.Error(codes.Internal, "something went wrong")
}

func (a *StackApp) DeleteApplication(ctx context.Context, req *connect.Request[v1.DeleteApplicationRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.DeleteApplication(ctx, req.Msg.GetModelUuid(), req.Msg.GetName(), req.Msg.GetDestroyStorage(), req.Msg.GetForce()); err != nil {
		return nil, err
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (a *StackApp) ExposeApplication(ctx context.Context, req *connect.Request[v1.ExposeApplicationRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.ExposeApplication(ctx, req.Msg.GetModelUuid(), req.Msg.GetName(), nil); err != nil {
		return nil, err
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (a *StackApp) AddApplicationUnits(ctx context.Context, req *connect.Request[v1.AddApplicationUnitsRequest]) (*connect.Response[v1.AddApplicationUnitsResponse], error) {
	placements, err := a.buildPlacements(ctx, req.Msg.GetModelUuid(), req.Msg.GetPlacements())
	if err != nil {
		return nil, err
	}
	mss, err := a.svc.AddApplicationsUnits(ctx, req.Msg.GetModelUuid(), req.Msg.GetName(), int(req.Msg.GetNumber()), placements)
	if err != nil {
		return nil, err
	}

	res := &v1.AddApplicationUnitsResponse{}
	res.SetMachines(toAddApplicationUnitsMachines(mss))
	return connect.NewResponse(res), nil
}

func (a *StackApp) ListIntegrations(ctx context.Context, req *connect.Request[v1.ListIntegrationsRequest]) (*connect.Response[v1.ListIntegrationsResponse], error) {
	rss, err := a.svc.ListIntegrations(ctx, req.Msg.GetModelUuid())
	if err != nil {
		return nil, err
	}
	res := &v1.ListIntegrationsResponse{}
	res.SetIntegrations(toIntegrations(rss))
	return connect.NewResponse(res), nil
}

func (a *StackApp) CreateIntegration(ctx context.Context, req *connect.Request[v1.CreateIntegrationRequest]) (*connect.Response[v1.Integration], error) {
	arr, err := a.svc.CreateIntegration(ctx, req.Msg.GetModelUuid(), req.Msg.GetEndpoints())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resultsToIntegration(arr)), nil
}

func (a *StackApp) DeleteIntegration(ctx context.Context, req *connect.Request[v1.DeleteIntegrationRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.DeleteIntegration(ctx, req.Msg.GetModelUuid(), int(req.Msg.GetId())); err != nil {
		return nil, err
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (a *StackApp) ListActions(ctx context.Context, req *connect.Request[v1.ListActionsRequest]) (*connect.Response[v1.ListActionsResponse], error) {
	m, err := a.svc.ListActions(ctx, req.Msg.GetModelUuid(), req.Msg.GetApplicationName())
	if err != nil {
		return nil, err
	}
	ret := &v1.ListActionsResponse{}
	ret.SetActions(toActions(m))
	return connect.NewResponse(ret), nil
}

func (a *StackApp) RunAction(ctx context.Context, req *connect.Request[v1.RunActionRequest]) (*connect.Response[v1.Action], error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (a *StackApp) buildPlacements(ctx context.Context, uuid string, ps []*v1.Placement) ([]*instance.Placement, error) {
	ret := make([]*instance.Placement, len(ps))
	for i := range ps {
		p, err := a.buildPlacement(ctx, uuid, ps[i])
		if err != nil {
			return nil, err
		}
		ret[i] = p
	}
	return ret, nil
}

func (a *StackApp) buildPlacement(ctx context.Context, uuid string, p *v1.Placement) (*instance.Placement, error) {
	directive, err := a.svc.MAASToJujuMachineID(ctx, uuid, p.GetMachineSystemId())
	if err != nil {
		return nil, err
	}
	return &instance.Placement{
		Scope:     p.GetScope(),
		Directive: directive,
	}, nil
}

func buildConstraints(c *v1.Constraint) constraints.Value {
	arch := c.GetArchitecture()
	cpuCores := c.GetCpuCores()
	mem := c.GetMemoryMb()
	tags := c.GetTags()

	return constraints.Value{
		Arch:     &arch,
		CpuCores: &cpuCores,
		Mem:      &mem,
		Tags:     &tags,
	}
}

func toModels(umss []*base.UserModelSummary) []*v1.Model {
	ret := make([]*v1.Model, len(umss))
	for i := range umss {
		ret[i] = toModel(umss[i])
	}
	return ret
}

func toModel(m *base.UserModelSummary) *v1.Model {
	ret := &v1.Model{}
	ret.SetUuid(m.UUID)
	ret.SetName(m.Name)
	ret.SetLife(string(m.Life))
	ret.SetStatus(string(m.Status.Status))
	if m.UserLastConnection != nil {
		ret.SetUpdatedAt(timestamppb.New(*m.UserLastConnection))
	}
	for _, c := range m.Counts {
		switch c.Entity {
		case "machines":
			ret.SetMachineCount(c.Count)
		case "cores":
			ret.SetCoreCount(c.Count)
		case "units":
			ret.SetUnitCount(c.Count)
		}
	}
	return ret
}

func modelInfoToModel(m *base.ModelInfo) *v1.Model {
	ret := &v1.Model{}
	ret.SetUuid(m.UUID)
	ret.SetName(m.Name)
	return ret
}

func toAddApplicationUnitsMachines(mss []params.MachineStatus) []*v1.AddApplicationUnitsResponse_Machine {
	ret := make([]*v1.AddApplicationUnitsResponse_Machine, len(mss))
	for i := range mss {
		ret[i] = toAddApplicationUnitsMachine(mss[i])
	}
	return ret
}

func toAddApplicationUnitsMachine(ms params.MachineStatus) *v1.AddApplicationUnitsResponse_Machine {
	ret := &v1.AddApplicationUnitsResponse_Machine{}
	ret.SetSystemId(ms.InstanceId.String())
	ret.SetHostname(ms.Hostname)
	return ret
}

func toIntegrations(rss []*params.RelationStatus) []*v1.Integration {
	ret := make([]*v1.Integration, len(rss))
	for i := range rss {
		ret[i] = toIntegration(rss[i])
	}
	return ret
}

func toIntegration(rs *params.RelationStatus) *v1.Integration {
	ret := &v1.Integration{}
	ret.SetId(int64(rs.Id))
	ret.SetInterface(rs.Interface)

	if len(rs.Endpoints) > 0 {
		ret.SetProvider(fmt.Sprintf("%s:%s", rs.Endpoints[0].ApplicationName, rs.Endpoints[0].Name))
		ret.SetRole(rs.Endpoints[0].Role)

		if len(rs.Endpoints) > 1 {
			ret.SetRequirer(fmt.Sprintf("%s:%s", rs.Endpoints[1].ApplicationName, rs.Endpoints[1].Name))
		}
	}
	return ret
}

// TODO: HOW
func resultsToIntegration(arr *params.AddRelationResults) *v1.Integration {
	ret := &v1.Integration{}
	return ret
}

func toActions(m map[string]action.ActionSpec) []*v1.Action {
	ret := []*v1.Action{}
	for k, v := range m {
		ret = append(ret, toAction(k, v))
	}
	return ret
}

func toAction(name string, spec action.ActionSpec) *v1.Action {
	ret := &v1.Action{}
	ret.SetName(name)
	ret.SetDescription(spec.Description)
	return ret
}

func toApplications(as map[string]params.ApplicationStatus, ms map[string]params.MachineStatus, cs map[string]map[string]any) []*v1.Application {
	ret := []*v1.Application{}
	for name := range as {
		app := as[name]
		ret = append(ret, toApplication(name, &app, ms, cs[name]))
	}
	return ret
}

func toApplicationStatus(s *params.DetailedStatus) *v1.Application_Status {
	ret := &v1.Application_Status{}
	ret.SetStatus(s.Status)
	ret.SetInfo(s.Info)
	since := s.Since
	if since != nil {
		ret.SetCreatedAt(timestamppb.New(*since))
	}
	return ret
}

func toApplicationUnit(name string, s *params.UnitStatus, ms map[string]params.MachineStatus) *v1.Application_Unit {
	ret := &v1.Application_Unit{}
	ret.SetName(name)
	ret.SetVersion(s.WorkloadVersion)
	ret.SetLeader(s.Leader)
	ret.SetIpAddress(s.Address + s.PublicAddress)
	ret.SetPorts(s.OpenedPorts)

	if m, ok := ms[s.Machine]; ok {
		ret.SetMachineSystemId(m.InstanceId.String())
	}

	ret.SetAgentStatus(toApplicationStatus(&s.AgentStatus))
	ret.SetWorkloadStatus(toApplicationStatus(&s.WorkloadStatus))

	subordinates := []*v1.Application_Unit{}
	for name := range s.Subordinates {
		unit := s.Subordinates[name]
		subordinates = append(subordinates, toApplicationUnit(name, &unit, ms))
	}

	ret.SetSubordinates(subordinates)

	return ret
}

func toApplication(name string, status *params.ApplicationStatus, ms map[string]params.MachineStatus, c map[string]any) *v1.Application {
	units := []*v1.Application_Unit{}
	for name := range status.Units {
		unit := status.Units[name]
		units = append(units, toApplicationUnit(name, &unit, ms))
	}

	ret := &v1.Application{}
	ret.SetName(name)
	ret.SetVersion(status.WorkloadVersion)
	ret.SetRevision(int64(status.CharmRev))
	ret.SetCharmName(status.Charm)
	ret.SetStatus(status.Status.Status)
	ret.SetInfo(status.Status.Info)

	since := status.Status.Since
	if since != nil {
		ret.SetCreatedAt(timestamppb.New(*since))
	}

	ret.SetUnits(units)

	configYAML, _ := jujuyaml.Marshal(c)
	ret.SetConfigYaml(string(configYAML))

	return ret
}
