package registry

import (
	"context"
	"strconv"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
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
	config, err := m.kubernetes.Config(scope)
	if err != nil {
		return nil, err
	}

	clientset, err := k8s.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	registryAddress, err := m.getRegistryURL(clientset)
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

func (m *Registry) getRegistryURL(clientset *k8s.Clientset) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	opts := metav1.GetOptions{}

	service, err := clientset.CoreV1().Services("registry").Get(ctx, "registry", opts)
	if err != nil {
		return "", err
	}

	// TODO: validate
	return service.Spec.ClusterIP + ":" + strconv.Itoa(int(service.Spec.Ports[0].Port)), nil
}
