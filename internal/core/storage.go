package core

import (
	"context"
	"errors"
	"sync"

	"github.com/juju/juju/api/client/action"
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

	uc.configs.Store(uuid, config)

	return config, nil
}

func (uc *StorageUseCase) newConfig(ctx context.Context, uuid, name string) (*StorageConfig, error) {
	// ceph-mon
	leader, err := uc.facility.GetLeader(ctx, uuid, name)
	if err != nil {
		return nil, err
	}

	// stdout
	ops, err := uc.action.ListOperations(ctx, uuid, leader, "juju-exec")
	if err != nil {
		return nil, err
	}
	return uc.getCephMonCommandOutput(ctx, uuid, ops)
}

func (uc *StorageUseCase) getCephMonCommandOutput(ctx context.Context, uuid string, ops []action.Operation) (*StorageConfig, error) {
	for i := range ops {
		for j := range ops[i].Actions {
			if ops[i].Actions[j].Action.Parameters["command"] != cephConfigCommand {
				continue
			}
			results, err := uc.action.ListResults(ctx, uuid, ops[i].Actions[j].Action.ID)
			if err != nil {
				return nil, err
			}
			for k := range results {
				stdout, ok := results[k].Output["stdout"]
				if !ok {
					continue
				}
				config, err := uc.extractStorageConfig([]byte(stdout.(string)))
				if err != nil {
					continue
				}
				if config.FSID == "" || config.MonHost == "" || config.Key == "" {
					continue
				}
				return config, nil
			}
		}
	}
	return nil, errors.New("ceph mon command output not found")
}

func (uc *StorageUseCase) extractStorageConfig(stdout []byte) (*StorageConfig, error) {
	file, err := ini.Load(stdout)
	if err != nil {
		return nil, err
	}
	return &StorageConfig{
		FSID:    file.Section("global").Key("fsid").String(),
		MonHost: file.Section("global").Key("mon_host").String(),
		Key:     file.Section("client.admin").Key("key").String(),
	}, nil
}
