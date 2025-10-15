package core

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type MON struct {
	Leader        bool
	Name          string
	Rank          uint64
	PublicAddress string

	MachineID       string
	MachineHostname string
}

type OSD struct {
	ID          int64
	Name        string
	Up          bool
	In          bool
	Exists      bool
	DeviceClass string
	Size        uint64
	Used        uint64
	PGCount     uint64
	Hostname    string

	MachineID       string
	MachineHostname string
}

type Pool struct {
	ID                  int64
	Name                string
	Updating            bool
	Type                string
	ECOverwrites        bool
	DataChunks          uint64
	CodingChunks        uint64
	ReplicatedSize      uint64
	QuotaBytes          uint64
	QuotaObjects        uint64
	UsedBytes           uint64
	UsedObjects         uint64
	MaxBytes            uint64
	PlacementGroupCount uint64
	PlacementGroupState map[string]int64
	CreatedAt           time.Time
	Applications        []string
}

type CephClusterRepo interface {
	ListMONs(ctx context.Context, config *CephConfig) ([]MON, error)
	ListOSDs(ctx context.Context, config *CephConfig) ([]OSD, error)
	DoSMART(ctx context.Context, config *CephConfig, who string) (map[string][]string, error)
	ListPools(ctx context.Context, config *CephConfig) ([]Pool, error)
	ListPoolsByApplication(ctx context.Context, config *CephConfig, application string) ([]Pool, error)
	CreatePool(ctx context.Context, config *CephConfig, pool, poolType string) error
	DeletePool(ctx context.Context, config *CephConfig, pool string) error
	EnableApplication(ctx context.Context, config *CephConfig, pool, application string) error
	GetParameter(ctx context.Context, config *CephConfig, pool, key string) (string, error)
	SetParameter(ctx context.Context, config *CephConfig, pool, key, value string) error
	SetQuota(ctx context.Context, config *CephConfig, pool string, maxBytes, maxObjects uint64) error
	GetQuota(ctx context.Context, config *CephConfig, pool string) (maxBytes, maxObjects uint64, err error)
	GetECProfile(ctx context.Context, config *CephConfig, name string) (k, m string, err error)
}

func (uc *StorageUseCase) ListMONs(ctx context.Context, scope, facility string) ([]MON, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	mons, err := uc.cephCluster.ListMONs(ctx, config)
	if err != nil {
		return nil, err
	}
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range mons {
		machine, _ := uc.getMachineByJujuMachine(machines, scope, mons[i].Name)
		if machine != nil {
			mons[i].MachineID = machine.SystemID
			mons[i].MachineHostname = machine.Hostname
		}
	}
	return mons, nil
}

func (uc *StorageUseCase) ListOSDs(ctx context.Context, scope, facility string) ([]OSD, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	osds, err := uc.cephCluster.ListOSDs(ctx, config)
	if err != nil {
		return nil, err
	}
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range osds {
		machine, _ := uc.getMachineByHostname(machines, osds[i].Hostname)
		if machine != nil {
			osds[i].MachineID = machine.SystemID
			osds[i].MachineHostname = machine.Hostname
		}
	}
	return osds, nil
}

func (uc *StorageUseCase) DoSMART(ctx context.Context, scope, facility, osd string) (map[string][]string, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	return uc.cephCluster.DoSMART(ctx, config, osd)
}

func (uc *StorageUseCase) ListPools(ctx context.Context, scope, facility, application string) ([]Pool, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	if application != "" {
		return uc.cephCluster.ListPoolsByApplication(ctx, config, application)
	}
	pools, err := uc.cephCluster.ListPools(ctx, config)
	if err != nil {
		return nil, err
	}
	for i := range pools {
		if pools[i].Type == "erasure" {
			ecOverwrites, _ := uc.cephCluster.GetParameter(ctx, config, pools[i].Name, "allow_ec_overwrites")
			if ecOverwrites == "true" {
				pools[i].ECOverwrites = true
			}
			ecProfile, _ := uc.cephCluster.GetParameter(ctx, config, pools[i].Name, "erasure_code_profile")
			if ecProfile != "" {
				k, m, _ := uc.cephCluster.GetECProfile(ctx, config, ecProfile)
				pools[i].DataChunks, _ = strconv.ParseUint(k, 10, 64)
				pools[i].CodingChunks, _ = strconv.ParseUint(m, 10, 64)
			}
		}
		maxBytes, maxObjects, _ := uc.cephCluster.GetQuota(ctx, config, pools[i].Name)
		pools[i].QuotaBytes = maxBytes
		pools[i].QuotaObjects = maxObjects
	}
	return pools, nil
}

func (uc *StorageUseCase) CreatePool(ctx context.Context, scope, facility, pool, poolType string, ecOverwrites bool, replicatedSize, quotaMaxBytes, quotaMaxObjects uint64, applications []string) (*Pool, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.cephCluster.CreatePool(ctx, config, pool, poolType); err != nil {
		return nil, err
	}
	if poolType == "erasure" && ecOverwrites {
		if err := uc.cephCluster.SetParameter(ctx, config, pool, "allow_ec_overwrites", "true"); err != nil {
			return nil, err
		}
	}
	if poolType == "replicated" && replicatedSize > 1 {
		if err := uc.cephCluster.SetParameter(ctx, config, pool, "size", strconv.FormatUint(replicatedSize, 10)); err != nil {
			return nil, err
		}
	}
	for _, app := range applications {
		if err := uc.cephCluster.EnableApplication(ctx, config, pool, app); err != nil {
			return nil, err
		}
	}
	if err := uc.cephCluster.SetQuota(ctx, config, pool, quotaMaxBytes, quotaMaxObjects); err != nil {
		return nil, err
	}
	return &Pool{
		Name: pool,
	}, nil
}

func (uc *StorageUseCase) UpdatePool(ctx context.Context, scope, facility, pool string, quotaMaxBytes, quotaMaxObjects uint64) (*Pool, error) {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}
	if err := uc.cephCluster.SetQuota(ctx, config, pool, quotaMaxBytes, quotaMaxObjects); err != nil {
		return nil, err
	}
	return &Pool{
		Name: pool,
	}, nil
}

func (uc *StorageUseCase) DeletePool(ctx context.Context, scope, facility, pool string) error {
	config, err := cephConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return err
	}
	return uc.cephCluster.DeletePool(ctx, config, pool)
}

func (uc *StorageUseCase) getMachineByJujuMachine(machines []Machine, scope, monName string) (*Machine, error) {
	for i := range machines {
		model, err := getJujuModelName(machines[i].WorkloadAnnotations)
		if err != nil {
			continue
		}
		if model != scope {
			continue
		}
		// machineID: 8
		machineID, err := getJujuMachineID(machines[i].WorkloadAnnotations)
		if err != nil {
			continue
		}
		// monName: juju-1eb21e-8-lxd-0
		monMachineIDTokens := strings.Split(monName, "-")
		if len(monMachineIDTokens) < 2 {
			continue
		}
		if machineID != monMachineIDTokens[2] {
			continue
		}
		return &machines[i], nil
	}
	return nil, fmt.Errorf("machine with mon %q not found", monName)
}

func (uc *StorageUseCase) getMachineByHostname(machines []Machine, hostname string) (*Machine, error) {
	for i := range machines {
		if machines[i].Hostname == hostname {
			return &machines[i], nil
		}
	}
	return nil, fmt.Errorf("machine with hostname %q not found", hostname)
}
