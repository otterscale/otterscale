package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"connectrpc.com/connect"
	pb "github.com/openhdc/openhdc/api/nexus/v1"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *NexusApp) ListNetworks(ctx context.Context, req *connect.Request[pb.ListNetworksRequest]) (*connect.Response[pb.ListNetworksResponse], error) {
	ss, err := a.svc.ListNetworks(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.ListNetworksResponse{}
	res.SetNetworks(toProtoNetworks(ss))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) CreateNetwork(ctx context.Context, req *connect.Request[pb.CreateNetworkRequest]) (*connect.Response[pb.Network], error) {
	n, err := a.svc.CreateNetwork(ctx, req.Msg.GetCidr(), req.Msg.GetGatewayIp(), req.Msg.GetDnsServers(), req.Msg.GetDhcpOn())
	if err != nil {
		return nil, err
	}
	res := toProtoNetwork(n)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) CreateIPRange(ctx context.Context, req *connect.Request[pb.CreateIPRangeRequest]) (*connect.Response[pb.Network_IPRange], error) {
	ipr, err := a.svc.CreateIPRange(ctx, int(req.Msg.GetSubnetId()), req.Msg.GetStartIp(), req.Msg.GetEndIp(), req.Msg.GetComment())
	if err != nil {
		return nil, err
	}
	res := toProtoIPRange(ipr)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) DeleteNetwork(ctx context.Context, req *connect.Request[pb.DeleteNetworkRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.DeleteNetwork(ctx, int(req.Msg.GetId())); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) DeleteIPRange(ctx context.Context, req *connect.Request[pb.DeleteIPRangeRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.DeleteIPRange(ctx, int(req.Msg.GetId())); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) UpdateFabric(ctx context.Context, req *connect.Request[pb.UpdateFabricRequest]) (*connect.Response[pb.Network_Fabric], error) {
	fab, err := a.svc.UpdateFabric(ctx, int(req.Msg.GetId()), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	res := toProtoFabric(fab)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) UpdateVLAN(ctx context.Context, req *connect.Request[pb.UpdateVLANRequest]) (*connect.Response[pb.Network_VLAN], error) {
	vlan, err := a.svc.UpdateVLAN(ctx, int(req.Msg.GetFabricId()), int(req.Msg.GetVid()), req.Msg.GetName(), int(req.Msg.GetMtu()), req.Msg.GetDescription(), req.Msg.GetDhcpOn())
	if err != nil {
		return nil, err
	}
	res := toProtoVLAN(vlan)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) UpdateSubnet(ctx context.Context, req *connect.Request[pb.UpdateSubnetRequest]) (*connect.Response[pb.Network_Subnet], error) {
	sub, err := a.svc.UpdateSubnet(ctx, int(req.Msg.GetId()), req.Msg.GetName(), req.Msg.GetCidr(), req.Msg.GetGatewayIp(), req.Msg.GetDnsServers(), req.Msg.GetDescription(), req.Msg.GetAllowDnsResolution())
	if err != nil {
		return nil, err
	}
	res := toProtoSubnet(sub)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) UpdateIPRange(ctx context.Context, req *connect.Request[pb.UpdateIPRangeRequest]) (*connect.Response[pb.Network_IPRange], error) {
	ipr, err := a.svc.UpdateIPRange(ctx, int(req.Msg.GetId()), req.Msg.GetStartIp(), req.Msg.GetEndIp(), req.Msg.GetComment())
	if err != nil {
		return nil, err
	}
	res := toProtoIPRange(ipr)
	return connect.NewResponse(res), nil
}

func toProtoNetworks(ns []model.Network) []*pb.Network {
	ret := []*pb.Network{}
	for i := range ns {
		ret = append(ret, toProtoNetwork(&ns[i]))
	}
	return ret
}

func toProtoNetwork(n *model.Network) *pb.Network {
	ret := &pb.Network{}
	ret.SetFabric(toProtoFabric(n.Fabric))
	ret.SetVlan(toProtoVLAN(n.VLAN))
	if n.Subnet != nil {
		ret.SetSubnet(toProtoSubnet(n.Subnet))
	}
	return ret
}

func toProtoIPAddresses(ipas []model.IPAddress) []*pb.Network_IPAddress {
	ret := []*pb.Network_IPAddress{}
	for i := range ipas {
		ret = append(ret, toProtoIPAddress(&ipas[i]))
	}
	return ret
}

func toProtoIPAddress(ipa *model.IPAddress) *pb.Network_IPAddress {
	ret := &pb.Network_IPAddress{}
	ret.SetType(model.AllocType(ipa.AllocType).String())
	ret.SetIp(ipa.IP.String())
	ret.SetUser(ipa.User)
	ret.SetMachineId(ipa.NodeSummary.SystemID)
	ret.SetNodeType(model.NodeType(ipa.NodeSummary.NodeType).String())
	ret.SetHostname(ipa.NodeSummary.Hostname)
	return ret
}

func toProtoIPRanges(iprs []model.IPRange) []*pb.Network_IPRange {
	ret := []*pb.Network_IPRange{}
	for i := range iprs {
		ret = append(ret, toProtoIPRange(&iprs[i]))
	}
	return ret
}

func toProtoIPRange(ipr *model.IPRange) *pb.Network_IPRange {
	ret := &pb.Network_IPRange{}
	ret.SetId(int64(ipr.ID))
	ret.SetType(ipr.Type)
	ret.SetStartIp(ipr.StartIP.String())
	ret.SetEndIp(ipr.EndIP.String())
	ret.SetComment(ipr.Comment)
	return ret
}

func toProtoStatistics(ns *model.NetworkStatistics) *pb.Network_Statistics {
	ret := &pb.Network_Statistics{}
	ret.SetAvailable(int64(ns.NumAvailable))
	ret.SetTotal(int64(ns.TotalAddresses))
	ret.SetUsagePercent(ns.UsageString)
	ret.SetAvailablePercent(ns.AvailableString)
	return ret
}

func toProtoFabric(f *model.Fabric) *pb.Network_Fabric {
	ret := &pb.Network_Fabric{}
	ret.SetId(int64(f.ID))
	ret.SetName(f.Name)
	return ret
}

func toProtoVLAN(v *model.VLAN) *pb.Network_VLAN {
	ret := &pb.Network_VLAN{}
	ret.SetId(int64(v.ID))
	ret.SetVid(int64(v.VID))
	ret.SetName(v.Name)
	ret.SetMtu(int64(v.MTU))
	ret.SetDescription(v.Description)
	ret.SetDhcpOn(v.DHCPOn)
	return ret
}

func toProtoSubnet(ns *model.NetworkSubnet) *pb.Network_Subnet {
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
