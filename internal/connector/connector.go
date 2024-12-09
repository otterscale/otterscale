package connector

import (
	"context"
	"errors"
	"io"

	"github.com/apache/arrow-go/v18/arrow"
	"golang.org/x/sync/errgroup"
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
	msgs := make(chan arrow.Record)
	defer close(msgs)
	eg, ctx := errgroup.WithContext(stream.Context())
	eg.Go(func() error {
		return c.opts.source.Read(ctx, msgs, WithTables(req.Tables), WithSkipTables(req.SkipTables))
	})
	eg.Go(func() error {
		for msg := range msgs {
			rec, err := FromRecord(msg)
			if err != nil {
				return err
			}
			res := &pb.PullResponse{
				Record: rec,
			}
			// TODO: CHECK SIZE
			// if proto.Size(res) > MaxMsgSize {
			// 	continue
			// }
			if err := stream.Send(res); err != nil {
				return err
			}
		}
		return nil
	})
	return eg.Wait()
}

func (c *Connector) Push(stream pb.Connector_PushServer) error {
	msgs := make(chan arrow.Record)
	defer close(msgs)
	eg, ctx := errgroup.WithContext(stream.Context())
	eg.Go(func() error {
		return c.opts.destination.Write(ctx, msgs)
	})
	eg.Go(func() error {
		for {
			req, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
			rec, err := ToRecord(req.GetRecord())
			if err != nil {
				return err
			}
			msgs <- rec
		}
		return nil
	})
	return eg.Wait()
}

type Source interface {
	Read(ctx context.Context, record chan<- arrow.Record, opts ...ReadOption) error
	Close(ctx context.Context) error
}

type Destination interface {
	Write(ctx context.Context, record chan<- arrow.Record) error
	Close(ctx context.Context) error
}
