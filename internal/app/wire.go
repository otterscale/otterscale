package app

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewApplicationService,
	NewConfigurationService,
	NewFacilityService,
	NewEssentialService,
	NewMachineService,
	NewNetworkService,
	NewScopeService,
	NewTagService,
)
