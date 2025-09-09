package core

import (
	"context"
	"testing"
	"time"

	"github.com/juju/juju/api/client/action"
	"github.com/stretchr/testify/assert"
)

type storageMockActionRepo struct {
	runCount int
}

// List is a mock implementation to satisfy the ActionRepo interface.
func (m *storageMockActionRepo) List(ctx context.Context, uuid, leader string) (map[string]ActionSpec, error) {
	return map[string]ActionSpec{}, nil
}

func (m *storageMockActionRepo) RunCommand(ctx context.Context, uuid, leader, command string) (string, error) {
	m.runCount++
	return "action-id", nil
}

func (m *storageMockActionRepo) RunAction(ctx context.Context, uuid, leader, action string, params map[string]any) (string, error) {
	return "action-id", nil
}

func (m *storageMockActionRepo) GetResult(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
	// Simulate completed action on second call
	if m.runCount > 1 {
		return &action.ActionResult{
			Status: "completed",
			Output: map[string]any{
				"stdout": `
[global]
fsid = test-fsid
mon_host = 1.2.3.4
[client.admin]
key = test-key
`,
			},
		}, nil
	}
	// Simulate RGW user list
	return &action.ActionResult{
		Status: "completed",
		Output: map[string]any{
			"stdout": `["otterscale"]`,
		},
	}, nil
}

func TestExtractStorageCephConfig(t *testing.T) {
	result := &action.ActionResult{
		Output: map[string]any{
			"stdout": `
[global]
fsid = test-fsid
mon_host = 1.2.3.4
[client.admin]
key = test-key
`,
		},
	}
	cfg, err := extractStorageCephConfig(result)
	assert.NoError(t, err)
	assert.Equal(t, "test-fsid", cfg.FSID)
	assert.Equal(t, "1.2.3.4", cfg.MONHost)
	assert.Equal(t, "test-key", cfg.Key)
}

func TestGetRGWCommand(t *testing.T) {
	result := &action.ActionResult{
		Output: map[string]any{
			"stdout": `["otterscale"]`,
		},
	}
	cmd, err := getRGWCommand(result)
	assert.NoError(t, err)
	assert.Equal(t, cephRGWUserInfoCommand, cmd)
}

func TestExtractStorageRGWConfig(t *testing.T) {
	result := &action.ActionResult{
		Output: map[string]any{
			"stdout": `{"keys":[{"access_key":"ak","secret_key":"sk"}]}`,
		},
	}
	cfg, err := extractStorageRGWConfig(result)
	assert.NoError(t, err)
	assert.Equal(t, "ak", cfg.AccessKey)
	assert.Equal(t, "sk", cfg.SecretKey)
}

func TestRunCommandAndWaitForActionCompleted(t *testing.T) {
	actionRepo := &storageMockActionRepo{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := runCommand(ctx, actionRepo, "uuid", "leader", "cmd")
	assert.NoError(t, err)
}

func TestRgwNameAndNfsName(t *testing.T) {
	assert.Equal(t, "foo-bar-radosgw", rgwName("foo-bar-mon"))
	assert.Equal(t, "foo-bar-nfs", nfsName("foo-bar-mon"))
}
