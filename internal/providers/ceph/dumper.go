package ceph

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ceph/go-ceph/rados"

	"github.com/otterscale/otterscale/internal/core/storage"
)

type cephTime struct {
	time.Time
}

func (ct *cephTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == `""` || str == "null" {
		return nil
	}

	t, err := time.Parse(`2006-01-02T15:04:05.000000-0700`, strings.ReplaceAll(str, `"`, ``))
	if err != nil {
		return err
	}

	ct.Time = t

	return nil
}

type cephSubvolumeTime struct {
	time.Time
}

func (ct *cephSubvolumeTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == `""` || str == "null" {
		return nil
	}

	t, err := time.Parse(`"2006-01-02 15:04:05"`, str)
	if err != nil {
		return err
	}

	ct.Time = t

	return nil
}

type monDump struct {
	MONs []struct {
		Name          string   `json:"name,omitempty"`
		Rank          uint64   `json:"rank,omitempty"`
		PublicAddress string   `json:"public_addr,omitempty"`
		Created       cephTime `json:"created,omitempty"`
	} `json:"mons,omitempty"`
}

type monStat struct {
	Leader string `json:"leader,omitempty"`
}

type osdDump struct {
	Pools []struct {
		ID                   int64          `json:"pool,omitempty"`
		Name                 string         `json:"pool_name,omitempty"`
		Type                 int            `json:"type,omitempty"`
		Size                 uint64         `json:"size,omitempty"`
		PgNum                uint64         `json:"pg_num,omitempty"`
		PgNumTarget          uint64         `json:"pg_num_target,omitempty"`
		PgPlacementNum       uint64         `json:"pg_placement_num,omitempty"`
		PgPlacementNumTarget uint64         `json:"pg_placement_num_target,omitempty"`
		ApplicationMetadata  map[string]any `json:"application_metadata,omitempty"`
		CreateTime           cephTime       `json:"create_time,omitempty"`
	} `json:"pools,omitempty"`
	OSDs []struct {
		ID int64 `json:"osd,omitempty"`
		Up int   `json:"up,omitempty"`
		In int   `json:"in,omitempty"`
	} `json:"osds,omitempty"`
}

type osdTree struct {
	Nodes []struct {
		ID       int64   `json:"id,omitempty"`
		Name     string  `json:"name,omitempty"`
		Type     string  `json:"type,omitempty"`
		Exists   int     `json:"exists,omitempty"`
		Children []int64 `json:"children,omitempty"`
	} `json:"nodes,omitempty"`
}

type osdDF struct {
	Nodes []struct {
		ID          int64  `json:"id,omitempty"`
		DeviceClass string `json:"device_class,omitempty"`
		Name        string `json:"name,omitempty"`
		KB          uint64 `json:"kb,omitempty"`
		KBUsed      uint64 `json:"kb_used,omitempty"`
		PGCount     uint64 `json:"pgs,omitempty"`
	} `json:"nodes,omitempty"`
}

type df struct {
	Pools []struct {
		Name  string `json:"name,omitempty"`
		ID    int64  `json:"id,omitempty"`
		Stats struct {
			UsedObjects uint64 `json:"objects,omitempty"`
			UsedBytes   uint64 `json:"bytes_used,omitempty"`
			MaxBytes    uint64 `json:"max_avail,omitempty"`
		} `json:"stats,omitempty"`
	} `json:"pools,omitempty"`
}

type poolQuota struct {
	QuotaMaxObjects uint64 `json:"quota_max_objects,omitempty"`
	QuotaMaxBytes   uint64 `json:"quota_max_bytes,omitempty"`
}

type ecProfile struct {
	K string `json:"k,omitempty"`
	M string `json:"m,omitempty"`
}

type osdSMART struct {
	Smartctl struct {
		Output []string `json:"output,omitempty"`
	} `json:"smartctl,omitempty"`
}

type pgDump struct {
	PGMap struct {
		PGStats []struct {
			ID    string `json:"pgid,omitempty"`
			State string `json:"state,omitempty"`
		} `json:"pg_stats,omitempty"`
	} `json:"pg_map,omitempty"`
}

type device struct {
	ID       string `json:"devid,omitempty"`
	Location []struct {
		Host string `json:"host,omitempty"`
		Dev  string `json:"dev,omitempty"`
		Path string `json:"path,omitempty"`
	} `json:"location,omitempty"`
	Daemons []string `json:"daemons,omitempty"`
}

