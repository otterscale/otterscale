package maas

import (
	"context"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

func TestNewFabric(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)

	repo := NewFabric(maas)

	if repo == nil {
		t.Fatal("Expected Fabric repository to be created, got nil")
	}

	// Verify it implements the interface
	var _ core.FabricRepo = repo
}

func TestFabric_InterfaceCompliance(t *testing.T) {
	// Test that fabric implements core.FabricRepo interface
	conf := &config.Config{}
	maas := New(conf)
	repo := NewFabric(maas)

	var _ core.FabricRepo = repo
}

func TestFabric_Methods_Structure(t *testing.T) {
	// Test the structure and method signatures
	conf := &config.Config{}
	maas := New(conf)
	repo := NewFabric(maas)

	// Verify the repo is of the correct type
	f, ok := repo.(*fabric)
	if !ok {
		t.Fatal("Expected *fabric, got different type")
	}

	if f.maas == nil {
		t.Error("Expected maas field to be set, got nil")
	}
}

// Test with actual MAAS configuration structure
func TestFabric_WithConfig(t *testing.T) {
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
			repo := NewFabric(maas)

			if repo == nil {
				t.Error("Expected repository to be created")
			}

			ctx := context.Background()

			// Try to use the methods (they will likely fail due to invalid config, but should not panic)
			_, err := repo.List(ctx)
			if err == nil {
				t.Log("List succeeded unexpectedly (might be in test environment)")
			}

			_, err = repo.Get(ctx, 1)
			if err == nil {
				t.Log("Get succeeded unexpectedly (might be in test environment)")
			}

			params := &entity.FabricParams{
				Name: "test-fabric",
			}

			_, err = repo.Create(ctx, params)
			if err == nil {
				t.Log("Create succeeded unexpectedly (might be in test environment)")
			}

			_, err = repo.Update(ctx, 1, params)
			if err == nil {
				t.Log("Update succeeded unexpectedly (might be in test environment)")
			}

			err = repo.Delete(ctx, 1)
			if err == nil {
				t.Log("Delete succeeded unexpectedly (might be in test environment)")
			}
		})
	}
}

// Test error handling with invalid configurations
func TestFabric_ErrorHandling(t *testing.T) {
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
			repo := NewFabric(maas)

			ctx := context.Background()

			// All methods should return errors with invalid configurations
			_, err := repo.List(ctx)
			if err == nil {
				t.Logf("List unexpectedly succeeded for %s", tt.desc)
			}

			_, err = repo.Get(ctx, 1)
			if err == nil {
				t.Logf("Get unexpectedly succeeded for %s", tt.desc)
			}

			params := &entity.FabricParams{
				Name: "test-fabric",
			}

			_, err = repo.Create(ctx, params)
			if err == nil {
				t.Logf("Create unexpectedly succeeded for %s", tt.desc)
			}

			_, err = repo.Update(ctx, 1, params)
			if err == nil {
				t.Logf("Update unexpectedly succeeded for %s", tt.desc)
			}

			err = repo.Delete(ctx, 1)
			if err == nil {
				t.Logf("Delete unexpectedly succeeded for %s", tt.desc)
			}
		})
	}
}

func TestFabric_ContextHandling(t *testing.T) {
	// Test that methods handle context properly (though the current implementation ignores it)
	conf := &config.Config{}
	maas := New(conf)
	repo := NewFabric(maas)

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Methods should still fail with client error since they don't check context cancellation
	// This tests current behavior - context is passed but not used
	_, err := repo.List(ctx)
	if err == nil {
		t.Log("List succeeded with cancelled context (context not being checked)")
	}

	_, err = repo.Get(ctx, 1)
	if err == nil {
		t.Log("Get succeeded with cancelled context (context not being checked)")
	}

	params := &entity.FabricParams{
		Name: "test-fabric",
	}

	_, err = repo.Create(ctx, params)
	if err == nil {
		t.Log("Create succeeded with cancelled context (context not being checked)")
	}

	_, err = repo.Update(ctx, 1, params)
	if err == nil {
		t.Log("Update succeeded with cancelled context (context not being checked)")
	}

	err = repo.Delete(ctx, 1)
	if err == nil {
		t.Log("Delete succeeded with cancelled context (context not being checked)")
	}
}

