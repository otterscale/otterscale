package juju

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/connector"

	"github.com/openhdc/openhdc/internal/env"
)

type Juju = api.Connection

// Default configuration values
const (
	defaultControllerAddress = "localhost:17070"
	defaultUsername          = "admin"
)

// NewConfig creates a new Juju configuration from environment variables
func NewConfig() (*connector.SimpleConfig, error) {
	config := &connector.SimpleConfig{
		ControllerAddresses: strings.Split(env.GetOrDefault(env.OPENHDC_JUJU_CONTROLLER_ADDRESSES, defaultControllerAddress), ","),
		Username:            env.GetOrDefault(env.OPENHDC_JUJU_USERNAME, defaultUsername),
		Password:            env.GetOrDefault(env.OPENHDC_JUJU_PASSWORD, ""),
	}

	if caCertPath := os.Getenv(env.OPENHDC_JUJU_CACERT_PATH); caCertPath != "" {
		caCert, err := loadCACert(caCertPath)
		if err != nil {
			return nil, err
		}
		config.CACert = caCert
	}

	return config, nil
}

// New creates a new Juju API connection
func New(cfg *connector.SimpleConfig) (Juju, error) {
	sc, err := connector.NewSimple(*cfg)
	if err != nil {
		return nil, err
	}

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return sc.Connect(ctx)
}

// newConnection creates a new Juju connection with the specified model UUID
func newConnection(modelUUID string) (Juju, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}
	cfg.ModelUUID = modelUUID
	return New(cfg)
}

// loadCACert loads the CA certificate from a file
func loadCACert(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func checkConnection(ctx context.Context, conn Juju) (Juju, error) {
	if !conn.IsBroken(ctx) {
		return conn, nil
	}
	return newConnection("")
}
