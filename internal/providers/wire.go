// Package providers aggregates all infrastructure-layer implementations
// (chisel, kubernetes, otterscale, cache) into a single Wire provider set.
package providers

import (
	"github.com/google/wire"

	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/providers/cache"
	"github.com/otterscale/otterscale/internal/providers/chisel"
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
	"github.com/otterscale/otterscale/internal/providers/manifest"
	"github.com/otterscale/otterscale/internal/providers/otterscale"
	"github.com/otterscale/otterscale/internal/transport"
)

// ProvideDiscoveryCache constructs a DiscoveryCache with the default TTL.
// This bridges the core.DiscoveryClient to the core.SchemaResolver
// interface via caching.
func ProvideDiscoveryCache(discovery core.DiscoveryClient) *cache.DiscoveryCache {
	return cache.NewDiscoveryCache(discovery, cache.DefaultTTL)
}

// ProviderSet is the Wire provider set for all external adapters.
var ProviderSet = wire.NewSet(
	chisel.NewService,
	wire.Bind(new(core.TunnelProvider), new(*chisel.Service)),
	wire.Bind(new(transport.TunnelService), new(*chisel.Service)),
	manifest.NewRenderer,
	wire.Bind(new(core.ManifestRenderer), new(*manifest.Renderer)),
	kubernetes.New,
	kubernetes.NewDiscoveryClient,
	kubernetes.NewResourceRepo,
	kubernetes.NewRuntimeRepo,
	otterscale.NewFleetRegistrar,
	ProvideDiscoveryCache,
	wire.Bind(new(core.SchemaResolver), new(*cache.DiscoveryCache)),
	wire.Bind(new(core.CacheEvictor), new(*cache.DiscoveryCache)),
)
