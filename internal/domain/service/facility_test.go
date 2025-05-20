package service

import (
	"context"
	"errors"
	"testing"

	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/rpc/params"
	"github.com/openhdc/otterscale/internal/domain/model"
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestNexusService_ListFacilities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJuju := mocks.NewMockJujuClient(ctrl)
	ns := &NexusService{client: mockJuju}

	ctx := context.Background()
	uuid := "test-uuid"

	status := &params.FullStatus{
		Applications: map[string]params.ApplicationStatus{
			"fac1": {Charm: "foo"},
			"fac2": {Charm: "bar"},
		},
	}
	mockJuju.EXPECT().Status(ctx, uuid, []string{"application", "*"}).Return(status, nil)

	fs, err := ns.ListFacilities(ctx, uuid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fs) != 2 {
		t.Errorf("expected 2 facilities, got %d", len(fs))
	}
}

func TestNexusService_ListFacilities_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJuju := mocks.NewMockJujuClient(ctrl)
	ns := &NexusService{client: mockJuju}

	ctx := context.Background()
	uuid := "test-uuid"

	mockJuju.EXPECT().Status(ctx, uuid, []string{"application", "*"}).Return(nil, errors.New("fail"))
	_, err := ns.ListFacilities(ctx, uuid)
	if err == nil {
		t.Error("expected error")
	}
}

func TestNexusService_GetFacility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJuju := mocks.NewMockJujuClient(ctrl)
	mockFacility := mocks.NewMockJujuApplication(ctrl)
	ns := &NexusService{client: mockJuju, facility: mockFacility}

	ctx := context.Background()
	uuid := "test-uuid"
	name := "fac1"

	status := &params.FullStatus{
		Applications: map[string]params.ApplicationStatus{
			"fac1": {Charm: "foo"},
		},
	}
	cfg := map[string]any{"foo": "bar"}

	mockJuju.EXPECT().Status(ctx, uuid, []string{"application", name}).Return(status, nil)
	mockFacility.EXPECT().GetConfig(ctx, uuid, name).Return(cfg, nil)

	fac, err := ns.GetFacility(ctx, uuid, name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fac == nil || fac.Name != name {
		t.Errorf("unexpected facility: %+v", fac)
	}
}

func TestNexusService_GetFacility_notfound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJuju := mocks.NewMockJujuClient(ctrl)
	ns := &NexusService{client: mockJuju}

	ctx := context.Background()
	uuid := "test-uuid"
	name := "fac1"

	status := &params.FullStatus{
		Applications: map[string]params.ApplicationStatus{
			"other": {Charm: "foo"},
		},
	}
	mockJuju.EXPECT().Status(ctx, uuid, []string{"application", name}).Return(status, nil)

	_, err := ns.GetFacility(ctx, uuid, name)
	if err == nil {
		t.Error("expected error for not found")
	}
}

// func TestNexusService_CreateFacility(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockFacility := mocks.NewMockJujuApplication(ctrl)
// 	ns := &NexusService{facility: mockFacility}

// 	ctx := context.Background()
// 	uuid := "test-uuid"
// 	name := "fac1"
// 	configYAML := "foo: bar"
// 	charmName := "mycharm"
// 	channel := "stable"
// 	revision := 1
// 	number := 1
// 	mps := []model.MachinePlacement{}
// 	mc := &model.MachineConstraint{}
// 	trust := true

// 	// imageBase and toPlacements are not mocked here, so this is a basic call
// 	mockFacility.EXPECT().Create(ctx, uuid, name, configYAML, charmName, channel, revision, number, gomock.Any(), gomock.Any(), gomock.Any(), trust).Return(&application.DeployInfo{}, nil)

// 	_, err := ns.CreateFacility(ctx, uuid, name, configYAML, charmName, channel, revision, number, mps, mc, trust)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// }

