package app

import (
	"context"
	"os"
	"time"

	"github.com/openhdc/openhdc/pkg/transport"
)

type Option func(*options)

type options struct {
	id      string
	name    string
	version string
	ctx     context.Context
	sigs    []os.Signal
	timeout time.Duration
	servers []*transport.Server
}

func ID(id string) Option {
	return func(o *options) { o.id = id }
}

func Name(name string) Option {
	return func(o *options) { o.name = name }
}

func Version(version string) Option {
	return func(o *options) { o.version = version }
}

func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

func Signals(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

func Timeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

func Servers(servers ...*transport.Server) Option {
	return func(o *options) { o.servers = servers }
}
