package service

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/openhdc/otterscale/internal/domain/model"
	"go.uber.org/mock/gomock"

	// Make sure this import path matches where your generated mocks are located.
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
)

func TestNexusService_ListNetworks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubnet := mocks.NewMockMAASSubnet(ctrl)
	mockFabric := mocks.NewMockMAASFabric(ctrl)
	mockIPRange := mocks.NewMockMAASIPRange(ctrl)

	// Prepare test data
	subnets := []model.Subnet{
		{ID: 1, CIDR: "10.0.0.0/24", VLAN: model.VLAN{ID: 100}},
	}
	fabrics := []model.Fabric{
		{ID: 1, VLANs: []model.VLAN{{ID: 100}, {ID: 200}}},
	}
	// stats := &model.SubnetStatistics{}
	ipAddrs := []model.IPAddress{}
	ipRanges := []model.IPRange{
		{StartIP: net.ParseIP("10.0.0.10"), EndIP: net.ParseIP("10.0.0.20")},
	}

	mockSubnet.EXPECT().List(gomock.Any()).Return(subnets, nil)
	// mockSubnet.EXPECT().GetStatistics(gomock.Any(), 1).Return(stats, nil)
	mockSubnet.EXPECT().GetIPAddresses(gomock.Any(), 1).Return(ipAddrs, nil)
	mockIPRange.EXPECT().List(gomock.Any()).Return(ipRanges, nil)
	mockFabric.EXPECT().List(gomock.Any()).Return(fabrics, nil)

	svc := &NexusService{
		subnet:  mockSubnet,
		fabric:  mockFabric,
		ipRange: mockIPRange,
	}

	networks, err := svc.ListNetworks(context.Background())
	if err != nil {
		t.Fatalf("ListNetworks returned error: %v", err)
	}
	if len(networks) == 0 {
		t.Error("expected at least one network")
	}
}

func TestNexusService_CreateNetwork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFabric := mocks.NewMockMAASFabric(ctrl)
	mockSubnet := mocks.NewMockMAASSubnet(ctrl)
	mockVLAN := mocks.NewMockMAASVLAN(ctrl)

	fabric := &model.Fabric{ID: 1, VLANs: []model.VLAN{{ID: 2, VID: 3}}}
	subnet := &model.Subnet{ID: 10, CIDR: "10.0.0.0/24", VLAN: model.VLAN{ID: 2, VID: 3}}
	// stats := &model.SubnetStatistics{}
	ipAddrs := []model.IPAddress{}
	ipRanges := []model.IPRange{}

	mockFabric.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fabric, nil)
	mockSubnet.EXPECT().Create(gomock.Any(), gomock.Any()).Return(subnet, nil)
	mockVLAN.EXPECT().Update(gomock.Any(), fabric.ID, fabric.VLANs[0].VID, gomock.Any()).Return(&fabric.VLANs[0], nil)
	// mockSubnet.EXPECT().GetStatistics(gomock.Any(), subnet.ID).Return(stats, nil)
	mockSubnet.EXPECT().GetIPAddresses(gomock.Any(), subnet.ID).Return(ipAddrs, nil)
	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	mockIPRange.EXPECT().List(gomock.Any()).Return(ipRanges, nil)

	svc := &NexusService{
		fabric:  mockFabric,
		subnet:  mockSubnet,
		vlan:    mockVLAN,
		ipRange: mockIPRange,
	}

	network, err := svc.CreateNetwork(context.Background(), "10.0.0.0/24", "10.0.0.1", []string{"8.8.8.8"}, true)
	if err != nil {
		t.Fatalf("CreateNetwork returned error: %v", err)
	}
	if network == nil || network.Fabric.ID != 1 {
		t.Error("unexpected network result")
	}
}

func TestNexusService_DeleteNetwork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFabric := mocks.NewMockMAASFabric(ctrl)
	mockSubnet := mocks.NewMockMAASSubnet(ctrl)

	fabric := &model.Fabric{ID: 1, VLANs: []model.VLAN{{ID: 2}}}
	subnets := []model.Subnet{{ID: 10, VLAN: model.VLAN{ID: 2}}}

	mockFabric.EXPECT().Get(gomock.Any(), 1).Return(fabric, nil)
	mockSubnet.EXPECT().List(gomock.Any()).Return(subnets, nil)
	mockSubnet.EXPECT().Delete(gomock.Any(), 10).Return(nil)
	mockFabric.EXPECT().Delete(gomock.Any(), 1).Return(nil)

	svc := &NexusService{
		fabric: mockFabric,
		subnet: mockSubnet,
	}

	err := svc.DeleteNetwork(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeleteNetwork returned error: %v", err)
	}
}

