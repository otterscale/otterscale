package app

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/openhdc/openhdc/api/stack/v1"
	"github.com/openhdc/openhdc/internal/domain/service"
)

type StackApp struct {
	v1.UnimplementedStackServiceServer

	svc *service.StackService
}

func NewStackApp(svc *service.StackService) *StackApp {
	return &StackApp{
		svc: svc,
	}
}

var _ v1.StackServiceServer = (*StackApp)(nil)

func (a *StackApp) UpdateNTPServers(ctx context.Context, req *v1.UpdateNTPServersRequest) (*v1.UpdateNTPServersResponse, error) {
	return &v1.UpdateNTPServersResponse{}, nil
}

func (a *StackApp) ListPackageRepositories(ctx context.Context, req *v1.ListPackageRepositoriesRequest) (*v1.ListPackageRepositoriesResponse, error) {
	return &v1.ListPackageRepositoriesResponse{}, nil
}

func (a *StackApp) UpdatePackageRepositories(ctx context.Context, req *v1.UpdatePackageRepositoryRequest) (*v1.PackageRepository, error) {
	return &v1.PackageRepository{}, nil
}

func (a *StackApp) ListNetworks(ctx context.Context, req *v1.ListNetworksRequest) (*v1.ListNetworksResponse, error) {
	return &v1.ListNetworksResponse{}, nil
}

func (a *StackApp) CreateNetwork(ctx context.Context, req *v1.CreateNetworkRequest) (*v1.Network, error) {
	return &v1.Network{}, nil
}

func (a *StackApp) UpdateNetwork(ctx context.Context, req *v1.UpdateNetworkRequest) (*v1.Network, error) {
	return &v1.Network{}, nil
}

func (a *StackApp) DeleteNetwork(ctx context.Context, req *v1.DeleteNetworkRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (a *StackApp) UpdateDHCP(ctx context.Context, req *v1.UpdateDHCPRequest) (*v1.VLAN, error) {
	return &v1.VLAN{}, nil
}

func (a *StackApp) UpdateIPRange(ctx context.Context, req *v1.UpdateIPRangeRequest) (*v1.Subnet, error) {
	return &v1.Subnet{}, nil
}

func (a *StackApp) ImportBootResources(ctx context.Context, req *v1.ImportBootResourcesRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (a *StackApp) PowerOnMachine(ctx context.Context, req *v1.PowerOnMachineRequest) (*v1.Machine, error) {
	return &v1.Machine{}, nil
}

func (a *StackApp) PowerOffMachine(ctx context.Context, req *v1.PowerOffMachineRequest) (*v1.Machine, error) {
	return &v1.Machine{}, nil
}

func (a *StackApp) CommissionMachine(ctx context.Context, req *v1.CommissionMachineRequest) (*v1.Machine, error) {
	return &v1.Machine{}, nil
}

func (a *StackApp) OverrideMachineFailedTest(ctx context.Context, req *v1.OverrideMachineFailedTestRequest) (*v1.Machine, error) {
	return &v1.Machine{}, nil
}
