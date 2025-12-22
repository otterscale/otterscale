package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"k8s.io/apimachinery/pkg/util/intstr"

	apppb "github.com/otterscale/otterscale/api/application/v1"
	pb "github.com/otterscale/otterscale/api/instance/v1"
	"github.com/otterscale/otterscale/api/instance/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/application/workload"
	"github.com/otterscale/otterscale/internal/core/instance/cdi"
	"github.com/otterscale/otterscale/internal/core/instance/vm"
	"github.com/otterscale/otterscale/internal/core/instance/vmi"
	"github.com/otterscale/otterscale/internal/core/instance/vnc"
)

type InstanceService struct {
	pbconnect.UnimplementedInstanceServiceHandler

	dataVolume             *cdi.UseCase
	virtualMachine         *vm.UseCase
	virtualMachineInstance *vmi.UseCase
	vnc                    *vnc.UseCase
}

func NewInstanceService(dataVolume *cdi.UseCase, virtualMachine *vm.UseCase, virtualMachineInstance *vmi.UseCase, vnc *vnc.UseCase) *InstanceService {
	return &InstanceService{
		dataVolume:             dataVolume,
		virtualMachine:         virtualMachine,
		virtualMachineInstance: virtualMachineInstance,
		vnc:                    vnc,
	}
}

var _ pbconnect.InstanceServiceHandler = (*InstanceService)(nil)

