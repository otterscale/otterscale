package juju

import (
	"sync"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/connector"

	"github.com/openhdc/otterscale/internal/config"
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

func (m *Juju) connection(uuid string) (api.Connection, error) {
	if v, ok := m.connections.Load(uuid); ok {
		conn := v.(api.Connection)
		if !conn.IsBroken() {
			return conn, nil
		}
		conn.Close()
	}

	conn, err := m.newConnection(uuid)
	if err != nil {
		return nil, err
	}

	m.connections.Store(uuid, conn)

	return conn, nil
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
