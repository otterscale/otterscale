package core

import (
	"context"
	"testing"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommon(t *testing.T) {
	serverRepo := &essMockServerRepo{}
	machineRepo := &essMockMachineRepo{}
	facilityRepo := &essMockFacilityRepo{}
	facilityOffersRepo := &mockFacilityOffersRepo{}
	conf := &config.Config{
		Juju: config.Juju{
			Username:   "test-user",
			Controller: "test-controller",
		},
	}

	uuid := "test-uuid"
	prefix := "test-prefix"
	configs := map[string]string{
		"ch:prometheus":               "{\"test-prefix-prometheus\":{}}",
		"ch:grafana":                  "{\"test-prefix-grafana\":{}}",
		"ch:elasticsearch":            "{\"test-prefix-elasticsearch\":{}}",
		"ch:filebeat":                 "{\"test-prefix-filebeat\":{}}",
		"ch:prometheus-node-exporter": "{\"test-prefix-prometheus-node-exporter\":{}}",
	}

	err := CreateCommon(context.Background(), serverRepo, machineRepo, facilityRepo, facilityOffersRepo, conf, uuid, prefix, configs)
	assert.NoError(t, err)
}

func TestNewCommonConfigs(t *testing.T) {
	prefix := "test-prefix"

	configs, err := newCommonConfigs(prefix)
	assert.NoError(t, err)

	assert.Contains(t, configs, "ch:ceph-csi")
}
