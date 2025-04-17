package data

import (
	"github.com/google/wire"

	"github.com/openhdc/openhdc/internal/data/juju"
	"github.com/openhdc/openhdc/internal/data/kube"
	"github.com/openhdc/openhdc/internal/data/maas"
	"github.com/openhdc/openhdc/internal/data/repo"
)

var ProviderSet = wire.NewSet(
	maas.New,
	maas.NewBootResource,
	maas.NewFabric,
	maas.NewIPRange,
	maas.NewMachine,
	maas.NewPackageRepository,
	maas.NewServer,
	maas.NewSubnet,
	maas.NewVLAN,
	juju.NewJujuMap,
	juju.NewAction,
	juju.NewApplication,
	juju.NewClient,
	juju.NewMachine,
	juju.NewModel,
	juju.NewModelConfig,
	kube.NewKubeMap,
	kube.NewHelmMap,
	kube.NewClient,
	kube.NewApps,
	kube.NewBatch,
	kube.NewCore,
	kube.NewStorage,
	kube.NewHelm,
	repo.NewConfig,
	repo.New,
)
