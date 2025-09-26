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
func instanceType_mustNewKube(t testing.TB, cfg *config.Config) *Kube {
	k, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create Kube instance: %v", err)
	}
	return k
}

/* -------------------------------------------------------------------------- *
 * Construction / interface compliance                                         *
 * -------------------------------------------------------------------------- */
func TestNewInstanceType(t *testing.T) {
	k := instanceType_mustNewKube(t, &config.Config{})
	repo := NewInstanceType(k)
	if repo == nil {
		t.Fatal("expected InstanceType repository to be created, got nil")
	}
	// Verify it implements the interface
	var _ oscore.KubeInstanceTypeRepo = repo
}

/* -------------------------------------------------------------------------- *
 * Structure validation                                                       *
 * -------------------------------------------------------------------------- */
func TestInstanceType_Structure(t *testing.T) {
	k := instanceType_mustNewKube(t, &config.Config{})
	repo := NewInstanceType(k)

	it, ok := repo.(*instanceType)
	if !ok {
		t.Fatal("expected *instanceType, got a different type")
	}
	if it.kube == nil {
		t.Error("expected kube field to be set, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Error handling – invalid config should cause client creation failure       *
 * -------------------------------------------------------------------------- */
func TestInstanceType_ErrorHandling(t *testing.T) {
	k := instanceType_mustNewKube(t, &config.Config{})
	repo := NewInstanceType(k)
	ctx := context.Background()

	// Invalid config will cause client creation to fail
	invalidCfg := &rest.Config{Host: "https://example.invalid"}

	cases := []struct {
		name string
		call func() error
	}{
		{"ListClusterWide", func() error {
			_, err := repo.ListClusterWide(ctx, invalidCfg)
			return err
		}},
		{"List", func() error {
			_, err := repo.List(ctx, invalidCfg, "default")
			return err
		}},
		{"Get", func() error {
			_, err := repo.Get(ctx, invalidCfg, "default", "test-it")
			return err
		}},
		{"Create", func() error {
			_, err := repo.Create(ctx, invalidCfg, "default", "test-it", 2, 4096)
			return err
		}},
		{"Delete", func() error {
			return repo.Delete(ctx, invalidCfg, "default", "test-it")
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
func TestInstanceType_EdgeCases(t *testing.T) {
	t.Run("nil_kube", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using instanceType with nil Kube")
			}
		}()

		it := &instanceType{kube: nil}
		ctx := context.Background()
		invalidCfg := &rest.Config{Host: "https://example.invalid"}

		// This should panic because kube is nil
		_, _ = it.ListClusterWide(ctx, invalidCfg)
	})
}

/* -------------------------------------------------------------------------- *
 * Benchmarks                                                                 *
 * -------------------------------------------------------------------------- */
func BenchmarkInstanceType_Creation(b *testing.B) {
	k := instanceType_mustNewKube(b, &config.Config{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewInstanceType(k); repo == nil {
			b.Fatal("failed to create instanceType repo")
		}
	}
}
