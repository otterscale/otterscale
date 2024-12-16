package openhdc

import (
	"context"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"

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
	kind     workload.Kind
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

func WithKind(kind workload.Kind) Option {
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
	Tables     []string
	SkipTables []string
}

func WithTables(tables []string) ReadOption {
	return func(o *ReadOptions) {
		o.Tables = tables
	}
}

func WithSkipTables(skipTables []string) ReadOption {
	return func(o *ReadOptions) {
		o.SkipTables = skipTables
	}
}

type WriteOption func(*WriteOptions)

type WriteOptions struct {
	// Kind WriteKind
	// Table string
}

// func WithWriteKind(kind WriteKind) WriteOption {
// 	return func(o *WriteOptions) {
// 		o.Kind = kind
// 	}
// }

// func WithTable(table string) WriteOption {
// 	return func(o *WriteOptions) {
// 		o.Table = table
// 	}
// }
