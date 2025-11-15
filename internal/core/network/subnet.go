package network

import (
	"context"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/subnet"
)

// Subnet represents a MAAS Subnet resource.
type Subnet = entity.Subnet

// IPAddress represents a MAAS IPAddress resource.
type IPAddress = subnet.IPAddress

// Statistics represents MAAS Statistics statistics.
type Statistics = subnet.Statistics

type SubnetRepo interface {
	List(ctx context.Context) ([]Subnet, error)
	Get(ctx context.Context, id int) (*Subnet, error)
	Create(ctx context.Context, fabricID, vlanID int, cidr, gatewayIP string, dnsServers []string) (*Subnet, error)
	Update(ctx context.Context, id int, name, cidr, gatewayIP string, dnsServers []string, description string, allowDNSResolution bool) (*Subnet, error)
	Delete(ctx context.Context, id int) error
	GetIPAddresses(ctx context.Context, id int) ([]IPAddress, error)
	GetStatistics(ctx context.Context, id int) (*Statistics, error)
}

func (uc *UseCase) UpdateSubnet(ctx context.Context, id int, name, cidr, gatewayIP string, dnsServers []string, description string, allowDNSResolution bool) (*SubnetData, error) {
	subnet, err := uc.subnet.Update(ctx, id, name, cidr, gatewayIP, dnsServers, description, allowDNSResolution)
	if err != nil {
		return nil, err
	}

	return uc.getSubnetData(ctx, subnet)
}
