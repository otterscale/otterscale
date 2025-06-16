package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/state"
	"gopkg.in/ini.v1"
)

const cephConfigCommand = "ceph config generate-minimal-conf && ceph auth get client.admin"

type StorageConfig struct {
	FSID    string
	MonHost string
	Key     string
}

type StorageUseCase struct {
	action   ActionRepo
	facility FacilityRepo
	pool     CephPoolRepo

	configs sync.Map
}

func NewStorageUseCase(action ActionRepo, facility FacilityRepo, pool CephPoolRepo) *StorageUseCase {
	return &StorageUseCase{
		action:   action,
		facility: facility,
		pool:     pool,
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
	// ceph-mon
	leader, err := uc.facility.GetLeader(ctx, uuid, name)
	if err != nil {
		return nil, err
	}
	id, err := uc.action.Run(ctx, uuid, leader, cephConfigCommand)
	if err != nil {
		return nil, err
	}
	result, err := uc.waitForActionCompleted(ctx, uuid, id)
	if err != nil {
		return nil, err
	}
	return uc.extractStorageConfig(result)
}

func (uc *StorageUseCase) waitForActionCompleted(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
	const tickInterval = time.Second
	const timeoutDuration = 5 * time.Second

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	timeout := time.After(timeoutDuration)
	for {
		select {
		case <-ticker.C:
			result, err := uc.action.GetResult(ctx, uuid, id)
			if err != nil {
				return nil, err
			}
			if result.Status == string(state.ActionCompleted) {
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

func (uc *StorageUseCase) extractStorageConfig(result *action.ActionResult) (*StorageConfig, error) {
	stdout, ok := result.Output["stdout"]
	if !ok {
		return nil, errors.New("ceph mon stdout not found")
	}
	file, err := ini.Load([]byte(stdout.(string)))
	if err != nil {
		return nil, err
	}
	fsID := file.Section("global").Key("fsid").String()
	if fsID == "" {
		return nil, errors.New("ceph mon stdout fsid not found")
	}
	monHost := file.Section("global").Key("mon_host").String()
	if monHost == "" {
		return nil, errors.New("ceph mon stdout mon_host not found")
	}
	key := file.Section("client.admin").Key("key").String()
	if key == "" {
		return nil, errors.New("ceph mon stdout key not found")
	}
	return &StorageConfig{
		FSID:    fsID,
		MonHost: monHost,
		Key:     key,
	}, nil
}
