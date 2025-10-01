package app

import (
	"context"
	"strings"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	clonev1beta1 "kubevirt.io/api/clone/v1beta1"
	virtv1 "kubevirt.io/api/core/v1"
	snapshotv1beta1 "kubevirt.io/api/snapshot/v1beta1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"

	apppb "github.com/otterscale/otterscale/api/application/v1"
	pb "github.com/otterscale/otterscale/api/virtual_machine/v1"
	"github.com/otterscale/otterscale/api/virtual_machine/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type VirtualMachineService struct {
	pbconnect.UnimplementedVirtualMachineServiceHandler

	uc *core.VirtualMachineUseCase
}

func NewVirtualMachineService(uc *core.VirtualMachineUseCase) *VirtualMachineService {
	return &VirtualMachineService{uc: uc}
}

var _ pbconnect.VirtualMachineServiceHandler = (*VirtualMachineService)(nil)

func (s *VirtualMachineService) CheckInfrastructureStatus(ctx context.Context, req *pb.CheckInfrastructureStatusRequest) (*pb.CheckInfrastructureStatusResponse, error) {
	result, err := s.uc.CheckInfrastructureStatus(ctx, req.GetScopeUuid(), req.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.CheckInfrastructureStatusResponse{}
	resp.SetResult(pb.CheckInfrastructureStatusResponse_Result(result))
	return resp, nil
}

func (s *VirtualMachineService) ListVirtualMachines(ctx context.Context, req *pb.ListVirtualMachinesRequest) (*pb.ListVirtualMachinesResponse, error) {
	vms, err := s.uc.ListVirtualMachines(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace())
	if err != nil {
		return nil, err
	}
	its, err := s.uc.ListInstanceTypes(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), true)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListVirtualMachinesResponse{}
	resp.SetVirtualMachines(toProtoVirtualMachines(vms, its))
	return resp, nil
}

func (s *VirtualMachineService) GetVirtualMachine(ctx context.Context, req *pb.GetVirtualMachineRequest) (*pb.VirtualMachine, error) {
	vm, err := s.uc.GetVirtualMachine(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}
	its, err := s.uc.ListInstanceTypes(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), true)
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachine(vm, its)
	return resp, nil
}

func (s *VirtualMachineService) CreateVirtualMachine(ctx context.Context, req *pb.CreateVirtualMachineRequest) (*pb.VirtualMachine, error) {
	vm, err := s.uc.CreateVirtualMachine(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), req.GetInstanceTypeName(), req.GetBootDataVolumeName(), req.GetStartupScript())
	if err != nil {
		return nil, err
	}
	its, err := s.uc.ListInstanceTypes(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), true)
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachine(vm, its)
	return resp, nil
}

func (s *VirtualMachineService) DeleteVirtualMachine(ctx context.Context, req *pb.DeleteVirtualMachineRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteVirtualMachine(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) AttachVirtualMachineDisk(ctx context.Context, req *pb.AttachVirtualMachineDiskRequest) (*pb.VirtualMachine_Disk, error) {
	disk, volume, err := s.uc.AttachVirtualMachineDisk(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), req.GetDataVolumeName())
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachineDisk(disk, []virtv1.Volume{*volume})
	return resp, nil
}

func (s *VirtualMachineService) DetachVirtualMachineDisk(ctx context.Context, req *pb.DetachVirtualMachineDiskRequest) (*emptypb.Empty, error) {
	if err := s.uc.DetachVirtualMachineDisk(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), req.GetDataVolumeName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) CreateVirtualMachineClone(ctx context.Context, req *pb.CreateVirtualMachineCloneRequest) (*pb.VirtualMachine_Clone, error) {
	clone, err := s.uc.CreateVirtualMachineClone(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), req.GetSourceVirtualMachineName(), req.GetTargetVirtualMachineName())
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachineClone(clone)
	return resp, nil
}

