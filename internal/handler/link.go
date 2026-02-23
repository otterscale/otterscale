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

	pb "github.com/otterscale/api/link/v1"
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
// certificate for mTLS. The response includes the server version so
// agents can detect mismatches and self-update.
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

// GetAgentManifest returns a multi-document YAML manifest for
// installing the otterscale agent on the caller's target cluster.
// The manifest includes a ClusterRoleBinding that grants the
// authenticated user cluster-admin access.
func (s *LinkService) GetAgentManifest(ctx context.Context, req *pb.GetAgentManifestRequest) (*pb.GetAgentManifestResponse, error) {
	userInfo, ok := core.UserInfoFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user info not found in context"))
	}

	cluster := req.GetCluster()

	manifest, err := s.link.GenerateAgentManifest(ctx, cluster, userInfo.Subject)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	url, err := s.link.IssueManifestURL(ctx, cluster, userInfo.Subject)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	resp := &pb.GetAgentManifestResponse{}
	resp.SetManifest(manifest)
	resp.SetUrl(url)
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
