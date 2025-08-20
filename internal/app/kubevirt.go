package app

import (
	"context"
    "strconv"

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
	vm, err := s.uc.CreateVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), toCoreMetadata(req.Msg.GetMetadata()), toCoreVirtualMachineSpec(req.Msg.GetSpec()))
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachine(vm)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) GetVirtualMachine(ctx context.Context, req *connect.Request[pb.GetVirtualMachineRequest]) (*connect.Response[pb.VirtualMachine], error) {
	vm, err := s.uc.GetVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachine(vm)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListVirtualMachines(ctx context.Context, req *connect.Request[pb.ListVirtualMachinesRequest]) (*connect.Response[pb.ListVirtualMachinesResponse], error) {
	vms, err := s.uc.ListVirtualMachines(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
	if err != nil {
		return nil, err
	}
	resp := &pb.ListVirtualMachinesResponse{}
	resp.SetVirtualMachines(toProtoVirtualMachines(vms))
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) UpdateVirtualMachine(ctx context.Context, req *connect.Request[pb.UpdateVirtualMachineRequest]) (*connect.Response[pb.VirtualMachine], error) {
	vm, err := s.uc.UpdateVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace(), req.Msg.GetInstanceTypeName(), req.Msg.GetNetworkName(), req.Msg.GetStartupScript(), req.Msg.GetLabels(), req.Msg.GetAnnotations(), req.Msg.GetDataVolumes(), toCoreDevices(req.Msg.GetDevices()))
	if err != nil {
		return nil, err
	}
	resp := toProtoVirtualMachine(vm)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) DeleteVirtualMachine(ctx context.Context, req *connect.Request[pb.DeleteVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.DeleteVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// Virtual Machine Control Operations
func (s *KubeVirtService) StartVirtualMachine(ctx context.Context, req *connect.Request[pb.StartVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.StartVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) StopVirtualMachine(ctx context.Context, req *connect.Request[pb.StopVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.StopVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) PauseVirtualMachine(ctx context.Context, req *connect.Request[pb.PauseVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.PauseVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) UnpauseVirtualMachine(ctx context.Context, req *connect.Request[pb.UnpauseVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.UnpauseVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// Virtual Machine Advanced Operations
func (s *KubeVirtService) CloneVirtualMachine(ctx context.Context, req *connect.Request[pb.CloneVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.CloneVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetTargetName(), req.Msg.GetTargetNamespace(), req.Msg.GetSourceName(), req.Msg.GetSourceNamespace()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) SnapshotVirtualMachine(ctx context.Context, req *connect.Request[pb.SnapshotVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.SnapshotVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace(), req.Msg.GetSnapshotName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) RestoreVirtualMachine(ctx context.Context, req *connect.Request[pb.RestoreVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.RestoreVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace(), req.Msg.GetSnapshotName()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) MigrateVirtualMachine(ctx context.Context, req *connect.Request[pb.MigrateVirtualMachineRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := s.uc.MigrateVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace(), req.Msg.GetTargetNode()); err != nil {
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
	dvs, err := s.uc.ListDataVolumes(ctx, req.Msg.GetScopeUuid(), req.Msg.GetNamespace(), req.Msg.GetFacilityName())
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
func toProtoVirtualMachines(vms []core.VirtualMachine) []*pb.VirtualMachine {
	ret := []*pb.VirtualMachine{}
	for i := range vms {
		ret = append(ret, toProtoVirtualMachine(&vms[i]))
	}
	return ret
}

func toProtoVirtualMachine(vm *core.VirtualMachine) *pb.VirtualMachine {
	ret := &pb.VirtualMachine{}
	/*
		ret.SetMetadata(toProtoMetadata(vm.Metadata))
		ret.SetSpec(toProtoVirtualMachineSpec(vm.Spec))
		ret.SetStatus(pb.VirtualMachine_Status(pb.VirtualMachine_Status_value[vm.Status]))
		ret.SetClones(toProtoOperations(vm.Clones))
		ret.SetSnapshots(toProtoOperations(vm.Snapshots))
		ret.SetMigrates(toProtoOperations(vm.Migrates))
		ret.SetRestores(toProtoOperations(vm.Restores))
	*/
	return ret
}

func toProtoMetadata(m core.Metadata) *pb.Metadata {
	ret := &pb.Metadata{}

	ret.SetName(m.Name)
	ret.SetNamespace(m.Namespace)
	ret.SetLabels(m.Labels)
	ret.SetAnnotations(m.Annotations)
	ret.SetCreatedAt(m.CreatedAt)
	ret.SetUpdatedAt(m.UpdatedAt)

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

func toProtoVirtualMachineSpec(spec core.KubeVirtVirtualMachineSpec) *pb.VirtualMachineSpec {
	ret := &pb.VirtualMachineSpec{}

	ret.SetInstanceTypeName(spec.InstanceTypeName)
	ret.SetNetworkName(spec.NetworkName)
	ret.SetStartupScript(spec.StartupScript)
	ret.SetDataVolumes(spec.DataVolumes)
	ret.SetDevices(toProtoDevices(spec.Devices))

	return ret
}

func toCoreVirtualMachineSpec(spec *pb.VirtualMachineSpec) core.VirtualMachineSpec {
	/*
		return core.KubeVirtVirtualMachineSpec{
			InstanceTypeName:    spec.GetInstanceTypeName(),
			NetworkName:   spec.GetNetworkName(),
			StartupScript: spec.GetStartupScript(),
			DataVolumes:   spec.GetDataVolumes(),
			Devices:       toCoreDevices(spec.GetDevices()),
		}
	*/
	return core.VirtualMachineSpec{}
}

func toProtoDevices(devices []core.Device) []*pb.Device {
	ret := []*pb.Device{}
	for _, d := range devices {
		dev := &pb.Device{}
		dev.SetName(d.Name)
		dev.SetType(d.Type)
		ret = append(ret, dev)
	}
	return ret
}

func toCoreDevices(devices []*pb.Device) []core.Device {
	ret := []core.Device{}
	for _, d := range devices {
		ret = append(ret, core.Device{
			Name: d.GetName(),
			Type: d.GetType(),
		})
	}
	return ret
}

func toProtoOperations(ops []core.Operation) []*pb.VirtualMachine_Operation {
	ret := []*pb.VirtualMachine_Operation{}
	for _, op := range ops {
		operation := &pb.VirtualMachine_Operation{}
		operation.SetName(op.Name)
		operation.SetType(op.Type)
		operation.SetDescription(op.Description)
		operation.SetCreatedAt(op.CreatedAt)
		operation.SetStatus(&pb.VirtualMachine_Operation_Result{})
		operation.GetStatus().SetStatus(pb.VirtualMachine_Operation_Result_Status(pb.VirtualMachine_Operation_Result_Status_value[op.Status.Status]))
		operation.GetStatus().SetMessage(op.Status.Message)
		operation.GetStatus().SetReason(op.Status.Reason)
		ret = append(ret, operation)
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
	ret.SetMetadata(toProtoMetadata(core.Metadata{
		Name:        dv.GetName(),
		Namespace:   dv.GetNamespace(),
		Labels:      dv.GetLabels(),
		Annotations: dv.GetAnnotations(),
		CreatedAt:   timestamppb.New(dv.CreationTimestamp.Time),
	}))
	source, sourceType, sizeBytes := core.ExtractDataVolumeInfo(dv)
	ret.SetType(sourceType)
	ret.SetSizeBytes(sizeBytes)
	ret.SetSource(source)
	return ret
}

func toCoreDataVolume(dv *pb.DataVolume) core.DataVolume {
	ret := core.DataVolume{}
	/*
		return core.KubeVirtDataVolume{
			Metadata:  toCoreMetadata(dv.GetMetadata()),
			Source:    dv.GetSource(),
			Type:      dv.GetType(),
			SizeBytes: dv.GetSizeBytes(),
		}*/
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
	ret.SetMetadata(toProtoMetadata(s.Metadata))

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

	ret.SetMetadata(toProtoMetadata(f.Metadata))
	ret.SetCpuCores(f.CpuCores)
	ret.SetMemoryBytes(f.MemoryBytes)

	return ret
}

func toCoreInstanceType(f *pb.InstanceType) core.InstanceType {
	return core.InstanceType{
		Metadata:    toCoreMetadata(f.GetMetadata()),
		CpuCores:    f.GetCpuCores(),
		MemoryBytes: f.GetMemoryBytes(),
	}
}
