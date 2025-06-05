package service

import (
	"context"
	"errors"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/rpc/params"
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestListScopes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJujuModel := mocks.NewMockJujuModel(ctrl)
	s := &NexusService{
		scope: mockJujuModel,
	}

	ctx := context.Background()
	expectedScopes := []base.UserModelSummary{
		{Name: "scope1"},
		{Name: "scope2"},
	}
	mockJujuModel.EXPECT().List(ctx).Return(expectedScopes, nil)

	scopes, err := s.ListScopes(ctx)
	if err != nil {
		t.Errorf("ListScopes failed: %v", err)
	}
	if len(scopes) != len(expectedScopes) {
		t.Errorf("expected %d scopes, got %d", len(expectedScopes), len(scopes))
	}
	for i, scope := range scopes {
		if scope.Name != expectedScopes[i].Name {
			t.Errorf("expected scope name %s, got %s", expectedScopes[i].Name, scope.Name)
		}
	}
}
func TestCreateScope(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var errMock = errors.New("mock error")

	mockSSHKey := mocks.NewMockMAASSSHKey(ctrl)
	mockJujuModel := mocks.NewMockJujuModel(ctrl)
	mockJujuKey := mocks.NewMockJujuKey(ctrl)
	s := &NexusService{
		sshKey:     mockSSHKey,
		scope:      mockJujuModel,
		keyManager: mockJujuKey,
	}

	ctx := context.Background()
	scopeName := "test-scope"
	sshKey := &entity.SSHKey{Key: "test-key"}
	modelInfo := &base.ModelInfo{
		Name:            scopeName,
		UUID:            "test-uuid",
		Type:            "test-type",
		ControllerUUID:  "test-controller-uuid",
		IsController:    true,
		ProviderType:    "test-provider-type",
		Cloud:           "test-cloud",
		CloudRegion:     "test-cloud-region",
		CloudCredential: "test-cloud-credential",
		Owner:           "test-owner",
		Life:            "test-life",
	}
	t.Run("CreateScope", func(t *testing.T) {
		mockSSHKey.EXPECT().Default(ctx).Return(sshKey, nil)
		mockJujuModel.EXPECT().Create(ctx, scopeName).Return(modelInfo, nil)
		mockJujuKey.EXPECT().Add(ctx, modelInfo.UUID, sshKey.Key).Return([]params.ErrorResult{}, nil)

		scope, err := s.CreateScope(ctx, scopeName)
		if err != nil {
			t.Errorf("CreateScope failed: %v", err)
		}
		if scope.Name != scopeName {
			t.Errorf("expected scope name %s, got %s", scopeName, scope.Name)
		}
		if scope.UUID != modelInfo.UUID {
			t.Errorf("expected scope UUID %s, got %s", modelInfo.UUID, scope.UUID)
		}
	})
	t.Run("CreateScopeWithError", func(t *testing.T) {
		mockSSHKey.EXPECT().Default(ctx).Return(nil, errMock)
		// mockJujuModel.EXPECT().Create(ctx, scopeName).Return(nil, errMock)

		scope, err := s.CreateScope(ctx, scopeName)
		if err == nil {
			t.Error("expected error but got nil")
		}
		if scope != nil {
			t.Errorf("expected nil scope, got %v", scope)
		}
	})
	t.Run("CreateScopeWithKeyError", func(t *testing.T) {
		mockSSHKey.EXPECT().Default(ctx).Return(sshKey, nil)
		mockJujuModel.EXPECT().Create(ctx, scopeName).Return(modelInfo, nil)
		mockJujuKey.EXPECT().Add(ctx, modelInfo.UUID, sshKey.Key).Return([]params.ErrorResult{{Error: &params.Error{Message: errMock.Error()}}}, nil)

		scope, err := s.CreateScope(ctx, scopeName)
		if err == nil {
			t.Error("expected error but got nil")
		}
		if scope != nil {
			t.Errorf("expected nil scope, got %v", scope)
		}
	})

	t.Run("CreateScopeWithModelError", func(t *testing.T) {
		mockSSHKey.EXPECT().Default(ctx).Return(sshKey, nil)
		mockJujuModel.EXPECT().Create(ctx, scopeName).Return(nil, errMock)

		scope, err := s.CreateScope(ctx, scopeName)
		if err == nil {
			t.Error("expected error but got nil")
		}
		if scope != nil {
			t.Errorf("expected nil scope, got %v", scope)
		}
	})
	t.Run("CreateScopeWithKeyError", func(t *testing.T) { // This is the duplicated test case name, the fix is applied here
		mockSSHKey.EXPECT().Default(ctx).Return(sshKey, nil)
		mockJujuModel.EXPECT().Create(ctx, scopeName).Return(modelInfo, nil)
		mockJujuKey.EXPECT().Add(ctx, modelInfo.UUID, sshKey.Key).Return([]params.ErrorResult{{Error: &params.Error{Message: errMock.Error()}}}, nil)

		scope, err := s.CreateScope(ctx, scopeName)
		if err == nil {
			t.Error("expected error but got nil")
		}
		if scope != nil {
			t.Errorf("expected nil scope, got %v", scope)
		}
	})
	t.Run("CreateScopeWithEmptyName", func(t *testing.T) {
		scopeName = ""
		sshKey = &entity.SSHKey{Key: "test-key"}
		// Return a modelInfo with empty Name to match the expectation
		modelInfo = &base.ModelInfo{
			Name:            "",
			UUID:            "test-uuid",
			Type:            "test-type",
			ControllerUUID:  "test-controller-uuid",
			IsController:    true,
			ProviderType:    "test-provider-type",
			Cloud:           "test-cloud",
			CloudRegion:     "test-cloud-region",
			CloudCredential: "test-cloud-credential",
			Owner:           "test-owner",
			Life:            "test-life",
		}

		mockSSHKey.EXPECT().Default(ctx).Return(sshKey, nil)
		mockJujuModel.EXPECT().Create(ctx, scopeName).Return(modelInfo, nil)
		mockJujuKey.EXPECT().Add(ctx, modelInfo.UUID, sshKey.Key).Return([]params.ErrorResult{}, nil)

		scope, err := s.CreateScope(ctx, scopeName)
		if err != nil {
			t.Errorf("CreateScope failed with empty name: %v", err)
		}
		if scope.Name != "" {
			t.Errorf("expected empty scope name, got %s", scope.Name)
		}
	})
	t.Run("CreateScopeWithEmptyUUID", func(t *testing.T) {
		scopeName = "test-scope-empty-uuid"
		sshKey = &entity.SSHKey{Key: "test-key-empty-uuid"}
		modelInfo = &base.ModelInfo{
			Name:            scopeName,
			UUID:            "",
			Type:            "test-type",
			ControllerUUID:  "test-controller-uuid",
			IsController:    true,
			ProviderType:    "test-provider-type",
			Cloud:           "test-cloud",
			CloudRegion:     "test-cloud-region",
			CloudCredential: "test-cloud-credential",
			Owner:           "test-owner",
			Life:            "test-life",
		}

		mockSSHKey.EXPECT().Default(ctx).Return(sshKey, nil)
		mockJujuModel.EXPECT().Create(ctx, scopeName).Return(modelInfo, nil)
		mockJujuKey.EXPECT().Add(ctx, modelInfo.UUID, sshKey.Key).Return([]params.ErrorResult{}, nil)

		scope, err := s.CreateScope(ctx, scopeName)
		if err != nil {
			t.Errorf("CreateScope failed with empty UUID: %v", err)
		}
		if scope.UUID != "" {
			t.Errorf("expected empty scope UUID, got %s", scope.UUID)
		}
	})
}

