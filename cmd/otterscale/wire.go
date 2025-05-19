//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/otterscale/internal/app"
	"github.com/openhdc/otterscale/internal/cmd"
	"github.com/openhdc/otterscale/internal/data"
	"github.com/openhdc/otterscale/internal/domain"
	"github.com/spf13/cobra"
)

func wireApp(string) (*cobra.Command, func(), error) {
	panic(wire.Build(cmd.ProviderSet, app.ProviderSet, domain.ProviderSet, data.ProviderSet))
}
