//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/openhdc/openhdc/pkg/app"
)

func wireApp() (*app.App, func(), error) {
	panic(wire.Build(newApp, ProviderSet))
}
