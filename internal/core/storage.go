package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/juju/juju/api/client/action"
	"golang.org/x/sync/errgroup"
	"gopkg.in/ini.v1"
)

const (
	cephConfigCommand        = "ceph config generate-minimal-conf && ceph auth get client.admin"
	cephRGWUserListCommand   = "radosgw-admin user list"
	cephRGWUserCreateCommand = "radosgw-admin user create --system --uid=otterscale --display-name=OtterScale --format json"
	cephRGWUserInfoCommand   = "radosgw-admin user info --uid=otterscale --format=json"
)

type StorageCephConfig struct {
	FSID    string
	MONHost string
	Key     string
}

type StorageRGWConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

type StorageConfig struct {
	*StorageCephConfig
	*StorageRGWConfig
}

type StorageUseCase struct {
	action   ActionRepo
	facility FacilityRepo
	cluster  CephClusterRepo
	rbd      CephRBDRepo
	fs       CephFSRepo
	rgw      CephRGWRepo
	machine  MachineRepo

	configs sync.Map
}

func NewStorageUseCase(action ActionRepo, facility FacilityRepo, cluster CephClusterRepo, rbd CephRBDRepo, fs CephFSRepo, rgw CephRGWRepo, machine MachineRepo) *StorageUseCase {
	return &StorageUseCase{
		action:   action,
		facility: facility,
		cluster:  cluster,
		rbd:      rbd,
		fs:       fs,
		rgw:      rgw,
		machine:  machine,
	}
}

func (uc *StorageUseCase) config(ctx context.Context, uuid, name string) (*StorageConfig, error) {
	key := uuid + "/" + name

	if v, ok := uc.configs.Load(key); ok {
		return v.(*StorageConfig), nil
	}

	config, err := uc.newConfig(ctx, uuid, name)
	if err != nil {
		return nil, err
	}

	uc.configs.Store(key, config)

	return config, nil
}

func (uc *StorageUseCase) newConfig(ctx context.Context, uuid, name string) (*StorageConfig, error) {
	var (
		leader   string
		endpoint string
	)
	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		monLeader, err := uc.facility.GetLeader(egctx, uuid, name) // ceph-mon
		if err != nil {
			return err
		}
		leader = monLeader
		return nil
	})
	eg.Go(func() error {
		leader, err := uc.facility.GetLeader(egctx, uuid, rgwName(name))
		if err != nil {
			return err
		}
		info, err := uc.facility.GetUnitInfo(egctx, uuid, leader)
		if err != nil {
			return err
		}
		endpoint = fmt.Sprintf("http://%s", info.PublicAddress)
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	var cephConfig *StorageCephConfig
	rgwConfig := &StorageRGWConfig{
		Endpoint: endpoint,
	}

	eg, egctx = errgroup.WithContext(ctx)
	eg.Go(func() error {
		result, err := uc.runCommand(egctx, uuid, leader, cephConfigCommand)
		if err != nil {
			return err
		}
		config, err := uc.extractStorageCephConfig(result)
		if err != nil {
			return err
		}
		cephConfig = config
		return nil
	})
	eg.Go(func() error {
		listResult, err := uc.runCommand(egctx, uuid, leader, cephRGWUserListCommand)
		if err != nil {
			return err
		}
		cmd, err := uc.getRGWCommand(listResult)
		if err != nil {
			return err
		}
		result, err := uc.runCommand(egctx, uuid, leader, cmd)
		if err != nil {
			return err
		}
		config, err := uc.extractStorageRGWConfig(result)
		if err != nil {
			return err
		}
		rgwConfig.AccessKey = config.AccessKey
		rgwConfig.SecretKey = config.SecretKey
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return &StorageConfig{
		StorageCephConfig: cephConfig,
		StorageRGWConfig:  rgwConfig,
	}, nil
}

func (uc *StorageUseCase) runCommand(ctx context.Context, uuid, leader, command string) (*action.ActionResult, error) {
	id, err := uc.action.RunCommand(ctx, uuid, leader, command)
	if err != nil {
		return nil, err
	}
	return waitForActionCompleted(ctx, uc.action, uuid, id, time.Second, time.Minute)
}

func (uc *StorageUseCase) runAction(ctx context.Context, uuid, leader, action string, params map[string]any) error {
	id, err := uc.action.RunAction(ctx, uuid, leader, action, params)
	if err != nil {
		return err
	}
	if _, err := waitForActionCompleted(ctx, uc.action, uuid, id, time.Second, time.Minute); err != nil {
		return err
	}
	return nil
}

func (uc *StorageUseCase) extractStorageCephConfig(result *action.ActionResult) (*StorageCephConfig, error) {
	stdout, ok := result.Output["stdout"]
	if !ok {
		return nil, errors.New("ceph config stdout not found")
	}
	file, err := ini.Load([]byte(stdout.(string)))
	if err != nil {
		return nil, err
	}
	fsID := file.Section("global").Key("fsid").String()
	if fsID == "" {
		return nil, errors.New("ceph config fsid not found")
	}
	monHost := file.Section("global").Key("mon_host").String()
	if monHost == "" {
		return nil, errors.New("ceph config mon_host not found")
	}
	key := file.Section("client.admin").Key("key").String()
	if key == "" {
		return nil, errors.New("ceph config key not found")
	}
	return &StorageCephConfig{
		FSID:    fsID,
		MONHost: monHost,
		Key:     key,
	}, nil
}

func (uc *StorageUseCase) getRGWCommand(result *action.ActionResult) (string, error) {
	stdout, ok := result.Output["stdout"]
	if !ok {
		return "", errors.New("rgw list config stdout not found")
	}
	var users []string
	if err := json.Unmarshal([]byte(stdout.(string)), &users); err != nil {
		return "", err
	}
	cmd := cephRGWUserCreateCommand
	if slices.Contains(users, "otterscale") {
		cmd = cephRGWUserInfoCommand
	}
	return cmd, nil
}

func (uc *StorageUseCase) extractStorageRGWConfig(result *action.ActionResult) (*StorageRGWConfig, error) {
	stdout, ok := result.Output["stdout"]
	if !ok {
		return nil, errors.New("rgw config stdout not found")
	}
	type Info struct {
		Keys []struct {
			AccessKey string `json:"access_key,omitempty"`
			SecretKey string `json:"secret_key,omitempty"`
		} `json:"keys,omitempty"`
	}
	var info Info
	if err := json.Unmarshal([]byte(stdout.(string)), &info); err != nil {
		return nil, err
	}
	if len(info.Keys) == 0 {
		return nil, errors.New("rgw config key not found")
	}
	accessKey := info.Keys[0].AccessKey
	if accessKey == "" {
		return nil, errors.New("rgw config access key not found")
	}
	secretKey := info.Keys[0].SecretKey
	if secretKey == "" {
		return nil, errors.New("rgw config secret key not found")
	}
	return &StorageRGWConfig{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, nil
}

func rgwName(monName string) string {
	tokens := strings.Split(monName, "-")
	lastIndex := len(tokens) - 1
	tokens[lastIndex] = "radosgw"
	return strings.Join(tokens, "-")
}

func nfsName(monName string) string {
	tokens := strings.Split(monName, "-")
	lastIndex := len(tokens) - 1
	tokens[lastIndex] = "nfs"
	return strings.Join(tokens, "-")
}
