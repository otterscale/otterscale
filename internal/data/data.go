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
	maas.NewBootResource,
	maas.NewFabric,
	maas.NewIPRange,
	maas.NewMachine,
	maas.NewPackageRepository,
	maas.NewServer,
	maas.NewSubnet,
	maas.NewVLAN,
	juju.NewConfig,
	juju.New,
	juju.NewModel,
	juju.NewModelConfig,
	kube.NewKubes,
	kube.NewNamespace,
	kube.NewCronJob,
	kube.NewJob,
	repo.NewConfig,
	repo.New,
)
