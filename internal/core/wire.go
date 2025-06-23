package core

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewApplicationUseCase,
	NewConfigurationUseCase,
	NewEnvironmentUseCase,
	NewFacilityUseCase,
	NewEssentialUseCase,
	NewMachineUseCase,
	NewNetworkUseCase,
	NewScopeUseCase,
	NewStorageUseCase,
	NewTagUseCase,
)
