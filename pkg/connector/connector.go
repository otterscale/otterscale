package connector

import (
	"context"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"github.com/openhdc/openhdc/pkg/adapter"
	"github.com/openhdc/openhdc/pkg/codec"
	"github.com/openhdc/openhdc/pkg/transport"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pb.ConnectorServer = (*Connector)(nil)

type Connector struct {
	pb.UnimplementedConnectorServer

	codec       codec.Codec
	adapter     adapter.Adapter
	kind        Kind
	source      Source
	destination Destination
	serverOpts  []transport.ServerOption

	Server *transport.Server
}

func New(codec codec.Codec, adapter adapter.Adapter, opts ...Option) *Connector {
	c := &Connector{
		codec:   codec,
		adapter: adapter,
	}
	for _, opt := range opts {
		opt(c)
	}

	c.Server = transport.NewServer(append(c.serverOpts, transport.Connector())...)
	pb.RegisterConnectorServer(c.Server, c)

	return c
}

func (c *Connector) Close(ctx context.Context, _ *pb.CloseRequest) (*emptypb.Empty, error) {
	if c.source != nil {
		return &emptypb.Empty{}, c.source.Close(ctx)
	}
	if c.destination != nil {
		return &emptypb.Empty{}, c.destination.Close(ctx)
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
