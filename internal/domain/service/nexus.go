package service

type NexusService struct {
	server            MAASServer
	packageRepository MAASPackageRepository
	bootResource      MAASBootResource
	fabric            MAASFabric
	vlan              MAASVLAN
	subnet            MAASSubnet
	ipRange           MAASIPRange
	machine           MAASMachine
	client            JujuClient
	machineManager    JujuMachine
	scope             JujuModel
	scopeConfig       JujuModelConfig
	facility          JujuApplication
	action            JujuAction
}

func NewNexusService(server MAASServer, packageRepository MAASPackageRepository, bootResource MAASBootResource, fabric MAASFabric, vlan MAASVLAN, subnet MAASSubnet, ipRange MAASIPRange, machine MAASMachine, client JujuClient, machineManager JujuMachine, scope JujuModel, scopeConfig JujuModelConfig, facility JujuApplication, action JujuAction) *NexusService {
	return &NexusService{
		server:            server,
		packageRepository: packageRepository,
		bootResource:      bootResource,
		fabric:            fabric,
		vlan:              vlan,
		subnet:            subnet,
		ipRange:           ipRange,
		machine:           machine,
		client:            client,
		machineManager:    machineManager,
		scope:             scope,
		scopeConfig:       scopeConfig,
		facility:          facility,
		action:            action,
	}
}
