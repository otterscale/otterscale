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

	pb "github.com/otterscale/otterscale/api/fleet/v1"
	"github.com/otterscale/otterscale/api/fleet/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

// FleetService implements the Fleet gRPC service. It handles cluster
// listing and agent registration.
type FleetService struct {
	pbconnect.UnimplementedFleetServiceHandler

	fleet *core.FleetUseCase
}

// NewFleetService returns a FleetService backed by the given use-case.
func NewFleetService(fleet *core.FleetUseCase) *FleetService {
	return &FleetService{
		fleet: fleet,
	}
}

var _ pbconnect.FleetServiceHandler = (*FleetService)(nil)

// ListClusters returns the names of all clusters that have a
// registered agent.
func (s *FleetService) ListClusters(ctx context.Context, req *pb.ListClustersRequest) (*pb.ListClustersResponse, error) {
	clusters := s.fleet.ListClusters(ctx)

	resp := &pb.ListClustersResponse{}
	resp.SetClusters(toProtoClusters(clusters))
	return resp, nil
}

// Register validates and signs the agent's CSR, allocates a tunnel
// endpoint, and returns the signed certificate together with the CA
// certificate for mTLS. The response includes the server version so
// agents can detect mismatches and self-update.
func (s *FleetService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	reg, err := s.fleet.RegisterCluster(ctx, req.GetCluster(), req.GetAgentId(), req.GetAgentVersion(), req.GetCsr())
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
func (s *FleetService) GetAgentManifest(ctx context.Context, req *pb.GetAgentManifestRequest) (*pb.GetAgentManifestResponse, error) {
	userInfo, ok := core.UserInfoFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user info not found in context"))
	}

	cluster := req.GetCluster()

	manifest, err := s.fleet.GenerateAgentManifest(ctx, cluster, userInfo.Subject)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	url, err := s.fleet.IssueManifestURL(ctx, cluster, userInfo.Subject)
	if err != nil {
		return nil, domainErrorToConnectError(err)
	}

	resp := &pb.GetAgentManifestResponse{}
	resp.SetManifest(manifest)
	resp.SetUrl(url)
	return resp, nil
}

// toProtoClusters converts a map of cluster names to Cluster domain
// objects into a sorted slice of protobuf Cluster messages. Results
// are sorted by name to ensure deterministic ordering.
func toProtoClusters(m map[string]core.Cluster) []*pb.Cluster {
	ret := make([]*pb.Cluster, 0, len(m))
	for name, cluster := range m {
		ret = append(ret, toProtoCluster(name, cluster))
	}
	slices.SortFunc(ret, func(a, b *pb.Cluster) int {
		return cmp.Compare(a.GetName(), b.GetName())
	})
	return ret
}

// toProtoCluster converts a cluster name and its domain object into a
// protobuf Cluster message.
func toProtoCluster(name string, cluster core.Cluster) *pb.Cluster {
	ret := &pb.Cluster{}
	ret.SetName(name)
	ret.SetAgentVersion(cluster.AgentVersion)
	return ret
}
