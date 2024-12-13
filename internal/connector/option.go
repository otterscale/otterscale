package connector

import (
	"github.com/openhdc/openhdc/api/workload/v1"
	"google.golang.org/grpc"
)

type Option func(*options)

type options struct {
	kind workload.Kind
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
