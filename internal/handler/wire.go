package handler

import (
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for ConnectRPC service handlers
// and the raw HTTP manifest handler.
var ProviderSet = wire.NewSet(NewFleetService, NewResourceService, NewRuntimeService, NewManifestHandler)
