package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/otterscale/otterscale/api/network/v1"
	"github.com/otterscale/otterscale/api/network/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/network"
)

type NetworkService struct {
	pbconnect.UnimplementedNetworkServiceHandler

	network *network.UseCase
}

func NewNetworkService(network *network.UseCase) *NetworkService {
	return &NetworkService{
		network: network,
	}
}

var _ pbconnect.NetworkServiceHandler = (*NetworkService)(nil)

func (s *NetworkService) ListNetworks(ctx context.Context, _ *pb.ListNetworksRequest) (*pb.ListNetworksResponse, error) {
	networks, err := s.network.ListNetworks(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListNetworksResponse{}
	resp.SetNetworks(toProtoNetworks(networks))
	return resp, nil
}

func (s *NetworkService) CreateNetwork(ctx context.Context, req *pb.CreateNetworkRequest) (*pb.Network, error) {
	network, err := s.network.CreateNetwork(ctx, req.GetCidr(), req.GetGatewayIp(), req.GetDnsServers(), req.GetDhcpOn())
	if err != nil {
		return nil, err
	}

	resp := toProtoNetwork(network)
	return resp, nil
}

func (s *NetworkService) CreateIPRange(ctx context.Context, req *pb.CreateIPRangeRequest) (*pb.Network_IPRange, error) {
	ipRange, err := s.network.CreateIPRange(ctx, int(req.GetSubnetId()), req.GetStartIp(), req.GetEndIp(), req.GetComment())
	if err != nil {
		return nil, err
	}

	resp := toProtoIPRange(ipRange)
	return resp, nil
}

func (s *NetworkService) DeleteNetwork(ctx context.Context, req *pb.DeleteNetworkRequest) (*emptypb.Empty, error) {
	if err := s.network.DeleteNetwork(ctx, int(req.GetId())); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *NetworkService) DeleteIPRange(ctx context.Context, req *pb.DeleteIPRangeRequest) (*emptypb.Empty, error) {
	if err := s.network.DeleteIPRange(ctx, int(req.GetId())); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *NetworkService) UpdateFabric(ctx context.Context, req *pb.UpdateFabricRequest) (*pb.Network_Fabric, error) {
	fabric, err := s.network.UpdateFabric(ctx, int(req.GetId()), req.GetName())
	if err != nil {
		return nil, err
	}

	resp := toProtoFabric(fabric)
	return resp, nil
}

func (s *NetworkService) UpdateVLAN(ctx context.Context, req *pb.UpdateVLANRequest) (*pb.Network_VLAN, error) {
	vlan, err := s.network.UpdateVLAN(ctx, int(req.GetFabricId()), int(req.GetVid()), req.GetName(), int(req.GetMtu()), req.GetDescription(), req.GetDhcpOn())
	if err != nil {
		return nil, err
	}

	resp := toProtoVLAN(vlan)
	return resp, nil
}

func (s *NetworkService) UpdateSubnet(ctx context.Context, req *pb.UpdateSubnetRequest) (*pb.Network_Subnet, error) {
	subnet, err := s.network.UpdateSubnet(ctx, int(req.GetId()), req.GetName(), req.GetCidr(), req.GetGatewayIp(), req.GetDnsServers(), req.GetDescription(), req.GetAllowDnsResolution())
	if err != nil {
		return nil, err
	}

	resp := toProtoSubnet(subnet)
	return resp, nil
}

func (s *NetworkService) UpdateIPRange(ctx context.Context, req *pb.UpdateIPRangeRequest) (*pb.Network_IPRange, error) {
	ipRange, err := s.network.UpdateIPRange(ctx, int(req.GetId()), req.GetStartIp(), req.GetEndIp(), req.GetComment())
	if err != nil {
		return nil, err
	}

	resp := toProtoIPRange(ipRange)
	return resp, nil
}

func toProtoNetworks(ns []network.Network) []*pb.Network {
	ret := []*pb.Network{}

	for i := range ns {
		ret = append(ret, toProtoNetwork(&ns[i]))
	}

	return ret
}

func toProtoNetwork(n *network.Network) *pb.Network {
	ret := &pb.Network{}
	ret.SetFabric(toProtoFabric(n.Fabric))
	ret.SetVlan(toProtoVLAN(n.VLAN))

	if n.Subnet != nil {
		ret.SetSubnet(toProtoSubnet(n.Subnet))
	}

	return ret
}

func toProtoIPAddresses(ipas []network.IPAddress) []*pb.Network_IPAddress {
	ret := []*pb.Network_IPAddress{}

	for i := range ipas {
		ret = append(ret, toProtoIPAddress(&ipas[i]))
	}

	return ret
}

func toProtoIPAddress(ipa *network.IPAddress) *pb.Network_IPAddress {
	ret := &pb.Network_IPAddress{}
	ret.SetType(network.AllocType(ipa.AllocType).String())
	ret.SetIp(ipa.IP.String())
	ret.SetUser(ipa.User)
	ret.SetMachineId(ipa.NodeSummary.SystemID)
	ret.SetNodeType(network.NodeType(ipa.NodeSummary.NodeType).String())
	ret.SetHostname(ipa.NodeSummary.Hostname)
	return ret
}

func toProtoIPRanges(iprs []network.IPRange) []*pb.Network_IPRange {
	ret := []*pb.Network_IPRange{}

	for i := range iprs {
		ret = append(ret, toProtoIPRange(&iprs[i]))
	}

	return ret
}

func toProtoIPRange(ipr *network.IPRange) *pb.Network_IPRange {
	ret := &pb.Network_IPRange{}
	ret.SetId(int64(ipr.ID))
	ret.SetType(ipr.Type)
	ret.SetStartIp(ipr.StartIP.String())
	ret.SetEndIp(ipr.EndIP.String())
	ret.SetComment(ipr.Comment)
	return ret
}

func toProtoStatistics(ns *network.Statistics) *pb.Network_Statistics {
	ret := &pb.Network_Statistics{}
	ret.SetAvailable(int64(ns.NumAvailable))
	ret.SetTotal(int64(ns.TotalAddresses))
	ret.SetUsagePercent(ns.UsageString)
	ret.SetAvailablePercent(ns.AvailableString)
	return ret
}

func toProtoFabric(f *network.Fabric) *pb.Network_Fabric {
	ret := &pb.Network_Fabric{}
	ret.SetId(int64(f.ID))
	ret.SetName(f.Name)
	return ret
}

func toProtoVLAN(v *network.VLAN) *pb.Network_VLAN {
	ret := &pb.Network_VLAN{}
	ret.SetId(int64(v.ID))
	ret.SetVid(int64(v.VID))
	ret.SetName(v.Name)
	ret.SetMtu(int64(v.MTU))
	ret.SetDescription(v.Description)
	ret.SetDhcpOn(v.DHCPOn)
	return ret
}

func toProtoSubnet(ns *network.NetworkSubnet) *pb.Network_Subnet {
	dnsServers := make([]string, len(ns.DNSServers))
	for i, dns := range ns.DNSServers {
		dnsServers[i] = dns.String()
	}

	ret := &pb.Network_Subnet{}
	ret.SetId(int64(ns.ID))
	ret.SetName(ns.Name)
	ret.SetCidr(ns.CIDR)
	ret.SetGatewayIp(ns.GatewayIP.String())
	ret.SetDnsServers(dnsServers)
	ret.SetDescription(ns.Description)
	ret.SetManagedAllocation(ns.Managed)
	ret.SetActiveDiscovery(ns.ActiveDiscovery)
	ret.SetAllowProxyAccess(ns.AllowProxy)
	ret.SetAllowDnsResolution(ns.AllowDNS)
	ret.SetIpAddresses(toProtoIPAddresses(ns.IPAddresses))
	ret.SetIpRanges(toProtoIPRanges(ns.IPRanges))
	ret.SetStatistics(toProtoStatistics(ns.Statistics))
	return ret
}
