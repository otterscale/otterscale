package maas

import (
	"github.com/canonical/gomaasclient/client"

	"github.com/openhdc/openhdc/internal/env"
)

type MAAS = client.Client

const (
	defaultURL     = "http://localhost:5240/MAAS/"
	defaultKey     = "http://localhost:5240/MAAS/"
	defaultVersion = "2.0"
)

func New() (*MAAS, error) {
	url := env.GetOrDefault(env.OPENHDC_MAAS_API_URL, defaultURL)
	key := env.GetOrDefault(env.OPENHDC_MAAS_API_KEY, "::")
	version := env.GetOrDefault(env.OPENHDC_MAAS_API_VERSION, defaultVersion)
	return client.GetClient(url, key, version)
}
