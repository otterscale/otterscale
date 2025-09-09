package core

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/stretchr/testify/assert"
)

// Mock FabricRepo
type mockFabricRepo struct {
	fabrics []Fabric
}

func (m *mockFabricRepo) List(ctx context.Context) ([]Fabric, error) {
	return m.fabrics, nil
}

func (m *mockFabricRepo) Get(ctx context.Context, id int) (*Fabric, error) {
	for _, f := range m.fabrics {
		if f.ID == id {
			return &f, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockFabricRepo) Create(ctx context.Context, params *entity.FabricParams) (*Fabric, error) {
	f := Fabric{ID: 3, Name: params.Name, VLANs: []entity.VLAN{{ID: 30, VID: 300}}}
	m.fabrics = append(m.fabrics, f)
	return &f, nil
}

func (m *mockFabricRepo) Update(ctx context.Context, id int, params *entity.FabricParams) (*Fabric, error) {
	return &Fabric{ID: id, Name: params.Name}, nil
}

func (m *mockFabricRepo) Delete(ctx context.Context, id int) error {
	return nil
}

// Mock VLANRepo
type mockVLANRepo struct{}

func (m *mockVLANRepo) Update(ctx context.Context, fabricID, vid int, params *entity.VLANParams) (*VLAN, error) {
	return &VLAN{ID: vid, Name: params.Name, MTU: params.MTU, Description: params.Description, DHCPOn: params.DHCPOn}, nil
}

// Mock SubnetRepo
type mockSubnetRepo struct {
	subnets []Subnet
}

func (m *mockSubnetRepo) List(ctx context.Context) ([]Subnet, error) {
	return m.subnets, nil
}

func (m *mockSubnetRepo) Get(ctx context.Context, id int) (*Subnet, error) {
	for _, s := range m.subnets {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockSubnetRepo) Create(ctx context.Context, params *entity.SubnetParams) (*Subnet, error) {
	return &Subnet{ID: 10, CIDR: params.CIDR, VLAN: entity.VLAN{ID: 20}}, nil
}

func (m *mockSubnetRepo) Update(ctx context.Context, id int, params *entity.SubnetParams) (*Subnet, error) {
	return &Subnet{ID: id, Name: params.Name, CIDR: params.CIDR}, nil
}

func (m *mockSubnetRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func (m *mockSubnetRepo) GetIPAddresses(ctx context.Context, id int) ([]IPAddress, error) {
	return []IPAddress{{IP: net.ParseIP("192.168.1.10")}}, nil
}

func (m *mockSubnetRepo) GetStatistics(ctx context.Context, id int) (*NetworkStatistics, error) {
	return &NetworkStatistics{}, nil
}

// Mock IPRangeRepo
type mockIPRangeRepo struct{}

func (m *mockIPRangeRepo) List(ctx context.Context) ([]IPRange, error) {
	return []IPRange{
		{StartIP: net.ParseIP("192.168.1.100"), EndIP: net.ParseIP("192.168.1.200")},
	}, nil
}

func (m *mockIPRangeRepo) Create(ctx context.Context, params *entity.IPRangeParams) (*IPRange, error) {
	return &IPRange{StartIP: net.ParseIP(params.StartIP), EndIP: net.ParseIP(params.EndIP), Comment: params.Comment}, nil
}

func (m *mockIPRangeRepo) Update(ctx context.Context, id int, params *entity.IPRangeParams) (*IPRange, error) {
	return &IPRange{StartIP: net.ParseIP(params.StartIP), EndIP: net.ParseIP(params.EndIP), Comment: params.Comment}, nil
}

func (m *mockIPRangeRepo) Delete(ctx context.Context, id int) error {
	return nil
}

func TestNetworkUseCase_ListNetworks(t *testing.T) {
	fabricRepo := &mockFabricRepo{
		fabrics: []Fabric{
			{
				ID:   1,
				Name: "fabric1",
				VLANs: []entity.VLAN{
					{ID: 10, VID: 100, Name: "vlan100"},
				},
			},
		},
	}
	subnetRepo := &mockSubnetRepo{
		subnets: []Subnet{
			{ID: 1, CIDR: "192.168.1.0/24", VLAN: entity.VLAN{ID: 10, VID: 100, Name: "vlan100"}},
		},
	}
	uc := NewNetworkUseCase(fabricRepo, &mockVLANRepo{}, subnetRepo, &mockIPRangeRepo{})
	networks, err := uc.ListNetworks(context.Background())
	assert.NoError(t, err)
	assert.Len(t, networks, 1)
	assert.Equal(t, "fabric1", networks[0].Fabric.Name)
	assert.Equal(t, "vlan100", networks[0].VLAN.Name)
	assert.NotNil(t, networks[0].Subnet)
}

func TestNetworkUseCase_CreateNetwork(t *testing.T) {
	fabricRepo := &mockFabricRepo{}
	subnetRepo := &mockSubnetRepo{}
	uc := NewNetworkUseCase(fabricRepo, &mockVLANRepo{}, subnetRepo, &mockIPRangeRepo{})
	network, err := uc.CreateNetwork(context.Background(), "10.0.0.0/24", "10.0.0.1", []string{"8.8.8.8"}, true)
	assert.NoError(t, err)
	assert.NotNil(t, network)
	assert.Equal(t, "10.0.0.0/24", network.Subnet.CIDR)
}

func TestNetworkUseCase_UpdateFabric(t *testing.T) {
	fabricRepo := &mockFabricRepo{}
	uc := NewNetworkUseCase(fabricRepo, &mockVLANRepo{}, &mockSubnetRepo{}, &mockIPRangeRepo{})
	fabric, err := uc.UpdateFabric(context.Background(), 1, "newname")
	assert.NoError(t, err)
	assert.Equal(t, "newname", fabric.Name)
}

func TestNetworkUseCase_UpdateVLAN(t *testing.T) {
	uc := NewNetworkUseCase(&mockFabricRepo{}, &mockVLANRepo{}, &mockSubnetRepo{}, &mockIPRangeRepo{})
	vlan, err := uc.UpdateVLAN(context.Background(), 1, 2, "vlan2", 1500, "desc", true)
	assert.NoError(t, err)
	assert.Equal(t, "vlan2", vlan.Name)
	assert.Equal(t, 1500, vlan.MTU)
	assert.Equal(t, "desc", vlan.Description)
	assert.True(t, vlan.DHCPOn)
}

func TestNetworkUseCase_UpdateSubnet(t *testing.T) {
	subnetRepo := &mockSubnetRepo{}
	uc := NewNetworkUseCase(&mockFabricRepo{}, &mockVLANRepo{}, subnetRepo, &mockIPRangeRepo{})
	ns, err := uc.UpdateSubnet(context.Background(), 1, "subnet1", "10.0.1.0/24", "10.0.1.1", []string{"1.1.1.1"}, "desc", true)
	assert.NoError(t, err)
	assert.Equal(t, "subnet1", ns.Name)
	assert.Equal(t, "10.0.1.0/24", ns.CIDR)
}

func TestNetworkUseCase_CreateIPRange(t *testing.T) {
	uc := NewNetworkUseCase(&mockFabricRepo{}, &mockVLANRepo{}, &mockSubnetRepo{}, &mockIPRangeRepo{})
	ipr, err := uc.CreateIPRange(context.Background(), 1, "192.168.1.100", "192.168.1.200", "test")
	assert.NoError(t, err)
	assert.Equal(t, "test", ipr.Comment)
}

func TestNetworkUseCase_DeleteNetwork(t *testing.T) {
	fabricRepo := &mockFabricRepo{
		fabrics: []Fabric{
			{
				ID:   1,
				Name: "fabric1",
				VLANs: []entity.VLAN{
					{ID: 10, VID: 100, Name: "vlan100"},
				},
			},
		},
	}
	subnetRepo := &mockSubnetRepo{
		subnets: []Subnet{
			{ID: 1, CIDR: "192.168.1.0/24", VLAN: entity.VLAN{ID: 10, VID: 100, Name: "vlan100"}},
		},
	}
	uc := NewNetworkUseCase(fabricRepo, &mockVLANRepo{}, subnetRepo, &mockIPRangeRepo{})
	err := uc.DeleteNetwork(context.Background(), 1)
	assert.NoError(t, err)
}

func TestNetworkUseCase_DeleteIPRange(t *testing.T) {
	uc := NewNetworkUseCase(&mockFabricRepo{}, &mockVLANRepo{}, &mockSubnetRepo{}, &mockIPRangeRepo{})
	err := uc.DeleteIPRange(context.Background(), 1)
	assert.NoError(t, err)
}

func TestNetworkUseCase_UpdateIPRange(t *testing.T) {
	uc := NewNetworkUseCase(&mockFabricRepo{}, &mockVLANRepo{}, &mockSubnetRepo{}, &mockIPRangeRepo{})
	ipr, err := uc.UpdateIPRange(context.Background(), 1, "192.168.1.101", "192.168.1.201", "update")
	assert.NoError(t, err)
	assert.Equal(t, "update", ipr.Comment)
}
