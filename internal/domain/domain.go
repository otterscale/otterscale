package domain

import (
	"github.com/google/wire"

	"github.com/openhdc/otterscale/internal/domain/service"
)

var ProviderSet = wire.NewSet(service.NewNexusService)
