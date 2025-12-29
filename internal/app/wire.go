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
	NewKubernetesService,
	NewMachineService,
	NewModelService,
	NewNetworkService,
	NewOrchestratorService,
	NewRegistryService,
	NewScopeService,
	NewStorageService,
)
