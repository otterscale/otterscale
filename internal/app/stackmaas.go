package app

import (
	"context"

	"connectrpc.com/connect"
	"github.com/canonical/gomaasclient/entity"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/openhdc/openhdc/api/stack/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *StackApp) ListNTPServers(ctx context.Context, req *connect.Request[v1.ListNTPServersRequest]) (*connect.Response[v1.ListNTPServersResponse], error) {
	ntpServers, err := a.svc.ListNTPServers(ctx)
	if err != nil {
		return nil, err
	}
	res := &v1.ListNTPServersResponse{}
	res.SetNtpServers(ntpServers)
	return connect.NewResponse(res), nil
}

// UpdateNTPServers updates NTP server configuration
func (a *StackApp) UpdateNTPServers(ctx context.Context, req *connect.Request[v1.UpdateNTPServersRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.UpdateNTPServers(ctx, req.Msg.GetNtpServers()); err != nil {
		return nil, err
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

// ListPackageRepositories retrieves package repositories
func (a *StackApp) ListPackageRepositories(ctx context.Context, req *connect.Request[v1.ListPackageRepositoriesRequest]) (*connect.Response[v1.ListPackageRepositoriesResponse], error) {
	prs, err := a.svc.ListPackageRepositories(ctx)
	if err != nil {
		return nil, err
	}
	res := &v1.ListPackageRepositoriesResponse{}
	res.SetPackageRepositories(toPackageRepositories(prs))
	return connect.NewResponse(res), nil
}

// UpdatePackageRepositoryURL updates a package repository URL
func (a *StackApp) UpdatePackageRepositoryURL(ctx context.Context, req *connect.Request[v1.UpdatePackageRepositoryURLRequest]) (*connect.Response[v1.PackageRepository], error) {
	pr, err := a.svc.UpdatePackageRepositoryURL(ctx, int(req.Msg.GetId()), req.Msg.GetUrl())
	if err != nil {
		return nil, err
	}
	if !req.Msg.GetSkipJuju() {
		ms, err := a.svc.ListModels(ctx)
		if err != nil {
			return nil, err
		}
		for _, m := range ms {
			if err := a.svc.SetModelConfigAPTMirror(ctx, m.UUID, req.Msg.GetUrl()); err != nil {
				return nil, err
			}
		}
	}
	return connect.NewResponse(toPackageRepository(pr)), nil
}

// Network management functions
func (a *StackApp) ListNetworks(ctx context.Context, req *connect.Request[v1.ListNetworksRequest]) (*connect.Response[v1.ListNetworksResponse], error) {
	ns, err := a.svc.ListNetworks(ctx)
	if err != nil {
		return nil, err
	}
	res := &v1.ListNetworksResponse{}
	res.SetNetworks(toNetworks(ns))
	return connect.NewResponse(res), nil
}

func (a *StackApp) CreateNetwork(ctx context.Context, req *connect.Request[v1.CreateNetworkRequest]) (*connect.Response[v1.Network], error) {
	fabricParams := &entity.FabricParams{}
	vlanParams := &entity.VLANParams{
		DHCPOn: req.Msg.GetDhcpOn(),
	}
	subnetParams := &entity.SubnetParams{
		CIDR:       req.Msg.GetCidr(),
		GatewayIP:  req.Msg.GetGatewayIp(),
		DNSServers: req.Msg.GetDnsServers(),
	}
	ipRangeParams := &entity.IPRangeParams{
		StartIP: req.Msg.GetStartIp(),
		EndIP:   req.Msg.GetEndIp(),
	}
	n, err := a.svc.CreateNetwork(ctx, fabricParams, vlanParams, subnetParams, ipRangeParams)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(toNetwork(n)), nil
}

func (a *StackApp) DeleteNetwork(ctx context.Context, req *connect.Request[v1.DeleteNetworkRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.DeleteNetwork(ctx, int(req.Msg.GetFabricId())); err != nil {
		return nil, err
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

// Network configuration functions
func (a *StackApp) UpdateFabric(ctx context.Context, req *connect.Request[v1.UpdateFabricRequest]) (*connect.Response[v1.Fabric], error) {
	params := &entity.FabricParams{
		Name: req.Msg.GetName(),
	}
	f, err := a.svc.UpdateFabric(ctx, int(req.Msg.GetId()), params)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(toFabric(f)), nil
}

func (a *StackApp) UpdateVLAN(ctx context.Context, req *connect.Request[v1.UpdateVLANRequest]) (*connect.Response[v1.VLAN], error) {
	params := &entity.VLANParams{
		Name:        req.Msg.GetName(),
		MTU:         int(req.Msg.GetMtu()),
		Description: req.Msg.GetDescription(),
		DHCPOn:      req.Msg.GetDhcpOn(),
	}
	v, err := a.svc.UpdateVLAN(ctx, int(req.Msg.GetFabricId()), int(req.Msg.GetVid()), params)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(toVLAN(v)), nil
}

func (a *StackApp) UpdateSubnet(ctx context.Context, req *connect.Request[v1.UpdateSubnetRequest]) (*connect.Response[v1.Subnet], error) {
	params := &entity.SubnetParams{
		Name:        req.Msg.GetName(),
		CIDR:        req.Msg.GetCidr(),
		GatewayIP:   req.Msg.GetGatewayIp(),
		DNSServers:  req.Msg.GetDnsServers(),
		Description: req.Msg.GetDescription(),
		AllowDNS:    req.Msg.GetAllowDnsResolution(),
	}
	s, err := a.svc.UpdateSubnet(ctx, int(req.Msg.GetId()), params)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(toSubnet(s)), nil
}

func (a *StackApp) UpdateIPRange(ctx context.Context, req *connect.Request[v1.UpdateIPRangeRequest]) (*connect.Response[v1.IPRange], error) {
	params := &entity.IPRangeParams{
		StartIP: req.Msg.GetStartIp(),
		EndIP:   req.Msg.GetEndIp(),
		Comment: req.Msg.GetComment(),
	}
	r, err := a.svc.UpdateIPRange(ctx, int(req.Msg.GetId()), params)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(toIPRange(r)), nil
}

// ListMachines retrieves machines
func (a *StackApp) ListMachines(ctx context.Context, req *connect.Request[v1.ListMachinesRequest]) (*connect.Response[v1.ListMachinesResponse], error) {
	ms, err := a.svc.ListMachines(ctx)
	if err != nil {
		return nil, err
	}
	res := &v1.ListMachinesResponse{}
	res.SetMachines(toMachines(ms))
	return connect.NewResponse(res), nil
}

func (a *StackApp) GetMachine(ctx context.Context, req *connect.Request[v1.GetMachineRequest]) (*connect.Response[v1.Machine], error) {
	m, err := a.svc.GetMachine(ctx, req.Msg.GetSystemId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(toMachine(m)), nil
}

// Boot resource management
func (a *StackApp) ImportBootResources(ctx context.Context, req *connect.Request[v1.ImportBootResourcesRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.ImportBootResources(ctx); err != nil {
		return nil, err
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

// Machine power management functions
func (a *StackApp) PowerOnMachine(ctx context.Context, req *connect.Request[v1.PowerOnMachineRequest]) (*connect.Response[v1.Machine], error) {
	params := &entity.MachinePowerOnParams{}
	m, err := a.svc.PowerOnMachine(ctx, req.Msg.GetSystemId(), params)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(toMachine(m)), nil
}

func (a *StackApp) PowerOffMachine(ctx context.Context, req *connect.Request[v1.PowerOffMachineRequest]) (*connect.Response[v1.Machine], error) {
	params := &entity.MachinePowerOffParams{}
	m, err := a.svc.PowerOffMachine(ctx, req.Msg.GetSystemId(), params)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(toMachine(m)), nil
}

// Machine provisioning functions
func (a *StackApp) CommissionMachine(ctx context.Context, req *connect.Request[v1.CommissionMachineRequest]) (*connect.Response[v1.Machine], error) {
	params := &entity.MachineCommissionParams{
		EnableSSH:      boolToInt(req.Msg.GetEnableSsh()),
		SkipBMCConfig:  boolToInt(req.Msg.GetSkipBmcConfig()),
		SkipNetworking: boolToInt(req.Msg.GetSkipNetworking()),
		SkipStorage:    boolToInt(req.Msg.GetSkipStorage()),
	}
	m, err := a.svc.CommissionMachine(ctx, req.Msg.GetSystemId(), params)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(toMachine(m)), nil
}

func toPackageRepositories(prs []*entity.PackageRepository) []*v1.PackageRepository {
	ret := make([]*v1.PackageRepository, len(prs))
	for i := range prs {
		ret[i] = toPackageRepository(prs[i])
	}
	return ret
}

func toPackageRepository(pr *entity.PackageRepository) *v1.PackageRepository {
	ret := &v1.PackageRepository{}
	ret.SetId(int64(pr.ID))
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
		if s.VLAN != nil {
			setting.SetVlan(toVLAN(s.VLAN))
		}
		if s.Subnet != nil {
			setting.SetSubnet(toSubnet(s.Subnet))
		}
		settings[i] = setting
	}
	ret := &v1.Network{}
	ret.SetFabric(toFabric(n.Fabric))
	ret.SetSettings(settings)
	return ret
}

func toFabric(f *entity.Fabric) *v1.Fabric {
	ret := &v1.Fabric{}
	ret.SetId(int64(f.ID))
	ret.SetName(f.Name)
	return ret
}

func toVLAN(v *entity.VLAN) *v1.VLAN {
	ret := &v1.VLAN{}
	ret.SetId(int64(v.ID))
	ret.SetVid(int64(v.VID))
	ret.SetName(v.Name)
	ret.SetMtu(int64(v.MTU))
	ret.SetDescription(v.Description)
	ret.SetDhcpOn(v.DHCPOn)
	return ret
}

func toSubnet(s *model.NetworkSubnet) *v1.Subnet {
	dnsServers := make([]string, len(s.DNSServers))
	for i, dns := range s.DNSServers {
		dnsServers[i] = dns.String()
	}
	ipAddresses := []*v1.Subnet_IPAddress{}
	for i := range s.IPAddresses {
		ipAddress := &v1.Subnet_IPAddress{}
		ipAddress.SetType(model.AllocType(s.IPAddresses[i].AllocType).String())
		ipAddress.SetIp(s.IPAddresses[i].IP.String())
		ipAddress.SetUser(s.IPAddresses[i].User)
		ipAddress.SetSystemId(s.IPAddresses[i].NodeSummary.SystemID)
		ipAddress.SetNodeType(model.NodeType(s.IPAddresses[i].NodeSummary.NodeType).String())
		ipAddress.SetHostname(s.IPAddresses[i].NodeSummary.Hostname)
		ipAddresses = append(ipAddresses, ipAddress)
	}
	reservedIPRanges := []*v1.Subnet_ReservedIPRange{}
	for i := range s.ReservedIPRanges {
		reservedIPRange := &v1.Subnet_ReservedIPRange{}
		reservedIPRange.SetPurposes(s.ReservedIPRanges[i].Purpose)
		reservedIPRange.SetStart(s.ReservedIPRanges[i].IPRange.Start.String())
		reservedIPRange.SetEnd(s.ReservedIPRanges[i].IPRange.End.String())
		reservedIPRange.SetCount(int64(s.ReservedIPRanges[i].IPRange.NumAddresses))
		reservedIPRanges = append(reservedIPRanges, reservedIPRange)
	}
	statistics := &v1.Subnet_Statistics{}
	statistics.SetAvailable(int64(s.Statistics.NumAvailable))
	statistics.SetTotal(int64(s.Statistics.TotalAddresses))
	statistics.SetUsagePercent(s.Statistics.UsageString)
	statistics.SetAvailablePercent(s.Statistics.AvailableString)
	ret := &v1.Subnet{}
	ret.SetId(int64(s.ID))
	ret.SetName(s.Name)
	ret.SetCidr(s.CIDR)
	ret.SetGatewayIp(s.GatewayIP.String())
	ret.SetDnsServers(dnsServers)
	ret.SetDescription(s.Description)
	ret.SetManagedAllocation(s.Managed)
	ret.SetActiveDiscovery(s.ActiveDiscovery)
	ret.SetAllowProxyAccess(s.AllowProxy)
	ret.SetAllowDnsResolution(s.AllowDNS)
	ret.SetIpAddresses(ipAddresses)
	ret.SetReservedIpRanges(reservedIPRanges)
	ret.SetStatistics(statistics)
	return ret
}

func toIPRange(r *entity.IPRange) *v1.IPRange {
	ret := &v1.IPRange{}
	ret.SetId(int64(r.ID))
	ret.SetStartIp(r.StartIP.String())
	ret.SetEndIp(r.EndIP.String())
	ret.SetType(r.Type)
	ret.SetComment(r.Comment)
	return ret
}

func toMachines(ms []*entity.Machine) []*v1.Machine {
	ret := make([]*v1.Machine, len(ms))
	for i := range ms {
		ret[i] = toMachine(ms[i])
	}
	return ret
}

//nolint:funlen
func toMachine(m *entity.Machine) *v1.Machine {
	ipAddresses := make([]string, len(m.IPAddresses))
	for i, ip := range m.IPAddresses {
		ipAddresses[i] = ip.String()
	}
	numaNodes := make([]*v1.Machine_NUMANode, len(m.NUMANodeSet))
	for i := range m.NUMANodeSet {
		ns := &m.NUMANodeSet[i]
		numaNode := &v1.Machine_NUMANode{}
		numaNode.SetIndex(int64(ns.Index))
		numaNode.SetCores(int64(len(ns.Cores)))
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
		networkInterface.SetLinkSpeed(int64(is.LinkSpeed))
		networkInterface.SetInterfaceSpeed(int64(is.InterfaceSpeed))
		networkInterface.SetType(is.Type)
		networkInterface.SetFabricName(is.VLAN.Fabric)
		networkInterface.SetFabricId(int64(is.VLAN.FabricID))
		networkInterface.SetVlanName(is.VLAN.Name)
		networkInterface.SetVlanId(int64(is.VLAN.ID))
		networkInterface.SetSubnetName(subnetName)
		networkInterface.SetSubnetId(int64(subnetID))
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
	ret.SetCpuSpeed(int64(m.CPUSpeed))
	ret.SetCpuCount(int64(m.CPUCount))
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
