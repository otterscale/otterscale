package maas

import (
	"github.com/canonical/gomaasclient/client"

	"github.com/otterscale/otterscale/internal/config"
)

type MAAS struct {
	conf *config.Config
}

func New(conf *config.Config) *MAAS {
	return &MAAS{
		conf: conf,
	}
}

func (m *MAAS) client() (*client.Client, error) {
	maas := m.conf.MAAS
	return client.GetClient(maas.URL, maas.Key, maas.Version)
}
