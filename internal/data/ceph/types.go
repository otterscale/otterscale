package ceph

import "time"

type monDump struct {
	MONs []struct {
		Name          string `json:"name,omitempty"`
		Rank          uint64 `json:"rank,omitempty"`
		PublicAddress string `json:"public_addr,omitempty"`
	} `json:"mons,omitempty"`
}

type monStat struct {
	Leader string `json:"leader,omitempty"`
}

type osdDump struct {
	Pools []struct {
		ID                  int64          `json:"pool,omitempty"`
		Name                string         `json:"pool_name,omitempty"`
		Type                int            `json:"type,omitempty"`
		Size                uint64         `json:"size,omitempty"`
		PGCount             uint64         `json:"pg_num,omitempty"`
		ApplicationMetadata map[string]any `json:"application_metadata,omitempty"`
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
			Objects   uint64 `json:"objects,omitempty"`
			BytesUsed uint64 `json:"bytes_used,omitempty"`
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
	Filesystems []struct {
		Mdsmap struct {
			Epoch                     int       `json:"epoch,omitempty"`
			Flags                     int       `json:"flags,omitempty"`
			EverAllowedFeatures       int       `json:"ever_allowed_features,omitempty"`
			ExplicitlyAllowedFeatures int       `json:"explicitly_allowed_features,omitempty"`
			Created                   time.Time `json:"created,omitempty"`
			Modified                  time.Time `json:"modified,omitempty"`
			Tableserver               int       `json:"tableserver,omitempty"`
			Root                      int       `json:"root,omitempty"`
			SessionTimeout            int       `json:"session_timeout,omitempty"`
			SessionAutoclose          int       `json:"session_autoclose,omitempty"`
			MaxFileSize               int64     `json:"max_file_size,omitempty"`
			MaxXattrSize              int       `json:"max_xattr_size,omitempty"`
			LastFailure               int       `json:"last_failure,omitempty"`
			LastFailureOsdEpoch       int       `json:"last_failure_osd_epoch,omitempty"`
			MaxMds                    int       `json:"max_mds,omitempty"`
			In                        []int     `json:"in,omitempty"`
			Failed                    []any     `json:"failed,omitempty"`
			Damaged                   []any     `json:"damaged,omitempty"`
			Stopped                   []any     `json:"stopped,omitempty"`
			DataPools                 []int     `json:"data_pools,omitempty"`
			MetadataPool              int       `json:"metadata_pool,omitempty"`
			Enabled                   bool      `json:"enabled,omitempty"`
			FsName                    string    `json:"fs_name,omitempty"`
			Balancer                  string    `json:"balancer,omitempty"`
			BalRankMask               string    `json:"bal_rank_mask,omitempty"`
			StandbyCountWanted        int       `json:"standby_count_wanted,omitempty"`
			QdbLeader                 int       `json:"qdb_leader,omitempty"`
			QdbCluster                []int     `json:"qdb_cluster,omitempty"`
		} `json:"mdsmap,omitempty"`
		ID int `json:"id,omitempty"`
	} `json:"filesystems,omitempty"`
}

type names struct {
	Name string `json:"name,omitempty"`
}

type subvolumeInfo struct {
	UID        int       `json:"uid,omitempty"`
	GID        int       `json:"gid,omitempty"`
	Mode       int       `json:"mode,omitempty"`
	DataPool   string    `json:"data_pool,omitempty"`
	Features   []string  `json:"features,omitempty"`
	Flavor     int       `json:"flavor,omitempty"`
	Type       string    `json:"type,omitempty"`
	Path       string    `json:"path,omitempty"`
	State      string    `json:"state,omitempty"`
	BytesPcent string    `json:"bytes_pcent,omitempty"`
	BytesQuota any       `json:"bytes_quota,omitempty"`
	BytesUsed  uint64    `json:"bytes_used,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	Ctime      time.Time `json:"ctime,omitempty"`
	Mtime      time.Time `json:"mtime,omitempty"`
	Atime      time.Time `json:"atime,omitempty"`
}

type subvolumeSnapshotInfo struct {
	DataPool         string    `json:"data_pool,omitempty"`
	HasPendingClones string    `json:"has_pending_clones,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
}

type subvolumeGroupInfo struct {
	UID        int       `json:"uid,omitempty"`
	GID        int       `json:"gid,omitempty"`
	Mode       int       `json:"mode,omitempty"`
	DataPool   string    `json:"data_pool,omitempty"`
	BytesPcent string    `json:"bytes_pcent,omitempty"`
	BytesQuota any       `json:"bytes_quota,omitempty"`
	BytesUsed  uint64    `json:"bytes_used,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	Ctime      time.Time `json:"ctime,omitempty"`
	Mtime      time.Time `json:"mtime,omitempty"`
	Atime      time.Time `json:"atime,omitempty"`
}
