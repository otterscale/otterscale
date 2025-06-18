package ceph

import (
	"context"
	"encoding/json"
	"slices"

	"github.com/ceph/go-ceph/rados"

	"github.com/openhdc/otterscale/internal/core"
)

type cluster struct {
	ceph *Ceph
}

func NewCluster(ceph *Ceph) core.CephClusterRepo {
	return &cluster{
		ceph: ceph,
	}
}

var _ core.CephClusterRepo = (*cluster)(nil)

func (r *cluster) ListMONs(ctx context.Context, config *core.StorageConfig) ([]core.MON, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	monDump, err := r.monDump(conn)
	if err != nil {
		return nil, err
	}
	return r.toMONs(monDump), nil
}

func (r *cluster) ListOSDs(ctx context.Context, config *core.StorageConfig) ([]core.OSD, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	osdDump, err := r.osdDump(conn)
	if err != nil {
		return nil, err
	}
	osdTree, err := r.osdTree(conn)
	if err != nil {
		return nil, err
	}
	return r.toOSDs(osdDump, osdTree), nil
}

func (r *cluster) ListPools(ctx context.Context, config *core.StorageConfig) ([]core.Pool, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	osdDump, err := r.osdDump(conn)
	if err != nil {
		return nil, err
	}
	return r.toPools(osdDump), nil
}

func (r *cluster) ListPoolsByApplication(ctx context.Context, config *core.StorageConfig, application string) ([]core.Pool, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	osdDump, err := r.osdDump(conn)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(r.toPools(osdDump), func(p core.Pool) bool {
		return !slices.Contains(p.Applications, application)
	}), nil
}

// func (r *pool) Create(ctx context.Context, config *core.StorageConfig, name string) (*core.Pool, error) {
// 	conn, err := r.ceph.connection(config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := conn.MakePool(name); err != nil {
// 		return nil, err
// 	}
// 	return &core.Pool{
// 		Name: name,
// 	}, nil
// }

// func (r *pool) Delete(ctx context.Context, config *core.StorageConfig, name string) error {
// 	conn, err := r.ceph.connection(config)
// 	if err != nil {
// 		return err
// 	}
// 	return conn.DeletePool(name)
// }

func (r *cluster) toMONs(d *monDump) []core.MON {
	ret := []core.MON{}
	for i := range d.Mons {
		ret = append(ret, core.MON{
			Name: d.Mons[i].Name,
		})
	}
	return ret
}

func (r *cluster) toOSDs(d *osdDump, t *osdTree) []core.OSD {
	ret := []core.OSD{}
	for i := range d.Osds {
		osd := core.OSD{}
		for j := range t.Nodes {
			if t.Nodes[j].Type == "osd" && d.Osds[i].Osd == t.Nodes[j].ID {
				osd.Name = t.Nodes[j].Name
				osd.DeviceClass = t.Nodes[j].DeviceClass
			}
		}
		ret = append(ret, osd)
	}
	return ret
}

func (r *cluster) toPools(d *osdDump) []core.Pool {
	ret := []core.Pool{}
	for i := range d.Pools {
		apps := []string{}
		for app := range d.Pools[i].ApplicationMetadata {
			apps = append(apps, app)
		}
		ret = append(ret, core.Pool{
			Name:         d.Pools[i].PoolName,
			Applications: apps,
		})
	}
	return ret
}

func (r *cluster) monDump(conn *rados.Conn) (*monDump, error) {
	cmd, err := json.Marshal(map[string]string{
		"prefix": "mon dump",
		"format": "json",
	})
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, err
	}
	var monDump monDump
	if err := json.Unmarshal(resp, &monDump); err != nil {
		return nil, err
	}
	return &monDump, nil
}

func (r *cluster) osdDump(conn *rados.Conn) (*osdDump, error) {
	cmd, err := json.Marshal(map[string]string{
		"prefix": "osd dump",
		"format": "json",
	})
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, err
	}
	var osdDump osdDump
	if err := json.Unmarshal(resp, &osdDump); err != nil {
		return nil, err
	}
	return &osdDump, nil
}

func (r *cluster) osdTree(conn *rados.Conn) (*osdTree, error) {
	cmd, err := json.Marshal(map[string]string{
		"prefix": "osd tree",
		"format": "json",
	})
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, err
	}
	var osdTree osdTree
	if err := json.Unmarshal(resp, &osdTree); err != nil {
		return nil, err
	}
	return &osdTree, nil
}
