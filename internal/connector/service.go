package connector

import (
	"context"
	"errors"
	"io"

	pb "github.com/openhdc/openhdc/api/connector/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pb.ConnectorServer = (*Service)(nil)

type Service struct {
	pb.UnimplementedConnectorServer

	opts      options
	connector Connector
}

func NewService(c Connector, opts ...Option) *Service {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &Service{
		opts:      o,
		connector: c,
	}
}

func (s *Service) Close(ctx context.Context, _ *pb.CloseRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.connector.Close(ctx)
}

func (s *Service) Pull(req *pb.PullRequest, stream pb.Connector_PullServer) error {
	msgs := make(chan *pb.Message)
	defer close(msgs)
	eg, ctx := errgroup.WithContext(stream.Context())
	eg.Go(func() error {
		// TODO: BETTER
		opts := ReadOptions{}
		for _, opt := range []ReadOption{
			WithNamespace(req.Namespace),
			WithTables(req.Tables),
			WithSkipTables(req.SkipTables),
		} {
			opt(&opts)
		}
		return s.connector.Read(ctx, msgs, opts)
	})
	eg.Go(func() error {
		for msg := range msgs {
			// TODO: CHECK SIZE
			// if proto.Size(res) > MaxMsgSize {
			// 	continue
			// }
			if err := stream.Send(msg); err != nil {
				return err
			}
		}
		return nil
	})
	return eg.Wait()
}

func (s *Service) Push(stream pb.Connector_PushServer) error {
	msgs := make(chan *pb.Message)
	defer close(msgs)
	eg, ctx := errgroup.WithContext(stream.Context())
	eg.Go(func() error {
		// TODO: BETTER
		opts := WriteOptions{}
		for _, opt := range []WriteOption{} {
			opt(&opts)
		}
		return s.connector.Write(ctx, msgs, opts)
	})
	eg.Go(func() error {
		for {
			msg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
			msgs <- msg
		}
		return nil
	})
	return eg.Wait()
}
