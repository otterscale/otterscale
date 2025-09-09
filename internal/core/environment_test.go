package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockEnvRepo struct{}

func (m *mockEnvRepo) ListEnvironments(ctx context.Context) ([]string, error) {
	return []string{"dev"}, nil
}

func TestEnvironment_ListEnvironments(t *testing.T) {
	repo := &mockEnvRepo{}
	envs, err := repo.ListEnvironments(context.Background())
	assert.NoError(t, err)
	assert.Contains(t, envs, "dev")
}
