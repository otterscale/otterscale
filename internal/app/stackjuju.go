package app

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"

	v1 "github.com/openhdc/openhdc/api/stack/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

// AddMachines adds machines to the model based on the provided parameters
func (a *StackApp) AddMachines(ctx context.Context, req *connect.Request[v1.AddMachinesRequest]) (*connect.Response[v1.AddMachinesResponse], error) {
	machineParams := buildMachineParams(req.Msg.GetParameters())

	machines, err := a.svc.AddMachine(ctx, req.Msg.GetModelUuid(), machineParams)
	if err != nil {
		return nil, err
	}

	res := &v1.AddMachinesResponse{}
	res.SetMachines(machines)
	return connect.NewResponse(res), nil
}

// buildMachineParams converts API machine parameters to Juju machine parameters
func buildMachineParams(parameters []*v1.AddMachinesRequest_Parameter) []params.AddMachineParams {
	machineParams := make([]params.AddMachineParams, 0, len(parameters))

	for _, param := range parameters {
		machineParam := params.AddMachineParams{
			Placement: &instance.Placement{
				Scope: param.GetPlacement(),
			},
		}

		if constraint := param.GetConstraint(); constraint != nil {
			machineParam.Constraints = buildConstraints(constraint)
		}

		machineParams = append(machineParams, machineParam)
	}

	return machineParams
}

// buildConstraints converts API constraints to Juju constraints
func buildConstraints(c *v1.AddMachinesRequest_Constraint) constraints.Value {
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

func (a *StackApp) GetModelConfigs(ctx context.Context, req *connect.Request[v1.GetModelConfigsRequest]) (*connect.Response[v1.GetModelConfigsResponse], error) {
	modelConfigs, err := a.svc.GetModelConfigs(ctx, req.Msg.GetUuid())
	if err != nil {
		return nil, err
	}
	mc := map[string]string{}
	for k, v := range modelConfigs {
		mc[k] = fmt.Sprintf("%v", v)
	}
	res := &v1.GetModelConfigsResponse{}
	res.SetModelConfigs(mc)
	return connect.NewResponse(res), nil
}

func (a *StackApp) ListApplications(ctx context.Context, req *connect.Request[v1.ListApplicationsRequest]) (*connect.Response[v1.ListApplicationsResponse], error) {
	return nil, nil
}

func (a *StackApp) CreateApplication(ctx context.Context, req *connect.Request[v1.CreateApplicationRequest]) (*connect.Response[v1.Application], error) {
	return nil, nil
}

func (a *StackApp) DeleteApplication(ctx context.Context, req *connect.Request[v1.DeleteApplicationRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, nil
}

func (a *StackApp) UpdateApplicationUnits(ctx context.Context, req *connect.Request[v1.UpdateApplicationUnitsRequest]) (*connect.Response[v1.Application], error) {
	return nil, nil
}

func (a *StackApp) UpdateApplicationRelation(ctx context.Context, req *connect.Request[v1.UpdateApplicationRelationRequest]) (*connect.Response[v1.Application], error) {
	return nil, nil
}

func (a *StackApp) UpdateApplicationConfig(ctx context.Context, req *connect.Request[v1.UpdateApplicationConfigRequest]) (*connect.Response[v1.Application], error) {
	return nil, nil
}

func (a *StackApp) ListIntegrations(ctx context.Context, req *connect.Request[v1.ListIntegrationsRequest]) (*connect.Response[v1.ListIntegrationsResponse], error) {
	return nil, nil
}

func (a *StackApp) ListActions(ctx context.Context, req *connect.Request[v1.ListActionsRequest]) (*connect.Response[v1.ListActionsResponse], error) {
	return nil, nil
}

func (a *StackApp) RunAction(ctx context.Context, req *connect.Request[v1.RunActionRequest]) (*connect.Response[v1.Action], error) {
	return nil, nil
}

func toModels(es []*model.Environment) []*v1.Model {
	ret := make([]*v1.Model, len(es))
	for i := range es {
		ret[i] = toModel(es[i])
	}
	return ret
}

func toModel(m *model.Environment) *v1.Model {
	return &v1.Model{}
}

func modelInfoToModel(m *base.ModelInfo) *v1.Model {
	return &v1.Model{}
}
