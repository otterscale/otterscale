package client

type Option func(*Client)

func WithName(name string) Option {
	return func(c *Client) {
		c.name = name
	}
}

func WithVersion(version string) Option {
	return func(c *Client) {
		c.version = version
	}
}

func WithPath(path string) Option {
	return func(c *Client) {
		c.path = path
	}
}
