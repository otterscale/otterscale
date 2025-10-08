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
	NewEssentialService,
	NewMachineService,
	NewNetworkService,
	NewPremiumService,
	NewScopeService,
	NewStorageService,
	NewTagService,
	NewVirtualMachineService,
)
