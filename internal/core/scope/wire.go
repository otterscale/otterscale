package scope

import (
	"github.com/google/wire"

	"github.com/otterscale/otterscale/internal/core/scope/scope"
)

var ProviderSet = wire.NewSet(
	scope.NewScopeUseCase,
)
