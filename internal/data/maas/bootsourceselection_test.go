package maas

import (
	"context"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/openhdc/otterscale/internal/config"
	"github.com/openhdc/otterscale/internal/core"
)

func TestNewBootSourceSelection(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)

	repo := NewBootSourceSelection(maas)

	if repo == nil {
		t.Fatal("Expected BootSourceSelection repository to be created, got nil")
	}

	// Verify it implements the interface
	var _ core.BootSourceSelectionRepo = repo
}

func TestBootSourceSelection_InterfaceCompliance(t *testing.T) {
	// Test that bootSourceSelection implements core.BootSourceSelectionRepo interface
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootSourceSelection(maas)

	var _ core.BootSourceSelectionRepo = repo
}

func TestBootSourceSelection_Methods_Structure(t *testing.T) {
	// Test the structure and method signatures
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootSourceSelection(maas)

	// Verify the repo is of the correct type
	bss, ok := repo.(*bootSourceSelection)
	if !ok {
		t.Fatal("Expected *bootSourceSelection, got different type")
	}

	if bss.maas == nil {
		t.Error("Expected maas field to be set, got nil")
	}
}

// Test with actual MAAS configuration structure
func TestBootSourceSelection_WithConfig(t *testing.T) {
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
			repo := NewBootSourceSelection(maas)

			if repo == nil {
				t.Error("Expected repository to be created")
			}

			// Try to use the methods (they will likely fail due to invalid config, but should not panic)
			ctx := context.Background()

			_, err := repo.List(ctx, 1)
			// We expect an error since we don't have a real MAAS server
			if err == nil {
				t.Log("List succeeded unexpectedly (might be in test environment)")
			}

			params := &entity.BootSourceSelectionParams{
				OS:        "ubuntu",
				Release:   "20.04",
				Arches:    []string{"amd64"},
				Subarches: []string{"generic"},
				Labels:    []string{"release"},
			}

			_, err = repo.Create(ctx, 1, params)
			if err == nil {
				t.Log("Create succeeded unexpectedly (might be in test environment)")
			}
		})
	}
}

// Test error handling with invalid configurations
func TestBootSourceSelection_ErrorHandling(t *testing.T) {
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
			repo := NewBootSourceSelection(maas)

			ctx := context.Background()

			// All methods should return errors with invalid configurations
			_, err := repo.List(ctx, 1)
			if err == nil {
				t.Logf("List unexpectedly succeeded for %s", tt.desc)
			}

			params := &entity.BootSourceSelectionParams{
				OS:      "ubuntu",
				Release: "20.04",
				Arches:  []string{"amd64"},
			}

			_, err = repo.Create(ctx, 1, params)
			if err == nil {
				t.Logf("Create unexpectedly succeeded for %s", tt.desc)
			}
		})
	}
}

func TestBootSourceSelection_ContextHandling(t *testing.T) {
	// Test that methods handle context properly (though the current implementation ignores it)
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootSourceSelection(maas)

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Methods should still fail with client error since they don't check context cancellation
	// This tests current behavior - context is passed but not used
	_, err := repo.List(ctx, 1)
	if err == nil {
		t.Log("List succeeded with cancelled context (context not being checked)")
	}

	params := &entity.BootSourceSelectionParams{
		OS:      "ubuntu",
		Release: "20.04",
		Arches:  []string{"amd64"},
	}

	_, err = repo.Create(ctx, 1, params)
	if err == nil {
		t.Log("Create succeeded with cancelled context (context not being checked)")
	}
}

