package kube

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/openhdc/openhdc/internal/env"
)

// Kubes represents a map of named Kubernetes client connections
type Kubes map[string]*kubernetes.Clientset

// NewKubes creates a new Kubernetes clientset map.
// If running in-cluster, it automatically adds the cluster client with key "(empty)"
func NewKubes() (Kubes, error) {
	k := make(Kubes)
	if isInCluster() {
		cfg, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		client, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return nil, err
		}
		k[""] = client
	}
	return k, nil
}

func (k Kubes) Add(cluster string, cfg *rest.Config) error {
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}
	k[cluster] = client
	return nil
}

func (k Kubes) Get(cluster string) (*kubernetes.Clientset, error) {
	if client, ok := k[cluster]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("kubernetes cluster %q not found", cluster)
}

func isInCluster() bool {
	val := os.Getenv(env.OPENHDC_IN_CLUSTER)
	inCluster, _ := strconv.ParseBool(strings.ToLower(val))
	return inCluster
}
