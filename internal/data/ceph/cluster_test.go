// file: cluster_test.go
package ceph

import (
	"context"
	"testing"

	"github.com/otterscale/otterscale/internal/core"
)

/*
=============================================================

	Helper – 建立 *zero‑value* Ceph (不會真的連線)

=============================================================
*/
func mustNewCeph(t testing.TB) *Ceph {
	// 零值 Ceph，connection 會因缺少資訊而失敗，
	// 正好用來測試「錯誤路徑」。
	return &Ceph{}
}

/*
=============================================================

	Helper – 期待 error 或 panic

=============================================================
*/
func expectErrorOrPanic(t *testing.T, fn func() error) {
	defer func() {
		if r := recover(); r != nil {
			// panic 也是接受的行為
			t.Logf("recovered panic as expected: %v", r)
		}
	}()

	if err := fn(); err != nil {
		// 回傳非 nil error 也是接受的行為
		return
	}
	// 兩者皆未出現 → 測試失敗
	t.Errorf("expected an error or panic, got nil")
}

/*
=============================================================

	Mock 資料（完全使用 types.go 中宣告的匿名結構，
	保留所有 json tag，以確保型別相容）

=============================================================
*/
var mockData = struct {
	monDump *monDump
	monStat *monStat
	osdDump *osdDump
	osdTree *osdTree
	osdDF   *osdDF
	pgDump  *pgDump
	dfAll   *df
}{
	/*-------------------------- MON --------------------------*/
	monDump: func() *monDump {
		m := &monDump{}
		m.MONs = []struct {
			Name          string   `json:"name,omitempty"`
			Rank          uint64   `json:"rank,omitempty"`
			PublicAddress string   `json:"public_addr,omitempty"`
			Created       CephTime `json:"created,omitempty"`
		}{
			{
				Name:          "mon-a",
				Rank:          0,
				PublicAddress: "10.0.0.1:6789",
				Created:       CephTime{},
			},
			{
				Name:          "mon-b",
				Rank:          1,
				PublicAddress: "10.0.0.2:6789",
				Created:       CephTime{},
			},
		}
		return m
	}(),
	monStat: &monStat{Leader: "mon-a"},

	/*-------------------------- OSD Dump --------------------------*/
	osdDump: func() *osdDump {
		od := &osdDump{}
		od.OSDs = []struct {
			ID int64 `json:"osd,omitempty"`
			Up int   `json:"up,omitempty"`
			In int   `json:"in,omitempty"`
		}{
			{ID: 0, Up: 1, In: 1},
			{ID: 1, Up: 0, In: 1},
		}
		od.Pools = []struct {
			ID                   int64          `json:"pool,omitempty"`
			Name                 string         `json:"pool_name,omitempty"`
			Type                 int            `json:"type,omitempty"`
			Size                 uint64         `json:"size,omitempty"`
			PgNum                uint64         `json:"pg_num,omitempty"`
			PgNumTarget          uint64         `json:"pg_num_target,omitempty"`
			PgPlacementNum       uint64         `json:"pg_placement_num,omitempty"`
			PgPlacementNumTarget uint64         `json:"pg_placement_num_target,omitempty"`
			ApplicationMetadata  map[string]any `json:"application_metadata,omitempty"`
			CreateTime           CephTime       `json:"create_time,omitempty"`
		}{
			{
				ID:                   5,
				Name:                 "pool-data",
				PgNum:                128,
				PgPlacementNum:       0,
				PgNumTarget:          128,
				PgPlacementNumTarget: 0,
				Size:                 3,
				Type:                 1,
				CreateTime:           CephTime{},
				ApplicationMetadata:  map[string]any{"rbd": nil},
			},
		}
		return od
	}(),

	/*-------------------------- OSD Tree --------------------------*/
	osdTree: func() *osdTree {
		ot := &osdTree{}
		ot.Nodes = []struct {
			ID       int64   `json:"id,omitempty"`
			Name     string  `json:"name,omitempty"`
			Type     string  `json:"type,omitempty"`
			Exists   int     `json:"exists,omitempty"`
			Children []int64 `json:"children,omitempty"`
		}{
			{ID: 0, Type: "osd", Exists: 1},
			{ID: -1, Type: "host", Name: "host-01", Children: []int64{0}},
		}
		return ot
	}(),

	/*-------------------------- OSD DF --------------------------*/
	osdDF: func() *osdDF {
		df := &osdDF{}
		df.Nodes = []struct {
			ID          int64  `json:"id,omitempty"`
			DeviceClass string `json:"device_class,omitempty"`
			Name        string `json:"name,omitempty"`
			KB          uint64 `json:"kb,omitempty"`
			KBUsed      uint64 `json:"kb_used,omitempty"`
			PGCount     uint64 `json:"pgs,omitempty"`
		}{
			{
				ID:          0,
				Name:        "osd.0",
				DeviceClass: "hdd",
				KB:          1024 * 1024,
				KBUsed:      512 * 1024,
				PGCount:     10,
			},
		}
		return df
	}(),

	/*-------------------------- PG Dump --------------------------*/
	pgDump: func() *pgDump {
		pg := &pgDump{}
		pg.PGMap.PGStats = []struct {
			ID    string `json:"pgid,omitempty"`
			State string `json:"state,omitempty"`
		}{
			{ID: "5.1234", State: "active+clean"},
			{ID: "5.1235", State: "active+clean"},
		}
		return pg
	}(),

	/*-------------------------- DF All --------------------------*/
	dfAll: func() *df {
		d := &df{}
		d.Pools = []struct {
			Name  string `json:"name,omitempty"`
			ID    int64  `json:"id,omitempty"`
			Stats struct {
				Objects   uint64 `json:"objects,omitempty"`
				BytesUsed uint64 `json:"bytes_used,omitempty"`
			} `json:"stats,omitempty"`
		}{
			{
				Name: "pool-data",
				ID:   5,
				Stats: struct {
					Objects   uint64 `json:"objects,omitempty"`
					BytesUsed uint64 `json:"bytes_used,omitempty"`
				}{
					Objects:   1000,
					BytesUsed: 100 << 30,
				},
			},
		}
		return d
	}(),
}

