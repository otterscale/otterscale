// Package providers aggregates all infrastructure-layer implementations
// (chisel, kubernetes, otterscale, cache) into a single Wire provider set.
package providers

import (
	"github.com/google/wire"

	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/providers/cache"
	"github.com/otterscale/otterscale/internal/providers/chisel"
	"github.com/otterscale/otterscale/internal/providers/harbor"
	"github.com/otterscale/otterscale/internal/providers/helm"
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
	"github.com/otterscale/otterscale/internal/providers/otterscale"
	"github.com/otterscale/otterscale/internal/providers/values"
	"github.com/otterscale/otterscale/internal/transport"
)

// ProvideDiscoveryCache bridges core.DiscoveryClient to core.SchemaResolver,
// caching at the default TTL.
func ProvideDiscoveryCache(discovery core.DiscoveryClient) *cache.DiscoveryCache {
	return cache.NewDiscoveryCache(discovery, cache.DefaultTTL)
}

var ProviderSet = wire.NewSet(
	chisel.NewService,
	wire.Bind(new(core.TunnelProvider), new(*chisel.Service)),
	wire.Bind(new(transport.TunnelService), new(*chisel.Service)),
	kubernetes.New,
	kubernetes.NewDiscoveryClient,
	kubernetes.NewResourceRepo,
	kubernetes.NewRuntimeRepo,
	otterscale.NewLinkRegistrar,
	helm.NewRepo,
	values.ProvideChartVersions,
	values.NewRenderer,
	wire.Bind(new(core.AgentValuesRenderer), new(*values.Renderer)),
	cache.NewAgentValuesStore,
	wire.Bind(new(core.AgentValuesStore), new(*cache.AgentValuesStore)),
	harbor.ProvideHarborClient,
	wire.Bind(new(core.HarborClient), new(*harbor.Client)),
	ProvideDiscoveryCache,
	wire.Bind(new(core.SchemaResolver), new(*cache.DiscoveryCache)),
	wire.Bind(new(core.CacheEvictor), new(*cache.DiscoveryCache)),
)
