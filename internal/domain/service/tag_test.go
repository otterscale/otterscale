package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/openhdc/otterscale/internal/domain/model"
	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
)

func TestNexusService_ListTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMAASTag := mocks.NewMockMAASTag(ctrl)
	s := &NexusService{tag: mockMAASTag}
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedTags := []model.Tag{
			{Name: "tag1", Comment: "comment1"},
			{Name: "tag2", Comment: "comment2"},
		}
		mockMAASTag.EXPECT().List(ctx).Return(expectedTags, nil)

		tags, err := s.ListTags(ctx)
		if err != nil {
			t.Fatalf("ListTags() error = %v, wantErr nil", err)
		}
		if !reflect.DeepEqual(tags, expectedTags) {
			t.Errorf("ListTags() got = %v, want %v", tags, expectedTags)
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("maas list error")
		mockMAASTag.EXPECT().List(ctx).Return(nil, expectedErr)

		_, err := s.ListTags(ctx)
		if !errors.Is(err, expectedErr) {
			t.Errorf("ListTags() error = %v, want %v", err, expectedErr)
		}
	})
}

func TestNexusService_GetTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMAASTag := mocks.NewMockMAASTag(ctrl)
	s := &NexusService{tag: mockMAASTag}
	ctx := context.Background()
	tagName := "test-tag"

	t.Run("success", func(t *testing.T) {
		expectedTag := &model.Tag{Name: tagName, Comment: "A test tag"}
		mockMAASTag.EXPECT().Get(ctx, tagName).Return(expectedTag, nil)

		tag, err := s.GetTag(ctx, tagName)
		if err != nil {
			t.Fatalf("GetTag() error = %v, wantErr nil", err)
		}
		if !reflect.DeepEqual(tag, expectedTag) {
			t.Errorf("GetTag() got = %v, want %v", tag, expectedTag)
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("maas get error")
		mockMAASTag.EXPECT().Get(ctx, tagName).Return(nil, expectedErr)

		_, err := s.GetTag(ctx, tagName)
		if !errors.Is(err, expectedErr) {
			t.Errorf("GetTag() error = %v, want %v", err, expectedErr)
		}
	})
}

func TestNexusService_CreateTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMAASTag := mocks.NewMockMAASTag(ctrl)
	s := &NexusService{tag: mockMAASTag}
	ctx := context.Background()
	tagName := "new-tag"
	tagComment := "A new comment"

	t.Run("success", func(t *testing.T) {
		expectedTag := &model.Tag{Name: tagName, Comment: tagComment}
		mockMAASTag.EXPECT().Create(ctx, tagName, tagComment).Return(expectedTag, nil)

		tag, err := s.CreateTag(ctx, tagName, tagComment)
		if err != nil {
			t.Fatalf("CreateTag() error = %v, wantErr nil", err)
		}
		if !reflect.DeepEqual(tag, expectedTag) {
			t.Errorf("CreateTag() got = %v, want %v", tag, expectedTag)
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("maas create error")
		mockMAASTag.EXPECT().Create(ctx, tagName, tagComment).Return(nil, expectedErr)

		_, err := s.CreateTag(ctx, tagName, tagComment)
		if !errors.Is(err, expectedErr) {
			t.Errorf("CreateTag() error = %v, want %v", err, expectedErr)
		}
	})
}

func TestNexusService_DeleteTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMAASTag := mocks.NewMockMAASTag(ctrl)
	s := &NexusService{tag: mockMAASTag}
	ctx := context.Background()
	tagName := "tag-to-delete"

	t.Run("success", func(t *testing.T) {
		mockMAASTag.EXPECT().Delete(ctx, tagName).Return(nil)

		err := s.DeleteTag(ctx, tagName)
		if err != nil {
			t.Fatalf("DeleteTag() error = %v, wantErr nil", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("maas delete error")
		mockMAASTag.EXPECT().Delete(ctx, tagName).Return(expectedErr)

		err := s.DeleteTag(ctx, tagName)
		if !errors.Is(err, expectedErr) {
			t.Errorf("DeleteTag() error = %v, want %v", err, expectedErr)
		}
	})
}

func TestNexusService_AddMachineTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMAASTag := mocks.NewMockMAASTag(ctrl)
	s := &NexusService{tag: mockMAASTag}
	ctx := context.Background()
	machineID := "machine-123"
	tagsToAdd := []string{"tagA", "tagB"}

	t.Run("success", func(t *testing.T) {
		mockMAASTag.EXPECT().AddMachines(gomock.Any(), "tagA", []string{machineID}).Return(nil)
		mockMAASTag.EXPECT().AddMachines(gomock.Any(), "tagB", []string{machineID}).Return(nil)

		err := s.AddMachineTags(ctx, machineID, tagsToAdd)
		if err != nil {
			t.Fatalf("AddMachineTags() error = %v, wantErr nil", err)
		}
	})

	t.Run("error_on_one_tag", func(t *testing.T) {
		expectedErr := errors.New("maas addmachines error for tagB")
		mockMAASTag.EXPECT().AddMachines(gomock.Any(), "tagA", []string{machineID}).Return(nil)
		mockMAASTag.EXPECT().AddMachines(gomock.Any(), "tagB", []string{machineID}).Return(expectedErr)

		err := s.AddMachineTags(ctx, machineID, tagsToAdd)
		if !errors.Is(err, expectedErr) {
			t.Errorf("AddMachineTags() error = %v, want %v", err, expectedErr)
		}
	})

	t.Run("empty_tags_list", func(t *testing.T) {
		// No calls to mockMAASTag.AddMachines are expected
		err := s.AddMachineTags(ctx, machineID, []string{})
		if err != nil {
			t.Fatalf("AddMachineTags() with empty list error = %v, wantErr nil", err)
		}
	})
}

func TestNexusService_RemoveMachineTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMAASTag := mocks.NewMockMAASTag(ctrl)
	s := &NexusService{tag: mockMAASTag}
	ctx := context.Background()
	machineID := "machine-xyz"
	tagsToRemove := []string{"tagX", "tagY"}

	t.Run("success", func(t *testing.T) {
		mockMAASTag.EXPECT().RemoveMachines(gomock.Any(), "tagX", []string{machineID}).Return(nil)
		mockMAASTag.EXPECT().RemoveMachines(gomock.Any(), "tagY", []string{machineID}).Return(nil)

		err := s.RemoveMachineTags(ctx, machineID, tagsToRemove)
		if err != nil {
			t.Fatalf("RemoveMachineTags() error = %v, wantErr nil", err)
		}
	})

	t.Run("error_on_one_tag", func(t *testing.T) {
		expectedErr := errors.New("maas removemachines error for tagY")
		mockMAASTag.EXPECT().RemoveMachines(gomock.Any(), "tagX", []string{machineID}).Return(nil)
		mockMAASTag.EXPECT().RemoveMachines(gomock.Any(), "tagY", []string{machineID}).Return(expectedErr)

		err := s.RemoveMachineTags(ctx, machineID, tagsToRemove)
		if !errors.Is(err, expectedErr) {
			t.Errorf("RemoveMachineTags() error = %v, want %v", err, expectedErr)
		}
	})

	t.Run("empty_tags_list", func(t *testing.T) {
		// No calls to mockMAASTag.RemoveMachines are expected
		err := s.RemoveMachineTags(ctx, machineID, []string{})
		if err != nil {
			t.Fatalf("RemoveMachineTags() with empty list error = %v, wantErr nil", err)
		}
	})
}
