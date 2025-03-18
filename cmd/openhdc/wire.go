//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/internal/app"
	"github.com/openhdc/openhdc/internal/cmd"
	"github.com/openhdc/openhdc/internal/data"
	"github.com/openhdc/openhdc/internal/domain"
	"github.com/spf13/cobra"
)

func wireApp(string, []openhdc.ServerOption) (*cobra.Command, func(), error) {
	panic(wire.Build(cmd.ProviderSet, app.ProviderSet, domain.ProviderSet, data.ProviderSet))
}
