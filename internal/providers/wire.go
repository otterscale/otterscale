// Package providers aggregates all infrastructure-layer implementations
// (chisel, kubernetes, otterscale, cache) into a single Wire provider set.
package providers

import (
	"log/slog"

	"github.com/google/wire"
	"k8s.io/client-go/rest"

	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/providers/cache"
	"github.com/otterscale/otterscale/internal/providers/chisel"
	"github.com/otterscale/otterscale/internal/providers/harbor"
	"github.com/otterscale/otterscale/internal/providers/helm"
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
	"github.com/otterscale/otterscale/internal/providers/manifest"
	"github.com/otterscale/otterscale/internal/providers/otterscale"
	"github.com/otterscale/otterscale/internal/providers/rancher"
	"github.com/otterscale/otterscale/internal/transport"
)

// ProvideDiscoveryCache constructs a DiscoveryCache with the default TTL.
// This bridges the core.DiscoveryClient to the core.SchemaResolver
// interface via caching.
func ProvideDiscoveryCache(discovery core.DiscoveryClient) *cache.DiscoveryCache {
	return cache.NewDiscoveryCache(discovery, cache.DefaultTTL)
}

// ProvideComposingSchemaResolver wraps the DiscoveryCache with a
// ComposingSchemaResolver that applies schema composition rules
// (e.g. injecting FluxCD HelmRelease schema into ModuleTemplate).
// Accepting the concrete *cache.DiscoveryCache lets Wire distinguish
// the upstream cache from the composed resolver.
func ProvideComposingSchemaResolver(dc *cache.DiscoveryCache) *core.ComposingSchemaResolver {
	return core.NewComposingSchemaResolver(dc)
}

// ProvideRancherStore keeps Rancher Project discovery optional at server
// startup. Configuration or client construction failures leave the Project
// cache unavailable without blocking unrelated HTTP and tunnel services.
func ProvideRancherStore() *rancher.Store {
	return provideRancherStore(kubernetes.ProvideInClusterConfig)
}

func provideRancherStore(provideConfig func() (*rest.Config, error)) *rancher.Store {
	config, err := provideConfig()
	if err != nil {
		slog.Error("Rancher Project informer disabled: Kubernetes config unavailable", "error", err)
		return rancher.NewUnavailableStore()
	}
	store, err := rancher.NewStore(config)
	if err != nil {
		slog.Error("Rancher Project informer disabled: client initialization failed", "error", err)
		return rancher.NewUnavailableStore()
	}
	return store
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
	otterscale.NewLinkRegistrar,
	ProvideRancherStore,
	wire.Bind(new(core.RancherProjectStore), new(*rancher.Store)),
	harbor.ProvideHarborClient,
	helm.NewRepo,
	ProvideDiscoveryCache,
	ProvideComposingSchemaResolver,
	wire.Bind(new(core.SchemaResolver), new(*core.ComposingSchemaResolver)),
	wire.Bind(new(core.CacheEvictor), new(*cache.DiscoveryCache)),
)
