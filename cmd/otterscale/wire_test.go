package main

import (
	"path/filepath"
	"testing"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

func TestWireServerWithoutKubernetesConfig(t *testing.T) {
	t.Setenv("KUBERNETES_SERVICE_HOST", "")
	t.Setenv("KUBERNETES_SERVICE_PORT", "")
	t.Setenv("OTTERSCALE_SERVER_EXTERNAL_URL", "https://server.example.com")
	t.Setenv("OTTERSCALE_SERVER_EXTERNAL_TUNNEL_URL", "https://tunnel.example.com")

	originalKubeconfig := clientcmd.RecommendedHomeFile
	clientcmd.RecommendedHomeFile = filepath.Join(t.TempDir(), "missing-kubeconfig")
	t.Cleanup(func() { clientcmd.RecommendedHomeFile = originalKubeconfig })

	conf, err := config.New()
	if err != nil {
		t.Fatalf("config.New: %v", err)
	}
	server, cleanup, err := wireServer(core.Version("test"), conf)
	if err != nil {
		t.Fatalf("wireServer must degrade when Kubernetes config is unavailable: %v", err)
	}
	defer cleanup()
	if server == nil {
		t.Fatal("wireServer returned a nil server")
	}
}
