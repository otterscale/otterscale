package service

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/juju/juju/api/client/action"

	// "github.com/juju/juju/api/client/application"
	// corebase "github.com/juju/juju/core/base"

	// "github.com/juju/juju/core/constraints"
	// "github.com/juju/juju/core/instance"
	"github.com/juju/juju/rpc/params"
	"go.uber.org/mock/gomock"
	jujuyaml "gopkg.in/yaml.v2"

	"github.com/openhdc/otterscale/internal/domain/model"
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	// Import for base.SeriesUbuntu to be recognized by mockgen if it were used directly in interfaces,
	// but for tests, we usually mock the direct dependencies.
	// _ "github.com/juju/juju/core/base"
	// _ "github.com/canonical/gomaasclient/entity" // For entity types if needed directly
)

func TestNexusService_ListFacilities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockJujuClient(ctrl)
	s := &NexusService{client: mockClient}

	ctx := context.Background()
	uuid := "test-uuid"

	t.Run("success", func(t *testing.T) {
		expectedStatus := &params.FullStatus{
			Applications: map[string]params.ApplicationStatus{
				"facility1": {Charm: "charm1", CharmProfile: "focal"},
				"facility2": {Charm: "charm2", CharmProfile: "bionic"},
			},
		}
		mockClient.EXPECT().Status(ctx, uuid, []string{"application", "*"}).Return(expectedStatus, nil)

		facilities, err := s.ListFacilities(ctx, uuid)
		if err != nil {
			t.Fatalf("ListFacilities() error = %v, wantErr nil", err)
		}
		if len(facilities) != 2 {
			t.Fatalf("ListFacilities() len = %d, want 2", len(facilities))
		}
		// Add more specific checks for facility content if needed
	})

	t.Run("client_status_error", func(t *testing.T) {
		mockClient.EXPECT().Status(ctx, uuid, []string{"application", "*"}).Return(nil, errExpected)
		_, err := s.ListFacilities(ctx, uuid)
		if !errors.Is(err, errExpected) {
			t.Fatalf("ListFacilities() error = %v, want errExpected %v", err, errExpected)
		}
	})
}

func TestNexusService_GetFacility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockJujuClient(ctrl)
	mockJujuApp := mocks.NewMockJujuApplication(ctrl)
	s := NewNexusService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockClient, nil, nil, nil, mockJujuApp, nil, nil, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()
	uuid := "test-uuid"
	facilityName := "test-facility"
	appStatus := params.ApplicationStatus{Charm: "test-charm"}
	configData := map[string]any{"key": "value"}
	configYAMLBytes, _ := jujuyaml.Marshal(configData)

	t.Run("success", func(t *testing.T) {
		expectedStatus := &params.FullStatus{
			Applications: map[string]params.ApplicationStatus{
				facilityName: appStatus,
			},
		}
		mockClient.EXPECT().Status(ctx, uuid, []string{"application", facilityName}).Return(expectedStatus, nil)
		mockJujuApp.EXPECT().GetConfig(ctx, uuid, facilityName).Return(configData, nil)

		facility, err := s.GetFacility(ctx, uuid, facilityName)
		if err != nil {
			t.Fatalf("GetFacility() unexpected error: %v", err)
		}
		if facility.Name != facilityName {
			t.Errorf("expected facility name %s, got %s", facilityName, facility.Name)
		}
		if !reflect.DeepEqual(facility.Status, &appStatus) {
			t.Errorf("expected status %+v, got %+v", &appStatus, facility.Status)
		}
		if facility.FacilityMetadata.ConfigYAML != string(configYAMLBytes) {
			t.Errorf("expected config YAML %s, got %s", string(configYAMLBytes), facility.FacilityMetadata.ConfigYAML)
		}
	})

	t.Run("error_client_status", func(t *testing.T) {
		mockClient.EXPECT().Status(ctx, uuid, []string{"application", facilityName}).Return(nil, errExpected)
		_, err := s.GetFacility(ctx, uuid, facilityName)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("error_get_config", func(t *testing.T) {
		expectedStatus := &params.FullStatus{
			Applications: map[string]params.ApplicationStatus{
				facilityName: appStatus,
			},
		}
		mockClient.EXPECT().Status(ctx, uuid, []string{"application", facilityName}).Return(expectedStatus, nil)
		mockJujuApp.EXPECT().GetConfig(ctx, uuid, facilityName).Return(nil, errExpected)
		_, err := s.GetFacility(ctx, uuid, facilityName)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		expectedStatus := &params.FullStatus{
			Applications: map[string]params.ApplicationStatus{
				"other-facility": appStatus, // Facility not found
			},
		}
		mockClient.EXPECT().Status(ctx, uuid, []string{"application", facilityName}).Return(expectedStatus, nil)
		_, err := s.GetFacility(ctx, uuid, facilityName)
		if err == nil {
			t.Fatal("expected not found error, got nil")
		}
		// Check for gRPC status code NotFound, though this is a domain service error
		// For simplicity, checking string containment.
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected not found error, got %v", err)
		}
	})
}

