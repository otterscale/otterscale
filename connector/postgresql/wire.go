//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/openhdc/connector/postgresql/client"
	"github.com/openhdc/openhdc/pkg/app"
	"github.com/openhdc/openhdc/pkg/connector"
	"github.com/openhdc/openhdc/pkg/transport"
)

func wireApp([]client.Option, []transport.ServerOption, []connector.Option) (*app.App, func(), error) {
	panic(wire.Build(newApp, ProviderSet))
}
