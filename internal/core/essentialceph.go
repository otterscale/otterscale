package core

import (
	"context"
)

const (
	charmCeph    = "ceph-mon"
	charmCephCSI = "ceph-csi"
)

var (
	cephCharms = []essentialCharm{
		{name: "ch:ceph-fs", lxd: true},
		{name: "ch:ceph-mon", lxd: true},
		{name: "ch:ceph-osd", lxd: false},
	}

	cephRelations = [][]string{
		{"ceph-fs:ceph-mds", "ceph-mon:mds"},
		{"ceph-osd:mon", "ceph-mon:osd"},
	}

	cephCSICharms = []essentialCharm{
		{name: "ch:ceph-csi", lxd: true},
	}

	cephCSIRelations = [][]string{
		{"ceph-csi", "ceph-mon"},
		{"ceph-csi", "kubernetes-control-plane"},
	}
)

func CreateCeph(ctx context.Context, serverRepo ServerRepo, machineRepo MachineRepo, facilityRepo FacilityRepo, uuid, machineID, prefix string, configs map[string]string) error {
	if err := createEssential(ctx, serverRepo, machineRepo, facilityRepo, uuid, machineID, prefix, cephCharms, configs); err != nil {
		return err
	}
	return createEssentialRelations(ctx, facilityRepo, uuid, toEndpointList(prefix, cephRelations))
}

func CreateCephCSI(ctx context.Context, serverRepo ServerRepo, machineRepo MachineRepo, facilityRepo FacilityRepo, uuid, prefix string, configs map[string]string) error {
	if err := createEssential(ctx, serverRepo, machineRepo, facilityRepo, uuid, "", prefix, cephCSICharms, configs); err != nil {
		return err
	}
	return createEssentialRelations(ctx, facilityRepo, uuid, toEndpointList(prefix, cephCSIRelations))
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

func newCephCSIConfigs(prefix string) (map[string]string, error) {
	configs := map[string]map[string]any{
		"ceph-csi": {
			"default-storage":      "ceph-ext4",
			"provisioner-replicas": 1,
		},
	}
	return NewCharmConfigs(prefix, configs)
}

func listCephs(ctx context.Context, scopeRepo ScopeRepo, clientRepo ClientRepo, uuid string) ([]Essential, error) {
	return listEssentials(ctx, scopeRepo, clientRepo, charmCeph, 2, uuid)
}
