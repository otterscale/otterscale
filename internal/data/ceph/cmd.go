package ceph

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/ceph/go-ceph/rados"
)

func osdCommand(conn *rados.Conn, osd int, m any) ([]byte, error) {
	cmd, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.OsdCommand(osd, [][]byte{cmd})
	if err != nil {
		return nil, errors.New(string(resp))
	}
	return resp, nil
}

func osdCommandWithUnmarshal(conn *rados.Conn, osd int, key string, m, v any) error {
	resp, err := osdCommand(conn, osd, m)
	if err != nil {
		return err
	}
	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return err
	}
	tmp, err := json.Marshal(result[key])
	if err != nil {
		return err
	}
	return json.Unmarshal(tmp, &v)
}

func monCommand(conn *rados.Conn, m any) ([]byte, error) {
	cmd, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	resp, _, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, errors.New(string(resp))
	}
	return resp, nil
}

func monCommandWithUnmarshal(conn *rados.Conn, m, v any) error {
	resp, err := monCommand(conn, m)
	if err != nil {
		return err
	}
	return json.Unmarshal(resp, &v)
}

func dumpMon(conn *rados.Conn) (*monDump, error) {
	cmd := map[string]string{
		"prefix": "mon dump",
		"format": "json",
	}
	var monDump monDump
	if err := monCommandWithUnmarshal(conn, cmd, &monDump); err != nil {
		return nil, err
	}
	return &monDump, nil
}

func dumpOSD(conn *rados.Conn) (*osdDump, error) {
	cmd := map[string]string{
		"prefix": "osd dump",
		"format": "json",
	}
	var osdDump osdDump
	if err := monCommandWithUnmarshal(conn, cmd, &osdDump); err != nil {
		return nil, err
	}
	return &osdDump, nil
}

func treeOSD(conn *rados.Conn) (*osdTree, error) {
	cmd := map[string]string{
		"prefix": "osd tree",
		"format": "json",
	}
	var osdTree osdTree
	if err := monCommandWithUnmarshal(conn, cmd, &osdTree); err != nil {
		return nil, err
	}
	return &osdTree, nil
}

