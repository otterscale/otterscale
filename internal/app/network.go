package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/otterscale/api/network/v1"
	"github.com/openhdc/otterscale/api/network/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/core"
	"github.com/openhdc/otterscale/internal/enum"
)

type NetworkService struct {
	pbconnect.UnimplementedNetworkServiceHandler

	uc *core.NetworkUseCase
}

func NewNetworkService(uc *core.NetworkUseCase) *NetworkService {
	return &NetworkService{uc: uc}
}

var _ pbconnect.NetworkServiceHandler = (*NetworkService)(nil)

func (s *NetworkService) ListNetworks(ctx context.Context, req *connect.Request[pb.ListNetworksRequest]) (*connect.Response[pb.ListNetworksResponse], error) {
	networks, err := s.uc.ListNetworks(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListNetworksResponse{}
	resp.SetNetworks(toProtoNetworks(networks))
	return connect.NewResponse(resp), nil
}

func (s *NetworkService) CreateNetwork(ctx context.Context, req *connect.Request[pb.CreateNetworkRequest]) (*connect.Response[pb.Network], error) {
	network, err := s.uc.CreateNetwork(ctx, req.Msg.GetCidr(), req.Msg.GetGatewayIp(), req.Msg.GetDnsServers(), req.Msg.GetDhcpOn())
	if err != nil {
		return nil, err
	}
	resp := toProtoNetwork(network)
	return connect.NewResponse(resp), nil
}

func (s *NetworkService) CreateIPRange(ctx context.Context, req *connect.Request[pb.CreateIPRangeRequest]) (*connect.Response[pb.Network_IPRange], error) {
	ipRange, err := s.uc.CreateIPRange(ctx, int(req.Msg.GetSubnetId()), req.Msg.GetStartIp(), req.Msg.GetEndIp(), req.Msg.GetComment())
	if err != nil {
		return nil, err
	}
	resp := toProtoIPRange(ipRange)
	return connect.NewResponse(resp), nil
}

func (s *NetworkService) DeleteNetwork(ctx context.Context, req *connect.Request[pb.DeleteNetworkRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteNetwork(ctx, int(req.Msg.GetId())); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *NetworkService) DeleteIPRange(ctx context.Context, req *connect.Request[pb.DeleteIPRangeRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteIPRange(ctx, int(req.Msg.GetId())); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *NetworkService) UpdateFabric(ctx context.Context, req *connect.Request[pb.UpdateFabricRequest]) (*connect.Response[pb.Network_Fabric], error) {
	fabric, err := s.uc.UpdateFabric(ctx, int(req.Msg.GetId()), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoFabric(fabric)
	return connect.NewResponse(resp), nil
}

func (s *NetworkService) UpdateVLAN(ctx context.Context, req *connect.Request[pb.UpdateVLANRequest]) (*connect.Response[pb.Network_VLAN], error) {
	vlan, err := s.uc.UpdateVLAN(ctx, int(req.Msg.GetFabricId()), int(req.Msg.GetVid()), req.Msg.GetName(), int(req.Msg.GetMtu()), req.Msg.GetDescription(), req.Msg.GetDhcpOn())
	if err != nil {
		return nil, err
	}
	resp := toProtoVLAN(vlan)
	return connect.NewResponse(resp), nil
}

func (s *NetworkService) UpdateSubnet(ctx context.Context, req *connect.Request[pb.UpdateSubnetRequest]) (*connect.Response[pb.Network_Subnet], error) {
	subnet, err := s.uc.UpdateSubnet(ctx, int(req.Msg.GetId()), req.Msg.GetName(), req.Msg.GetCidr(), req.Msg.GetGatewayIp(), req.Msg.GetDnsServers(), req.Msg.GetDescription(), req.Msg.GetAllowDnsResolution())
	if err != nil {
		return nil, err
	}
	resp := toProtoSubnet(subnet)
	return connect.NewResponse(resp), nil
}

func (s *NetworkService) UpdateIPRange(ctx context.Context, req *connect.Request[pb.UpdateIPRangeRequest]) (*connect.Response[pb.Network_IPRange], error) {
	ipRange, err := s.uc.UpdateIPRange(ctx, int(req.Msg.GetId()), req.Msg.GetStartIp(), req.Msg.GetEndIp(), req.Msg.GetComment())
	if err != nil {
		return nil, err
	}
	resp := toProtoIPRange(ipRange)
	return connect.NewResponse(resp), nil
}

func toProtoNetworks(ns []core.Network) []*pb.Network {
	ret := []*pb.Network{}
	for i := range ns {
		ret = append(ret, toProtoNetwork(&ns[i]))
	}
	return ret
}

func toProtoNetwork(n *core.Network) *pb.Network {
	ret := &pb.Network{}
	ret.SetFabric(toProtoFabric(n.Fabric))
	ret.SetVlan(toProtoVLAN(n.VLAN))
	if n.Subnet != nil {
		ret.SetSubnet(toProtoSubnet(n.Subnet))
	}
	return ret
}

func toProtoIPAddresses(ipas []core.IPAddress) []*pb.Network_IPAddress {
	ret := []*pb.Network_IPAddress{}
	for i := range ipas {
		ret = append(ret, toProtoIPAddress(&ipas[i]))
	}
	return ret
}

func toProtoIPAddress(ipa *core.IPAddress) *pb.Network_IPAddress {
	ret := &pb.Network_IPAddress{}
	ret.SetType(enum.AllocType(ipa.AllocType).String())
	ret.SetIp(ipa.IP.String())
	ret.SetUser(ipa.User)
	ret.SetMachineId(ipa.NodeSummary.SystemID)
	ret.SetNodeType(enum.NodeType(ipa.NodeSummary.NodeType).String())
	ret.SetHostname(ipa.NodeSummary.Hostname)
	return ret
}

func toProtoIPRanges(iprs []core.IPRange) []*pb.Network_IPRange {
	ret := []*pb.Network_IPRange{}
	for i := range iprs {
		ret = append(ret, toProtoIPRange(&iprs[i]))
	}
	return ret
}

func toProtoIPRange(ipr *core.IPRange) *pb.Network_IPRange {
	ret := &pb.Network_IPRange{}
	ret.SetId(int64(ipr.ID))
	ret.SetType(ipr.Type)
	ret.SetStartIp(ipr.StartIP.String())
	ret.SetEndIp(ipr.EndIP.String())
	ret.SetComment(ipr.Comment)
	return ret
}

func toProtoStatistics(ns *core.NetworkStatistics) *pb.Network_Statistics {
	ret := &pb.Network_Statistics{}
	ret.SetAvailable(int64(ns.NumAvailable))
	ret.SetTotal(int64(ns.TotalAddresses))
	ret.SetUsagePercent(ns.UsageString)
	ret.SetAvailablePercent(ns.AvailableString)
	return ret
}

func toProtoFabric(f *core.Fabric) *pb.Network_Fabric {
	ret := &pb.Network_Fabric{}
	ret.SetId(int64(f.ID))
	ret.SetName(f.Name)
	return ret
}

func toProtoVLAN(v *core.VLAN) *pb.Network_VLAN {
	ret := &pb.Network_VLAN{}
	ret.SetId(int64(v.ID))
	ret.SetVid(int64(v.VID))
	ret.SetName(v.Name)
	ret.SetMtu(int64(v.MTU))
	ret.SetDescription(v.Description)
	ret.SetDhcpOn(v.DHCPOn)
	return ret
}

func toProtoSubnet(ns *core.NetworkSubnet) *pb.Network_Subnet {
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
