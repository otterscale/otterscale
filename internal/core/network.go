package core

import (
	"context"
	"net"
	"strconv"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/subnet"
)

type (
	Fabric            = entity.Fabric
	VLAN              = entity.VLAN
	Subnet            = entity.Subnet
	IPRange           = entity.IPRange
	IPAddress         = subnet.IPAddress
	NetworkStatistics = subnet.Statistics
)

type Network struct {
	*Fabric
	*VLAN
	Subnet *NetworkSubnet
}

type NetworkSubnet struct {
	*Subnet
	Statistics  *NetworkStatistics
	IPAddresses []IPAddress
	IPRanges    []IPRange
}

type FabricRepo interface {
	List(ctx context.Context) ([]Fabric, error)
	Get(ctx context.Context, id int) (*Fabric, error)
	Create(ctx context.Context, params *entity.FabricParams) (*Fabric, error)
	Update(ctx context.Context, id int, params *entity.FabricParams) (*Fabric, error)
	Delete(ctx context.Context, id int) error
}

type VLANRepo interface {
	Update(ctx context.Context, fabricID, vid int, params *entity.VLANParams) (*VLAN, error)
}

type SubnetRepo interface {
	List(ctx context.Context) ([]Subnet, error)
	Get(ctx context.Context, id int) (*Subnet, error)
	Create(ctx context.Context, params *entity.SubnetParams) (*Subnet, error)
	Update(ctx context.Context, id int, params *entity.SubnetParams) (*Subnet, error)
	Delete(ctx context.Context, id int) error
	GetIPAddresses(ctx context.Context, id int) ([]IPAddress, error)
	GetStatistics(ctx context.Context, id int) (*NetworkStatistics, error)
}

type IPRangeRepo interface {
	List(ctx context.Context) ([]IPRange, error)
	Create(ctx context.Context, params *entity.IPRangeParams) (*IPRange, error)
	Update(ctx context.Context, id int, params *entity.IPRangeParams) (*IPRange, error)
	Delete(ctx context.Context, id int) error
}

type NetworkUseCase struct {
	fabric  FabricRepo
	vlan    VLANRepo
	subnet  SubnetRepo
	ipRange IPRangeRepo
}

func NewNetworkUseCase(fabric FabricRepo, vlan VLANRepo, subnet SubnetRepo, ipRange IPRangeRepo) *NetworkUseCase {
	return &NetworkUseCase{
		fabric:  fabric,
		vlan:    vlan,
		subnet:  subnet,
		ipRange: ipRange,
	}
}

func (uc *NetworkUseCase) ListNetworks(ctx context.Context) ([]Network, error) {
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

func (uc *NetworkUseCase) CreateNetwork(ctx context.Context, cidr, gatewayIP string, dnsServers []string, dhcpOn bool) (*Network, error) {
	fabricParams := &entity.FabricParams{}
	fabric, err := uc.fabric.Create(ctx, fabricParams)
	if err != nil {
		return nil, err
	}

	vlan := fabric.VLANs[0]

	subnetParams := &entity.SubnetParams{
		CIDR:       cidr,
		GatewayIP:  gatewayIP,
		DNSServers: dnsServers,
		Fabric:     strconv.Itoa(fabric.ID),
		VLAN:       strconv.Itoa(vlan.ID),
	}
	subnet, err := uc.subnet.Create(ctx, subnetParams)
	if err != nil {
		return nil, err
	}

	if dhcpOn {
		vlanParams := &entity.VLANParams{}
		if _, err := uc.vlan.Update(ctx, fabric.ID, vlan.VID, vlanParams); err != nil {
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

func (uc *NetworkUseCase) CreateIPRange(ctx context.Context, subnetID int, startIP, endIP, comment string) (*IPRange, error) {
	return createIPRange(ctx, uc.ipRange, subnetID, startIP, endIP, comment)
}

func (uc *NetworkUseCase) DeleteNetwork(ctx context.Context, id int) error {
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

func (uc *NetworkUseCase) DeleteIPRange(ctx context.Context, id int) error {
	return uc.ipRange.Delete(ctx, id)
}

func (uc *NetworkUseCase) UpdateFabric(ctx context.Context, id int, name string) (*Fabric, error) {
	params := &entity.FabricParams{
		Name: name,
	}
	return uc.fabric.Update(ctx, id, params)
}

func (uc *NetworkUseCase) UpdateVLAN(ctx context.Context, fabricID, vid int, name string, mtu int, description string, dhcpOn bool) (*VLAN, error) {
	params := &entity.VLANParams{
		Name:        name,
		MTU:         mtu,
		Description: description,
		DHCPOn:      dhcpOn,
	}
	return uc.vlan.Update(ctx, fabricID, vid, params)
}

func (uc *NetworkUseCase) UpdateSubnet(ctx context.Context, id int, name, cidr, gatewayIP string, dnsServers []string, description string, allowDNSResolution bool) (*NetworkSubnet, error) {
	params := &entity.SubnetParams{
		Name:        name,
		CIDR:        cidr,
		GatewayIP:   gatewayIP,
		DNSServers:  dnsServers,
		Description: description,
		AllowDNS:    allowDNSResolution,
	}
	subnet, err := uc.subnet.Update(ctx, id, params)
	if err != nil {
		return nil, err
	}
	return uc.getNetworkSubnet(ctx, subnet)
}

func (uc *NetworkUseCase) UpdateIPRange(ctx context.Context, id int, startIP, endIP, comment string) (*IPRange, error) {
	params := &entity.IPRangeParams{
		StartIP: startIP,
		EndIP:   endIP,
		Comment: comment,
	}
	return uc.ipRange.Update(ctx, id, params)
}

func (uc *NetworkUseCase) getNetworkSubnet(ctx context.Context, subnet *Subnet) (*NetworkSubnet, error) {
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
