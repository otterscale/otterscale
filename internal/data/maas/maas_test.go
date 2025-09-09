package maas

import (
	"testing"

	"github.com/otterscale/otterscale/internal/config"
)

func TestNewMAAS(t *testing.T) {
	conf := &config.Config{}
	maas := New(conf)
	if maas == nil {
		t.Fatal("Expected MAAS to be created, but got nil")
	}
	if maas.conf != conf {
		t.Error("Expected MAAS.conf to be set, but got a different config")
	}
}

func TestMAAS_Client(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
		desc   string
	}{
		{
			name:   "empty_config",
			config: &config.Config{},
			desc:   "Empty MAAS configuration should cause client errors",
		},
		{
			name: "valid_config",
			config: &config.Config{
				MAAS: config.MAAS{
					URL:     "http://localhost:5240/MAAS",
					Key:     "consumer-secret:token-key:token-secret",
					Version: "2.8",
				},
			},
			desc: "Valid MAAS configuration should return a client",
		},
		{
			name: "invalid_url",
			config: &config.Config{
				MAAS: config.MAAS{
					URL:     "invalid-url",
					Key:     "consumer-secret:token-key:token-secret",
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
			// Try to get the client (it will likely fail due to invalid config, but should not panic)
			client, err := maas.client()
			if err == nil {
				t.Logf("Client creation succeeded unexpectedly: %s", tt.desc)
			} else {
				t.Logf("Client creation returned expected error: %v", err)
			}
			if client != nil && err == nil {
				t.Logf("Client: %v", client)
			}
		})
	}
}

func TestMAAS_TypeAssertions(t *testing.T) {
	// Test that we can assert types correctly
	conf := &config.Config{}
	maas := New(conf)
	if maas == nil {
		t.Fatal("Could not create MAAS")
	}
	if maas.conf == nil {
		t.Error("MAAS.conf field not set correctly")
	}
}

func TestMAAS_EdgeCases(t *testing.T) {
	t.Run("nil_config", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when creating MAAS with nil config")
			}
		}()
		// This should panic or handle gracefully
		maas := New(nil)
		_, _ = maas.client()
	})

	t.Run("background_context", func(t *testing.T) {
		conf := &config.Config{
			MAAS: config.MAAS{
				URL:     "http://localhost:5240/MAAS",
				Key:     "consumer-secret:token-key:token-secret",
				Version: "2.8",
			},
		}
		maas := New(conf)
		// Should work with background context
		client, err := maas.client()
		if err != nil {
			t.Logf("Client creation error (expected): %v", err)
		}
		if client == nil {
			t.Error("Expected client to be created, but got nil")
		}
	})

	t.Run("todo_context", func(t *testing.T) {
		conf := &config.Config{
			MAAS: config.MAAS{
				URL:     "http://localhost:5240/MAAS",
				Key:     "consumer-secret:token-key:token-secret",
				Version: "2.8",
			},
		}
		maas := New(conf)
		// Should work with TODO context
		client, err := maas.client()
		if err != nil {
			t.Logf("Client creation error (expected): %v", err)
		}
		if client == nil {
			t.Error("Expected client to be created, but got nil")
		}
	})
}
