package core

import (
	"context"
)

const charmCeph = "ceph-mon"

var (
	cephCharms = []EssentialCharm{
		{Name: "ch:ceph-mon", Channel: "squid/stable", LXD: true},
		{Name: "ch:ceph-osd", Channel: "squid/stable", Machine: true},
		{Name: "ch:ceph-fs", Channel: "squid/stable", LXD: true},
		{Name: "ch:ceph-radosgw", Channel: "squid/stable", LXD: true},
		{Name: "ch:ceph-nfs", Channel: "squid/stable", LXD: true},
		{Name: "ch:hacluster", Channel: "2.4/stable", Subordinate: true},
	}

	cephRelations = [][]string{
		{"ceph-mon:client", "ceph-nfs:ceph-client"},
		{"ceph-mon:mds", "ceph-fs:ceph-mds"},
		{"ceph-mon:osd", "ceph-osd:mon"},
		{"ceph-mon:radosgw", "ceph-radosgw:mon"},
		{"hacluster:ha", "ceph-nfs:ha"},
	}
)

func CreateCeph(ctx context.Context, serverRepo ServerRepo, machineRepo MachineRepo, facilityRepo FacilityRepo, tagRepo TagRepo, uuid, machineID, prefix string, configs map[string]string) error {
	if err := createEssential(ctx, serverRepo, machineRepo, facilityRepo, tagRepo, uuid, machineID, prefix, cephCharms, configs); err != nil {
		return err
	}
	return createEssentialRelations(ctx, facilityRepo, uuid, toEndpointList(prefix, cephRelations))
}

func GetCephCharms() []EssentialCharm {
	return cephCharms
}

func newCephConfigs(prefix, osdDevices, vip string) (map[string]string, error) {
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
			"ceph-osd-replication-count": 2,
		},
		"ceph-nfs": {
			"vip": vip,
		},
		"ceph-radosgw": {
			"ceph-osd-replication-count": 2,
		},
		"hacluster": {
			"cluster_count": 1,
		},
	}
	return NewCharmConfigs(prefix, configs)
}

func listCephs(ctx context.Context, scopeRepo ScopeRepo, clientRepo ClientRepo, uuid string) ([]Essential, error) {
	return listEssentials(ctx, scopeRepo, clientRepo, charmCeph, 2, uuid)
}