func TestNexusService_UpdateFacility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFacility := mocks.NewMockJujuApplication(ctrl)
	mockJuju := mocks.NewMockJujuClient(ctrl)
	ns := &NexusService{facility: mockFacility, client: mockJuju}

	ctx := context.Background()
	uuid := "test-uuid"
	name := "fac1"
	configYAML := "foo: bar"

	mockFacility.EXPECT().Update(ctx, uuid, name, configYAML).Return(nil)
	status := &params.FullStatus{
		Applications: map[string]params.ApplicationStatus{
			"fac1": {Charm: "foo"},
		},
	}
	mockJuju.EXPECT().Status(ctx, uuid, []string{"application", name}).Return(status, nil)
	mockFacility.EXPECT().GetConfig(ctx, uuid, name).Return(map[string]any{}, nil)

	_, err := ns.UpdateFacility(ctx, uuid, name, configYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNexusService_DeleteFacility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFacility := mocks.NewMockJujuApplication(ctrl)
	ns := &NexusService{facility: mockFacility}

	ctx := context.Background()
	uuid := "test-uuid"
	name := "fac1"

	mockFacility.EXPECT().Delete(ctx, uuid, name, true, false).Return(nil)
	err := ns.DeleteFacility(ctx, uuid, name, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNexusService_ExposeFacility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFacility := mocks.NewMockJujuApplication(ctrl)
	ns := &NexusService{facility: mockFacility}

	ctx := context.Background()
	uuid := "test-uuid"
	name := "fac1"

	mockFacility.EXPECT().Expose(ctx, uuid, name, nil).Return(nil)
	err := ns.ExposeFacility(ctx, uuid, name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNexusService_AddFacilityUnits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFacility := mocks.NewMockJujuApplication(ctrl)
	ns := &NexusService{facility: mockFacility}

	ctx := context.Background()
	uuid := "test-uuid"
	name := "fac1"
	number := 2
	mps := []model.MachinePlacement{}

	mockFacility.EXPECT().AddUnits(ctx, uuid, name, number, gomock.Any()).Return([]string{"unit1", "unit2"}, nil)
	units, err := ns.AddFacilityUnits(ctx, uuid, name, number, mps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(units) != 2 {
		t.Errorf("expected 2 units, got %d", len(units))
	}
}

func TestNexusService_ListCharms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCharmHub := mocks.NewMockJujuCharmHub(ctrl)
	ns := &NexusService{charmhub: mockCharmHub}

	ctx := context.Background()
	mockCharmHub.EXPECT().List(ctx).Return([]model.Charm{
		{Type: "charm", Result: model.CharmResult{DeployableOn: []string{}}},
		{Type: "bundle", Result: model.CharmResult{DeployableOn: []string{}}},
		{Type: "charm", Result: model.CharmResult{DeployableOn: []string{"kubernetes"}}},
	}, nil)

	charms, err := ns.ListCharms(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(charms) != 1 {
		t.Errorf("expected 1 charm, got %d", len(charms))
	}
}

func TestNexusService_GetCharm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCharmHub := mocks.NewMockJujuCharmHub(ctrl)
	ns := &NexusService{charmhub: mockCharmHub}

	ctx := context.Background()
	mockCharmHub.EXPECT().Get(ctx, "foo").Return(&model.Charm{Name: "foo"}, nil)

	charm, err := ns.GetCharm(ctx, "foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if charm == nil || charm.Name != "foo" {
		t.Errorf("unexpected charm: %+v", charm)
	}
}

func TestNexusService_ListArtifacts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCharmHub := mocks.NewMockJujuCharmHub(ctrl)
	ns := &NexusService{charmhub: mockCharmHub}

	ctx := context.Background()
	// Use empty struct, as model.CharmArtifact likely has no Name field
	mockCharmHub.EXPECT().ListArtifacts(ctx, "foo").Return([]model.CharmArtifact{{}}, nil)

	arts, err := ns.ListArtifacts(ctx, "foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(arts) != 1 {
		t.Errorf("unexpected artifacts: %+v", arts)
	}
}

func TestNexusService_ListActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAction := mocks.NewMockJujuAction(ctrl)
	ns := &NexusService{action: mockAction}

	ctx := context.Background()
	uuid := "test-uuid"
	appName := "app1"

	mockAction.EXPECT().List(ctx, uuid, appName).Return(map[string]action.ActionSpec{
		"act1": {},
	}, nil)

	actions, err := ns.ListActions(ctx, uuid, appName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(actions) != 1 || actions[0].Name != "act1" {
		t.Errorf("unexpected actions: %+v", actions)
	}
}
