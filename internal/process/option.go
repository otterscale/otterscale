package process

type Option func(*options)

type options struct {
	name    string
	version string
	path    string
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

func WithPath(path string) Option {
	return func(o *options) {
		o.path = path
	}
}
