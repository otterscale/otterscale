package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/openhdc/openhdc/api/stack/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
	"github.com/openhdc/openhdc/internal/domain/service"
)

// StackApp implements the StackServiceServer interface
type StackApp struct {
	v1.UnimplementedStackServiceServer
	svc *service.StackService
}

// NewStackApp creates a new StackApp instance
func NewStackApp(svc *service.StackService) *StackApp {
	return &StackApp{svc: svc}
}

// Ensure StackApp implements the StackServiceServer interface
var _ v1.StackServiceServer = (*StackApp)(nil)

// UpdateNTPServers updates NTP server configuration
func (a *StackApp) UpdateNTPServers(ctx context.Context, req *v1.UpdateNTPServersRequest) (*emptypb.Empty, error) {
	if err := a.svc.UpdateNTPServers(ctx, req.GetNtpServers()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListPackageRepositories retrieves package repositories
func (a *StackApp) ListPackageRepositories(ctx context.Context, req *v1.ListPackageRepositoriesRequest) (*v1.ListPackageRepositoriesResponse, error) {
	pageSize, err := getPageSize(req.GetPageSize())
	if err != nil {
		return nil, err
	}
	pageToken, err := getPageToken(req.GetPageToken())
	if err != nil {
		return nil, err
	}
	prs, nextPageToken, err := a.svc.ListPackageRepositories(ctx, pageSize, pageToken)
	if err != nil {
		return nil, err
	}
	ret := &v1.ListPackageRepositoriesResponse{}
	ret.SetPackageRepositories(toPackageRepositories(prs))
	ret.SetNextPageToken(nextPageToken)
	return ret, nil
}

// UpdatePackageRepositoryURL updates a package repository URL
func (a *StackApp) UpdatePackageRepositoryURL(ctx context.Context, req *v1.UpdatePackageRepositoryURLRequest) (*v1.PackageRepository, error) {
	pr, err := a.svc.UpdatePackageRepositoryURL(ctx, int(req.GetId()), req.GetUrl())
	if err != nil {
		return nil, err
	}
	return toPackageRepository(pr), nil
}

// Network management functions
func (a *StackApp) ListNetworks(ctx context.Context, req *v1.ListNetworksRequest) (*v1.ListNetworksResponse, error) {
	pageSize, err := getPageSize(req.GetPageSize())
	if err != nil {
		return nil, err
	}
	pageToken, err := getPageToken(req.GetPageToken())
	if err != nil {
		return nil, err
	}
	ns, nextPageToken, err := a.svc.ListNetworks(ctx, pageSize, pageToken)
	if err != nil {
		return nil, err
	}
	ret := &v1.ListNetworksResponse{}
	ret.SetNetworks(toNetworks(ns))
	ret.SetNextPageToken(nextPageToken)
	return ret, nil
}

func (a *StackApp) CreateNetwork(ctx context.Context, req *v1.CreateNetworkRequest) (*v1.Network, error) {
	fabricParams := &model.FabricParams{}
	vlanParams := &model.VLANParams{
		DHCPOn: req.GetDhcpOn(),
	}
	subnetParams := &model.SubnetParams{
		CIDR:       req.GetCidr(),
		GatewayIP:  req.GetGatewayIp(),
		DNSServers: req.GetDnsServers(),
	}
	ipRangeParams := &model.IPRangeParams{
		StartIP: req.GetStartIp(),
		EndIP:   req.GetEndIp(),
	}
	n, err := a.svc.CreateNetwork(ctx, fabricParams, vlanParams, subnetParams, ipRangeParams)
	if err != nil {
		return nil, err
	}
	return toNetwork(n), nil
}

func (a *StackApp) DeleteNetwork(ctx context.Context, req *v1.DeleteNetworkRequest) (*emptypb.Empty, error) {
	if err := a.svc.DeleteNetwork(ctx, int(req.GetFabricId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Network configuration functions
func (a *StackApp) UpdateFabric(ctx context.Context, req *v1.UpdateFabricRequest) (*v1.Fabric, error) {
	params := &model.FabricParams{
		Name: req.GetName(),
	}
	f, err := a.svc.UpdateFabric(ctx, int(req.GetId()), params)
	if err != nil {
		return nil, err
	}
	return toFabric(f), nil
}

func (a *StackApp) UpdateVLAN(ctx context.Context, req *v1.UpdateVLANRequest) (*v1.VLAN, error) {
	params := &model.VLANParams{
		Name:        req.GetName(),
		MTU:         int(req.GetMtu()),
		Description: req.GetDescription(),
		DHCPOn:      req.GetDhcpOn(),
	}
	v, err := a.svc.UpdateVLAN(ctx, int(req.GetFabricId()), int(req.GetVid()), params)
	if err != nil {
		return nil, err
	}
	return toVLAN(v), nil
}

func (a *StackApp) UpdateSubnet(ctx context.Context, req *v1.UpdateSubnetRequest) (*v1.Subnet, error) {
	params := &model.SubnetParams{
		Name:        req.GetName(),
		CIDR:        req.GetCidr(),
		GatewayIP:   req.GetGatewayIp(),
		DNSServers:  req.GetDnsServers(),
		Description: req.GetDescription(),
		AllowDNS:    req.GetAllowDnsResolution(),
	}
	s, err := a.svc.UpdateSubnet(ctx, int(req.GetId()), params)
	if err != nil {
		return nil, err
	}
	return toSubnet(s), nil
}

func (a *StackApp) UpdateIPRange(ctx context.Context, req *v1.UpdateIPRangeRequest) (*v1.IPRange, error) {
	params := &model.IPRangeParams{
		StartIP: req.GetStartIp(),
		EndIP:   req.GetEndIp(),
		Comment: req.GetComment(),
	}
	r, err := a.svc.UpdateIPRange(ctx, int(req.GetId()), params)
	if err != nil {
		return nil, err
	}
	return toIPRange(r), nil
}

// ListMachines retrieves machines
func (a *StackApp) ListMachines(ctx context.Context, req *v1.ListMachinesRequest) (*v1.ListMachinesResponse, error) {
	pageSize, err := getPageSize(req.GetPageSize())
	if err != nil {
		return nil, err
	}
	pageToken, err := getPageToken(req.GetPageToken())
	if err != nil {
		return nil, err
	}
	ms, nextPageToken, err := a.svc.ListMachines(ctx, pageSize, pageToken)
	if err != nil {
		return nil, err
	}
	ret := &v1.ListMachinesResponse{}
	ret.SetMachines(toMachines(ms))
	ret.SetNextPageToken(nextPageToken)
	return ret, nil
}

// Boot resource management
func (a *StackApp) ImportBootResources(ctx context.Context, req *v1.ImportBootResourcesRequest) (*emptypb.Empty, error) {
	if err := a.svc.ImportBootResources(ctx); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Machine power management functions
func (a *StackApp) PowerOnMachine(ctx context.Context, req *v1.PowerOnMachineRequest) (*v1.Machine, error) {
	params := &model.MachinePowerOnParams{}
	m, err := a.svc.PowerOnMachine(ctx, req.GetSystemId(), params)
	if err != nil {
		return nil, err
	}
	return toMachine(m), nil
}

func (a *StackApp) PowerOffMachine(ctx context.Context, req *v1.PowerOffMachineRequest) (*v1.Machine, error) {
	params := &model.MachinePowerOffParams{}
	m, err := a.svc.PowerOffMachine(ctx, req.GetSystemId(), params)
	if err != nil {
		return nil, err
	}
	return toMachine(m), nil
}

// Machine provisioning functions
func (a *StackApp) CommissionMachine(ctx context.Context, req *v1.CommissionMachineRequest) (*v1.Machine, error) {
	params := &model.MachineCommissionParams{
		EnableSSH:      boolToInt(req.GetEnableSsh()),
		SkipBMCConfig:  boolToInt(req.GetSkipBmcConfig()),
		SkipNetworking: boolToInt(req.GetSkipNetworking()),
		SkipStorage:    boolToInt(req.GetSkipStorage()),
	}
	m, err := a.svc.CommissionMachine(ctx, req.GetSystemId(), params)
	if err != nil {
		return nil, err
	}
	return toMachine(m), nil
}

func toPackageRepositories(prs []*model.PackageRepository) []*v1.PackageRepository {
	ret := make([]*v1.PackageRepository, len(prs))
	for i := range prs {
		ret[i] = toPackageRepository(prs[i])
	}
	return ret
}

func toPackageRepository(pr *model.PackageRepository) *v1.PackageRepository {
	ret := &v1.PackageRepository{}
	ret.SetId(int32(pr.ID)) //nolint:gosec
	ret.SetName(pr.Name)
	ret.SetUrl(pr.URL)
	ret.SetEnabled(pr.Enabled)
	return ret
}

func toNetworks(ns []*model.Network) []*v1.Network {
	ret := make([]*v1.Network, len(ns))
	for i := range ns {
		ret[i] = toNetwork(ns[i])
	}
	return ret
}

func toNetwork(n *model.Network) *v1.Network {
	settings := make([]*v1.Network_Setting, len(n.Settings))
	for i := range n.Settings {
		s := n.Settings[i]
		setting := &v1.Network_Setting{}
		setting.SetVlan(toVLAN(s.VLAN))
		setting.SetSubnet(toSubnet(s.Subnet))
		setting.SetIpRange(toIPRange(s.IPRange))
		settings[i] = setting
	}
	ret := &v1.Network{}
	ret.SetFabric(toFabric(n.Fabric))
	ret.SetSettings(settings)
	return ret
}

func toFabric(f *model.Fabric) *v1.Fabric {
	ret := &v1.Fabric{}
	ret.SetId(int32(f.ID)) //nolint:gosec
	ret.SetName(f.Name)
	return ret
}

func toVLAN(v *model.VLAN) *v1.VLAN {
	ret := &v1.VLAN{}
	ret.SetId(int32(v.ID))   //nolint:gosec
	ret.SetVid(int32(v.VID)) //nolint:gosec
	ret.SetName(v.Name)
	ret.SetMtu(int32(v.MTU)) //nolint:gosec
	ret.SetDescription(v.Description)
	ret.SetDhcpOn(v.DHCPOn)
	return ret
}

func toSubnet(s *model.Subnet) *v1.Subnet {
	dnsServers := make([]string, len(s.DNSServers))
	for i, dns := range s.DNSServers {
		dnsServers[i] = dns.String()
	}
	ret := &v1.Subnet{}
	ret.SetId(int32(s.ID)) //nolint:gosec
	ret.SetName(s.Name)
	ret.SetCidr(s.CIDR)
	ret.SetGatewayIp(s.GatewayIP.String())
	ret.SetDnsServers(dnsServers)
	ret.SetDescription(s.Description)
	ret.SetManagedAllocation(s.Managed)
	ret.SetActiveDiscovery(s.ActiveDiscovery)
	ret.SetAllowProxyAccess(s.AllowProxy)
	ret.SetAllowDnsResolution(s.AllowDNS)
	return ret
}

func toIPRange(r *model.IPRange) *v1.IPRange {
	ret := &v1.IPRange{}
	ret.SetId(int32(r.ID)) //nolint:gosec
	ret.SetStartIp(r.StartIP.String())
	ret.SetEndIp(r.EndIP.String())
	ret.SetType(r.Type)
	ret.SetComment(r.Comment)
	return ret
}

func toMachines(ms []*model.Machine) []*v1.Machine {
	ret := make([]*v1.Machine, len(ms))
	for i := range ms {
		ret[i] = toMachine(ms[i])
	}
	return ret
}

//nolint:funlen
func toMachine(m *model.Machine) *v1.Machine {
	ipAddresses := make([]string, len(m.IPAddresses))
	for i, ip := range m.IPAddresses {
		ipAddresses[i] = ip.String()
	}
	numaNodes := make([]*v1.Machine_NUMANode, len(m.NUMANodeSet))
	for i := range m.NUMANodeSet {
		ns := &m.NUMANodeSet[i]
		numaNode := &v1.Machine_NUMANode{}
		numaNode.SetIndex(int32(ns.Index))      //nolint:gosec
		numaNode.SetCores(int32(len(ns.Cores))) //nolint:gosec
		numaNode.SetMemory(int64(ns.Memory))
		numaNodes[i] = numaNode
	}
	blockDevices := make([]*v1.Machine_BlockDevice, len(m.BlockDeviceSet))
	for i := range m.BlockDeviceSet {
		bds := &m.BlockDeviceSet[i]
		blockDevice := &v1.Machine_BlockDevice{}
		blockDevice.SetBootDisk(bds.ID == m.BootDisk.ID)
		blockDevice.SetName(bds.Name)
		blockDevice.SetSerial(bds.Serial)
		blockDevice.SetModel(bds.Model)
		blockDevice.SetFirmwareVersion(bds.FirmwareVersion)
		blockDevice.SetSize(bds.Size)
		blockDevice.SetType(bds.Type)
		blockDevice.SetTags(bds.Tags)
		blockDevice.SetUsedFor(bds.UsedFor)
		blockDevices[i] = blockDevice
	}
	networkInterfaces := make([]*v1.Machine_NetworkInterface, len(m.InterfaceSet))
	for i := range m.InterfaceSet {
		is := &m.InterfaceSet[i]
		subnetName := ""
		subnetID := 0
		ipAdress := ""
		for j := range is.Links {
			link := &is.Links[j]
			subnetName = link.Subnet.Name
			subnetID = link.Subnet.ID
			ipAdress = link.IPAddress
			break
		}
		networkInterface := &v1.Machine_NetworkInterface{}
		networkInterface.SetBootInterface(is.ID == m.BootInterface.ID)
		networkInterface.SetName(is.Name)
		networkInterface.SetMacAddress(is.MACAddress)
		networkInterface.SetLinkConnected(is.LinkConnected)
		networkInterface.SetLinkSpeed(int32(is.LinkSpeed))           //nolint:gosec
		networkInterface.SetInterfaceSpeed(int32(is.InterfaceSpeed)) //nolint:gosec
		networkInterface.SetType(is.Type)
		networkInterface.SetFabricName(is.VLAN.Fabric)        //nolint:gosec
		networkInterface.SetFabricId(int32(is.VLAN.FabricID)) //nolint:gosec
		networkInterface.SetVlanName(is.VLAN.Name)
		networkInterface.SetVlanId(int32(is.VLAN.ID)) //nolint:gosec
		networkInterface.SetSubnetName(subnetName)
		networkInterface.SetSubnetId(int32(subnetID)) //nolint:gosec
		networkInterface.SetIpAddress(ipAdress)
		networkInterface.SetDhcpOn(is.VLAN.DHCPOn)
		networkInterfaces[i] = networkInterface
	}
	ret := &v1.Machine{}
	ret.SetSystemId(m.SystemID)
	ret.SetHardwareUuid(m.HardwareUUID)
	ret.SetHostname(m.Hostname)
	ret.SetFqdn(m.FQDN)
	ret.SetTags(m.TagNames)
	ret.SetDescription(m.Description)
	ret.SetStatus(m.StatusName)
	ret.SetPowerState(m.PowerState)
	ret.SetPowerType(m.PowerType)
	ret.SetOsystem(m.OSystem)
	ret.SetDistroSeries(m.DistroSeries)
	ret.SetHweKernel(m.HWEKernel)
	ret.SetArchitecture(m.Architecture)
	ret.SetCpuSpeed(int32(m.CPUSpeed)) //nolint:gosec
	ret.SetCpuCount(int32(m.CPUCount)) //nolint:gosec
	ret.SetMemory(m.Memory)
	ret.SetStorage(m.Storage)
	ret.SetIpAddresses(ipAddresses)
	ret.SetWorkloadAnnotations(m.WorkloadAnnotations)
	ret.SetHardwareInformation(m.HardwareInfo)
	ret.SetBiosBootMethod(m.BiosBootMethod)
	ret.SetNumaNodes(numaNodes)
	ret.SetBlockDevices(blockDevices)
	ret.SetNetworkInterfaces(networkInterfaces)
	return ret
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
