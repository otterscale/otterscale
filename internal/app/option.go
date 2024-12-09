package app

import (
	"context"
	"os"
	"time"

	"github.com/openhdc/openhdc/internal/connector"
)

type Option func(*options)

type options struct {
	id      string
	name    string
	version string
	ctx     context.Context
	sigs    []os.Signal
	timeout time.Duration
	servers []*connector.Server
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithVersion(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithSignals(sigs ...os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

func WithServers(servers ...*connector.Server) Option {
	return func(o *options) {
		o.servers = servers
	}
}
