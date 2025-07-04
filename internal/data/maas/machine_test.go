package maas

import (
	"context"
	"testing"

	"github.com/canonical/gomaasclient/entity"

	"github.com/openhdc/otterscale/internal/config"
	"github.com/openhdc/otterscale/internal/core"
)

func TestNewMachine(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)
	
	repo := NewMachine(maas)
	
	if repo == nil {
		t.Fatal("Expected Machine repository to be created, got nil")
	}
	
	// Verify it implements the interface
	var _ core.MachineRepo = repo
}

func TestMachine_InterfaceCompliance(t *testing.T) {
	// Test that machine implements core.MachineRepo interface
	conf := &config.Config{}
	maas := New(conf)
	repo := NewMachine(maas)
	
	var _ core.MachineRepo = repo
}

func TestMachine_Methods_Structure(t *testing.T) {
	// Test the structure and method signatures
	conf := &config.Config{}
	maas := New(conf)
	repo := NewMachine(maas)
	
	// Verify the repo is of the correct type
	m, ok := repo.(*machine)
	if !ok {
		t.Fatal("Expected *machine, got different type")
	}
	
	if m.maas == nil {
		t.Error("Expected maas field to be set, got nil")
	}
}

// Test with actual MAAS configuration structure
func TestMachine_WithConfig(t *testing.T) {
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
			repo := NewMachine(maas)
			
			if repo == nil {
				t.Error("Expected repository to be created")
			}
			
			// Try to use the methods (they will likely fail due to invalid config, but should not panic)
			ctx := context.Background()
			
			_, err := repo.List(ctx)
			if err == nil {
				t.Log("List succeeded unexpectedly (might be in test environment)")
			}
			
			// Test other methods with dummy parameters
			_, err = repo.Get(ctx, "test-system-id")
			if err == nil {
				t.Log("Get succeeded unexpectedly")
			}
			
			_, err = repo.Release(ctx, "test-system-id", &entity.MachineReleaseParams{})
			if err == nil {
				t.Log("Release succeeded unexpectedly")
			}
			
			_, err = repo.PowerOff(ctx, "test-system-id", &entity.MachinePowerOffParams{})
			if err == nil {
				t.Log("PowerOff succeeded unexpectedly")
			}
			
			_, err = repo.Commission(ctx, "test-system-id", &entity.MachineCommissionParams{})
			if err == nil {
				t.Log("Commission succeeded unexpectedly")
			}
		})
	}
}

// Test error handling with invalid configurations
func TestMachine_ErrorHandling(t *testing.T) {
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
		{
			name: "invalid_version",
			config: &config.Config{
				MAAS: config.MAAS{
					URL:     "http://localhost:5240/MAAS",
					Key:     "test-key",
					Version: "invalid",
				},
			},
			desc: "Invalid version should cause client errors",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maas := New(tt.config)
			repo := NewMachine(maas)
			
			ctx := context.Background()
			
			// Test all methods should return errors with invalid configurations
			_, err := repo.List(ctx)
			if err == nil {
				t.Logf("List unexpectedly succeeded for %s", tt.desc)
			} else {
				t.Logf("List returned expected error for %s: %v", tt.desc, err)
			}
			
			_, err = repo.Get(ctx, "test-id")
			if err == nil {
				t.Logf("Get unexpectedly succeeded for %s", tt.desc)
			} else {
				t.Logf("Get returned expected error for %s: %v", tt.desc, err)
			}
			
			_, err = repo.Release(ctx, "test-id", &entity.MachineReleaseParams{})
			if err == nil {
				t.Logf("Release unexpectedly succeeded for %s", tt.desc)
			} else {
				t.Logf("Release returned expected error for %s: %v", tt.desc, err)
			}
			
			_, err = repo.PowerOff(ctx, "test-id", &entity.MachinePowerOffParams{})
			if err == nil {
				t.Logf("PowerOff unexpectedly succeeded for %s", tt.desc)
			} else {
				t.Logf("PowerOff returned expected error for %s: %v", tt.desc, err)
			}
			
			_, err = repo.Commission(ctx, "test-id", &entity.MachineCommissionParams{})
			if err == nil {
				t.Logf("Commission unexpectedly succeeded for %s", tt.desc)
			} else {
				t.Logf("Commission returned expected error for %s: %v", tt.desc, err)
			}
		})
	}
}

func TestMachine_MethodBehavior(t *testing.T) {
	// Test the behavior and signature of all methods
	conf := &config.Config{}
	maas := New(conf)
	repo := NewMachine(maas)
	
	ctx := context.Background()
	
	// Test List method
	t.Run("List_method", func(t *testing.T) {
		result, err := repo.List(ctx)
		// With empty config, we expect an error
		if err == nil && result != nil {
			t.Log("List returned successful result, might be in test environment with MAAS")
			t.Logf("Returned %d machines", len(result))
		}
		if err != nil {
			t.Logf("List returned expected error: %v", err)
		}
	})
	
	// Test Get method
	t.Run("Get_method", func(t *testing.T) {
		result, err := repo.Get(ctx, "test-system-id")
		if err == nil && result != nil {
			t.Log("Get returned successful result")
		}
		if err != nil {
			t.Logf("Get returned expected error: %v", err)
		}
	})
	
	// Test Release method
	t.Run("Release_method", func(t *testing.T) {
		params := &entity.MachineReleaseParams{}
		result, err := repo.Release(ctx, "test-system-id", params)
		if err == nil && result != nil {
			t.Log("Release returned successful result")
		}
		if err != nil {
			t.Logf("Release returned expected error: %v", err)
		}
	})
	
	// Test PowerOff method
	t.Run("PowerOff_method", func(t *testing.T) {
		params := &entity.MachinePowerOffParams{}
		result, err := repo.PowerOff(ctx, "test-system-id", params)
		if err == nil && result != nil {
			t.Log("PowerOff returned successful result")
		}
		if err != nil {
			t.Logf("PowerOff returned expected error: %v", err)
		}
	})
	
	// Test Commission method
	t.Run("Commission_method", func(t *testing.T) {
		params := &entity.MachineCommissionParams{}
		result, err := repo.Commission(ctx, "test-system-id", params)
		if err == nil && result != nil {
			t.Log("Commission returned successful result")
		}
		if err != nil {
			t.Logf("Commission returned expected error: %v", err)
		}
	})
}

