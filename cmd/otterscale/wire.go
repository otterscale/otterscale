//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/otterscale/otterscale/internal/bootstrap"
	"github.com/otterscale/otterscale/internal/cmd"
	"github.com/otterscale/otterscale/internal/cmd/agent"
	"github.com/otterscale/otterscale/internal/cmd/server"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/handler"
	"github.com/otterscale/otterscale/internal/providers"
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
	"github.com/otterscale/otterscale/internal/providers/manifest"
	"github.com/spf13/cobra"
)

// wireCmd assembles the root Cobra command with configuration loaded.
func wireCmd() (*cobra.Command, func(), error) {
	panic(wire.Build(newCmd, config.ProviderSet))
}

// wireServer assembles a fully wired Server with all gRPC services,
// use-cases, and infrastructure providers. The version parameter is
// provided by the caller and flows through Wire to LinkUseCase.
// The config parameter provides the CA directory for persistent CA
// material via provideCA.
func wireServer(v core.Version, conf *config.Config) (*server.Server, func(), error) {
	panic(wire.Build(cmd.ProviderSet, handler.ProviderSet, core.ProviderSet, providers.ProviderSet, provideCA, manifest.ProvideAgentManifestConfig))
}

// wireAgent assembles a fully wired Agent with its handler, link
// registrar, and bootstrapper. The version parameter is provided by
// the caller and flows through Wire to both LinkRegistrar and Agent.
func wireAgent(v core.Version) (*agent.Agent, func(), error) {
	panic(wire.Build(cmd.ProviderSet, providers.ProviderSet, bootstrap.ProviderSet, kubernetes.ProvideInClusterConfig))
}
