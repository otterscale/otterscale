// Package handler implements the ConnectRPC service handlers that form
// the server's public API. Each handler translates between protobuf
// messages and the domain use-cases defined in package core.
package handler

import (
	"cmp"
	"context"
	"errors"
	"slices"

	"connectrpc.com/connect"

	pb "github.com/otterscale/otterscale/api/link/v1"

	"github.com/otterscale/otterscale/internal/core"
)

// LinkService handles cluster listing and agent registration.
type LinkService struct {
	pb.UnimplementedLinkServiceHandler

	link *core.LinkUseCase
}

func NewLinkService(link *core.LinkUseCase) *LinkService {
	return &LinkService{
		link: link,
	}
}

var _ pb.LinkServiceHandler = (*LinkService)(nil)

// ListLinks returns every cluster with a registered agent.
func (s *LinkService) ListLinks(ctx context.Context, _ *pb.ListLinksRequest) (*pb.ListLinksResponse, error) {
	links := s.link.ListLinks(ctx)

	resp := &pb.ListLinksResponse{}
	resp.SetLinks(toProtoLinks(links))
	return resp, nil
}

// Register signs the agent's CSR, allocates a tunnel endpoint, and returns the
// certificate with the CA certificate for mTLS, plus the server version for
// diagnostics.
func (s *LinkService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	reg, err := s.link.RegisterCluster(ctx, &core.RegistrationRequest{
		Cluster:      req.GetCluster(),
		AgentID:      req.GetAgentId(),
		AgentVersion: req.GetAgentVersion(),
		JoinToken:    req.GetJoinToken(),
		CSRPEM:       req.GetCsr(),
	})
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	resp := &pb.RegisterResponse{}
	resp.SetEndpoint(reg.Endpoint)
	resp.SetCertificate(reg.Certificate)
	resp.SetCaCertificate(reg.CACertificate)
	resp.SetTunnelUser(reg.TunnelUser)
	resp.SetTunnelPassword(reg.TunnelPassword)
	resp.SetServerVersion(reg.ServerVersion)
	return resp, nil
}

// IssueJoinToken derives the join token for a cluster, so an import flow can
// obtain one over the API instead of exec-ing `otterscale join token` inside
// the server pod.
//
// Restricted to the admin group: what this returns authorizes claiming a
// cluster, and thereby cluster-admin on it.
func (s *LinkService) IssueJoinToken(ctx context.Context, req *pb.IssueJoinTokenRequest) (*pb.IssueJoinTokenResponse, error) {
	userInfo, ok := core.UserInfoFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user info not found in context"))
	}
	if !core.IsAdmin(userInfo.Groups) {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("caller is not a member of the admin group"))
	}

	token, err := s.link.IssueJoinToken(ctx, req.GetCluster())
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	resp := &pb.IssueJoinTokenResponse{}
	resp.SetJoinToken(token)
	return resp, nil
}

// toProtoLinks sorts by cluster name, for deterministic ordering.
func toProtoLinks(m map[string]core.Link) []*pb.Link {
	ret := make([]*pb.Link, 0, len(m))
	for cluster, link := range m {
		ret = append(ret, toProtoLink(cluster, link))
	}
	slices.SortFunc(ret, func(a, b *pb.Link) int {
		return cmp.Compare(a.GetCluster(), b.GetCluster())
	})
	return ret
}

func toProtoLink(cluster string, link core.Link) *pb.Link {
	ret := &pb.Link{}
	ret.SetCluster(cluster)
	ret.SetAgentVersion(link.AgentVersion)
	return ret
}
