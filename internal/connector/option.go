package connector

type Option func(*Connector)

func WithName(name string) Option {
	return func(c *Connector) {
		c.name = name
	}
}

func WithVersion(version string) Option {
	return func(c *Connector) {
		c.version = version
	}
}

func WithPath(path string) Option {
	return func(c *Connector) {
		c.path = path
	}
}