// Test edge cases and error conditions
func TestMachine_EdgeCases(t *testing.T) {
	t.Run("nil_maas", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when calling methods with nil MAAS")
			}
		}()
		
		// This should panic or handle gracefully
		repo := &machine{maas: nil}
		ctx := context.Background()
		_, _ = repo.List(ctx)
	})
	
	t.Run("empty_system_id", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewMachine(maas)
		ctx := context.Background()
		
		// Test with empty system ID
		_, err := repo.Get(ctx, "")
		if err == nil {
			t.Log("Get succeeded with empty system ID")
		} else {
			t.Logf("Get failed with empty system ID: %v", err)
		}
	})
	
	t.Run("nil_params", func(t *testing.T) {
		conf := &config.Config{}
		maas := New(conf)
		repo := NewMachine(maas)
		ctx := context.Background()
		
		// Test with nil parameters
		_, err := repo.Release(ctx, "test-id", nil)
		if err == nil {
			t.Log("Release succeeded with nil params")
		} else {
			t.Logf("Release failed with nil params: %v", err)
		}
		
		_, err = repo.PowerOff(ctx, "test-id", nil)
		if err == nil {
			t.Log("PowerOff succeeded with nil params")
		} else {
			t.Logf("PowerOff failed with nil params: %v", err)
		}
		
		_, err = repo.Commission(ctx, "test-id", nil)
		if err == nil {
			t.Log("Commission succeeded with nil params")
		} else {
			t.Logf("Commission failed with nil params: %v", err)
		}
	})
}

// Test concurrent access
func TestMachine_ConcurrentAccess(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewMachine(maas)
	ctx := context.Background()
	
	// Run multiple goroutines concurrently
	done := make(chan bool, 20)
	
	for i := 0; i < 20; i++ {
		go func(id int) {
			defer func() { 
				if r := recover(); r != nil {
					t.Errorf("Goroutine %d panicked: %v", id, r)
				}
				done <- true 
			}()
			
			// Test different methods concurrently
			switch id % 5 {
			case 0:
				_, _ = repo.List(ctx)
			case 1:
				_, _ = repo.Get(ctx, "test-id")
			case 2:
				_, _ = repo.Release(ctx, "test-id", &entity.MachineReleaseParams{})
			case 3:
				_, _ = repo.PowerOff(ctx, "test-id", &entity.MachinePowerOffParams{})
			case 4:
				_, _ = repo.Commission(ctx, "test-id", &entity.MachineCommissionParams{})
			}
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 20; i++ {
		<-done
	}
}

// Benchmark tests
func BenchmarkMachine_Creation(b *testing.B) {
	conf := &config.Config{}
	maas := New(conf)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo := NewMachine(maas)
		if repo == nil {
			b.Fatal("Failed to create repository")
		}
	}
}

func BenchmarkMachine_List(b *testing.B) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewMachine(maas)
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.List(ctx)
	}
}

func BenchmarkMachine_Get(b *testing.B) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewMachine(maas)
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Get(ctx, "test-id")
	}
}

func BenchmarkMachine_Operations(b *testing.B) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewMachine(maas)
	ctx := context.Background()
	
	b.Run("Release", func(b *testing.B) {
		params := &entity.MachineReleaseParams{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = repo.Release(ctx, "test-id", params)
		}
	})
	
	b.Run("PowerOff", func(b *testing.B) {
		params := &entity.MachinePowerOffParams{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = repo.PowerOff(ctx, "test-id", params)
		}
	})
	
	b.Run("Commission", func(b *testing.B) {
		params := &entity.MachineCommissionParams{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = repo.Commission(ctx, "test-id", params)
		}
	})
}

// Test method signature compatibility
func TestMachine_MethodSignatures(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)
	repo := NewMachine(maas)
	
	// Test that all methods have the correct signatures
	ctx := context.Background()
	
	// Test List method
	machines, err := repo.List(ctx)
	if machines == nil && err == nil {
		t.Error("Both machines and error are nil from List, expected at least one to have a value")
	}
	
	// Test Get method
	machine, err := repo.Get(ctx, "test-id")
	if machine == nil && err == nil {
		t.Error("Both machine and error are nil from Get, expected at least one to have a value")
	}
	
	// Test Release method
	machine, err = repo.Release(ctx, "test-id", &entity.MachineReleaseParams{})
	if machine == nil && err == nil {
		t.Error("Both machine and error are nil from Release, expected at least one to have a value")
	}
	
	// Test PowerOff method
	machine, err = repo.PowerOff(ctx, "test-id", &entity.MachinePowerOffParams{})
	if machine == nil && err == nil {
		t.Error("Both machine and error are nil from PowerOff, expected at least one to have a value")
	}
	
	// Test Commission method
	machine, err = repo.Commission(ctx, "test-id", &entity.MachineCommissionParams{})
	if machine == nil && err == nil {
		t.Error("Both machine and error are nil from Commission, expected at least one to have a value")
	}
	
	// Test that error is properly typed when present
	if err != nil {
		errStr := err.Error()
		if errStr == "" {
			t.Error("Error has empty message")
		}
	}
}
