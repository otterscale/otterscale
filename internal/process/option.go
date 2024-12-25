package process

import (
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/api/workload/v1"
)

type Option func(*options)

type options struct {
	kind    property.WorkloadKind
	name    string
	version string
	path    string
	sync    *workload.Sync
	spec    *structpb.Struct
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

func WithSync(sync *workload.Sync) Option {
	return func(o *options) {
		o.sync = sync
	}
}

func WithSpec(spec *structpb.Struct) Option {
	return func(o *options) {
		o.spec = spec
	}
}
