package app

import (
	"context"

	"connectrpc.com/connect"

	pb "github.com/otterscale/otterscale/api/scope/v1"
	"github.com/otterscale/otterscale/api/scope/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type ScopeService struct {
	pbconnect.UnimplementedScopeServiceHandler

	uc *core.ScopeUseCase
}

func NewScopeService(uc *core.ScopeUseCase) *ScopeService {
	return &ScopeService{uc: uc}
}

var _ pbconnect.ScopeServiceHandler = (*ScopeService)(nil)

func (s *ScopeService) ListScopes(ctx context.Context, req *connect.Request[pb.ListScopesRequest]) (*connect.Response[pb.ListScopesResponse], error) {
	scopes, err := s.uc.ListScopes(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListScopesResponse{}
	resp.SetScopes(toProtoScopes(scopes))
	return connect.NewResponse(resp), nil
}

func (s *ScopeService) CreateScope(ctx context.Context, req *connect.Request[pb.CreateScopeRequest]) (*connect.Response[pb.Scope], error) {
	scope, err := s.uc.CreateScope(ctx, req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoScope(scope)
	return connect.NewResponse(resp), nil
}

func toProtoScopes(ss []core.Scope) []*pb.Scope {
	ret := []*pb.Scope{}
	for i := range ss {
		ret = append(ret, toProtoScope(&ss[i]))
	}
	return ret
}

func toProtoScope(s *core.Scope) *pb.Scope {
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
