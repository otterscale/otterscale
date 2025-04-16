package service

type NexusService struct {
	server            MAASServer
	packageRepository MAASPackageRepository
	bootResource      MAASBootResource
	machine           MAASMachine
	client            JujuClient
	machineManager    JujuMachine
	scope             JujuModel
	scopeConfig       JujuModelConfig
}

func NewNexusService(
	server MAASServer,
	packageRepository MAASPackageRepository,
	bootResource MAASBootResource,
	machine MAASMachine,
	client JujuClient,
	machineManager JujuMachine,
	scope JujuModel,
	scopeConfig JujuModelConfig,
) *NexusService {
	return &NexusService{
		server:            server,
		packageRepository: packageRepository,
		bootResource:      bootResource,
		machine:           machine,
		client:            client,
		machineManager:    machineManager,
		scope:             scope,
		scopeConfig:       scopeConfig,
	}
}
