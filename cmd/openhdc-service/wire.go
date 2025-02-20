//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/internal/service/app"
	"github.com/openhdc/openhdc/internal/service/domain"
	"github.com/openhdc/openhdc/internal/service/infra"
	"github.com/openhdc/openhdc/internal/service/server"
)

func wireApp([]openhdc.ServerOption) (*openhdc.App, func(), error) {
	panic(wire.Build(newApp, server.ProviderSet, app.ProviderSet, domain.ProviderSet, infra.ProviderSet))
}
