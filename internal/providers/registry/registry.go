package registry

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"oras.land/oras-go/v2/registry/remote"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
)

type Registry struct {
	conf       *config.Config
	kubernetes *kubernetes.Kubernetes

	registries sync.Map
}

func New(conf *config.Config, kubernetes *kubernetes.Kubernetes) (*Registry, error) {
	return &Registry{
		conf:       conf,
		kubernetes: kubernetes,
	}, nil
}

func (m *Registry) client(scope string) (*remote.Registry, error) {
	if v, ok := m.registries.Load(scope); ok {
		return v.(*remote.Registry), nil
	}

	client, err := m.newClient(scope)
	if err != nil {
		return nil, err
	}

	m.registries.Store(scope, client)

	return client, nil
}

func (m *Registry) newClient(scope string) (*remote.Registry, error) {
	registryAddress, err := m.getRegistryURL(scope)
	if err != nil {
		return nil, err
	}

	registry, err := remote.NewRegistry(registryAddress)
	if err != nil {
		return nil, err
	}

	registry.PlainHTTP = true

	return registry, nil
}

func (m *Registry) getRegistryURL(scope string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ip, err := m.kubernetes.InternalIP(ctx, scope)
	if err != nil {
		return "", err
	}

	service, err := m.kubernetes.GetService(ctx, scope, "registry", "registry")
	if err != nil {
		return "", err
	}

	ports := service.Spec.Ports

	if len(ports) == 0 {
		return "", fmt.Errorf("registry service has no ports defined")
	}

	return ip + ":" + strconv.Itoa(int(ports[0].NodePort)), nil
}