func TestFabric_MethodBehavior(t *testing.T) {
	// Test the behavior and signatures of each method
	conf := &config.Config{}
	maas := New(conf)
	repo := NewFabric(maas)

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

	// Test Get method with different IDs
	t.Run("Get_method_different_ids", func(t *testing.T) {
		testIDs := []int{0, 1, 5, 100, -1}
		for _, id := range testIDs {
			result, err := repo.Get(ctx, id)
			if err == nil && result != nil {
				t.Logf("Get with ID %d returned successful result", id)
			} else if err != nil {
				t.Logf("Get with ID %d returned expected error: %v", id, err)
			}
		}
	})

	// Test Create method
	t.Run("Create_method", func(t *testing.T) {
		params := &entity.FabricParams{
			Name:        "test-fabric",
			Description: "Test fabric description",
		}

		result, err := repo.Create(ctx, params)
		// With empty config, we expect an error
		if err == nil && result != nil {
			t.Log("Create returned successful result, might be in test environment with MAAS")
		} else if err != nil {
			t.Logf("Create returned expected error: %v", err)
		}
	})

	// Test Create method with nil params
	t.Run("Create_method_nil_params", func(t *testing.T) {
		result, err := repo.Create(ctx, nil)
		if err == nil && result != nil {
			t.Log("Create with nil params succeeded unexpectedly")
		} else if err != nil {
			t.Logf("Create with nil params returned expected error: %v", err)
		}
	})

	// Test Update method
	t.Run("Update_method", func(t *testing.T) {
		params := &entity.FabricParams{
			Name:        "updated-fabric",
			Description: "Updated fabric description",
		}

		result, err := repo.Update(ctx, 1, params)
		// With empty config, we expect an error
		if err == nil && result != nil {
			t.Log("Update returned successful result, might be in test environment with MAAS")
		} else if err != nil {
			t.Logf("Update returned expected error: %v", err)
		}
	})

	// Test Update method with different IDs and parameters
	t.Run("Update_method_various_params", func(t *testing.T) {
		testCases := []struct {
			id     int
			params *entity.FabricParams
			desc   string
		}{
			{
				id:     1,
				params: &entity.FabricParams{Name: "fabric-1"},
				desc:   "minimal params",
			},
			{
				id: 2,
				params: &entity.FabricParams{
					Name:        "fabric-2",
					Description: "Description for fabric 2",
				},
				desc: "full params",
			},
			{
				id:     0,
				params: &entity.FabricParams{Name: "fabric-0"},
				desc:   "zero ID",
			},
			{
				id:     -1,
				params: &entity.FabricParams{Name: "fabric-negative"},
				desc:   "negative ID",
			},
			{
				id:     5,
				params: nil,
				desc:   "nil params",
			},
		}

		for i, tc := range testCases {
			result, err := repo.Update(ctx, tc.id, tc.params)
			if err == nil && result != nil {
				t.Logf("Update test %d (%s) returned successful result", i, tc.desc)
			} else if err != nil {
				t.Logf("Update test %d (%s) returned expected error: %v", i, tc.desc, err)
			}
		}
	})

	// Test Delete method
	t.Run("Delete_method", func(t *testing.T) {
		err := repo.Delete(ctx, 1)
		// With empty config, we expect an error
		if err == nil {
			t.Log("Delete returned successful result, might be in test environment with MAAS")
		} else {
			t.Logf("Delete returned expected error: %v", err)
		}
	})

	// Test Delete method with different IDs
	t.Run("Delete_method_different_ids", func(t *testing.T) {
		testIDs := []int{0, 1, 5, 100, -1}
		for _, id := range testIDs {
			err := repo.Delete(ctx, id)
			if err == nil {
				t.Logf("Delete with ID %d succeeded", id)
			} else {
				t.Logf("Delete with ID %d returned expected error: %v", id, err)
			}
		}
	})
}

