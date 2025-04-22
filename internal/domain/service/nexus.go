package service

type NexusService struct {
	server              MAASServer
	packageRepository   MAASPackageRepository
	bootResource        MAASBootResource
	bootSource          MAASBootSource
	bootSourceSelection MAASBootSourceSelection
	fabric              MAASFabric
	vlan                MAASVLAN
	subnet              MAASSubnet
	ipRange             MAASIPRange
	machine             MAASMachine
	tag                 MAASTag
	client              JujuClient
	machineManager      JujuMachine
	scope               JujuModel
	scopeConfig         JujuModelConfig
	facility            JujuApplication
	action              JujuAction
	charmhub            JujuCharmHub
	kubernetes          KubeClient
	apps                KubeApps
	batch               KubeBatch
	core                KubeCore
	storage             KubeStorage
	helm                KubeHelm
}

func NewNexusService(server MAASServer, packageRepository MAASPackageRepository, bootResource MAASBootResource, bootSource MAASBootSource, bootSourceSelection MAASBootSourceSelection, fabric MAASFabric, vlan MAASVLAN, subnet MAASSubnet, ipRange MAASIPRange, machine MAASMachine, tag MAASTag, client JujuClient, machineManager JujuMachine, scope JujuModel, scopeConfig JujuModelConfig, facility JujuApplication, action JujuAction, charmhub JujuCharmHub, kubernetes KubeClient, apps KubeApps, batch KubeBatch, core KubeCore, storage KubeStorage, helm KubeHelm) *NexusService {
	return &NexusService{
		server:              server,
		packageRepository:   packageRepository,
		bootResource:        bootResource,
		bootSource:          bootSource,
		bootSourceSelection: bootSourceSelection,
		fabric:              fabric,
		vlan:                vlan,
		subnet:              subnet,
		ipRange:             ipRange,
		machine:             machine,
		tag:                 tag,
		client:              client,
		machineManager:      machineManager,
		scope:               scope,
		scopeConfig:         scopeConfig,
		facility:            facility,
		action:              action,
		charmhub:            charmhub,
		kubernetes:          kubernetes,
		apps:                apps,
		batch:               batch,
		core:                core,
		storage:             storage,
		helm:                helm,
	}
}
