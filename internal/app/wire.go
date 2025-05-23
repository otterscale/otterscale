package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewApplicationService,
	NewConfigurationService,
	NewFacilityService,
	NewGeneralService,
	NewMachineService,
	NewNetworkService,
	NewScopeService,
	NewTagService,
)
