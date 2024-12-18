//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/connector/csv/client"
)

func wireApp([]openhdc.ServerOption, []client.Option) (*app.App, func(), error) {
	panic(wire.Build(newApp, ProviderSet))
}
