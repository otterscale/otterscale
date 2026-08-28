// Package cmd defines the Cobra subcommands (server, agent) and their
// Wire provider sets. It bridges configuration, dependency injection,
// and the transport/application layers.
package cmd

import (
	"github.com/google/wire"

	"github.com/otterscale/otterscale/internal/cmd/agent"
	"github.com/otterscale/otterscale/internal/cmd/server"
)

// ProviderSet exposes the Agent and Server constructors and their handlers.
var ProviderSet = wire.NewSet(
	agent.NewAgent,
	agent.NewHandler,
	agent.ProvideEnrolmentToken,
	server.NewServer,
	server.NewHandler,
	server.ProvideBackgroundListeners,
	server.ProvideEnrolment,
)