func TestNexusService_CreateFacility(t *testing.T) {

}

func TestNexusService_UpdateFacility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJujuApp := mocks.NewMockJujuApplication(ctrl)
	mockClient := mocks.NewMockJujuClient(ctrl) // For the GetFacility call
	s := NewNexusService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockClient, nil, nil, nil, mockJujuApp, nil, nil, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()
	uuid := "test-uuid"
	name := "test-facility"
	configYAML := "new_key: new_value"

	t.Run("success", func(t *testing.T) {
		mockJujuApp.EXPECT().Update(ctx, uuid, name, configYAML).Return(nil)
		// Mocks for the subsequent GetFacility call
		appStatus := params.ApplicationStatus{Charm: "updated-charm"}
		configData := map[string]any{"new_key": "new_value"}
		expectedStatus := &params.FullStatus{Applications: map[string]params.ApplicationStatus{name: appStatus}}
		mockClient.EXPECT().Status(ctx, uuid, []string{"application", name}).Return(expectedStatus, nil)
		mockJujuApp.EXPECT().GetConfig(ctx, uuid, name).Return(configData, nil)

		facility, err := s.UpdateFacility(ctx, uuid, name, configYAML)
		if err != nil {
			t.Fatalf("UpdateFacility() unexpected error: %v", err)
		}
		if facility.Name != name {
			t.Errorf("expected facility name %s, got %s", name, facility.Name)
		}
	})

	t.Run("error_juju_update", func(t *testing.T) {
		mockJujuApp.EXPECT().Update(ctx, uuid, name, configYAML).Return(errExpected)
		_, err := s.UpdateFacility(ctx, uuid, name, configYAML)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})

	t.Run("error_get_facility_after_update", func(t *testing.T) {
		mockJujuApp.EXPECT().Update(ctx, uuid, name, configYAML).Return(nil)
		mockClient.EXPECT().Status(ctx, uuid, []string{"application", name}).Return(nil, errExpected) // GetFacility fails
		_, err := s.UpdateFacility(ctx, uuid, name, configYAML)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})
}

func TestNexusService_DeleteFacility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJujuApp := mocks.NewMockJujuApplication(ctrl)
	s := NewNexusService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockJujuApp, nil, nil, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()
	uuid := "test-uuid"
	name := "test-facility"
	destroyStorage := true
	force := false

	t.Run("success", func(t *testing.T) {
		mockJujuApp.EXPECT().Delete(ctx, uuid, name, destroyStorage, force).Return(nil)
		err := s.DeleteFacility(ctx, uuid, name, destroyStorage, force)
		if err != nil {
			t.Fatalf("DeleteFacility() unexpected error: %v", err)
		}
	})

	t.Run("error_juju_delete", func(t *testing.T) {
		mockJujuApp.EXPECT().Delete(ctx, uuid, name, destroyStorage, force).Return(errExpected)
		err := s.DeleteFacility(ctx, uuid, name, destroyStorage, force)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})
}

func TestNexusService_ExposeFacility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJujuApp := mocks.NewMockJujuApplication(ctrl)
	s := NewNexusService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockJujuApp, nil, nil, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()
	uuid := "test-uuid"
	name := "test-facility"

	t.Run("success", func(t *testing.T) {
		mockJujuApp.EXPECT().Expose(ctx, uuid, name, nil).Return(nil)
		err := s.ExposeFacility(ctx, uuid, name)
		if err != nil {
			t.Fatalf("ExposeFacility() unexpected error: %v", err)
		}
	})

	t.Run("error_juju_expose", func(t *testing.T) {
		mockJujuApp.EXPECT().Expose(ctx, uuid, name, nil).Return(errExpected)
		err := s.ExposeFacility(ctx, uuid, name)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})
}

func TestNexusService_AddFacilityUnits(t *testing.T) {

}

func TestNexusService_ListActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAction := mocks.NewMockJujuAction(ctrl)
	s := NewNexusService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockAction, nil, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()
	uuid := "test-uuid"
	appName := "test-app"
	actionSpec := action.ActionSpec{Description: "test action"}

	t.Run("success", func(t *testing.T) {
		expectedActions := map[string]action.ActionSpec{"action1": actionSpec}
		mockAction.EXPECT().List(ctx, uuid, appName).Return(expectedActions, nil)

		actions, err := s.ListActions(ctx, uuid, appName)
		if err != nil {
			t.Fatalf("ListActions() unexpected error: %v", err)
		}
		if len(actions) != 1 {
			t.Fatalf("expected 1 action, got %d", len(actions))
		}
		if actions[0].Name != "action1" || !reflect.DeepEqual(actions[0].Spec, &actionSpec) {
			t.Errorf("expected action 'action1' with spec %+v, got %+v", &actionSpec, actions[0])
		}
	})

	t.Run("error_action_list", func(t *testing.T) {
		mockAction.EXPECT().List(ctx, uuid, appName).Return(nil, errExpected)
		_, err := s.ListActions(ctx, uuid, appName)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})
}

