// Package handler implements the ConnectRPC service handlers that form
// the server's public API. Each handler translates between protobuf
// messages and the domain use-cases defined in package core.
package handler

import (
	"cmp"
	"context"
	"slices"

	pb "github.com/otterscale/otterscale/api/link/v1"

	"github.com/otterscale/otterscale/internal/core"
)

// LinkService implements the Link gRPC service. It handles cluster
// listing and agent registration.
type LinkService struct {
	pb.UnimplementedLinkServiceHandler

	link *core.LinkUseCase
}

// NewLinkService returns a LinkService backed by the given use-case.
func NewLinkService(link *core.LinkUseCase) *LinkService {
	return &LinkService{
		link: link,
	}
}

var _ pb.LinkServiceHandler = (*LinkService)(nil)

// ListLinks returns the names of all clusters that have a
// registered agent.
func (s *LinkService) ListLinks(ctx context.Context, _ *pb.ListLinksRequest) (*pb.ListLinksResponse, error) {
	links := s.link.ListLinks(ctx)

	resp := &pb.ListLinksResponse{}
	resp.SetLinks(toProtoLinks(links))
	return resp, nil
}

// Register validates and signs the agent's CSR, allocates a tunnel
// endpoint, and returns the signed certificate together with the CA
// certificate for mTLS. The response includes the server version for
// diagnostics.
func (s *LinkService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	reg, err := s.link.RegisterCluster(ctx, req.GetCluster(), req.GetAgentId(), req.GetAgentVersion(), req.GetCsr())
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	resp := &pb.RegisterResponse{}
	resp.SetEndpoint(reg.Endpoint)
	resp.SetCertificate(reg.Certificate)
	resp.SetCaCertificate(reg.CACertificate)
	resp.SetServerVersion(reg.ServerVersion)
	return resp, nil
}

// toProtoLinks converts a map of cluster names to Link domain
// objects into a sorted slice of protobuf Link messages. Results
// are sorted by name to ensure deterministic ordering.
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

// toProtoLink converts a cluster name and its domain object into a
// protobuf Link message.
func toProtoLink(cluster string, link core.Link) *pb.Link {
	ret := &pb.Link{}
	ret.SetCluster(cluster)
	ret.SetAgentVersion(link.AgentVersion)
	return ret
}
