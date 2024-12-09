//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/openhdc/connector/postgresql/client"
	"github.com/openhdc/openhdc/internal/app"
	"github.com/openhdc/openhdc/internal/connector"
)

func wireApp([]client.Option, []connector.Option, []connector.ServerOption) (*app.App, func(), error) {
	panic(wire.Build(newApp, ProviderSet))
}
