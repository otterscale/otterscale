package app

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewApplicationService,
	NewBootstrapService,
	NewConfigurationService,
	NewEnvironmentService,
	NewFacilityService,
	NewInstanceService,
	NewMachineService,
	NewModelService,
	NewNetworkService,
	NewOrchestratorService,
	NewRegistryService,
	NewResourceService,
	NewScopeService,
	NewStorageService,
)
