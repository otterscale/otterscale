package core

import (
	"context"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/node"
	"github.com/stretchr/testify/assert"
)

func TestCreateCeph(t *testing.T) {
	serverRepo := &essMockServerRepo{}
	machineRepo := &essMockMachineRepo{
		machines: []Machine{
			{
				Machine: &entity.Machine{
					SystemID:            "test-machine",
					Status:              node.StatusDeployed,
					WorkloadAnnotations: map[string]string{"juju-machine-id": "1"},
				},
			},
		},
	}
	facilityRepo := &essMockFacilityRepo{}
	tagRepo := &mockTagRepo{}

	err := CreateCeph(context.Background(), serverRepo, machineRepo, facilityRepo, tagRepo, "test-uuid", "test-machine", "test-prefix", map[string]string{
		"ch:ceph-mon": "{\"test-prefix-ceph-mon\":{}}",
		"ch:ceph-osd": "{\"test-prefix-ceph-osd\":{}}",
	})
	assert.NoError(t, err)
}

func TestNewCephConfigs(t *testing.T) {
	prefix := "test-prefix"
	osdDevices := "/dev/sdb /dev/sdc"
	nfsVIP := "192.168.1.100"

	configs, err := newCephConfigs(prefix, osdDevices, nfsVIP)
	assert.NoError(t, err)

	assert.Contains(t, configs, "ch:ceph-mon")
	assert.Contains(t, configs, "ch:ceph-osd")
	assert.Contains(t, configs, "ch:ceph-nfs")
}
