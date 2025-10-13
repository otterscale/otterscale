package core

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewBootstrapUseCase,
	NewBISTUseCase,
	NewChartUseCase,
	NewConfigurationUseCase,
	NewEnvironmentUseCase,
	NewFacilityUseCase,
	NewKubernetesUseCase,
	NewMachineUseCase,
	NewModelUseCase,
	NewNetworkUseCase,
	NewOrchestratorUseCase,
	NewReleaseUseCase,
	NewScopeUseCase,
	NewStorageUseCase,
	NewVirtualMachineUseCase,
)
