//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/internal/biz"
	"github.com/openhdc/openhdc/internal/data"
	"github.com/openhdc/openhdc/internal/server"
	"github.com/openhdc/openhdc/internal/service"
)

func wireApp([]openhdc.ServerOption) (*openhdc.App, func(), error) {
	panic(wire.Build(newApp, server.ProviderSet, service.ProviderSet, biz.ProviderSet, data.ProviderSet))
}
