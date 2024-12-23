package process

import (
	"github.com/openhdc/openhdc/api/property/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

type Option func(*options)

type options struct {
	name     string
	version  string
	path     string
	syncMode property.SyncMode
	spec     *structpb.Struct
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

func WithSyncMode(syncMode property.SyncMode) Option {
	return func(o *options) {
		o.syncMode = syncMode
	}
}

func WithSpec(spec *structpb.Struct) Option {
	return func(o *options) {
		o.spec = spec
	}
}
