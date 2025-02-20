//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/openhdc"
	"github.com/openhdc/openhdc/internal/service"
	"github.com/openhdc/openhdc/internal/service/app"
	"github.com/openhdc/openhdc/internal/service/domain"
	"github.com/openhdc/openhdc/internal/service/infra"
)

func wireApp([]openhdc.ServerOption) (*openhdc.App, func(), error) {
	panic(wire.Build(newApp, service.ProviderSet, app.ProviderSet, domain.ProviderSet, infra.ProviderSet))
}
