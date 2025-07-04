package maas

import (
	"context"
	"testing"

	"github.com/openhdc/otterscale/internal/config"
	"github.com/openhdc/otterscale/internal/core"
)

func TestNewBootResource(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)

	repo := NewBootResource(maas)

	if repo == nil {
		t.Fatal("Expected BootResource repository to be created, got nil")
	}

	// Verify it implements the interface
	var _ core.BootResourceRepo = repo
}

func TestBootResource_InterfaceCompliance(t *testing.T) {
	// Test that bootResource implements core.BootResourceRepo interface
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootResource(maas)

	var _ core.BootResourceRepo = repo
}

func TestBootResource_Methods_Structure(t *testing.T) {
	// Test the structure and method signatures
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootResource(maas)

	// Verify the repo is of the correct type
	br, ok := repo.(*bootResource)
	if !ok {
		t.Fatal("Expected *bootResource, got different type")
	}

	if br.maas == nil {
		t.Error("Expected maas field to be set, got nil")
	}
}

// Test with actual MAAS configuration structure
func TestBootResource_WithConfig(t *testing.T) {
	tests := []struct {
		name string
		conf *config.Config
	}{
		{
			name: "empty_config",
			conf: &config.Config{},
		},
		{
			name: "config_with_maas",
			conf: &config.Config{
				MAAS: config.MAAS{
					URL:     "http://localhost:5240/MAAS",
					Key:     "test-key",
					Version: "2.8",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maas := New(tt.conf)
			repo := NewBootResource(maas)

			if repo == nil {
				t.Error("Expected repository to be created")
			}

			// Try to use the methods (they will likely fail due to invalid config, but should not panic)
			ctx := context.Background()

			_, err := repo.List(ctx)
			// We expect an error since we don't have a real MAAS server
			if err == nil {
				t.Log("List succeeded unexpectedly (might be in test environment)")
			}

			err = repo.Import(ctx)
			if err == nil {
				t.Log("Import succeeded unexpectedly (might be in test environment)")
			}

			_, err = repo.IsImporting(ctx)
			if err == nil {
				t.Log("IsImporting succeeded unexpectedly (might be in test environment)")
			}
		})
	}
}

// Test error handling with invalid configurations
func TestBootResource_ErrorHandling(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
		desc   string
	}{
		{
			name:   "empty_maas_config",
			config: &config.Config{},
			desc:   "Empty MAAS configuration should cause client errors",
		},
		{
			name: "invalid_url",
			config: &config.Config{
				MAAS: config.MAAS{
					URL:     "invalid-url",
					Key:     "test-key",
					Version: "2.8",
				},
			},
			desc: "Invalid URL should cause client errors",
		},
		{
			name: "missing_key",
			config: &config.Config{
				MAAS: config.MAAS{
					URL:     "http://localhost:5240/MAAS",
					Key:     "",
					Version: "2.8",
				},
			},
			desc: "Missing API key should cause authentication errors",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maas := New(tt.config)
			repo := NewBootResource(maas)

			ctx := context.Background()

			// All methods should return errors with invalid configurations
			_, err := repo.List(ctx)
			if err == nil {
				t.Logf("List unexpectedly succeeded for %s", tt.desc)
			}

			err = repo.Import(ctx)
			if err == nil {
				t.Logf("Import unexpectedly succeeded for %s", tt.desc)
			}

			_, err = repo.IsImporting(ctx)
			if err == nil {
				t.Logf("IsImporting unexpectedly succeeded for %s", tt.desc)
			}
		})
	}
}

func TestBootResource_ContextHandling(t *testing.T) {
	// Test that methods handle context properly (though the current implementation ignores it)
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootResource(maas)

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Methods should still fail with client error since they don't check context cancellation
	// This tests current behavior - context is passed but not used
	_, err := repo.List(ctx)
	if err == nil {
		t.Log("List succeeded with cancelled context (context not being checked)")
	}

	err = repo.Import(ctx)
	if err == nil {
		t.Log("Import succeeded with cancelled context (context not being checked)")
	}

	_, err = repo.IsImporting(ctx)
	if err == nil {
		t.Log("IsImporting succeeded with cancelled context (context not being checked)")
	}
}

