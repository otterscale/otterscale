package kube

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/openhdc/openhdc/internal/env"
)

type kube struct {
	*kubernetes.Clientset
	*rest.Config
}

// Kubes represents a map of named Kubernetes client connections
type KubeMap map[string]*kube

// NewKubes creates a new Kubernetes clientset map.
// If running in-cluster, it automatically adds the cluster client with key "(empty)"
func NewKubeMap() (KubeMap, error) {
	k := make(map[string]*kube)
	if isInCluster() {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
		k[""] = &kube{
			Clientset: clientset,
			Config:    config,
		}
	}
	return k, nil
}

func (k KubeMap) Add(cluster string, config *rest.Config) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	k[cluster] = &kube{
		Clientset: clientset,
		Config:    config,
	}
	return nil
}

func (k KubeMap) get(cluster string) (*kube, error) {
	if kube, ok := k[cluster]; ok {
		return kube, nil
	}
	return nil, fmt.Errorf("kubernetes cluster %q not found", cluster)
}

func (k KubeMap) GetKubeClientset(cluster string) (*kubernetes.Clientset, error) {
	kube, err := k.get(cluster)
	if err != nil {
		return nil, err
	}
	return kube.Clientset, nil
}

func (k KubeMap) GetHelmConfig(cluster, namespace string) (*action.Configuration, error) {
	kube, err := k.get(cluster)
	if err != nil {
		return nil, err
	}

	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.APIServer = &kube.Config.Host
	configFlags.BearerToken = &kube.Config.BearerToken
	configFlags.CAFile = &kube.Config.CAFile
	configFlags.CertFile = &kube.Config.TLSClientConfig.CertFile
	configFlags.KeyFile = &kube.Config.TLSClientConfig.KeyFile
	configFlags.Insecure = &kube.Config.Insecure
	configFlags.Namespace = &namespace

	config := new(action.Configuration)
	if err := config.Init(configFlags, namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}
	return config, nil
}

func isInCluster() bool {
	val := os.Getenv(env.OPENHDC_IN_CLUSTER)
	inCluster, _ := strconv.ParseBool(strings.ToLower(val))
	return inCluster
}
