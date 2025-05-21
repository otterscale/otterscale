package maas

import (
	"github.com/canonical/gomaasclient/client"

	"github.com/openhdc/otterscale/internal/config"
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
	maas := m.conf.GetMaas()
	return client.GetClient(maas.GetUrl(), maas.GetKey(), maas.GetVersion())
}
