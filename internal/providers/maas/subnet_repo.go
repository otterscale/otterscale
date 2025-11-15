package maas

import (
	"context"
	"strconv"

	"github.com/canonical/gomaasclient/entity"

	"github.com/otterscale/otterscale/internal/core/network"
)

// Note: MAAS API do not support context.
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

	return client.Subnets.Get()
}

func (r *subnetRepo) Get(_ context.Context, id int) (*network.Subnet, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.Subnet.Get(id)
}

func (r *subnetRepo) Create(_ context.Context, fabricID, vlanID int, cidr, gatewayIP string, dnsServers []string) (*network.Subnet, error) {
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

	return client.Subnets.Create(params)
}

func (r *subnetRepo) Update(_ context.Context, id int, name, cidr, gatewayIP string, dnsServers []string, description string, allowDNSResolution bool) (*network.Subnet, error) {
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

	return client.Subnet.Update(id, params)
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

	return client.Subnet.GetIPAddresses(id)
}

func (r *subnetRepo) GetStatistics(_ context.Context, id int) (*network.Statistics, error) {
	client, err := r.maas.Client()
	if err != nil {
		return nil, err
	}

	return client.Subnet.GetStatistics(id)
}
