// Package cmd defines the Cobra subcommands (server, agent) and their
// Wire provider sets. It bridges configuration, dependency injection,
// and the transport/application layers.
package cmd

import (
	"github.com/google/wire"

	"github.com/otterscale/otterscale/internal/cmd/agent"
	"github.com/otterscale/otterscale/internal/cmd/server"
)

// ProviderSet is the Wire provider set for the CLI layer. It exposes
// the Agent and Server constructors plus their handlers.
var ProviderSet = wire.NewSet(
	agent.NewAgent,
	agent.NewHandler,
	agent.NewUpdater,
	server.NewServer,
	server.NewHandler,
	server.ProvideBackgroundListeners,
)
