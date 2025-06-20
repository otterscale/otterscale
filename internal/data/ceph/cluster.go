package ceph

import (
	"context"
	"encoding/json"
	"slices"
	"strconv"
	"strings"

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

func (r *cluster) DoSMART(ctx context.Context, config *core.StorageConfig, who string) (map[string][]string, error) {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return nil, err
	}
	deviceListByDaemon, err := r.deviceListByDaemon(conn, who)
	if err != nil {
		return nil, err
	}
	osd := r.getOSDNumber(who)
	outputs := map[string][]string{}
	for _, device := range deviceListByDaemon {
		resp, err := r.osdSMART(conn, osd, device.Devid)
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

func (r *cluster) CreatePool(ctx context.Context, config *core.StorageConfig, poolName, poolType string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return r.osdPoolCreate(conn, poolName, poolType)
}

func (r *cluster) DeletePool(ctx context.Context, config *core.StorageConfig, poolName string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return r.osdPoolDelete(conn, poolName)
}

func (r *cluster) EnableApplication(ctx context.Context, config *core.StorageConfig, poolName, application string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return r.osdPoolApplicationEnable(conn, poolName, application)
}

func (r *cluster) SetParameter(ctx context.Context, config *core.StorageConfig, poolName, key, value string) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	return r.osdPoolSet(conn, poolName, key, value)
}

func (r *cluster) SetQuota(ctx context.Context, config *core.StorageConfig, poolName string, maxBytes, maxObjects int) error {
	conn, err := r.ceph.connection(config)
	if err != nil {
		return err
	}
	if err := r.osdPoolSetQuota(conn, poolName, "max_bytes", maxBytes); err != nil {
		return err
	}
	return r.osdPoolSetQuota(conn, poolName, "max_objects", maxObjects)
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

func (r *cluster) deviceListByDaemon(conn *rados.Conn, who string) ([]deviceListByDaemon, error) {
	cmd, err := json.Marshal(map[string]string{
		"prefix": "device ls-by-daemon",
		"who":    who,
		"format": "json",
	})
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, err
	}
	var deviceListByDaemon []deviceListByDaemon
	if err := json.Unmarshal(resp, &deviceListByDaemon); err != nil {
		return nil, err
	}
	return deviceListByDaemon, nil
}

func (r *cluster) osdSMART(conn *rados.Conn, osd int, deviceID string) (*osdSMART, error) {
	cmd, err := json.Marshal(map[string]string{
		"prefix": "smart",
		"devid":  deviceID,
		"format": "json",
	})
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.OsdCommand(osd, [][]byte{cmd})
	if err != nil {
		return nil, err
	}
	var deviceList map[string]any
	if err := json.Unmarshal(resp, &deviceList); err != nil {
		return nil, err
	}
	device, err := json.Marshal(deviceList[deviceID])
	if err != nil {
		return nil, err
	}
	var osdSMART osdSMART
	if err := json.Unmarshal(device, &osdSMART); err != nil {
		return nil, err
	}
	return &osdSMART, nil
}

func (r *cluster) osdPoolCreate(conn *rados.Conn, poolName, poolType string) error {
	m := map[string]string{
		"prefix":    "osd pool create",
		"pool":      poolName,
		"pool_type": poolType,
		"format":    "json",
	}
	if poolType == "erasure" {
		m["erasure_code_profile"] = "default"
	}
	cmd, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, _, err := conn.MonCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (r *cluster) osdPoolDelete(conn *rados.Conn, poolName string) error {
	m := map[string]any{
		"prefix":                      "osd pool delete",
		"pool":                        poolName,
		"pool2":                       poolName,
		"yes_i_really_really_mean_it": true,
		"format":                      "json",
	}
	cmd, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, _, err := conn.MonCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (r *cluster) osdPoolApplicationEnable(conn *rados.Conn, poolName, app string) error {
	cmd, err := json.Marshal(map[string]any{
		"prefix":               "osd pool application enable",
		"pool":                 poolName,
		"app":                  app,
		"yes_i_really_mean_it": true,
		"format":               "json",
	})
	if err != nil {
		return err
	}
	if _, _, err := conn.MonCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (r *cluster) osdPoolSet(conn *rados.Conn, poolName, key, value string) error {
	cmd, err := json.Marshal(map[string]string{
		"prefix": "osd pool set",
		"pool":   poolName,
		"var":    key,
		"val":    value,
		"format": "json",
	})
	if err != nil {
		return err
	}
	if _, _, err := conn.MonCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (r *cluster) osdPoolSetQuota(conn *rados.Conn, poolName, field string, value int) error {
	cmd, err := json.Marshal(map[string]string{
		"prefix": "osd pool set-quota",
		"pool":   poolName,
		"field":  field,
		"val":    strconv.Itoa(value),
		"format": "json",
	})
	if err != nil {
		return err
	}
	if _, _, err := conn.MonCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (r *cluster) getOSDNumber(osd string) int {
	token := strings.Split(osd, ".")
	if len(token) > 1 {
		number, _ := strconv.Atoi(token[1])
		return number
	}
	return 0
}
