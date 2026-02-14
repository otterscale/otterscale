package config

import "github.com/google/wire"

// ProviderSet is the Wire provider set for configuration.
var ProviderSet = wire.NewSet(New)