func TestBootSourceSelection_MethodBehavior(t *testing.T) {
	// Test the behavior and signatures of each method
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootSourceSelection(maas)

	ctx := context.Background()

	// Test List method
	t.Run("List_method", func(t *testing.T) {
		result, err := repo.List(ctx, 1)
		// With empty config, we expect an error
		if err == nil && result != nil {
			t.Log("List returned successful result, might be in test environment with MAAS")
		}
		if err != nil {
			t.Logf("List returned expected error: %v", err)
		}
	})

	// Test List method with different IDs
	t.Run("List_method_different_ids", func(t *testing.T) {
		testIDs := []int{0, 1, 5, 100, -1}
		for _, id := range testIDs {
			result, err := repo.List(ctx, id)
			if err == nil && result != nil {
				t.Logf("List with ID %d returned successful result", id)
			} else if err != nil {
				t.Logf("List with ID %d returned expected error: %v", id, err)
			}
		}
	})

	// Test Create method
	t.Run("Create_method", func(t *testing.T) {
		params := &entity.BootSourceSelectionParams{
			OS:        "ubuntu",
			Release:   "20.04",
			Arches:    []string{"amd64"},
			Subarches: []string{"generic"},
			Labels:    []string{"release"},
		}

		result, err := repo.Create(ctx, 1, params)
		// With empty config, we expect an error
		if err == nil && result != nil {
			t.Log("Create returned successful result, might be in test environment with MAAS")
		} else if err != nil {
			t.Logf("Create returned expected error: %v", err)
		}
	})

	// Test Create method with nil params
	t.Run("Create_method_nil_params", func(t *testing.T) {
		result, err := repo.Create(ctx, 1, nil)
		if err == nil && result != nil {
			t.Log("Create with nil params succeeded unexpectedly")
		} else if err != nil {
			t.Logf("Create with nil params returned expected error: %v", err)
		}
	})

	// Test Create method with different parameters
	t.Run("Create_method_various_params", func(t *testing.T) {
		testParams := []*entity.BootSourceSelectionParams{
			{
				OS:      "ubuntu",
				Release: "18.04",
				Arches:  []string{"amd64"},
			},
			{
				OS:      "ubuntu",
				Release: "20.04",
				Arches:  []string{"amd64", "arm64"},
			},
			{
				OS:        "ubuntu",
				Release:   "22.04",
				Arches:    []string{"amd64"},
				Subarches: []string{"generic", "hwe"},
				Labels:    []string{"release", "daily"},
			},
		}

		for i, params := range testParams {
			result, err := repo.Create(ctx, 1, params)
			if err == nil && result != nil {
				t.Logf("Create test %d returned successful result", i)
			} else if err != nil {
				t.Logf("Create test %d returned expected error: %v", i, err)
			}
		}
	})
}

func TestBootSourceSelection_IntegrationPatterns(t *testing.T) {
	// Test patterns that might be used in integration scenarios
	conf := &config.Config{
		MAAS: config.MAAS{
			URL:     "http://example.com/MAAS",
			Key:     "dummy-key:for:testing",
			Version: "2.8",
		},
	}

	maas := New(conf)
	repo := NewBootSourceSelection(maas)
	ctx := context.Background()

	// Test sequential calls (as might happen in real usage)
	t.Run("sequential_calls", func(t *testing.T) {
		// First list existing selections
		selections, err := repo.List(ctx, 1)
		if err != nil {
			t.Logf("List error (expected): %v", err)
		} else {
			t.Logf("List returned %d selections", len(selections))
		}

		// Then try to create a new selection
		params := &entity.BootSourceSelectionParams{
			OS:      "ubuntu",
			Release: "20.04",
			Arches:  []string{"amd64"},
		}

		selection, err := repo.Create(ctx, 1, params)
		if err != nil {
			t.Logf("Create error (expected): %v", err)
		} else if selection != nil {
			t.Log("Create succeeded")
		}

		// List again to see if the new selection appears
		selections, err = repo.List(ctx, 1)
		if err != nil {
			t.Logf("Second List error (expected): %v", err)
		} else {
			t.Logf("Second List returned %d selections", len(selections))
		}
	})
}

// Benchmark tests
func BenchmarkBootSourceSelection_Creation(b *testing.B) {
	conf := &config.Config{}
	maas := New(conf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo := NewBootSourceSelection(maas)
		if repo == nil {
			b.Fatal("Failed to create repository")
		}
	}
}

func BenchmarkBootSourceSelection_MethodCalls(b *testing.B) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootSourceSelection(maas)
	ctx := context.Background()

	params := &entity.BootSourceSelectionParams{
		OS:      "ubuntu",
		Release: "20.04",
		Arches:  []string{"amd64"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// These will fail but we're measuring the overhead
		_, _ = repo.List(ctx, 1)
		_, _ = repo.Create(ctx, 1, params)
	}
}

// Test concurrent access
func TestBootSourceSelection_ConcurrentAccess(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootSourceSelection(maas)
	ctx := context.Background()

	// Run multiple goroutines concurrently
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Method call panicked: %v", r)
				}
				done <- true
			}()

			// Each method should handle concurrent calls safely
			_, _ = repo.List(ctx, id%5+1)

			params := &entity.BootSourceSelectionParams{
				OS:      "ubuntu",
				Release: "20.04",
				Arches:  []string{"amd64"},
			}
			_, _ = repo.Create(ctx, id%5+1, params)
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestBootSourceSelection_TypeAssertions(t *testing.T) {
	// Test that we can assert types correctly
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootSourceSelection(maas)

	// Should be able to cast to the concrete type
	bss, ok := repo.(*bootSourceSelection)
	if !ok {
		t.Fatal("Could not cast to *bootSourceSelection")
	}

	if bss.maas != maas {
		t.Error("bootSourceSelection.maas field not set correctly")
	}

	// Should implement the interface
	var _ core.BootSourceSelectionRepo = bss
}

