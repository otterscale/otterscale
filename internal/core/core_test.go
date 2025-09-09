package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type coreMockMachineRepo struct{}

func (m *coreMockMachineRepo) Get(ctx context.Context, name string) (string, error) {
	return "focal", nil
}

func (m *coreMockMachineRepo) Update(ctx context.Context, name, value string) error {
	return nil
}

func TestBoolToInt(t *testing.T) {
	assert.Equal(t, 1, boolToInt(true))
	assert.Equal(t, 0, boolToInt(false))
}

func TestDefaultBase(t *testing.T) {
	mb := &coreMockMachineRepo{}
	b, err := defaultBase(context.Background(), mb)
	assert.NoError(t, err)
	assert.Equal(t, "ubuntu", b.OS)
}

func TestGetJujuModelUUID(t *testing.T) {
	m := map[string]string{"juju-model-uuid": "uuid-123"}
	uuid, err := getJujuModelUUID(m)
	assert.NoError(t, err)
	assert.Equal(t, "uuid-123", uuid)
	_, err = getJujuModelUUID(map[string]string{})
	assert.Error(t, err)
}

func TestGetJujuMachineID(t *testing.T) {
	m := map[string]string{"juju-machine-id": "juju-0-1"}
	id, err := getJujuMachineID(m)
	assert.NoError(t, err)
	assert.Equal(t, "1", id)
	_, err = getJujuMachineID(map[string]string{})
	assert.Error(t, err)
}

func TestToPlacement(t *testing.T) {
	p := &MachinePlacement{LXD: true}
	placement := toPlacement(p, "foo")
	assert.Equal(t, "lxd", placement.Scope)
	assert.Equal(t, "foo", placement.Directive)
}

func TestToConstraint(t *testing.T) {
	c := &MachineConstraint{
		Architecture: "amd64",
		CPUCores:     2,
		MemoryMB:     4096,
		Tags:         []string{"foo"},
	}
	val := toConstraint(c)
	assert.NotNil(t, val.Arch)
	assert.NotNil(t, val.CpuCores)
	assert.NotNil(t, val.Mem)
	assert.NotNil(t, val.Tags)
}
