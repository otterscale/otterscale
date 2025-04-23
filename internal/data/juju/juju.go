package juju

import (
	"os"
	"strings"
	"sync"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/connector"

	"github.com/openhdc/openhdc/internal/env"
)

const (
	defaultControllerAddress = "localhost:17070"
	defaultUsername          = "admin"
)

// map[string]api.Connection
type JujuMap struct {
	*sync.Map
}

func NewJujuMap() (JujuMap, error) {
	return JujuMap{&sync.Map{}}, nil
}

func (m JujuMap) Get(uuid string) (api.Connection, error) {
	if v, ok := m.Load(uuid); ok {
		conn := v.(api.Connection)
		if !conn.IsBroken() {
			return conn, nil
		}
		conn.Close()
	}

	conn, err := newConnection(uuid)
	if err != nil {
		return nil, err
	}
	m.Store(uuid, conn)

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

	return sc.Connect()
}

func loadCACert(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
