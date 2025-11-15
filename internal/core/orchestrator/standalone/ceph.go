package standalone

import (
	"strings"

	"github.com/otterscale/otterscale/internal/core/machine/tag"
)

type ceph struct {
	Scope      string
	OSDDevices []string
	NFSVIP     string
}

func newCeph(scope string, osdDevices []string, nfsVIP string) base {
	return &ceph{
		Scope:      scope,
		OSDDevices: osdDevices,
		NFSVIP:     nfsVIP,
	}
}

func (c *ceph) Charms() []charm {
	return []charm{
		{Name: "ch:ceph-mon", Channel: "squid/stable", PlacementScope: "lxd"},
		{Name: "ch:ceph-osd", Channel: "squid/stable", PlacementScope: "#"},
		{Name: "ch:ceph-fs", Channel: "squid/stable", PlacementScope: "lxd"},
		{Name: "ch:ceph-radosgw", Channel: "squid/stable", PlacementScope: "lxd"},
		{Name: "ch:ceph-nfs", Channel: "squid/stable", PlacementScope: "lxd"},
		{Name: "ch:hacluster", Channel: "2.8/stable", Subordinate: true},
	}
}

func (c *ceph) Configs() (string, error) {
	configs := map[string]map[string]any{
		"ceph-mon": {
			"monitor-count":       1,
			"expected-osd-count":  1,
			"config-flags":        `{ "global": {"osd_pool_default_size": 1, "osd_pool_default_min_size": 1, "mon_allow_pool_size_one": true} }`,
			"enable-perf-metrics": true,
		},
		"ceph-osd": {
			"osd-devices": strings.Join(c.OSDDevices, " "),
		},
		"ceph-fs": {
			"ceph-osd-replication-count": 2,
		},
		"ceph-nfs": {
			"vip": c.NFSVIP,
		},
		"ceph-radosgw": {
			"ceph-osd-replication-count": 2,
		},
		"hacluster": {
			"cluster_count": 1,
		},
	}

	return buildConfigs(c.Scope, configs)
}

func (c *ceph) Relations() [][]string {
	relations := [][]string{
		{"ceph-mon:client", "ceph-nfs:ceph-client"},
		{"ceph-mon:mds", "ceph-fs:ceph-mds"},
		{"ceph-mon:osd", "ceph-osd:mon"},
		{"ceph-mon:radosgw", "ceph-radosgw:mon"},
		{"hacluster:ha", "ceph-nfs:ha"},
	}

	return buildRelations(c.Scope, relations)
}

func (c *ceph) Tags() []string {
	return []string{
		tag.Ceph,
		tag.CephMON,
		tag.CephOSD,
	}
}