/*
=============================================================

	1️⃣ Construction / interface compliance

=============================================================
*/
func TestNewCluster(t *testing.T) {
	ceph := mustNewCeph(t)
	repo := NewCluster(ceph)
	if repo == nil {
		t.Fatal("expected Cluster repository to be created, got nil")
	}
	// compile‑time interface check
	var _ core.CephClusterRepo = repo
}

/*
=============================================================

	2️⃣ Error path – connection failure (每個公開方法測試一次)

=============================================================
*/
func TestCluster_ListMONs_ConnectionError(t *testing.T) {
	repo := NewCluster(mustNewCeph(t))
	expectErrorOrPanic(t, func() error {
		_, err := repo.ListMONs(context.Background(), &core.StorageConfig{})
		return err
	})
}

func TestCluster_ListOSDs_ConnectionError(t *testing.T) {
	repo := NewCluster(mustNewCeph(t))
	expectErrorOrPanic(t, func() error {
		_, err := repo.ListOSDs(context.Background(), &core.StorageConfig{})
		return err
	})
}

func TestCluster_DoSMART_ConnectionError(t *testing.T) {
	repo := NewCluster(mustNewCeph(t))
	expectErrorOrPanic(t, func() error {
		_, err := repo.DoSMART(context.Background(), &core.StorageConfig{}, "osd.0")
		return err
	})
}

func TestCluster_ListPools_ConnectionError(t *testing.T) {
	repo := NewCluster(mustNewCeph(t))
	expectErrorOrPanic(t, func() error {
		_, err := repo.ListPools(context.Background(), &core.StorageConfig{})
		return err
	})
}

func TestCluster_ListPoolsByApplication_ConnectionError(t *testing.T) {
	repo := NewCluster(mustNewCeph(t))
	expectErrorOrPanic(t, func() error {
		_, err := repo.ListPoolsByApplication(context.Background(),
			&core.StorageConfig{}, "rbd")
		return err
	})
}