func TestNexusService_ListCharms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCharmHub := mocks.NewMockJujuCharmHub(ctrl)
	s := NewNexusService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockCharmHub, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()

	t.Run("success_with_filtering", func(t *testing.T) {
		allCharms := []model.Charm{
			{Name: "k8s-charm", Type: "charm", Result: model.CharmResult{DeployableOn: []string{"kubernetes"}}},       // Excluded: deployable on k8s
			{Name: "machine-charm", Type: "charm", Result: model.CharmResult{DeployableOn: []string{"aws", "azure"}}}, // Included
			{Name: "bundle-charm", Type: "bundle"}, // Excluded: type not charm
			{Name: "another-machine-charm", Type: "charm", Result: model.CharmResult{DeployableOn: []string{"lxd"}}}, // Included
		}
		mockCharmHub.EXPECT().List(ctx).Return(allCharms, nil)

		charms, err := s.ListCharms(ctx)
		if err != nil {
			t.Fatalf("ListCharms() unexpected error: %v", err)
		}
		if len(charms) != 2 {
			t.Fatalf("expected 2 charms after filtering, got %d", len(charms))
		}
		if charms[0].Name != "machine-charm" || charms[1].Name != "another-machine-charm" {
			t.Errorf("unexpected charms returned: %+v", charms)
		}
	})

	t.Run("error_charmhub_list", func(t *testing.T) {
		mockCharmHub.EXPECT().List(ctx).Return(nil, errExpected)
		_, err := s.ListCharms(ctx)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})
}

func TestNexusService_GetCharm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCharmHub := mocks.NewMockJujuCharmHub(ctrl)
	s := NewNexusService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockCharmHub, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()
	charmName := "test-charm"
	expectedCharm := &model.Charm{Name: charmName, Type: "charm"}

	t.Run("success", func(t *testing.T) {
		mockCharmHub.EXPECT().Get(ctx, charmName).Return(expectedCharm, nil)
		charm, err := s.GetCharm(ctx, charmName)
		if err != nil {
			t.Fatalf("GetCharm() unexpected error: %v", err)
		}
		if !reflect.DeepEqual(charm, expectedCharm) {
			t.Errorf("expected charm %+v, got %+v", expectedCharm, charm)
		}
	})

	t.Run("error_charmhub_get", func(t *testing.T) {
		mockCharmHub.EXPECT().Get(ctx, charmName).Return(nil, errExpected)
		_, err := s.GetCharm(ctx, charmName)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})
}

func TestNexusService_ListArtifacts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCharmHub := mocks.NewMockJujuCharmHub(ctrl)
	s := NewNexusService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockCharmHub, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()
	charmName := "test-charm"
	expectedArtifacts := []model.CharmArtifact{{Channel: model.CharmChannel{Name: "stable"}}}

	t.Run("success", func(t *testing.T) {
		mockCharmHub.EXPECT().ListArtifacts(ctx, charmName).Return(expectedArtifacts, nil)
		artifacts, err := s.ListArtifacts(ctx, charmName)
		if err != nil {
			t.Fatalf("ListArtifacts() unexpected error: %v", err)
		}
		if !reflect.DeepEqual(artifacts, expectedArtifacts) {
			t.Errorf("expected artifacts %+v, got %+v", expectedArtifacts, artifacts)
		}
	})

	t.Run("error_charmhub_list_artifacts", func(t *testing.T) {
		mockCharmHub.EXPECT().ListArtifacts(ctx, charmName).Return(nil, errExpected)
		_, err := s.ListArtifacts(ctx, charmName)
		if !errors.Is(err, errExpected) {
			t.Fatalf("expected error %v, got %v", errExpected, err)
		}
	})
}

// Helper to convert model.MachineConstraint to constraints.Value for mocking
// This is a simplified version. A more complete one would handle all fields.
// func toMockConstraint(mc *model.MachineConstraint) *constraints.Value {
// 	if mc == nil {
// 		return nil
// 	}
// 	cons := constraints.Value{}
// 	if mc.CPUCores != 0 {
// 		cores := uint64(mc.CPUCores)
// 		cons.CpuCores = &cores
// 	}
// 	// Add other constraint conversions as needed for tests
// 	return &cons
// }

// // Helper to convert model.MachinePlacement to instance.Placement for mocking
// // This is a simplified version.
// func toMockPlacement(mp *model.MachinePlacement, directive string) *instance.Placement {
// 	if mp == nil {
// 		return nil
// 	}
// 	return &instance.Placement{
// 		Scope:     mp.MachineID,
// 		Directive: directive, // Directive comes from maasToJujuMachineMap
// 	}
// }
