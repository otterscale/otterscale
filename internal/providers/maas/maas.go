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

func (m *MAAS) Client() (*client.Client, error) {
	return client.GetClient(
		m.conf.MAASURL(),
		m.conf.MAASKey(),
		m.conf.MAASVersion(),
	)
}
