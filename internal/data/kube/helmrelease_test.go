package kube

import (
	"testing"

	"connectrpc.com/connect"
	"k8s.io/client-go/rest"

	"github.com/otterscale/otterscale/internal/config"
	oscore "github.com/otterscale/otterscale/internal/core"
)

// -----------------------------------------------------------------------------
// Construction / interface compliance
// -----------------------------------------------------------------------------
func TestNewHelmRelease(t *testing.T) {
	k := helm_mustNewKube(t, &config.Config{})
	rl, err := NewHelmRelease(k)
	if err != nil {
		t.Fatalf("NewHelmRelease returned error: %v", err)
	}
	if rl == nil {
		t.Fatal("NewHelmRelease returned nil")
	}
	var _ oscore.ReleaseRepo = rl
}

func TestHelmRelease_InterfaceCompliance(t *testing.T) {
	k := helm_mustNewKube(t, &config.Config{})
	rl, _ := NewHelmRelease(k)
	var _ oscore.ReleaseRepo = rl
}

// -----------------------------------------------------------------------------
// Structure validation
// -----------------------------------------------------------------------------
func TestHelmRelease_Structure(t *testing.T) {
	k := helm_mustNewKube(t, &config.Config{})
	rl, _ := NewHelmRelease(k)

	h, ok := rl.(*helmRelease)
	if !ok {
		t.Fatalf("expected *helmRelease, got %T", rl)
	}
	if h.kube == nil {
		t.Error("helmRelease.kube field is nil")
	}
}

// -----------------------------------------------------------------------------
// List – expected to fail because there is no real cluster
// -----------------------------------------------------------------------------
func TestHelmRelease_List_Error(t *testing.T) {
	k := helm_mustNewKube(t, &config.Config{})
	rl, _ := NewHelmRelease(k)

	restCfg := &rest.Config{Host: "https://example.invalid"}
	_, err := rl.List(restCfg, "default")
	if err == nil {
		t.Error("expected error from List with non‑existent cluster, got nil")
	}
}

// -----------------------------------------------------------------------------
// Install – invalid release name (validation happens before any Helm work)
// -----------------------------------------------------------------------------
func TestHelmRelease_Install_InvalidName(t *testing.T) {
	k := helm_mustNewKube(t, &config.Config{})
	rl, _ := NewHelmRelease(k)

	restCfg := &rest.Config{Host: "https://example.invalid"}
	_, err := rl.Install(restCfg, "default", "bad name", false, "any/chart", nil)

	if err == nil {
		t.Fatal("expected error for invalid release name, got nil")
	}
	if connect.CodeOf(err) != connect.CodeInvalidArgument {
		t.Fatalf("expected connect.CodeInvalidArgument, got %v", connect.CodeOf(err))
	}
}

// -----------------------------------------------------------------------------
// Install – chart cannot be located
// -----------------------------------------------------------------------------
func TestHelmRelease_Install_ChartNotFound(t *testing.T) {
	k := helm_mustNewKube(t, &config.Config{})
	rl, _ := NewHelmRelease(k)

	restCfg := &rest.Config{Host: "https://example.invalid"}
	_, err := rl.Install(restCfg, "default", "valid-name", false, "nonexistent/chart", nil)

	if err == nil {
		t.Fatal("expected error when chart cannot be located, got nil")
	}
}

// -----------------------------------------------------------------------------
// Uninstall – expected error (no cluster)
// -----------------------------------------------------------------------------
func TestHelmRelease_Uninstall_Error(t *testing.T) {
	k := helm_mustNewKube(t, &config.Config{})
	rl, _ := NewHelmRelease(k)

	restCfg := &rest.Config{Host: "https://example.invalid"}
	_, err := rl.Uninstall(restCfg, "default", "some-release", false)
	if err == nil {
		t.Error("expected error from Uninstall with non‑existent cluster, got nil")
	}
}

// -----------------------------------------------------------------------------
// Upgrade – chart not found
// -----------------------------------------------------------------------------
func TestHelmRelease_Upgrade_ChartNotFound(t *testing.T) {
	k := helm_mustNewKube(t, &config.Config{})
	rl, _ := NewHelmRelease(k)

	restCfg := &rest.Config{Host: "https://example.invalid"}
	_, err := rl.Upgrade(restCfg, "default", "rel", false, "missing/chart", nil)
	if err == nil {
		t.Fatal("expected error when chart cannot be located in Upgrade, got nil")
	}
}

// -----------------------------------------------------------------------------
// Rollback – expected error (no cluster)
// -----------------------------------------------------------------------------
func TestHelmRelease_Rollback_Error(t *testing.T) {
	k := helm_mustNewKube(t, &config.Config{})
	rl, _ := NewHelmRelease(k)

	restCfg := &rest.Config{Host: "https://example.invalid"}
	err := rl.Rollback(restCfg, "default", "rel", false)
	if err == nil {
		t.Error("expected error from Rollback with non‑existent cluster, got nil")
	}
}

// -----------------------------------------------------------------------------
// GetValues – expected error (no cluster)
// -----------------------------------------------------------------------------
func TestHelmRelease_GetValues_Error(t *testing.T) {
	k := helm_mustNewKube(t, &config.Config{})
	rl, _ := NewHelmRelease(k)

	restCfg := &rest.Config{Host: "https://example.invalid"}
	_, err := rl.GetValues(restCfg, "default", "rel")
	if err == nil {
		t.Error("expected error from GetValues with non‑existent cluster, got nil")
	}
}

// -----------------------------------------------------------------------------
// Edge case – nil Kube should panic on any method call
// -----------------------------------------------------------------------------
func TestHelmRelease_EdgeCases_NilKube(t *testing.T) {
	// Helper that calls all public methods and expects a panic.
	callAll := func() {
		h := &helmRelease{kube: nil}
		cfg := &rest.Config{Host: "https://example.invalid"}

		_, _ = h.List(cfg, "default")
		_, _ = h.Install(cfg, "default", "name", false, "chart/ref", nil)
		_, _ = h.Uninstall(cfg, "default", "name", false)
		_, _ = h.Upgrade(cfg, "default", "name", false, "chart/ref", nil)
		_ = h.Rollback(cfg, "default", "name", false)
		_, _ = h.GetValues(cfg, "default", "name")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when using helmRelease with nil Kube")
		}
	}()
	callAll()
}

// -----------------------------------------------------------------------------
// Benchmarks – creation and List call overhead
// -----------------------------------------------------------------------------
func BenchmarkHelmRelease_Creation(b *testing.B) {
	cfg := &config.Config{}
	k := helm_mustNewKube(b, cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := NewHelmRelease(k); err != nil {
			b.Fatalf("creation error: %v", err)
		}
	}
}

func BenchmarkHelmRelease_List(b *testing.B) {
	cfg := &config.Config{}
	k := helm_mustNewKube(b, cfg)
	rl, _ := NewHelmRelease(k)
	restCfg := &rest.Config{Host: "https://example.invalid"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = rl.List(restCfg, "default")
	}
}