func (s *VirtualMachineService) DeleteVirtualMachineClone(ctx context.Context, req *pb.DeleteVirtualMachineCloneRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteVirtualMachineClone(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) CreateVirtualMachineSnapshot(ctx context.Context, req *pb.CreateVirtualMachineSnapshotRequest) (*pb.VirtualMachine_Snapshot, error) {
	snapshot, err := s.uc.CreateVirtualMachineSnapshot(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), req.GetVirtualMachineName())
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachineSnapshot(snapshot)
	return resp, nil
}

func (s *VirtualMachineService) DeleteVirtualMachineSnapshot(ctx context.Context, req *pb.DeleteVirtualMachineSnapshotRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteVirtualMachineSnapshot(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) CreateVirtualMachineRestore(ctx context.Context, req *pb.CreateVirtualMachineRestoreRequest) (*pb.VirtualMachine_Restore, error) {
	restore, err := s.uc.CreateVirtualMachineRestore(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), req.GetVirtualMachineName(), req.GetSnapshotName())
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachineRestore(restore)
	return resp, nil
}

func (s *VirtualMachineService) DeleteVirtualMachineRestore(ctx context.Context, req *pb.DeleteVirtualMachineRestoreRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteVirtualMachineRestore(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) StartVirtualMachine(ctx context.Context, req *pb.StartVirtualMachineRequest) (*emptypb.Empty, error) {
	if err := s.uc.StartVirtualMachine(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) StopVirtualMachine(ctx context.Context, req *pb.StopVirtualMachineRequest) (*emptypb.Empty, error) {
	if err := s.uc.StopVirtualMachine(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) RestartVirtualMachine(ctx context.Context, req *pb.RestartVirtualMachineRequest) (*emptypb.Empty, error) {
	if err := s.uc.RestartVirtualMachine(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) PauseInstance(ctx context.Context, req *pb.PauseInstanceRequest) (*emptypb.Empty, error) {
	if err := s.uc.PauseInstance(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) ResumeInstance(ctx context.Context, req *pb.ResumeInstanceRequest) (*emptypb.Empty, error) {
	if err := s.uc.ResumeInstance(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) MigrateInstance(ctx context.Context, req *pb.MigrateInstanceRequest) (*emptypb.Empty, error) {
	if err := s.uc.MigrateInstance(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), req.GetHostname()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) ListDataVolumes(ctx context.Context, req *pb.ListDataVolumesRequest) (*pb.ListDataVolumesResponse, error) {
	its, err := s.uc.ListDataVolumes(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetBootImage())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListDataVolumesResponse{}
	resp.SetDataVolumes(toProtoDataVolumes(its))
	return resp, nil
}

func (s *VirtualMachineService) GetDataVolume(ctx context.Context, req *pb.GetDataVolumeRequest) (*pb.DataVolume, error) {
	it, err := s.uc.GetDataVolume(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoDataVolume(it)
	return resp, nil
}

func (s *VirtualMachineService) CreateDataVolume(ctx context.Context, req *pb.CreateDataVolumeRequest) (*pb.DataVolume, error) {
	src := req.GetSource()
	it, err := s.uc.CreateDataVolume(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), core.SourceType(src.GetType()), src.GetData(), req.GetSizeBytes(), req.GetBootImage())
	if err != nil {
		return nil, err
	}
	resp := toProtoDataVolume(it)
	return resp, nil
}

func (s *VirtualMachineService) DeleteDataVolume(ctx context.Context, req *pb.DeleteDataVolumeRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteDataVolume(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) ExtendDataVolume(ctx context.Context, req *pb.ExtendDataVolumeRequest) (*emptypb.Empty, error) {
	if err := s.uc.ExtendDataVolume(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), req.GetSizeBytes()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) ListInstanceTypes(ctx context.Context, req *pb.ListInstanceTypesRequest) (*pb.ListInstanceTypesResponse, error) {
	its, err := s.uc.ListInstanceTypes(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetIncludeClusterWide())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListInstanceTypesResponse{}
	resp.SetInstanceTypes(toProtoInstanceTypes(its))
	return resp, nil
}

func (s *VirtualMachineService) GetInstanceType(ctx context.Context, req *pb.GetInstanceTypeRequest) (*pb.InstanceType, error) {
	it, err := s.uc.GetInstanceType(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoInstanceType(it)
	return resp, nil
}

func (s *VirtualMachineService) CreateInstanceType(ctx context.Context, req *pb.CreateInstanceTypeRequest) (*pb.InstanceType, error) {
	it, err := s.uc.CreateInstanceType(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), req.GetCpuCores(), req.GetMemoryBytes())
	if err != nil {
		return nil, err
	}
	resp := toProtoInstanceType(it)
	return resp, nil
}

func (s *VirtualMachineService) DeleteInstanceType(ctx context.Context, req *pb.DeleteInstanceTypeRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteInstanceType(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func (s *VirtualMachineService) CreateVirtualMachineService(ctx context.Context, req *pb.CreateVirtualMachineServiceRequest) (*apppb.Application_Service, error) {
	svc, err := s.uc.CreateVirtualMachineService(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), req.GetVirtualMachineName(), toPorts(req.GetPorts()))
	if err != nil {
		return nil, err
	}
	resp := toProtoService(svc)
	return resp, nil
}

func (s *VirtualMachineService) UpdateVirtualMachineService(ctx context.Context, req *pb.UpdateVirtualMachineServiceRequest) (*apppb.Application_Service, error) {
	svc, err := s.uc.UpdateVirtualMachineService(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName(), toPorts(req.GetPorts()))
	if err != nil {
		return nil, err
	}
	resp := toProtoService(svc)
	return resp, nil
}

func (s *VirtualMachineService) DeleteVirtualMachineService(ctx context.Context, req *pb.DeleteVirtualMachineServiceRequest) (*emptypb.Empty, error) {
	if err := s.uc.DeleteVirtualMachineService(ctx, req.GetScopeUuid(), req.GetFacilityName(), req.GetNamespace(), req.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return resp, nil
}

func toPorts(ps []*apppb.Application_Service_Port) []corev1.ServicePort {
	ret := []corev1.ServicePort{}
	for i := range ps {
		ret = append(ret, toPort(ps[i]))
	}
	return ret
}

func toPort(p *apppb.Application_Service_Port) corev1.ServicePort {
	return corev1.ServicePort{
		Name:       p.GetName(),
		Port:       p.GetPort(),
		NodePort:   p.GetNodePort(),
		Protocol:   corev1.Protocol(strings.ToUpper(p.GetProtocol())),
		TargetPort: intstr.Parse(p.GetTargetPort()),
	}
}

func toProtoVirtualMachineDiskVolumeSource(v *virtv1.VolumeSource) *pb.VirtualMachine_Disk_Volume_Source {
	ret := &pb.VirtualMachine_Disk_Volume_Source{}
	if v.DataVolume != nil {
		ret.SetType(pb.VirtualMachine_Disk_Volume_Source_DATA_VOLUME)
		ret.SetData(v.DataVolume.Name)
	} else if v.CloudInitNoCloud != nil {
		ret.SetType(pb.VirtualMachine_Disk_Volume_Source_CLOUD_INIT_NO_CLOUD)
		ret.SetData(v.CloudInitNoCloud.UserData)
	}
	return ret
}

func toProtoVirtualMachineDiskVolume(v *virtv1.Volume) *pb.VirtualMachine_Disk_Volume {
	ret := &pb.VirtualMachine_Disk_Volume{}
	ret.SetName(v.Name)
	ret.SetSource(toProtoVirtualMachineDiskVolumeSource(&v.VolumeSource))
	return ret
}

func toProtoVirtualMachineDisks(ds []virtv1.Disk, vs []virtv1.Volume) []*pb.VirtualMachine_Disk {
	ret := []*pb.VirtualMachine_Disk{}
	for i := range ds {
		ret = append(ret, toProtoVirtualMachineDisk(&ds[i], vs))
	}
	return ret
}

func toProtoVirtualMachineDisk(d *virtv1.Disk, vs []virtv1.Volume) *pb.VirtualMachine_Disk {
	ret := &pb.VirtualMachine_Disk{}
	ret.SetName(d.Name)
	disk := d.Disk
	if disk != nil {
		ret.SetBus(convertDiskBusToProto(disk.Bus))
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

func convertDiskBusToProto(bus virtv1.DiskBus) pb.VirtualMachine_Disk_Bus {
	switch bus {
	case virtv1.DiskBusVirtio:
		return pb.VirtualMachine_Disk_VIRTIO
	case virtv1.DiskBusSATA:
		return pb.VirtualMachine_Disk_SATA
	case virtv1.DiskBusSCSI:
		return pb.VirtualMachine_Disk_SCSI
	case virtv1.DiskBusUSB:
		return pb.VirtualMachine_Disk_USB
	default:
		return pb.VirtualMachine_Disk_VIRTIO
	}
}

func toProtoVirtualMachines(vmds []core.VirtualMachineData, its []core.VirtualMachineInstanceTypeData) []*pb.VirtualMachine {
	ret := []*pb.VirtualMachine{}
	for i := range vmds {
		ret = append(ret, toProtoVirtualMachine(&vmds[i], its))
	}
	return ret
}

func toProtoVirtualMachine(vmd *core.VirtualMachineData, its []core.VirtualMachineInstanceTypeData) *pb.VirtualMachine {
	ret := &pb.VirtualMachine{}
	ret.SetName(vmd.VirtualMachine.Name)
	ret.SetNamespace(vmd.VirtualMachine.Namespace)

	instanceType := vmd.VirtualMachine.Spec.Instancetype
	if instanceType != nil {
		for _, it := range its {
			if (it.ClusterWide && it.VirtualMachineClusterInstanceType.Name == instanceType.Name) ||
				(!it.ClusterWide && it.VirtualMachineInstanceType.Namespace == vmd.VirtualMachine.Namespace && it.VirtualMachineInstanceType.Name == instanceType.Name) {
				ret.SetInstanceType(toProtoInstanceType(&it))
				break
			}
		}
	}

	ret.SetStatus(string(vmd.VirtualMachine.Status.PrintableStatus))
	ret.SetReady(vmd.VirtualMachine.Status.Ready)

	instance := vmd.VirtualMachineInstance
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

	ret.SetCreatedAt(timestamppb.New(vmd.VirtualMachine.CreationTimestamp.Time))
	ret.SetServices(toProtoServices(vmd.Services))
	ret.SetDisks(toProtoVirtualMachineDisks(vmd.VirtualMachine.Spec.Template.Spec.Domain.Devices.Disks, vmd.VirtualMachine.Spec.Template.Spec.Volumes))
	ret.SetClones(toProtoVirtualMachineClones(vmd.Clones))
	ret.SetSnapshots(toProtoVirtualMachineSnapshots(vmd.Snapshots))
	ret.SetRestores(toProtoVirtualMachineRestores(vmd.Restores))
	return ret
}

func toProtoApplicationConditionFromClone(c *clonev1beta1.Condition) *apppb.Application_Condition {
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

func toProtoApplicationConditionFromSnapshot(c *snapshotv1beta1.Condition) *apppb.Application_Condition {
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

func toProtoVirtualMachineClones(cs []core.VirtualMachineClone) []*pb.VirtualMachine_Clone {
	ret := []*pb.VirtualMachine_Clone{}
	for i := range cs {
		ret = append(ret, toProtoVirtualMachineClone(&cs[i]))
	}
	return ret
}

func toProtoVirtualMachineClone(c *core.VirtualMachineClone) *pb.VirtualMachine_Clone {
	ret := &pb.VirtualMachine_Clone{}
	ret.SetName(c.Name)
	ret.SetNamespace(c.Namespace)
	ret.SetSourceName(c.Spec.Source.Name)
	ret.SetTargetName(c.Spec.Target.Name)
	ret.SetPhase(string(c.Status.Phase))
	ret.SetCreatedAt(timestamppb.New(c.CreationTimestamp.Time))
	if len(c.Status.Conditions) > 0 {
		index := len(c.Status.Conditions) - 1
		ret.SetLastCondition(toProtoApplicationConditionFromClone(&c.Status.Conditions[index]))
	}
	return ret
}

func toProtoVirtualMachineSnapshots(ss []core.VirtualMachineSnapshot) []*pb.VirtualMachine_Snapshot {
	ret := []*pb.VirtualMachine_Snapshot{}
	for i := range ss {
		ret = append(ret, toProtoVirtualMachineSnapshot(&ss[i]))
	}
	return ret
}

func toProtoVirtualMachineSnapshot(s *core.VirtualMachineSnapshot) *pb.VirtualMachine_Snapshot {
	ret := &pb.VirtualMachine_Snapshot{}
	ret.SetName(s.Name)
	ret.SetNamespace(s.Namespace)
	ret.SetSourceName(s.Spec.Source.Name)
	if s.Status != nil {
		ret.SetPhase(string(s.Status.Phase))
		if s.Status.ReadyToUse != nil {
			ret.SetReadyToUse(*s.Status.ReadyToUse)
		}
		if len(s.Status.Conditions) > 0 {
			index := len(s.Status.Conditions) - 1
			ret.SetLastCondition(toProtoApplicationConditionFromSnapshot(&s.Status.Conditions[index]))
		}
	}
	ret.SetCreatedAt(timestamppb.New(s.CreationTimestamp.Time))
	return ret
}

func toProtoVirtualMachineRestores(rs []core.VirtualMachineRestore) []*pb.VirtualMachine_Restore {
	ret := []*pb.VirtualMachine_Restore{}
	for i := range rs {
		ret = append(ret, toProtoVirtualMachineRestore(&rs[i]))
	}
	return ret
}

func toProtoVirtualMachineRestore(r *core.VirtualMachineRestore) *pb.VirtualMachine_Restore {
	ret := &pb.VirtualMachine_Restore{}
	ret.SetName(r.Name)
	ret.SetNamespace(r.Namespace)
	ret.SetTargetName(r.Spec.Target.Name)
	if r.Status != nil && r.Status.Complete != nil {
		ret.SetComplete(*r.Status.Complete)
		if len(r.Status.Conditions) > 0 {
			index := len(r.Status.Conditions) - 1
			ret.SetLastCondition(toProtoApplicationConditionFromSnapshot(&r.Status.Conditions[index]))
		}
	}
	ret.SetCreatedAt(timestamppb.New(r.CreationTimestamp.Time))
	return ret
}

func getDataVolumeSize(spec *cdiv1beta1.DataVolumeSpec) int64 {
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

func extractStorageSize(requests corev1.ResourceList) int64 {
	if requests == nil {
		return 0
	}
	if s, ok := requests[corev1.ResourceStorage]; ok {
		if v, ok := s.AsInt64(); ok {
			return v
		}
	}
	return 0
}

func toProtoDataVolumeSource(s *cdiv1beta1.DataVolumeSource) *pb.DataVolume_Source {
	ret := &pb.DataVolume_Source{}
	switch {
	case s.Blank != nil:
		ret.SetType(pb.DataVolume_Source_BLANK_IMAGE)
		ret.SetData("")
	case s.HTTP != nil:
		ret.SetType(pb.DataVolume_Source_HTTP_URL)
		ret.SetData(s.HTTP.URL)
	case s.PVC != nil:
		ret.SetType(pb.DataVolume_Source_EXISTING_PERSISTENT_VOLUME_CLAIM)
		ret.SetData(s.PVC.Name)
	}
	return ret
}

func toProtoDataVolumeCondition(c *cdiv1beta1.DataVolumeCondition) *pb.DataVolume_Condition {
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

func toProtoDataVolumes(its []core.DataVolumeWithStorage) []*pb.DataVolume {
	ret := []*pb.DataVolume{}
	for i := range its {
		ret = append(ret, toProtoDataVolume(&its[i]))
	}
	return ret
}

func toProtoDataVolume(it *core.DataVolumeWithStorage) *pb.DataVolume {
	ret := &pb.DataVolume{}
	ret.SetName(it.Name)
	ret.SetNamespace(it.Namespace)
	ret.SetSource(toProtoDataVolumeSource(it.Spec.Source))
	ret.SetBootImage(it.Labels[core.DataVolumeBootImageLabel] == "true")
	ret.SetPhase(string(it.Status.Phase))
	ret.SetProgress(string(it.Status.Progress))
	ret.SetSizeBytes(getDataVolumeSize(&it.Spec))
	if it.Storage != nil {
		ret.SetPersistentVolumeClaim(toProtoPersistentVolumeClaim(it.Storage))
	}
	if len(it.Status.Conditions) > 0 {
		index := len(it.Status.Conditions) - 1
		ret.SetLastCondition(toProtoDataVolumeCondition(&it.Status.Conditions[index]))
	}
	return ret
}

func toProtoClusterInstanceTypes(its []core.VirtualMachineClusterInstanceType) []*pb.InstanceType {
	ret := []*pb.InstanceType{}
	for i := range its {
		ret = append(ret, toProtoClusterInstanceType(&its[i]))
	}
	return ret
}

func toProtoClusterInstanceType(it *core.VirtualMachineClusterInstanceType) *pb.InstanceType {
	ret := &pb.InstanceType{}
	ret.SetName(it.Name)
	ret.SetNamespace(it.Namespace)
	ret.SetCpuCores(it.Spec.CPU.Guest)
	ret.SetMemoryBytes(it.Spec.Memory.Guest.Value())
	return ret
}

func toProtoInstanceTypes(its []core.VirtualMachineInstanceTypeData) []*pb.InstanceType {
	ret := []*pb.InstanceType{}
	for i := range its {
		ret = append(ret, toProtoInstanceType(&its[i]))
	}
	return ret
}

func toProtoInstanceType(it *core.VirtualMachineInstanceTypeData) *pb.InstanceType {
	if it.ClusterWide {
		return toProtoInstanceTypeFromClusterInstanceType(it.VirtualMachineClusterInstanceType)
	}
	return toProtoInstanceTypeFromInstanceType(it.VirtualMachineInstanceType)
}

func toProtoInstanceTypeFromClusterInstanceType(it *core.VirtualMachineClusterInstanceType) *pb.InstanceType {
	ret := &pb.InstanceType{}
	ret.SetName(it.Name)
	ret.SetNamespace(it.Namespace)
	ret.SetCpuCores(it.Spec.CPU.Guest)
	ret.SetMemoryBytes(it.Spec.Memory.Guest.Value())
	ret.SetClusterWide(true)
	ret.SetCreatedAt(timestamppb.New(it.CreationTimestamp.Time))
	return ret
}

func toProtoInstanceTypeFromInstanceType(it *core.VirtualMachineInstanceType) *pb.InstanceType {
	ret := &pb.InstanceType{}
	ret.SetName(it.Name)
	ret.SetNamespace(it.Namespace)
	ret.SetCpuCores(it.Spec.CPU.Guest)
	ret.SetMemoryBytes(it.Spec.Memory.Guest.Value())
	ret.SetCreatedAt(timestamppb.New(it.CreationTimestamp.Time))
	return ret
}
