package service

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/subnet"
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/rpc/params"

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
	List(ctx context.Context) ([]*entity.PackageRepository, error)
	Update(ctx context.Context, id int, params *entity.PackageRepositoryParams) (*entity.PackageRepository, error)
}

// MAASFabric represents fabric operations
type MAASFabric interface {
	List(ctx context.Context) ([]*entity.Fabric, error)
	Get(ctx context.Context, id int) (*entity.Fabric, error)
	Create(ctx context.Context, params *entity.FabricParams) (*entity.Fabric, error)
	Update(ctx context.Context, id int, params *entity.FabricParams) (*entity.Fabric, error)
	Delete(ctx context.Context, id int) error
}

// MAASVLAN represents VLAN operations
type MAASVLAN interface {
	Update(ctx context.Context, fabricID, vid int, params *entity.VLANParams) (*entity.VLAN, error)
}

// MAASSubnet represents subnet operations
type MAASSubnet interface {
	List(ctx context.Context) ([]*entity.Subnet, error)
	Get(ctx context.Context, id int) (*entity.Subnet, error)
	Create(ctx context.Context, params *entity.SubnetParams) (*entity.Subnet, error)
	Update(ctx context.Context, id int, params *entity.SubnetParams) (*entity.Subnet, error)
	Delete(ctx context.Context, id int) error
	GetIPAddresses(ctx context.Context, id int) ([]subnet.IPAddress, error)
	GetReservedIPRanges(ctx context.Context, id int) ([]subnet.ReservedIPRange, error)
	GetUnreservedIPRanges(ctx context.Context, id int) ([]subnet.IPRange, error)
	GetStatistics(ctx context.Context, id int) (*subnet.Statistics, error)
}

// MAASIPRange represents IP range operations
type MAASIPRange interface {
	List(ctx context.Context) ([]*entity.IPRange, error)
	Create(ctx context.Context, params *entity.IPRangeParams) (*entity.IPRange, error)
	Update(ctx context.Context, id int, params *entity.IPRangeParams) (*entity.IPRange, error)
}

// MAASBootResource represents boot resource operations
type MAASBootResource interface {
	Import(ctx context.Context) error
}

// MAASMachine represents machine operations
type MAASMachine interface {
	List(ctx context.Context) ([]*entity.Machine, error)
	Get(ctx context.Context, systemID string) (*entity.Machine, error)
	PowerOn(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error)
	PowerOff(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error)
	Commission(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error)
}

type JujuClient interface {
	Status(ctx context.Context, uuid string) (*params.FullStatus, error)
}

type JujuMachine interface {
	AddMachines(ctx context.Context, uuid string, params []params.AddMachineParams) ([]params.AddMachinesResult, error)
}

type JujuModel interface {
	List(ctx context.Context) ([]*base.UserModelSummary, error)
	Create(ctx context.Context, name string) (*base.ModelInfo, error)
}

type JujuModelConfig interface {
	List(ctx context.Context, uuid string) (map[string]interface{}, error)
	Set(ctx context.Context, uuid string, config map[string]interface{}) error
	Unset(ctx context.Context, uuid string, keys ...string) error
}

type JujuApplication interface {
	ResolveUnitErrors(ctx context.Context, uuid string, units []string) error
	CreateRelation(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error)
	DeleteRelation(ctx context.Context, uuid string, id int) error
}

type JujuAction interface {
	List(ctx context.Context, uuid, appName string) (map[string]action.ActionSpec, error)
}

// StackService coordinates operations across multiple MAAS resources
type StackService struct {
	server            MAASServer
	packageRepository MAASPackageRepository
	fabric            MAASFabric
	vlan              MAASVLAN
	subnet            MAASSubnet
	ipRange           MAASIPRange
	bootResource      MAASBootResource
	client            JujuClient
	machine           MAASMachine
	jujuMachine       JujuMachine
	model             JujuModel
	modelConfig       JujuModelConfig
	application       JujuApplication
	action            JujuAction
}

// NewStackService creates a new instance of StackService
func NewStackService(
	server MAASServer,
	packageRepository MAASPackageRepository,
	fabric MAASFabric,
	vlan MAASVLAN,
	subnet MAASSubnet,
	ipRange MAASIPRange,
	bootResource MAASBootResource,
	machine MAASMachine,
	jujuMachine JujuMachine,
	client JujuClient,
	model JujuModel,
	modelConfig JujuModelConfig,
	application JujuApplication,
	action JujuAction,
) *StackService {
	return &StackService{
		server:            server,
		packageRepository: packageRepository,
		fabric:            fabric,
		vlan:              vlan,
		subnet:            subnet,
		ipRange:           ipRange,
		bootResource:      bootResource,
		machine:           machine,
		client:            client,
		jujuMachine:       jujuMachine,
		model:             model,
		modelConfig:       modelConfig,
		application:       application,
		action:            action,
	}
}

