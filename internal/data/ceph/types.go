package ceph

import "time"

type monDump struct {
	Epoch             int       `json:"epoch,omitempty"`
	Fsid              string    `json:"fsid,omitempty"`
	Modified          time.Time `json:"modified,omitempty"`
	Created           time.Time `json:"created,omitempty"`
	MinMonRelease     int       `json:"min_mon_release,omitempty"`
	MinMonReleaseName string    `json:"min_mon_release_name,omitempty"`
	ElectionStrategy  int       `json:"election_strategy,omitempty"`
	DisallowedLeaders string    `json:"disallowed_leaders: ,omitempty"`
	StretchMode       bool      `json:"stretch_mode,omitempty"`
	TiebreakerMon     string    `json:"tiebreaker_mon,omitempty"`
	RemovedRanks      string    `json:"removed_ranks: ,omitempty"`
	Features          struct {
		Persistent []string `json:"persistent,omitempty"`
		Optional   []any    `json:"optional,omitempty"`
	} `json:"features,omitempty"`
	Mons []struct {
		Rank        int    `json:"rank,omitempty"`
		Name        string `json:"name,omitempty"`
		PublicAddrs struct {
			Addrvec []struct {
				Type  string `json:"type,omitempty"`
				Addr  string `json:"addr,omitempty"`
				Nonce int    `json:"nonce,omitempty"`
			} `json:"addrvec,omitempty"`
		} `json:"public_addrs,omitempty"`
		Addr          string `json:"addr,omitempty"`
		PublicAddr    string `json:"public_addr,omitempty"`
		Priority      int    `json:"priority,omitempty"`
		Weight        int    `json:"weight,omitempty"`
		CrushLocation string `json:"crush_location,omitempty"`
	} `json:"mons,omitempty"`
	Quorum []int `json:"quorum,omitempty"`
}

