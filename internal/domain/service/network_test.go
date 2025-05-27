package service

import (
	"context"
	"errors"
	"net"
	"strconv"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/subnet"
	"github.com/openhdc/otterscale/internal/domain/model"
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestListNetworks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubnet := mocks.NewMockMAASSubnet(ctrl)
	mockFabric := mocks.NewMockMAASFabric(ctrl)
	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	s := &NexusService{subnet: mockSubnet, fabric: mockFabric, ipRange: mockIPRange}

	ctx := context.Background()

	// Case 1: Success
	returnedSubnets := []entity.Subnet{
		{ID: 1, CIDR: "192.168.1.0/24", VLAN: entity.VLAN{ID: 101}},
		{ID: 2, CIDR: "192.168.2.0/24", VLAN: entity.VLAN{ID: 102}},
	}

	returnedFabrics := []entity.Fabric{
		{ID: 10, VLANs: []entity.VLAN{{ID: 101, VID: 101}}},
		{ID: 20, VLANs: []entity.VLAN{{ID: 102, VID: 102}, {ID: 103, VID: 103}}},
	}

	returnedIPRanges := []entity.IPRange{} // Empty IPRanges

	mockSubnet.EXPECT().List(ctx).Return(returnedSubnets, nil)
	mockFabric.EXPECT().List(ctx).Return(returnedFabrics, nil)

	mockSubnet.EXPECT().GetStatistics(ctx, 1).Return(&subnet.Statistics{}, nil)
	mockSubnet.EXPECT().GetIPAddresses(ctx, 1).Return([]subnet.IPAddress{}, nil)
	mockIPRange.EXPECT().List(ctx).Return(returnedIPRanges, nil)

	mockSubnet.EXPECT().GetStatistics(ctx, 2).Return(&subnet.Statistics{}, nil)
	mockSubnet.EXPECT().GetIPAddresses(ctx, 2).Return([]subnet.IPAddress{}, nil)
	mockIPRange.EXPECT().List(ctx).Return(returnedIPRanges, nil)

	networks, err := s.ListNetworks(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(networks) != 3 {
		t.Errorf("Expected 3 networks, got %d", len(networks))
	}

	// Case 2: Error getting subnets
	mockSubnet.EXPECT().List(ctx).Return(nil, errors.New("subnet error"))
	_, err = s.ListNetworks(ctx)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Case 3: Error getting fabrics
	mockSubnet.EXPECT().List(ctx).Return([]entity.Subnet{}, nil)
	mockFabric.EXPECT().List(ctx).Return(nil, errors.New("fabric error"))
	_, err = s.ListNetworks(ctx)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

}

func TestCreateNetwork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFabric := mocks.NewMockMAASFabric(ctrl)
	mockSubnet := mocks.NewMockMAASSubnet(ctrl)
	mockVLAN := mocks.NewMockMAASVLAN(ctrl)
	mockIPRange := mocks.NewMockMAASIPRange(ctrl)

	s := &NexusService{fabric: mockFabric, subnet: mockSubnet, vlan: mockVLAN, ipRange: mockIPRange}
	ctx := context.Background()

	// Case 1: Success with DHCP off
	fabricParams := &model.FabricParams{}
	returnedFabric := &entity.Fabric{ID: 10, VLANs: []entity.VLAN{{ID: 101, VID: 101}}}
	subnetParams := &model.SubnetParams{CIDR: "192.168.1.0/24", GatewayIP: "192.168.1.1", DNSServers: []string{"8.8.8.8"}, Fabric: strconv.Itoa(returnedFabric.ID), VLAN: strconv.Itoa(returnedFabric.VLANs[0].ID)}

	mockFabric.EXPECT().Create(ctx, fabricParams).Return(returnedFabric, nil)

	returnedSubnet := &entity.Subnet{ID: 1, CIDR: "192.168.1.0/24", VLAN: returnedFabric.VLANs[0]}

	mockSubnet.EXPECT().Create(ctx, subnetParams).Return(returnedSubnet, nil)

	mockSubnet.EXPECT().GetStatistics(ctx, returnedSubnet.ID).Return(&subnet.Statistics{}, nil)
	mockSubnet.EXPECT().GetIPAddresses(ctx, returnedSubnet.ID).Return([]subnet.IPAddress{}, nil)
	mockIPRange.EXPECT().List(ctx).Return([]entity.IPRange{}, nil)

	_, err := s.CreateNetwork(ctx, "192.168.1.0/24", "192.168.1.1", []string{"8.8.8.8"}, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Case 2: Success with DHCP on
	mockFabric.EXPECT().Create(ctx, fabricParams).Return(returnedFabric, nil)
	mockSubnet.EXPECT().Create(ctx, subnetParams).Return(returnedSubnet, nil)
	mockVLAN.EXPECT().Update(ctx, returnedFabric.ID, returnedFabric.VLANs[0].VID, &model.VLANParams{}).Return(&returnedFabric.VLANs[0], nil)

	mockSubnet.EXPECT().GetStatistics(ctx, returnedSubnet.ID).Return(&subnet.Statistics{}, nil)
	mockSubnet.EXPECT().GetIPAddresses(ctx, returnedSubnet.ID).Return([]subnet.IPAddress{}, nil)
	mockIPRange.EXPECT().List(ctx).Return([]entity.IPRange{}, nil)

	_, err = s.CreateNetwork(ctx, "192.168.1.0/24", "192.168.1.1", []string{"8.8.8.8"}, true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Case 3: Error creating fabric
	mockFabric.EXPECT().Create(ctx, fabricParams).Return(nil, errors.New("fabric error"))
	_, err = s.CreateNetwork(ctx, "192.168.1.0/24", "192.168.1.1", []string{"8.8.8.8"}, false)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// ... Add more cases as needed for error handling in subnet.Create, vlan.Update etc.

}

func TestCreateIPRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	s := &NexusService{ipRange: mockIPRange}
	ctx := context.Background()

	params := &model.IPRangeParams{
		Type:    "reserved",
		Subnet:  "1",
		StartIP: "192.168.1.10",
		EndIP:   "192.168.1.20",
		Comment: "test range",
	}

	// Case 1: Success
	mockIPRange.EXPECT().Create(ctx, params).Return(&entity.IPRange{}, nil)
	_, err := s.CreateIPRange(ctx, 1, "192.168.1.10", "192.168.1.20", "test range")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Case 2: Error creating IP range
	mockIPRange.EXPECT().Create(ctx, params).Return(nil, errors.New("ip range error"))
	_, err = s.CreateIPRange(ctx, 1, "192.168.1.10", "192.168.1.20", "test range")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestDeleteNetwork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFabric := mocks.NewMockMAASFabric(ctrl)
	mockSubnet := mocks.NewMockMAASSubnet(ctrl)
	s := &NexusService{fabric: mockFabric, subnet: mockSubnet}
	ctx := context.Background()

	// Case 1: Success
	fabric := &entity.Fabric{ID: 10, VLANs: []entity.VLAN{{ID: 101}}}
	subnets := []entity.Subnet{{ID: 1, VLAN: fabric.VLANs[0]}}

	mockFabric.EXPECT().Get(ctx, 10).Return(fabric, nil)
	mockSubnet.EXPECT().List(ctx).Return(subnets, nil)
	mockSubnet.EXPECT().Delete(ctx, 1).Return(nil)
	mockFabric.EXPECT().Delete(ctx, 10).Return(nil)

	err := s.DeleteNetwork(ctx, 10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Case 2: Error getting fabric
	mockFabric.EXPECT().Get(ctx, 10).Return(nil, errors.New("fabric error"))
	err = s.DeleteNetwork(ctx, 10)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// ... Add more error cases as needed for subnet.List, subnet.Delete, and fabric.Delete
}

func TestDeleteIPRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	s := &NexusService{ipRange: mockIPRange}
	ctx := context.Background()

	// Case 1: Success
	mockIPRange.EXPECT().Delete(ctx, 1).Return(nil)
	err := s.DeleteIPRange(ctx, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Case 2: Error deleting IP range
	mockIPRange.EXPECT().Delete(ctx, 1).Return(errors.New("ip range error"))
	err = s.DeleteIPRange(ctx, 1)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGetNetworkSubnet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubnet := mocks.NewMockMAASSubnet(ctrl)
	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	s := &NexusService{subnet: mockSubnet, ipRange: mockIPRange}
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		subnetEntity := &model.Subnet{ID: 1, CIDR: "192.168.1.0/24"}
		mockSubnet.EXPECT().GetStatistics(ctx, subnetEntity.ID).Return(&subnet.Statistics{}, nil)
		mockSubnet.EXPECT().GetIPAddresses(ctx, subnetEntity.ID).Return([]subnet.IPAddress{}, nil)
		mockIPRange.EXPECT().List(ctx).Return([]entity.IPRange{}, nil)

		_, err := s.getNetworkSubnet(ctx, subnetEntity)
		if err != nil {
			t.Fatalf("getNetworkSubnet() unexpected error = %v", err)
		}
	})

	t.Run("Invalid CIDR", func(t *testing.T) {
		// Use a different ID for this test case to avoid conflicts if mocks were not perfectly isolated,
		// though with subtests and fresh expectations, it's less critical.
		invalidSubnetEntity := &model.Subnet{ID: 2, CIDR: "invalid"}

		// Expectations for calls that happen before ParseCIDR fails
		mockSubnet.EXPECT().GetStatistics(ctx, invalidSubnetEntity.ID).Return(&subnet.Statistics{}, nil)
		mockSubnet.EXPECT().GetIPAddresses(ctx, invalidSubnetEntity.ID).Return([]subnet.IPAddress{}, nil)
		mockIPRange.EXPECT().List(ctx).Return([]entity.IPRange{}, nil)

		_, err := s.getNetworkSubnet(ctx, invalidSubnetEntity)
		if err == nil {
			t.Errorf("Expected error for invalid CIDR, got nil")
		}
		// You can add more specific error checking here if needed, e.g.,
		// if !strings.Contains(err.Error(), "invalid CIDR address") {
		// 	t.Errorf("Expected CIDR parsing error, got: %v", err)
		// }
	})

	t.Run("GetStatistics error", func(t *testing.T) {
		subnetEntity := &model.Subnet{ID: 3, CIDR: "192.168.3.0/24"}
		expectedErr := errors.New("GetStatistics error")
		mockSubnet.EXPECT().GetStatistics(ctx, subnetEntity.ID).Return(nil, expectedErr)
		// GetIPAddresses and List should not be called after this error

		_, err := s.getNetworkSubnet(ctx, subnetEntity)
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("GetIPAddresses error", func(t *testing.T) {
		subnetEntity := &model.Subnet{ID: 4, CIDR: "192.168.4.0/24"}
		expectedErr := errors.New("GetIPAddresses error")
		mockSubnet.EXPECT().GetStatistics(ctx, subnetEntity.ID).Return(&subnet.Statistics{}, nil)
		mockSubnet.EXPECT().GetIPAddresses(ctx, subnetEntity.ID).Return(nil, expectedErr)
		// List should not be called after this error

		_, err := s.getNetworkSubnet(ctx, subnetEntity)
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("IPRange List error", func(t *testing.T) {
		subnetEntity := &model.Subnet{ID: 5, CIDR: "192.168.5.0/24"}
		expectedErr := errors.New("IPRange List error")
		mockSubnet.EXPECT().GetStatistics(ctx, subnetEntity.ID).Return(&subnet.Statistics{}, nil)
		mockSubnet.EXPECT().GetIPAddresses(ctx, subnetEntity.ID).Return([]subnet.IPAddress{}, nil)
		mockIPRange.EXPECT().List(ctx).Return(nil, expectedErr)

		_, err := s.getNetworkSubnet(ctx, subnetEntity)
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestUpdateFabric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFabric := mocks.NewMockMAASFabric(ctrl)
	s := &NexusService{fabric: mockFabric}
	ctx := context.Background()

	// Case 1: Success
	params := &model.FabricParams{Name: "new_fabric_name"}
	updatedFabric := &entity.Fabric{ID: 1, Name: "new_fabric_name"}
	mockFabric.EXPECT().Update(ctx, 1, params).Return(updatedFabric, nil)

	_, err := s.UpdateFabric(ctx, 1, "new_fabric_name")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Case 2: Error updating fabric
	mockFabric.EXPECT().Update(ctx, 1, params).Return(nil, errors.New("update fabric error"))
	_, err = s.UpdateFabric(ctx, 1, "new_fabric_name")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUpdateVLAN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVLAN := mocks.NewMockMAASVLAN(ctrl)
	s := &NexusService{vlan: mockVLAN}
	ctx := context.Background()

	params := &model.VLANParams{
		Name:        "new_vlan_name",
		MTU:         1500,
		Description: "new description",
		DHCPOn:      true,
	}
	updatedVLAN := &entity.VLAN{VID: 101, Name: "new_vlan_name", MTU: 1500, Description: "new description", DHCPOn: true}

	// Case 1: Success
	mockVLAN.EXPECT().Update(ctx, 1, 101, params).Return(updatedVLAN, nil)
	_, err := s.UpdateVLAN(ctx, 1, 101, "new_vlan_name", 1500, "new description", true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Case 2: Error updating VLAN
	mockVLAN.EXPECT().Update(ctx, 1, 101, params).Return(nil, errors.New("update vlan error"))
	_, err = s.UpdateVLAN(ctx, 1, 101, "new_vlan_name", 1500, "new description", true)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUpdateSubnet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubnet := mocks.NewMockMAASSubnet(ctrl)
	mockIPRange := mocks.NewMockMAASIPRange(ctrl)

	s := &NexusService{subnet: mockSubnet, ipRange: mockIPRange}
	ctx := context.Background()

	params := &model.SubnetParams{
		Name:        "new_subnet_name",
		CIDR:        "192.168.1.0/24",
		GatewayIP:   "192.168.1.1",
		DNSServers:  []string{"8.8.8.8"},
		Description: "new description",
		AllowDNS:    true,
	}
	updatedSubnet := &entity.Subnet{ID: 1, Name: "new_subnet_name", CIDR: "192.168.1.0/24"}

	// Case 1: Success
	mockSubnet.EXPECT().Update(ctx, 1, params).Return(updatedSubnet, nil)
	mockSubnet.EXPECT().GetStatistics(ctx, updatedSubnet.ID).Return(&subnet.Statistics{}, nil)
	mockSubnet.EXPECT().GetIPAddresses(ctx, updatedSubnet.ID).Return([]subnet.IPAddress{}, nil)
	mockIPRange.EXPECT().List(ctx).Return([]entity.IPRange{}, nil)

	_, err := s.UpdateSubnet(ctx, 1, "new_subnet_name", "192.168.1.0/24", "192.168.1.1", []string{"8.8.8.8"}, "new description", true)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Case 2: Error updating subnet
	mockSubnet.EXPECT().Update(ctx, 1, params).Return(nil, errors.New("update subnet error"))
	_, err = s.UpdateSubnet(ctx, 1, "new_subnet_name", "192.168.1.0/24", "192.168.1.1", []string{"8.8.8.8"}, "new description", true)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Case 3: Error getting Network Subnet after successful subnet update
	mockSubnet.EXPECT().Update(ctx, 1, params).Return(updatedSubnet, nil)
	mockSubnet.EXPECT().GetStatistics(ctx, updatedSubnet.ID).Return(nil, errors.New("error getting statistics")) // Example error

	_, err = s.UpdateSubnet(ctx, 1, "new_subnet_name", "192.168.1.0/24", "192.168.1.1", []string{"8.8.8.8"}, "new description", true)
	if err == nil {
		t.Errorf("Expected error in getNetworkSubnet, got nil")
	}

}

func TestUpdateIPRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	s := &NexusService{ipRange: mockIPRange}
	ctx := context.Background()

	// Case 1: Success
	params := &model.IPRangeParams{
		StartIP: "192.168.1.20",    // Updated start IP
		EndIP:   "192.168.1.30",    // Updated end IP
		Comment: "updated comment", // Include the comment here!
	}

	updatedIPRange := &entity.IPRange{ID: 1, StartIP: net.ParseIP(params.StartIP), EndIP: net.ParseIP(params.EndIP), Comment: params.Comment}
	mockIPRange.EXPECT().Update(ctx, 1, params).Return(updatedIPRange, nil) // Only one expectation now

	_, err := s.UpdateIPRange(ctx, 1, "192.168.1.20", "192.168.1.30", "updated comment")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Case 2: Error updating IP range
	params = &model.IPRangeParams{ // New params for this case, just like the successful one.
		StartIP: "192.168.1.20",
		EndIP:   "192.168.1.30",
		Comment: "updated comment", // Including the comment here.
	}

	mockIPRange.EXPECT().Update(ctx, 1, params).Return(nil, errors.New("update ip range error")) // Corrected params

	_, err = s.UpdateIPRange(ctx, 1, "192.168.1.20", "192.168.1.30", "updated comment")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
