package core

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewApplicationUseCase,
	NewConfigurationUseCase,
	NewFacilityUseCase,
	NewGeneralUseCase,
	NewMachineUseCase,
	NewNetworkUseCase,
	NewScopeUseCase,
	NewTagUseCase,
)
