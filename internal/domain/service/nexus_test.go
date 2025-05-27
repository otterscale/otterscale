package service

import (
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
)

func TestNewNexusService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mocks.NewMockMAASServer(ctrl)
	mockPackageRepository := mocks.NewMockMAASPackageRepository(ctrl)
	mockBootResource := mocks.NewMockMAASBootResource(ctrl)
	mockBootSource := mocks.NewMockMAASBootSource(ctrl)
	mockBootSourceSelection := mocks.NewMockMAASBootSourceSelection(ctrl)
	mockFabric := mocks.NewMockMAASFabric(ctrl)
	mockVLAN := mocks.NewMockMAASVLAN(ctrl)
	mockSubnet := mocks.NewMockMAASSubnet(ctrl)
	mockIPRange := mocks.NewMockMAASIPRange(ctrl)
	mockMachine := mocks.NewMockMAASMachine(ctrl)
	mockTag := mocks.NewMockMAASTag(ctrl)
	mockSSHKey := mocks.NewMockMAASSSHKey(ctrl)
	mockKeyManager := mocks.NewMockJujuKey(ctrl)
	mockClient := mocks.NewMockJujuClient(ctrl)
	mockMachineManager := mocks.NewMockJujuMachine(ctrl)
	mockScope := mocks.NewMockJujuModel(ctrl)
	mockScopeConfig := mocks.NewMockJujuModelConfig(ctrl)
	mockFacility := mocks.NewMockJujuApplication(ctrl)
	mockAction := mocks.NewMockJujuAction(ctrl)
	mockCharmhub := mocks.NewMockJujuCharmHub(ctrl)
	mockKubernetes := mocks.NewMockKubeClient(ctrl)
	mockApps := mocks.NewMockKubeApps(ctrl)
	mockBatch := mocks.NewMockKubeBatch(ctrl)
	mockCore := mocks.NewMockKubeCore(ctrl)
	mockStorage := mocks.NewMockKubeStorage(ctrl)
	mockHelm := mocks.NewMockKubeHelm(ctrl)

	service := NewNexusService(
		mockServer,
		mockPackageRepository,
		mockBootResource,
		mockBootSource,
		mockBootSourceSelection,
		mockFabric,
		mockVLAN,
		mockSubnet,
		mockIPRange,
		mockMachine,
		mockTag,
		mockSSHKey,
		mockKeyManager,
		mockClient,
		mockMachineManager,
		mockScope,
		mockScopeConfig,
		mockFacility,
		mockAction,
		mockCharmhub,
		mockKubernetes,
		mockApps,
		mockBatch,
		mockCore,
		mockStorage,
		mockHelm,
	)

	if service == nil {
		t.Fatal("NewNexusService returned nil")
	}

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
	}{
		{"server", mockServer, service.server},
		{"packageRepository", mockPackageRepository, service.packageRepository},
		{"bootResource", mockBootResource, service.bootResource},
		{"bootSource", mockBootSource, service.bootSource},
		{"bootSourceSelection", mockBootSourceSelection, service.bootSourceSelection},
		{"fabric", mockFabric, service.fabric},
		{"vlan", mockVLAN, service.vlan},
		{"subnet", mockSubnet, service.subnet},
		{"ipRange", mockIPRange, service.ipRange},
		{"machine", mockMachine, service.machine},
		{"tag", mockTag, service.tag},
		{"sshKey", mockSSHKey, service.sshKey},
		{"keyManager", mockKeyManager, service.keyManager},
		{"client", mockClient, service.client},
		{"machineManager", mockMachineManager, service.machineManager},
		{"scope", mockScope, service.scope},
		{"scopeConfig", mockScopeConfig, service.scopeConfig},
		{"facility", mockFacility, service.facility},
		{"action", mockAction, service.action},
		{"charmhub", mockCharmhub, service.charmhub},
		{"kubernetes", mockKubernetes, service.kubernetes},
		{"apps", mockApps, service.apps},
		{"batch", mockBatch, service.batch},
		{"core", mockCore, service.core},
		{"storage", mockStorage, service.storage},
		{"helm", mockHelm, service.helm},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Using reflect.DeepEqual for interfaces might not be ideal if we want to check pointer equality.
			// For mocks, checking if they are the same instance is more direct.
			if tt.actual != tt.expected {
				t.Errorf("field %s: expected %v, got %v", tt.name, reflect.TypeOf(tt.expected), reflect.TypeOf(tt.actual))
			}
		})
	}
}
