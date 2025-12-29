package kubernetes

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core/application/cluster"
	"github.com/otterscale/otterscale/internal/core/scope"
	"github.com/otterscale/otterscale/internal/mux/impersonation"
	"github.com/otterscale/otterscale/internal/providers/juju"
)

type Kubernetes struct {
	conf *config.Config
	juju *juju.Juju

	configs       sync.Map
	clientsets    sync.Map
	extClientsets sync.Map
}

func New(conf *config.Config, juju *juju.Juju) (*Kubernetes, error) {
	return &Kubernetes{
		conf: conf,
		juju: juju,
	}, nil
}

func (m *Kubernetes) Config(scope string) (*rest.Config, error) {
	if v, ok := m.configs.Load(scope); ok {
		return v.(*rest.Config), nil
	}

	config, err := m.getConfig(scope)
	if err != nil {
		return nil, err
	}

	m.configs.Store(scope, config)

	return config, nil
}

func (m *Kubernetes) dynamic(ctx context.Context, cluster string) (*dynamic.DynamicClient, error) {
	userSub, ok := impersonation.GetSubject(ctx)
	if !ok {
		return nil, fmt.Errorf("user sub not found in context")
	}

	config, err := m.Config(cluster)
	if err != nil {
		return nil, err
	}

	userConfig := rest.CopyConfig(config)

	userConfig.Impersonate = rest.ImpersonationConfig{
		UserName: userSub,
	}

	return dynamic.NewForConfig(userConfig)
}

func (m *Kubernetes) discovery(cluster string) (*discovery.DiscoveryClient, error) {
	config, err := m.Config(cluster)
	if err != nil {
		return nil, err
	}

	return discovery.NewDiscoveryClientForConfig(config)
}

func (m *Kubernetes) InternalIP(ctx context.Context, scope string) (string, error) {
	controlPlanes, err := m.listControlPlanes(ctx, scope)
	if err != nil {
		return "", err
	}

	for i := range controlPlanes {
		if !isNodeReady(&controlPlanes[i]) {
			continue
		}

		if ip := getInternalIP(&controlPlanes[i]); ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("no control plane node with InternalIP found")
}

func (m *Kubernetes) GetService(ctx context.Context, scope, namespace, name string) (*corev1.Service, error) {
	clientset, err := m.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.CoreV1().Services(namespace).Get(ctx, name, opts)
}

func (m *Kubernetes) newMicroK8sConfig() (*rest.Config, error) {
	kubeConfig, err := base64.StdEncoding.DecodeString(m.conf.MicroK8sConfig())
	if err != nil {
		return nil, err
	}

	configAPI, err := clientcmd.Load(kubeConfig)
	if err != nil {
		return nil, err
	}

	return clientcmd.NewDefaultClientConfig(*configAPI, &clientcmd.ConfigOverrides{}).ClientConfig()
}

func (m *Kubernetes) getKubeConfig(ctx context.Context, scope, name string) (*api.Config, error) {
	result, err := m.juju.Run(ctx, scope, name, "get-kubeconfig", nil)
	if err != nil {
		return nil, err
	}
	return clientcmd.Load([]byte(result["kubeconfig"].(string)))
}

func (m *Kubernetes) writeCAToFile(caData []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", "otterscale-ca-*.crt")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(caData); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func (m *Kubernetes) newConfig(scope string) (*rest.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	kubeConfig, err := m.getKubeConfig(ctx, scope, "kubernetes-control-plane")
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.NewDefaultClientConfig(*kubeConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}

	// Write CA data to temp file for helm
	if config.CAData != nil {
		fileName, err := m.writeCAToFile(config.CAData)
		if err != nil {
			return nil, err
		}
		config.CAFile = fileName
	}

	return config, nil
}

func (m *Kubernetes) getConfig(scopeName string) (*rest.Config, error) {
	if scopeName != scope.ReservedName {
		return m.newConfig(scopeName)
	}

	if os.Getenv("OTTERSCALE_CONTAINER") == "true" {
		return rest.InClusterConfig()
	}

	return m.newMicroK8sConfig()
}

func (m *Kubernetes) clientset(scope string) (*kubernetes.Clientset, error) {
	if v, ok := m.clientsets.Load(scope); ok {
		return v.(*kubernetes.Clientset), nil
	}

	config, err := m.Config(scope)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.clientsets.Store(scope, clientset)

	return clientset, nil
}

func (m *Kubernetes) extClientset(scope string) (*clientset.Clientset, error) {
	if v, ok := m.extClientsets.Load(scope); ok {
		return v.(*clientset.Clientset), nil
	}

	config, err := m.Config(scope)
	if err != nil {
		return nil, err
	}

	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.extClientsets.Store(scope, clientset)

	return clientset, nil
}

func (m *Kubernetes) listControlPlanes(ctx context.Context, scope string) ([]cluster.Node, error) {
	clientset, err := m.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/control-plane",
	}

	list, err := clientset.CoreV1().Nodes().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func isNodeReady(node *cluster.Node) bool {
	for i := range node.Status.Conditions {
		if node.Status.Conditions[i].Type == corev1.NodeReady &&
			node.Status.Conditions[i].Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func getInternalIP(node *cluster.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}

func (m *Kubernetes) ValidateKubeConfig(kubeconfig string) (bool, error) {
	if kubeconfig == "" {
		return false, fmt.Errorf("kubeconfig is empty")
	}

	kubeconfig_decode, err := base64.StdEncoding.DecodeString(kubeconfig)
	if err != nil {
		return false, err
	}

	configAPI, err := clientcmd.Load([]byte(kubeconfig_decode))
	if err != nil {
		return false, err
	}

	config, err := clientcmd.NewDefaultClientConfig(*configAPI, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return false, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, err
	}

	_, err = clientset.ServerVersion()
	if err != nil {
		return false, err
	}

	return true, nil
}