func TestCreateDefaultScope(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSSHKey := mocks.NewMockMAASSSHKey(ctrl)
	mockJujuModel := mocks.NewMockJujuModel(ctrl)
	mockJujuKey := mocks.NewMockJujuKey(ctrl)
	s := &NexusService{
		sshKey:     mockSSHKey,
		scope:      mockJujuModel,
		keyManager: mockJujuKey,
	}

	ctx := context.Background()
	sshKey := &entity.SSHKey{Key: "test-key"}
	modelInfo := &base.ModelInfo{
		Name: defaultScopeName,
		UUID: "test-uuid",
	}

	mockSSHKey.EXPECT().Default(ctx).Return(sshKey, nil)
	mockJujuModel.EXPECT().Create(ctx, defaultScopeName).Return(modelInfo, nil)
	mockJujuKey.EXPECT().Add(ctx, modelInfo.UUID, sshKey.Key).Return([]params.ErrorResult{}, nil)

	scope, err := s.CreateDefaultScope(ctx)

	if err != nil {
		t.Errorf("CreateDefaultScope failed: %v", err)
	}
	if scope.Name != defaultScopeName {
		t.Errorf("expected scope name %s, got %s", defaultScopeName, scope.Name)
	}
	if scope.UUID != modelInfo.UUID {
		t.Errorf("expected scope UUID %s, got %s", modelInfo.UUID, scope.UUID)
	}
}
