//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/connectors/csv/client"
)

func wireApp([]openhdc.ServerOption, []client.Option) (*openhdc.App, func(), error) {
	panic(wire.Build(newApp, ProviderSet))
}
