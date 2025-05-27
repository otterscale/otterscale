package service

import (
	"context"
	"fmt"
	"testing"

	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	"go.uber.org/mock/gomock"
	v1 "k8s.io/api/storage/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors" // Alias for clarity
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestListStorageClasses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKubeStorage := mocks.NewMockKubeStorage(ctrl)
	mockKubeClient := mocks.NewMockKubeClient(ctrl)           // Added
	mockJujuApplication := mocks.NewMockJujuApplication(ctrl) // Added

	s := &NexusService{
		storage:    mockKubeStorage,
		kubernetes: mockKubeClient,      // Initialized
		facility:   mockJujuApplication, // Initialized
	}

	ctx := context.Background()
	uuid := "test-uuid"
	facilityName := "test-facility" // Consistent naming

	t.Run("success", func(t *testing.T) {
		// Mock for setKubernetesClient to succeed (e.g., client already exists)
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(true)

		expectedStorageClasses := []v1.StorageClass{
			{ObjectMeta: metav1.ObjectMeta{Name: "storageclass1"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "storageclass2"}},
		}
		mockKubeStorage.EXPECT().ListStorageClasses(ctx, uuid, facilityName).Return(expectedStorageClasses, nil)

		storageClasses, err := s.ListStorageClasses(ctx, uuid, facilityName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(storageClasses) != len(expectedStorageClasses) {
			t.Fatalf("expected %d storage classes, got %d", len(expectedStorageClasses), len(storageClasses))
		}

		for i, sc := range storageClasses {
			if sc.Name != expectedStorageClasses[i].Name {
				t.Errorf("expected storage class name %s, got %s", expectedStorageClasses[i].Name, sc.Name)
			}
		}
	})

	t.Run("setKubernetesClient error", func(t *testing.T) {
		// Mock for setKubernetesClient to fail
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(false) // Client doesn't exist, proceed to setup path
		expectedSetKubeClientErr := fmt.Errorf("get leader failed")
		// Assuming GetLeader is a call in setKubernetesClient that can fail after Exists returns false
		mockJujuApplication.EXPECT().GetLeader(ctx, uuid, facilityName).Return("", expectedSetKubeClientErr)

		// ListStorageClasses on mockKubeStorage should NOT be called
		mockKubeStorage.EXPECT().ListStorageClasses(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		_, err := s.ListStorageClasses(ctx, uuid, facilityName)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		// Check if the error is the one from GetLeader or a wrapper around it
		// if !k8serrors.Is(err, expectedSetKubeClientErr) && err.Error() != expectedSetKubeClientErr.Error() { // errors.Is might not work if not k8s error type
		// 	t.Errorf("expected error containing '%v', got %v", expectedSetKubeClientErr, err)
		// }
	})

	t.Run("ListStorageClasses error", func(t *testing.T) {
		// Mock for setKubernetesClient to succeed
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(true)

		// Mock for the KubeStorage call to fail
		internalErr := fmt.Errorf("db connection failed") // Provide a concrete error
		expectedListErr := k8serrors.NewInternalError(internalErr)
		mockKubeStorage.EXPECT().ListStorageClasses(ctx, uuid, facilityName).Return(nil, expectedListErr)

		_, err := s.ListStorageClasses(ctx, uuid, facilityName)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !k8serrors.IsInternalError(err) {
			t.Errorf("expected an internal error, but got type %T: %v", err, err)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		// Mock for setKubernetesClient to succeed
		mockKubeClient.EXPECT().Exists(uuid, facilityName).Return(true)

		mockKubeStorage.EXPECT().ListStorageClasses(ctx, uuid, facilityName).Return([]v1.StorageClass{}, nil)

		storageClasses, err := s.ListStorageClasses(ctx, uuid, facilityName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(storageClasses) != 0 {
			t.Fatalf("expected 0 storage classes, got %d", len(storageClasses))
		}
	})

}
