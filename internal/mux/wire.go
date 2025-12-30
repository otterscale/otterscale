package mux

import (
	"github.com/google/wire"

	"github.com/otterscale/otterscale/internal/app"
)

var ProviderSet = wire.NewSet(NewBootstrap, NewServe, app.ProviderSet)
