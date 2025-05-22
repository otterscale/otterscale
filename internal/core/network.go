package core

import (
	"context"
	"net"
	"strconv"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/subnet"
	"github.com/openhdc/otterscale/internal/domain/model"
)

type (
	Fabric            = entity.Fabric
	FabricParams      = entity.FabricParams
	VLAN              = entity.VLAN
	VLANParams        = entity.VLANParams
	Subnet            = entity.Subnet
	SubnetParams      = entity.SubnetParams
	IPRange           = entity.IPRange
	IPRangeParams     = entity.IPRangeParams
	IPAddress         = subnet.IPAddress
	NetworkStatistics = subnet.Statistics
)

type Network struct {
	*entity.Fabric
	*entity.VLAN
	Subnet *NetworkSubnet
}

type NetworkSubnet struct {
	*Subnet
	*subnet.Statistics
	IPAddresses []subnet.IPAddress
	IPRanges    []entity.IPRange
}

type FabricRepo interface {
	List(ctx context.Context) ([]entity.Fabric, error)
	Get(ctx context.Context, id int) (*entity.Fabric, error)
	Create(ctx context.Context, params *entity.FabricParams) (*entity.Fabric, error)
	Update(ctx context.Context, id int, params *entity.FabricParams) (*entity.Fabric, error)
	Delete(ctx context.Context, id int) error
}

type VLANRepo interface {
	Update(ctx context.Context, fabricID, vid int, params *entity.VLANParams) (*entity.VLAN, error)
}

type SubnetRepo interface {
	List(ctx context.Context) ([]entity.Subnet, error)
	Get(ctx context.Context, id int) (*entity.Subnet, error)
	Create(ctx context.Context, params *entity.SubnetParams) (*entity.Subnet, error)
	Update(ctx context.Context, id int, params *entity.SubnetParams) (*entity.Subnet, error)
	Delete(ctx context.Context, id int) error
	GetIPAddresses(ctx context.Context, id int) ([]subnet.IPAddress, error)
	GetReservedIPRanges(ctx context.Context, id int) ([]subnet.ReservedIPRange, error)
	GetUnreservedIPRanges(ctx context.Context, id int) ([]subnet.IPRange, error)
	GetStatistics(ctx context.Context, id int) (*subnet.Statistics, error)
}

type IPRangeRepo interface {
	List(ctx context.Context) ([]entity.IPRange, error)
	Create(ctx context.Context, params *entity.IPRangeParams) (*entity.IPRange, error)
	Update(ctx context.Context, id int, params *entity.IPRangeParams) (*entity.IPRange, error)
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

func (uc *NetworkUseCase) ListNetworks(ctx context.Context) ([]model.Network, error) {
	subnets, err := uc.subnet.List(ctx)
	if err != nil {
		return nil, err
	}

	networkSubnets := make([]model.NetworkSubnet, len(subnets))
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

func (uc *NetworkUseCase) CreateNetwork(ctx context.Context, cidr, gatewayIP string, dnsServers []string, dhcpOn bool) (*model.Network, error) {
	fabricParams := &model.FabricParams{}
	fabric, err := uc.fabric.Create(ctx, fabricParams)
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
	subnet, err := uc.subnet.Create(ctx, subnetParams)
	if err != nil {
		return nil, err
	}

	if dhcpOn {
		vlanParams := &model.VLANParams{}
		if _, err := uc.vlan.Update(ctx, fabric.ID, vlan.VID, vlanParams); err != nil {
			return nil, err
		}
	}

	networkSubnet, err := uc.getNetworkSubnet(ctx, subnet)
	if err != nil {
		return nil, err
	}

	return &model.Network{
		Fabric: fabric,
		VLAN:   &vlan,
		Subnet: networkSubnet,
	}, nil
}

func (uc *NetworkUseCase) CreateIPRange(ctx context.Context, subnetID int, startIP, endIP, comment string) (*model.IPRange, error) {
	params := &model.IPRangeParams{
		Type:    "reserved",
		Subnet:  strconv.Itoa(subnetID),
		StartIP: startIP,
		EndIP:   endIP,
		Comment: comment,
	}
	return uc.ipRange.Create(ctx, params)
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

func (uc *NetworkUseCase) UpdateFabric(ctx context.Context, id int, name string) (*model.Fabric, error) {
	params := &model.FabricParams{
		Name: name,
	}
	return uc.fabric.Update(ctx, id, params)
}

func (uc *NetworkUseCase) UpdateVLAN(ctx context.Context, fabricID, vid int, name string, mtu int, description string, dhcpOn bool) (*model.VLAN, error) {
	params := &model.VLANParams{
		Name:        name,
		MTU:         mtu,
		Description: description,
		DHCPOn:      dhcpOn,
	}
	return uc.vlan.Update(ctx, fabricID, vid, params)
}

func (uc *NetworkUseCase) UpdateSubnet(ctx context.Context, id int, name, cidr, gatewayIP string, dnsServers []string, description string, allowDNSResolution bool) (*model.NetworkSubnet, error) {
	params := &model.SubnetParams{
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

func (uc *NetworkUseCase) UpdateIPRange(ctx context.Context, id int, startIP, endIP, comment string) (*model.IPRange, error) {
	params := &model.IPRangeParams{
		StartIP: startIP,
		EndIP:   endIP,
		Comment: comment,
	}
	return uc.ipRange.Update(ctx, id, params)
}

func (uc *NetworkUseCase) getNetworkSubnet(ctx context.Context, subnet *model.Subnet) (*model.NetworkSubnet, error) {
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
