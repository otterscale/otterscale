package mux

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewBootstrap, NewServe)