func createOSDPool(conn *rados.Conn, pool, poolType string) error {
	cmd := map[string]string{
		"prefix":    "osd pool create",
		"pool":      pool,
		"pool_type": poolType,
		"format":    "json",
	}
	if poolType == "erasure" {
		cmd["erasure_code_profile"] = "default"
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func deleteOSDPool(conn *rados.Conn, pool string) error {
	cmd := map[string]any{
		"prefix":                      "osd pool delete",
		"pool":                        pool,
		"pool2":                       pool,
		"yes_i_really_really_mean_it": true,
		"format":                      "json",
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func enableOSDPoolApplication(conn *rados.Conn, pool, application string) error {
	cmd := map[string]any{
		"prefix":               "osd pool application enable",
		"pool":                 pool,
		"app":                  application,
		"yes_i_really_mean_it": true,
		"format":               "json",
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func setOSDPool(conn *rados.Conn, pool, key, value string) error {
	cmd := map[string]string{
		"prefix": "osd pool set",
		"pool":   pool,
		"var":    key,
		"val":    value,
		"format": "json",
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func setOSDPoolQuota(conn *rados.Conn, pool, field string, value uint64) error {
	cmd := map[string]string{
		"prefix": "osd pool set-quota",
		"pool":   pool,
		"field":  field,
		"val":    strconv.FormatUint(value, 10),
		"format": "json",
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func smartOSD(conn *rados.Conn, osd int, deviceID string) (*osdSMART, error) {
	cmd := map[string]string{
		"prefix": "smart",
		"devid":  deviceID,
		"format": "json",
	}
	var osdSMART osdSMART
	if err := osdCommandWithUnmarshal(conn, osd, deviceID, cmd, &osdSMART); err != nil {
		return nil, err
	}
	return &osdSMART, nil
}

func listDevices(conn *rados.Conn, who string) ([]device, error) {
	cmd := map[string]string{
		"prefix": "device ls-by-daemon",
		"who":    who,
		"format": "json",
	}
	var devices []device
	if err := monCommandWithUnmarshal(conn, cmd, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func dumpFS(conn *rados.Conn) (*fsDump, error) {
	cmd := map[string]string{
		"prefix": "fs dump",
		"format": "json",
	}
	var fsDump fsDump
	if err := monCommandWithUnmarshal(conn, cmd, &fsDump); err != nil {
		return nil, err
	}
	return &fsDump, nil
}

func listSubvolumes(conn *rados.Conn, volume, group string) ([]subvolume, error) {
	cmd := map[string]string{
		"prefix":   "fs subvolume ls",
		"vol_name": volume,
		"format":   "json",
	}
	if group != "" {
		cmd["group_name"] = group
	}
	var subvolumes []subvolume
	if err := monCommandWithUnmarshal(conn, cmd, &subvolumes); err != nil {
		return nil, err
	}
	return subvolumes, nil
}

func getSubvolume(conn *rados.Conn, volume, subvolume, group string) (*subvolumeInfo, error) {
	cmd := map[string]string{
		"prefix":   "fs subvolume info",
		"vol_name": volume,
		"sub_name": subvolume,
		"format":   "json",
	}
	if group != "" {
		cmd["group_name"] = group
	}
	var info subvolumeInfo
	if err := monCommandWithUnmarshal(conn, cmd, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func createSubvolume(conn *rados.Conn, volume, subvolume, group string, size uint64) error {
	cmd := map[string]any{
		"prefix":   "fs subvolume create",
		"vol_name": volume,
		"sub_name": subvolume,
		"size":     size,
		"format":   "json",
	}
	if group != "" {
		cmd["group_name"] = group
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func resizeSubvolume(conn *rados.Conn, volume, subvolume, group string, size uint64) error {
	cmd := map[string]any{
		"prefix":   "fs subvolume resize",
		"vol_name": volume,
		"sub_name": subvolume,
		"new_size": size,
		"format":   "json",
	}
	if group != "" {
		cmd["group_name"] = group
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func removeSubvolume(conn *rados.Conn, volume, subvolume, group string) error {
	cmd := map[string]any{
		"prefix":   "fs subvolume rm",
		"vol_name": volume,
		"sub_name": subvolume,
		"format":   "json",
	}
	if group != "" {
		cmd["group_name"] = group
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func getSubvolumeSnapshot(conn *rados.Conn, volume, subvolume, group, snapshot string) (*subvolumeSnapshotInfo, error) {
	cmd := map[string]string{
		"prefix":    "fs subvolume snapshot info",
		"vol_name":  volume,
		"sub_name":  subvolume,
		"snap_name": snapshot,
		"format":    "json",
	}
	if group != "" {
		cmd["group_name"] = group
	}
	var info subvolumeSnapshotInfo
	if err := monCommandWithUnmarshal(conn, cmd, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func createSubvolumeSnapshot(conn *rados.Conn, volume, subvolume, group, snapshot string) error {
	cmd := map[string]any{
		"prefix":    "fs subvolume snapshot create",
		"vol_name":  volume,
		"sub_name":  subvolume,
		"snap_name": snapshot,
		"format":    "json",
	}
	if group != "" {
		cmd["group_name"] = group
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func removeSubvolumeSnapshot(conn *rados.Conn, volume, subvolume, group, snapshot string) error {
	cmd := map[string]any{
		"prefix":    "fs subvolume snapshot rm",
		"vol_name":  volume,
		"sub_name":  subvolume,
		"snap_name": snapshot,
		"format":    "json",
	}
	if group != "" {
		cmd["group_name"] = group
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func listSubvolumeGroups(conn *rados.Conn, volume string) ([]subvolumeGroup, error) {
	cmd := map[string]string{
		"prefix":   "fs subvolumegroup ls",
		"vol_name": volume,
		"format":   "json",
	}
	var groups []subvolumeGroup
	if err := monCommandWithUnmarshal(conn, cmd, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func getSubvolumeGroup(conn *rados.Conn, volume, group string) (*subvolumeGroupInfo, error) {
	cmd := map[string]string{
		"prefix":     "fs subvolumegroup info",
		"vol_name":   volume,
		"group_name": group,
		"format":     "json",
	}
	var info subvolumeGroupInfo
	if err := monCommandWithUnmarshal(conn, cmd, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func createSubvolumeGroup(conn *rados.Conn, volume, group string, size uint64) error {
	cmd := map[string]any{
		"prefix":     "fs subvolumegroup create",
		"vol_name":   volume,
		"group_name": group,
		"size":       size,
		"format":     "json",
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func resizeSubvolumeGroup(conn *rados.Conn, volume, group string, size uint64) error {
	cmd := map[string]any{
		"prefix":     "fs subvolumegroup resize",
		"vol_name":   volume,
		"group_name": group,
		"new_size":   size,
		"format":     "json",
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}

func removeSubvolumeGroup(conn *rados.Conn, volume, group string) error {
	cmd := map[string]any{
		"prefix":     "fs subvolumegroup rm",
		"vol_name":   volume,
		"group_name": group,
		"format":     "json",
	}
	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}
	return nil
}
