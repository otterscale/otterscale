package juju

import (
	"sync"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/connector"

	"github.com/openhdc/otterscale/internal/config"
)

type Juju struct {
	conf        *config.ConfigSet
	connections sync.Map
}

func New(conf *config.ConfigSet) *Juju {
	return &Juju{
		conf: conf,
	}
}

func (m *Juju) newConnection(uuid string) (api.Connection, error) {
	juju := m.conf.GetJuju()
	opts := connector.SimpleConfig{
		ModelUUID:           uuid,
		ControllerAddresses: juju.GetControllerAddresses(),
		Username:            juju.GetUsername(),
		Password:            juju.GetPassword(),
		CACert:              juju.GetCaCert(),
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
	juju := m.conf.GetJuju()
	if juju != nil {
		return juju.GetUsername()
	}
	return ""
}

func (m *Juju) cloudName() string {
	juju := m.conf.GetJuju()
	if juju != nil {
		return juju.GetCloudName()
	}
	return ""
}

func (m *Juju) cloudRegion() string {
	juju := m.conf.GetJuju()
	if juju != nil {
		return juju.GetCloudRegion()
	}
	return ""
}

func (m *Juju) charmhubAPIURL() string {
	juju := m.conf.GetJuju()
	if juju != nil {
		return juju.GetCharmhubApiUrl()
	}
	return ""
}
