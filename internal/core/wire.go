package core

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewLinkUseCase,
	NewAgentValuesUseCase,
	NewResourceUseCase,
	NewRuntimeUseCase,
	NewSessionStore,
)
