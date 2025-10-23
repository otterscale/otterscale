package juju

import (
	"fmt"
	"sync"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/client/modelmanager"
	"github.com/juju/juju/api/connector"

	"github.com/otterscale/otterscale/internal/config"
)

type Juju struct {
	conf        *config.Config
	connections sync.Map
}

func New(conf *config.Config) *Juju {
	return &Juju{
		conf: conf,
	}
}

func (m *Juju) newConnection(uuid string) (api.Connection, error) {
	juju := m.conf.Juju
	opts := connector.SimpleConfig{
		ModelUUID:           uuid,
		ControllerAddresses: juju.ControllerAddresses,
		Username:            juju.Username,
		Password:            juju.Password,
		CACert:              juju.CACert,
	}
	sc, err := connector.NewSimple(opts)
	if err != nil {
		return nil, err
	}
	return sc.Connect()
}

func (m *Juju) connection(scope string) (api.Connection, error) {
	if v, ok := m.connections.Load(scope); ok {
		conn := v.(api.Connection)
		if !conn.IsBroken() {
			return conn, nil
		}
		conn.Close()
	}

	uuid, err := m.getModelUUID(scope)
	if err != nil {
		return nil, err
	}

	conn, err := m.newConnection(uuid)
	if err != nil {
		return nil, err
	}

	m.connections.Store(scope, conn)

	return conn, nil
}

func (m *Juju) getModelUUID(scope string) (string, error) {
	if scope == "controller" {
		return "", nil
	}

	conn, err := m.newConnection("")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	models, err := modelmanager.NewClient(conn).ListModels(m.username())
	if err != nil {
		return "", err
	}
	for _, model := range models {
		if model.Name == scope {
			return model.UUID, nil
		}
	}
	return "", fmt.Errorf("scope %q not found", scope)
}

func (m *Juju) username() string {
	return m.conf.Juju.Username
}

func (m *Juju) cloudName() string {
	return m.conf.Juju.CloudName
}

func (m *Juju) cloudRegion() string {
	return m.conf.Juju.CloudRegion
}

func (m *Juju) charmhubAPIURL() string {
	return m.conf.Juju.CharmhubAPIURL
}
