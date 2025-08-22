package data

import (
	"github.com/google/wire"

	"github.com/openhdc/otterscale/internal/data/ceph"
	"github.com/openhdc/otterscale/internal/data/juju"
	"github.com/openhdc/otterscale/internal/data/kube"
	"github.com/openhdc/otterscale/internal/data/maas"
)

var ProviderSet = wire.NewSet(
	maas.New,
	maas.NewBootResource,
	maas.NewBootSource,
	maas.NewBootSourceSelection,
	maas.NewEvent,
	maas.NewFabric,
	maas.NewIPRange,
	maas.NewMachine,
	maas.NewPackageRepository,
	maas.NewServer,
	maas.NewSSHKey,
	maas.NewSubnet,
	maas.NewTag,
	maas.NewVLAN,
	juju.New,
	juju.NewAction,
	juju.NewApplication,
	juju.NewApplicationOffers,
	juju.NewCharm,
	juju.NewClient,
	juju.NewKey,
	juju.NewMachine,
	juju.NewModel,
	juju.NewModelConfig,
	kube.New,
	kube.NewApps,
	kube.NewBatch,
	kube.NewCore,
	kube.NewHelmRelease,
	kube.NewHelmChart,
	kube.NewStorage,
	ceph.New,
	ceph.NewCluster,
	ceph.NewRBD,
	ceph.NewRGW,
	ceph.NewFS,
)
