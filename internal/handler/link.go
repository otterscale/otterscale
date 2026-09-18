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
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/link/v1"

	"github.com/otterscale/otterscale/internal/core"
)

// LinkService handles cluster listing, agent registration, and what an
// operator needs to install an agent on a joining cluster.
type LinkService struct {
	pb.UnimplementedLinkServiceHandler

	link        *core.LinkUseCase
	agentValues *core.AgentValuesUseCase
}

func NewLinkService(link *core.LinkUseCase, agentValues *core.AgentValuesUseCase) *LinkService {
	return &LinkService{
		link:        link,
		agentValues: agentValues,
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

// IssueAgentValues renders the Helm override values that install the agent on
// a joining cluster, plus a URL serving the same bytes.
//
// Restricted to the admin group: the result embeds a join token, which claims
// the cluster it names, and it binds the caller to cluster-admin on it.
func (s *LinkService) IssueAgentValues(ctx context.Context, req *pb.IssueAgentValuesRequest) (*pb.IssueAgentValuesResponse, error) {
	userInfo, err := requireAdmin(ctx)
	if err != nil {
		return nil, err
	}

	result, err := s.agentValues.Issue(ctx, &core.AgentValuesRequest{
		Cluster:     req.GetCluster(),
		Subject:     userInfo.Subject,
		ExtraUsers:  req.GetExtraUsers(),
		ClusterInfo: toCoreClusterInfo(req),
	})
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	resp := &pb.IssueAgentValuesResponse{}
	resp.SetValues(result.YAML)
	resp.SetUrl(result.URL)
	resp.SetUrlExpiresAt(timestamppb.New(result.ExpiresAt))
	return resp, nil
}

// requireAdmin is the gate on a procedure that hands out cluster-admin.
func requireAdmin(ctx context.Context) (core.UserInfo, error) {
	userInfo, ok := core.UserInfoFromContext(ctx)
	if !ok {
		return core.UserInfo{}, connect.NewError(connect.CodeUnauthenticated, errors.New("user info not found in context"))
	}
	if !core.IsAdmin(userInfo.Groups) {
		return core.UserInfo{}, connect.NewError(connect.CodePermissionDenied, errors.New("caller is not a member of the admin group"))
	}
	return userInfo, nil
}

// toCoreClusterInfo returns nil when the caller sent no cluster_info, which
// the use case rejects. Presence is why that field is a message: a flat
// string could not tell "not supplied" from "supplied empty".
func toCoreClusterInfo(req *pb.IssueAgentValuesRequest) *core.AgentClusterInfo {
	if !req.HasClusterInfo() {
		return nil
	}
	info := req.GetClusterInfo()
	return &core.AgentClusterInfo{
		ExternalAddress: info.GetExternalAddress(),
		NodePortRange:   info.GetNodePortRange(),
		InferenceURL:    info.GetInferenceUrl(),
	}
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
