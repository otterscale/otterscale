package app

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewApplicationService,
	NewConfigurationService,
	NewEnvironmentService,
	NewFacilityService,
	NewEssentialService,
	NewMachineService,
	NewNetworkService,
	NewScopeService,
	NewStorageService,
	NewTagService,
)
