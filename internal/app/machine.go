package app

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/openhdc/api/nexus/v1"
	"github.com/openhdc/openhdc/internal/domain/model"
)

func (a *NexusApp) ListMachines(ctx context.Context, req *connect.Request[pb.ListMachinesRequest]) (*connect.Response[pb.ListMachinesResponse], error) {
	ms, err := a.svc.ListMachines(ctx, req.Msg.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	res := &pb.ListMachinesResponse{}
	res.SetMachines(toProtoMachines(ms))
	return connect.NewResponse(res), nil
}

func (a *NexusApp) GetMachine(ctx context.Context, req *connect.Request[pb.GetMachineRequest]) (*connect.Response[pb.Machine], error) {
	m, err := a.svc.GetMachine(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	res := toProtoMachine(m)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) CreateMachine(ctx context.Context, req *connect.Request[pb.CreateMachineRequest]) (*connect.Response[pb.Machine], error) {
	m, err := a.svc.CreateMachine(ctx, req.Msg.GetId(), req.Msg.GetEnableSsh(), req.Msg.GetSkipBmcConfig(), req.Msg.GetSkipNetworking(), req.Msg.GetSkipStorage(), req.Msg.GetScopeUuid(), req.Msg.GetTags())
	if err != nil {
		return nil, err
	}
	res := toProtoMachine(m)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) DeleteMachine(ctx context.Context, req *connect.Request[pb.DeleteMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.DeleteMachine(ctx, req.Msg.GetId(), req.Msg.GetForce()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) PowerOnMachine(ctx context.Context, req *connect.Request[pb.PowerOnMachineRequest]) (*connect.Response[pb.Machine], error) {
	m, err := a.svc.PowerOnMachine(ctx, req.Msg.GetId(), req.Msg.GetComment())
	if err != nil {
		return nil, err
	}
	res := toProtoMachine(m)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) PowerOffMachine(ctx context.Context, req *connect.Request[pb.PowerOffMachineRequest]) (*connect.Response[pb.Machine], error) {
	m, err := a.svc.PowerOffMachine(ctx, req.Msg.GetId(), req.Msg.GetComment())
	if err != nil {
		return nil, err
	}
	res := toProtoMachine(m)
	return connect.NewResponse(res), nil
}

func (a *NexusApp) AddMachineTags(ctx context.Context, req *connect.Request[pb.AddMachineTagsRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.AddMachineTags(ctx, req.Msg.GetId(), req.Msg.GetTags()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func (a *NexusApp) RemoveMachineTags(ctx context.Context, req *connect.Request[pb.RemoveMachineTagsRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := a.svc.RemoveMachineTags(ctx, req.Msg.GetId(), req.Msg.GetTags()); err != nil {
		return nil, err
	}
	res := &emptypb.Empty{}
	return connect.NewResponse(res), nil
}

func toProtoMachines(ms []model.Machine) []*pb.Machine {
	ret := []*pb.Machine{}
	for i := range ms {
		ret = append(ret, toProtoMachine(&ms[i]))
	}
	return ret
}

func toProtoMachine(m *model.Machine) *pb.Machine {
	ipAddresses := make([]string, len(m.IPAddresses))
	for i, ip := range m.IPAddresses {
		ipAddresses[i] = ip.String()
	}
	ret := &pb.Machine{}
	ret.SetId(m.SystemID)
	ret.SetHardwareUuid(m.HardwareUUID)
	ret.SetHostname(m.Hostname)
	ret.SetFqdn(m.FQDN)
	ret.SetTags(m.TagNames)
	ret.SetDescription(m.Description)
	ret.SetStatus(m.StatusName)
	ret.SetStatusMessage(m.StatusMessage)
	ret.SetPowerState(m.PowerState)
	ret.SetPowerType(m.PowerType)
	ret.SetOsystem(m.OSystem)
	ret.SetDistroSeries(m.DistroSeries)
	ret.SetHweKernel(m.HWEKernel)
	ret.SetArchitecture(m.Architecture)
	ret.SetCpuSpeed(int64(m.CPUSpeed))
	ret.SetCpuCount(int64(m.CPUCount))
	ret.SetMemoryMb(m.Memory)
	ret.SetStorageMb(m.Storage)
	ret.SetIpAddresses(ipAddresses)
	ret.SetWorkloadAnnotations(m.WorkloadAnnotations)
	ret.SetHardwareInformation(m.HardwareInfo)
	ret.SetBiosBootMethod(m.BiosBootMethod)
	ret.SetNumaNodes(toProtoNUMANodes(m.NUMANodeSet))
	ret.SetBlockDevices(toProtoBlockDevices(m.BlockDeviceSet, m.BootDisk.ID))
	ret.SetNetworkInterfaces(toProtoNetworkInterfaces(m.InterfaceSet, m.BootInterface.ID))
	return ret
}

func toProtoNUMANodes(ns []model.NUMANode) []*pb.Machine_NUMANode {
	ret := []*pb.Machine_NUMANode{}
	for i := range ns {
		ret = append(ret, toProtoNUMANode(&ns[i]))
	}
	return ret
}

func toProtoNUMANode(n *model.NUMANode) *pb.Machine_NUMANode {
	ret := &pb.Machine_NUMANode{}
	ret.SetIndex(int64(n.Index))
	ret.SetCpuCores(int64(len(n.Cores)))
	ret.SetMemoryMb(int64(n.Memory))
	return ret
}

func toProtoBlockDevices(bds []model.BlockDevice, bootDiskID int) []*pb.Machine_BlockDevice {
	ret := []*pb.Machine_BlockDevice{}
	for i := range bds {
		ret = append(ret, toProtoBlockDevice(&bds[i], bootDiskID))
	}
	return ret
}

func toProtoBlockDevice(bd *model.BlockDevice, bootDiskID int) *pb.Machine_BlockDevice {
	ret := &pb.Machine_BlockDevice{}
	ret.SetBootDisk(bd.ID == bootDiskID)
	ret.SetName(bd.Name)
	ret.SetSerial(bd.Serial)
	ret.SetModel(bd.Model)
	ret.SetFirmwareVersion(bd.FirmwareVersion)
	ret.SetStorageMb(float64(bd.Size) / 1000 / 1000) //nolint:mnd
	ret.SetType(bd.Type)
	ret.SetTags(bd.Tags)
	ret.SetUsedFor(bd.UsedFor)
	return ret
}

func toProtoNetworkInterfaces(nis []model.NetworkInterface, bootInterfaceID int) []*pb.Machine_NetworkInterface {
	ret := []*pb.Machine_NetworkInterface{}
	for i := range nis {
		ret = append(ret, toProtoNetworkInterface(&nis[i], bootInterfaceID))
	}
	return ret
}

func toProtoNetworkInterface(ni *model.NetworkInterface, bootInterfaceID int) *pb.Machine_NetworkInterface {
	subnetName := ""
	subnetID := 0
	ipAdress := ""
	for i := range ni.Links {
		subnetName = ni.Links[i].Subnet.Name
		subnetID = ni.Links[i].Subnet.ID
		ipAdress = ni.Links[i].IPAddress
		break
	}
	ret := &pb.Machine_NetworkInterface{}
	ret.SetBootInterface(ni.ID == bootInterfaceID)
	ret.SetName(ni.Name)
	ret.SetMacAddress(ni.MACAddress)
	ret.SetLinkConnected(ni.LinkConnected)
	ret.SetLinkSpeed(int64(ni.LinkSpeed))
	ret.SetInterfaceSpeed(int64(ni.InterfaceSpeed))
	ret.SetType(ni.Type)
	ret.SetFabricName(ni.VLAN.Fabric)
	ret.SetFabricId(int64(ni.VLAN.FabricID))
	ret.SetVlanName(ni.VLAN.Name)
	ret.SetVlanId(int64(ni.VLAN.ID))
	ret.SetSubnetName(subnetName)
	ret.SetSubnetId(int64(subnetID))
	ret.SetIpAddress(ipAdress)
	ret.SetDhcpOn(ni.VLAN.DHCPOn)
	return ret
}

func toModelFactors(fs []*pb.Machine_Factor) []model.MachineFactor {
	ret := []model.MachineFactor{}
	for i := range fs {
		ret = append(ret, *toModelFactor(fs[i]))
	}
	return ret
}

func toModelFactor(f *pb.Machine_Factor) *model.MachineFactor {
	ret := &model.MachineFactor{}
	p := f.GetPlacement()
	if p != nil {
		ret.MachinePlacement = toModelPlacement(p)
	}
	c := f.GetConstraint()
	if c != nil {
		ret.MachineConstraint = toModelConstraint(c)
	}
	return ret
}

func toModelPlacements(ps []*pb.Machine_Placement) []model.MachinePlacement {
	ret := []model.MachinePlacement{}
	for i := range ps {
		ret = append(ret, *toModelPlacement(ps[i]))
	}
	return ret
}

func toModelPlacement(p *pb.Machine_Placement) *model.MachinePlacement {
	return &model.MachinePlacement{
		LXD:       p.GetLxd(),
		KVM:       p.GetKvm(),
		Machine:   p.GetMachine(),
		MachineID: p.GetMachineId(),
	}
}

func toModelConstraint(c *pb.Machine_Constraint) *model.MachineConstraint {
	return &model.MachineConstraint{
		Architecture: c.GetArchitecture(),
		CPUCores:     c.GetCpuCores(),
		MemoryMB:     c.GetMemoryMb(),
		Tags:         c.GetTags(),
	}
}
