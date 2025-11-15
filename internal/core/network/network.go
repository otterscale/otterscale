package network

import (
	"context"
	"net"
)

type Network struct {
	Fabric *Fabric
	VLAN   *VLAN
	Subnet *NetworkSubnet
}

type NetworkSubnet struct {
	*Subnet
	Statistics  *Statistics
	IPAddresses []IPAddress
	IPRanges    []IPRange
}

type UseCase struct {
	fabric  FabricRepo
	subnet  SubnetRepo
	ipRange IPRangeRepo
	vlan    VLANRepo
}

func NewUseCase(fabric FabricRepo, subnet SubnetRepo, ipRange IPRangeRepo, vlan VLANRepo) *UseCase {
	return &UseCase{
		fabric:  fabric,
		subnet:  subnet,
		ipRange: ipRange,
		vlan:    vlan,
	}
}

func (uc *UseCase) ListNetworks(ctx context.Context) ([]Network, error) {
	subnets, err := uc.subnet.List(ctx)
	if err != nil {
		return nil, err
	}

	networkSubnets := make([]NetworkSubnet, len(subnets))
	for i := range subnets {
		ns, err := uc.getNetworkSubnet(ctx, &subnets[i])
		if err != nil {
			return nil, err
		}
		networkSubnets[i] = *ns
	}

	fabrics, err := uc.fabric.List(ctx)
	if err != nil {
		return nil, err
	}

	networks := []Network{}
	for i := range fabrics {
		for j := range fabrics[i].VLANs {
			exist := false

			for k := range networkSubnets {
				if networkSubnets[k].VLAN.ID != fabrics[i].VLANs[j].ID { // not only one
					continue
				}

				networks = append(networks, Network{
					Fabric: &fabrics[i],
					VLAN:   &fabrics[i].VLANs[j],
					Subnet: &networkSubnets[k],
				})

				exist = true
			}

			if !exist {
				networks = append(networks, Network{
					Fabric: &fabrics[i],
					VLAN:   &fabrics[i].VLANs[j],
				})
			}
		}
	}

	return networks, nil
}

func (uc *UseCase) CreateNetwork(ctx context.Context, cidr, gatewayIP string, dnsServers []string, dhcpOn bool) (*Network, error) {
	fabric, err := uc.fabric.Create(ctx)
	if err != nil {
		return nil, err
	}

	vlan := fabric.VLANs[0] // default VLAN

	subnet, err := uc.subnet.Create(ctx, fabric.ID, vlan.ID, cidr, gatewayIP, dnsServers)
	if err != nil {
		return nil, err
	}

	if dhcpOn {
		if _, err := uc.vlan.Update(ctx, fabric.ID, vlan.VID, vlan.Name, vlan.MTU, vlan.Description, true); err != nil {
			return nil, err
		}
	}

	networkSubnet, err := uc.getNetworkSubnet(ctx, subnet)
	if err != nil {
		return nil, err
	}

	return &Network{
		Fabric: fabric,
		VLAN:   &vlan,
		Subnet: networkSubnet,
	}, nil
}

func (uc *UseCase) DeleteNetwork(ctx context.Context, id int) error {
	fabric, err := uc.fabric.Get(ctx, id)
	if err != nil {
		return err
	}

	subnets, err := uc.subnet.List(ctx)
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
			if err := uc.subnet.Delete(ctx, subnets[i].ID); err != nil {
				return err
			}
		}
	}

	return uc.fabric.Delete(ctx, id)
}

func (uc *UseCase) getNetworkSubnet(ctx context.Context, subnet *Subnet) (*NetworkSubnet, error) {
	statistics, err := uc.subnet.GetStatistics(ctx, subnet.ID)
	if err != nil {
		return nil, err
	}

	ipAddresses, err := uc.subnet.GetIPAddresses(ctx, subnet.ID)
	if err != nil {
		return nil, err
	}

	allIPRanges, err := uc.ipRange.List(ctx)
	if err != nil {
		return nil, err
	}

	_, ipNet, err := net.ParseCIDR(subnet.CIDR)
	if err != nil {
		return nil, err
	}

	ipRanges := []IPRange{}

	for i := range allIPRanges {
		if ipNet.Contains(allIPRanges[i].StartIP) && ipNet.Contains(allIPRanges[i].EndIP) {
			ipRanges = append(ipRanges, allIPRanges[i])
		}
	}

	return &NetworkSubnet{
		Subnet:      subnet,
		Statistics:  statistics,
		IPAddresses: ipAddresses,
		IPRanges:    ipRanges,
	}, nil
}
