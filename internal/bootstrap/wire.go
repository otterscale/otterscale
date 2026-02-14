package bootstrap

import "github.com/google/wire"

// ProviderSet is the Wire provider set for the bootstrap package.
var ProviderSet = wire.NewSet(New)
