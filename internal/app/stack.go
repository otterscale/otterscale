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
	n, err := a.svc.CreateNetwork(ctx, &model.FabricParams{}, &model.SubnetParams{}, &model.IPRangeParams{})
	if err != nil {
		return nil, err
	}
	return toNetwork(n), nil
}

func (a *StackApp) DeleteNetwork(ctx context.Context, req *v1.DeleteNetworkRequest) (*emptypb.Empty, error) {
	if err := a.svc.DeleteNetwork(ctx, int(req.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Network configuration functions
func (a *StackApp) UpdateFabric(ctx context.Context, req *v1.UpdateFabricRequest) (*v1.Fabric, error) {
	f, err := a.svc.UpdateFabric(ctx, int(req.GetId()), &model.FabricParams{})
	if err != nil {
		return nil, err
	}
	return toFabric(f), nil
}

func (a *StackApp) UpdateVLAN(ctx context.Context, req *v1.UpdateVLANRequest) (*v1.VLAN, error) {
	v, err := a.svc.UpdateVLAN(ctx, int(req.GetFabricId()), int(req.GetVid()), &model.VLANParams{})
	if err != nil {
		return nil, err
	}
	return toVLAN(v), nil
}

func (a *StackApp) UpdateSubnet(ctx context.Context, req *v1.UpdateSubnetRequest) (*v1.Subnet, error) {
	s, err := a.svc.UpdateSubnet(ctx, int(req.GetId()), &model.SubnetParams{})
	if err != nil {
		return nil, err
	}
	return toSubnet(s), nil
}

func (a *StackApp) UpdateIPRange(ctx context.Context, req *v1.UpdateIPRangeRequest) (*v1.IPRange, error) {
	r, err := a.svc.UpdateIPRange(ctx, int(req.GetId()), &model.IPRangeParams{})
	if err != nil {
		return nil, err
	}
	return toIPRange(r), nil
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
	m, err := a.svc.PowerOnMachine(ctx, req.GetSystemId(), &model.MachinePowerOnParams{})
	if err != nil {
		return nil, err
	}
	return toMachine(m), nil
}

func (a *StackApp) PowerOffMachine(ctx context.Context, req *v1.PowerOffMachineRequest) (*v1.Machine, error) {
	m, err := a.svc.PowerOffMachine(ctx, req.GetSystemId(), &model.MachinePowerOffParams{})
	if err != nil {
		return nil, err
	}
	return toMachine(m), nil
}

// Machine provisioning functions
func (a *StackApp) CommissionMachine(ctx context.Context, req *v1.CommissionMachineRequest) (*v1.Machine, error) {
	m, err := a.svc.CommissionMachine(ctx, req.GetSystemId(), &model.MachineCommissionParams{})
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
	return &v1.PackageRepository{}
}

func toNetworks(ns []*model.Network) []*v1.Network {
	ret := make([]*v1.Network, len(ns))
	for i := range ns {
		ret[i] = toNetwork(ns[i])
	}
	return ret
}

func toNetwork(n *model.Network) *v1.Network {
	return &v1.Network{}
}

func toFabric(f *model.Fabric) *v1.Fabric {
	return &v1.Fabric{}
}

func toVLAN(v *model.VLAN) *v1.VLAN {
	return &v1.VLAN{}
}

func toSubnet(s *model.Subnet) *v1.Subnet {
	return &v1.Subnet{}
}

func toIPRange(r *model.IPRange) *v1.IPRange {
	return &v1.IPRange{}
}

func toMachine(m *model.Machine) *v1.Machine {
	return &v1.Machine{}
}