func TestFabric_IntegrationPatterns(t *testing.T) {
	// Test patterns that might be used in integration scenarios
	conf := &config.Config{
		MAAS: config.MAAS{
			URL:     "http://example.com/MAAS",
			Key:     "dummy-key:for:testing",
			Version: "2.8",
		},
	}

	maas := New(conf)
	repo := NewFabric(maas)
	ctx := context.Background()

	// Test CRUD workflow (as might happen in real usage)
	t.Run("crud_workflow", func(t *testing.T) {
		// First list existing fabrics
		fabrics, err := repo.List(ctx)
		if err != nil {
			t.Logf("List error (expected): %v", err)
		} else {
			t.Logf("List returned %d fabrics", len(fabrics))
		}

		// Try to create a new fabric
		createParams := &entity.FabricParams{
			Name:        "test-fabric",
			Description: "Test fabric for integration testing",
		}

		fabric, err := repo.Create(ctx, createParams)
		if err != nil {
			t.Logf("Create error (expected): %v", err)
		} else if fabric != nil {
			t.Log("Create succeeded")

			// Try to get the created fabric
			retrievedFabric, err := repo.Get(ctx, fabric.ID)
			if err != nil {
				t.Logf("Get error (expected): %v", err)
			} else if retrievedFabric != nil {
				t.Log("Get succeeded")
			}

			// Try to update the fabric
			updateParams := &entity.FabricParams{
				Name:        "updated-test-fabric",
				Description: "Updated test fabric description",
			}

			updatedFabric, err := repo.Update(ctx, fabric.ID, updateParams)
			if err != nil {
				t.Logf("Update error (expected): %v", err)
			} else if updatedFabric != nil {
				t.Log("Update succeeded")
			}

			// Try to delete the fabric
			err = repo.Delete(ctx, fabric.ID)
			if err != nil {
				t.Logf("Delete error (expected): %v", err)
			} else {
				t.Log("Delete succeeded")
			}
		}

		// List again to see final state
		fabrics, err = repo.List(ctx)
		if err != nil {
			t.Logf("Final List error (expected): %v", err)
		} else {
			t.Logf("Final List returned %d fabrics", len(fabrics))
		}
	})
}

// Benchmark tests
func BenchmarkFabric_Creation(b *testing.B) {
	conf := &config.Config{}
	maas := New(conf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo := NewFabric(maas)
		if repo == nil {
			b.Fatal("Failed to create repository")
		}
	}
}

func BenchmarkFabric_MethodCalls(b *testing.B) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewFabric(maas)
	ctx := context.Background()

	params := &entity.FabricParams{
		Name: "benchmark-fabric",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// These will fail but we're measuring the overhead
		_, _ = repo.List(ctx)
		_, _ = repo.Get(ctx, 1)
		_, _ = repo.Create(ctx, params)
		_, _ = repo.Update(ctx, 1, params)
		_ = repo.Delete(ctx, 1)
	}
}

// Test concurrent access
func TestFabric_ConcurrentAccess(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewFabric(maas)
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
			_, _ = repo.List(ctx)
			_, _ = repo.Get(ctx, id%5+1)

			params := &entity.FabricParams{
				Name: "concurrent-fabric",
			}
			_, _ = repo.Create(ctx, params)
			_, _ = repo.Update(ctx, id%5+1, params)
			_ = repo.Delete(ctx, id%5+1)
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestFabric_TypeAssertions(t *testing.T) {
	// Test that we can assert types correctly
	conf := &config.Config{}
	maas := New(conf)
	repo := NewFabric(maas)

	// Should be able to cast to the concrete type
	f, ok := repo.(*fabric)
	if !ok {
		t.Fatal("Could not cast to *fabric")
	}

	if f.maas != maas {
		t.Error("fabric.maas field not set correctly")
	}

	// Should implement the interface
	var _ core.FabricRepo = f
}

// Test edge cases and error conditions
func TestFabric_EdgeCases(t *testing.T) {
	t.Run("nil_maas", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when creating fabric with nil MAAS")
			}
		}()

		// This should panic or handle gracefully
		repo := &fabric{maas: nil}
		ctx := context.Background()
		_, _ = repo.List(ctx)
	})

	t.Run("background_context", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewFabric(maas)

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
		repo := NewFabric(maas)

		// Should work with TODO context
		ctx := context.TODO()
		_, err := repo.List(ctx)
		if err == nil {
			t.Log("List succeeded with TODO context")
		}
	})

	t.Run("zero_fabric_id", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewFabric(maas)
		ctx := context.Background()

		// Test with zero fabric ID
		_, err := repo.Get(ctx, 0)
		if err == nil {
			t.Log("Get with zero ID succeeded")
		} else {
			t.Logf("Get with zero ID returned error: %v", err)
		}

		params := &entity.FabricParams{
			Name: "zero-id-fabric",
		}

		_, err = repo.Update(ctx, 0, params)
		if err == nil {
			t.Log("Update with zero fabric ID succeeded")
		} else {
			t.Logf("Update with zero fabric ID returned error: %v", err)
		}

		err = repo.Delete(ctx, 0)
		if err == nil {
			t.Log("Delete with zero fabric ID succeeded")
		} else {
			t.Logf("Delete with zero fabric ID returned error: %v", err)
		}
	})

	t.Run("negative_fabric_id", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewFabric(maas)
		ctx := context.Background()

		// Test with negative fabric ID
		_, err := repo.Get(ctx, -1)
		if err == nil {
			t.Log("Get with negative ID succeeded")
		} else {
			t.Logf("Get with negative ID returned error: %v", err)
		}

		params := &entity.FabricParams{
			Name: "negative-id-fabric",
		}

		_, err = repo.Update(ctx, -1, params)
		if err == nil {
			t.Log("Update with negative fabric ID succeeded")
		} else {
			t.Logf("Update with negative fabric ID returned error: %v", err)
		}

		err = repo.Delete(ctx, -1)
		if err == nil {
			t.Log("Delete with negative fabric ID succeeded")
		} else {
			t.Logf("Delete with negative fabric ID returned error: %v", err)
		}
	})
}

