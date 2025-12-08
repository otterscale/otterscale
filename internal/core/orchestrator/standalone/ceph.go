package standalone

import (
	"encoding/json"
	"strings"

	"github.com/otterscale/otterscale/internal/core/machine/tag"
	"github.com/otterscale/otterscale/internal/core/versions"
)

type cephGlobalConfig struct {
	OsdPoolDefaultSize    int  `json:"osd_pool_default_size"`
	OsdPoolDefaultMinSize int  `json:"osd_pool_default_min_size"`
	MonAllowPoolSizeOne   bool `json:"mon_allow_pool_size_one"`
	MonAllowPoolDelete    bool `json:"mon_allow_pool_delete"`
}

type cephConfigFlags struct {
	Global cephGlobalConfig `json:"global"`
}

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
		{Name: "ch:ceph-mon", Channel: versions.Ceph, PlacementScope: "lxd"},
		{Name: "ch:ceph-osd", Channel: versions.Ceph, PlacementScope: "#"},
		{Name: "ch:ceph-fs", Channel: versions.Ceph, PlacementScope: "lxd"},
		{Name: "ch:ceph-radosgw", Channel: versions.Ceph, PlacementScope: "lxd"},
		{Name: "ch:ceph-nfs", Channel: versions.Ceph, PlacementScope: "lxd"},
		{Name: "ch:hacluster", Channel: versions.HACluster, Subordinate: true},
	}
}

func (c *ceph) Config(charmName string) (string, error) {
	configFlags, _ := json.Marshal(cephConfigFlags{
		Global: cephGlobalConfig{
			OsdPoolDefaultSize:    1,
			OsdPoolDefaultMinSize: 1,
			MonAllowPoolSizeOne:   true,
			MonAllowPoolDelete:    true,
		},
	})

	configs := map[string]map[string]any{
		"ceph-mon": {
			"config-flags":        string(configFlags),
			"enable-perf-metrics": true,
			"expected-osd-count":  1,
			"monitor-count":       1,
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

	return buildConfig(c.Scope, charmName, configs)
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
