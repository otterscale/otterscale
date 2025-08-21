package app

import (
	"context"
	"strconv"
	"strings"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/openhdc/otterscale/api/kubevirt/v1"
	"github.com/openhdc/otterscale/api/kubevirt/v1/pbconnect"
	"github.com/openhdc/otterscale/internal/core"
)

type KubeVirtService struct {
	pbconnect.UnimplementedKubeVirtServiceHandler

	uc *core.KubeVirtUseCase
}

func NewKubeVirtService(uc *core.KubeVirtUseCase) *KubeVirtService {
	return &KubeVirtService{uc: uc}
}

var _ pbconnect.KubeVirtServiceHandler = (*KubeVirtService)(nil)

// Virtual Machine Operations
func (s *KubeVirtService) CreateVirtualMachine(ctx context.Context, req *connect.Request[pb.CreateVirtualMachineRequest]) (*connect.Response[pb.VirtualMachine], error) {
	vm, err := s.uc.CreateVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetNetworkName(), req.Msg.GetStartupScript(), req.Msg.GetLabels(), req.Msg.GetAnnotations(), toCoreVirtualMachineResource(req.Msg.GetCustom(), req.Msg.GetInstancetype()), toCoreDiskDevices(req.Msg.GetDisks()))
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachine(vm, nil)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) GetVirtualMachine(ctx context.Context, req *connect.Request[pb.GetVirtualMachineRequest]) (*connect.Response[pb.VirtualMachine], error) {
	vm, vmi, err := s.uc.GetVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachine(vm, vmi)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListVirtualMachines(ctx context.Context, req *connect.Request[pb.ListVirtualMachinesRequest]) (*connect.Response[pb.ListVirtualMachinesResponse], error) {
	vms, vmis, err := s.uc.ListVirtualMachines(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListVirtualMachinesResponse{}
	resp.SetVirtualMachines(toProtoVirtualMachines(vms, vmis))
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) UpdateVirtualMachine(ctx context.Context, req *connect.Request[pb.UpdateVirtualMachineRequest]) (*connect.Response[pb.VirtualMachine], error) {
	vm, vmi, err := s.uc.UpdateVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetNetworkName(), req.Msg.GetStartupScript(), req.Msg.GetLabels(), req.Msg.GetAnnotations(), toCoreDiskDevices(req.Msg.GetDisks()))
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachine(vm, vmi)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) DeleteVirtualMachine(ctx context.Context, req *connect.Request[pb.DeleteVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
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
	if err := s.uc.CloneVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetTargetNamespace(), req.Msg.GetTargetName(), req.Msg.GetSourceNamespace(), req.Msg.GetSourceName(), req.Msg.GetDescription()); err != nil {
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
	if err := s.uc.RestoreVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetSnapshotName(), req.Msg.GetDescription()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) MigrateVirtualMachine(ctx context.Context, req *connect.Request[pb.MigrateVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.MigrateVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), req.Msg.GetTargetNode(), req.Msg.GetDescription()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// GetVirtualMachineClone retrieves a virtual machine clone
func (s *KubeVirtService) GetVirtualMachineClone(ctx context.Context, req *connect.Request[pb.GetVirtualMachineCloneRequest]) (*connect.Response[pb.VirtualMachineOperation], error) {
	op, err := s.uc.GetVirtualMachineClone(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := fromVirtualMachineClone(op)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListVirtualMachineClones(ctx context.Context, req *connect.Request[pb.ListVirtualMachineClonesRequest]) (*connect.Response[pb.ListVirtualMachineOperations], error) {
	resp := &pb.ListVirtualMachineOperations{}
	ops, err := s.uc.ListVirtualMachineClones(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp.SetOps(fromVirtualMachineClones(ops))
	return connect.NewResponse(resp), nil
}

// DeleteVirtualMachineClone deletes a virtual machine clone
func (s *KubeVirtService) DeleteVirtualMachineClone(ctx context.Context, req *connect.Request[pb.DeleteVirtualMachineCloneRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.uc.DeleteVirtualMachineClone(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// GetVirtualMachineSnapshot retrieves a virtual machine snapshot
func (s *KubeVirtService) GetVirtualMachineSnapshot(ctx context.Context, req *connect.Request[pb.GetVirtualMachineSnapshotRequest]) (*connect.Response[pb.VirtualMachineOperation], error) {
	op, err := s.uc.GetVirtualMachineSnapshot(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := fromVirtualMachineSnapshot(op)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListVirtualMachineSnapshots(ctx context.Context, req *connect.Request[pb.ListVirtualMachineSnapshotsRequest]) (*connect.Response[pb.ListVirtualMachineOperations], error) {
	resp := &pb.ListVirtualMachineOperations{}
	ops, err := s.uc.ListVirtualMachineSnapshots(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp.SetOps(fromVirtualMachineSnapshots(ops))
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

// GetVirtualMachineRestore retrieves a virtual machine restore operation
func (s *KubeVirtService) GetVirtualMachineRestore(ctx context.Context, req *connect.Request[pb.GetVirtualMachineRestoreRequest]) (*connect.Response[pb.VirtualMachineOperation], error) {
	op, err := s.uc.GetVirtualMachineRestore(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := fromVirtualMachineRestore(op)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListVirtualMachineRestores(ctx context.Context, req *connect.Request[pb.ListVirtualMachineRestoresRequest]) (*connect.Response[pb.ListVirtualMachineOperations], error) {
	resp := &pb.ListVirtualMachineOperations{}
	ops, err := s.uc.ListVirtualMachineRestores(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp.SetOps(fromVirtualMachineRestores(ops))
	return connect.NewResponse(resp), nil
}

// DeleteVirtualMachineRestore deletes a virtual machine restore operation
func (s *KubeVirtService) DeleteVirtualMachineRestore(ctx context.Context, req *connect.Request[pb.DeleteVirtualMachineRestoreRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.uc.DeleteVirtualMachineRestore(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// GetVirtualMachineMigrate retrieves a virtual machine Migrate operation
func (s *KubeVirtService) GetVirtualMachineMigrate(ctx context.Context, req *connect.Request[pb.GetVirtualMachineMigrateRequest]) (*connect.Response[pb.VirtualMachineOperation], error) {
	op, err := s.uc.GetVirtualMachineMigrate(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := fromVirtualMachineMigrate(op)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListVirtualMachineMigrates(ctx context.Context, req *connect.Request[pb.ListVirtualMachineMigratesRequest]) (*connect.Response[pb.ListVirtualMachineOperations], error) {
	resp := &pb.ListVirtualMachineOperations{}
	ops, err := s.uc.ListVirtualMachineMigrates(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp.SetOps(fromVirtualMachineMigrates(ops))
	return connect.NewResponse(resp), nil
}

// DeleteVirtualMachineMigrate deletes a virtual machine Migrate operation
func (s *KubeVirtService) DeleteVirtualMachineMigrate(ctx context.Context, req *connect.Request[pb.DeleteVirtualMachineMigrateRequest]) (*connect.Response[emptypb.Empty], error) {
	err := s.uc.DeleteVirtualMachineMigrate(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
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
		req.Msg.GetDataVolume().GetSizeBytes())
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

// VMService Operations
func (s *KubeVirtService) CreateVMService(ctx context.Context, req *connect.Request[pb.CreateVMServiceRequest]) (*connect.Response[pb.KubeVirtVMService], error) {
	vmservice, err := s.uc.CreateVMService(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetVMService().GetMetadata().GetNamespace(), req.Msg.GetVMService().GetMetadata().GetName(), toCoreVMService(req.Msg.GetVMService()))
	if err != nil {
		return nil, err
	}

	resp := toProtoKubeVirtVMService(vmservice)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) GetVMService(ctx context.Context, req *connect.Request[pb.GetVMServiceRequest]) (*connect.Response[pb.KubeVirtVMService], error) {
	vmservice, err := s.uc.GetVMService(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName())
	if err != nil {
		return nil, err
	}
	resp := toProtoKubeVirtVMService(vmservice)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListVMServices(ctx context.Context, req *connect.Request[pb.ListVMServicesRequest]) (*connect.Response[pb.ListVMServicesResponse], error) {
	vmservices, err := s.uc.ListVMServices(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListVMServicesResponse{}
	resp.SetVMServices(toProtoKubeVirtVMServices(vmservices))
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) UpdateVMService(ctx context.Context, req *connect.Request[pb.UpdateVMServiceRequest]) (*connect.Response[pb.KubeVirtVMService], error) {
	vmservice, err := s.uc.UpdateVMService(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName(), toCoreVMService(req.Msg.GetVMService()))
	if err != nil {
		return nil, err
	}
	resp := toProtoKubeVirtVMService(vmservice)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) DeleteVMService(ctx context.Context, req *connect.Request[pb.DeleteVMServiceRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteVMService(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace(), req.Msg.GetName()); err != nil {
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
func toProtoVirtualMachines(vms []core.VirtualMachine, vmis []core.VirtualMachineInstance) []*pb.VirtualMachine {
	ret := []*pb.VirtualMachine{}
	for i := 0; i < len(vmis); i++ {
		for j := 0; j < len(vms); j++ {
			if vms[j].Name == vmis[i].Name {
				ret = append(ret, toProtoVirtualMachine(&vms[j], &vmis[i]))
				break
			}
		}
	}
	return ret
}

func toProtoVirtualMachine(vm *core.VirtualMachine, vmi *core.VirtualMachineInstance) *pb.VirtualMachine {
	ret := &pb.VirtualMachine{}

	ret.SetMetadata(fromVirtualMachine(vm))
	ret.SetSpec(toProtoVirtualMachineSpec(vm))
	if vmi != nil {
		ret.SetStatus(toProtoVirtualMachineStatus(string(vmi.Status.Phase)))
	}
	// [TODO] snapshots
	// ret.SetSnapshots(toProtoOperations(vm.Snapshots))
	return ret
}

func fromVirtualMachine(vm *core.VirtualMachine) *pb.Metadata {
	return toProtoMetadata(vm.GetName(), vm.GetNamespace(), vm.GetLabels(), vm.GetAnnotations(), vm.CreationTimestamp.Time, vm.GetAnnotations()["otterscale/update-timestamp"])
}

func fromDataVolume(dv *core.DataVolume) *pb.Metadata {
	return toProtoMetadata(dv.GetName(), dv.GetNamespace(), dv.GetLabels(), dv.GetAnnotations(), dv.CreationTimestamp.Time, dv.GetAnnotations()["otterscale/update-timestamp"])
}

func fromVirtualMachineService(s *core.KubeVirtVMService) *pb.Metadata {
	return toProtoMetadata(s.Metadata.Name, s.Metadata.Namespace, s.Metadata.Labels, s.Metadata.Annotations, s.Metadata.CreatedAt.AsTime(), s.Metadata.Annotations["otterscale/update-timestamp"])
}

func fromInstanceType(t *core.InstanceType) *pb.Metadata {
	return toProtoMetadata(t.Metadata.Name, t.Metadata.Namespace, t.Metadata.Labels, t.Metadata.Annotations, t.Metadata.CreatedAt.AsTime(), t.Metadata.Annotations["otterscale/update-timestamp"])
}

func toProtoMetadata(name, namespace string, labels, annotations map[string]string, creationTimestamp time.Time, updateTimestamp string) *pb.Metadata {
	ret := &pb.Metadata{}

	ret.SetName(name)
	ret.SetNamespace(namespace)
	ret.SetLabels(labels)
	ret.SetAnnotations(annotations)
	ret.SetCreatedAt(timestamppb.New(creationTimestamp))
	parsedUpdateTime, err := time.Parse(time.RFC3339, updateTimestamp)
	if err == nil {
		ret.SetUpdatedAt(timestamppb.New(parsedUpdateTime))
	}

	return ret
}

func toCoreMetadata(m *pb.Metadata) core.Metadata {
	return core.Metadata{
		Name:        m.GetName(),
		Namespace:   m.GetNamespace(),
		Labels:      m.GetLabels(),
		Annotations: m.GetAnnotations(),
		CreatedAt:   m.GetCreatedAt(),
		UpdatedAt:   m.GetUpdatedAt(),
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

func toProtoVirtualMachineSpec(vm *core.VirtualMachine) *pb.VirtualMachineSpec {
	ret := &pb.VirtualMachineSpec{}

	ret.SetStartupScript(toProtoVirutalMachineScripts(vm.Spec.Template.Spec.Volumes))
	ret.SetResoureces(toProtoVirtualMachineResources(vm))
	// [TODO] if multiple network in the future
	ret.SetNetworkName(vm.Spec.Template.Spec.Networks[0].Name)

	return ret
}

func toProtoVirtualMachineStatus(s string) pb.VirtualMachine_Status {
	v, ok := pb.VirtualMachine_Status_value[strings.ToUpper(s)]
	if ok {
		return pb.VirtualMachine_Status(v)
	}
	return pb.VirtualMachine_RUNNING
}

func toProtoVirutalMachineScripts(volume []core.KubeVirtVolume) string {
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

func toCoreDiskDevices(disks []*pb.VirtualMachineDisk) []core.DiskDevice {
	ret := []core.DiskDevice{}
	for i := range disks {
		ret = append(ret, core.DiskDevice{
			Name:     disks[i].GetName(),
			DiskType: pb.VirtualMachineDiskType_name[int32(disks[i].GetDiskType())],
			Bus:      disks[i].GetBus(),
			Data:     disks[i].GetData(),
		})
	}
	return ret
}

func fromVirtualMachineClones(clones []core.VirtualMachineClone) []*pb.VirtualMachineOperation {
	ret := []*pb.VirtualMachineOperation{}
	for i := range clones {
		ret = append(ret, toProtoVirtualMachineOperation(pb.VirtualMachineOperation_CLONE, clones[i].Namespace, clones[i].Name, clones[i].Namespace, clones[i].Spec.Source.Name, clones[i].GetAnnotations()["otterscale.io/clone-description"], clones[i].CreationTimestamp.Time))
	}
	return ret
}

func fromVirtualMachineClone(clone *core.VirtualMachineClone) *pb.VirtualMachineOperation {
	return toProtoVirtualMachineOperation(pb.VirtualMachineOperation_CLONE, clone.Namespace, clone.Name, clone.Namespace, clone.Spec.Source.Name, clone.GetAnnotations()["otterscale.io/clone-description"], clone.CreationTimestamp.Time)
}

func fromVirtualMachineSnapshots(snapshots []core.VirtualMachineSnapshot) []*pb.VirtualMachineOperation {
	ret := []*pb.VirtualMachineOperation{}
	for i := range snapshots {
		ret = append(ret, toProtoVirtualMachineOperation(pb.VirtualMachineOperation_SNAPSHOT, snapshots[i].Namespace, snapshots[i].Name, snapshots[i].Namespace, snapshots[i].Spec.Source.Name, snapshots[i].GetAnnotations()["otterscale.io/snapshot-description"], snapshots[i].CreationTimestamp.Time))
	}
	return ret
}

func fromVirtualMachineSnapshot(snapshot *core.VirtualMachineSnapshot) *pb.VirtualMachineOperation {
	return toProtoVirtualMachineOperation(pb.VirtualMachineOperation_SNAPSHOT, snapshot.Namespace, snapshot.Name, snapshot.Namespace, snapshot.Spec.Source.Name, snapshot.GetAnnotations()["otterscale.io/snapshot-description"], snapshot.CreationTimestamp.Time)
}

func fromVirtualMachineRestores(restores []core.VirtualMachineRestore) []*pb.VirtualMachineOperation {
	ret := []*pb.VirtualMachineOperation{}
	for i := range restores {
		ret = append(ret, toProtoVirtualMachineOperation(pb.VirtualMachineOperation_RESTORE, restores[i].Namespace, restores[i].Spec.Target.Name, restores[i].Namespace, restores[i].Spec.VirtualMachineSnapshotName, restores[i].GetAnnotations()["otterscale.io/restore-description"], restores[i].CreationTimestamp.Time))
	}
	return ret
}

func fromVirtualMachineRestore(restore *core.VirtualMachineRestore) *pb.VirtualMachineOperation {
	return toProtoVirtualMachineOperation(pb.VirtualMachineOperation_RESTORE, restore.Namespace, restore.Spec.Target.Name, restore.Namespace, restore.Spec.VirtualMachineSnapshotName, restore.GetAnnotations()["otterscale.io/restore-description"], restore.CreationTimestamp.Time)
}

func fromVirtualMachineMigrates(migrates []core.VirtualMachineInstanceMigration) []*pb.VirtualMachineOperation {
	ret := []*pb.VirtualMachineOperation{}
	for i := range migrates {
		ret = append(ret, toProtoVirtualMachineOperation(pb.VirtualMachineOperation_MIGRATE, migrates[i].Namespace, migrates[i].Name, migrates[i].Namespace, migrates[i].Spec.VMIName, migrates[i].GetAnnotations()["otterscale.io/migrate-description"], migrates[i].CreationTimestamp.Time))
	}
	return ret
}

func fromVirtualMachineMigrate(migrate *core.VirtualMachineInstanceMigration) *pb.VirtualMachineOperation {
	return toProtoVirtualMachineOperation(pb.VirtualMachineOperation_MIGRATE, migrate.Namespace, migrate.Name, migrate.Namespace, migrate.Spec.VMIName, migrate.GetAnnotations()["otterscale.io/migrate-description"], migrate.CreationTimestamp.Time)
}

// toProtoVirtualMachineOperation converts core VirtualMachineOperation to protobuf
func toProtoVirtualMachineOperation(t pb.VirtualMachineOperation_Type, namespace, name, srcNamespace, srcName, description string, creationTimestamp time.Time) *pb.VirtualMachineOperation {
	ret := &pb.VirtualMachineOperation{}
	ret.SetType(t)
	ret.SetName(name)
	ret.SetNamespace(namespace)
	ret.SetSourceName(srcName)
	ret.SetDescription(description)
	ret.SetSourceNamespace(srcNamespace)
	ret.SetCreatedAt(timestamppb.New(creationTimestamp))

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
	ret.SetMetadata(fromDataVolume(dv))
	source, sourceType, sizeBytes := core.ExtractDataVolumeInfo(dv)
	ret.SetType(sourceType)
	ret.SetSizeBytes(sizeBytes)
	ret.SetSource(source)
	return ret
}

func toProtoKubeVirtVMServices(vmservices []core.KubeVirtVMService) []*pb.KubeVirtVMService {
	ret := make([]*pb.KubeVirtVMService, 0, len(vmservices))
	for i := range vmservices {
		ret = append(ret, toProtoKubeVirtVMService(&vmservices[i]))
	}
	return ret
}

func toProtoKubeVirtVMService(s *core.KubeVirtVMService) *pb.KubeVirtVMService {
	ret := &pb.KubeVirtVMService{}
	ret.SetMetadata(fromVirtualMachineService(s))

	spec := &pb.KubeVirtVMServiceSpec{}
	if s.Selector != nil {
		if v, ok := s.Selector["kubevirt.io/vm"]; ok && v != "" {
			spec.SetVMName(v)
		}
	}
	spec.SetSelector(s.Selector)

	spec.SetType(pb.KubeVirtVMServiceSpec_Type(pb.KubeVirtVMServiceSpec_Type_value["CLUSTER_IP"]))
	ports := make([]*pb.ServicePort, 0, len(s.Ports))
	for _, p := range s.Ports {
		sp := &pb.ServicePort{}
		sp.SetName(p.Name)
		sp.SetPort(p.Port)
		sp.SetTargetPort(strconv.Itoa(int(p.TargetPort)))
		sp.SetProtocol(pb.KubeVirtVMServiceSpec_Protocol(pb.KubeVirtVMServiceSpec_Protocol_value["TCP"]))
		sp.SetNodePort(p.NodePort)
		ports = append(ports, sp)
	}
	spec.SetPorts(ports)

	ret.SetSpec(spec)

	ret.SetStatus(&pb.KubeVirtVMServiceStatus{})

	return ret
}

func toCoreVMService(n *pb.KubeVirtVMService) core.KubeVirtVMService {
	ret := core.KubeVirtVMService{
		Metadata: toCoreMetadata(n.GetMetadata()),
	}

	if n.GetSpec() != nil {
		selector := map[string]string{}
		if vm := n.GetSpec().GetVMName(); vm != "" {
			selector["kubevirt.io/vm"] = vm
		}
		for k, v := range n.GetSpec().GetSelector() {
			selector[k] = v
		}
		ret.Selector = selector

		ports := make([]core.KubeVirtVMServicePort, 0, len(n.GetSpec().GetPorts()))
		for _, p := range n.GetSpec().GetPorts() {
			tp := int32(0)
			if s := p.GetTargetPort(); s != "" {
				if v, err := strconv.Atoi(s); err == nil {
					tp = int32(v)
				}
			}
			ports = append(ports, core.KubeVirtVMServicePort{
				Name:       p.GetName(),
				Port:       p.GetPort(),
				NodePort:   p.GetNodePort(),
				TargetPort: tp,
			})
		}
		ret.Ports = ports
	}
	return ret
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

func toCoreInstanceType(f *pb.InstanceType) core.InstanceType {
	return core.InstanceType{
		Metadata:    toCoreMetadata(f.GetMetadata()),
		CPUCores:    f.GetCpuCores(),
		MemoryBytes: f.GetMemoryBytes(),
	}
}
