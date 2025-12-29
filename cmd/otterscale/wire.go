//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/mux"
	"github.com/otterscale/otterscale/internal/providers"
	"github.com/spf13/cobra"
)

func wireCmd(bool) (*cobra.Command, func(), error) {
	panic(wire.Build(newCmd, mux.ProviderSet, core.ProviderSet, providers.ProviderSet, config.ProviderSet))
}