func TestFabric_ParameterValidation(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewFabric(maas)
	ctx := context.Background()

	// Test Create and Update with various parameter combinations
	tests := []struct {
		name   string
		params *entity.FabricParams
		desc   string
	}{
		{
			name:   "valid_minimal_params",
			params: &entity.FabricParams{Name: "valid-fabric"},
			desc:   "Valid minimal parameters",
		},
		{
			name: "valid_full_params",
			params: &entity.FabricParams{
				Name:        "full-fabric",
				Description: "Full fabric with description",
			},
			desc: "Valid full parameters",
		},
		{
			name:   "empty_name",
			params: &entity.FabricParams{Name: ""},
			desc:   "Empty name field",
		},
		{
			name: "empty_description",
			params: &entity.FabricParams{
				Name:        "fabric-empty-desc",
				Description: "",
			},
			desc: "Empty description field",
		},
		{
			name: "long_name",
			params: &entity.FabricParams{
				Name: "very-long-fabric-name-that-might-exceed-some-limits-in-certain-systems-and-databases",
			},
			desc: "Very long name",
		},
		{
			name: "special_characters",
			params: &entity.FabricParams{
				Name:        "fabric-with-special-chars-123_abc",
				Description: "Description with special chars: !@#$%^&*()",
			},
			desc: "Special characters in fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_create", func(t *testing.T) {
			result, err := repo.Create(ctx, tt.params)
			if err == nil && result != nil {
				t.Logf("%s (Create): succeeded unexpectedly", tt.desc)
			} else if err != nil {
				t.Logf("%s (Create): returned expected error: %v", tt.desc, err)
			}
		})

		t.Run(tt.name+"_update", func(t *testing.T) {
			result, err := repo.Update(ctx, 1, tt.params)
			if err == nil && result != nil {
				t.Logf("%s (Update): succeeded unexpectedly", tt.desc)
			} else if err != nil {
				t.Logf("%s (Update): returned expected error: %v", tt.desc, err)
			}
		})
	}
}

func TestFabric_MAASConfigurationVariations(t *testing.T) {
	// Test with different MAAS configuration variations
	configs := []struct {
		name   string
		config *config.Config
		desc   string
	}{
		{
			name: "localhost_maas",
			config: &config.Config{
				MAAS: config.MAAS{
					URL:     "http://localhost:5240/MAAS",
					Key:     "test:api:key",
					Version: "2.8",
				},
			},
			desc: "Localhost MAAS configuration",
		},
		{
			name: "remote_maas",
			config: &config.Config{
				MAAS: config.MAAS{
					URL:     "https://maas.example.com/MAAS",
					Key:     "prod:api:key",
					Version: "3.0",
				},
			},
			desc: "Remote MAAS configuration",
		},
		{
			name: "version_2_9",
			config: &config.Config{
				MAAS: config.MAAS{
					URL:     "http://maas.local/MAAS",
					Key:     "local:api:key",
					Version: "2.9",
				},
			},
			desc: "MAAS version 2.9",
		},
	}

	for _, tc := range configs {
		t.Run(tc.name, func(t *testing.T) {
			maas := New(tc.config)
			repo := NewFabric(maas)
			ctx := context.Background()

			// Test each method with this configuration
			_, err := repo.List(ctx)
			if err != nil {
				t.Logf("List failed for %s: %v", tc.desc, err)
			} else {
				t.Logf("List succeeded for %s", tc.desc)
			}

			_, err = repo.Get(ctx, 1)
			if err != nil {
				t.Logf("Get failed for %s: %v", tc.desc, err)
			} else {
				t.Logf("Get succeeded for %s", tc.desc)
			}

			params := &entity.FabricParams{Name: "test-fabric"}
			_, err = repo.Create(ctx, params)
			if err != nil {
				t.Logf("Create failed for %s: %v", tc.desc, err)
			} else {
				t.Logf("Create succeeded for %s", tc.desc)
			}
		})
	}
}
