package kube

import (
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const ns = "openhdc"

var (
	ttl       = int32(86400) //nolint:mnd
	succLimit = int32(5)     //nolint:mnd
	failLimit = int32(3)     //nolint:mnd
)

func newConfig() (*rest.Config, error) {
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	if isGoRun {
		path := filepath.Join(homedir.HomeDir(), ".kube", "config")
		return clientcmd.BuildConfigFromFlags("", path)
	}
	return rest.InClusterConfig()
}

func NewClientset() (*kubernetes.Clientset, error) {
	config, err := newConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
