package process

import (
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/openhdc/openhdc/api/property/v1"
)

type Option func(*options)

type options struct {
	kind     property.WorkloadKind
	name     string
	version  string
	path     string
	syncMode property.SyncMode
	cursor   string
	spec     *structpb.Struct
}

func WithKind(kind property.WorkloadKind) Option {
	return func(o *options) {
		o.kind = kind
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

func WithCursor(cursor string) Option {
	return func(o *options) {
		o.cursor = cursor
	}
}

func WithSpec(spec *structpb.Struct) Option {
	return func(o *options) {
		o.spec = spec
	}
}