type fsDump struct {
	FileSystems []struct {
		MDSMap struct {
			FileSystemName string   `json:"fs_name,omitempty"`
			Created        cephTime `json:"created,omitempty"`
		} `json:"mdsmap,omitempty"`
		ID int64 `json:"id,omitempty"`
	} `json:"filesystems,omitempty"`
}

type names struct {
	Name string `json:"name,omitempty"`
}

type Bytes uint64
type UnixMode uint32

type SubvolumeState string
type Feature string

type subvolumeInfo struct {
	Path  string         `json:"path,omitempty"`
	State SubvolumeState `json:"state,omitempty"`
	UID   uint32         `json:"uid,omitempty"`
	GID   uint32         `json:"gid,omitempty"`
	Mode  UnixMode       `json:"mode,omitempty"`

	BytesPercent string `json:"bytes_percent,omitempty"`
	BytesUsed    Bytes  `json:"bytes_used,omitempty"`
	BytesQuota   *Bytes `json:"bytes_quota,omitempty"` // nil = unlimited

	DataPool      string `json:"data_pool,omitempty"`
	PoolNamespace string `json:"pool_namespace,omitempty"`

	Atime     cephSubvolumeTime `json:"atime,omitempty"`
	Mtime     cephSubvolumeTime `json:"mtime,omitempty"`
	Ctime     cephSubvolumeTime `json:"ctime,omitempty"`
	CreatedAt cephSubvolumeTime `json:"created_at,omitempty"`

	Features []Feature `json:"features,omitempty"`
}

type AccessType string
type Squash string
type SecType string

const (
	NoneSquash       Squash = "None"
	RootSquash       Squash = "Root"
	AllSquash        Squash = "All"
	RootIDSquash     Squash = "RootId"
	NoRootSquash            = NoneSquash
	Unspecifiedquash Squash = ""
)

// SecType indicates the kind of security/authentication to be used by an export.

const (
	SysSec   SecType = "sys"
	NoneSec  SecType = "none"
	Krb5Sec  SecType = "krb5"
	Krb5iSec SecType = "krb5i"
	Krb5pSec SecType = "krb5p"
)

// cephNFSFSALInfo describes NFS-Ganesha specific FSAL properties of an export.
type cephNFSFSALInfo struct {
	Name           string `json:"name,omitempty"`
	UserID         string `json:"user_id,omitempty"`
	FileSystemName string `json:"fs_name,omitempty"`
}

type cephNFSClientInfo struct {
	Addresses  []string `json:"addresses,omitempty"`
	AccessType string   `json:"access_type,omitempty"`
	Squash     string   `json:"squash,omitempty"`
}

type cephNFSExport struct {
	ExportID      uint64              `json:"export_id,omitempty"`
	Path          string              `json:"path,omitempty"`
	ClusterID     string              `json:"cluster_id,omitempty"`
	PseudoPath    string              `json:"pseudo,omitempty"`
	AccessType    string              `json:"access_type,omitempty"`
	Squash        string              `json:"squash,omitempty"`
	SecurityLabel bool                `json:"security_label,omitempty"`
	Protocols     []int               `json:"protocols,omitempty"`
	Transports    []string            `json:"transports,omitempty"`
	FSAL          cephNFSFSALInfo     `json:"fsal,omitempty"`
	Clients       []cephNFSClientInfo `json:"clients,omitempty"`
	SecType       []string            `json:"sectype,omitempty"`
}

type subvolumeSnapshotInfo struct {
	DataPool         string            `json:"data_pool,omitempty"`
	HasPendingClones string            `json:"has_pending_clones,omitempty"`
	CreatedAt        cephSubvolumeTime `json:"created_at,omitempty"`
}

type subvolumeGroupInfo struct {
	DataPool   string            `json:"data_pool,omitempty"`
	Mode       int               `json:"mode,omitempty"`
	BytesQuota any               `json:"bytes_quota,omitempty"`
	BytesUsed  uint64            `json:"bytes_used,omitempty"`
	CreatedAt  cephSubvolumeTime `json:"created_at,omitempty"`
}

func osdCommand(conn *rados.Conn, osd int, m any) ([]byte, error) {
	cmd, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	resp, status, err := conn.OsdCommand(osd, [][]byte{cmd})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", status, err)
	}

	return resp, nil
}

