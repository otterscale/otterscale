package ceph

import (
	"context"
	"slices"
	"strconv"
	"strings"

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
	monDump, err := dumpMon(conn)
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
	osdDump, err := dumpOSD(conn)
	if err != nil {
		return nil, err
	}
	osdTree, err := treeOSD(conn)
	if err != nil {
		return nil, err
	}
	return r.toOSDs(osdDump, osdTree), nil
}

func (r *cluster) DoSMART(ctx context.Context, config *core.StorageConfig, who string) (map[string][]string, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	devices, err := listDevices(conn, who)
	if err != nil {
		return nil, err
	}
	osd := r.getOSDNumber(who)
	outputs := map[string][]string{}
	for _, device := range devices {
		resp, err := smartOSD(conn, osd, device.Devid)
		if err != nil {
			return nil, err
		}
		outputs[device.Devid] = resp.Smartctl.Output
	}
	return outputs, nil
}

func (r *cluster) ListPools(ctx context.Context, config *core.StorageConfig) ([]core.Pool, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	osdDump, err := dumpOSD(conn)
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
	osdDump, err := dumpOSD(conn)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(r.toPools(osdDump), func(p core.Pool) bool {
		return !slices.Contains(p.Applications, application)
	}), nil
}

func (r *cluster) CreatePool(ctx context.Context, config *core.StorageConfig, pool, poolType string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return createOSDPool(conn, pool, poolType)
}

func (r *cluster) DeletePool(ctx context.Context, config *core.StorageConfig, pool string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return deleteOSDPool(conn, pool)
}

func (r *cluster) EnableApplication(ctx context.Context, config *core.StorageConfig, pool, application string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return enableOSDPoolApplication(conn, pool, application)
}

func (r *cluster) SetParameter(ctx context.Context, config *core.StorageConfig, pool, key, value string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return setOSDPool(conn, pool, key, value)
}

func (r *cluster) SetQuota(ctx context.Context, config *core.StorageConfig, pool string, maxBytes, maxObjects uint64) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	if err := setOSDPoolQuota(conn, pool, "max_bytes", maxBytes); err != nil {
		return err
	}
	return setOSDPoolQuota(conn, pool, "max_objects", maxObjects)
}

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

func (r *cluster) getOSDNumber(osd string) int {
	token := strings.Split(osd, ".")
	if len(token) > 1 {
		number, _ := strconv.Atoi(token[1])
		return number
	}
	return 0
}
