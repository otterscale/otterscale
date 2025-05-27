//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/otterscale/internal/app"
	"github.com/openhdc/otterscale/internal/config"
	"github.com/openhdc/otterscale/internal/core"
	"github.com/openhdc/otterscale/internal/data"
	"github.com/openhdc/otterscale/internal/mux"
	"github.com/spf13/cobra"
)

func wireCmd() (*cobra.Command, func(), error) {
	panic(wire.Build(newCmd, mux.ProviderSet, app.ProviderSet, core.ProviderSet, data.ProviderSet, config.ProviderSet))
}
