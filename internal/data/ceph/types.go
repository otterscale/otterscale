package ceph

import (
	"strings"
	"time"
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

type subvolumeInfo struct {
	Path       string            `json:"path,omitempty"`
	DataPool   string            `json:"data_pool,omitempty"`
	Mode       int               `json:"mode,omitempty"`
	BytesQuota any               `json:"bytes_quota,omitempty"`
	BytesUsed  uint64            `json:"bytes_used,omitempty"`
	CreatedAt  cephSubvolumeTime `json:"created_at,omitempty"`
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
