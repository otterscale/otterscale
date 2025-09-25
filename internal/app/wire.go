package app

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewApplicationService,
	NewBISTService,
	NewConfigurationService,
	NewEnvironmentService,
	NewFacilityService,
	NewEssentialService,
	NewMachineService,
	NewNetworkService,
	NewPremiumService,
	NewScopeService,
	NewStorageService,
	NewTagService,
	NewVirtualMachineService,
)
