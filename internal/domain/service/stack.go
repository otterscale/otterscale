package service

import (
	"context"
	"errors"

	"github.com/openhdc/openhdc/internal/domain/model"
)

// MAAS API interfaces grouped by resource type

// MAASServer represents MAAS server configuration operations
type MAASServer interface {
	Get(ctx context.Context, name string) (string, error)
	Update(ctx context.Context, name, value string) error
}

// MAASPackageRepository represents package repository operations
type MAASPackageRepository interface {
	List(ctx context.Context) ([]*model.PackageRepository, error)
	Update(ctx context.Context, id int, params *model.PackageRepositoryParams) (*model.PackageRepository, error)
}

// MAASNetworkingInterfaces groups all networking-related interfaces
type MAASNetworking struct {
	Fabric  MAASFabric
	VLAN    MAASVLAN
	Subnet  MAASSubnet
	IPRange MAASIPRange
}

// MAASFabric represents fabric operations
type MAASFabric interface {
	List(ctx context.Context) ([]*model.Fabric, error)
	Get(ctx context.Context, id int) (*model.Fabric, error)
	Create(ctx context.Context, params *model.FabricParams) (*model.Fabric, error)
	Update(ctx context.Context, id int, params *model.FabricParams) (*model.Fabric, error)
	Delete(ctx context.Context, id int) error
}

// MAASVLAN represents VLAN operations
type MAASVLAN interface {
	Update(ctx context.Context, fabricID, vid int, params *model.VLANParams) (*model.VLAN, error)
}

// MAASSubnet represents subnet operations
type MAASSubnet interface {
	List(ctx context.Context) ([]*model.Subnet, error)
	Get(ctx context.Context, id int) (*model.Subnet, error)
	Create(ctx context.Context, params *model.SubnetParams) (*model.Subnet, error)
	Update(ctx context.Context, id int, params *model.SubnetParams) (*model.Subnet, error)
	Delete(ctx context.Context, id int) error
}

// MAASIPRange represents IP range operations
type MAASIPRange interface {
	List(ctx context.Context) ([]*model.IPRange, error)
	Create(ctx context.Context, params *model.IPRangeParams) (*model.IPRange, error)
	Update(ctx context.Context, id int, params *model.IPRangeParams) (*model.IPRange, error)
}

// MAASMachine represents machine operations (placeholder for future implementation)
type MAASMachine interface {
}

// StackService coordinates operations across multiple MAAS resources
type StackService struct {
	server            MAASServer
	packageRepository MAASPackageRepository
	fabric            MAASFabric
	vlan              MAASVLAN
	subnet            MAASSubnet
	ipRange           MAASIPRange
}

// NewStackService creates a new instance of StackService
func NewStackService(server MAASServer, packageRepository MAASPackageRepository, fabric MAASFabric, vlan MAASVLAN, subnet MAASSubnet, ipRange MAASIPRange) *StackService {
	return &StackService{
		server:            server,
		packageRepository: packageRepository,
		fabric:            fabric,
		vlan:              vlan,
		subnet:            subnet,
		ipRange:           ipRange,
	}
}

// UpdateNTPServers updates the NTP servers configuration
func (s *StackService) UpdateNTPServers(ctx context.Context, ntpServers string) error {
	return s.server.Update(ctx, "ntp_servers", ntpServers)
}

// Package repository operations

// ListPackageRepositories returns all package repositories
func (s *StackService) ListPackageRepositories(ctx context.Context) ([]*model.PackageRepository, error) {
	return s.packageRepository.List(ctx)
}

// UpdatePackageRepositoryURL updates a package repository URL
func (s *StackService) UpdatePackageRepositoryURL(ctx context.Context, id int, url string) (*model.PackageRepository, error) {
	// TODO: UPDATE JUJU ALSO
	params := &model.PackageRepositoryParams{
		URL: url,
	}
	return s.packageRepository.Update(ctx, id, params)
}

// Network operations

// ListNetworks returns all networks with their associated resources
func (s *StackService) ListNetworks(ctx context.Context) ([]*model.Network, error) {
	// Get all required resources
	subnets, err := s.subnet.List(ctx)
	if err != nil {
		return nil, err
	}

	ipRanges, err := s.ipRange.List(ctx)
	if err != nil {
		return nil, err
	}

	fabrics, err := s.fabric.List(ctx)
	if err != nil {
		return nil, err
	}

	// Build mapping of VLANs to network settings
	vlanToNetworkSetting := vlanToNetworkSettingMap(subnets, ipRanges)

	// Convert fabrics to networks
	networks := make([]*model.Network, len(fabrics))
	for i, fabric := range fabrics {
		networks[i] = toNetwork(fabric, vlanToNetworkSetting)
	}

	return networks, nil
}

