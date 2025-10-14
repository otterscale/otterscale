package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/otterscale/otterscale/api/machine/v1"
	"github.com/otterscale/otterscale/api/machine/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type MachineService struct {
	pbconnect.UnimplementedMachineServiceHandler

	machine *core.MachineUseCase
}

func NewMachineService(machine *core.MachineUseCase) *MachineService {
	return &MachineService{
		machine: machine,
	}
}

var _ pbconnect.MachineServiceHandler = (*MachineService)(nil)

func (s *MachineService) ListMachines(ctx context.Context, req *pb.ListMachinesRequest) (*pb.ListMachinesResponse, error) {
	machines, err := s.machine.ListMachines(ctx, req.GetScopeUuid())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListMachinesResponse{}
	resp.SetMachines(toProtoMachines(machines))
	return resp, nil
}

func (s *MachineService) GetMachine(ctx context.Context, req *pb.GetMachineRequest) (*pb.Machine, error) {
	machine, err := s.machine.GetMachine(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	resp := toProtoMachine(machine)
	return resp, nil
}

func (s *MachineService) CreateMachine(ctx context.Context, req *pb.CreateMachineRequest) (*pb.Machine, error) {
	machine, err := s.machine.CreateMachine(ctx, req.GetId(), req.GetEnableSsh(), req.GetSkipBmcConfig(), req.GetSkipNetworking(), req.GetSkipStorage(), req.GetScopeUuid(), req.GetTags())
	if err != nil {
		return nil, err
	}
	resp := toProtoMachine(machine)
	return resp, nil
}

func (s *MachineService) DeleteMachine(ctx context.Context, req *pb.DeleteMachineRequest) (*emptypb.Empty, error) {
	if err := s.machine.DeleteMachine(ctx, req.GetId(), req.GetForce(), req.GetPurgeDisk()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *MachineService) PowerOffMachine(ctx context.Context, req *pb.PowerOffMachineRequest) (*pb.Machine, error) {
	machine, err := s.machine.PowerOffMachine(ctx, req.GetId(), req.GetComment())
	if err != nil {
		return nil, err
	}
	resp := toProtoMachine(machine)
	return resp, nil
}

func (s *MachineService) AddMachineTags(ctx context.Context, req *pb.AddMachineTagsRequest) (*emptypb.Empty, error) {
	if err := s.machine.AddMachineTags(ctx, req.GetId(), req.GetTags()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *MachineService) RemoveMachineTags(ctx context.Context, req *pb.RemoveMachineTagsRequest) (*emptypb.Empty, error) {
	if err := s.machine.RemoveMachineTags(ctx, req.GetId(), req.GetTags()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *MachineService) ListTags(ctx context.Context, _ *pb.ListTagsRequest) (*pb.ListTagsResponse, error) {
	tags, err := s.machine.ListTags(ctx)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListTagsResponse{}
	resp.SetTags(toProtoTags(tags))
	return resp, nil
}

func (s *MachineService) GetTag(ctx context.Context, req *pb.GetTagRequest) (*pb.Tag, error) {
	tag, err := s.machine.GetTag(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoTag(tag)
	return resp, nil
}

func (s *MachineService) CreateTag(ctx context.Context, req *pb.CreateTagRequest) (*pb.Tag, error) {
	tag, err := s.machine.CreateTag(ctx, req.GetName(), req.GetComment())
	if err != nil {
		return nil, err
	}
	resp := toProtoTag(tag)
	return resp, nil
}

func (s *MachineService) DeleteTag(ctx context.Context, req *pb.DeleteTagRequest) (*emptypb.Empty, error) {
	if err := s.machine.DeleteTag(ctx, req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func toProtoMachines(ms []core.Machine) []*pb.Machine {
	ret := []*pb.Machine{}
	for i := range ms {
		ret = append(ret, toProtoMachine(&ms[i]))
	}
	return ret
}

func toProtoMachine(m *core.Machine) *pb.Machine {
	ipAddresses := make([]string, len(m.IPAddresses))
	for i, ip := range m.IPAddresses {
		ipAddresses[i] = ip.String()
	}
	ret := &pb.Machine{}
	ret.SetId(m.SystemID)
	if !m.LastCommissioned.IsZero() {
		ret.SetLastCommissioned(timestamppb.New(m.LastCommissioned))
	}
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
	ret.SetGpuDevices(toProtoNodeDevices(m.GPUs))
	return ret
}

func toProtoNUMANodes(ns []core.NUMANode) []*pb.Machine_NUMANode {
	ret := []*pb.Machine_NUMANode{}
	for i := range ns {
		ret = append(ret, toProtoNUMANode(&ns[i]))
	}
	return ret
}

func toProtoNUMANode(n *core.NUMANode) *pb.Machine_NUMANode {
	ret := &pb.Machine_NUMANode{}
	ret.SetIndex(int64(n.Index))
	ret.SetCpuCores(int64(len(n.Cores)))
	ret.SetMemoryMb(int64(n.Memory))
	return ret
}

func toProtoBlockDevices(bds []core.BlockDevice, bootDiskID int) []*pb.Machine_BlockDevice {
	ret := []*pb.Machine_BlockDevice{}
	for i := range bds {
		ret = append(ret, toProtoBlockDevice(&bds[i], bootDiskID))
	}
	return ret
}

func toProtoBlockDevice(bd *core.BlockDevice, bootDiskID int) *pb.Machine_BlockDevice {
	ret := &pb.Machine_BlockDevice{}
	ret.SetBootDisk(bd.ID == bootDiskID)
	ret.SetName(bd.Name)
	ret.SetSerial(bd.Serial)
	ret.SetModel(bd.Model)
	ret.SetFirmwareVersion(bd.FirmwareVersion)
	ret.SetStorageMb(float64(bd.Size) / 1000 / 1000) //nolint:mnd // convert to MB
	ret.SetType(bd.Type)
	ret.SetTags(bd.Tags)
	ret.SetUsedFor(bd.UsedFor)
	return ret
}

func toProtoNetworkInterfaces(nis []core.NetworkInterface, bootInterfaceID int) []*pb.Machine_NetworkInterface {
	ret := []*pb.Machine_NetworkInterface{}
	for i := range nis {
		ret = append(ret, toProtoNetworkInterface(&nis[i], bootInterfaceID))
	}
	return ret
}

func toProtoNetworkInterface(ni *core.NetworkInterface, bootInterfaceID int) *pb.Machine_NetworkInterface {
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

func toProtoNodeDevices(ns []core.NodeDevice) []*pb.Machine_NodeDevice {
	ret := []*pb.Machine_NodeDevice{}
	for i := range ns {
		ret = append(ret, toProtoNodeDevice(&ns[i]))
	}
	return ret
}

func toProtoNodeDevice(n *core.NodeDevice) *pb.Machine_NodeDevice {
	ret := &pb.Machine_NodeDevice{}
	ret.SetVendorId(n.VendorID)
	ret.SetVendorName(n.VendorName)
	ret.SetProductId(n.ProductID)
	ret.SetProductName(n.ProductName)
	ret.SetBusName(n.BusName)
	ret.SetPciAddress(n.PCIAddress)
	return ret
}

func toProtoTags(ts []core.Tag) []*pb.Tag {
	ret := []*pb.Tag{}
	for i := range ts {
		ret = append(ret, toProtoTag(&ts[i]))
	}
	return ret
}

func toProtoTag(t *core.Tag) *pb.Tag {
	ret := &pb.Tag{}
	ret.SetName(t.Name)
	ret.SetComment(t.Comment)
	return ret
}
