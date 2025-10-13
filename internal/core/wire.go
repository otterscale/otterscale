package core

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewApplicationUseCase,
	NewBootstrapUseCase,
	NewBISTUseCase,
	NewConfigurationUseCase,
	NewEnvironmentUseCase,
	NewFacilityUseCase,
	NewLargeLanguageModelUseCase,
	NewMachineUseCase,
	NewNetworkUseCase,
	NewOrchestratorUseCase,
	NewScopeUseCase,
	NewStorageUseCase,
	NewVirtualMachineUseCase,
)
