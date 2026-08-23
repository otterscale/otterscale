package handler

import (
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for ConnectRPC service handlers
// and the Prometheus reverse proxy.
var ProviderSet = wire.NewSet(NewLinkService, NewResourceService, NewRuntimeService, NewProxyHandler)
