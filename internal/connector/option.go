package connector

import "github.com/openhdc/openhdc/internal/transport"

type Option func(*options)

type options struct {
	kind        Kind
	source      Source
	destination Destination
	serverOpts  []transport.ServerOption
}

func WithKind(kind Kind) Option {
	return func(o *options) {
		o.kind = kind
	}
}

func WithSource(source Source) Option {
	return func(o *options) {
		o.source = source
	}
}

func WithDestination(destination Destination) Option {
	return func(o *options) {
		o.destination = destination
	}
}

func WithServerOptions(serverOpts []transport.ServerOption) Option {
	return func(o *options) {
		o.serverOpts = serverOpts
	}
}
