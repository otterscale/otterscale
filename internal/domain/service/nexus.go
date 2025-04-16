package service

type NexusService struct {
	server            MAASServer
	packageRepository MAASPackageRepository
	bootResource      MAASBootResource
	scope             JujuModel
	scopeConfig       JujuModelConfig
}

func NewNexusService(
	server MAASServer,
	packageRepository MAASPackageRepository,
	bootResource MAASBootResource,
	scope JujuModel,
	scopeConfig JujuModelConfig,
) *NexusService {
	return &NexusService{
		server:            server,
		packageRepository: packageRepository,
		bootResource:      bootResource,
		scope:             scope,
		scopeConfig:       scopeConfig,
	}
}
