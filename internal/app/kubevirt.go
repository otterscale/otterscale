package app

import (
	"context"

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
	vm, err := s.uc.UpdateVirtualMachine(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace(), req.Msg.GetFlavorName(), req.Msg.GetNetworkName(), req.Msg.GetStartupScript(), req.Msg.GetLabels(), req.Msg.GetAnnotations(), req.Msg.GetDataVolumes(), toCoreDevices(req.Msg.GetDevices()))
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
	dv, err := s.uc.CreateDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetDataVolume().GetMetadata().GetNamespace(), req.Msg.GetDataVolume().GetMetadata().GetName(), toCoreDataVolume(req.Msg.GetDataVolume()))
	if err != nil {
		return nil, err
	}
	resp := toProtoDataVolume(dv)
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) GetDataVolume(ctx context.Context, req *connect.Request[pb.GetDataVolumeRequest]) (*connect.Response[pb.DataVolume], error) {
	dv, err := s.uc.GetDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace())
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
	if err := s.uc.DeleteDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace()); err != nil {
		return nil, err
	}
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ExtendDataVolume(ctx context.Context, req *connect.Request[pb.ExtendDataVolumeRequest]) (*connect.Response[emptypb.Empty], error) {
	/*
		if err := s.uc.ExtendDataVolume(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace(), req.Msg.GetSizeBytes()); err != nil {
			return nil, err
		}
	*/
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// Network Operations
func (s *KubeVirtService) CreateNetwork(ctx context.Context, req *connect.Request[pb.CreateNetworkRequest]) (*connect.Response[pb.KubeVirtNetwork], error) {
	/*
		network, err := s.uc.CreateNetwork(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNetwork().GetMetadata().GetNamespace(), req.Msg.GetNetwork().GetMetadata().GetName(), toCoreNetwork(req.Msg.GetNetwork()))
		if err != nil {
			return nil, err
		}

		resp := toProtoNetwork(network)
	*/
	resp := &pb.KubeVirtNetwork{}

	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) GetNetwork(ctx context.Context, req *connect.Request[pb.GetNetworkRequest]) (*connect.Response[pb.KubeVirtNetwork], error) {
	/*
		network, err := s.uc.GetNetwork(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace())
		if err != nil {
			return nil, err
		}
		resp := toProtoNetwork(network)
	*/
	resp := &pb.KubeVirtNetwork{}

	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListNetworks(ctx context.Context, req *connect.Request[pb.ListNetworksRequest]) (*connect.Response[pb.ListNetworksResponse], error) {
	/*
		networks, err := s.uc.ListNetworks(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
		if err != nil {
			return nil, err
		}
		resp := &pb.ListNetworksResponse{}
		resp.SetNetworks(toProtoNetworks(networks))
	*/
	resp := &pb.ListNetworksResponse{}

	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) UpdateNetwork(ctx context.Context, req *connect.Request[pb.UpdateNetworkRequest]) (*connect.Response[pb.KubeVirtNetwork], error) {
	/*
		network, err := s.uc.UpdateNetwork(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), toCoreNetwork(req.Msg.GetNetwork()))
		if err != nil {
			return nil, err
		}
		resp := toProtoNetwork(network)
	*/
	resp := &pb.KubeVirtNetwork{}

	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) DeleteNetwork(ctx context.Context, req *connect.Request[pb.DeleteNetworkRequest]) (*connect.Response[emptypb.Empty], error) {
	/*
		if err := s.uc.DeleteNetwork(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace()); err != nil {
			return nil, err
		}
	*/
	resp := &emptypb.Empty{}
	return connect.NewResponse(resp), nil
}

// Flavor Operations
func (s *KubeVirtService) CreateFlavor(ctx context.Context, req *connect.Request[pb.CreateFlavorRequest]) (*connect.Response[pb.Flavor], error) {
	/*
		flavor, err := s.uc.CreateFlavor(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), toCoreFlavor(req.Msg.GetFlavor()))
		if err != nil {
			return nil, err
		}
		resp := toProtoFlavor(flavor)
	*/
	resp := &pb.Flavor{}
	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) GetFlavor(ctx context.Context, req *connect.Request[pb.GetFlavorRequest]) (*connect.Response[pb.Flavor], error) {
	/*
		flavor, err := s.uc.GetFlavor(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace())
		if err != nil {
			return nil, err
		}
		resp := toProtoFlavor(flavor)
	*/
	resp := &pb.Flavor{}

	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) ListFlavors(ctx context.Context, req *connect.Request[pb.ListFlavorsRequest]) (*connect.Response[pb.ListFlavorsResponse], error) {
	/*
		flavors, err := s.uc.ListFlavors(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetNamespace())
		if err != nil {
			return nil, err
		}
		resp := &pb.ListFlavorsResponse{}
		resp.Flavors = toProtoFlavors(flavors)
	*/
	resp := &pb.ListFlavorsResponse{}

	return connect.NewResponse(resp), nil
}

func (s *KubeVirtService) DeleteFlavor(ctx context.Context, req *connect.Request[pb.DeleteFlavorRequest]) (*connect.Response[emptypb.Empty], error) {
	/*
		if err := s.uc.DeleteFlavor(ctx, req.Msg.GetScopeUuid(), req.Msg.GetFacilityName(), req.Msg.GetName(), req.Msg.GetNamespace()); err != nil {
			return nil, err
		}
	*/
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

	ret.SetFlavorName(spec.FlavorName)
	ret.SetNetworkName(spec.NetworkName)
	ret.SetStartupScript(spec.StartupScript)
	ret.SetDataVolumes(spec.DataVolumes)
	ret.SetDevices(toProtoDevices(spec.Devices))

	return ret
}

func toCoreVirtualMachineSpec(spec *pb.VirtualMachineSpec) core.VirtualMachineSpec {
	/*
		return core.KubeVirtVirtualMachineSpec{
			FlavorName:    spec.GetFlavorName(),
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

func toProtoKubeVirtNetworks(networks []core.KubeVirtNetwork) []*pb.KubeVirtNetwork {
	ret := []*pb.KubeVirtNetwork{}
	for i := range networks {
		ret = append(ret, toProtoKubeVirtNetwork(&networks[i]))
	}
	return ret
}

func toProtoKubeVirtNetwork(n *core.KubeVirtNetwork) *pb.KubeVirtNetwork {
	ret := &pb.KubeVirtNetwork{}
	ret.SetMetadata(toProtoMetadata(n.Metadata))
	ret.SetServiceType(n.ServiceType)
	ret.SetPort(n.Port)
	ret.SetNodePort(n.NodePort)
	ret.SetContainerPort(n.ContainerPort)

	return ret
}

func toCoreNetwork(n *pb.KubeVirtNetwork) core.KubeVirtNetwork {
	return core.KubeVirtNetwork{
		Metadata:      toCoreMetadata(n.GetMetadata()),
		ServiceType:   n.GetServiceType(),
		Port:          n.GetPort(),
		NodePort:      n.GetNodePort(),
		ContainerPort: n.GetContainerPort(),
	}
}

func toProtoFlavors(flavors []core.Flavor) []*pb.Flavor {
	ret := []*pb.Flavor{}
	for i := range flavors {
		ret = append(ret, toProtoFlavor(&flavors[i]))
	}
	return ret
}

func toProtoFlavor(f *core.Flavor) *pb.Flavor {
	ret := &pb.Flavor{}

	ret.SetMetadata(toProtoMetadata(f.Metadata))
	ret.SetCpuCores(f.CpuCores)
	ret.SetMemoryBytes(f.MemoryBytes)

	return ret
}

func toCoreFlavor(f *pb.Flavor) core.Flavor {
	return core.Flavor{
		Metadata:    toCoreMetadata(f.GetMetadata()),
		CpuCores:    f.GetCpuCores(),
		MemoryBytes: f.GetMemoryBytes(),
	}
}