func (s *InstanceService) ListVirtualMachines(ctx context.Context, req *pb.ListVirtualMachinesRequest) (*pb.ListVirtualMachinesResponse, error) {
	vms, err := s.virtualMachine.ListVirtualMachines(ctx, req.GetScope(), req.GetNamespace())
	if err != nil {
		return nil, err
	}

	its, err := s.virtualMachineInstance.ListInstanceTypes(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListVirtualMachinesResponse{}
	resp.SetVirtualMachines(toProtoVirtualMachines(vms, its))
	return resp, nil
}

func (s *InstanceService) GetVirtualMachine(ctx context.Context, req *pb.GetVirtualMachineRequest) (*pb.VirtualMachine, error) {
	vm, err := s.virtualMachine.GetVirtualMachine(ctx, req.GetScope(), req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}

	its, err := s.virtualMachineInstance.ListInstanceTypes(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := toProtoVirtualMachine(vm, its)
	return resp, nil
}

func (s *InstanceService) CreateVirtualMachine(ctx context.Context, req *pb.CreateVirtualMachineRequest) (*pb.VirtualMachine, error) {
	vm, err := s.virtualMachine.CreateVirtualMachine(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetInstanceTypeName(), req.GetBootDataVolumeName(), req.GetStartupScript())
	if err != nil {
		return nil, err
	}

	its, err := s.virtualMachineInstance.ListInstanceTypes(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := toProtoVirtualMachine(vm, its)
	return resp, nil
}

func (s *InstanceService) DeleteVirtualMachine(ctx context.Context, req *pb.DeleteVirtualMachineRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachine.DeleteVirtualMachine(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) AttachVirtualMachineDisk(ctx context.Context, req *pb.AttachVirtualMachineDiskRequest) (*pb.VirtualMachine_Disk, error) {
	disk, volume, err := s.virtualMachine.AttachVirtualMachineDisk(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetDataVolumeName())
	if err != nil {
		return nil, err
	}

	resp := toProtoVirtualMachineDisk(disk, []vm.VirtualMachineVolume{*volume})
	return resp, nil
}

func (s *InstanceService) DetachVirtualMachineDisk(ctx context.Context, req *pb.DetachVirtualMachineDiskRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachine.DetachVirtualMachineDisk(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetDataVolumeName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) CreateVirtualMachineClone(ctx context.Context, req *pb.CreateVirtualMachineCloneRequest) (*pb.VirtualMachine_Clone, error) {
	clone, err := s.virtualMachine.CreateVirtualMachineClone(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetSourceVirtualMachineName(), req.GetTargetVirtualMachineName())
	if err != nil {
		return nil, err
	}

	resp := toProtoVirtualMachineClone(clone)
	return resp, nil
}

func (s *InstanceService) DeleteVirtualMachineClone(ctx context.Context, req *pb.DeleteVirtualMachineCloneRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachine.DeleteVirtualMachineClone(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) CreateVirtualMachineSnapshot(ctx context.Context, req *pb.CreateVirtualMachineSnapshotRequest) (*pb.VirtualMachine_Snapshot, error) {
	snapshot, err := s.virtualMachine.CreateVirtualMachineSnapshot(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetVirtualMachineName())
	if err != nil {
		return nil, err
	}

	resp := toProtoVirtualMachineSnapshot(snapshot)
	return resp, nil
}

func (s *InstanceService) DeleteVirtualMachineSnapshot(ctx context.Context, req *pb.DeleteVirtualMachineSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachine.DeleteVirtualMachineSnapshot(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) CreateVirtualMachineRestore(ctx context.Context, req *pb.CreateVirtualMachineRestoreRequest) (*pb.VirtualMachine_Restore, error) {
	restore, err := s.virtualMachine.CreateVirtualMachineRestore(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetVirtualMachineName(), req.GetSnapshotName())
	if err != nil {
		return nil, err
	}

	resp := toProtoVirtualMachineRestore(restore)
	return resp, nil
}

func (s *InstanceService) DeleteVirtualMachineRestore(ctx context.Context, req *pb.DeleteVirtualMachineRestoreRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachine.DeleteVirtualMachineRestore(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) StartVirtualMachine(ctx context.Context, req *pb.StartVirtualMachineRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachine.StartVirtualMachine(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) StopVirtualMachine(ctx context.Context, req *pb.StopVirtualMachineRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachine.StopVirtualMachine(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) RestartVirtualMachine(ctx context.Context, req *pb.RestartVirtualMachineRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachine.RestartVirtualMachine(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) PauseInstance(ctx context.Context, req *pb.PauseInstanceRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachineInstance.PauseInstance(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) ResumeInstance(ctx context.Context, req *pb.ResumeInstanceRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachineInstance.ResumeInstance(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) MigrateInstance(ctx context.Context, req *pb.MigrateInstanceRequest) (*emptypb.Empty, error) {
	if _, err := s.virtualMachineInstance.MigrateInstance(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetHostname()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) VNCInstance(_ context.Context, req *pb.VNCInstanceRequest) (*pb.VNCInstanceResponse, error) {
	sessionID, err := s.vnc.CreateVNCSession(req.GetScope(), req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}

	resp := &pb.VNCInstanceResponse{}
	resp.SetSessionId(sessionID)
	return resp, nil
}

func (s *InstanceService) ListDataVolumes(ctx context.Context, req *pb.ListDataVolumesRequest) (*pb.ListDataVolumesResponse, error) {
	bootImage := filterToBootImagePointer(req.GetFilter())

	its, err := s.dataVolume.ListDataVolumes(ctx, req.GetScope(), req.GetNamespace(), bootImage)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListDataVolumesResponse{}
	resp.SetDataVolumes(toProtoDataVolumes(its))
	return resp, nil
}

func (s *InstanceService) GetDataVolume(ctx context.Context, req *pb.GetDataVolumeRequest) (*pb.DataVolume, error) {
	it, err := s.dataVolume.GetDataVolume(ctx, req.GetScope(), req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}

	resp := toProtoDataVolume(it)
	return resp, nil
}

func (s *InstanceService) CreateDataVolume(ctx context.Context, req *pb.CreateDataVolumeRequest) (*pb.DataVolume, error) {
	src := req.GetSource()
	if src == nil {
		return nil, fmt.Errorf("data volume source is required")
	}

	it, err := s.dataVolume.CreateDataVolume(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), cdi.DataVolumeSourceType(src.GetType()), src.GetData(), req.GetSizeBytes(), req.GetBootImage())
	if err != nil {
		return nil, err
	}

	resp := toProtoDataVolume(it)
	return resp, nil
}

func (s *InstanceService) DeleteDataVolume(ctx context.Context, req *pb.DeleteDataVolumeRequest) (*emptypb.Empty, error) {
	if err := s.dataVolume.DeleteDataVolume(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) ExtendDataVolume(ctx context.Context, req *pb.ExtendDataVolumeRequest) (*emptypb.Empty, error) {
	if err := s.dataVolume.ExtendDataVolume(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetSizeBytes()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) ListInstanceTypes(ctx context.Context, req *pb.ListInstanceTypesRequest) (*pb.ListInstanceTypesResponse, error) {
	its, err := s.virtualMachineInstance.ListInstanceTypes(ctx, req.GetScope())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListInstanceTypesResponse{}
	resp.SetInstanceTypes(toProtoInstanceTypes(its))
	return resp, nil
}

func (s *InstanceService) GetInstanceType(ctx context.Context, req *pb.GetInstanceTypeRequest) (*pb.InstanceType, error) {
	it, err := s.virtualMachineInstance.GetInstanceType(ctx, req.GetScope(), req.GetName())
	if err != nil {
		return nil, err
	}

	resp := toProtoInstanceType(it)
	return resp, nil
}

func (s *InstanceService) CreateInstanceType(ctx context.Context, req *pb.CreateInstanceTypeRequest) (*pb.InstanceType, error) {
	it, err := s.virtualMachineInstance.CreateInstanceType(ctx, req.GetScope(), req.GetName(), req.GetCpuCores(), req.GetMemoryBytes())
	if err != nil {
		return nil, err
	}

	resp := toProtoInstanceType(it)
	return resp, nil
}

func (s *InstanceService) DeleteInstanceType(ctx context.Context, req *pb.DeleteInstanceTypeRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachineInstance.DeleteInstanceType(ctx, req.GetScope(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) CreateVirtualMachineService(ctx context.Context, req *pb.CreateVirtualMachineServiceRequest) (*apppb.Application_Service, error) {
	svc, err := s.virtualMachine.CreateVirtualMachineService(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), req.GetVirtualMachineName(), toPorts(req.GetPorts()))
	if err != nil {
		return nil, err
	}

	resp := toProtoService(svc)
	return resp, nil
}

func (s *InstanceService) UpdateVirtualMachineService(ctx context.Context, req *pb.UpdateVirtualMachineServiceRequest) (*apppb.Application_Service, error) {
	svc, err := s.virtualMachine.UpdateVirtualMachineService(ctx, req.GetScope(), req.GetNamespace(), req.GetName(), toPorts(req.GetPorts()))
	if err != nil {
		return nil, err
	}

	resp := toProtoService(svc)
	return resp, nil
}

func (s *InstanceService) DeleteVirtualMachineService(ctx context.Context, req *pb.DeleteVirtualMachineServiceRequest) (*emptypb.Empty, error) {
	if err := s.virtualMachine.DeleteVirtualMachineService(ctx, req.GetScope(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *InstanceService) VNCPathPrefix() string {
	return s.vnc.VNCPathPrefix()
}

func (s *InstanceService) VNCHandler() http.HandlerFunc {
	return s.vnc.VNCHandler
}

func filterToBootImagePointer(filter pb.DataVolumeFilter) *bool {
	switch filter {
	case pb.DataVolumeFilter_BOOTABLE:
		val := true
		return &val
	case pb.DataVolumeFilter_NON_BOOTABLE:
		val := false
		return &val
	case pb.DataVolumeFilter_ALL:
		return nil
	}
}

func toPorts(ps []*apppb.Application_Service_Port) []service.Port {
	ret := []service.Port{}

	for i := range ps {
		ret = append(ret, toPort(ps[i]))
	}

	return ret
}

func toPort(p *apppb.Application_Service_Port) service.Port {
	return service.Port{
		Name:       p.GetName(),
		Port:       p.GetPort(),
		NodePort:   p.GetNodePort(),
		Protocol:   service.Protocol(strings.ToUpper(p.GetProtocol())),
		TargetPort: intstr.Parse(p.GetTargetPort()),
	}
}

func toProtoVirtualMachineDiskVolumeSource(v *vm.VirtualMachineVolumeSource) *pb.VirtualMachine_Disk_Volume_Source {
	ret := &pb.VirtualMachine_Disk_Volume_Source{}

	if v.DataVolume != nil {
		ret.SetType(pb.VirtualMachine_Disk_Volume_Source_TYPE_DATA_VOLUME)
		ret.SetData(v.DataVolume.Name)
	} else if v.CloudInitNoCloud != nil {
		ret.SetType(pb.VirtualMachine_Disk_Volume_Source_TYPE_CLOUD_INIT_NO_CLOUD)
		ret.SetData(v.CloudInitNoCloud.UserData)
	}

	return ret
}

func toProtoVirtualMachineDiskVolume(v *vm.VirtualMachineVolume) *pb.VirtualMachine_Disk_Volume {
	ret := &pb.VirtualMachine_Disk_Volume{}
	ret.SetName(v.Name)
	ret.SetSource(toProtoVirtualMachineDiskVolumeSource(&v.VolumeSource))
	return ret
}

func toProtoVirtualMachineDisks(ds []vm.VirtualMachineDisk, vs []vm.VirtualMachineVolume) []*pb.VirtualMachine_Disk {
	ret := []*pb.VirtualMachine_Disk{}

	for i := range ds {
		ret = append(ret, toProtoVirtualMachineDisk(&ds[i], vs))
	}

	return ret
}

func toProtoVirtualMachineDisk(d *vm.VirtualMachineDisk, vs []vm.VirtualMachineVolume) *pb.VirtualMachine_Disk {
	ret := &pb.VirtualMachine_Disk{}
	ret.SetName(d.Name)

	disk := d.Disk
	if disk != nil {
		ret.SetBus(toProtoVirtualMachineDiskBus(disk.Bus))
	}

	bootOrder := d.BootOrder
	if bootOrder != nil {
		ret.SetBootOrder(uint32(*bootOrder)) //nolint:gosec // ignore
	}

	for i := range vs {
		if vs[i].Name == d.Name {
			ret.SetVolume(toProtoVirtualMachineDiskVolume(&vs[i]))
		}
	}

	return ret
}

func toProtoVirtualMachineDiskBus(bus vm.VirtualMachineDiskBus) pb.VirtualMachine_Disk_Bus {
	switch bus {
	case vm.VirtualMachineDiskBusVirtio:
		return pb.VirtualMachine_Disk_BUS_VIRTIO

	case vm.VirtualMachineDiskBusSATA:
		return pb.VirtualMachine_Disk_BUS_SATA

	case vm.VirtualMachineDiskBusSCSI:
		return pb.VirtualMachine_Disk_BUS_SCSI

	case vm.VirtualMachineDiskBusUSB:
		return pb.VirtualMachine_Disk_BUS_USB

	default:
		return pb.VirtualMachine_Disk_BUS_VIRTIO
	}
}

func toProtoVirtualMachines(vmds []vm.VirtualMachineData, its []vmi.VirtualMachineInstanceTypeData) []*pb.VirtualMachine {
	ret := []*pb.VirtualMachine{}

	for i := range vmds {
		ret = append(ret, toProtoVirtualMachine(&vmds[i], its))
	}

	return ret
}

func toProtoVirtualMachine(vmd *vm.VirtualMachineData, its []vmi.VirtualMachineInstanceTypeData) *pb.VirtualMachine {
	ret := &pb.VirtualMachine{}
	ret.SetName(vmd.Name)
	ret.SetNamespace(vmd.Namespace)

	instanceType := vmd.Spec.Instancetype

	if instanceType != nil {
		for _, it := range its {
			if (it.ClusterWide && it.Type.Name == instanceType.Name) || (!it.ClusterWide && it.Type.Namespace == vmd.Namespace && it.Type.Name == instanceType.Name) {
				ret.SetInstanceType(toProtoInstanceType(&it))
				break
			}
		}
	}

	ret.SetStatus(string(vmd.Status.PrintableStatus))
	ret.SetReady(vmd.Status.Ready)

	instance := vmd.Instance

	if instance != nil {
		ret.SetInstancePhase(string(instance.Status.Phase))
	}

	machine := vmd.Machine

	if machine != nil {
		ret.SetMachineId(machine.SystemID)
		ret.SetHostname(machine.Hostname)

		ipAddresses := make([]string, len(machine.IPAddresses))

		for i, ip := range machine.IPAddresses {
			ipAddresses[i] = ip.String()
		}

		ret.SetIpAddresses(ipAddresses)
	}

	ret.SetCreatedAt(timestamppb.New(vmd.CreationTimestamp.Time))
	ret.SetServices(toProtoServices(vmd.Services))
	ret.SetDisks(toProtoVirtualMachineDisks(vmd.Spec.Template.Spec.Domain.Devices.Disks, vmd.Spec.Template.Spec.Volumes))
	ret.SetClones(toProtoVirtualMachineClones(vmd.Clones))
	ret.SetSnapshots(toProtoVirtualMachineSnapshots(vmd.Snapshots))
	ret.SetRestores(toProtoVirtualMachineRestores(vmd.Restores))

	return ret
}

func toProtoApplicationConditionFromClone(c *vm.VirtualMachineCloneCondition) *apppb.Application_Condition {
	ret := &apppb.Application_Condition{}
	ret.SetType(string(c.Type))
	ret.SetStatus(string(c.Status))
	ret.SetReason((c.Reason))
	ret.SetMessage((c.Message))

	if !c.LastProbeTime.IsZero() {
		ret.SetTransitionedAt(timestamppb.New(c.LastProbeTime.Time))
	}

	if !c.LastTransitionTime.IsZero() {
		ret.SetTransitionedAt(timestamppb.New(c.LastTransitionTime.Time))
	}

	return ret
}

func toProtoApplicationConditionFromSnapshot(c *vm.VirtualMachineSnapshotCondition) *apppb.Application_Condition {
	ret := &apppb.Application_Condition{}
	ret.SetType(string(c.Type))
	ret.SetStatus(string(c.Status))
	ret.SetReason((c.Reason))
	ret.SetMessage((c.Message))

	if !c.LastProbeTime.IsZero() {
		ret.SetTransitionedAt(timestamppb.New(c.LastProbeTime.Time))
	}

	if !c.LastTransitionTime.IsZero() {
		ret.SetTransitionedAt(timestamppb.New(c.LastTransitionTime.Time))
	}

	return ret
}

func toProtoVirtualMachineClones(cs []vm.VirtualMachineClone) []*pb.VirtualMachine_Clone {
	ret := []*pb.VirtualMachine_Clone{}

	for i := range cs {
		ret = append(ret, toProtoVirtualMachineClone(&cs[i]))
	}

	return ret
}

func toProtoVirtualMachineClone(c *vm.VirtualMachineClone) *pb.VirtualMachine_Clone {
	ret := &pb.VirtualMachine_Clone{}
	ret.SetName(c.Name)
	ret.SetNamespace(c.Namespace)
	ret.SetSourceName(c.Spec.Source.Name)
	ret.SetTargetName(c.Spec.Target.Name)
	ret.SetPhase(string(c.Status.Phase))
	ret.SetCreatedAt(timestamppb.New(c.CreationTimestamp.Time))

	conditions := c.Status.Conditions

	if len(conditions) > 0 {
		index := len(conditions) - 1
		ret.SetLastCondition(toProtoApplicationConditionFromClone(&c.Status.Conditions[index]))
	}

	return ret
}

func toProtoVirtualMachineSnapshots(ss []vm.VirtualMachineSnapshot) []*pb.VirtualMachine_Snapshot {
	ret := []*pb.VirtualMachine_Snapshot{}

	for i := range ss {
		ret = append(ret, toProtoVirtualMachineSnapshot(&ss[i]))
	}

	return ret
}

func toProtoVirtualMachineSnapshot(s *vm.VirtualMachineSnapshot) *pb.VirtualMachine_Snapshot {
	ret := &pb.VirtualMachine_Snapshot{}
	ret.SetName(s.Name)
	ret.SetNamespace(s.Namespace)
	ret.SetSourceName(s.Spec.Source.Name)

	if s.Status != nil {
		ret.SetPhase(string(s.Status.Phase))

		if s.Status.ReadyToUse != nil {
			ret.SetReadyToUse(*s.Status.ReadyToUse)
		}

		conditions := s.Status.Conditions

		if len(conditions) > 0 {
			index := len(conditions) - 1
			ret.SetLastCondition(toProtoApplicationConditionFromSnapshot(&s.Status.Conditions[index]))
		}
	}

	ret.SetCreatedAt(timestamppb.New(s.CreationTimestamp.Time))

	return ret
}

func toProtoVirtualMachineRestores(rs []vm.VirtualMachineRestore) []*pb.VirtualMachine_Restore {
	ret := []*pb.VirtualMachine_Restore{}

	for i := range rs {
		ret = append(ret, toProtoVirtualMachineRestore(&rs[i]))
	}

	return ret
}

func toProtoVirtualMachineRestore(r *vm.VirtualMachineRestore) *pb.VirtualMachine_Restore {
	ret := &pb.VirtualMachine_Restore{}
	ret.SetName(r.Name)
	ret.SetNamespace(r.Namespace)
	ret.SetTargetName(r.Spec.Target.Name)

	if r.Status != nil && r.Status.Complete != nil {
		ret.SetComplete(*r.Status.Complete)

		conditions := r.Status.Conditions

		if len(conditions) > 0 {
			index := len(conditions) - 1
			ret.SetLastCondition(toProtoApplicationConditionFromSnapshot(&r.Status.Conditions[index]))
		}
	}

	ret.SetCreatedAt(timestamppb.New(r.CreationTimestamp.Time))

	return ret
}

func getDataVolumeSize(spec *cdi.DataVolumeSpec) int64 {
	if spec.PVC != nil {
		if size := extractStorageSize(spec.PVC.Resources.Requests); size > 0 {
			return size
		}
	}

	if spec.Storage != nil {
		if size := extractStorageSize(spec.Storage.Resources.Requests); size > 0 {
			return size
		}
	}

	return 0
}

func extractStorageSize(requests workload.ResourceList) int64 {
	if requests == nil {
		return 0
	}

	if s, ok := requests[workload.ResourceStorage]; ok {
		if v, ok := s.AsInt64(); ok {
			return v
		}
	}

	return 0
}

func toProtoDataVolumeSource(s *cdi.DataVolumeSource) *pb.DataVolume_Source {
	ret := &pb.DataVolume_Source{}

	switch {
	case s.Blank != nil:
		ret.SetType(pb.DataVolume_Source_TYPE_BLANK_IMAGE)
		ret.SetData("")

	case s.HTTP != nil:
		ret.SetType(pb.DataVolume_Source_TYPE_HTTP_URL)
		ret.SetData(s.HTTP.URL)

	case s.PVC != nil:
		ret.SetType(pb.DataVolume_Source_TYPE_EXISTING_PERSISTENT_VOLUME_CLAIM)
		ret.SetData(s.PVC.Name)
	}

	return ret
}

func toProtoDataVolumeCondition(c *cdi.DataVolumeCondition) *pb.DataVolume_Condition {
	ret := &pb.DataVolume_Condition{}
	ret.SetType(string(c.Type))
	ret.SetStatus(string(c.Status))
	ret.SetReason((c.Reason))
	ret.SetMessage((c.Message))

	if !c.LastHeartbeatTime.IsZero() {
		ret.SetHeartbeatAt(timestamppb.New(c.LastHeartbeatTime.Time))
	}

	if !c.LastTransitionTime.IsZero() {
		ret.SetTransitionedAt(timestamppb.New(c.LastTransitionTime.Time))
	}

	return ret
}

func toProtoDataVolumes(its []cdi.DataVolumePersistent) []*pb.DataVolume {
	ret := []*pb.DataVolume{}

	for i := range its {
		ret = append(ret, toProtoDataVolume(&its[i]))
	}

	return ret
}

func toProtoDataVolume(it *cdi.DataVolumePersistent) *pb.DataVolume {
	ret := &pb.DataVolume{}
	ret.SetName(it.Name)
	ret.SetNamespace(it.Namespace)
	ret.SetSource(toProtoDataVolumeSource(it.Spec.Source))
	ret.SetBootImage(it.BootImage)
	ret.SetPhase(string(it.Status.Phase))
	ret.SetProgress(string(it.Status.Progress))
	ret.SetSizeBytes(getDataVolumeSize(&it.Spec))

	if it.Persistent != nil {
		ret.SetPersistentVolumeClaim(toProtoPersistentVolumeClaim(it.Persistent))
	}

	conditions := it.Status.Conditions

	if len(conditions) > 0 {
		index := len(conditions) - 1
		ret.SetLastCondition(toProtoDataVolumeCondition(&it.Status.Conditions[index]))
	}

	return ret
}

func toProtoInstanceTypes(its []vmi.VirtualMachineInstanceTypeData) []*pb.InstanceType {
	ret := []*pb.InstanceType{}

	for i := range its {
		ret = append(ret, toProtoInstanceType(&its[i]))
	}

	return ret
}

func toProtoInstanceType(it *vmi.VirtualMachineInstanceTypeData) *pb.InstanceType {
	ret := &pb.InstanceType{}
	ret.SetName(it.Type.Name)
	ret.SetCpuCores(it.Type.Spec.CPU.Guest)
	ret.SetMemoryBytes(it.Type.Spec.Memory.Guest.Value())
	ret.SetClusterWide(it.ClusterWide)
	ret.SetCreatedAt(timestamppb.New(it.Type.CreationTimestamp.Time))
	return ret
}
