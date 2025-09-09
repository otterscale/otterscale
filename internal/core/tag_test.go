package core

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock TagRepo for testing
type mockTagRepo struct {
	tags []Tag
}

func (m *mockTagRepo) List(ctx context.Context) ([]Tag, error) {
	return m.tags, nil
}

func (m *mockTagRepo) Get(ctx context.Context, name string) (*Tag, error) {
	for _, tag := range m.tags {
		if tag.Name == name {
			return &tag, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockTagRepo) Create(ctx context.Context, name, comment string) (*Tag, error) {
	tag := Tag{Name: name, Comment: comment}
	m.tags = append(m.tags, tag)
	return &tag, nil
}

func (m *mockTagRepo) Delete(ctx context.Context, name string) error {
	for i, tag := range m.tags {
		if tag.Name == name {
			m.tags = append(m.tags[:i], m.tags[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockTagRepo) AddMachines(ctx context.Context, name string, machineIDs []string) error {
	return nil
}

func (m *mockTagRepo) RemoveMachines(ctx context.Context, name string, machineIDs []string) error {
	return nil
}

func TestTagUseCase_ListTags(t *testing.T) {
	mock := &mockTagRepo{tags: []Tag{{Name: "tag1"}, {Name: "tag2"}}}
	uc := NewTagUseCase(mock)
	tags, err := uc.ListTags(context.Background())
	assert.NoError(t, err)
	assert.Len(t, tags, 2)
}

func TestTagUseCase_GetTag(t *testing.T) {
	mock := &mockTagRepo{tags: []Tag{{Name: "tag1"}}}
	uc := NewTagUseCase(mock)
	tag, err := uc.GetTag(context.Background(), "tag1")
	assert.NoError(t, err)
	assert.Equal(t, "tag1", tag.Name)

	// test not found
	tag, err = uc.GetTag(context.Background(), "notfound")
	assert.Error(t, err)
	assert.Nil(t, tag)
}

func TestTagUseCase_CreateTag(t *testing.T) {
	mock := &mockTagRepo{}
	uc := NewTagUseCase(mock)
	tag, err := uc.CreateTag(context.Background(), "tag3", "comment")
	assert.NoError(t, err)
	assert.Equal(t, "tag3", tag.Name)
	assert.Equal(t, "comment", tag.Comment)
}

func TestTagUseCase_DeleteTag(t *testing.T) {
	mock := &mockTagRepo{tags: []Tag{{Name: "tag1"}, {Name: "tag2"}}}
	uc := NewTagUseCase(mock)
	err := uc.DeleteTag(context.Background(), "tag1")
	assert.NoError(t, err)
	assert.Len(t, mock.tags, 1)

	// test delete not found
	err = uc.DeleteTag(context.Background(), "notfound")
	assert.Error(t, err)
}
