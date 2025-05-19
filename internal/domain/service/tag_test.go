package service

import (
	"context"
	"errors"
	"testing"

	"github.com/openhdc/otterscale/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) List(ctx context.Context) ([]model.Tag, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Tag), args.Error(1)
}

func (m *MockTagRepository) Get(ctx context.Context, name string) (*model.Tag, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockTagRepository) Create(ctx context.Context, name, comment string) (*model.Tag, error) {
	args := m.Called(ctx, name, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockTagRepository) Delete(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockTagRepository) AddMachines(ctx context.Context, tagName string, machineIDs []string) error {
	args := m.Called(ctx, tagName, machineIDs)
	return args.Error(0)
}

func (m *MockTagRepository) RemoveMachines(ctx context.Context, tagName string, machineIDs []string) error {
	args := m.Called(ctx, tagName, machineIDs)
	return args.Error(0)
}

func TestNexusService_ListTags(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := &NexusService{tag: mockRepo}
	ctx := context.Background()

	expectedTags := []model.Tag{
		{Name: "tag1", Comment: "Test tag 1"},
		{Name: "tag2", Comment: "Test tag 2"},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("List", ctx).Return(expectedTags, nil).Once()

		tags, err := service.ListTags(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedTags, tags)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("repository error")
		mockRepo.On("List", ctx).Return([]model.Tag{}, expectedErr).Once()

		tags, err := service.ListTags(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, tags)
		mockRepo.AssertExpectations(t)
	})
}

func TestNexusService_GetTag(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := &NexusService{tag: mockRepo}
	ctx := context.Background()
	tagName := "test-tag"

	t.Run("success", func(t *testing.T) {
		expectedTag := &model.Tag{Name: tagName, Comment: "Test tag"}
		mockRepo.On("Get", ctx, tagName).Return(expectedTag, nil).Once()

		tag, err := service.GetTag(ctx, tagName)

		assert.NoError(t, err)
		assert.Equal(t, expectedTag, tag)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		expectedErr := errors.New("tag not found")
		mockRepo.On("Get", ctx, tagName).Return((*model.Tag)(nil), expectedErr).Once()

		tag, err := service.GetTag(ctx, tagName)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, tag)
		mockRepo.AssertExpectations(t)
	})
}

func TestNexusService_CreateTag(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := &NexusService{tag: mockRepo}
	ctx := context.Background()
	tagName := "new-tag"
	comment := "New test tag"

	t.Run("success", func(t *testing.T) {
		expectedTag := &model.Tag{Name: tagName, Comment: comment}
		mockRepo.On("Create", ctx, tagName, comment).Return(expectedTag, nil).Once()

		tag, err := service.CreateTag(ctx, tagName, comment)

		assert.NoError(t, err)
		assert.Equal(t, expectedTag, tag)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("failed to create tag")
		mockRepo.On("Create", ctx, tagName, comment).Return((*model.Tag)(nil), expectedErr).Once()

		tag, err := service.CreateTag(ctx, tagName, comment)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, tag)
		mockRepo.AssertExpectations(t)
	})
}

func TestNexusService_DeleteTag(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := &NexusService{tag: mockRepo}
	ctx := context.Background()
	tagName := "delete-tag"

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", ctx, tagName).Return(nil).Once()

		err := service.DeleteTag(ctx, tagName)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("failed to delete tag")
		mockRepo.On("Delete", ctx, tagName).Return(expectedErr).Once()

		err := service.DeleteTag(ctx, tagName)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNexusService_AddMachineTags(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := &NexusService{tag: mockRepo}
	ctx := context.Background()
	machineID := "machine-123"
	tags := []string{"tag1", "tag2"}

	t.Run("success", func(t *testing.T) {
		for _, tag := range tags {
			mockRepo.On("AddMachines", mock.Anything, tag, []string{machineID}).Return(nil).Once()
		}

		err := service.AddMachineTags(ctx, machineID, tags)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("partial error", func(t *testing.T) {
		expectedErr := errors.New("failed to add machine to tag")
		mockRepo.On("AddMachines", mock.Anything, tags[0], []string{machineID}).Return(nil).Once()
		mockRepo.On("AddMachines", mock.Anything, tags[1], []string{machineID}).Return(expectedErr).Once()

		err := service.AddMachineTags(ctx, machineID, tags)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty tags", func(t *testing.T) {
		err := service.AddMachineTags(ctx, machineID, []string{})
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t) // No calls expected
	})
}

func TestNexusService_RemoveMachineTags(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := &NexusService{tag: mockRepo}
	ctx := context.Background()
	machineID := "machine-123"
	tags := []string{"tag1", "tag2"}

	t.Run("success", func(t *testing.T) {
		for _, tag := range tags {
			mockRepo.On("RemoveMachines", mock.Anything, tag, []string{machineID}).Return(nil).Once()
		}

		err := service.RemoveMachineTags(ctx, machineID, tags)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("partial error", func(t *testing.T) {
		expectedErr := errors.New("failed to remove machine from tag")
		mockRepo.On("RemoveMachines", mock.Anything, tags[0], []string{machineID}).Return(nil).Once()
		mockRepo.On("RemoveMachines", mock.Anything, tags[1], []string{machineID}).Return(expectedErr).Once()

		err := service.RemoveMachineTags(ctx, machineID, tags)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}
