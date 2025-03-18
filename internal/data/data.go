package data

import (
	"github.com/google/wire"

	"github.com/openhdc/openhdc/internal/data/juju"
	"github.com/openhdc/openhdc/internal/data/kube"
	"github.com/openhdc/openhdc/internal/data/maas"
	"github.com/openhdc/openhdc/internal/data/repo"
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
	repo.NewConfig,
	repo.New,
	repo.NewUser,
)
