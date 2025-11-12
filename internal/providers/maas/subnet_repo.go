package maas

import (
	"context"
	"strconv"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/subnet"

	"github.com/otterscale/otterscale/internal/core/network"
)

type subnetRepo struct {
	maas *MAAS
}

func NewSubnetRepo(maas *MAAS) network.SubnetRepo {
	return &subnetRepo{
		maas: maas,
	}
}

var _ network.SubnetRepo = (*subnetRepo)(nil)

func (r *subnetRepo) List(_ context.Context) ([]network.Subnet, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	subnets, err := client.Subnets.Get()
	if err != nil {
		return nil, err
	}

	return r.toSubnets(subnets), nil
}

func (r *subnetRepo) Get(_ context.Context, id int) (*network.Subnet, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	subnet, err := client.Subnet.Get(id)
	if err != nil {
		return nil, err
	}

	return r.toSubnet(subnet), nil
}

func (r *subnetRepo) Create(ctx context.Context, fabricID, vlanID int, cidr, gatewayIP string, dnsServers []string) (*network.Subnet, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.SubnetParams{
		Fabric:     strconv.Itoa(fabricID),
		VLAN:       strconv.Itoa(vlanID),
		CIDR:       cidr,
		GatewayIP:  gatewayIP,
		DNSServers: dnsServers,
	}

	subnet, err := client.Subnets.Create(params)
	if err != nil {
		return nil, err
	}

	return r.toSubnet(subnet), nil
}

func (r *subnetRepo) Update(ctx context.Context, id int, name, cidr, gatewayIP string, dnsServers []string, description string, allowDNSResolution bool) (*network.Subnet, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	params := &entity.SubnetParams{
		Name:        name,
		CIDR:        cidr,
		GatewayIP:   gatewayIP,
		DNSServers:  dnsServers,
		Description: description,
		AllowDNS:    allowDNSResolution,
	}

	subnet, err := client.Subnet.Update(id, params)
	if err != nil {
		return nil, err
	}

	return r.toSubnet(subnet), nil
}

func (r *subnetRepo) Delete(_ context.Context, id int) error {
	client, err := r.maas.Client()
	if err != nil {
		return err
	}

	return client.Subnet.Delete(id)
}

func (r *subnetRepo) GetIPAddresses(_ context.Context, id int) ([]network.IPAddress, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	ips, err := client.Subnet.GetIPAddresses(id)
	if err != nil {
		return nil, err
	}

	return r.toIPAddresses(ips), nil
}

func (r *subnetRepo) GetStatistics(_ context.Context, id int) (*network.Statistics, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	statistics, err := client.Subnet.GetStatistics(id)
	if err != nil {
		return nil, err
	}

	return r.toStatistics(statistics), nil
}

func (r *subnetRepo) toSubnet(s *entity.Subnet) *network.Subnet {
	return &network.Subnet{}
}

func (r *subnetRepo) toSubnets(ss []entity.Subnet) []network.Subnet {
	ret := make([]network.Subnet, 0, len(ss))

	return ret
}

func (r *subnetRepo) toIPAddress(ip *subnet.IPAddress) *network.IPAddress {
	return &network.IPAddress{}
}

func (r *subnetRepo) toIPAddresses(ips []subnet.IPAddress) []network.IPAddress {
	ret := make([]network.IPAddress, 0, len(ips))

	return ret
}

func (r *subnetRepo) toStatistics(s *subnet.Statistics) *network.Statistics {
	return &network.Statistics{}
}
