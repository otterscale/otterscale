package maas

import (
	"github.com/canonical/gomaasclient/client"

	"github.com/openhdc/otterscale/internal/config"
)

type MAAS struct {
	configset *config.ConfigSet
}

func New(conf *config.Config) *MAAS {
	return &MAAS{
		configset: conf.ConfigSet,
	}
}

func (m *MAAS) client() (*client.Client, error) {
	maas := m.configset.GetMaas()
	return client.GetClient(maas.GetUrl(), maas.GetKey(), maas.GetVersion())
}
