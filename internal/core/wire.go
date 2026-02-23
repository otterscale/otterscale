package core

import (
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for all domain use-cases.
var ProviderSet = wire.NewSet(
	NewLinkUseCase,
	NewResourceUseCase,
	NewRuntimeUseCase,
	NewSessionStore,
)
