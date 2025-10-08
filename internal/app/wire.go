package app

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewApplicationService,
	NewConfigurationService,
	NewEnvironmentService,
	NewFacilityService,
	NewLargeLanguageModelService,
	NewMachineService,
	NewNetworkService,
	NewOrchestratorService,
	NewScopeService,
	NewStorageService,
	NewVirtualMachineService,
)
