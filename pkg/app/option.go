package app

import "github.com/openhdc/openhdc/pkg/transport"

type Option func(*App)

func Servers(srv ...transport.Server) Option {
	return func(o *App) { o.servers = srv }
}
