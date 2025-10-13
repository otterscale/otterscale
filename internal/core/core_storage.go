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

var storageConfigMap sync.Map

func storageConfig(ctx context.Context, facility FacilityRepo, action ActionRepo, uuid, name string) (*StorageConfig, error) {
	key := uuid + "/" + name

	if v, ok := storageConfigMap.Load(key); ok {
		return v.(*StorageConfig), nil
	}

	config, err := newStorageConfig(ctx, facility, action, uuid, name)
	if err != nil {
		return nil, err
	}

	storageConfigMap.Store(key, config)

	return config, nil
}

func newStorageConfig(ctx context.Context, facility FacilityRepo, action ActionRepo, uuid, name string) (*StorageConfig, error) {
	var (
		leader   string
		endpoint string
	)
	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		monLeader, err := facility.GetLeader(egctx, uuid, name) // ceph-mon
		if err != nil {
			return err
		}
		leader = monLeader
		return nil
	})
	eg.Go(func() error {
		leader, err := facility.GetLeader(egctx, uuid, rgwName(name))
		if err != nil {
			return err
		}
		info, err := facility.GetUnitInfo(egctx, uuid, leader)
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
		result, err := runCommand(egctx, action, uuid, leader, "ceph config generate-minimal-conf && ceph auth get client.admin")
		if err != nil {
			return err
		}
		config, err := extractStorageCephConfig(result)
		if err != nil {
			return err
		}
		cephConfig = config
		return nil
	})
	eg.Go(func() error {
		listResult, err := runCommand(egctx, action, uuid, leader, "radosgw-admin user list")
		if err != nil {
			return err
		}
		cmd, err := getRGWCommand(listResult)
		if err != nil {
			return err
		}
		result, err := runCommand(egctx, action, uuid, leader, cmd)
		if err != nil {
			return err
		}
		config, err := extractStorageRGWConfig(result)
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

func extractStorageCephConfig(result *action.ActionResult) (*StorageCephConfig, error) {
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

func getRGWCommand(result *action.ActionResult) (string, error) {
	stdout, ok := result.Output["stdout"]
	if !ok {
		return "", errors.New("rgw list config stdout not found")
	}
	var users []string
	if err := json.Unmarshal([]byte(stdout.(string)), &users); err != nil {
		return "", err
	}
	cmd := "radosgw-admin user create --system --uid=otterscale --display-name=OtterScale --format json"
	if slices.Contains(users, "otterscale") {
		cmd = "radosgw-admin user info --uid=otterscale --format=json"
	}
	return cmd, nil
}

func extractStorageRGWConfig(result *action.ActionResult) (*StorageRGWConfig, error) {
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

func runCommand(ctx context.Context, actionRepo ActionRepo, uuid, leader, command string) (*action.ActionResult, error) {
	id, err := actionRepo.RunCommand(ctx, uuid, leader, command)
	if err != nil {
		return nil, err
	}
	return waitForActionCompleted(ctx, actionRepo, uuid, id, time.Second, time.Minute)
}

func runAction(ctx context.Context, actionRepo ActionRepo, uuid, leader, action string, params map[string]any) (*action.ActionResult, error) {
	id, err := actionRepo.RunAction(ctx, uuid, leader, action, params)
	if err != nil {
		return nil, err
	}
	return waitForActionCompleted(ctx, actionRepo, uuid, id, time.Second, time.Minute)
}

func waitForActionCompleted(ctx context.Context, actionRepo ActionRepo, uuid, id string, tickInterval, timeoutDuration time.Duration) (*action.ActionResult, error) {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	timeout := time.After(timeoutDuration)
	for {
		select {
		case <-ticker.C:
			result, err := actionRepo.GetResult(ctx, uuid, id)
			if err != nil {
				return nil, err
			}
			if result.Status == "completed" { // state.ActionCompleted
				return result, nil
			}
			continue

		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for action %s to become completed", id)

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
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
