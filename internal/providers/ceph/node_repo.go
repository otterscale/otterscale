package ceph

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"github.com/otterscale/otterscale/internal/core/storage"
)

// Note: Ceph API do not support context.
type nodeRepo struct {
	ceph *Ceph
}

func NewNodeRepo(ceph *Ceph) storage.NodeRepo {
	return &nodeRepo{
		ceph: ceph,
	}
}

var _ storage.NodeRepo = (*nodeRepo)(nil)

func (r *nodeRepo) ListMonitors(_ context.Context, scope string) ([]storage.Monitor, error) {
	conn, err := r.ceph.connection(scope)
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

	return r.toMonitors(monDump, monStat), nil
}

func (r *nodeRepo) ListObjectStorageDaemons(_ context.Context, scope string) ([]storage.ObjectStorageDaemon, error) {
	conn, err := r.ceph.connection(scope)
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

func (r *nodeRepo) DoSMART(_ context.Context, scope, who string) (map[string][]string, error) {
	conn, err := r.ceph.connection(scope)
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

func (r *nodeRepo) Config(scope string) (host, id, key string, err error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return "", "", "", err
	}

	host, err = conn.GetConfigOption("mon_host")
	if err != nil {
		return "", "", "", err
	}

	id, err = conn.GetConfigOption("fsid")
	if err != nil {
		return "", "", "", err
	}

	key, err = conn.GetConfigOption("key")
	if err != nil {
		return "", "", "", err
	}

	return host, id, key, nil
}

func (r *nodeRepo) toMonitors(d *monDump, s *monStat) []storage.Monitor {
	ret := []storage.Monitor{}

	for i := range d.MONs {
		ret = append(ret, storage.Monitor{
			Leader:        d.MONs[i].Name == s.Leader,
			Name:          d.MONs[i].Name,
			Rank:          d.MONs[i].Rank,
			PublicAddress: d.MONs[i].PublicAddress,
		})
	}

	return ret
}

func (r *nodeRepo) toOSDs(d *osdDump, t *osdTree, df *osdDF) []storage.ObjectStorageDaemon {
	ret := []storage.ObjectStorageDaemon{}

	for i := range df.Nodes {
		osd := storage.ObjectStorageDaemon{
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

func (r *nodeRepo) getOSDNumber(osd string) int {
	token := strings.Split(osd, ".")

	if len(token) > 1 {
		number, _ := strconv.Atoi(token[1])
		return number
	}

	return 0
}
