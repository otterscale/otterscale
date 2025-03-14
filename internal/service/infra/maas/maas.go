package maas

import (
	"github.com/canonical/gomaasclient/client"

	"github.com/openhdc/openhdc/internal/env"
)

type MAAS = client.Client

// Default configuration values
const (
	defaultURL     = "http://localhost:5240/MAAS/"
	defaultVersion = "2.0"
)

// Config holds MAAS client configuration
type Config struct {
	URL     string
	Key     string
	Version string
}

// NewConfig creates a new Config with values from environment variables or defaults
func NewConfig() *Config {
	return &Config{
		URL:     env.GetOrDefault(env.OPENHDC_MAAS_API_URL, defaultURL),
		Key:     env.GetOrDefault(env.OPENHDC_MAAS_API_KEY, ""),
		Version: env.GetOrDefault(env.OPENHDC_MAAS_API_VERSION, defaultVersion),
	}
}

// New creates a new MAAS client with the configured parameters
func New(cfg *Config) (*MAAS, error) {
	return client.GetClient(cfg.URL, cfg.Key, cfg.Version)
}
