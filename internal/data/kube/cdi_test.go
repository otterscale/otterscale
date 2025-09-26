package kube

import (
	"context"
	"testing"

	"k8s.io/client-go/rest"

	"github.com/otterscale/otterscale/internal/config"
	oscore "github.com/otterscale/otterscale/internal/core"
)

/* -------------------------------------------------------------------------- *
 * Helper – works for both normal tests (testing.T) and benchmarks (testing.B) *
 * -------------------------------------------------------------------------- */
func cdi_mustNewKube(t testing.TB, cfg *config.Config) *Kube {
	k, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create Kube instance: %v", err)
	}
	return k
}

/* -------------------------------------------------------------------------- *
 * Construction / interface compliance                                         *
 * -------------------------------------------------------------------------- */
func TestNewCDI(t *testing.T) {
	k := cdi_mustNewKube(t, &config.Config{})
	repo := NewCDI(k)
	if repo == nil {
		t.Fatal("expected CDI repository to be created, got nil")
	}
	// Verify it implements the interface
	var _ oscore.KubeCDIRepo = repo
}

/* -------------------------------------------------------------------------- *
 * Structure validation                                                       *
 * -------------------------------------------------------------------------- */
func TestCDI_Structure(t *testing.T) {
	k := cdi_mustNewKube(t, &config.Config{})
	repo := NewCDI(k)

	c, ok := repo.(*cdi)
	if !ok {
		t.Fatal("expected *cdi, got a different type")
	}
	if c.kube == nil {
		t.Error("expected kube field to be set, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Error handling – invalid config should cause client creation failure       *
 * -------------------------------------------------------------------------- */
func TestCDI_ErrorHandling(t *testing.T) {
	k := cdi_mustNewKube(t, &config.Config{})
	repo := NewCDI(k)
	ctx := context.Background()

	// Invalid config will cause client creation to fail
	invalidCfg := &rest.Config{Host: "https://example.invalid"}

	cases := []struct {
		name string
		call func() error
	}{
		{"ListDataVolumes", func() error {
			_, err := repo.ListDataVolumes(ctx, invalidCfg, "default", false)
			return err
		}},
		{"GetDataVolume", func() error {
			_, err := repo.GetDataVolume(ctx, invalidCfg, "default", "test-dv")
			return err
		}},
		{"CreateDataVolume", func() error {
			_, err := repo.CreateDataVolume(ctx, invalidCfg, "default", "test-dv", oscore.SourceType(0), "src-data", 1073741824, false)
			return err
		}},
		{"DeleteDataVolume", func() error {
			return repo.DeleteDataVolume(ctx, invalidCfg, "default", "test-dv")
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); err == nil {
				t.Errorf("expected error for %s with invalid config, got nil", tc.name)
			} else {
				t.Logf("%s returned expected error: %v", tc.name, err)
			}
		})
	}
}

/* -------------------------------------------------------------------------- *
 * Edge cases – nil Kube should panic                                         *
 * -------------------------------------------------------------------------- */
func TestCDI_EdgeCases(t *testing.T) {
	t.Run("nil_kube", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using cdi with nil Kube")
			}
		}()

		c := &cdi{kube: nil}
		ctx := context.Background()
		invalidCfg := &rest.Config{Host: "https://example.invalid"}

		// This should panic because kube is nil
		_, _ = c.ListDataVolumes(ctx, invalidCfg, "default", false)
	})
}

/* -------------------------------------------------------------------------- *
 * Benchmarks                                                                 *
 * -------------------------------------------------------------------------- */
func BenchmarkCDI_Creation(b *testing.B) {
	k := cdi_mustNewKube(b, &config.Config{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewCDI(k); repo == nil {
			b.Fatal("failed to create cdi repo")
		}
	}
}
