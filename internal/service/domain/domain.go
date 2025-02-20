package domain

import (
	"github.com/google/wire"

	"github.com/openhdc/openhdc/internal/service/domain/service"
)

var ProviderSet = wire.NewSet(service.NewUserService)
