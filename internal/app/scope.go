package app

import (
	"context"

	"connectrpc.com/connect"

	pb "github.com/openhdc/openhdc/api/nexus/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *NexusApp) ListScopes(ctx context.Context, req *connect.Request[pb.ListScopesRequest]) (*connect.Response[pb.ListScopesResponse], error) {
	ss, err := a.svc.ListScopes(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.ListScopesResponse{}
	res.SetScopes(toProtoScopes(ss))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) CreateScope(ctx context.Context, req *connect.Request[pb.CreateScopeRequest]) (*connect.Response[pb.Scope], error) {
	s, err := a.svc.CreateScope(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	res := toProtoScope(s)
	return connect.NewResponse(res), nil
}

func toProtoScopes(ss []model.Scope) []*pb.Scope {
	ret := []*pb.Scope{}
	for i := range ss {
		ret = append(ret, toProtoScope(&ss[i]))
	}
	return ret
}

func toProtoScope(s *model.Scope) *pb.Scope {
	ret := &pb.Scope{}
	ret.SetUuid(s.UUID)
	ret.SetName(s.Name)
	ret.SetType(string(s.Type))
	ret.SetProviderType(s.ProviderType)
	ret.SetLife(string(s.Life))
	ret.SetStatus(string(s.Status.Status))
	ret.SetAgentVersion(s.AgentVersion.String())
	ret.SetController(s.IsController)
	for _, c := range s.Counts {
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
