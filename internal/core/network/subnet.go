package network

import (
	"context"
)

type Subnet struct {
	ID   int
	VLAN VLAN
	CIDR string
}

type IPAddress struct {
	ID int
}

type Statistics struct {
	ID int
}

type SubnetRepo interface {
	List(ctx context.Context) ([]Subnet, error)
	Get(ctx context.Context, id int) (*Subnet, error)
	Create(ctx context.Context, fabricID, vlanID int, cidr, gatewayIP string, dnsServers []string) (*Subnet, error)
	Update(ctx context.Context, id int, name, cidr, gatewayIP string, dnsServers []string, description string, allowDNSResolution bool) (*Subnet, error)
	Delete(ctx context.Context, id int) error
	GetIPAddresses(ctx context.Context, id int) ([]IPAddress, error)
	GetStatistics(ctx context.Context, id int) (*Statistics, error)
}

func (uc *NetworkUseCase) UpdateSubnet(ctx context.Context, id int, name, cidr, gatewayIP string, dnsServers []string, description string, allowDNSResolution bool) (*NetworkSubnet, error) {
	subnet, err := uc.subnet.Update(ctx, id, name, cidr, gatewayIP, dnsServers, description, allowDNSResolution)
	if err != nil {
		return nil, err
	}

	return uc.getNetworkSubnet(ctx, subnet)
}
