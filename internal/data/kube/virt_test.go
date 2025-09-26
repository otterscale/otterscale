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
func virt_mustNewKube(t testing.TB, cfg *config.Config) *Kube {
	k, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create Kube instance: %v", err)
	}
	return k
}

/* -------------------------------------------------------------------------- *
 * Construction / interface compliance                                         *
 * -------------------------------------------------------------------------- */
func TestNewVirt(t *testing.T) {
	k := virt_mustNewKube(t, &config.Config{})
	repo := NewVirt(k)
	if repo == nil {
		t.Fatal("expected Virt repository to be created, got nil")
	}
	// Verify it implements the interface
	var _ oscore.KubeVirtRepo = repo
}

/* -------------------------------------------------------------------------- *
 * Structure validation                                                       *
 * -------------------------------------------------------------------------- */
func TestVirt_Structure(t *testing.T) {
	k := virt_mustNewKube(t, &config.Config{})
	repo := NewVirt(k)

	v, ok := repo.(*virt)
	if !ok {
		t.Fatal("expected *virt, got a different type")
	}
	if v.kube == nil {
		t.Error("expected kube field to be set, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Error handling – invalid config should cause client creation failure       *
 * -------------------------------------------------------------------------- */
func TestVirt_ErrorHandling(t *testing.T) {
	k := virt_mustNewKube(t, &config.Config{})
	repo := NewVirt(k)
	ctx := context.Background()

	// Invalid config will cause client creation to fail
	invalidCfg := &rest.Config{Host: "https://example.invalid"}

	cases := []struct {
		name string
		call func() error
	}{
		{"ListVirtualMachines", func() error {
			_, err := repo.ListVirtualMachines(ctx, invalidCfg, "default")
			return err
		}},
		{"GetVirtualMachine", func() error {
			_, err := repo.GetVirtualMachine(ctx, invalidCfg, "default", "test-vm")
			return err
		}},
		{"CreateVirtualMachine", func() error {
			_, err := repo.CreateVirtualMachine(ctx, invalidCfg, "default", "test-vm", "instance-type", "boot-dv", "startup-script")
			return err
		}},
		{"UpdateVirtualMachine", func() error {
			_, err := repo.UpdateVirtualMachine(ctx, invalidCfg, "default", &oscore.VirtualMachine{})
			return err
		}},
		{"DeleteVirtualMachine", func() error {
			return repo.DeleteVirtualMachine(ctx, invalidCfg, "default", "test-vm")
		}},
		{"StartVirtualMachine", func() error {
			return repo.StartVirtualMachine(ctx, invalidCfg, "default", "test-vm")
		}},
		{"StopVirtualMachine", func() error {
			return repo.StopVirtualMachine(ctx, invalidCfg, "default", "test-vm")
		}},
		{"RestartVirtualMachine", func() error {
			return repo.RestartVirtualMachine(ctx, invalidCfg, "default", "test-vm")
		}},
		{"ListInstances", func() error {
			_, err := repo.ListInstances(ctx, invalidCfg, "default")
			return err
		}},
		{"GetInstance", func() error {
			_, err := repo.GetInstance(ctx, invalidCfg, "default", "test-instance")
			return err
		}},
		{"PauseInstance", func() error {
			return repo.PauseInstance(ctx, invalidCfg, "default", "test-instance")
		}},
		{"UnpauseInstance", func() error {
			return repo.UnpauseInstance(ctx, invalidCfg, "default", "test-instance")
		}},
		{"MigrateInstance", func() error {
			_, err := repo.MigrateInstance(ctx, invalidCfg, "default", "test-instance", "hostname")
			return err
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
func TestVirt_EdgeCases(t *testing.T) {
	t.Run("nil_kube", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when using virt with nil Kube")
			}
		}()

		v := &virt{kube: nil}
		ctx := context.Background()
		invalidCfg := &rest.Config{Host: "https://example.invalid"}

		// This should panic because kube is nil
		_, _ = v.ListVirtualMachines(ctx, invalidCfg, "default")
	})
}

/* -------------------------------------------------------------------------- *
 * Benchmarks                                                                 *
 * -------------------------------------------------------------------------- */
func BenchmarkVirt_Creation(b *testing.B) {
	k := virt_mustNewKube(b, &config.Config{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if repo := NewVirt(k); repo == nil {
			b.Fatal("failed to create virt repo")
		}
	}
}