func osdCommandWithKeyUnmarshal(conn *rados.Conn, osd int, key string, m, v any) error {
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

	resp, status, err := conn.MonCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", status, err)
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

func monCommandWithKeyUnmarshal(conn *rados.Conn, key string, m, v any) error {
	resp, err := monCommand(conn, m)
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

func statMon(conn *rados.Conn) (*monStat, error) {
	cmd := map[string]string{
		"prefix": "mon stat",
		"format": "json",
	}

	var monStat monStat

	if err := monCommandWithUnmarshal(conn, cmd, &monStat); err != nil {
		return nil, err
	}

	return &monStat, nil
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

func dfOSD(conn *rados.Conn) (*osdDF, error) {
	cmd := map[string]string{
		"prefix": "osd df",
		"format": "json",
	}

	var osdDF osdDF

	if err := monCommandWithUnmarshal(conn, cmd, &osdDF); err != nil {
		return nil, err
	}

	return &osdDF, nil
}

func dfAll(conn *rados.Conn) (*df, error) {
	cmd := map[string]string{
		"prefix": "df",
		"format": "json",
	}

	var df df

	if err := monCommandWithUnmarshal(conn, cmd, &df); err != nil {
		return nil, err
	}

	return &df, nil
}

func dumpPG(conn *rados.Conn) (*pgDump, error) {
	cmd := map[string]string{
		"prefix": "pg dump",
		"format": "json",
	}

	var pgDump pgDump

	if err := monCommandWithUnmarshal(conn, cmd, &pgDump); err != nil {
		return nil, err
	}

	return &pgDump, nil
}

func createOSDPool(conn *rados.Conn, pool string, poolType storage.PoolType) error {
	cmd := map[string]string{
		"prefix":    "osd pool create",
		"pool":      pool,
		"pool_type": poolType.String(),
		"format":    "json",
	}

	if poolType == storage.PoolTypeErasure {
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

func getOSDPool(conn *rados.Conn, pool, key string) (string, error) {
	cmd := map[string]string{
		"prefix": "osd pool get",
		"pool":   pool,
		"var":    key,
		"format": "json",
	}

	var v any

	if err := monCommandWithKeyUnmarshal(conn, key, cmd, &v); err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", v), nil
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

func getOSDPoolQuota(conn *rados.Conn, pool string) (*poolQuota, error) {
	cmd := map[string]string{
		"prefix": "osd pool get-quota",
		"pool":   pool,
		"format": "json",
	}

	var poolQuota poolQuota

	if err := monCommandWithUnmarshal(conn, cmd, &poolQuota); err != nil {
		return nil, err
	}

	return &poolQuota, nil
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

func getECProfile(conn *rados.Conn, name string) (*ecProfile, error) {
	cmd := map[string]string{
		"prefix": "osd erasure-code-profile get",
		"name":   name,
		"format": "json",
	}

	var ecProfile ecProfile

	if err := monCommandWithUnmarshal(conn, cmd, &ecProfile); err != nil {
		return nil, err
	}

	return &ecProfile, nil
}

func smartOSD(conn *rados.Conn, osd int, deviceID string) (*osdSMART, error) {
	cmd := map[string]string{
		"prefix": "smart",
		"devid":  deviceID,
		"format": "json",
	}

	var osdSMART osdSMART

	if err := osdCommandWithKeyUnmarshal(conn, osd, deviceID, cmd, &osdSMART); err != nil {
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

func listSubvolumes(conn *rados.Conn, volume, group string) ([]names, error) {
	cmd := map[string]string{
		"prefix":   "fs subvolume ls",
		"vol_name": volume,
		"format":   "json",
	}

	if group != "" {
		cmd["group_name"] = group
	}

	var names []names

	if err := monCommandWithUnmarshal(conn, cmd, &names); err != nil {
		return nil, err
	}

	return names, nil
}

func getSubvolume(conn *rados.Conn, volume, group, subvolume string) (*subvolumeInfo, error) {
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

func createSubvolume(conn *rados.Conn, volume, group, subvolume string, size *uint64, uid, gid, mode *uint32, poolLayout *string, isNamespaceIsolated *bool) error {
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
	if size != nil {
		cmd["size"] = *size
	}
	if uid != nil {
		cmd["uid"] = *uid
	}
	if gid != nil {
		cmd["gid"] = *gid
	}
	if mode != nil {
		cmd["mode"] = *mode
	}
	if poolLayout != nil && *poolLayout != "" {
		cmd["pool_layout"] = *poolLayout
	}
	if isNamespaceIsolated != nil {
		cmd["namespace_isolated"] = *isNamespaceIsolated
	}

	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}

	return nil
}

func resizeSubvolume(conn *rados.Conn, volume, group, subvolume string, newSize uint64, noShrink *bool) error {
	cmd := map[string]any{
		"prefix":   "fs subvolume resize",
		"vol_name": volume,
		"sub_name": subvolume,
		"new_size": newSize,
		"format":   "json",
	}

	if group != "" {
		cmd["group_name"] = group
	}

	if noShrink != nil {
		cmd["no_shrink"] = *noShrink
	}

	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}

	return nil
}

func removeSubvolume(conn *rados.Conn, volume, group, subvolume string, isForce *bool) error {
	cmd := map[string]any{
		"prefix":   "fs subvolume rm",
		"vol_name": volume,
		"sub_name": subvolume,
		"format":   "json",
	}

	if group != "" {
		cmd["group_name"] = group
	}

	if isForce != nil {
		cmd["force"] = *isForce
	}

	if _, err := monCommand(conn, cmd); err != nil {
		return err
	}

	return nil
}

func listSubvolumeSnapshots(conn *rados.Conn, volume, subvolume, group string) ([]names, error) {
	cmd := map[string]string{
		"prefix":   "fs subvolume snapshot ls",
		"vol_name": volume,
		"sub_name": subvolume,
		"format":   "json",
	}

	if group != "" {
		cmd["group_name"] = group
	}

	var names []names

	if err := monCommandWithUnmarshal(conn, cmd, &names); err != nil {
		return nil, err
	}

	return names, nil
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

func listSubvolumeGroups(conn *rados.Conn, volume string) ([]names, error) {
	cmd := map[string]string{
		"prefix":   "fs subvolumegroup ls",
		"vol_name": volume,
		"format":   "json",
	}

	var names []names

	if err := monCommandWithUnmarshal(conn, cmd, &names); err != nil {
		return nil, err
	}

	return names, nil
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

func mgrCommand(conn *rados.Conn, m any) ([]byte, error) {
	cmd, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	resp, status, err := conn.MgrCommand([][]byte{cmd})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", status, err)
	}
	return resp, nil
}

func mgrCommandWithUnmarshal(conn *rados.Conn, m, v any) error {
	resp, err := mgrCommand(conn, m)
	if err != nil {
		return err
	}
	return json.Unmarshal(resp, &v)
}

func mgrCommandWithKeyUnmarshal(conn *rados.Conn, key string, m, v any) error {
	resp, err := mgrCommand(conn, m)
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

func mgrCommandWithInbuf(conn *rados.Conn, m any, inbuf []byte) ([]byte, error) {
	cmd, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	resp, status, err := conn.MgrCommand([][]byte{cmd, inbuf})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", status, err)
	}
	return resp, nil
}

func mgrCommandWithInbufUnmarshal(conn *rados.Conn, m any, inbuf []byte, v any) error {
	resp, err := mgrCommandWithInbuf(conn, m, inbuf)
	if err != nil {
		return err
	}
	return json.Unmarshal(resp, &v)
}

func listDetailedNFSExports(conn *rados.Conn, clusterID string) ([]cephNFSExport, error) {
	cmd := map[string]any{
		"prefix":     "nfs export ls",
		"detailed":   "true",
		"format":     "json",
		"cluster_id": clusterID,
	}

	var l []cephNFSExport
	if err := mgrCommandWithUnmarshal(conn, cmd, &l); err != nil {
		return nil, err
	}
	return l, nil
}

func applyNFSExport(conn *rados.Conn, clusterID string, spec *cephNFSExport) error {
	inbuf, err := json.Marshal(spec)
	if err != nil {
		return err
	}

	cmd := map[string]any{
		"prefix":     "nfs export apply",
		"cluster_id": clusterID,
		"format":     "json",
	}

	_, err = mgrCommandWithInbuf(conn, cmd, inbuf)
	return err
}

func getNFSExport(conn *rados.Conn, clusterID, pseudoPath string) (*cephNFSExport, error) {
	cmd := map[string]string{
		"prefix":      "nfs export info",
		"cluster_id":  clusterID,
		"pseudo_path": pseudoPath,
		"format":      "json",
	}

	var info cephNFSExport

	if err := mgrCommandWithUnmarshal(conn, cmd, &info); err != nil {
		cmd["pseudo"] = pseudoPath
		delete(cmd, "pseudo_path")
		if err2 := mgrCommandWithUnmarshal(conn, cmd, &info); err2 != nil {
			return nil, err
		}
	}
	return &info, nil
}

func removeNFSExport(conn *rados.Conn, clusterID, pseudoPath string) error {
	cmd := map[string]any{
		"prefix":      "nfs export rm",
		"cluster_id":  clusterID,
		"format":      "json",
		"pseudo_path": pseudoPath,
	}

	if _, err := mgrCommand(conn, cmd); err != nil {
		cmd["pseudo"] = pseudoPath
		delete(cmd, "pseudo_path")
		_, err2 := mgrCommand(conn, cmd)
		return err2
	}
	return nil
}