type osdDump struct {
	Epoch                  int      `json:"epoch,omitempty"`
	Fsid                   string   `json:"fsid,omitempty"`
	Created                string   `json:"created,omitempty"`
	Modified               string   `json:"modified,omitempty"`
	LastUpChange           string   `json:"last_up_change,omitempty"`
	LastInChange           string   `json:"last_in_change,omitempty"`
	Flags                  string   `json:"flags,omitempty"`
	FlagsNum               int      `json:"flags_num,omitempty"`
	FlagsSet               []string `json:"flags_set,omitempty"`
	CrushVersion           int      `json:"crush_version,omitempty"`
	FullRatio              float64  `json:"full_ratio,omitempty"`
	BackfillfullRatio      float64  `json:"backfillfull_ratio,omitempty"`
	NearfullRatio          float64  `json:"nearfull_ratio,omitempty"`
	ClusterSnapshot        string   `json:"cluster_snapshot,omitempty"`
	PoolMax                int      `json:"pool_max,omitempty"`
	MaxOsd                 int      `json:"max_osd,omitempty"`
	RequireMinCompatClient string   `json:"require_min_compat_client,omitempty"`
	MinCompatClient        string   `json:"min_compat_client,omitempty"`
	RequireOsdRelease      string   `json:"require_osd_release,omitempty"`
	AllowCrimson           bool     `json:"allow_crimson,omitempty"`
	Pools                  []struct {
		Pool                              int    `json:"pool,omitempty"`
		PoolName                          string `json:"pool_name,omitempty"`
		CreateTime                        string `json:"create_time,omitempty"`
		Flags                             int    `json:"flags,omitempty"`
		FlagsNames                        string `json:"flags_names,omitempty"`
		Type                              int    `json:"type,omitempty"`
		Size                              int    `json:"size,omitempty"`
		MinSize                           int    `json:"min_size,omitempty"`
		CrushRule                         int    `json:"crush_rule,omitempty"`
		PeeringCrushBucketCount           int    `json:"peering_crush_bucket_count,omitempty"`
		PeeringCrushBucketTarget          int    `json:"peering_crush_bucket_target,omitempty"`
		PeeringCrushBucketBarrier         int    `json:"peering_crush_bucket_barrier,omitempty"`
		PeeringCrushBucketMandatoryMember int64  `json:"peering_crush_bucket_mandatory_member,omitempty"`
		ObjectHash                        int    `json:"object_hash,omitempty"`
		PgAutoscaleMode                   string `json:"pg_autoscale_mode,omitempty"`
		PgNum                             int    `json:"pg_num,omitempty"`
		PgPlacementNum                    int    `json:"pg_placement_num,omitempty"`
		PgPlacementNumTarget              int    `json:"pg_placement_num_target,omitempty"`
		PgNumTarget                       int    `json:"pg_num_target,omitempty"`
		PgNumPending                      int    `json:"pg_num_pending,omitempty"`
		LastPgMergeMeta                   struct {
			SourcePgid       string `json:"source_pgid,omitempty"`
			ReadyEpoch       int    `json:"ready_epoch,omitempty"`
			LastEpochStarted int    `json:"last_epoch_started,omitempty"`
			LastEpochClean   int    `json:"last_epoch_clean,omitempty"`
			SourceVersion    string `json:"source_version,omitempty"`
			TargetVersion    string `json:"target_version,omitempty"`
		} `json:"last_pg_merge_meta,omitempty"`
		LastChange                     string `json:"last_change,omitempty"`
		LastForceOpResend              string `json:"last_force_op_resend,omitempty"`
		LastForceOpResendPrenautilus   string `json:"last_force_op_resend_prenautilus,omitempty"`
		LastForceOpResendPreluminous   string `json:"last_force_op_resend_preluminous,omitempty"`
		Auid                           int    `json:"auid,omitempty"`
		SnapMode                       string `json:"snap_mode,omitempty"`
		SnapSeq                        int    `json:"snap_seq,omitempty"`
		SnapEpoch                      int    `json:"snap_epoch,omitempty"`
		PoolSnaps                      []any  `json:"pool_snaps,omitempty"`
		RemovedSnaps                   string `json:"removed_snaps,omitempty"`
		QuotaMaxBytes                  int    `json:"quota_max_bytes,omitempty"`
		QuotaMaxObjects                int    `json:"quota_max_objects,omitempty"`
		Tiers                          []any  `json:"tiers,omitempty"`
		TierOf                         int    `json:"tier_of,omitempty"`
		ReadTier                       int    `json:"read_tier,omitempty"`
		WriteTier                      int    `json:"write_tier,omitempty"`
		CacheMode                      string `json:"cache_mode,omitempty"`
		TargetMaxBytes                 int    `json:"target_max_bytes,omitempty"`
		TargetMaxObjects               int    `json:"target_max_objects,omitempty"`
		CacheTargetDirtyRatioMicro     int    `json:"cache_target_dirty_ratio_micro,omitempty"`
		CacheTargetDirtyHighRatioMicro int    `json:"cache_target_dirty_high_ratio_micro,omitempty"`
		CacheTargetFullRatioMicro      int    `json:"cache_target_full_ratio_micro,omitempty"`
		CacheMinFlushAge               int    `json:"cache_min_flush_age,omitempty"`
		CacheMinEvictAge               int    `json:"cache_min_evict_age,omitempty"`
		ErasureCodeProfile             string `json:"erasure_code_profile,omitempty"`
		HitSetParams                   struct {
			Type string `json:"type,omitempty"`
		} `json:"hit_set_params,omitempty"`
		HitSetPeriod              int            `json:"hit_set_period,omitempty"`
		HitSetCount               int            `json:"hit_set_count,omitempty"`
		UseGmtHitset              bool           `json:"use_gmt_hitset,omitempty"`
		MinReadRecencyForPromote  int            `json:"min_read_recency_for_promote,omitempty"`
		MinWriteRecencyForPromote int            `json:"min_write_recency_for_promote,omitempty"`
		HitSetGradeDecayRate      int            `json:"hit_set_grade_decay_rate,omitempty"`
		HitSetSearchLastN         int            `json:"hit_set_search_last_n,omitempty"`
		GradeTable                []any          `json:"grade_table,omitempty"`
		StripeWidth               int            `json:"stripe_width,omitempty"`
		ExpectedNumObjects        int            `json:"expected_num_objects,omitempty"`
		FastRead                  bool           `json:"fast_read,omitempty"`
		Options                   any            `json:"options,omitempty"`
		ApplicationMetadata       map[string]any `json:"application_metadata,omitempty"`
		ReadBalance               struct {
			ScoreType                      string `json:"score_type,omitempty"`
			ScoreActing                    int    `json:"score_acting,omitempty"`
			ScoreStable                    int    `json:"score_stable,omitempty"`
			OptimalScore                   int    `json:"optimal_score,omitempty"`
			RawScoreActing                 int    `json:"raw_score_acting,omitempty"`
			RawScoreStable                 int    `json:"raw_score_stable,omitempty"`
			PrimaryAffinityWeighted        int    `json:"primary_affinity_weighted,omitempty"`
			AveragePrimaryAffinity         int    `json:"average_primary_affinity,omitempty"`
			AveragePrimaryAffinityWeighted int    `json:"average_primary_affinity_weighted,omitempty"`
		} `json:"read_balance,omitempty"`
	} `json:"pools,omitempty"`
	Osds []struct {
		Osd             int    `json:"osd,omitempty"`
		UUID            string `json:"uuid,omitempty"`
		Up              int    `json:"up,omitempty"`
		In              int    `json:"in,omitempty"`
		Weight          int    `json:"weight,omitempty"`
		PrimaryAffinity int    `json:"primary_affinity,omitempty"`
		LastCleanBegin  int    `json:"last_clean_begin,omitempty"`
		LastCleanEnd    int    `json:"last_clean_end,omitempty"`
		UpFrom          int    `json:"up_from,omitempty"`
		UpThru          int    `json:"up_thru,omitempty"`
		DownAt          int    `json:"down_at,omitempty"`
		LostAt          int    `json:"lost_at,omitempty"`
		PublicAddrs     struct {
			Addrvec []struct {
				Type  string `json:"type,omitempty"`
				Addr  string `json:"addr,omitempty"`
				Nonce int    `json:"nonce,omitempty"`
			} `json:"addrvec,omitempty"`
		} `json:"public_addrs,omitempty"`
		ClusterAddrs struct {
			Addrvec []struct {
				Type  string `json:"type,omitempty"`
				Addr  string `json:"addr,omitempty"`
				Nonce int    `json:"nonce,omitempty"`
			} `json:"addrvec,omitempty"`
		} `json:"cluster_addrs,omitempty"`
		HeartbeatBackAddrs struct {
			Addrvec []struct {
				Type  string `json:"type,omitempty"`
				Addr  string `json:"addr,omitempty"`
				Nonce int    `json:"nonce,omitempty"`
			} `json:"addrvec,omitempty"`
		} `json:"heartbeat_back_addrs,omitempty"`
		HeartbeatFrontAddrs struct {
			Addrvec []struct {
				Type  string `json:"type,omitempty"`
				Addr  string `json:"addr,omitempty"`
				Nonce int    `json:"nonce,omitempty"`
			} `json:"addrvec,omitempty"`
		} `json:"heartbeat_front_addrs,omitempty"`
		PublicAddr         string   `json:"public_addr,omitempty"`
		ClusterAddr        string   `json:"cluster_addr,omitempty"`
		HeartbeatBackAddr  string   `json:"heartbeat_back_addr,omitempty"`
		HeartbeatFrontAddr string   `json:"heartbeat_front_addr,omitempty"`
		State              []string `json:"state,omitempty"`
	} `json:"osds,omitempty"`
	OsdXinfo []struct {
		Osd                  int    `json:"osd,omitempty"`
		DownStamp            string `json:"down_stamp,omitempty"`
		LaggyProbability     int    `json:"laggy_probability,omitempty"`
		LaggyInterval        int    `json:"laggy_interval,omitempty"`
		Features             int64  `json:"features,omitempty"`
		OldWeight            int    `json:"old_weight,omitempty"`
		LastPurgedSnapsScrub string `json:"last_purged_snaps_scrub,omitempty"`
		DeadEpoch            int    `json:"dead_epoch,omitempty"`
	} `json:"osd_xinfo,omitempty"`
}

