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
	juju.NewAction,
	juju.NewApplication,
	juju.NewClient,
	juju.NewMachine,
	juju.NewModel,
	juju.NewModelConfig,
	kube.NewKubes,
	kube.NewClient,
	kube.NewCronJob,
	kube.NewDeployment,
	kube.NewJob,
	kube.NewPersistentVolumeClaim,
	kube.NewPod,
	kube.NewService,
	kube.NewNamespace,
	repo.NewConfig,
	repo.New,
)
