package core

import (
	"context"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/canonical/gomaasclient/entity/node"
	"github.com/stretchr/testify/assert"
)

func TestCreateKubernetes(t *testing.T) {
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

	err := CreateKubernetes(context.Background(), serverRepo, machineRepo, facilityRepo, tagRepo, "test-uuid", "test-machine", "test-prefix", map[string]string{
		"ch:kubernetes-control-plane": "{\"test-prefix-kubernetes-control-plane\":{}}",
		"ch:kubernetes-worker":        "{\"test-prefix-kubernetes-worker\":{}}",
	})
	assert.NoError(t, err)
}

func TestNewKubernetesConfigs(t *testing.T) {
	prefix := "test-prefix"
	vip := "192.168.1.100"
	cidr := "10.244.0.0/16"

	configs, err := newKubernetesConfigs(prefix, vip, cidr)
	assert.NoError(t, err)

	assert.Contains(t, configs, "ch:kubernetes-control-plane")
	assert.Contains(t, configs, "ch:keepalived")
	assert.Contains(t, configs, "ch:calico")
	assert.Contains(t, configs, "ch:containerd")
	assert.Contains(t, configs, "ch:kubeapi-load-balancer")
}

func TestFormatCharmFunctions(t *testing.T) {
	appName, ok := formatAppCharm("ch:amd64/kubernetes-control-plane-123")
	assert.True(t, ok)
	assert.Equal(t, "kubernetes-control-plane", appName)

	_, ok = formatAppCharm("invalid-name")
	assert.False(t, ok)

	essCharm := formatEssentialCharm("ch:kubernetes-control-plane")
	assert.Equal(t, "kubernetes-control-plane", essCharm)
}

func TestToEssentialName(t *testing.T) {
	name := toEssentialName("test-prefix", "ch:charm-name")
	assert.Equal(t, "test-prefix-charm-name", name)

	name = toEssentialName("test-prefix", "simple-name")
	assert.Equal(t, "test-prefix-simple-name", name)
}
