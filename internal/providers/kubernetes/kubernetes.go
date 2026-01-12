package kubernetes

import (
	"context"
	"fmt"
	"os"
	"sync"

	"k8s.io/client-go/tools/clientcmd"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core/application/cluster"
	"github.com/otterscale/otterscale/internal/core/scope"
)

var kubeConfigPath string = "/root/.kube/%s/config"

type Kubernetes struct {
	conf *config.Config

	kubeConfigPath string
	configs        sync.Map
	clientsets     sync.Map
	apiClientsets  sync.Map
}

func New(conf *config.Config) (*Kubernetes, error) {
	return &Kubernetes{
		conf:           conf,
		kubeConfigPath: kubeConfigPath,
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

func (m *Kubernetes) APIExtensionsClientset(scope string) (*clientset.Clientset, error) {
	return m.apiClientset(scope)
}

func (m *Kubernetes) Clientset(scope string) (*kubernetes.Clientset, error) {
	return m.clientset(scope)
}

func (m *Kubernetes) clientset(scope string) (*kubernetes.Clientset, error) {
	if v, ok := m.clientsets.Load(scope); ok {
		return v.(*kubernetes.Clientset), nil
	}

	config, err := m.newConfig(scope)
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

func (m *Kubernetes) getKubeConfig(path string) (*api.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return clientcmd.Load(data)
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
	kubeConfigPath := fmt.Sprintf(m.kubeConfigPath, scope)

	kubeConfig, err := m.getKubeConfig(kubeConfigPath)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.NewDefaultClientConfig(*kubeConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}

	// Write CA data to temp file if needed, but keep CAData for client compatibility
	if config.CAData != nil && config.CAFile == "" {
		fileName, err := m.writeCAToFile(config.CAData)
		if err != nil {
			return nil, err
		}
		config.CAFile = fileName
		// Keep CAData - don't clear it, both can coexist
	}

	return config, nil
}

func (m *Kubernetes) newConfigForAPIExtensions(scope string) (*rest.Config, error) {
	// For apiextensions clientset with client-cert auth, we need to ensure TLS is configured properly
	kubeConfigPath := fmt.Sprintf(m.kubeConfigPath, scope)

	kubeConfig, err := m.getKubeConfig(kubeConfigPath)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.NewDefaultClientConfig(*kubeConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}

	// For client certificate auth, ensure Insecure is false and CAData is present
	config.Insecure = false

	// Write CA data to temp file
	if config.CAData != nil && config.CAFile == "" {
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

	return rest.InClusterConfig()
}

func (m *Kubernetes) apiClientset(scope string) (*clientset.Clientset, error) {
	if v, ok := m.apiClientsets.Load(scope); ok {
		return v.(*clientset.Clientset), nil
	}

	config, err := m.newConfigForAPIExtensions(scope)
	if err != nil {
		return nil, err
	}

	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m.apiClientsets.Store(scope, clientset)

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
