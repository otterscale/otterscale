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

const (
	defaultControllerAddress = "localhost:17070"
	defaultUsername          = "admin"
)

type JujuMap map[string]api.Connection

func NewJujuMap() (JujuMap, error) {
	return JujuMap{}, nil
}

func (m JujuMap) Get(ctx context.Context, uuid string) (api.Connection, error) {
	if conn, ok := m[uuid]; ok {
		if !conn.IsBroken(ctx) {
			return conn, nil
		}
		conn.Close()
	}

	conn, err := newConnection(uuid)
	if err != nil {
		return nil, err
	}
	m[uuid] = conn

	return conn, nil
}

func newConnection(uuid string) (api.Connection, error) {
	cfg := &connector.SimpleConfig{
		ModelUUID:           uuid,
		ControllerAddresses: strings.Split(env.GetOrDefault(env.OPENHDC_JUJU_CONTROLLER_ADDRESSES, defaultControllerAddress), ","),
		Username:            env.GetOrDefault(env.OPENHDC_JUJU_USERNAME, defaultUsername),
		Password:            env.GetOrDefault(env.OPENHDC_JUJU_PASSWORD, ""),
	}

	if path := os.Getenv(env.OPENHDC_JUJU_CACERT_PATH); path != "" {
		caCert, err := loadCACert(path)
		if err != nil {
			return nil, err
		}
		cfg.CACert = caCert
	}

	sc, err := connector.NewSimple(*cfg)
	if err != nil {
		return nil, err
	}

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return sc.Connect(ctx)
}

func loadCACert(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
