package kubernetes

import (
	"log/slog"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ProvideInClusterConfig is a Wire provider that returns a
// *rest.Config for in-cluster Kubernetes API access. It falls back to
// the user's kubeconfig for local development.
func ProvideInClusterConfig() (*rest.Config, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		slog.Warn("in-cluster config not available, falling back to kubeconfig", "error", err)
		return clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	}
	return cfg, nil
}