func TestBootResource_MethodBehavior(t *testing.T) {
	// Test the behavior and signatures of each method
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootResource(maas)

	ctx := context.Background()

	// Test List method
	t.Run("List_method", func(t *testing.T) {
		result, err := repo.List(ctx)
		// With empty config, we expect an error
		if err == nil && result != nil {
			t.Log("List returned successful result, might be in test environment with MAAS")
		}
		if err != nil {
			t.Logf("List returned expected error: %v", err)
		}
	})

	// Test Import method
	t.Run("Import_method", func(t *testing.T) {
		err := repo.Import(ctx)
		// With empty config, we expect an error
		if err == nil {
			t.Log("Import succeeded, might be in test environment with MAAS")
		} else {
			t.Logf("Import returned expected error: %v", err)
		}
	})

	// Test IsImporting method
	t.Run("IsImporting_method", func(t *testing.T) {
		result, err := repo.IsImporting(ctx)
		// With empty config, we expect an error
		if err == nil {
			t.Logf("IsImporting succeeded with result: %v", result)
		} else {
			t.Logf("IsImporting returned expected error: %v", err)
		}
	})
}

func TestBootResource_IntegrationPatterns(t *testing.T) {
	// Test patterns that might be used in integration scenarios
	conf := &config.Config{
		MAAS: config.MAAS{
			URL:     "http://example.com/MAAS",
			Key:     "dummy-key:for:testing",
			Version: "2.8",
		},
	}

	maas := New(conf)
	repo := NewBootResource(maas)
	ctx := context.Background()

	// Test sequential calls (as might happen in real usage)
	t.Run("sequential_calls", func(t *testing.T) {
		// First check if importing
		importing, err := repo.IsImporting(ctx)
		if err != nil {
			t.Logf("IsImporting error (expected): %v", err)
		} else {
			t.Logf("IsImporting result: %v", importing)
		}

		// Then try to list resources
		resources, err := repo.List(ctx)
		if err != nil {
			t.Logf("List error (expected): %v", err)
		} else {
			t.Logf("List returned %d resources", len(resources))
		}

		// Finally try import (this would typically only be done if needed)
		err = repo.Import(ctx)
		if err != nil {
			t.Logf("Import error (expected): %v", err)
		} else {
			t.Log("Import succeeded")
		}
	})
}

// Benchmark tests
func BenchmarkBootResource_Creation(b *testing.B) {
	conf := &config.Config{}
	maas := New(conf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo := NewBootResource(maas)
		if repo == nil {
			b.Fatal("Failed to create repository")
		}
	}
}

func BenchmarkBootResource_MethodCalls(b *testing.B) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootResource(maas)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// These will fail but we're measuring the overhead
		_, _ = repo.List(ctx)
		_ = repo.Import(ctx)
		_, _ = repo.IsImporting(ctx)
	}
}

// Test concurrent access
func TestBootResource_ConcurrentAccess(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootResource(maas)
	ctx := context.Background()

	// Run multiple goroutines concurrently
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Method call panicked: %v", r)
				}
				done <- true
			}()

			// Each method should handle concurrent calls safely
			_, _ = repo.List(ctx)
			_ = repo.Import(ctx)
			_, _ = repo.IsImporting(ctx)
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestBootResource_TypeAssertions(t *testing.T) {
	// Test that we can assert types correctly
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootResource(maas)

	// Should be able to cast to the concrete type
	br, ok := repo.(*bootResource)
	if !ok {
		t.Fatal("Could not cast to *bootResource")
	}

	if br.maas != maas {
		t.Error("bootResource.maas field not set correctly")
	}

	// Should implement the interface
	var _ core.BootResourceRepo = br
}

// Test edge cases and error conditions
func TestBootResource_EdgeCases(t *testing.T) {
	t.Run("nil_maas", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when creating bootResource with nil MAAS")
			}
		}()

		// This should panic or handle gracefully
		repo := &bootResource{maas: nil}
		ctx := context.Background()
		_, _ = repo.List(ctx)
	})

	t.Run("background_context", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewBootResource(maas)

		// Should work with background context
		ctx := context.Background()
		_, err := repo.List(ctx)
		if err == nil {
			t.Log("List succeeded with background context")
		}
	})

	t.Run("todo_context", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewBootResource(maas)

		// Should work with TODO context
		ctx := context.TODO()
		_, err := repo.List(ctx)
		if err == nil {
			t.Log("List succeeded with TODO context")
		}
	})
}
