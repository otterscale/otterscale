package connector

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/internal/adapter"
	"github.com/openhdc/openhdc/internal/codec"
	"github.com/openhdc/openhdc/internal/transport"
)

var _ pb.ConnectorServer = (*Connector)(nil)

type Connector struct {
	pb.UnimplementedConnectorServer

	opts    options
	codec   codec.Codec
	adapter adapter.Adapter
	server  *transport.Server
}

func New(codec codec.Codec, adapter adapter.Adapter, opts ...Option) *Connector {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	c := &Connector{
		opts:    o,
		codec:   codec,
		adapter: adapter,
	}
	c.server = transport.NewServer(append(o.serverOpts, transport.Connector())...)
	pb.RegisterConnectorServer(c.server, c)
	return c
}

func (c *Connector) Server() *transport.Server {
	return c.server
}

func (c *Connector) Close(ctx context.Context, _ *pb.CloseRequest) (*emptypb.Empty, error) {
	if c.opts.source != nil {
		return &emptypb.Empty{}, c.opts.source.Close(ctx)
	}
	if c.opts.destination != nil {
		return &emptypb.Empty{}, c.opts.destination.Close(ctx)
	}
	return &emptypb.Empty{}, nil
}

func (c *Connector) Pull(req *pb.PullRequest, stream pb.Connector_PullServer) error {
	// ctx := stream.Context()
	// records := make(chan arrow.Record)
	return nil
}

func (c *Connector) Push(stream pb.Connector_PushServer) error {
	// ctx := stream.Context()
	// return s.destination.Write(ctx)
	return nil
}
