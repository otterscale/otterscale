package core

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewBootstrapUseCase,
	NewBISTUseCase,
	NewChartUseCase,
	NewConfigurationUseCase,
	NewEnvironmentUseCase,
	NewFacilityUseCase,
	NewInstanceUseCase,
	NewKubernetesUseCase,
	NewMachineUseCase,
	NewModelUseCase,
	NewNetworkUseCase,
	NewOrchestratorUseCase,
	NewReleaseUseCase,
	NewScopeUseCase,
	NewStorageUseCase,
)
