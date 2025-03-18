package kube

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/openhdc/openhdc/internal/env"
)

var ErrClusterNotFound = errors.New("cluster not found")

// Kubes represents a map of named Kubernetes client connections
type Kubes map[string]*kubernetes.Clientset

// NewKubes creates a new Kubernetes clientset map.
// If running in-cluster, it automatically adds the cluster client with key "f"
func NewKubes() (Kubes, error) {
	kb := make(Kubes)

	if isInCluster() {
		cfg, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		cs, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return nil, err
		}

		kb["f"] = cs
	}

	return kb, nil
}

// isInCluster returns true if the application should use in-cluster configuration.
func isInCluster() bool {
	val := os.Getenv(env.OPENHDC_IN_CLUSTER)
	inCluster, _ := strconv.ParseBool(strings.ToLower(val))
	return inCluster
}

// Add creates and adds a new Kubernetes client to the map with the specified name
func (k Kubes) Add(name, host, bearerToken string) error {
	client, err := newClient(host, bearerToken)
	if err != nil {
		return err
	}
	k[name] = client
	return nil
}

func (k Kubes) Get(name string) (*kubernetes.Clientset, error) {
	if client, ok := k[name]; ok {
		return client, nil
	}
	return nil, ErrClusterNotFound
}

// newClient creates a new Kubernetes clientset from the provided configuration.
func newClient(host, bearerToken string) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(&rest.Config{
		Host:        host,
		BearerToken: bearerToken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: !strings.HasPrefix(host, "https"),
		},
	})
}