func TestNexusService_CreateIPRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	ipRange := &model.IPRange{ID: 1}
	mockIPRange.EXPECT().Create(gomock.Any(), gomock.Any()).Return(ipRange, nil)

	svc := &NexusService{
		ipRange: mockIPRange,
	}

	result, err := svc.CreateIPRange(context.Background(), 1, "10.0.0.10", "10.0.0.20", "test")
	if err != nil {
		t.Fatalf("CreateIPRange returned error: %v", err)
	}
	if result.ID != 1 {
		t.Error("unexpected IPRange result")
	}
}

func TestNexusService_DeleteIPRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	mockIPRange.EXPECT().Delete(gomock.Any(), 1).Return(nil)

	svc := &NexusService{
		ipRange: mockIPRange,
	}

	err := svc.DeleteIPRange(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeleteIPRange returned error: %v", err)
	}
}

func TestNexusService_UpdateFabric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFabric := mocks.NewMockMAASFabric(ctrl)
	fabric := &model.Fabric{ID: 1}
	mockFabric.EXPECT().Update(gomock.Any(), 1, gomock.Any()).Return(fabric, nil)

	svc := &NexusService{
		fabric: mockFabric,
	}

	result, err := svc.UpdateFabric(context.Background(), 1, "new-fabric")
	if err != nil {
		t.Fatalf("UpdateFabric returned error: %v", err)
	}
	if result.ID != 1 {
		t.Error("unexpected Fabric result")
	}
}

func TestNexusService_UpdateVLAN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVLAN := mocks.NewMockMAASVLAN(ctrl)
	vlan := &model.VLAN{ID: 2}
	mockVLAN.EXPECT().Update(gomock.Any(), 1, 2, gomock.Any()).Return(vlan, nil)

	svc := &NexusService{
		vlan: mockVLAN,
	}

	result, err := svc.UpdateVLAN(context.Background(), 1, 2, "vlan", 1500, "desc", true)
	if err != nil {
		t.Fatalf("UpdateVLAN returned error: %v", err)
	}
	if result.ID != 2 {
		t.Error("unexpected VLAN result")
	}
}

func TestNexusService_UpdateSubnet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubnet := mocks.NewMockMAASSubnet(ctrl)
	subnet := &model.Subnet{ID: 10, CIDR: "10.0.0.0/24"}
	// stats := &model.SubnetStatistics{}
	ipAddrs := []model.IPAddress{}
	ipRanges := []model.IPRange{}

	mockSubnet.EXPECT().Update(gomock.Any(), 10, gomock.Any()).Return(subnet, nil)
	// mockSubnet.EXPECT().GetStatistics(gomock.Any(), 10).Return(stats, nil)
	mockSubnet.EXPECT().GetIPAddresses(gomock.Any(), 10).Return(ipAddrs, nil)
	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	mockIPRange.EXPECT().List(gomock.Any()).Return(ipRanges, nil)

	svc := &NexusService{
		subnet:  mockSubnet,
		ipRange: mockIPRange,
	}

	result, err := svc.UpdateSubnet(context.Background(), 10, "subnet", "10.0.0.0/24", "10.0.0.1", []string{"8.8.8.8"}, "desc", true)
	if err != nil {
		t.Fatalf("UpdateSubnet returned error: %v", err)
	}
	if result.Subnet.ID != 10 {
		t.Error("unexpected Subnet result")
	}
}

func TestNexusService_UpdateIPRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	ipRange := &model.IPRange{ID: 1}
	mockIPRange.EXPECT().Update(gomock.Any(), 1, gomock.Any()).Return(ipRange, nil)

	svc := &NexusService{
		ipRange: mockIPRange,
	}

	result, err := svc.UpdateIPRange(context.Background(), 1, "10.0.0.10", "10.0.0.20", "comment")
	if err != nil {
		t.Fatalf("UpdateIPRange returned error: %v", err)
	}
	if result.ID != 1 {
		t.Error("unexpected IPRange result")
	}
}

func TestNexusService_getNetworkSubnet_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubnet := mocks.NewMockMAASSubnet(ctrl)
	mockIPRange := mocks.NewMockMAASIPRange(ctrl)

	mockSubnet.EXPECT().GetStatistics(gomock.Any(), 10).Return(nil, errors.New("fail"))

	svc := &NexusService{
		subnet:  mockSubnet,
		ipRange: mockIPRange,
	}

	_, err := svc.getNetworkSubnet(context.Background(), &model.Subnet{ID: 10})
	if err == nil {
		t.Error("expected error from getNetworkSubnet")
	}
}