func (s *StackService) ListNTPServers(ctx context.Context) ([]string, error) {
	ntpServers, err := s.server.Get(ctx, "ntp_servers")
	if err != nil {
		return nil, err
	}
	return strings.Split(ntpServers, ","), nil
}

// UpdateNTPServers updates the NTP servers configuration
func (s *StackService) UpdateNTPServers(ctx context.Context, ntpServers []string) error {
	return s.server.Update(ctx, "ntp_servers", strings.Join(ntpServers, ","))
}

// Package repository operations

// ListPackageRepositories returns all package repositories
func (s *StackService) ListPackageRepositories(ctx context.Context) ([]*entity.PackageRepository, error) {
	return s.packageRepository.List(ctx)
}

// UpdatePackageRepositoryURL updates a package repository URL
func (s *StackService) UpdatePackageRepositoryURL(ctx context.Context, id int, url string) (*entity.PackageRepository, error) {
	// TODO: UPDATE JUJU ALSO
	params := &entity.PackageRepositoryParams{
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

	networkSubnets, err := s.getNetworkSubnets(ctx, subnets...)
	if err != nil {
		return nil, err
	}

	fabrics, err := s.fabric.List(ctx)
	if err != nil {
		return nil, err
	}

	// Convert fabrics to networks
	ret := make([]*model.Network, len(fabrics))
	for i, fabric := range fabrics {
		ret[i] = toNetwork(fabric, networkSubnets)
	}

	return ret, nil
}

// CreateNetwork creates a new network with associated resources
func (s *StackService) CreateNetwork(ctx context.Context, fabricParams *entity.FabricParams, vlanParams *entity.VLANParams, subnetParams *entity.SubnetParams, ipRangeParams *entity.IPRangeParams) (*model.Network, error) {
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
	if _, err = s.ipRange.Create(ctx, ipRangeParams); err != nil {
		return nil, err
	}

	// Update DHCP On
	if _, err := s.vlan.Update(ctx, fabric.ID, vlan.VID, vlanParams); err != nil {
		return nil, err
	}

	subnets, err := s.getNetworkSubnets(ctx, subnet)
	if err != nil {
		return nil, err
	}

	return toNetwork(fabric, subnets), nil
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

	subnets, err := s.subnet.List(ctx)
	if err != nil {
		return err
	}

	for _, subnet := range subnets {
		if err := s.subnet.Delete(ctx, subnet.ID); err != nil {
			return err
		}
	}

	// Finally delete the fabric
	return s.fabric.Delete(ctx, id)
}

// Resource update operations

// UpdateFabric updates fabric properties
func (s *StackService) UpdateFabric(ctx context.Context, id int, params *entity.FabricParams) (*entity.Fabric, error) {
	return s.fabric.Update(ctx, id, params)
}

// UpdateVLAN updates VLAN properties
func (s *StackService) UpdateVLAN(ctx context.Context, fabricID, vid int, params *entity.VLANParams) (*entity.VLAN, error) {
	return s.vlan.Update(ctx, fabricID, vid, params)
}

// UpdateSubnet updates subnet properties and returns the updated subnet with associated information
func (s *StackService) UpdateSubnet(ctx context.Context, id int, params *entity.SubnetParams) (*model.NetworkSubnet, error) {
	subnet, err := s.subnet.Update(ctx, id, params)
	if err != nil {
		return nil, err
	}
	return s.getNetworkSubnet(ctx, subnet)
}

// UpdateIPRange updates IP range properties
func (s *StackService) UpdateIPRange(ctx context.Context, id int, params *entity.IPRangeParams) (*entity.IPRange, error) {
	return s.ipRange.Update(ctx, id, params)
}

// Boot resource operations

// ImportBootResources triggers the import of boot resources
func (s *StackService) ImportBootResources(ctx context.Context) error {
	return s.bootResource.Import(ctx)
}

// Machine operations

// ListMachines returns all machines with their associated resources
func (s *StackService) ListMachines(ctx context.Context) ([]*entity.Machine, error) {
	return s.machine.List(ctx)
}

func (s *StackService) GetMachine(ctx context.Context, systemID string) (*entity.Machine, error) {
	return s.machine.Get(ctx, systemID)
}

func (s *StackService) AddMachine(ctx context.Context, uuid string, params []params.AddMachineParams) ([]string, error) {
	rs, err := s.jujuMachine.AddMachines(ctx, uuid, params)
	if err != nil {
		return nil, err
	}
	machines := make([]string, len(rs))
	for i, r := range rs {
		machines[i] = r.Machine
	}
	return machines, nil
}

// PowerOnMachine powers on a machine identified by systemID
func (s *StackService) PowerOnMachine(ctx context.Context, systemID string, params *entity.MachinePowerOnParams) (*entity.Machine, error) {
	return s.machine.PowerOn(ctx, systemID, params)
}

// PowerOffMachine powers off a machine identified by systemID
func (s *StackService) PowerOffMachine(ctx context.Context, systemID string, params *entity.MachinePowerOffParams) (*entity.Machine, error) {
	return s.machine.PowerOff(ctx, systemID, params)
}

// CommissionMachine commissions a machine identified by systemID
func (s *StackService) CommissionMachine(ctx context.Context, systemID string, params *entity.MachineCommissionParams) (*entity.Machine, error) {
	return s.machine.Commission(ctx, systemID, params)
}

func (s *StackService) ListModels(ctx context.Context) ([]*base.UserModelSummary, error) {
	return s.model.List(ctx)
}

func (s *StackService) CreateModel(ctx context.Context, name string) (*base.ModelInfo, error) {
	return s.model.Create(ctx, name)
}

func (s *StackService) GetModelConfigs(ctx context.Context, uuid string) (map[string]any, error) {
	return s.modelConfig.List(ctx, uuid)
}

func (s *StackService) SetModelConfigAPTMirror(ctx context.Context, uuid, value string) error {
	return s.modelConfig.Set(ctx, uuid, map[string]any{
		"apt-mirror": value,
	})
}

func (s *StackService) ListApplications(ctx context.Context, uuid string) (map[string]params.ApplicationStatus, error) {
	status, err := s.client.Status(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return status.Applications, nil
}

func (s *StackService) ListJujuMachines(ctx context.Context, uuid string) (map[string]params.MachineStatus, error) {
	status, err := s.client.Status(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return status.Machines, nil
}

func (s *StackService) ListIntegrations(ctx context.Context, uuid string) ([]*params.RelationStatus, error) {
	status, err := s.client.Status(ctx, uuid)
	if err != nil {
		return nil, err
	}
	ret := make([]*params.RelationStatus, len(status.Relations))
	for i := range status.Relations {
		ret[i] = &status.Relations[i]
	}
	return ret, nil
}

func (s *StackService) CreateIntegration(ctx context.Context, uuid string, endpoints []string) (*params.AddRelationResults, error) {
	return s.application.CreateRelation(ctx, uuid, endpoints)
}

func (s *StackService) DeleteIntegration(ctx context.Context, uuid string, id int) error {
	return s.application.DeleteRelation(ctx, uuid, id)
}

func (s *StackService) ListActions(ctx context.Context, uuid, appName string) (map[string]action.ActionSpec, error) {
	return s.action.List(ctx, uuid, appName)
}

func (s *StackService) getNetworkSubnet(ctx context.Context, subnet *entity.Subnet) (*model.NetworkSubnet, error) {
	ipAddresses, err := s.subnet.GetIPAddresses(ctx, subnet.ID)
	if err != nil {
		return nil, err
	}
	ipRanges, err := s.ipRange.List(ctx)
	if err != nil {
		return nil, err
	}
	_, ipNet, err := net.ParseCIDR(subnet.CIDR)
	if err != nil {
		return nil, err
	}
	reservedIPRanges := []*entity.IPRange{}
	for _, ipRange := range ipRanges {
		if ipNet.Contains(ipRange.StartIP) && ipNet.Contains(ipRange.EndIP) {
			reservedIPRanges = append(reservedIPRanges, ipRange)
		}
	}
	statistics, err := s.subnet.GetStatistics(ctx, subnet.ID)
	if err != nil {
		return nil, err
	}
	return &model.NetworkSubnet{
		Subnet:           subnet,
		IPAddresses:      ipAddresses,
		ReservedIPRanges: reservedIPRanges,
		Statistics:       statistics,
	}, nil
}

func (s *StackService) getNetworkSubnets(ctx context.Context, subnets ...*entity.Subnet) (map[int]*model.NetworkSubnet, error) {
	ret := map[int]*model.NetworkSubnet{}
	for _, subnet := range subnets {
		ns, err := s.getNetworkSubnet(ctx, subnet)
		if err != nil {
			return nil, err
		}
		ret[subnet.VLAN.ID] = ns
	}
	return ret, nil
}

// Helper functions

// toNetwork converts a fabric to a network object
func toNetwork(fabric *entity.Fabric, subnets map[int]*model.NetworkSubnet) *model.Network {
	settings := make([]*model.NetworkSetting, 0, len(fabric.VLANs))

	for i := range fabric.VLANs {
		setting := &model.NetworkSetting{
			VLAN: &fabric.VLANs[i],
		}
		if subnet, ok := subnets[fabric.VLANs[i].ID]; ok {
			setting.Subnet = subnet
		}
		settings = append(settings, setting)
	}

	return &model.Network{
		Fabric:   fabric,
		Settings: settings,
	}
}
