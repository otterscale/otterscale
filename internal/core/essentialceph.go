package core

import (
	"context"
)

const charmCeph = "ceph-mon"

var (
	cephCharms = []EssentialCharm{
		{Name: "ch:ceph-fs", LXD: true},
		{Name: "ch:ceph-mon", LXD: true},
		{Name: "ch:ceph-osd", LXD: false},
		{Name: "ch:ceph-dashboard", LXD: false},
	}

	cephRelations = [][]string{
		{"ceph-fs:ceph-mds", "ceph-mon:mds"},
		{"ceph-osd:mon", "ceph-mon:osd"},
		{"ceph-dashboard:dashboard", "ceph-mon:dashboard"},
		{"easyrsa:client", "ceph-dashboard:certificates"},
	}
)

func CreateCeph(ctx context.Context, serverRepo ServerRepo, machineRepo MachineRepo, facilityRepo FacilityRepo, uuid, machineID, prefix string, configs map[string]string) error {
	if err := createEssential(ctx, serverRepo, machineRepo, facilityRepo, uuid, machineID, prefix, cephCharms, configs); err != nil {
		return err
	}
	return createEssentialRelations(ctx, facilityRepo, uuid, toEndpointList(prefix, cephRelations))
}

func GetCephCharms() []EssentialCharm {
	return cephCharms
}

func newCephConfigs(prefix, osdDevices string) (map[string]string, error) {
	configs := map[string]map[string]any{
		"ceph-mon": {
			"monitor-count":      1,
			"expected-osd-count": 1,
			"config-flags":       `{ "global": {"osd_pool_default_size": 1, "osd_pool_default_min_size": 1, "mon_allow_pool_size_one": true} }`,
		},
		"ceph-osd": {
			"osd-devices": osdDevices,
		},
		"ceph-fs": {
			"ceph-osd-replication-count": 1,
		},
	}
	return NewCharmConfigs(prefix, configs)
}

func listCephs(ctx context.Context, scopeRepo ScopeRepo, clientRepo ClientRepo, uuid string) ([]Essential, error) {
	return listEssentials(ctx, scopeRepo, clientRepo, charmCeph, 2, uuid)
}
