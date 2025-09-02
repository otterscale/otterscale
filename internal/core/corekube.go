package core

import (
	"context"
	"os"
	"sync"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const kubeConfigAction = "get-kubeconfig"

var kubeConfigs sync.Map

func kubeConfig(ctx context.Context, facility FacilityRepo, action ActionRepo, uuid, name string) (*rest.Config, error) {
	key := uuid + "/" + name

	if v, ok := kubeConfigs.Load(key); ok {
		return v.(*rest.Config), nil
	}

	config, err := newKubeConfig(ctx, facility, action, uuid, name)
	if err != nil {
		return nil, err
	}

	kubeConfigs.Store(key, config)

	return config, nil
}

func newKubeConfig(ctx context.Context, facility FacilityRepo, action ActionRepo, uuid, name string) (*rest.Config, error) {
	// kubernetes-control-plane
	leader, err := facility.GetLeader(ctx, uuid, name)
	if err != nil {
		return nil, err
	}

	result, err := runAction(ctx, action, uuid, leader, kubeConfigAction, nil)
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
