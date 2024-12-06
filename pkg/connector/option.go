package connector

import "github.com/openhdc/openhdc/pkg/transport"

type Option func(*Connector)

func WithKind(kind Kind) Option {
	return func(c *Connector) {
		c.kind = kind
	}
}

func WithSource(source Source) Option {
	return func(c *Connector) {
		c.source = source
	}
}

func WithDestination(destination Destination) Option {
	return func(c *Connector) {
		c.destination = destination
	}
}

func WithServerOptions(serverOpts []transport.ServerOption) Option {
	return func(c *Connector) {
		c.serverOpts = serverOpts
	}
}
