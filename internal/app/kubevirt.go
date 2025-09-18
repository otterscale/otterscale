package app

import (
	"context"
	"strings"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	virtCorev1 "kubevirt.io/api/core/v1"

	pb "github.com/otterscale/otterscale/api/kubevirt/v1"
	"github.com/otterscale/otterscale/api/kubevirt/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type KubeVirtService struct {
	pbconnect.UnimplementedKubeVirtServiceHandler

	uc *core.KubeVirtUseCase
}

func NewKubeVirtService(uc *core.KubeVirtUseCase) *KubeVirtService {
	return &KubeVirtService{uc: uc}
}

var _ pbconnect.KubeVirtServiceHandler = (*KubeVirtService)(nil)

// ListNamespaces returns a list of namespaces
func (s *KubeVirtService) ListNamespaces(ctx context.Context, req *connect.Request[pb.ListNamespacesRequest]) (*connect.Response[pb.ListNamespacesResponse], error) {
	namespaces, err := s.uc.ListNamespaces(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListNamespacesResponse{}
	namespaceNames := make([]string, len(namespaces))
	for i, ns := range namespaces {
		namespaceNames[i] = ns.Name
	}
	resp.SetNamespaces(namespaceNames)
	return connect.NewResponse(resp), nil
}

// ListPersistentVolumeClaims returns a list of persistent volume claims
func (s *KubeVirtService) ListBootablePersistentVolumeClaims(ctx context.Context, req *connect.Request[pb.ListPersistentVolumeClaimsRequest]) (*connect.Response[pb.ListPersistentVolumeClaimsResponse], error) {
	pvcs, err := s.uc.ListBootablePersistentVolumeClaims(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListPersistentVolumeClaimsResponse{}
	pvcList := make([]*pb.PersistentVolumeClaim, 0, len(pvcs))

	for i := range pvcs {
		pbPvc := &pb.PersistentVolumeClaim{}

		// Set metadata
		pbPvc.SetName(pvcs[i].GetName())
		pbPvc.SetNamespace(pvcs[i].GetNamespace())
		if pvcs[i].Spec.StorageClassName != nil {
			pbPvc.SetStorageClass(*pvcs[i].Spec.StorageClassName)
		}
		if pvcs[i].Spec.Resources.Requests != nil {
			if storageReq, ok := pvcs[i].Spec.Resources.Requests[corev1.ResourceStorage]; ok {
				pbPvc.SetSizeBytes(storageReq.Value())
			}
		}
		if len(pvcs[i].Spec.AccessModes) > 0 {
			pbPvc.SetAccessMode(string(pvcs[i].Spec.AccessModes[0]))
		}
		pbPvc.SetStatus(string(pvcs[i].Status.Phase))

		pvcList = append(pvcList, pbPvc)
	}

	resp.SetPersistentVolumeClaims(pvcList)
	return connect.NewResponse(resp), nil
}

// Virtual Machine Operations
func (s *KubeVirtService) CreateVirtualMachine(ctx context.Context, req *connect.Request[pb.CreateVirtualMachineRequest]) (*connect.Response[pb.VirtualMachineResponse], error) {
	// Extract disks with their DataVolumeSource information
	disks := req.Msg.GetDisks()
	diskDevices, dataVolumeSources := toCoreDiskDevicesWithDataVolumes(disks)

	vm, err := s.uc.CreateVirtualMachine(
		ctx,
		req.Msg.GetScopeUuid(),
		req.Msg.GetFacilityName(),
		req.Msg.GetNamespace(),
		req.Msg.GetName(),
		req.Msg.GetNetworkName(),
		req.Msg.GetStartupScript(),
		req.Msg.GetLabels(),
		toCoreVirtualMachineResource(req.Msg.GetCustom(), req.Msg.GetInstancetypeName()),
		diskDevices,
		dataVolumeSources,
	)
	if err != nil {
		return nil, err
	}
	vmInfo := &core.VirtualMachineInfo{
		VM: vm,
	}
	resp := toProtoVirtualMachineResponse(vmInfo)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListVirtualMachines(ctx context.Context, req *connect.Request[pb.ListVirtualMachinesRequest]) (*connect.Response[pb.ListVirtualMachinesResponse], error) {
	vmInfoList, err := s.uc.ListVirtualMachines(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListVirtualMachinesResponse{}
	resp.SetVirtualMachines(toProtoVirtualMachineResponses(vmInfoList))
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) UpdateVirtualMachine(ctx context.Context, req *connect.Request[pb.UpdateVirtualMachineRequest]) (*connect.Response[pb.VirtualMachineResponse], error) {
	vm, vmi, err := s.uc.UpdateVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetNetworkName(), req.Msg.GetLabels(), nil)
	if err != nil {
		return nil, err
	}

	vmInfo := &core.VirtualMachineInfo{
		VM:  vm,
		VMI: vmi,
	}
	resp := toProtoVirtualMachineResponse(vmInfo)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) DeleteVirtualMachine(ctx context.Context, req *connect.Request[pb.DeleteVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) CreateVirtualMachineDisk(ctx context.Context, req *connect.Request[pb.CreateVirtualMachineDiskRequest]) (*connect.Response[emptypb.Empty], error) {
	scopeUUID := req.Msg.GetScopeUuid()
	facilityName := req.Msg.GetFacilityName()
	vmName := req.Msg.GetVmName()
	vmNamespace := req.Msg.GetVmNamespace()
	pbDisk := req.Msg.GetDisk()

	disk := fromProtoVirtualMachineDisk(pbDisk)

	var dvInfo *core.DataVolumeInfo
	if strings.ToLower(disk.DiskType) == core.TYPEDATAVOLUME {
		if pbDisk.GetDataVolume() != nil {
			dvInfo = &core.DataVolumeInfo{
				Name:       pbDisk.GetName(),
				SourceType: strings.ToLower(pbDisk.GetDataVolume().GetType().String()),
				Source:     pbDisk.GetDataVolume().GetSource(),
				SizeBytes:  pbDisk.GetDataVolume().GetSizeBytes(),
				IsBootable: pbDisk.GetIsBootable(),
			}
		}
	}

	err := s.uc.CreateVirtualMachineDisk(ctx, scopeUUID, facilityName, vmNamespace, vmName, disk, dvInfo)
	if err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// Virtual Machine Control Operations
func (s *KubeVirtService) StartVirtualMachine(ctx context.Context, req *connect.Request[pb.StartVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.StartVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) StopVirtualMachine(ctx context.Context, req *connect.Request[pb.StopVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.StopVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) PauseVirtualMachine(ctx context.Context, req *connect.Request[pb.PauseVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.PauseVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ResumeVirtualMachine(ctx context.Context, req *connect.Request[pb.ResumeVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.UnpauseVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// Virtual Machine Advanced Operations
func (s *KubeVirtService) CloneVirtualMachine(ctx context.Context, req *connect.Request[pb.CloneVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.CloneVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetTargetNamespace(), req.Msg.GetTargetName(), req.Msg.GetSourceNamespace(), req.Msg.GetSourceName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) SnapshotVirtualMachine(ctx context.Context, req *connect.Request[pb.SnapshotVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.SnapshotVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetSnapshotName(), req.Msg.GetDescription()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) RestoreVirtualMachine(ctx context.Context, req *connect.Request[pb.RestoreVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.RestoreVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetSnapshotName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) MigrateVirtualMachine(ctx context.Context, req *connect.Request[pb.MigrateVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.MigrateVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetTargetNode()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// GetVirtualMachineSnapshot retrieves a virtual machine snapshot
func (s *KubeVirtService) GetVirtualMachineSnapshot(ctx context.Context, req *connect.Request[pb.GetVirtualMachineSnapshotRequest]) (*connect.Response[pb.VirtualMachineSnapshot], error) {
	snapshot, err := s.uc.GetVirtualMachineSnapshot(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachineSnapshot(snapshot)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListVirtualMachineSnapshots(ctx context.Context, req *connect.Request[pb.ListVirtualMachineSnapshotsRequest]) (*connect.Response[pb.ListVirtualMachineSnapshotsResponse], error) {
	resp := &pb.ListVirtualMachineSnapshotsResponse{}
	snapshots, err := s.uc.ListVirtualMachineSnapshots(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetVmName())
	if err != nil {
		return nil, err
	}
	resp.SetSnapshots(toProtoVirtualMachineSnapshots(snapshots))
	return connect.NewResponse(resp), nil
}

// DeleteVirtualMachineSnapshot deletes a virtual machine snapshot
func (s *KubeVirtService) DeleteVirtualMachineSnapshot(ctx context.Context, req *connect.Request[pb.DeleteVirtualMachineSnapshotRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.uc.DeleteVirtualMachineSnapshot(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// Data Volume Operations
func (s *KubeVirtService) CreateDataVolume(ctx context.Context, req *connect.Request[pb.CreateDataVolumeRequest]) (*connect.Response[pb.DataVolume], error) {
	dv, err := s.uc.CreateDataVolume(ctx, req.Msg.GetScopeUuid(),
		req.Msg.GetFacilityName(),
		req.Msg.GetDataVolume().GetMetadata().GetNamespace(),
		req.Msg.GetDataVolume().GetMetadata().GetName(),
		req.Msg.GetDataVolume().GetType(),
		req.Msg.GetDataVolume().GetSource(),
		req.Msg.GetDataVolume().GetSizeBytes(),
		req.Msg.GetIsBootable())
	if err != nil {
		return nil, err
	}

	resp := toProtoDataVolume(dv)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) GetDataVolume(ctx context.Context, req *connect.Request[pb.GetDataVolumeRequest]) (*connect.Response[pb.DataVolume], error) {
	dv, err := s.uc.GetDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoDataVolume(dv)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListDataVolumes(ctx context.Context, req *connect.Request[pb.ListDataVolumesRequest]) (*connect.Response[pb.ListDataVolumesResponse], error) {
	dvs, err := s.uc.ListDataVolumes(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListDataVolumesResponse{}
	resp.SetDatavolumes(toProtoDataVolumes(dvs))
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) DeleteDataVolume(ctx context.Context, req *connect.Request[pb.DeleteDataVolumeRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ExtendDataVolume(ctx context.Context, req *connect.Request[pb.ExtendDataVolumeRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.ExtendDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetSizeBytes()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// VirtualMachineService Operations
func (s *KubeVirtService) CreateVirtualMachineService(ctx context.Context, req *connect.Request[pb.CreateVirtualMachineServiceRequest]) (*connect.Response[pb.VirtualMachineService], error) {
	vmsvc, err := s.uc.CreateVirtualMachineService(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetVirtualMachineService().GetMetadata().GetNamespace(), req.Msg.GetVirtualMachineService().GetMetadata().GetName(), toCoreVirtualMachineService(req.Msg.GetVirtualMachineService()))
	if err != nil {
		return nil, err
	}

	resp := toProtoVirtualMachineService(vmsvc)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) GetVirtualMachineService(ctx context.Context, req *connect.Request[pb.GetVirtualMachineServiceRequest]) (*connect.Response[pb.VirtualMachineService], error) {
	vmsvc, err := s.uc.GetVirtualMachineService(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachineService(vmsvc)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListVirtualMachineServices(ctx context.Context, req *connect.Request[pb.ListVirtualMachineServicesRequest]) (*connect.Response[pb.ListVirtualMachineServicesResponse], error) {
	vmsvcs, err := s.uc.ListVirtualMachineServices(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListVirtualMachineServicesResponse{}
	resp.SetVirtualMachineServices(toProtoVirtualMachineServices(vmsvcs))
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) UpdateVirtualMachineService(ctx context.Context, req *connect.Request[pb.UpdateVirtualMachineServiceRequest]) (*connect.Response[pb.VirtualMachineService], error) {
	vmsvc, err := s.uc.UpdateVirtualMachineService(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), toCoreVirtualMachineService(req.Msg.GetVirtualMachineService()))
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachineService(vmsvc)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) DeleteVirtualMachineService(ctx context.Context, req *connect.Request[pb.DeleteVirtualMachineServiceRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteVirtualMachineService(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
		return nil, err
	}

	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// InstanceType Operations
func (s *KubeVirtService) CreateInstanceType(ctx context.Context, req *connect.Request[pb.CreateInstanceTypeRequest]) (*connect.Response[pb.InstanceType], error) {
	InstanceType, err := s.uc.CreateInstanceType(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), toCoreInstanceType(req.Msg.GetInstanceType()))
	if err != nil {
		return nil, err
	}
	resp := toProtoInstanceType(InstanceType)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) GetInstanceType(ctx context.Context, req *connect.Request[pb.GetInstanceTypeRequest]) (*connect.Response[pb.InstanceType], error) {
	InstanceType, err := s.uc.GetInstanceType(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoInstanceType(InstanceType)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListInstanceTypes(ctx context.Context, req *connect.Request[pb.ListInstanceTypesRequest]) (*connect.Response[pb.ListInstanceTypesResponse], error) {
	InstanceType, err := s.uc.ListInstanceTypes(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListInstanceTypesResponse{}
	resp.SetInstanceTypes(toProtoInstanceTypes(InstanceType))

	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) DeleteInstanceType(ctx context.Context, req *connect.Request[pb.DeleteInstanceTypeRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteInstanceType(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// Conversion functions
func toProtoVirtualMachineResponses(vmInfoList []core.VirtualMachineInfo) []*pb.VirtualMachineResponse {
	ret := []*pb.VirtualMachineResponse{}
	for i := 0; i < len(vmInfoList); i++ {
		ret = append(ret, toProtoVirtualMachineResponse(&vmInfoList[i]))
	}
	return ret
}

func toProtoVirtualMachineResponse(vmInfo *core.VirtualMachineInfo) *pb.VirtualMachineResponse {
	ret := &pb.VirtualMachineResponse{}
	if vmInfo.VM == nil {
		return ret
	}

	vm := vmInfo.VM
	vmi := vmInfo.VMI

	ret.SetMetadata(fromVirtualMachine(vm))
	ret.SetStartupScript(toProtoVirtualMachineScripts(vm.Spec.Template.Spec.Volumes))
	ret.SetResources(toProtoVirtualMachineResources(vm))
	if vmi != nil {
		ret.SetNodeName(vmi.Status.NodeName)
	}
	if vmInfo.Pod != nil {
		ret.SetPodName(vmInfo.Pod.GetName())
	}
	if len(vm.Spec.Template.Spec.Networks) > 0 {
		ret.SetNetworkName(vm.Spec.Template.Spec.Networks[0].Name)
	}

	ret.SetMaasId(vmInfo.SystemID)

	volumeMap := buildVolumeMap(vm)
	disks := make([]*pb.VirtualMachineResponse_DiskInfo, 0, len(vm.Spec.Template.Spec.Domain.Devices.Disks))

	for i := range vm.Spec.Template.Spec.Domain.Devices.Disks {
		disk := &vm.Spec.Template.Spec.Domain.Devices.Disks[i]
		diskInfo := toProtoVirtualMachineDisk(disk, volumeMap, vm.GetNamespace())

		if diskInfo.GetDiskType() == pb.VirtualMachineResponse_DiskInfo_DATAVOLUME {
			dvName := volumeMap[disk.Name].VolumeSource.DataVolume.Name

			for _, dv := range vmInfo.DataVolume {
				if dv.Name == dvName {
					dvSource := diskInfo.GetDataVolume()
					if dvSource == nil {
						dvSource = &pb.DataVolumeSource{}
					}

					source, sourceType, _, _, sizeBytes := core.ExtractDataVolumeInfo(dv)

					switch strings.ToLower(sourceType) {
					case "http":
						dvSource.SetType(pb.DataVolumeSource_HTTP)
					case "blank":
						dvSource.SetType(pb.DataVolumeSource_BLANK)
					case "pvc":
						dvSource.SetType(pb.DataVolumeSource_PVC)
					default:
						dvSource.SetType(pb.DataVolumeSource_UNSPECIFIED)
					}

					dvSource.SetSource(source)
					dvSource.SetSizeBytes(sizeBytes)

					diskInfo.SetDataVolume(dvSource)
					break
				}
			}
		}
		disks = append(disks, diskInfo)
	}

	ret.SetDisks(disks)

	if vmInfo.Service != nil {
		services := []*pb.VirtualMachineResponse_ServiceInfo{}
		for _, port := range vmInfo.Service.Spec.Ports {
			service := &pb.VirtualMachineResponse_ServiceInfo{}
			if vmInfo.Service.Spec.Type == corev1.ServiceTypeNodePort {
				service.SetStype(pb.VirtualMachineResponse_ServiceInfo_NODE_PORT)
			} else if vmInfo.Service.Spec.Type == corev1.ServiceTypeLoadBalancer {
				service.SetStype(pb.VirtualMachineResponse_ServiceInfo_LOAD_BALANCER)
			} else {
				service.SetStype(pb.VirtualMachineResponse_ServiceInfo_UNSPECIFIED)
			}
			service.SetPort(int64(port.Port))
			service.SetNodePort(int64(port.NodePort))
			services = append(services, service)
		}
		ret.SetServices(services)
	}

	v, ok := pb.VirtualMachineResponseStatus_value[strings.ToUpper(string(vm.Status.PrintableStatus))]
	if ok {
		ret.SetStatusPhase(pb.VirtualMachineResponseStatus(v))
	} else {
		ret.SetStatusPhase(pb.VirtualMachineResponseStatus(pb.VirtualMachineResponse_UNKNOWN))
	}

	return ret
}

func fromVirtualMachine(vm *core.VirtualMachine) *pb.Metadata {
	return toProtoMetadata(vm.GetNamespace(), vm.GetName(), vm.GetLabels(), vm.CreationTimestamp.Time, vm.GetAnnotations()["otterscale.io/last-updated"])
}

func fromDataVolume(dv *core.DataVolume) *pb.Metadata {
	return toProtoMetadata(dv.GetNamespace(), dv.GetName(), dv.GetLabels(), dv.CreationTimestamp.Time, dv.GetAnnotations()["otterscale.io/last-updated"])
}

func fromInstanceType(t *core.InstanceType) *pb.Metadata {
	return toProtoMetadata(t.Metadata.Namespace, t.Metadata.Name, t.Metadata.Labels, t.Metadata.CreatedAt.AsTime(), t.Metadata.Annotations["otterscale.io/last-updated"])
}

func toProtoMetadata(namespace, name string, labels map[string]string, creationTimestamp time.Time, updateTimestamp string) *pb.Metadata {
	ret := &pb.Metadata{}

	ret.SetName(name)
	ret.SetNamespace(namespace)
	ret.SetLabels(labels)
	ret.SetCreatedAt(timestamppb.New(creationTimestamp))
	parsedUpdateTime, err := time.Parse(time.RFC3339, updateTimestamp)
	if err == nil {
		ret.SetUpdatedAt(timestamppb.New(parsedUpdateTime))
	}

	return ret
}

func toCoreMetadata(m *pb.Metadata) core.Metadata {
	return core.Metadata{
		Name:      m.GetName(),
		Namespace: m.GetNamespace(),
		Labels:    m.GetLabels(),
		CreatedAt: m.GetCreatedAt(),
		UpdatedAt: m.GetUpdatedAt(),
	}
}

func toCoreVirtualMachineResource(r *pb.VirtualMachineResources, instanceName string) core.VirtualMachineResources {
	ret := core.VirtualMachineResources{}
	if r != nil {
		ret.CPUcores = r.GetCpuCores()
		ret.MemoryBytes = r.GetMemoryBytes()
	} else {
		ret.InstanceName = instanceName
	}

	return ret
}

func toProtoVirtualMachineScripts(volume []core.KubeVirtVolume) string {
	var userData string
	for i := range volume {
		if volume[i].CloudInitNoCloud != nil { // is a cloud-init volume
			userData = volume[i].CloudInitNoCloud.UserData
			break
		}
	}
	return userData
}

func toProtoVirtualMachineResources(vm *core.VirtualMachine) *pb.VirtualMachineResources {
	ret := &pb.VirtualMachineResources{}
	ret.SetCpuCores(vm.Spec.Template.Spec.Domain.CPU.Cores)
	ret.SetMemoryBytes(vm.Spec.Template.Spec.Domain.Resources.Requests.Memory().Value())
	return ret
}

func buildVolumeMap(vm *core.VirtualMachine) map[string]virtCorev1.Volume {
	m := make(map[string]virtCorev1.Volume, len(vm.Spec.Template.Spec.Volumes))
	for i := range vm.Spec.Template.Spec.Volumes {
		m[vm.Spec.Template.Spec.Volumes[i].Name] = vm.Spec.Template.Spec.Volumes[i]
	}
	return m
}

func toProtoVirtualMachineDisk(disk *virtCorev1.Disk, vmVols map[string]virtCorev1.Volume, namespace string) *pb.VirtualMachineResponse_DiskInfo {
	ret := &pb.VirtualMachineResponse_DiskInfo{}
	ret.SetName(disk.Name)
	ret.SetNamespace(namespace)
	if disk.BootOrder != nil {
		ret.SetBootOrder(int64(*disk.BootOrder))
	}
	if disk.Disk != nil {
		v, ok := pb.VirtualMachineDiskBus_value[strings.ToUpper(string(disk.Disk.Bus))]
		if ok {
			ret.SetBusType(pb.VirtualMachineResponse_DiskInfoBus(v))
		} else {
			ret.SetBusType(pb.VirtualMachineResponse_DiskInfoBus(pb.VirtualMachineDisk_VIRTIO))
		}
	}
	ret.SetDiskType(pb.VirtualMachineResponse_DiskInfo_UNSPECIFIED)

	if vol, ok := vmVols[disk.Name]; ok {
		switch {
		case vol.VolumeSource.DataVolume != nil:
			ret.SetDiskType(pb.VirtualMachineResponse_DiskInfo_DATAVOLUME)

			dvSource := &pb.DataVolumeSource{}

			dvSource.SetType(pb.DataVolumeSource_PVC) // 預設為 PVC 類型
			dvSource.SetSource(vol.VolumeSource.DataVolume.Name)

			ret.SetDataVolume(dvSource)

		case vol.VolumeSource.PersistentVolumeClaim != nil:
			ret.SetDiskType(pb.VirtualMachineResponse_DiskInfo_PERSISTENTVOLUMECLAIM)
			ret.SetSource(vol.VolumeSource.PersistentVolumeClaim.ClaimName)
		case vol.VolumeSource.ConfigMap != nil:
			ret.SetDiskType(pb.VirtualMachineResponse_DiskInfo_CONFIGMAP)
			ret.SetSource(vol.VolumeSource.ConfigMap.Name)
		case vol.VolumeSource.Secret != nil:
			ret.SetDiskType(pb.VirtualMachineResponse_DiskInfo_SECRET)
			ret.SetSource(vol.VolumeSource.Secret.SecretName)
		case vol.VolumeSource.CloudInitNoCloud != nil:
			ret.SetDiskType(pb.VirtualMachineResponse_DiskInfo_CLOUDINITNOCLOUD)
		}
	}

	return ret
}

func fromProtoVirtualMachineDisk(disk *pb.VirtualMachineDisk) core.DiskDevice {
	var diskType string
	switch disk.GetDiskType() {
	case pb.VirtualMachineDisk_DATAVOLUME:
		diskType = core.TYPEDATAVOLUME
	case pb.VirtualMachineDisk_PERSISTENTVOLUMECLAIM:
		diskType = core.TYPEPERSISTENTVOLUMECLAIM
	case pb.VirtualMachineDisk_CONFIGMAP:
		diskType = core.TYPECONFIGMAP
	case pb.VirtualMachineDisk_SECRET:
		diskType = core.TYPESECRET
	case pb.VirtualMachineDisk_CLOUDINITNOCLOUD:
		diskType = core.TYPECLOUDINITNOCLOUD
	default:
		diskType = ""
	}

	var busStr string
	switch disk.GetBusType() {
	case pb.VirtualMachineDisk_SATA:
		busStr = "sata"
	case pb.VirtualMachineDisk_SCSI:
		busStr = "scsi"
	case pb.VirtualMachineDisk_VIRTIO:
		busStr = "virtio"
	default:
		busStr = "virtio"
	}

	return core.DiskDevice{
		Name:     disk.GetName(),
		DiskType: diskType,
		Bus:      busStr,
		Data:     disk.GetSource(),
	}
}

// toCoreDiskDevicesWithDataVolumes extracts both DiskDevice and DataVolumeSource information
func toCoreDiskDevicesWithDataVolumes(disks []*pb.VirtualMachineDisk) ([]core.DiskDevice, map[string]*core.DataVolumeInfo) {
	diskDevices := []core.DiskDevice{}
	dataVolumeSources := make(map[string]*core.DataVolumeInfo)

	for i := range disks {
		// Convert pb.VirtualMachineDisk_bus to string
		var busStr string
		switch disks[i].GetBusType() {
		case pb.VirtualMachineDisk_SATA:
			busStr = "sata"
		case pb.VirtualMachineDisk_SCSI:
			busStr = "scsi"
		case pb.VirtualMachineDisk_VIRTIO:
			busStr = "virtio"
		default:
			busStr = "virtio"
		}

		// Get source based on source_data oneof field
		var source string
		if disks[i].GetSource() != "" {
			source = disks[i].GetSource()
		}

		// Create DiskDevice
		diskDevice := core.DiskDevice{
			Name:     disks[i].GetName(),
			DiskType: pb.VirtualMachineDiskType_name[int32(disks[i].GetDiskType())],
			Bus:      busStr,
			Data:     source,
		}
		diskDevices = append(diskDevices, diskDevice)

		// If disk has a DataVolumeSource, extract it
		if disks[i].GetDiskType() == pb.VirtualMachineDisk_DATAVOLUME {
			// Check if disk has DataVolumeSource info
			if dv := disks[i].GetDataVolume(); dv != nil {
				sourceType := mapDataVolumeSourceType(dv.GetType())
				dataVolumeSources[disks[i].GetName()] = &core.DataVolumeInfo{
					Name:       disks[i].GetName(), // Use disk name as DV name if not specified
					SourceType: sourceType,
					Source:     dv.GetSource(),
					SizeBytes:  dv.GetSizeBytes(),
					IsBootable: disks[i].GetIsBootable(),
				}
			}
		}
	}

	return diskDevices, dataVolumeSources
}

// mapDataVolumeSourceType maps protobuf DataVolumeSource.Type to string constants used in core
func mapDataVolumeSourceType(sourceType pb.DataVolumeSource_Type) string {
	switch sourceType {
	case pb.DataVolumeSource_HTTP:
		return "HTTP"
	case pb.DataVolumeSource_BLANK:
		return "Blank"
	case pb.DataVolumeSource_PVC:
		return "PVC"
	default:
		return "Blank" // Default to Blank if unspecified
	}
}

func toProtoVirtualMachineSnapshots(snapshots []core.VirtualMachineSnapshot) []*pb.VirtualMachineSnapshot {
	ret := []*pb.VirtualMachineSnapshot{}
	for i := range snapshots {
		ret = append(ret, toProtoVirtualMachineSnapshot(&snapshots[i]))
	}
	return ret
}

// toProtoVirtualMachineSnapshot converts core VirtualMachineSnapshot to protobuf
func toProtoVirtualMachineSnapshot(snapshot *core.VirtualMachineSnapshot) *pb.VirtualMachineSnapshot {
	ret := &pb.VirtualMachineSnapshot{}
	ret.SetName(snapshot.GetName())
	ret.SetNamespace(snapshot.GetNamespace())
	ret.SetSourceName(snapshot.Spec.Source.Name)
	ret.SetDescription(snapshot.GetAnnotations()["otterscale.io/snapshot-description"])
	ret.SetSourceNamespace(snapshot.GetNamespace())
	ret.SetCreatedAt(timestamppb.New(snapshot.CreationTimestamp.Time))

	v, ok := pb.VirtualMachineSnapshotStatus_value[strings.ToUpper(string(snapshot.Status.Phase))]
	if ok {
		ret.SetStatusPhase(pb.VirtualMachineSnapshotStatus(v))
	} else {
		ret.SetStatusPhase(pb.VirtualMachineSnapshotStatus(pb.VirtualMachineSnapshot_UNKNOWN))
	}

	for _, cond := range snapshot.Status.Conditions {
		if string(cond.Status) == "True" {
			ret.SetLastConditionMessage(cond.Message)
			ret.SetLastConditionReason(cond.Reason)
			break
		}
	}
	return ret
}

func toProtoDataVolumes(dvs []core.DataVolume) []*pb.DataVolume {
	ret := []*pb.DataVolume{}
	for i := range dvs {
		ret = append(ret, toProtoDataVolume(&dvs[i]))
	}
	return ret
}

func toProtoDataVolume(dv *core.DataVolume) *pb.DataVolume {
	ret := &pb.DataVolume{}
	source, sourceType, accessMode, storageClassName, sizeBytes := core.ExtractDataVolumeInfo(dv)

	ret.SetMetadata(fromDataVolume(dv))
	ret.SetSource(source)
	ret.SetType(sourceType)
	ret.SetSizeBytes(sizeBytes)
	ret.SetAccessMode(accessMode)
	ret.SetStorageClass(storageClassName)

	ret.SetStatusPhase(string(dv.Status.Phase))
	ret.SetStatusProgress(string(dv.Status.Progress))
	ret.SetStatusClaimName(dv.Status.ClaimName)

	if dv.Status.Conditions != nil {
		conditionMessage, conditionReason, conditionStatus := core.GetDataVolumeConditions(dv)
		ret.SetLastConditionMessage(conditionMessage)
		ret.SetLastConditionReason(conditionReason)
		ret.SetLastConditionStatus(conditionStatus)
	}
	return ret
}

func toProtoVirtualMachineServices(vmsvcs []core.VirtualMachineService) []*pb.VirtualMachineService {
	ret := make([]*pb.VirtualMachineService, 0, len(vmsvcs))
	for i := range vmsvcs {
		ret = append(ret, toProtoVirtualMachineService(&vmsvcs[i]))
	}
	return ret
}

func toProtoVirtualMachineService(vmsvc *core.VirtualMachineService) *pb.VirtualMachineService {
	ret := &pb.VirtualMachineService{}
	meta := &pb.Metadata{}
	meta.SetName(vmsvc.GetName())
	meta.SetNamespace(vmsvc.GetNamespace())
	meta.SetLabels(vmsvc.GetLabels())
	meta.SetCreatedAt(timestamppb.New(vmsvc.GetCreationTimestamp().Time))
	if ann := vmsvc.GetAnnotations(); ann != nil {
		if ts := ann["otterscale.io/last-updated"]; ts != "" {
			if t, err := time.Parse(time.RFC3339, ts); err == nil {
				meta.SetUpdatedAt(timestamppb.New(t))
			}
		}
	}

	sp := &pb.VirtualMachineServiceSpec{}
	switch vmsvc.Spec.Type {
	case corev1.ServiceTypeNodePort:
		sp.SetType(pb.VirtualMachineServiceSpec_NODE_PORT)
	case corev1.ServiceTypeLoadBalancer:
		sp.SetType(pb.VirtualMachineServiceSpec_LOAD_BALANCER)
	default:
		sp.SetType(pb.VirtualMachineServiceSpec_TYPE_UNSPECIFIED)
	}
	sp.SetSelector(vmsvc.Spec.Selector)
	if vm, ok := vmsvc.Spec.Selector["otterscale.io/virtualmachine"]; ok {
		sp.SetVirtualMachineName(vm)
	}
	ports := make([]*pb.ServicePort, 0, len(vmsvc.Spec.Ports))
	for _, p := range vmsvc.Spec.Ports {
		pp := &pb.ServicePort{}
		pp.SetName(p.Name)
		pp.SetPort(p.Port)
		if p.Protocol == corev1.ProtocolUDP {
			pp.SetProtocol(pb.ServicePort_UDP)
		} else {
			pp.SetProtocol(pb.ServicePort_TCP)
		}
		pp.SetNodePort(p.NodePort)
		ports = append(ports, pp)
	}
	sp.SetPorts(ports)
	ret.SetSpec(sp)
	st := &pb.VirtualMachineServiceStatus{}
	if vmsvc.Spec.ClusterIP != "" {
		st.SetClusterIp(vmsvc.Spec.ClusterIP)
	}
	if len(vmsvc.Spec.ClusterIPs) > 0 {
		st.SetClusterIps(vmsvc.Spec.ClusterIPs)
	}
	if ingress := vmsvc.Status.LoadBalancer.Ingress; len(ingress) > 0 {
		addrs := make([]string, 0, len(ingress))
		for _, in := range ingress {
			if in.IP != "" {
				addrs = append(addrs, in.IP)
			} else if in.Hostname != "" {
				addrs = append(addrs, in.Hostname)
			}
		}
		st.SetLoadBalancerIngress(addrs)
	}
	ret.SetStatus(st)

	return ret
}

func toCoreVirtualMachineService(vmsvc *pb.VirtualMachineService) *core.VirtualMachineServiceSpec {
	var spec corev1.ServiceSpec
	switch vmsvc.GetSpec().GetType() {
	case pb.VirtualMachineServiceSpec_NODE_PORT:
		spec.Type = corev1.ServiceTypeNodePort
	case pb.VirtualMachineServiceSpec_LOAD_BALANCER:
		spec.Type = corev1.ServiceTypeLoadBalancer
	default:
		spec.Type = corev1.ServiceTypeLoadBalancer
	}

	vmName := vmsvc.GetSpec().GetVirtualMachineName()
	if vmName != "" {
		spec.Selector = map[string]string{"otterscale.io/virtualmachine": vmName}
	} else {
		spec.Selector = map[string]string{}
	}

	for _, p := range vmsvc.GetSpec().GetPorts() {
		sp := corev1.ServicePort{
			Name:       p.GetName(),
			Port:       p.GetPort(),
			TargetPort: intstr.FromInt(int(p.GetPort())),
			Protocol:   corev1.ProtocolTCP,
		}
		if p.GetProtocol() == pb.ServicePort_UDP {
			sp.Protocol = corev1.ProtocolUDP
		}
		if spec.Type == corev1.ServiceTypeNodePort && p.GetNodePort() > 0 {
			sp.NodePort = p.GetNodePort()
		}
		spec.Ports = append(spec.Ports, sp)
	}

	return &spec
}

func toProtoInstanceTypes(flavors []core.InstanceType) []*pb.InstanceType {
	ret := []*pb.InstanceType{}
	for i := range flavors {
		ret = append(ret, toProtoInstanceType(&flavors[i]))
	}
	return ret
}

func toProtoInstanceType(f *core.InstanceType) *pb.InstanceType {
	ret := &pb.InstanceType{}

	ret.SetMetadata(fromInstanceType(f))
	ret.SetCpuCores(f.CPUCores)
	ret.SetMemoryBytes(f.MemoryBytes)

	return ret
}

func toCoreInstanceType(f *pb.InstanceType) *core.InstanceType {
	return &core.InstanceType{
		Metadata:    toCoreMetadata(f.GetMetadata()),
		CPUCores:    f.GetCpuCores(),
		MemoryBytes: f.GetMemoryBytes(),
	}
}
