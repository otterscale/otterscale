package juju

import (
	"github.com/juju/juju/api"
	"github.com/juju/juju/api/connector"

	"github.com/openhdc/otterscale/internal/config"
)

type Juju struct {
	conf *config.Config
}

func New(conf *config.Config) *Juju {
	return &Juju{
		conf: conf,
	}
}

func (m *Juju) connection(uuid string) (api.Connection, error) {
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
