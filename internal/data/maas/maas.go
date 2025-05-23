package maas

import (
	"log"

	"github.com/canonical/gomaasclient/client"

	"github.com/openhdc/otterscale/internal/config"
)

type MAAS struct {
	conf *config.ConfigSet
}

func New(conf *config.ConfigSet) *MAAS {
	return &MAAS{
		conf: conf,
	}
}

func (m *MAAS) client() (*client.Client, error) {
	maas := m.conf.GetMaas()
	log.Println(maas.GetUrl(), maas.GetKey(), maas.GetVersion())
	return client.GetClient(maas.GetUrl(), maas.GetKey(), maas.GetVersion())
}
