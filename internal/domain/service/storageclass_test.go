package service

import (
	"context"
	"testing"

	"github.com/openhdc/otterscale/internal/domain/model"
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestNexusService_ListStorageClasses_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockKubeStorage(ctrl)
	mockKube := mocks.NewMockKubeClient(ctrl)

	ns := &NexusService{
		storage:    mockStorage,
		kubernetes: mockKube,
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facility := "test-facility"

	mockKube.EXPECT().Exists(uuid, facility).Return(true)
	mockStorage.EXPECT().ListStorageClasses(ctx, uuid, facility).Return([]model.StorageClass{{}}, nil)

	classes, err := ns.ListStorageClasses(ctx, uuid, facility)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(classes) != 1 {
		t.Errorf("expected 1 storage class, got %d", len(classes))
	}
}

// Variable to hold the function that can be replaced in tests
var SetKubernetesClient = func(ns *NexusService, ctx context.Context, uuid, facility string) error {
	return ns.setKubernetesClient(ctx, uuid, facility)
}
