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
	monStat, err := statMon(conn)
	if err != nil {
		return nil, err
	}
	return r.toMONs(monDump, monStat), nil
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
	osdDF, err := dfOSD(conn)
	if err != nil {
		return nil, err
	}
	return r.toOSDs(osdDump, osdTree, osdDF), nil
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
		resp, err := smartOSD(conn, osd, device.ID)
		if err != nil {
			return nil, err
		}
		outputs[device.ID] = resp.Smartctl.Output
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
	dumpPG, err := dumpPG(conn)
	if err != nil {
		return nil, err
	}
	dfAll, err := dfAll(conn)
	if err != nil {
		return nil, err
	}
	return r.toPools(osdDump, dumpPG, dfAll), nil
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
	dumpPG, err := dumpPG(conn)
	if err != nil {
		return nil, err
	}
	dfAll, err := dfAll(conn)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(r.toPools(osdDump, dumpPG, dfAll), func(p core.Pool) bool {
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

func (r *cluster) GetParameter(ctx context.Context, config *core.StorageConfig, pool, key string) (string, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return "", err
	}
	return getOSDPool(conn, pool, key)
}

func (r *cluster) SetParameter(ctx context.Context, config *core.StorageConfig, pool, key, value string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return setOSDPool(conn, pool, key, value)
}

func (r *cluster) GetQuota(ctx context.Context, config *core.StorageConfig, pool string) (maxBytes, maxObjects uint64, err error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return 0, 0, err
	}
	quota, err := getOSDPoolQuota(conn, pool)
	if err != nil {
		return 0, 0, err
	}
	return quota.QuotaMaxBytes, quota.QuotaMaxObjects, nil
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

func (r *cluster) GetECProfile(ctx context.Context, config *core.StorageConfig, name string) (k, m string, err error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return "", "", err
	}
	profile, err := getECProfile(conn, name)
	if err != nil {
		return "", "", err
	}
	return profile.K, profile.M, nil
}

func (r *cluster) toMONs(d *monDump, s *monStat) []core.MON {
	ret := []core.MON{}
	for i := range d.MONs {
		ret = append(ret, core.MON{
			Leader:        d.MONs[i].Name == s.Leader,
			Name:          d.MONs[i].Name,
			Rank:          d.MONs[i].Rank,
			PublicAddress: d.MONs[i].PublicAddress,
		})
	}
	return ret
}

func (r *cluster) toOSDs(d *osdDump, t *osdTree, df *osdDF) []core.OSD {
	ret := []core.OSD{}
	for i := range df.Nodes {
		osd := core.OSD{
			ID:          df.Nodes[i].ID,
			Name:        df.Nodes[i].Name,
			DeviceClass: df.Nodes[i].DeviceClass,
			Size:        df.Nodes[i].KB * 1024,
			Used:        df.Nodes[i].KBUsed * 1024,
			PGCount:     df.Nodes[i].PGCount,
		}
		for j := range t.Nodes {
			if t.Nodes[j].Type == "osd" && t.Nodes[j].ID == df.Nodes[i].ID && t.Nodes[j].Exists == 1 {
				osd.Exists = true
			}
			if t.Nodes[j].Type == "host" && slices.Contains(t.Nodes[j].Children, df.Nodes[i].ID) {
				osd.Hostname = t.Nodes[j].Name
			}
		}
		for j := range d.OSDs {
			if d.OSDs[j].ID != df.Nodes[i].ID {
				continue
			}
			if d.OSDs[j].Up == 1 {
				osd.Up = true
			}
			if d.OSDs[j].In == 1 {
				osd.In = true
			}
			break
		}
		ret = append(ret, osd)
	}
	return ret
}

func (r *cluster) toPools(d *osdDump, pd *pgDump, df *df) []core.Pool {
	ret := []core.Pool{}
	for i := range d.Pools {
		pool := core.Pool{
			ID:                  d.Pools[i].ID,
			Name:                d.Pools[i].Name,
			ReplicatedSize:      d.Pools[i].Size,
			PlacementGroupCount: d.Pools[i].PGCount,
			PlacementGroupState: map[string]int64{},
			CreatedAt:           d.Pools[i].CreateTime.Time,
		}
		switch d.Pools[i].Type {
		case 1:
			pool.Type = "replicated"
		case 3:
			pool.Type = "erasure"
		}
		for j := range df.Pools {
			if d.Pools[i].ID != df.Pools[j].ID {
				continue
			}
			pool.UsedBytes = df.Pools[j].Stats.BytesUsed
			pool.UsedObjects = df.Pools[j].Stats.Objects
		}
		for j := range pd.PGMap.PGStats {
			id := strings.Split(pd.PGMap.PGStats[j].ID, ".")[0]
			if strconv.FormatInt(d.Pools[i].ID, 10) != id {
				continue
			}
			state := pd.PGMap.PGStats[j].State
			pool.PlacementGroupState[state]++
		}
		for app := range d.Pools[i].ApplicationMetadata {
			pool.Applications = append(pool.Applications, app)
		}
		ret = append(ret, pool)
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