// CreateNetwork creates a new network with associated resources
func (s *StackService) CreateNetwork(ctx context.Context, fabricParams *model.FabricParams, subnetParams *model.SubnetParams, ipRangeParams *model.IPRangeParams) (*model.Network, error) {
	// Create fabric first
	fabric, err := s.fabric.Create(ctx, fabricParams)
	if err != nil {
		return nil, err
	}

	// Ensure fabric has at least one VLAN
	if len(fabric.VLANs) == 0 {
		return nil, errors.New("created fabric has no VLANs")
	}

	// Use the first VLAN for subnet creation
	vlan := &fabric.VLANs[0]

	// Set subnet parameters based on created fabric
	subnetParams.Fabric = fabric.Name
	subnetParams.VLAN = vlan.Name

	// Create subnet on the default VLAN
	subnet, err := s.subnet.Create(ctx, subnetParams)
	if err != nil {
		return nil, err
	}

	// Create IP range for the subnet
	ipRangeParams.Subnet = subnet.Name
	ipRange, err := s.ipRange.Create(ctx, ipRangeParams)
	if err != nil {
		return nil, err
	}

	// Build network model from created resources
	vlanToNetworkSetting := map[int]*model.NetworkSetting{
		vlan.ID: {
			VLAN:    vlan,
			Subnet:  subnet,
			IPRange: ipRange,
		},
	}

	return toNetwork(fabric, vlanToNetworkSetting), nil
}

// DeleteNetwork deletes a network and all associated resources
func (s *StackService) DeleteNetwork(ctx context.Context, id int) error {
	// Get fabric first to verify it exists
	fabric, err := s.fabric.Get(ctx, id)
	if err != nil {
		return err
	}

	// No VLANs, just delete the fabric
	if len(fabric.VLANs) == 0 {
		return s.fabric.Delete(ctx, id)
	}

	// Get network resources to identify what needs to be deleted
	networkSettings, err := s.getNetworkSettingsForFabric(ctx, fabric)
	if err != nil {
		return err
	}

	// Delete all associated subnets first
	for _, ns := range networkSettings {
		if ns.Subnet != nil {
			if err := s.subnet.Delete(ctx, ns.Subnet.ID); err != nil {
				return err
			}
		}
	}

	// Finally delete the fabric
	return s.fabric.Delete(ctx, id)
}

// getNetworkSettingsForFabric retrieves network settings for a fabric
func (s *StackService) getNetworkSettingsForFabric(ctx context.Context, fabric *model.Fabric) (map[int]*model.NetworkSetting, error) {
	subnets, err := s.subnet.List(ctx)
	if err != nil {
		return nil, err
	}

	ipRanges, err := s.ipRange.List(ctx)
	if err != nil {
		return nil, err
	}

	return vlanToNetworkSettingMap(subnets, ipRanges), nil
}

// Resource update operations

// UpdateFabric updates fabric properties
func (s *StackService) UpdateFabric(ctx context.Context, id int, params *model.FabricParams) (*model.Fabric, error) {
	return s.fabric.Update(ctx, id, params)
}

// UpdateVLAN updates VLAN properties
func (s *StackService) UpdateVLAN(ctx context.Context, fabricID, vid int, params *model.VLANParams) (*model.VLAN, error) {
	return s.vlan.Update(ctx, fabricID, vid, params)
}

// UpdateSubnet updates subnet properties
func (s *StackService) UpdateSubnet(ctx context.Context, id int, params *model.SubnetParams) (*model.Subnet, error) {
	return s.subnet.Update(ctx, id, params)
}

// UpdateIPRange updates IP range properties
func (s *StackService) UpdateIPRange(ctx context.Context, id int, params *model.IPRangeParams) (*model.IPRange, error) {
	return s.ipRange.Update(ctx, id, params)
}

// Helper functions

// vlanToNetworkSettingMap creates a mapping of VLAN IDs to network settings
func vlanToNetworkSettingMap(subnets []*model.Subnet, ipRanges []*model.IPRange) map[int]*model.NetworkSetting {
	vlanToNetworkSetting := make(map[int]*model.NetworkSetting, len(subnets))

	// Create subnet lookup by name for faster IPRange matching
	subnetsByName := make(map[string]*model.Subnet)
	for i := range subnets {
		subnetsByName[subnets[i].Name] = subnets[i]
	}

	// Map IP ranges to subnets
	ipRangeBySubnet := make(map[string]*model.IPRange)
	for i := range ipRanges {
		if ipRanges[i].Subnet.Name != "" {
			ipRangeBySubnet[ipRanges[i].Subnet.Name] = ipRanges[i]
		}
	}

	// Build the network settings map
	for i := range subnets {
		subnet := subnets[i]
		vlan := &subnet.VLAN
		ipRange := ipRangeBySubnet[subnet.Name]

		vlanToNetworkSetting[vlan.ID] = &model.NetworkSetting{
			VLAN:    vlan,
			Subnet:  subnet,
			IPRange: ipRange,
		}
	}

	return vlanToNetworkSetting
}

// toNetwork converts a fabric to a network object
func toNetwork(fabric *model.Fabric, vlanToNetworkSetting map[int]*model.NetworkSetting) *model.Network {
	settings := make([]*model.NetworkSetting, 0, len(fabric.VLANs))

	for i := range fabric.VLANs {
		if setting, ok := vlanToNetworkSetting[fabric.VLANs[i].ID]; ok {
			settings = append(settings, setting)
			continue
		}
		// Add VLAN without subnet/IPRange
		settings = append(settings, &model.NetworkSetting{
			VLAN: &fabric.VLANs[i],
		})
	}

	return &model.Network{
		ID:          fabric.ID,
		Name:        fabric.Name,
		ClassType:   fabric.ClassType,
		ResourceURI: fabric.ResourceURI,
		Settings:    settings,
	}
}
