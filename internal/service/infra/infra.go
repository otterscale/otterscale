package infra

import (
	"github.com/google/wire"

	"github.com/openhdc/openhdc/internal/service/infra/juju"
	"github.com/openhdc/openhdc/internal/service/infra/kube"
	"github.com/openhdc/openhdc/internal/service/infra/maas"
)

var ProviderSet = wire.NewSet(
	maas.NewConfig,
	maas.New,
	juju.NewConfig,
	juju.New,
	kube.NewKubes,
	kube.NewNamespace,
	kube.NewCronJob,
	kube.NewJob,
)
