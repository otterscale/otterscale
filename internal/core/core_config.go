package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/juju/juju/api/client/action"
	"golang.org/x/sync/errgroup"
	"gopkg.in/ini.v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeConfigMap sync.Map
	cephConfigMap sync.Map
)

func kubeConfig(ctx context.Context, facility FacilityRepo, action ActionRepo, scope, name string) (*rest.Config, error) {
	key := scope + "/" + name

	if v, ok := kubeConfigMap.Load(key); ok {
		return v.(*rest.Config), nil
	}

	config, err := newKubeConfig(ctx, facility, action, scope, name)
	if err != nil {
		return nil, err
	}

	kubeConfigMap.Store(key, config)

	return config, nil
}

func newKubeConfig(ctx context.Context, facility FacilityRepo, action ActionRepo, scope, name string) (*rest.Config, error) {
	// kubernetes-control-plane
	leader, err := facility.GetLeader(ctx, scope, name)
	if err != nil {
		return nil, err
	}

	result, err := runAction(ctx, action, scope, leader, "get-kubeconfig", nil)
	if err != nil {
		return nil, err
	}

	configAPI, err := clientcmd.Load([]byte(result.Output["kubeconfig"].(string)))
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.NewDefaultClientConfig(*configAPI, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}

	// Write CA data to temp file for helm
	if config.CAData != nil {
		tmpFile, err := os.CreateTemp("", "otterscale-ca-*.crt")
		if err != nil {
			return nil, err
		}
		defer tmpFile.Close()

		if _, err := tmpFile.Write(config.CAData); err != nil {
			return nil, err
		}
		config.CAFile = tmpFile.Name()
	}

	return config, nil
}

func cephConfig(ctx context.Context, facility FacilityRepo, action ActionRepo, scope, name string) (*CephConfig, error) {
	key := scope + "/" + name

	if v, ok := cephConfigMap.Load(key); ok {
		return v.(*CephConfig), nil
	}

	config, err := newCephConfig(ctx, facility, action, scope, name)
	if err != nil {
		return nil, err
	}

	cephConfigMap.Store(key, config)

	return config, nil
}

func newCephConfig(ctx context.Context, facility FacilityRepo, action ActionRepo, scope, name string) (*CephConfig, error) {
	var (
		leader   string
		endpoint string
	)
	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		monLeader, err := facility.GetLeader(egctx, scope, name) // ceph-mon
		if err != nil {
			return err
		}
		leader = monLeader
		return nil
	})
	eg.Go(func() error {
		leader, err := facility.GetLeader(egctx, scope, rgwName(name))
		if err != nil {
			return err
		}
		info, err := facility.GetUnitInfo(egctx, scope, leader)
		if err != nil {
			return err
		}
		endpoint = fmt.Sprintf("http://%s", info.PublicAddress)
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	var cephConfig *CephClusterConfig
	rgwConfig := &CephObjectConfig{
		Endpoint: endpoint,
	}

	eg, egctx = errgroup.WithContext(ctx)
	eg.Go(func() error {
		result, err := runCommand(egctx, action, scope, leader, "ceph config generate-minimal-conf && ceph auth get client.admin")
		if err != nil {
			return err
		}
		config, err := extractCephClusterConfig(result)
		if err != nil {
			return err
		}
		cephConfig = config
		return nil
	})
	eg.Go(func() error {
		listResult, err := runCommand(egctx, action, scope, leader, "radosgw-admin user list")
		if err != nil {
			return err
		}
		cmd, err := getRGWCommand(listResult)
		if err != nil {
			return err
		}
		result, err := runCommand(egctx, action, scope, leader, cmd)
		if err != nil {
			return err
		}
		config, err := extractCephObjectConfig(result)
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
	return &CephConfig{
		CephClusterConfig: cephConfig,
		CephObjectConfig:  rgwConfig,
	}, nil
}

func extractCephClusterConfig(result *action.ActionResult) (*CephClusterConfig, error) {
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
	return &CephClusterConfig{
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

func extractCephObjectConfig(result *action.ActionResult) (*CephObjectConfig, error) {
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
	return &CephObjectConfig{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, nil
}

func runCommand(ctx context.Context, actionRepo ActionRepo, scope, leader, command string) (*action.ActionResult, error) {
	id, err := actionRepo.RunCommand(ctx, scope, leader, command)
	if err != nil {
		return nil, err
	}
	return waitForActionCompleted(ctx, actionRepo, scope, id, time.Second, time.Minute)
}

func runAction(ctx context.Context, actionRepo ActionRepo, scope, leader, action string, params map[string]any) (*action.ActionResult, error) {
	id, err := actionRepo.RunAction(ctx, scope, leader, action, params)
	if err != nil {
		return nil, err
	}
	return waitForActionCompleted(ctx, actionRepo, scope, id, time.Second, time.Minute)
}

func waitForActionCompleted(ctx context.Context, actionRepo ActionRepo, scope, id string, tickInterval, timeoutDuration time.Duration) (*action.ActionResult, error) {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	timeout := time.After(timeoutDuration)
	for {
		select {
		case <-ticker.C:
			result, err := actionRepo.GetResult(ctx, scope, id)
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
