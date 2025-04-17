package service

import (
	"context"
	"net"
	"strconv"

	"github.com/openhdc/openhdc/internal/domain/model"
)

func (s *NexusService) ListNetworks(ctx context.Context) ([]model.Network, error) {
	subnets, err := s.subnet.List(ctx)
	if err != nil {
		return nil, err
	}

	networkSubnets := make([]model.NetworkSubnet, len(subnets))
	for i := range subnets {
		ns, err := s.getNetworkSubnet(ctx, &subnets[i])
		if err != nil {
			return nil, err
		}
		networkSubnets[i] = *ns
	}

	fabrics, err := s.fabric.List(ctx)
	if err != nil {
		return nil, err
	}

	networks := []model.Network{}
	for i := range fabrics {
		for j := range fabrics[i].VLANs {
			exist := false
			for k := range networkSubnets {
				if networkSubnets[k].VLAN.ID != fabrics[i].VLANs[j].ID { // not only one
					continue
				}
				networks = append(networks, model.Network{
					Fabric: &fabrics[i],
					VLAN:   &fabrics[i].VLANs[j],
					Subnet: &networkSubnets[k],
				})
				exist = true
			}
			if !exist {
				networks = append(networks, model.Network{
					Fabric: &fabrics[i],
					VLAN:   &fabrics[i].VLANs[j],
				})
			}
		}
	}
	return networks, nil
}

func (s *NexusService) CreateNetwork(ctx context.Context, cidr, gatewayIP string, dnsServers []string, dhcpOn bool) (*model.Network, error) {
	fabricParams := &model.FabricParams{}
	fabric, err := s.fabric.Create(ctx, fabricParams)
	if err != nil {
		return nil, err
	}

	vlan := fabric.VLANs[0]

	subnetParams := &model.SubnetParams{
		CIDR:       cidr,
		GatewayIP:  gatewayIP,
		DNSServers: dnsServers,
		Fabric:     strconv.Itoa(fabric.ID),
		VLAN:       strconv.Itoa(vlan.ID),
	}
	subnet, err := s.subnet.Create(ctx, subnetParams)
	if err != nil {
		return nil, err
	}

	if dhcpOn {
		vlanParams := &model.VLANParams{}
		if _, err := s.vlan.Update(ctx, fabric.ID, vlan.VID, vlanParams); err != nil {
			return nil, err
		}
	}

	ns, err := s.getNetworkSubnet(ctx, subnet)
	if err != nil {
		return nil, err
	}

	return &model.Network{
		Fabric: fabric,
		VLAN:   &vlan,
		Subnet: ns,
	}, nil
}

func (s *NexusService) CreateIPRange(ctx context.Context, subnetID int, startIP, endIP, comment string) (*model.IPRange, error) {
	params := &model.IPRangeParams{
		Type:    "reserved",
		Subnet:  strconv.Itoa(subnetID),
		StartIP: startIP,
		EndIP:   endIP,
		Comment: comment,
	}
	return s.ipRange.Create(ctx, params)
}

func (s *NexusService) DeleteNetwork(ctx context.Context, id int) error {
	fabric, err := s.fabric.Get(ctx, id)
	if err != nil {
		return err
	}

	subnets, err := s.subnet.List(ctx)
	if err != nil {
		return err
	}

	for i := range subnets {
		exist := false
		for j := range fabric.VLANs {
			if subnets[i].VLAN.ID != fabric.VLANs[j].ID {
				continue
			}
			exist = true
		}
		if exist {
			if err := s.subnet.Delete(ctx, subnets[i].ID); err != nil {
				return err
			}
		}
	}

	return s.fabric.Delete(ctx, id)
}

func (s *NexusService) DeleteIPRange(ctx context.Context, id int) error {
	return s.ipRange.Delete(ctx, id)
}

func (s *NexusService) UpdateFabric(ctx context.Context, id int, name string) (*model.Fabric, error) {
	params := &model.FabricParams{
		Name: name,
	}
	return s.fabric.Update(ctx, id, params)
}

func (s *NexusService) UpdateVLAN(ctx context.Context, fabricID, vid int, name string, mtu int, description string, dhcpOn bool) (*model.VLAN, error) {
	params := &model.VLANParams{
		Name:        name,
		MTU:         mtu,
		Description: description,
		DHCPOn:      dhcpOn,
	}
	return s.vlan.Update(ctx, fabricID, vid, params)
}

func (s *NexusService) UpdateSubnet(ctx context.Context, id int, name, cidr, gatewayIP string, dnsServers []string, description string, allowDNSResolution bool) (*model.NetworkSubnet, error) {
	params := &model.SubnetParams{
		Name:        name,
		CIDR:        cidr,
		GatewayIP:   gatewayIP,
		DNSServers:  dnsServers,
		Description: description,
		AllowDNS:    allowDNSResolution,
	}
	subnet, err := s.subnet.Update(ctx, id, params)
	if err != nil {
		return nil, err
	}
	return s.getNetworkSubnet(ctx, subnet)
}

func (s *NexusService) UpdateIPRange(ctx context.Context, id int, startIP, endIP, comment string) (*model.IPRange, error) {
	params := &model.IPRangeParams{
		StartIP: startIP,
		EndIP:   endIP,
		Comment: comment,
	}
	return s.ipRange.Update(ctx, id, params)
}

func (s *NexusService) getNetworkSubnet(ctx context.Context, subnet *model.Subnet) (*model.NetworkSubnet, error) {
	statistics, err := s.subnet.GetStatistics(ctx, subnet.ID)
	if err != nil {
		return nil, err
	}
	ipAddresses, err := s.subnet.GetIPAddresses(ctx, subnet.ID)
	if err != nil {
		return nil, err
	}
	allIPRanges, err := s.ipRange.List(ctx)
	if err != nil {
		return nil, err
	}
	_, ipNet, err := net.ParseCIDR(subnet.CIDR)
	if err != nil {
		return nil, err
	}
	ipRanges := []model.IPRange{}
	for i := range allIPRanges {
		if ipNet.Contains(allIPRanges[i].StartIP) && ipNet.Contains(allIPRanges[i].EndIP) {
			ipRanges = append(ipRanges, allIPRanges[i])
		}
	}
	return &model.NetworkSubnet{
		Subnet:      subnet,
		Statistics:  statistics,
		IPAddresses: ipAddresses,
		IPRanges:    ipRanges,
	}, nil
}
