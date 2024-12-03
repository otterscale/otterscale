package app

import "github.com/openhdc/openhdc/pkg/transport"

type App struct {
	servers []*transport.Server
}

func New(opts ...Option) *App {
	return nil
}

func (a *App) Run() error {
	return nil
}
