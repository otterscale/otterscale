package domain

import (
	"github.com/google/wire"

	"github.com/openhdc/openhdc/internal/domain/service"
)

var ProviderSet = wire.NewSet(
	service.NewKubeService,
	service.NewStackService,
)
