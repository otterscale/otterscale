package openhdc

import (
	"context"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/openhdc/openhdc/api/property/v1"
	"github.com/openhdc/openhdc/api/workload/v1"
)

type Option func(*options)

type options struct {
	id       string
	name     string
	version  string
	ctx      context.Context
	sigs     []os.Signal
	timeout  time.Duration
	servers  []*Server
	logLevel slog.Leveler
	kind     property.WorkloadKind
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

func WithServers(servers ...*Server) Option {
	return func(o *options) {
		o.servers = servers
	}
}

func WithLogLevel(logLevel slog.Leveler) Option {
	return func(o *options) {
		o.logLevel = logLevel
	}
}

func WithKind(kind property.WorkloadKind) Option {
	return func(o *options) {
		o.kind = kind
	}
}

type ServerOption func(*serverOptions)

type serverOptions struct {
	network  string
	address  string
	grpcOpts []grpc.ServerOption
}

func WithNetwork(network string) ServerOption {
	return func(o *serverOptions) {
		o.network = network
	}
}

func WithAddress(address string) ServerOption {
	return func(o *serverOptions) {
		o.address = address
	}
}

func WithGrpcServerOptions(grpcOpts ...grpc.ServerOption) ServerOption {
	return func(o *serverOptions) {
		o.grpcOpts = append(o.grpcOpts, grpcOpts...)
	}
}

type ReadOption func(*ReadOptions)

type ReadOptions struct {
	BatchSize int64
	Keys      []string
	SkipKeys  []string
	Options   []*workload.Sync_Option
}

func WithBatchSize(batchSize int64) ReadOption {
	return func(o *ReadOptions) {
		o.BatchSize = batchSize
	}
}

func WithKeys(keys []string) ReadOption {
	return func(o *ReadOptions) {
		o.Keys = keys
	}
}

func WithSkipKeys(skipKeys []string) ReadOption {
	return func(o *ReadOptions) {
		o.SkipKeys = skipKeys
	}
}

func WithOptions(options []*workload.Sync_Option) ReadOption {
	return func(o *ReadOptions) {
		o.Options = options
	}
}
