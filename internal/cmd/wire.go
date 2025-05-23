package cmd

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(New)