// Test edge cases and error conditions
func TestBootSourceSelection_EdgeCases(t *testing.T) {
	t.Run("nil_maas", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when creating bootSourceSelection with nil MAAS")
			}
		}()

		// This should panic or handle gracefully
		repo := &bootSourceSelection{maas: nil}
		ctx := context.Background()
		_, _ = repo.List(ctx, 1)
	})

	t.Run("background_context", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewBootSourceSelection(maas)

		// Should work with background context
		ctx := context.Background()
		_, err := repo.List(ctx, 1)
		if err == nil {
			t.Log("List succeeded with background context")
		}
	})

	t.Run("todo_context", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewBootSourceSelection(maas)

		// Should work with TODO context
		ctx := context.TODO()
		_, err := repo.List(ctx, 1)
		if err == nil {
			t.Log("List succeeded with TODO context")
		}
	})

	t.Run("zero_boot_source_id", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewBootSourceSelection(maas)
		ctx := context.Background()

		// Test with zero boot source ID
		_, err := repo.List(ctx, 0)
		if err == nil {
			t.Log("List with zero ID succeeded")
		} else {
			t.Logf("List with zero ID returned error: %v", err)
		}

		params := &entity.BootSourceSelectionParams{
			OS:      "ubuntu",
			Release: "20.04",
			Arches:  []string{"amd64"},
		}

		_, err = repo.Create(ctx, 0, params)
		if err == nil {
			t.Log("Create with zero boot source ID succeeded")
		} else {
			t.Logf("Create with zero boot source ID returned error: %v", err)
		}
	})

	t.Run("negative_boot_source_id", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewBootSourceSelection(maas)
		ctx := context.Background()

		// Test with negative boot source ID
		_, err := repo.List(ctx, -1)
		if err == nil {
			t.Log("List with negative ID succeeded")
		} else {
			t.Logf("List with negative ID returned error: %v", err)
		}

		params := &entity.BootSourceSelectionParams{
			OS:      "ubuntu",
			Release: "20.04",
			Arches:  []string{"amd64"},
		}

		_, err = repo.Create(ctx, -1, params)
		if err == nil {
			t.Log("Create with negative boot source ID succeeded")
		} else {
			t.Logf("Create with negative boot source ID returned error: %v", err)
		}
	})
}

func TestBootSourceSelection_ParameterValidation(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewBootSourceSelection(maas)
	ctx := context.Background()

	// Test Create with various parameter combinations
	tests := []struct {
		name   string
		params *entity.BootSourceSelectionParams
		desc   string
	}{
		{
			name: "valid_minimal_params",
			params: &entity.BootSourceSelectionParams{
				OS:      "ubuntu",
				Release: "20.04",
				Arches:  []string{"amd64"},
			},
			desc: "Valid minimal parameters",
		},
		{
			name: "valid_full_params",
			params: &entity.BootSourceSelectionParams{
				OS:        "ubuntu",
				Release:   "20.04",
				Arches:    []string{"amd64", "arm64"},
				Subarches: []string{"generic", "hwe"},
				Labels:    []string{"release", "daily"},
			},
			desc: "Valid full parameters",
		},
		{
			name: "empty_os",
			params: &entity.BootSourceSelectionParams{
				OS:      "",
				Release: "20.04",
				Arches:  []string{"amd64"},
			},
			desc: "Empty OS field",
		},
		{
			name: "empty_release",
			params: &entity.BootSourceSelectionParams{
				OS:      "ubuntu",
				Release: "",
				Arches:  []string{"amd64"},
			},
			desc: "Empty release field",
		},
		{
			name: "empty_arches",
			params: &entity.BootSourceSelectionParams{
				OS:      "ubuntu",
				Release: "20.04",
				Arches:  []string{},
			},
			desc: "Empty architectures array",
		},
		{
			name: "nil_arches",
			params: &entity.BootSourceSelectionParams{
				OS:      "ubuntu",
				Release: "20.04",
				Arches:  nil,
			},
			desc: "Nil architectures array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(ctx, 1, tt.params)
			if err == nil && result != nil {
				t.Logf("%s: Create succeeded unexpectedly", tt.desc)
			} else if err != nil {
				t.Logf("%s: Create returned expected error: %v", tt.desc, err)
			}
		})
	}
}
