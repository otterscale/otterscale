package bootstrap

import (
	"github.com/google/wire"

	"github.com/otterscale/otterscale/internal/core/bootstrap/status"
)

var ProviderSet = wire.NewSet(
	status.NewStatusUseCase,
)
