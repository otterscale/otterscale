//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/openhdc/internal/service"
	"github.com/openhdc/openhdc/internal/service/app"
	"github.com/openhdc/openhdc/internal/service/domain"
	"github.com/openhdc/openhdc/internal/service/infra"
	"github.com/pocketbase/pocketbase"
)

func wireApp() (*pocketbase.PocketBase, func(), error) {
	panic(wire.Build(newApp, service.ProviderSet, app.ProviderSet, domain.ProviderSet, infra.ProviderSet))
}