type osdTree struct {
	Nodes []struct {
		ID          int    `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Type        string `json:"type,omitempty"`
		TypeID      int    `json:"type_id,omitempty"`
		Children    []int  `json:"children,omitempty"`
		PoolWeights struct {
		} `json:"pool_weights,omitempty"`
		DeviceClass     string  `json:"device_class,omitempty"`
		CrushWeight     float64 `json:"crush_weight,omitempty"`
		Depth           int     `json:"depth,omitempty"`
		Exists          int     `json:"exists,omitempty"`
		Status          string  `json:"status,omitempty"`
		Reweight        int     `json:"reweight,omitempty"`
		PrimaryAffinity int     `json:"primary_affinity,omitempty"`
	} `json:"nodes,omitempty"`
	Stray []any `json:"stray,omitempty"`
}

type osdSMART struct {
	Device struct {
		InfoName string `json:"info_name,omitempty"`
		Name     string `json:"name,omitempty"`
		Protocol string `json:"protocol,omitempty"`
		Type     string `json:"type,omitempty"`
	} `json:"device,omitempty"`
	DeviceType struct {
		Name      string `json:"name,omitempty"`
		ScsiValue int    `json:"scsi_value,omitempty"`
	} `json:"device_type,omitempty"`
	JSONFormatVersion []int `json:"json_format_version,omitempty"`
	LocalTime         struct {
		Asctime string `json:"asctime,omitempty"`
		TimeT   int    `json:"time_t,omitempty"`
	} `json:"local_time,omitempty"`
	LogicalBlockSize                          int    `json:"logical_block_size,omitempty"`
	ModelName                                 string `json:"model_name,omitempty"`
	NvmeSmartHealthInformationAddLogError     string `json:"nvme_smart_health_information_add_log_error,omitempty"`
	NvmeSmartHealthInformationAddLogErrorCode int    `json:"nvme_smart_health_information_add_log_error_code,omitempty"`
	NvmeVendor                                string `json:"nvme_vendor,omitempty"`
	Product                                   string `json:"product,omitempty"`
	Revision                                  string `json:"revision,omitempty"`
	RotationRate                              int    `json:"rotation_rate,omitempty"`
	ScsiVersion                               string `json:"scsi_version,omitempty"`
	Smartctl                                  struct {
		Argv         []string `json:"argv,omitempty"`
		BuildInfo    string   `json:"build_info,omitempty"`
		ExitStatus   int      `json:"exit_status,omitempty"`
		Output       []string `json:"output,omitempty"`
		PlatformInfo string   `json:"platform_info,omitempty"`
		SvnRevision  string   `json:"svn_revision,omitempty"`
		Version      []int    `json:"version,omitempty"`
	} `json:"smartctl,omitempty"`
	Temperature struct {
		Current   int `json:"current,omitempty"`
		DriveTrip int `json:"drive_trip,omitempty"`
	} `json:"temperature,omitempty"`
	UserCapacity struct {
		Blocks int   `json:"blocks,omitempty"`
		Bytes  int64 `json:"bytes,omitempty"`
	} `json:"user_capacity,omitempty"`
	Vendor string `json:"vendor,omitempty"`
}

type device struct {
	Devid    string `json:"devid,omitempty"`
	Location []struct {
		Host string `json:"host,omitempty"`
		Dev  string `json:"dev,omitempty"`
		Path string `json:"path,omitempty"`
	} `json:"location,omitempty"`
	Daemons []string `json:"daemons,omitempty"`
}

type fsDump struct {
	Epoch        int    `json:"epoch,omitempty"`
	Btime        string `json:"btime,omitempty"`
	DefaultFscid int    `json:"default_fscid,omitempty"`
	Filesystems  []struct {
		Mdsmap struct {
			Epoch                     int    `json:"epoch,omitempty"`
			Flags                     int    `json:"flags,omitempty"`
			EverAllowedFeatures       int    `json:"ever_allowed_features,omitempty"`
			ExplicitlyAllowedFeatures int    `json:"explicitly_allowed_features,omitempty"`
			Created                   string `json:"created,omitempty"`
			Modified                  string `json:"modified,omitempty"`
			Tableserver               int    `json:"tableserver,omitempty"`
			Root                      int    `json:"root,omitempty"`
			SessionTimeout            int    `json:"session_timeout,omitempty"`
			SessionAutoclose          int    `json:"session_autoclose,omitempty"`
			MaxFileSize               int64  `json:"max_file_size,omitempty"`
			MaxXattrSize              int    `json:"max_xattr_size,omitempty"`
			LastFailure               int    `json:"last_failure,omitempty"`
			LastFailureOsdEpoch       int    `json:"last_failure_osd_epoch,omitempty"`
			MaxMds                    int    `json:"max_mds,omitempty"`
			In                        []int  `json:"in,omitempty"`
			Failed                    []any  `json:"failed,omitempty"`
			Damaged                   []any  `json:"damaged,omitempty"`
			Stopped                   []any  `json:"stopped,omitempty"`
			DataPools                 []int  `json:"data_pools,omitempty"`
			MetadataPool              int    `json:"metadata_pool,omitempty"`
			Enabled                   bool   `json:"enabled,omitempty"`
			FsName                    string `json:"fs_name,omitempty"`
			Balancer                  string `json:"balancer,omitempty"`
			BalRankMask               string `json:"bal_rank_mask,omitempty"`
			StandbyCountWanted        int    `json:"standby_count_wanted,omitempty"`
			QdbLeader                 int    `json:"qdb_leader,omitempty"`
			QdbCluster                []int  `json:"qdb_cluster,omitempty"`
		} `json:"mdsmap,omitempty"`
		ID int `json:"id,omitempty"`
	} `json:"filesystems,omitempty"`
}

type subvolume struct {
	Name string `json:"name,omitempty"`
}

type subvolumeInfo struct {
	Atime         string   `json:"atime,omitempty"`
	BytesPcent    string   `json:"bytes_pcent,omitempty"`
	BytesQuota    any      `json:"bytes_quota,omitempty"`
	BytesUsed     any      `json:"bytes_used,omitempty"`
	CreatedAt     string   `json:"created_at,omitempty"`
	Ctime         string   `json:"ctime,omitempty"`
	DataPool      string   `json:"data_pool,omitempty"`
	Features      []string `json:"features,omitempty"`
	Flavor        int      `json:"flavor,omitempty"`
	Gid           int      `json:"gid,omitempty"`
	Mode          int      `json:"mode,omitempty"`
	MonAddrs      []string `json:"mon_addrs,omitempty"`
	Mtime         string   `json:"mtime,omitempty"`
	Path          string   `json:"path,omitempty"`
	PoolNamespace string   `json:"pool_namespace,omitempty"`
	State         string   `json:"state,omitempty"`
	Type          string   `json:"type,omitempty"`
	UID           int      `json:"uid,omitempty"`
}

type subvolumeSnapshotInfo struct {
	CreatedAt        string `json:"created_at,omitempty"`
	DataPool         string `json:"data_pool,omitempty"`
	HasPendingClones string `json:"has_pending_clones,omitempty"`
}

type subvolumeGroup struct {
	Name string `json:"name,omitempty"`
}

type subvolumeGroupInfo struct {
	Atime      string   `json:"atime,omitempty"`
	BytesPcent string   `json:"bytes_pcent,omitempty"`
	BytesQuota any      `json:"bytes_quota,omitempty"`
	BytesUsed  any      `json:"bytes_used,omitempty"`
	CreatedAt  string   `json:"created_at,omitempty"`
	Ctime      string   `json:"ctime,omitempty"`
	DataPool   string   `json:"data_pool,omitempty"`
	Gid        int      `json:"gid,omitempty"`
	Mode       int      `json:"mode,omitempty"`
	MonAddrs   []string `json:"mon_addrs,omitempty"`
	Mtime      string   `json:"mtime,omitempty"`
	UID        int      `json:"uid,omitempty"`
}
